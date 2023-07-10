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

package generator

import (
	"net/http"

	"github.com/horizoncd/horizon/pkg/token/models"
)

type CodeGenerateInfo struct {
	Token   models.Token
	Request *http.Request
}

type AuthorizationCodeGenerator interface {
	Generate(info *CodeGenerateInfo) string
}

type AccessTokenCodeGenerator interface {
	Generate(info *CodeGenerateInfo) string
}

type RefreshTokenCodeGenerator interface {
	Generate(info *CodeGenerateInfo) string
}
