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

package helper // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type ByteSize int64

func (h *ByteSize) UnmarshalJSON(raw []byte) error {
	return h.unmarshalShared(func(i interface{}) error {
		return json.Unmarshal(raw, &i)
	})
}

func (h *ByteSize) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return h.unmarshalShared(unmarshal)
}

func (h *ByteSize) UnmarshalText(text []byte) (err error) {
	slice := make([]byte, 1, 2+len(text))
	slice[0] = byte('"')
	slice = append(slice, text...)
	slice = append(slice, byte('"'))
	return h.UnmarshalJSON(slice)
}

var byteSizeRegex = regexp.MustCompile(`^([0-9]+\.?[0-9]*)\s*([kKmMgGtTpP]i?[bB])?$`)

func (h *ByteSize) unmarshalShared(unmarshal func(interface{}) error) error {
	var intType int64
	if err := unmarshal(&intType); err == nil {
		*h = ByteSize(intType)
		return nil
	}

	var floatType float64
	if err := unmarshal(&floatType); err == nil {
		*h = ByteSize(int64(floatType))
		return nil
	}

	var stringType string
	if err := unmarshal(&stringType); err != nil {
		return fmt.Errorf("failed to unmarshal to int64, float64, or string: %w", err)
	}

	matches := byteSizeRegex.FindStringSubmatch(stringType)
	if matches == nil {
		return fmt.Errorf("invalid byte size '%s'", stringType)
	}

	numeral, err := strconv.ParseFloat(matches[1], 32)
	if err != nil {
		return fmt.Errorf("invalid numeric base '%s'", matches[1])
	}

	var multiplier float64
	switch strings.ToLower(matches[2]) {
	case "":
		multiplier = 1
	case "kb":
		multiplier = 1000
	case "kib":
		multiplier = 1024
	case "mb":
		multiplier = 1000 * 1000
	case "mib":
		multiplier = 1024 * 1024
	case "gb":
		multiplier = 1000 * 1000 * 1000
	case "gib":
		multiplier = 1024 * 1024 * 1024
	case "tb":
		multiplier = 1000 * 1000 * 1000 * 1000
	case "tib":
		multiplier = 1024 * 1024 * 1024 * 1024
	case "pb":
		multiplier = 1000 * 1000 * 1000 * 1000 * 1000
	case "pib":
		multiplier = 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return fmt.Errorf("invalid unit '%s'", matches[2])
	}

	*h = ByteSize(int64(multiplier * numeral))
	return nil
}
