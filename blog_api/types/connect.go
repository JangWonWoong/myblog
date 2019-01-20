package types

import (
	"gopkg.in/mgo.v2"
)

func dbConnect() (*mgo.Session, error) {
	db, err := mgo.Dial("mongodb://localhost:27017/")
	if err != nil {
		return nil, err
	}
	return db, nil
}