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
		originalFileName := args[0]
		templateName := path.Base(wd)
		templateDirPath := fmt.Sprintf("%s/%s", config.TemplateDir, templateName)
		if _, err := os.Stat(templateDirPath); os.IsNotExist(err) {
			err := os.MkdirAll(templateDirPath, os.ModePerm)
			if err != nil {
				log.Fatalf("unable to make %s directory", templateDirPath)
			}
		}
		templateFilePath := fmt.Sprintf("%s/%s", templateDirPath, originalFileName)
		if _, err := os.Stat(templateFilePath); !os.IsNotExist(err) {
			log.Fatalf("already %s exists", templateFilePath)
		}

		originalFilePath := fmt.Sprintf("%s/%s", wd, originalFileName)
		err = os.Rename(originalFilePath, templateFilePath)
		if _, err := os.Stat(templateFilePath); os.IsNotExist(err) {
			log.Fatalf("unable to move from %s to %s", originalFilePath, templateFilePath)
		}

		err = os.Symlink(templateDirPath, fmt.Sprintf("./%s", originalFileName))
		if err != nil {
			log.Fatalf("unable to symlink from %s to %s", templateFilePath, fmt.Sprintf("./%s", originalFileName))
		}

		_, _ = fmt.Printf("register %s template to %s, and symlink\n", templateName, templateFilePath)
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
