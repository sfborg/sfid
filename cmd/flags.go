package cmd

import (
	"fmt"
	"os"

	sfid "github.com/sfborg/sfid/pkg"
	"github.com/sfborg/sfid/pkg/config"
	"github.com/spf13/cobra"
)

type flagFunc func(cmd *cobra.Command)

func versionFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("version")
	if b {
		version := sfid.GetVersion()
		fmt.Printf(
			"\nVersion: %s\nBuild:   %s\n\n",
			version.Version,
			version.Build,
		)
		os.Exit(0)
	}
}

func jobsNumFlag(cmd *cobra.Command) {
	i, _ := cmd.Flags().GetInt("jobs-number")
	if i > 0 {
		opts = append(opts, config.OptJobsNum(i))
	}
}
