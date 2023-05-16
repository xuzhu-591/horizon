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

package common

import triggers "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"

const (
	ClusterQueryEnvironment   = "environment"
	ClusterQueryName          = "filter"
	ClusterQueryByUser        = "userID"
	ClusterQueryByTemplate    = "template"
	ClusterQueryByRelease     = "templateRelease"
	ClusterQueryByRegion      = "region"
	ClusterQueryTagSelector   = "tagSelector"
	ClusterQueryScope         = "scope"
	ClusterQueryMergePatch    = "mergePatch"
	ClusterQueryInheritConfig = "inheritConfig"
	ClusterQueryTargetBranch  = "targetBranch"
	ClusterQueryTargetCommit  = "targetCommit"
	ClusterQueryTargetTag     = "targetTag"
	ClusterQueryContainerName = "containerName"
	ClusterQueryPodName       = "podName"
	ClusterQueryTailLines     = "tailLines"
	ClusterQueryStart         = "start"
	ClusterQueryEnd           = "end"
	ClusterQueryExtraOwner    = "extraOwner"
	ClusterQueryHard          = "hard"

	// ClusterQueryIsFavorite is used to query cluster with favorite for current user only.
	ClusterQueryIsFavorite = "isFavorite"
	// ClusterQueryWithFavorite is used to query cluster with favorite field inside.
	ClusterQueryWithFavorite = "withFavorite"
	ClusterQueryAction       = "action"

	ClusterQueryByGVK = "gvk"

	ClusterQueryResourceName = "resourceName"
)

const (
	ClusterClusterLabelKey = "cloudnative.music.netease.com/cluster"
	ClusterRestartTimeKey  = "cloudnative.music.netease.com/user-restart-time"
)

// status of cluster
const (
	ClusterStatusEmpty    = ""
	ClusterStatusFreeing  = "Freeing"
	ClusterStatusFreed    = "Freed"
	ClusterStatusDeleting = "Deleting"
	ClusterStatusCreating = "Creating"
)

const (
	TektonTriggersEventIDKey = triggers.GroupName + triggers.EventIDLabelKey
)
