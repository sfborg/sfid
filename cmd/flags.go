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

func uuidFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("uuid")
	if b {
		opts = append(opts, config.OptWithUUID(b))
	}
}

func shaFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("sha256")
	if b {
		opts = append(opts, config.OptWithSha(b))
	}
}

func recursiveFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("recursive")
	if b {
		opts = append(opts, config.OptRecursive(b))
	}
}

func gnFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("gn_namespace")
	if b {
		opts = append(opts, config.OptNameSpace("globalnames.org"))
	}
}

func twFlag(cmd *cobra.Command) {
	b, _ := cmd.Flags().GetBool("tw_namespace")
	if b {
		opts = append(opts, config.OptNameSpace("taxonworks.org"))
	}
}
