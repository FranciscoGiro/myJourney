package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID 				primitive.ObjectID 		`json:"_id" bson:"_id"`
	User_id 		string 					`json:"user_id" bson:"user_id"`
	Url 			string 		    		`json:"url" bson:"url,omitempty"`
	City 			string 	        		`json:"city" bson:"city"`
	Country 		string 	        		`json:"country" bson:"country"`
	Lat 			float64 	        	`json:"lat" bson:"lat"`
	Lng 			float64 	        	`json:"lng" bson:"lng"`
	Date	        time.Time	 		    `json:"date" bson:"date"`
	CreatedAt 		time.Time 		        `json:"createdAt" bson:"createdAt"`
}