package common

import (
	"flag"
	"fmt"

	cs "github.com/icemarkom/certsync"
)

func ProgramVersion(cfg cs.Config) {
	fmt.Fprintf(flag.CommandLine.Output(), "Version: %s\n Commit: %s\n", cfg.Version, cfg.GitCommit)
}

func ProgramUsage(cfg cs.Config) {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", cfg.BinaryName)
	flag.PrintDefaults()
	fmt.Fprintln(flag.CommandLine.Output())
	ProgramVersion(cfg)
}
