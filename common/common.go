// Copyright 2021 CertSync Contributors
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
//
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"flag"
	"fmt"

	cs "github.com/icemarkom/certsync"
)

func ProgramVersion(cfg *cs.Config) {
	fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n Commit: %s\n", cfg.Version, cfg.GitCommit)
}

func ProgramUsage(cfg *cs.Config) {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", cfg.BinaryName)
	flag.PrintDefaults()
	fmt.Fprintln(flag.CommandLine.Output())
	ProgramVersion(cfg)
}
