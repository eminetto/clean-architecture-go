package entity

import (
	"gopkg.in/mgo.v2/bson"
)

//ID type
type ID bson.ObjectId

//ToString convert an ID in a string
func (i ID) String() string {
	return bson.ObjectId(i).Hex()
}

// MarshalJSON will marshal ID to Json
func (i ID) MarshalJSON() ([]byte, error) {
	return bson.ObjectId(i).MarshalJSON()
}

// UnmarshalJSON will convert a string to an ID
func (i *ID) UnmarshalJSON(data []byte) error {
	s := string(data)
	s = s[1 : len(s)-1]
	if bson.IsObjectIdHex(s) {
		*i = ID(bson.ObjectIdHex(s))
	}

	return nil
}

// GetBSON implements bson.Getter.
func (i ID) GetBSON() (interface{}, error) {
	if i == "" {
		return "", nil
	}
	return bson.ObjectId(i), nil
}

// SetBSON implements bson.Setter.
func (i *ID) SetBSON(raw bson.Raw) error {
	decoded := new(string)
	bsonErr := raw.Unmarshal(decoded)
	if bsonErr == nil {
		*i = ID(bson.ObjectId(*decoded))
		return nil
	}
	return bsonErr
}

//StringToID convert a string to an ID
func StringToID(s string) ID {
	return ID(bson.ObjectIdHex(s))
}

//IsValidID check if is a valid ID
func IsValidID(s string) bool {
	return bson.IsObjectIdHex(s)
}

//NewID create a new id
func NewID() ID {
	return StringToID(bson.NewObjectId().Hex())
}
