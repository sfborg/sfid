/*
Copyright © 2024 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/sfborg/sfid/ent"
	sfid "github.com/sfborg/sfid/pkg"
	"github.com/sfborg/sfid/pkg/config"
	"github.com/spf13/cobra"
)

var opts []config.Option

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sfid",
	Short: "traverses a directory and generates UUID v5 and Sha256 hashes",
	Run: func(cmd *cobra.Command, args []string) {
		versionFlag(cmd)
		flags := []flagFunc{
			jobsNumFlag, shaFlag, uuidFlag, recursiveFlag, gnFlag, twFlag,
		}
		for _, v := range flags {
			v(cmd)
		}

		cfg := config.New(opts...)

		if len(args) != 1 {
			cmd.Help()
			os.Exit(0)
		}
		sf := sfid.New(cfg)
		chOut := make(chan *ent.Output)

		go func() {
			for o := range chOut {
				fmt.Println(o.String())
			}
		}()

		err := sf.Process(args[0], chOut)
		if err != nil {
			slog.Error("Cannot process", "input", args[0], "error", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("gn_namespace", "g", false, "globalnames.org namespace for UUID v5")
	rootCmd.Flags().BoolP("tw_namespace", "t", false, "taxonworks.org namespace for UUID v5")
	rootCmd.Flags().IntP("jobs_number", "j", 0, "jobs number")
	rootCmd.Flags().BoolP("version", "V", false, "show version of the app")
	rootCmd.Flags().BoolP("uuid", "u", false, "include UUID v5")
	rootCmd.Flags().BoolP("sha256", "s", false, "include SHA256 hash")
	rootCmd.Flags().BoolP("recursive", "r", false, "traverse a directory recursively")

}
