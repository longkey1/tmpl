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
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "A brief description of your command",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a color argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		tn := path.Base(wd)
		td := fmt.Sprintf("%s/%s", config.TemplateDir, tn)
		if _, err := os.Stat(td); os.IsNotExist(err) {
			err := os.MkdirAll(td, os.ModePerm)
			if err != nil {
				log.Fatalf("unable to make %s directory", td)
			}
		}
		tf := fmt.Sprintf("%s/%s", td, args[0])
		if _, err := os.Stat(tf); !os.IsNotExist(err) {
			log.Fatalf("already %s exists", tf)
		}

		err = os.Rename(fmt.Sprintf("%s/%s", wd, tf), fmt.Sprintf("%s/%s", td, tf))
		if _, err := os.Stat(tf); !os.IsNotExist(err) {
			log.Fatalf("unable to move from %s to %s", fmt.Sprintf("%s/%s", wd, tf), fmt.Sprintf("%s/%s", td, tf))
		}

		err = os.Symlink(fmt.Sprintf("%s/%s", td, tf), fmt.Sprintf("./%s", tf))
		if err != nil {
			log.Fatalf("unable to symlink from %s to %s", fmt.Sprintf("%s/%s", td, tf), fmt.Sprintf("./%s", tf))
		}

		_, _ = fmt.Printf("made %s template to %s, and symlink\n", tf, tn)
	},
}

func init() {
	rootCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
