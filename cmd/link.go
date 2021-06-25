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

	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
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
			err := os.Symlink(fmt.Sprintf("%s/%s", td, file.Name()), fmt.Sprintf("./%s", file.Name()))
			if err == nil {
				_, _ = fmt.Printf("created %s symlink\n", file.Name())
			} else {
				_, _ = fmt.Printf("not created %s symlink. %s\n", file.Name(), err)
			}
		}

		_, err = os.Stat(".git/info/exclude")
		if err != nil {
			return
		}

		excludeCurrent, err :=  os.Open(".git/info/exclude")
		if err != nil {
			return
		}
		defer func(excludeCurrent *os.File) {
			err := excludeCurrent.Close()
			if err != nil {
				log.Fatalf("unable to close .git/info/exclude. %v" ,err)
			}
		}(excludeCurrent)

		excludeNew, err := os.Create(".git/info/exclude.new")
		if err != nil {
			log.Fatalf("unable to create .git/info/exclude.new. %v" ,err)
		}
		defer func(excludeNew *os.File) {
			err := excludeNew.Close()
			if err != nil {
				log.Fatalf("enable to close .git/info/exclude.new. %v" ,err)
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
				_, err := excludeNew.WriteString(scanner.Text()+"\n")
				if err != nil {
					log.Fatalf("unable to write to .git/info/exclude.new. %v" ,err)
				}
			}
		}

		if len(files) > 0 {
			_, err = excludeNew.WriteString("###> tmpl ###\n")
			if err != nil {
				log.Fatalf("unable to write start line to .git/info/exclude.new. %v" ,err)
			}
			for _, file := range files {
				_, err = excludeNew.WriteString(fmt.Sprintf("/%s\n", file.Name()))
				if err != nil {
					log.Fatalf("unable to write start line to .git/info/exclude.new. %v" ,err)
				}
			}
			_, err = excludeNew.WriteString("###< tmpl ###\n")
			if err != nil {
				log.Fatalf("unable to write end line to .git/info/exclude.new. %v" ,err)
			}
		}

		err = os.Rename(excludeNew.Name(), excludeCurrent.Name())
		if err != nil {
			log.Fatalf("unable to rename from %s to %s. %v" ,excludeCurrent.Name(), excludeNew.Name(), err)
		}
	},
}

func init() {
	rootCmd.AddCommand(linkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
