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
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/cobra"
)

// unlinkCmd represents the clean command
var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		tn := path.Base(wd)
		if len(args) > 0 {
			tn = args[0]
		}
		td := fmt.Sprintf("%s/%s", config.TemplateDir, tn)
		files, err := ioutil.ReadDir(td)
		if err != nil {
			_ = fmt.Errorf("unable to decode into struct, %v", err)
			os.Exit(1)
		}
		for _, file := range files {
			err := os.Remove(fmt.Sprintf("./%s", file.Name()))
			if err == nil {
				_, _ = fmt.Printf("removed %s symlink\n", file.Name())
			} else {
				_, _ = fmt.Printf("not removed %s symlink\n", file.Name())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(unlinkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unlinkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unlinkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}