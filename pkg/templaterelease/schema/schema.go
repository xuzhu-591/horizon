package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	gitlabfty "g.hz.netease.com/horizon/pkg/gitlab/factory"
	"g.hz.netease.com/horizon/pkg/templaterelease/manager"
	"g.hz.netease.com/horizon/pkg/util/errors"
	"g.hz.netease.com/horizon/pkg/util/wlog"
)

var (
	Gtr = newSchemaGetter()
)

// Getter provides some functions for template schema
type Getter interface {
	// GetTemplateSchema get schema for specified template release
	GetTemplateSchema(ctx context.Context, templateName, releaseName string) (*Schemas, error)
}

type Schemas struct {
	Application *Schema
	Pipeline    *Schema
}

type Schema struct {
	JSONSchema map[string]interface{}
	UISchema   map[string]interface{}
}

const (
	// json schema file path
	_pipelineSchemaPath    = "schema/pipeline.schema.json"
	_applicationSchemaPath = "schema/application.schema.json"
	// ui schema file path
	_pipelineUISchemaPath    = "schema/pipeline.ui.schema.json"
	_applicationUISchemaPath = "schema/application.ui.schema.json"

	// ErrCodeReleaseNotFound  ReleaseNotFound error code
	_errCodeReleaseNotFound = errors.ErrorCode("ReleaseNotFound")
)

type getter struct {
	gitlabFactory      gitlabfty.Factory
	templateReleaseMgr manager.Manager
}

func newSchemaGetter() Getter {
	return &getter{
		gitlabFactory:      gitlabfty.Fty,
		templateReleaseMgr: manager.Mgr,
	}
}

func (g *getter) GetTemplateSchema(ctx context.Context,
	templateName, releaseName string) (_ *Schemas, err error) {
	const op = "template service: getTemplateSchema"
	defer wlog.Start(ctx, op).Stop(func() string { return wlog.ByErr(err) })

	tr, err := g.templateReleaseMgr.GetByTemplateNameAndRelease(ctx, templateName, releaseName)
	if err != nil {
		return nil, err
	}
	if tr == nil {
		return nil, errors.E(op, http.StatusNotFound, _errCodeReleaseNotFound,
			fmt.Sprintf("the release %v of template %v is not found", releaseName, templateName))
	}
	gitlabLib, err := g.gitlabFactory.GetByName(ctx, tr.GitlabName)
	if err != nil {
		return nil, err
	}

	// 1. read file concurrently
	var err1, err2, err3, err4 error
	var pipelineSchemaBytes, applicationSchemaBytes, pipelineUISchemaBytes, applicationUISchemaBytes []byte
	var wgReadFile sync.WaitGroup
	wgReadFile.Add(4)
	readFile := func(b *[]byte, err *error, filePath string) {
		defer wgReadFile.Done()
		bytes, e := gitlabLib.GetFile(ctx, tr.GitlabProject, tr.Name, filePath)
		*b = bytes
		*err = e
	}
	go readFile(&pipelineSchemaBytes, &err1, _pipelineSchemaPath)
	go readFile(&applicationSchemaBytes, &err2, _applicationSchemaPath)
	go readFile(&pipelineUISchemaBytes, &err3, _pipelineUISchemaPath)
	go readFile(&applicationUISchemaBytes, &err4, _applicationUISchemaPath)
	wgReadFile.Wait()

	for _, err := range []error{err1, err2, err3, err4} {
		if err != nil {
			return nil, err
		}
	}

	// 2. unmarshal concurrently
	var pipelineSchema, applicationSchema, pipelineUISchema, applicationUISchema map[string]interface{}
	var wgUnmarshal sync.WaitGroup
	wgUnmarshal.Add(4)
	unmarshal := func(b []byte, m *map[string]interface{}, err *error) {
		defer wgUnmarshal.Done()
		if e := json.Unmarshal(b, &m); e != nil {
			*err = e
		}
	}
	go unmarshal(pipelineSchemaBytes, &pipelineSchema, &err1)
	go unmarshal(applicationSchemaBytes, &applicationSchema, &err2)
	go unmarshal(pipelineUISchemaBytes, &pipelineUISchema, &err3)
	go unmarshal(applicationUISchemaBytes, &applicationUISchema, &err4)
	wgUnmarshal.Wait()

	for _, err := range []error{err1, err2, err3, err4} {
		if err != nil {
			return nil, err
		}
	}

	return &Schemas{
		Application: &Schema{
			JSONSchema: applicationSchema,
			UISchema:   applicationUISchema,
		},
		Pipeline: &Schema{
			JSONSchema: pipelineSchema,
			UISchema:   pipelineUISchema,
		},
	}, nil
}
