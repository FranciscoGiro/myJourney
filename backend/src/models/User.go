package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID 				primitive.ObjectID 		`json:"_id" bson:"_id,omitempty"`
	Name 			string 					`json:"name" bson:"name"`
	Password 		string 		    		`json:"password" bson:"password"`
	Email 			string 	        		`json:"email" bson:"email"`
	Role			string	 				`json:"role" bson:"role"`
	CreatedAt 		time.Time 		        `json:"createdAt" bson:"createdAt"`
	UpdatedAt 		time.Time 		        `json:"updatedAt" bson:"updatedAt"`
}