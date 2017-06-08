package server

import (
	"image"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	trainingSet = "training_set"
	// ImageSize is the size of the image
	ImageSize = 32
)

// InsertExamples inserts the given examples into the database in one bluk
// operation
func InsertExamples(db *mgo.Database, examples []interface{}) error {
	bulk := db.C(trainingSet).Bulk()
	bulk.Insert(examples...)
	bulk.Unordered()
	_, err := bulk.Run()
	return err
}

// InsertExample simply inserts a given exmaple in the "training_set"
// collection in the database
func InsertExample(db *mgo.Database, example interface{}) error {
	err := db.C(trainingSet).Insert(example)
	return err
}

// CreateExample creates a MongoDB training example record
func CreateExample(img image.Image, charCode int) interface{} {
	record := bson.M{
		"data":     GetPixels(img),
		"charcode": charCode,
	}

	return record
}
