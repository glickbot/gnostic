// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// gnostic_lint is a tool for analyzing OpenAPI descriptions.
//
// It scans an API description and checks it against a set of
// coding style guidelines.
//
// Results are returned in a JSON structure.
package main

import (
	"encoding/json"
	"path"
	"strings"

	plugins "github.com/googleapis/gnostic/plugins"
)

type DocumentLinter interface {
	Run()
}

// This is the main function for the plugin.
func main() {
	env, err := plugins.NewEnvironment()
	env.RespondAndExitIfError(err)

	var linter DocumentLinter
	if env.Request.Openapi2 != nil {
		linter = NewDocumentLinterV2(env.Request.Openapi2)
	}

	if env.Request.Openapi3 != nil {
		linter = NewDocumentLinterV3(env.Request.Openapi3)
	}

	if linter != nil {
		linter.Run()
		// Return the analysis results with an appropriate filename.
		// Results are in files named "lint.json" in the same relative
		// locations as the description source files.
		file := &plugins.File{}
		file.Name = strings.Replace(env.Request.SourceName, path.Base(env.Request.SourceName), "lint.json", -1)
		file.Data, err = json.MarshalIndent(linter, "", "  ")
		file.Data = append(file.Data, []byte("\n")...)
		env.RespondAndExitIfError(err)
		env.Response.Files = append(env.Response.Files, file)
	}

	env.RespondAndExit()
}