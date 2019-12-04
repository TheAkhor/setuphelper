package models

import (
	"log"
	"setuphelper/api/utilities"

	//"github.com/labstack/echo/middleware"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	UserModel struct {
		//ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		ID string `json:"_id,omitempty" bson:"_id,omitempty"`
		//UserID       string         `json:"userID"`
		UserName     string         `json:"userName"`
		Password     string         `json:"password"`
		SecurityList []UserSecurity `json:"securityList"`
		ContactID    string         `json:"contactID"`

		CreationTime
	}

	UserModelTest struct {
		ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		UserID       string             `json:"userID"`
		UserName     string             `json:"userName"`
		Password     string             `json:"password"`
		SecurityList []UserSecurity     `json:"securityList"`
		ContactID    string             `json:"contactID"`
	}

	UserSecurity struct {
		Level string `json:"level"`
	}
)

var (
	users           = map[string]*UserModel{}
	userseq         = 1
	userTableConfig = &TableConfig{"users"}
)

func DeleteUser(id string) {
	delete(users, id)
}

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

func GetUserList() map[string]*UserModel {
	return users
}

//func (m *UserModel) SetID(id string) {
func (m *UserModel) SetID(id primitive.ObjectID) {
	//m.ID = id.Hex()
	m.ID = id.String()
	//m.ID = id
	//m.UserID = id.Hex()
}

func (m *UserModel) GetID() string {
	//id, _ := primitive.ObjectIDFromHex(m.ID)
	return m.ID
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

//func (m *UserModel) Update(update bson.D) error {
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

func (m *UserModel) UpdateTest() error {
	filter := bson.M{"_id": m.GetID()}

	//m.CreationTime.UpdateTimes()

	//result := utilities.DatabaseObj.FindOneAndUpdate(userTableConfig.Name, filter, update)
	result := utilities.DatabaseObj.FindOneAndReplace(userTableConfig.Name, filter, m)

	if result.Err() != nil {
		return result.Err()
	}

	model := &UserModel{}

	utilities.PrintDebug("Decode", model)

	if err := result.Decode(model); err != nil {
		utilities.PrintDebug("Decode Issue", err.Error())

	}

	//result.Decode(m)

	return nil
}
