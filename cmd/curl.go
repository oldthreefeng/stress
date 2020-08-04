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
	"github.com/oldthreefeng/stress/pkg"

	"github.com/spf13/cobra"
)

var (
	curlExample = `
	
	# stress curl file to test 
	stress curl -f /tmp/curl.txt

	# stress curl file read from stdin 
 	cat a.txt | stress curl -f -

`
)

// curlCmd represents the curl command
var curlCmd = &cobra.Command{
	Use:   "curl",
	Short: "stress curl -f file ",
	Long: `stress curl is usr curl file to build stress http testing`,
	Example: curlExample,
	Run: func(cmd *cobra.Command, args []string) {
		StartCurl()
	},
}

func init() {
	rootCmd.AddCommand(curlCmd)
	rootCmd.PersistentFlags().StringVarP(&pkg.Path, "path", "f","" , "read curl file to build test")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// curlCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// curlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func StartCurl() {
	request := pkg.NewDefaultRequest()
	list := pkg.GetRequestListFromFile(pkg.Path)
	for k, v := range list {
		fmt.Printf("%d step", k)
		v.Print()
	}
	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", pkg.Concurrency, pkg.Number)
	Dispose(pkg.Concurrency, pkg.Number, request)
}
