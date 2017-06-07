package database

import (
	"image"
	"image/draw"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	IMAGE_SIZE = 32
)

type TrainingExample struct {
	Data     []float32
	Size     int
	Charcode int
}

const (
	TRAINING_SET = "training_set"
)

func InsertExamples(db *mgo.Database, examples []interface{}) error {
	bulk := db.C(TRAINING_SET).Bulk()
	bulk.Insert(examples...)
	bulk.Unordered()
	_, err := bulk.Run()
	return err
}

func InsertExample(db *mgo.Database, example interface{}) error {
	err := db.C(TRAINING_SET).Insert(example)
	return err
}

func CreateExample(img image.Image, charCode int) interface{} {
	rect := img.Bounds()
	rawImg := image.NewGray(rect)
	draw.Draw(rawImg, rect, img, rect.Min, draw.Src)
	pixels := rawImg.Pix

	example := make([]float32, len(pixels))

	// Convert into array of intensities - between 0 & 1
	for i, pixel := range pixels {
		example[i] = float32(pixel) / 255
	}

	// duplicate := TrainingExample{}
	// err := db.
	// 	C(TRAINING_SET).
	// 	Find(bson.M{
	// 		"data": example,
	// 	}).
	// 	One(&duplicate)

	// If no error: duplicate was found
	// if err == nil {
	// 	return errors.New("Duplicate image found")
	// }

	record := bson.M{
		"data":     example,
		"size":     IMAGE_SIZE,
		"charcode": charCode,
	}

	return record
}
