// Copyright 2017 Authors of Cilium
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

package proxy

import (
	"github.com/cilium/cilium/pkg/labels"
	"github.com/cilium/cilium/pkg/lock"
	"github.com/cilium/cilium/pkg/policy"
	"github.com/cilium/cilium/pkg/proxy/accesslog"
)

type proxySourceMocker struct {
	lock.RWMutex
	id       uint64
	ipv4     string
	ipv6     string
	labels   []string
	identity policy.NumericIdentity
}

func (m *proxySourceMocker) GetEndpointInfo() *accesslog.EndpointInfo {
	return &accesslog.EndpointInfo{
		ID:           m.id,
		IPv4:         m.ipv4,
		IPv6:         m.ipv6,
		Identity:     uint64(m.identity),
		LabelsSHA256: labels.NewLabelsFromModel(m.labels).SHA256Sum(),
		Labels:       m.labels,
	}
}

func (m *proxySourceMocker) ResolveIdentity(policy.NumericIdentity) *policy.Identity {
	identity := policy.NewIdentity()
	identity.ID = m.identity
	identity.Labels = labels.NewLabelsFromModel(m.labels)
	identity.LabelsSHA256 = identity.Labels.SHA256Sum()
	return identity
}
