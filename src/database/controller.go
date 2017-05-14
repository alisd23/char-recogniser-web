package database

import (
	"errors"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	IMAGE_SIZE = 32
)

type TrainingExample struct {
	Data     []byte
	Size     int
	Charcode int
}

const (
	TRAINING_SET = "training_set"
)

func InsertExample(db *mgo.Database, img []byte, charCode int) error {
	duplicate := TrainingExample{}
	err := db.
		C(TRAINING_SET).
		Find(bson.M{
			"data": img,
		}).
		One(&duplicate)

	// If no error: duplicate was found
	if err == nil {
		return errors.New("Duplicate image found")
	}

	record := TrainingExample{
		Data:     img,
		Size:     IMAGE_SIZE,
		Charcode: charCode,
	}

	err = db.C(TRAINING_SET).Insert(record)

	return err
}
