// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package confgenerator

import "github.com/GoogleCloudPlatform/ops-agent/confgenerator/fluentbit"

// LoggingReceiverAgent provides the logs generated by the agent itself.
// It is never referenced in the config file, and instead is forcibly added in confgenerator.go.
// Therefore, it does not need to implement any interfaces.
type LoggingReceiverAgent struct {
	UserAgent string
	OS        string
}

func (r LoggingReceiverAgent) Components(tag string) []fluentbit.Component {
	tails := []fluentbit.Component{}
	if r.OS != "windows" {
		tails = append(tails, fluentbit.Tail{
			Tag:          tag + "-collectd",
			IncludePaths: []string{"${logs_dir}/metrics-module.log"},
		}.Component())
	}
	tails = append(tails, fluentbit.Tail{
		Tag:          tag + "-fluent-bit",
		IncludePaths: []string{"${logs_dir}/logging-module.log"},
	}.Component())
	output := fluentbit.Stackdriver{
		Match:     "ops-agent-fluent-bit|ops-agent-collectd",
		Workers:   8,
		UserAgent: r.UserAgent,
	}.Component()
	return append(tails, output)
}

// intentionally not registered as a component because this is not created by users
