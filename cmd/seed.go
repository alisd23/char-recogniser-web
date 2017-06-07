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
	"char-recogniser-go/src/database"
	"char-recogniser-go/src/server"
	"fmt"
	"image/png"
	"io/ioutil"
	"path/filepath"
	"strconv"

	shuffle "github.com/shogo82148/go-shuffle"
	mgo "gopkg.in/mgo.v2"

	"github.com/spf13/cobra"
)

var seedSourceDir string

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with training examples in pure []byte form",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runSeedTask()
	},
}

const LOCAL_URL = "localhost:27017"

func createImageRecord(path string, db *mgo.Database) interface{} {
	dirname := filepath.Base(filepath.Dir(path))
	charCode, err := strconv.ParseInt(dirname, 10, 8)

	// Check directory charCode value is valid
	if err != nil {
		fmt.Printf("[INVALID DIRECTORY] %v - Expected a char code, received: %v\n", path, charCode)
		return nil
	}

	// Read file []byte into variable
	bytesArray, _ := ioutil.ReadFile(path)

	reader := bytes.NewReader(bytesArray)
	img, err := png.Decode(reader)

	if err != nil {
		fmt.Printf("[READ IMAGE] %v - error: %v\n", path, err)
		return nil
	}

	normalisedImage := server.NormaliseImage(img)

	// Insert image into DB
	return database.CreateExample(normalisedImage, int(charCode))
}

func seedBatch(db *mgo.Database, paths []string, batchNo int, done chan<- bool) {
	records := make([]interface{}, 0)
	for _, imgPath := range paths {
		record := createImageRecord(imgPath, db)
		if record != nil {
			records = append(records, record)
		}
	}
	database.InsertExamples(db, records)
	fmt.Printf("Chunk %v processed\n", batchNo)
	done <- true
}

func runSeedTask() {
	db, err := database.Connect(LOCAL_URL)

	if err != nil {
		return
	}

	sourceDirAbs, _ := filepath.Abs(seedSourceDir)
	imgPaths, err := filepath.Glob(sourceDirAbs + "/*/*")

	// Shuffle so we insert in a random order
	shuffle.Strings(imgPaths)

	if err != nil || len(imgPaths) == 0 {
		fmt.Println("[INVALID PATHS] Error reading training-set directory or no images found")
		return
	}

	done := make(chan bool)

	// For each file in directory, process image and save new image in form:
	// training-set/:character:/:index:.png
	noOfChunks := 150
	count := len(imgPaths)
	chunkSize := count / noOfChunks
	remainder := count % chunkSize

	fmt.Printf("No of chunks: %v\n", noOfChunks)
	fmt.Printf("Chunk size: %v\n", chunkSize)

	for i := 0; i < noOfChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		go seedBatch(db, imgPaths[start:end], i+1, done)
	}
	if remainder > 0 {
		start := noOfChunks * chunkSize
		end := len(imgPaths)
		go seedBatch(db, imgPaths[start:end], noOfChunks+1, done)
		noOfChunks++
	}

	for i := 0; i < noOfChunks; i++ {
		<-done
	}
}

func init() {
	RootCmd.AddCommand(seedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	seedCmd.
		Flags().
		StringVarP(&seedSourceDir, "sourceDir", "s", "", "Directory of source images for processing")
}
