package models

import (
	"log"
	"setuphelper/api/utilities"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (

	//ContactModel - Contact Struct
	ContactModel struct {
		//ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
		ID        string `json:"_id,omitempty" bson:"_id,omitempty"`
		ContactID string `json:"contactID"`
		Type      string `json:"type"`
		Active    bool   `json:"active"`
		CreationTime

		Name      string `json:"name"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Address   string `json:"address"`
		City      string `json:"city"`
		Postal    string `json:"postal"`
		Country   string `json:"country"`

		Birthdate string `json:"birthdate"`

		PhoneList  []PhoneNumber `json:"phoneList"`
		EmailLists []Email       `json:"emailList"`
		PhotoList  []Photo       `json:"photoList"`
	}

	//PhoneNumber - Phone Number Information
	PhoneNumber struct {
		Type   string `json:"type"`
		Number string `json:"number"`
	}

	//Email - Email Information
	Email struct {
		Email string `json:"email"`
	}

	//Photo - Photo Information
	Photo struct {
		File string `json:"file"`
	}
)

var (
	contactTableConfig = &TableConfig{"contacts"}
)

// GetID - Return the unique ID
//func (m *ContactModel) SetID(id string) {
func (m *ContactModel) SetID(id primitive.ObjectID) {
	m.ID = id.String()
}

// GetID - Return the unique ID
func (m *ContactModel) GetID() string {
	utilities.PrintDebug("Contact Model GET ID", m.ID, m)

	return m.ID
}

// Save - Check the incoming id to see if the model already exists
// Update or create new based on this data.
func (m *ContactModel) Save() error {
	log.Print("COntact SAVE----")

	m.CreationTime.UpdateTimes()

	//utilities.PrintDebug("created At", m.CreatedAt)

	log.Print("Contact Model Save", m.ContactID, m)
	utilities.PrintDebug(m)

	utilities.DatabaseObj.Insert(contactTableConfig.Name, m)
	return nil
}
