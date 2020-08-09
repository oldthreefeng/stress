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
	"bytes"
	"fmt"
	"github.com/oldthreefeng/stress/pkg"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	rootExample = `	
	# stress curl file to test 
	stress -f utils/curl.txt

	# stress curl file read from stdin 
 	cat utils/curl.txt | stress -f -

	# stress concurrency 10 & 10 times
	stress -c 10 -n 10 -f  utils/curl.txt

	# stress cli url
	stress -c 10 -n 100 -u https://www.baidu.com
	
	# curl.txt example
	cat utils/curl.txt
curl 'https://www.baidu.com' -H 'User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:79.0) Gecko/20100101 Firefox/79.0' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8' -H 'Accept-Language: zh-CN,en-US;q=0.7,en;q=0.3' --compressed -H 'Connection: keep-alive' -H 'Upgrade-Insecure-Requests: 1' -H 'Cache-Control: max-age=0' -H 'TE: Trailers'
`
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stress",
	Short: "stress is a test cli for http and websocket stress written by golang",
	Long: `stress is a test cli for http and websocket stress written by golang, 
go 实现的压测工具，每个用户用一个协程的方式模拟，最大限度的利用 CPU 资源`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Example: rootExample,
	Run: func(cmd *cobra.Command, args []string) {
		if pkg.Path != "" {
			if pkg.Path == "-" {
				err := ReadStdin()
				if err != nil {
					fmt.Println("Read stdin err: ", err)
					return
				}
			}
			StartCurl()
		} else {
			Start()
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		// 没有请求， 也没有curl文件
		if pkg.RequestUrl == "" && pkg.Path == "" {
			fmt.Println(VersionStr)
			os.Exit(-1)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file for stress (default is $HOME/.stress.yaml)")
	rootCmd.PersistentFlags().Uint64VarP(&pkg.Concurrency, "concurrency", "c", 1, "并发数")
	rootCmd.PersistentFlags().Uint64VarP(&pkg.Number, "number", "n", 1, "单协程的请求数")
	rootCmd.PersistentFlags().StringVarP(&pkg.RequestUrl, "requestUrl", "u", "", "curl文件路径")
	rootCmd.PersistentFlags().StringVarP(&pkg.Path, "path", "f", "", "read curl file to build test")

	rootCmd.PersistentFlags().StringVarP(&pkg.VerifyStr, "verify", "v", pkg.DefaultVerifyCode, " verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json")
	rootCmd.PersistentFlags().StringVarP(&pkg.Body, "data", "", "", "http post data same with --data-raw")
	rootCmd.PersistentFlags().StringVarP(&pkg.Body, "data-raw", "", "", "http post data same with --data")
	rootCmd.PersistentFlags().StringSliceVarP(&pkg.Header, "header", "H", []string{}, "http post data")
	rootCmd.PersistentFlags().BoolVarP(&pkg.Debug, "debug", "d", false, "debug 模式")
	rootCmd.PersistentFlags().BoolVarP(&pkg.Compressed, "compressed", "", false, "使用gzip压缩算法去请求 。同curl --compressed gzip")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".stress" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".stress")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// Start is for cli url stress
func Start() {
	request, err := pkg.NewRequest(pkg.RequestUrl, pkg.VerifyStr, 0, pkg.Debug, pkg.Header, pkg.Body)
	if err != nil {
		log.Fatal()
		return
	}
	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", pkg.Concurrency, pkg.Number)
	request.Print()
	Dispose(pkg.Concurrency, pkg.Number, request)
}

// StartCurl is for curl file
func StartCurl() {
	request := pkg.NewDefaultRequest()
	list := pkg.GetRequestListFromFile(pkg.Path)
	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", pkg.Concurrency, pkg.Number)
	for k, v := range list {
		fmt.Printf("%d step: \n", k)
		v.Print()
	}
	Dispose(pkg.Concurrency, pkg.Number, request)
}

// ReadStdin is read curl file from stdin
func ReadStdin() error {
	var b bytes.Buffer
	_, err := b.ReadFrom(os.Stdin)
	if err != nil {
		return err
	}
	pkg.Path = "/tmp/curl.tmp"
	return ioutil.WriteFile(pkg.Path, b.Bytes(), 0660)
}
