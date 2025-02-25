/*
Copyright 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringerGCloudPubSubResourceName(t *testing.T) {
	testCases := []struct {
		name         string
		input        GCloudPubSubResourceName
		expectOutput string
	}{{
		name: "Valid resource name",
		input: GCloudPubSubResourceName{
			Project:    "p",
			Collection: "c",
			Resource:   "r",
		},
		expectOutput: "projects/p/c/r",
	}, {
		name: "Invalid resource name",
		input: GCloudPubSubResourceName{
			Project: "",
		},
		expectOutput: "",
	}}

	for _, tc := range testCases {
		//nolint:scopelint
		t.Run(tc.name, func(t *testing.T) {
			s := tc.input.String()
			assert.Equal(t, tc.expectOutput, s)
		})
	}
}

func TestMarshalGCloudPubSubResourceName(t *testing.T) {
	testCases := []struct {
		name              string
		input             GCloudPubSubResourceName
		expectOutput      string
		expectErrContains string
	}{{
		name: "All fields are filled in",
		input: GCloudPubSubResourceName{
			Project:    "p",
			Collection: "c",
			Resource:   "r",
		},
		expectOutput: `"projects/p/c/r"`,
	}, {
		name: "Some required fields are empty",
		input: GCloudPubSubResourceName{
			Project: "p",
		},
		expectErrContains: "resource name contains empty attributes",
	}, {
		name:              "All fields are empty",
		input:             GCloudPubSubResourceName{},
		expectErrContains: "resource name contains empty attributes",
	}}

	for _, tc := range testCases {
		//nolint:scopelint
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.input)

			assert.Equal(t, tc.expectOutput, string(b))

			if errStr := tc.expectErrContains; errStr != "" {
				assert.Contains(t, err.Error(), errStr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUnmarshalGCloudPubSubResourceName(t *testing.T) {
	testCases := []struct {
		name              string
		input             string
		expectOutput      GCloudPubSubResourceName
		expectErrContains string
	}{{
		name:  "Valid resource name format",
		input: `"projects/p/c/r"`,
		expectOutput: GCloudPubSubResourceName{
			Project:    "p",
			Collection: "c",
			Resource:   "r",
		},
	}, {
		name:              "Some required fields are empty",
		input:             `"projects//c/r"`,
		expectErrContains: "resource name contains empty attributes",
	}, {
		name:              "Invalid format",
		input:             `"projects/p"`,
		expectErrContains: "does not match expected format",
	}, {
		name:              "Invalid input",
		input:             `not_a_resource_name`,
		expectErrContains: "invalid character",
	}}

	for _, tc := range testCases {
		//nolint:scopelint
		t.Run(tc.name, func(t *testing.T) {
			resName := &GCloudPubSubResourceName{}
			err := json.Unmarshal([]byte(tc.input), resName)

			assert.Equal(t, tc.expectOutput, *resName)

			if errStr := tc.expectErrContains; errStr != "" {
				assert.Contains(t, err.Error(), errStr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
