package pipelinerun

import (
	"context"
	"fmt"
	"strings"

	"g.hz.netease.com/horizon/core/common"
	appmanager "g.hz.netease.com/horizon/pkg/application/manager"
	"g.hz.netease.com/horizon/pkg/cluster/code"
	"g.hz.netease.com/horizon/pkg/cluster/gitrepo"
	clustermanager "g.hz.netease.com/horizon/pkg/cluster/manager"
	clustermodels "g.hz.netease.com/horizon/pkg/cluster/models"
	"g.hz.netease.com/horizon/pkg/cluster/tekton/factory"
	"g.hz.netease.com/horizon/pkg/cluster/tekton/log"
	envmanager "g.hz.netease.com/horizon/pkg/environment/manager"
	prmanager "g.hz.netease.com/horizon/pkg/pipelinerun/manager"
	prmodels "g.hz.netease.com/horizon/pkg/pipelinerun/models"
	"g.hz.netease.com/horizon/pkg/util/errors"
	"g.hz.netease.com/horizon/pkg/util/wlog"
)

type Controller interface {
	GetPrLog(ctx context.Context, prID uint) (*Log, error)
	GetClusterLatestLog(ctx context.Context, clusterID uint) (*Log, error)
	GetDiff(ctx context.Context, pipelinerunID uint) (*GetDiffResponse, error)
}

type controller struct {
	prMgr          prmanager.Manager
	clusterMgr     clustermanager.Manager
	applicationMgr appmanager.Manager
	envMgr         envmanager.Manager
	tektonFty      factory.Factory
	commitGetter   code.CommitGetter
	clusterGitRepo gitrepo.ClusterGitRepo
}

var _ Controller = (*controller)(nil)

func NewController(tektonFty factory.Factory) Controller {
	return &controller{
		prMgr:          prmanager.Mgr,
		clusterMgr:     clustermanager.Mgr,
		applicationMgr: appmanager.Mgr,
		envMgr:         envmanager.Mgr,
		tektonFty:      tektonFty,
	}
}

type Log struct {
	LogChannel <-chan log.Log
	ErrChannel <-chan error

	LogBytes []byte
}

func (c *controller) GetPrLog(ctx context.Context, prID uint) (_ *Log, err error) {
	const op = "pipelinerun controller: get pipelinerun log"
	defer wlog.Start(ctx, op).Stop(func() string { return wlog.ByErr(err) })

	pr, err := c.prMgr.GetByID(ctx, prID)
	if err != nil {
		return nil, errors.E(op, err)
	}

	cluster, err := c.clusterMgr.GetByID(ctx, pr.ClusterID)
	if err != nil {
		return nil, errors.E(op, err)
	}
	application, err := c.applicationMgr.GetByID(ctx, cluster.ApplicationID)
	if err != nil {
		return nil, errors.E(op, err)
	}
	er, err := c.envMgr.GetEnvironmentRegionByID(ctx, cluster.EnvironmentRegionID)
	if err != nil {
		return nil, errors.E(op, err)
	}

	// only builddeploy have logs
	if pr.Action != prmodels.ActionBuildDeploy {
		return nil, errors.E(op, fmt.Errorf("%v action has no log", pr.Action))
	}

	return c.getPrLog(ctx, pr, cluster, application.Name, er.EnvironmentName)
}

func (c *controller) GetClusterLatestLog(ctx context.Context, clusterID uint) (_ *Log, err error) {
	const op = "pipelinerun controller: get cluster latest log"
	defer wlog.Start(ctx, op).Stop(func() string { return wlog.ByErr(err) })

	pr, err := c.prMgr.GetLatestByClusterIDAndAction(ctx, clusterID, prmodels.ActionBuildDeploy)
	if err != nil {
		return nil, errors.E(op, err)
	}
	if pr == nil {
		return nil, errors.E(op, fmt.Errorf("no builddeploy pipelinerun"))
	}

	cluster, err := c.clusterMgr.GetByID(ctx, clusterID)
	if err != nil {
		return nil, errors.E(op, err)
	}
	application, err := c.applicationMgr.GetByID(ctx, cluster.ApplicationID)
	if err != nil {
		return nil, errors.E(op, err)
	}
	er, err := c.envMgr.GetEnvironmentRegionByID(ctx, cluster.EnvironmentRegionID)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return c.getPrLog(ctx, pr, cluster, application.Name, er.EnvironmentName)
}

func (c *controller) getPrLog(ctx context.Context, pr *prmodels.Pipelinerun, cluster *clustermodels.Cluster,
	application, environment string) (_ *Log, err error) {
	const op = "pipeline controller: get pipelinerun log"
	defer wlog.Start(ctx, op).Stop(func() string { return wlog.ByErr(err) })

	// if pr.PrObject is empty, get pipelinerun log in k8s
	if pr.PrObject == "" {
		tektonClient, err := c.tektonFty.GetTekton(environment)
		if err != nil {
			return nil, errors.E(op, err)
		}

		logCh, errCh, err := tektonClient.GetPipelineRunLogByID(ctx, cluster.Name, cluster.ID, pr.ID)
		if err != nil {
			return nil, errors.E(op, err)
		}
		return &Log{
			LogChannel: logCh,
			ErrChannel: errCh,
		}, nil
	}

	// else, get log from s3
	tektonCollector, err := c.tektonFty.GetTektonCollector(environment)
	if err != nil {
		return nil, errors.E(op, err)
	}
	// TODO(gjq): get pipelinerun log by pr.Object
	logBytes, err := tektonCollector.GetLatestPipelineRunLog(ctx, application, cluster.Name)
	if err != nil {
		return nil, errors.E(op, err)
	}
	return &Log{
		LogBytes: logBytes,
	}, nil
}

func (c *controller) GetDiff(ctx context.Context, pipelinerunID uint) (_ *GetDiffResponse, err error) {
	const op = "pipelinerun controller: get pipelinerun diff"
	defer wlog.Start(ctx, op).Stop(func() string { return wlog.ByErr(err) })

	// 1. get pipeline
	pipelinerun, err := c.prMgr.GetByID(ctx, pipelinerunID)
	if err != nil {
		return nil, err
	}

	// 2. get cluster and application
	cluster, err := c.clusterMgr.GetByID(ctx, pipelinerun.ClusterID)
	if err != nil {
		return nil, err
	}
	application, err := c.applicationMgr.GetByID(ctx, cluster.ApplicationID)
	if err != nil {
		return nil, err
	}

	// 3. get code diff
	var codeDiff *CodeInfo
	if pipelinerun.GitURL != "" && pipelinerun.GitCommit != "" &&
		pipelinerun.GitBranch != "" {
		commit, err := c.commitGetter.GetCommit(ctx, pipelinerun.GitURL, pipelinerun.GitCommit)
		if err != nil {
			return nil, err
		}
		var historyLink string
		if strings.HasPrefix(pipelinerun.GitURL, common.InternalGitSSHPrefix) {
			httpURL := common.InternalSSHToHTTPURL(pipelinerun.GitURL)
			historyLink = httpURL + common.CommitHistoryMiddle + pipelinerun.GitCommit
		}
		codeDiff = &CodeInfo{
			Branch:    pipelinerun.GitBranch,
			CommitID:  pipelinerun.GitCommit,
			CommitMsg: commit.Message,
			Link:      historyLink,
		}
	}

	// 4. get config diff
	var configDiff *ConfigDiff
	if pipelinerun.ConfigCommit != "" && pipelinerun.LastConfigCommit != "" {
		diff, err := c.clusterGitRepo.CompareConfig(ctx, application.Name, cluster.Name,
			&pipelinerun.LastConfigCommit, &pipelinerun.ConfigCommit)
		if err != nil {
			return nil, err
		}
		configDiff = &ConfigDiff{
			From: pipelinerun.LastConfigCommit,
			To:   pipelinerun.ConfigCommit,
			Diff: diff,
		}
	}

	return &GetDiffResponse{
		CodeInfo:   codeDiff,
		ConfigDiff: configDiff,
	}, nil
}
