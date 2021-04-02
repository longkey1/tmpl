/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type Config struct {
	TemplateDir string `mapstructure:"template_dir"`
}

var config Config

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: "0.0.2",
	Use:   "tmpl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%#v\n", config.TemplateDir)
		dir := fmt.Sprintf("%s/%s", config.TemplateDir, args[0])
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			_ = fmt.Errorf("unable to decode into struct, %v", err)
			os.Exit(1)
		}
		for _, file := range files {
			err := os.Symlink(fmt.Sprintf("%s/%s", dir, file.Name()), fmt.Sprintf("./%s", file.Name()))
			if err == nil {
				_, _ = fmt.Printf("created %s symlink\n", file.Name())
			} else {
				_, _ = fmt.Printf("could not created %s symlink\n", file.Name())
			}
		}
	},
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires template name")
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tmpl.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile == "" {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		cfgFile = home+"/.config/tmpl/config.toml"
	}
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		config.TemplateDir = filepath.Dir(cfgFile)+"/templates"
	} else {
		if err := viper.Unmarshal(&config); err != nil {
			_ = fmt.Errorf("unable to decode into struct, %v", err)
			os.Exit(1)
		}
	}
}
