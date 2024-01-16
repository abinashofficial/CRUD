package model

import "go.mongodb.org/mongo-driver/mongo"

var StudentInfo = []TestCase{{Name: "Abi", MobNumber: 994048389}, {Name: "Arun", MobNumber: 987436652}}

type TestCase struct {
	Name      string
	MobNumber int
}

type Requests struct {
	Name string `json:"Name"`
}

type OpenApplication struct {
	AppName        string              `json:"app_name" bson:"app_name"`
	AppKey         string              `json:"app_key" bson:"app_key"`
	AppSecret      string              `json:"app_secret" bson:"app_secret"`
	Brand          string              `json:"brand" bson:"brand"`
	Description    string              `json:"description" bson:"description"`
	PermissionsMap map[string][]string `json:"permissions_map" bson:"permissions_map"`
	Email          string              `json:"email" bson:"email"`
	Source         string              `json:"source" bson:"source"`
}

type MongoClient struct {
	PrimaryClient   *mongo.Client
	SecondaryClient *mongo.Client
}

type MongoError struct {
	PrimaryClientError   error
	SecondaryClientError error
}
