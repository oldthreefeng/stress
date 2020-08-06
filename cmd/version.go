/*
Copyright © 2020 louisehong

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
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "stress version",
	Long:  `stress build INFO`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(VersionStr)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var (
	Version    = "latest"
	Build      = ""
	BuildTime  = ""
	BuildBy    = "louisehong"
	VersionStr = fmt.Sprintf("stress version: %v, build git hash: %v ,\ngo version: %v ,Build Time : %v\nBuildBy: %v", Version, Build, runtime.Version(), BuildTime, BuildBy)
)
