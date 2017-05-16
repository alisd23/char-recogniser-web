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
	"char-recogniser-go/src/server"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

var sourceDir string

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
		runProcessImagesTask()
	},
}

func processImage(filePath, character, name string, done chan bool) {
	imgBytes, err := ioutil.ReadFile(filePath)

	// Send true down done channel when function returns
	notifyDone := func() {
		done <- true
	}
	defer notifyDone()

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

func runProcessImagesTask() {
	sourceDirAbs, _ := filepath.Abs(sourceDir)
	imgPaths, err := filepath.Glob(sourceDirAbs + "/*/*")

	if err != nil {
		fmt.Println("[INVALID PATHS]")
		return
	}

	imgCounts := map[string]int{}
	goroutines := 0
	done := make(chan bool)

	// For each file in directory, process image and save new image in form:
	// training-set/:character:/:index:.png
	for _, imgPath := range imgPaths {
		charCode := filepath.Base(filepath.Dir(imgPath))
		_, err := strconv.ParseInt(charCode, 10, 8)

		if err != nil {
			fmt.Println("[INVALID DIRECTORY] Expected a char code, received: ", charCode)
			continue
		}

		go processImage(
			imgPath,
			charCode,
			strconv.FormatInt(int64(imgCounts[charCode]), 10),
			done,
		)

		imgCounts[charCode]++
		goroutines++
	}

	for i := 0; i < goroutines; i++ {
		<-done
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
}
