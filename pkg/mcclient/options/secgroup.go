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

package options

import (
	"yunion.io/x/jsonutils"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/secrules"

	"yunion.io/x/onecloud/pkg/apis"
)

type SecgroupListOptions struct {
	BaseListOptions

	Equals    string `help:"Secgroup ID or Name, filter secgroups whose rules equals the specified one"`
	Server    string `help:"Filter secgroups bound to specified server"`
	Ip        string `help:"Filter secgroup by ip"`
	Ports     string `help:"Filter secgroup by ports"`
	Direction string `help:"Filter secgroup by ports" choices:"all|in|out"`
}

func (opts *SecgroupListOptions) Params() (jsonutils.JSONObject, error) {
	return ListStructToParams(opts)
}

type SecgroupCreateOptions struct {
	BaseCreateOptions
	Rules []string `help:"security rule to create"`
}

func (opts *SecgroupCreateOptions) Params() (jsonutils.JSONObject, error) {
	params := jsonutils.Marshal(opts).(*jsonutils.JSONDict)
	params.Remove("rules")
	rules := []secrules.SecurityRule{}
	for i, ruleStr := range opts.Rules {
		rule, err := secrules.ParseSecurityRule(ruleStr)
		if err != nil {
			return nil, errors.Wrapf(err, "ParseSecurityRule(%s)", ruleStr)
		}
		rule.Priority = i + 1
		rules = append(rules, *rule)
	}
	if len(rules) > 0 {
		params.Add(jsonutils.Marshal(rules), "rules")
	}
	return params, nil
}

type SecgroupIdOptions struct {
	ID string `help:"ID or Name of security group destination"`
}

func (opts *SecgroupIdOptions) GetId() string {
	return opts.ID
}

func (opts *SecgroupIdOptions) Params() (jsonutils.JSONObject, error) {
	return nil, nil
}

type SecgroupMergeOptions struct {
	SecgroupIdOptions
	SECGROUPS []string `help:"source IDs or Names of secgroup"`
}

func (opts *SecgroupMergeOptions) Params() (jsonutils.JSONObject, error) {
	return jsonutils.Marshal(map[string][]string{"secgruops": opts.SECGROUPS}), nil
}

type SecgroupsAddRuleOptions struct {
	SecgroupIdOptions
	DIRECTION   string `help:"Direction of rule" choices:"in|out"`
	PROTOCOL    string `help:"Protocol of rule" choices:"any|tcp|udp|icmp"`
	ACTION      string `help:"Actin of rule" choices:"allow|deny"`
	PRIORITY    int    `help:"Priority for rule, range 1 ~ 100"`
	Cidr        string `help:"IP or CIRD for rule"`
	Description string `help:"Desciption for rule"`
	Ports       string `help:"Port for rule"`
}

func (opts *SecgroupsAddRuleOptions) Params() (jsonutils.JSONObject, error) {
	params := jsonutils.Marshal(opts).(*jsonutils.JSONDict)
	params.Remove("id")
	return params, nil
}

type SecurityGroupCacheOptions struct {
	SecgroupIdOptions
	VPC     string `help:"ID or Name of vpc"`
	Classic *bool  `help:"Is classic vpc"`
}

func (opts *SecurityGroupCacheOptions) Params() (jsonutils.JSONObject, error) {
	params := jsonutils.Marshal(opts).(*jsonutils.JSONDict)
	params.Remove("id")
	return params, nil
}

type SecurityGroupUncacheSecurityGroup struct {
	SecgroupIdOptions
	CACHE string `help:"ID of secgroup cache"`
}

func (opts *SecurityGroupUncacheSecurityGroup) Params() (jsonutils.JSONObject, error) {
	params := jsonutils.Marshal(opts).(*jsonutils.JSONDict)
	params.Remove("id")
	return params, nil
}

type SecgroupChangeOwnerOptions struct {
	SecgroupIdOptions
	apis.ProjectizedResourceInput
}
