// Copyright 2019 Yunion
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

package models

import (
	"context"
	"database/sql"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/reflectutils"
	"yunion.io/x/sqlchemy"

	api "yunion.io/x/onecloud/pkg/apis/cloudid"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/util/stringutils2"
)

type SCloudgroupResourceBaseManager struct {
}

type SCloudgroupResourceBase struct {
	CloudgroupId string `width:"36" charset:"ascii" nullable:"false" list:"user" create:"required"`
}

func (self *SCloudgroupResourceBase) GetCloudgroup() (*SCloudgroup, error) {
	group, err := CloudgroupManager.FetchById(self.CloudgroupId)
	if err != nil {
		return nil, errors.Wrap(err, "FetchById")
	}
	return group.(*SCloudgroup), nil
}

func (manager *SCloudgroupResourceBaseManager) ListItemFilter(ctx context.Context, q *sqlchemy.SQuery, groupCred mcclient.TokenCredential, query api.CloudgroupResourceListInput) (*sqlchemy.SQuery, error) {
	if len(query.CloudgroupId) > 0 {
		group, err := CloudgroupManager.FetchByIdOrName(ctx, nil, query.CloudgroupId)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, httperrors.NewResourceNotFoundError2("cloudgroup", query.CloudgroupId)
			}
			return nil, httperrors.NewGeneralError(err)
		}
		q = q.Equals("cloudgroup_id", group.GetId())
	}
	return q, nil
}

func (manager *SCloudgroupResourceBaseManager) FetchCustomizeColumns(
	ctx context.Context,
	groupCred mcclient.TokenCredential,
	query jsonutils.JSONObject,
	objs []interface{},
	fields stringutils2.SSortedStrings,
	isList bool,
) []api.CloudgroupResourceDetails {
	rows := make([]api.CloudgroupResourceDetails, len(objs))
	groupIds := make([]string, len(objs))
	for i := range objs {
		var base *SCloudgroupResourceBase
		err := reflectutils.FindAnonymouStructPointer(objs[i], &base)
		if err != nil {
			log.Errorf("Cannot find SCloudgroupResourceBase in %#v: %s", objs[i], err)
		} else if base != nil && len(base.CloudgroupId) > 0 {
			groupIds[i] = base.CloudgroupId
		}
	}
	groupMaps, err := db.FetchIdNameMap2(CloudgroupManager, groupIds)
	if err != nil {
		return rows
	}
	for i := range rows {
		rows[i].Cloudgroup, _ = groupMaps[groupIds[i]]
	}
	return rows
}
