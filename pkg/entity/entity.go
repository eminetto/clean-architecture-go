package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ID type
type ID = primitive.ObjectID

//StringToID convert a string to an ID
func StringToID(s string) ID {
	res, err := primitive.ObjectIDFromHex(s)
	if err != nil {
		return primitive.NilObjectID
	}
	return res
}

//NewEmptyID return an empty ID
func NewEmptyID() ID {
	return primitive.NilObjectID
}

//IsValidID check if is a valid ID
func IsValidID(s string) bool {
	_, err := primitive.ObjectIDFromHex(s)
	if err == nil {
		return true
	}
	return false
}

//NewID create a new id
func NewID() ID {
	return primitive.NewObjectID()
}
