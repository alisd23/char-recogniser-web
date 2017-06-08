package database

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

var database = "char-recogniser"

// Connect connects to the MongoDB database
func Connect(url string) (*mgo.Database, error) {
	session, err := mgo.Dial(url)

	if err != nil {
		fmt.Println("[MONGO] Connection error: ", err)
		return nil, err
	}

	db := session.DB(database)
	return db, nil
}
