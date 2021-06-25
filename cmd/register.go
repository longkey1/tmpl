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
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
)

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
		templateDirPath := filepath.Join(config.TemplateDir, templateName)
		if _, err := os.Stat(templateDirPath); os.IsNotExist(err) {
			err := os.MkdirAll(templateDirPath, os.ModePerm)
			if err != nil {
				log.Fatalf("unable to make %s directory", templateDirPath)
			}
		}
		templateFilePath := filepath.Join(templateDirPath, originalFileName)
		if _, err := os.Stat(templateFilePath); !os.IsNotExist(err) {
			log.Fatalf("already %s exists", templateFilePath)
		}

		originalFilePath := filepath.Join(wd, originalFileName)
		err = filepath.Walk(originalFilePath, func(orgPath string, info fs.FileInfo, err error) error {
			rel, err := filepath.Rel(originalFilePath, orgPath)
			if err != nil {
				return err
			}
			tmpPath := filepath.Join(templateFilePath, rel)
			if info.IsDir() {
				err := os.MkdirAll(tmpPath, os.ModePerm)
				if err != nil {
					return err
				}
			} else {
				src, err := os.Open(orgPath)
				if err != nil {
					return err
				}
				defer func(src *os.File) {
					err := src.Close()
					if err != nil {
						log.Fatalf("unable to close %s for copy", src.Name())
					}
				}(src)

				dst, err := os.Create(tmpPath)
				if err != nil {
					return err
				}
				defer func(dst *os.File) {
					err := dst.Close()
					if err != nil {
						log.Fatalf("unable to close %s for copy", src.Name())
					}
				}(dst)

				_, err = io.Copy(dst, src)
				if  err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("unable to copy %s, %v", originalFilePath, err)
		}

		_, _ = fmt.Printf("register %s template to %s\n", templateName, templateFilePath)
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
