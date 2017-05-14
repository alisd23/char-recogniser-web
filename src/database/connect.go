package database

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

var DATABASE = "char-recogniser"

func Connect(url string) (*mgo.Database, error) {
	session, err := mgo.Dial(url)

	if err != nil {
		fmt.Println("[MONGO] Connection error: ", err)
		return nil, err
	}

	db := session.DB(DATABASE)
	return db, nil
}
