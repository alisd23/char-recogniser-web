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
	"char-recogniser-go/src/database"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
)

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

func runSeedTask() {
	db, err := database.Connect(LOCAL_URL)

	if err != nil {
		return
	}

	sourceDirAbs, _ := filepath.Abs("training-set")
	imgPaths, err := filepath.Glob(sourceDirAbs + "/*/*")

	if err != nil || len(imgPaths) == 0 {
		fmt.Println("[INVALID PATHS] Error reading training-set directory or no images found")
		return
	}

	// For each file in directory, process image and save new image in form:
	// training-set/:character:/:index:.png
	for _, imgPath := range imgPaths {
		dirname := filepath.Base(filepath.Dir(imgPath))
		charCode, err := strconv.ParseInt(dirname, 10, 8)

		// CHeck directory charCode value is valid
		if err != nil {
			fmt.Printf("[INVALID DIRECTORY] %v - Expected a char code, received: %v\n", imgPath, charCode)
			continue
		}

		// Read file []byte into variable
		bytes, err := ioutil.ReadFile(imgPath)

		if err != nil {
			fmt.Printf("[READ IMAGE] %v - error: %v\n", imgPath, err)
			continue
		}

		// Insert image into DB
		err = database.InsertExample(db, bytes, int(charCode))

		if err != nil {
			fmt.Printf("[DATABASE INSERT] %v - error: %v\n", imgPath, err)
		} else {
			fmt.Printf("[DATABASE INSERT] Image inserted: %v\n", imgPath)
		}
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
}
