package models

import (
	"context"
	"log"
	"setuphelper/api/utilities"
	"strings"

	//"github.com/labstack/echo/middleware"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	//UserModel - User Values
	UserModel struct {
		//ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		ID string `json:"_id,omitempty" bson:"_id,omitempty"`
		//UserID       string         `json:"userID"`
		UserName     string         `json:"userName"`
		Password     string         `json:"password"`
		FirstName    string         `json:"firstName"`
		LastName     string         `json:"lastName"`
		SecurityList []UserSecurity `json:"securityList"`
		ContactID    string         `json:"contactID"`

		CreationTime
	}

	//UserSecurty - UserSecurity Levels
	UserSecurity struct {
		Level string `json:"level"`
	}
)

var (
	//local variable use to set the collection name for
	//queries against the user model
	userTableConfig = &TableConfig{"users"}
)

// GetUser - Return the userModel of the Object ID
func GetUser(id primitive.ObjectID) (UserModel, error) {
	// Create a BSON ObjectID by passing string to ObjectIDFromHex() method
	filter := bson.M{"_id": id.String()}

	//console.log("Filter", filter)

	// Create a empty UsermModel we can use as to store the decoded results of the find.
	var userModel UserModel
	//err := utilities.DatabaseObj.FindOne(userTableConfig.Name, filter, &userModel)
	err := utilities.DatabaseObj.FindOne(userTableConfig.Name, filter, &userModel)

	utilities.PrintDebug("Get User", userModel, err)
	return userModel, err
}

// GetUserList - Return the userModel of the Object ID
func GetUserList() []UserModel {
	filter := bson.M{}
	cursor, err := utilities.DatabaseObj.Find(userTableConfig.Name, filter)
	defer cursor.Close(context.Background())

	if err != nil {
		utilities.PrintDebug("Get User List Error", err)
		//return []&UserModel{}
	}

	var userList []UserModel
	for cursor.Next(context.Background()) {
		// To decode into a struct, use cursor.Decode()
		userModel := UserModel{}

		err := cursor.Decode(&userModel)
		if err != nil {
			log.Fatal("GetUserList", err)
		}

		userList = append(userList, userModel)
	}

	return userList
}

// Delete - Save the User
func (m *UserModel) Delete() error {
	log.Print("User Model Delete")

	filter := bson.M{"_id": m.GetID()}
	result := utilities.DatabaseObj.FindOneAndDelete(userTableConfig.Name, filter)

	if result.Err() != nil {
		utilities.PrintDebug("Delete Error", result.Err)
	}

	return result.Err()
}

//GetID - Return the ID string
func (m *UserModel) GetID() string {
	//id, _ := primitive.ObjectIDFromHex(m.ID)
	return m.ID
}

//GetFullName - Return the full name of the user
func (m *UserModel) GetFullName() string {
	return strings.Join([]string{m.FirstName, m.LastName}, " ")
}

//IsAllowedToLogin - Determine if the user is allowed to login
func (m *UserModel) IsAllowedToLogin() bool {
	//id, _ := primitive.ObjectIDFromHex(m.ID)
	filter := bson.M{
		"username": m.UserName,
		"password": m.Password}

	var model UserModel
	err := utilities.DatabaseObj.FindOne(userTableConfig.Name, filter, &model)

	if err != nil {
		utilities.PrintDebug("User not found in DB", filter)
		return false
	}

	return true
}

//SetID - Set the ID
func (m *UserModel) SetID(id primitive.ObjectID) {
	//m.ID = id.Hex()
	m.ID = id.String()
	//m.ID = id
	//m.UserID = id.Hex()
}

// SetContactID - Set the contact ID of the saved model
func (m *UserModel) SetContactID(id string) {
	m.ContactID = id
	//return m.ContactID
}

// Save - Save the User
func (m *UserModel) Save() error {
	log.Print("User Model Save")

	m.CreationTime.UpdateTimes()

	utilities.DatabaseObj.Insert(userTableConfig.Name, m)

	return nil
}

//Update - Save the current model
func (m *UserModel) Update() error {
	//id, _ := primitive.ObjectIDFromHex(m.GetID())

	//filter := bson.M{"_id": id}
	// GetID returns a String so use that string to filter
	filter := bson.M{"_id": m.GetID()}
	m.CreationTime.UpdateTimes()

	update := bson.D{{"$set", m}}

	result := utilities.DatabaseObj.FindOneAndUpdate(userTableConfig.Name, filter, update)

	var userModel UserModel
	result.Decode(&userModel)

	//utilities.PrintDebug("RESULT", result.Err(), filter, userModel)
	utilities.PrintDebug("User Model Update", userModel)
	//utilities.PrintDebug("RESULT bson", update)
	return result.Err()
}
