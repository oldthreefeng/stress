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
	"log"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stress",
	Short: "stress is a test cli for http and websocket stress written by golang",
	Long: `stress is a test cli for http and websocket stress written by golang, 
go 实现的压测工具，每个用户用一个协程的方式模拟，最大限度的利用 CPU 资源`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		Start()
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
	rootCmd.PersistentFlags().Uint64VarP(&pkg.Concurrency, "concurrency", "c" ,1, "并发数")
	rootCmd.PersistentFlags().Uint64VarP(&pkg.Number, "number", "n" ,1, "单协程的请求数")
	rootCmd.PersistentFlags().StringVarP(&pkg.RequestUrl, "requestUrl", "u", "", "curl文件路径")

	rootCmd.PersistentFlags().StringVarP(&pkg.VerifyStr, "verify", "v","" , " verify 验证方法 在server/verify中 http 支持:statusCode、json webSocket支持:json")
	rootCmd.PersistentFlags().StringVarP(&pkg.Body, "data", "","" , "http post data")
	rootCmd.PersistentFlags().StringSliceVarP(&pkg.Header, "header", "H",[]string{}, "http post data")
	rootCmd.PersistentFlags().BoolVarP(&pkg.Debug, "debug", "d",false, "debug 模式")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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

func Start()  {
	request, err := pkg.NewRequest(pkg.RequestUrl,pkg.VerifyStr,0,pkg.Debug,pkg.Header,pkg.Body)
	if err != nil {
		log.Fatal()
		return
	}
	fmt.Printf("\n 开始启动  并发数:%d 请求数:%d 请求参数: \n", pkg.Concurrency, pkg.Number)
	request.Print()
	Dispose(pkg.Concurrency,pkg.Number,request)
}