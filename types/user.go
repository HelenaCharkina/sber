package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	Name       string             `json:"name"`
	Job        string             `json:"job"`
	EmployedAt string             `json:"employed_at"`
	Level      int                `json:"-" bson:"level"`
	Employees  []User             `json:"employees" bson:"result"`
}
