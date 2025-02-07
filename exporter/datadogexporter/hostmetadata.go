// Copyright The OpenTelemetry Authors
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

package datadogexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogexporter"

import (
	"fmt"
	"strings"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/datadogexporter/internal/metadata"
)

// getHostTags gets the host tags extracted from the configuration.
func getHostTags(c *Config) []string {
	tags := c.HostMetadata.Tags

	if len(tags) == 0 {
		//lint:ignore SA1019 Will be removed when environment variable detection is removed
		tags = strings.Split(c.EnvVarTags, " ") //nolint
	}

	if c.Env != "none" {
		tags = append(tags, fmt.Sprintf("env:%s", c.Env))
	}
	return tags
}

// newMetadataConfigfromConfig creates a new metadata pusher config from the main
func newMetadataConfigfromConfig(cfg *Config) metadata.PusherConfig {
	return metadata.PusherConfig{
		ConfigHostname:      cfg.Hostname,
		ConfigTags:          getHostTags(cfg),
		MetricsEndpoint:     cfg.Metrics.Endpoint,
		APIKey:              cfg.API.Key,
		UseResourceMetadata: cfg.HostMetadata.HostnameSource == HostnameSourceFirstResource,
		InsecureSkipVerify:  cfg.TLSSetting.InsecureSkipVerify,
		TimeoutSettings:     cfg.TimeoutSettings,
		RetrySettings:       cfg.RetrySettings,
	}
}
