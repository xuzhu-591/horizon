// Copyright © 2023 Horizoncd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"context"

	"github.com/horizoncd/horizon/lib/q"
	"github.com/horizoncd/horizon/pkg/pipelinerun/dao"
	"github.com/horizoncd/horizon/pkg/pipelinerun/models"
	"gorm.io/gorm"
)

// nolint
// -package=mock_manager
//
//go:generate mockgen -source=$GOFILE -destination=../../../mock/pkg/pipelinerun/manager/mock_manager.go
type Manager interface {
	Create(ctx context.Context, pipelinerun *models.Pipelinerun) (*models.Pipelinerun, error)
	GetByID(ctx context.Context, pipelinerunID uint) (*models.Pipelinerun, error)
	GetByCIEventID(ctx context.Context, ciEventID string) (*models.Pipelinerun, error)
	GetByClusterID(ctx context.Context, clusterID uint, canRollback bool,
		query q.Query) (int, []*models.Pipelinerun, error)
	GetFirstCanRollbackPipelinerun(ctx context.Context, clusterID uint) (*models.Pipelinerun, error)
	DeleteByID(ctx context.Context, pipelinerunID uint) error
	DeleteByClusterID(ctx context.Context, clusterID uint) error
	UpdateConfigCommitByID(ctx context.Context, pipelinerunID uint, commit string) error
	GetLatestByClusterIDAndActions(ctx context.Context, clusterID uint, actions ...string) (*models.Pipelinerun, error)
	GetLatestByClusterIDAndActionAndStatus(ctx context.Context, clusterID uint, action,
		status string) (*models.Pipelinerun, error)
	GetLatestSuccessByClusterID(ctx context.Context, clusterID uint) (*models.Pipelinerun, error)
	UpdateStatusByID(ctx context.Context, pipelinerunID uint, result models.PipelineStatus) error
	UpdateCIEventIDByID(ctx context.Context, pipelinerunID uint, ciEventID string) error
	// UpdateResultByID  update the pipelinerun restore result
	UpdateResultByID(ctx context.Context, pipelinerunID uint, result *models.Result) error
}

type manager struct {
	dao dao.DAO
}

func New(db *gorm.DB) Manager {
	return &manager{
		dao: dao.NewDAO(db),
	}
}

func (m *manager) Create(ctx context.Context, pipelinerun *models.Pipelinerun) (*models.Pipelinerun, error) {
	return m.dao.Create(ctx, pipelinerun)
}

func (m *manager) GetByID(ctx context.Context, pipelinerunID uint) (*models.Pipelinerun, error) {
	return m.dao.GetByID(ctx, pipelinerunID)
}

func (m *manager) GetByCIEventID(ctx context.Context, ciEventID string) (*models.Pipelinerun, error) {
	return m.dao.GetByCIEventID(ctx, ciEventID)
}

func (m *manager) DeleteByID(ctx context.Context, pipelinerunID uint) error {
	return m.dao.DeleteByID(ctx, pipelinerunID)
}

func (m *manager) DeleteByClusterID(ctx context.Context, clusterID uint) error {
	return m.dao.DeleteByClusterID(ctx, clusterID)
}

func (m *manager) UpdateConfigCommitByID(ctx context.Context, pipelinerunID uint, commit string) error {
	return m.dao.UpdateConfigCommitByID(ctx, pipelinerunID, commit)
}

func (m *manager) GetLatestByClusterIDAndActions(ctx context.Context,
	clusterID uint, actions ...string) (*models.Pipelinerun, error) {
	return m.dao.GetLatestByClusterIDAndActions(ctx, clusterID, actions...)
}

func (m *manager) GetLatestSuccessByClusterID(ctx context.Context, clusterID uint) (*models.Pipelinerun, error) {
	return m.dao.GetLatestSuccessByClusterID(ctx, clusterID)
}

func (m *manager) UpdateStatusByID(ctx context.Context, pipelinerunID uint, result models.PipelineStatus) error {
	return m.dao.UpdateStatusByID(ctx, pipelinerunID, result)
}

func (m *manager) UpdateCIEventIDByID(ctx context.Context, pipelinerunID uint, ciEventID string) error {
	return m.dao.UpdateCIEventIDByID(ctx, pipelinerunID, ciEventID)
}

func (m *manager) UpdateResultByID(ctx context.Context, pipelinerunID uint, result *models.Result) error {
	return m.dao.UpdateResultByID(ctx, pipelinerunID, result)
}

func (m *manager) GetByClusterID(ctx context.Context,
	clusterID uint, canRollback bool, query q.Query) (int, []*models.Pipelinerun, error) {
	return m.dao.GetByClusterID(ctx, clusterID, canRollback, query)
}

func (m *manager) GetFirstCanRollbackPipelinerun(ctx context.Context, clusterID uint) (*models.Pipelinerun, error) {
	return m.dao.GetFirstCanRollbackPipelinerun(ctx, clusterID)
}

func (m *manager) GetLatestByClusterIDAndActionAndStatus(ctx context.Context, clusterID uint, action,
	status string) (*models.Pipelinerun, error) {
	return m.dao.GetLatestByClusterIDAndActionAndStatus(ctx, clusterID, action, status)
}
