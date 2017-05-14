// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"image/png"
	"io/ioutil"
	"letter-recogniser-go/src/server"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var sourceDir string
var outputDir string

// processImagesCmd represents the processImages command
var processImagesCmd = &cobra.Command{
	Use:   "processImages",
	Short: "Command for normalising character training set images",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO - Probably want to save these images in []byte form to mongdb/cassandra
		// or something - so it's easy to retrieve later
		// NOTE - Need to MAKE SURE images are not saved to DB more than once!
		fmt.Println("[PROCESSING IMAGES] - Directory: ", sourceDir)
		runTask()
	},
}

func processImage(filePath, character, name string) {
	imgBytes, err := ioutil.ReadFile(filePath)

	if err != nil {
		fmt.Printf("[READ IMAGE] %v/%v - Error: %v\n", character, name, err)
		return
	}

	reader := bytes.NewReader(imgBytes)
	img, err := png.Decode(reader)

	if err != nil {
		fmt.Printf("[DECODE IMAGE] %v/%v - Error: %v\n", character, name, err)
		return
	}

	// Normalise image
	normalisedImg := server.NormaliseImage(img)
	outputPath, _ := filepath.Abs(filepath.Join("training-set", character, name+".png"))

	// Create directory if necessary - then save image
	os.Mkdir(filepath.Dir(outputPath), 0777)
	err = server.SaveImage(normalisedImg, outputPath)

	if err != nil {
		fmt.Printf("[SAVE IMAGE] %v/%v - Error: %v\n", character, name, err)
	} else {
		fmt.Printf("[IMAGE PROCESSED] %v\n", outputPath)
	}
}

func runTask() {
	sourceDirAbs, _ := filepath.Abs(sourceDir)
	directories, err := filepath.Glob(sourceDirAbs + "/*")

	if err != nil {
		fmt.Println("[READ DIRECTORIES] - error: ", err)
		return
	}

	// Get all directory (character) names
	for _, path := range directories {
		dirname := filepath.Base(path)
		characterCode, _ := strconv.ParseInt(dirname, 10, 8)

		// Search each character directory to get all image paths
		imgPaths, err := filepath.Glob(path + "/*")

		if err != nil {
			fmt.Println("[INVALID DIRECTORY] - ", dirname)
			continue
		}

		// For each file in directory, process image and save new image in form:
		// training-set/:character:/:index:.png
		for index, imgPath := range imgPaths {
			processImage(
				imgPath,
				strconv.FormatInt(int64(characterCode), 10),
				strconv.FormatInt(int64(index), 10),
			)
		}
	}
}

func init() {
	RootCmd.AddCommand(processImagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// processImagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processImagesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	processImagesCmd.
		Flags().
		StringVarP(&sourceDir, "sourceDir", "s", "", "Directory of source images for processing")

	processImagesCmd.
		Flags().
		StringVarP(&outputDir, "outputDir", "o", "", "Directory where process images are saved")
}
