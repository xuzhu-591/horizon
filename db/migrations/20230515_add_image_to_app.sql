-- Copyright © 2023 Horizoncd.
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

ALTER TABLE tb_application
ADD COLUMN `image` varchar(256) DEFAULT NULL COMMENT 'artifact image url'
COMMENT 'artifact image url for the application';

ALTER TABLE tb_cluster
ADD COLUMN `image` varchar(256) DEFAULT NULL COMMENT 'artifact image url'
COMMENT 'artifact image url for the cluster';
