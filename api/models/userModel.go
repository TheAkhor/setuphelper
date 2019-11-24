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
		ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		UserID       string             `json:"userID"`
		UserName     string             `json:"userName"`
		Password     string             `json:"password"`
		SecurityList []UserSecurity     `json:"securityList"`
		ContactID    string             `json:"contactID"`

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
	filter := bson.M{"_id": id}

	// Create a empty UsermModel we can use as to store the decoded results of the find.
	var userModel UserModel
	err := utilities.DatabaseObj.FindOne(userTableConfig.Name, filter, &userModel)

	utilities.PrintDebug("Get User", userModel, err)
	return userModel, nil
}

func GetUserList() map[string]*UserModel {
	return users
}

//func (m *UserModel) SetID(id string) {
func (m *UserModel) SetID(id primitive.ObjectID) {
	m.ID = id
	m.UserID = id.Hex()
}

func (m *UserModel) GetUserID() string {
	return m.UserID
}

// SetContactID - Set the contact ID of the saved model
func (m *UserModel) SetContactID(id string) {
	m.ContactID = id
	//return m.ContactID
}

// Save - Save the User
func (m *UserModel) Save() error {
	log.Print("User Model Save")

	if m.UserID != "" {
		m.UserID = string(userseq)
		userseq++
	}
	users[string(userseq)] = m

	m.CreationTime.UpdateTimes()

	//utilities.PrintDebug("created At", m.CreatedAt)

	log.Print("Contact User Save", m.UserID, m)
	utilities.PrintDebug(m)

	utilities.DatabaseObj.Insert(userTableConfig.Name, m)

	return nil
}
