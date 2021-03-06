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
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

// unlinkCmd represents the clean command
var unlinkCmd = &cobra.Command{
	Use:   "unlink",
	Short: "remove symbolic link or hard link to template files",
	Run: func(cmd *cobra.Command, args []string) {
		wd, err := os.Getwd()
		cobra.CheckErr(err)
		tn := path.Base(wd)

		t, err := cmd.Flags().GetString("target")
		cobra.CheckErr(err)
		if len(t) > 0 {
			tn = t
		}

		td := fmt.Sprintf("%s/%s", config.TemplateDir, tn)
		files, err := ioutil.ReadDir(td)
		cobra.CheckErr(err)
		for _, file := range files {
			err := os.Remove(fmt.Sprintf("./%s", file.Name()))
			if err == nil {
				_, _ = fmt.Printf("removed %s link\n", file.Name())
			} else {
				_, _ = fmt.Printf("not removed %s link\n", file.Name())
			}
		}

		gitInfoExclude := filepath.Join(".git", "info", "exclude")
		gitInfoExcludeNew := filepath.Join(".git", "info", "excludeNew")
		_, err = os.Stat(gitInfoExclude)
		if err != nil {
			return
		}

		excludeCurrent, err := os.Open(gitInfoExclude)
		if err != nil {
			return
		}
		defer func(excludeCurrent *os.File) {
			err := excludeCurrent.Close()
			if err != nil {
				log.Fatalf("unable to close %s. %v", gitInfoExclude, err)
			}
		}(excludeCurrent)

		excludeNew, err := os.Create(gitInfoExcludeNew)
		if err != nil {
			log.Fatalf("unable to create %s, %v", gitInfoExcludeNew, err)
		}
		defer func(excludeNew *os.File) {
			err := excludeNew.Close()
			if err != nil {
				log.Fatalf("enable to close %s, %v", gitInfoExcludeNew, err)
			}
		}(excludeNew)

		scanner := bufio.NewScanner(excludeCurrent)
		for scanner.Scan() {
			if scanner.Text() == "###> tmpl ###" {
				for scanner.Scan() {
					if scanner.Text() == "###< tmpl ###" {
						continue
					}
				}
			} else {
				_, err := excludeNew.WriteString(scanner.Text() + "\n")
				if err != nil {
					log.Fatalf("unable to write to %s, %v", gitInfoExcludeNew, err)
				}
			}
		}

		err = os.Rename(excludeNew.Name(), excludeCurrent.Name())
		if err != nil {
			log.Fatalf("unable to rename from %s to %s. %v", excludeCurrent.Name(), excludeNew.Name(), err)
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
	unlinkCmd.Flags().StringP("target", "t", currentDirname(), "template name")
}
