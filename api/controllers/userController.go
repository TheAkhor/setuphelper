package controllers

import (
	"fmt"
	"log"
	"net/http"

	"setuphelper/api/models"
	"setuphelper/api/utilities"

	"github.com/labstack/echo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/labstack/echo/middleware"
)

type (
	controllerModel struct {
		User    models.UserModel    `json:user`
		Contact models.ContactModel `json:"contact"`
	}
	Test struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
)

//----------
// Handlers
//----------

// CreateUser - Create a user
func CreateUser(c echo.Context) error {
	log.Print("Create User")

	// Create a empty userModel we can use to Decode context into (bind function)
	model := &controllerModel{}

	//Decode the incoming data into the model
	if err := c.Bind(model); err != nil {
		log.Print("Bind Issue", err)
		return c.JSON(http.StatusNoContent, err)
	}

	log.Print("user", fmt.Sprintf("%#v", model))
	log.Print("contact", fmt.Sprintf("%#v", model.Contact))

	// Save the Contact Model first so we can set the contactID in the user
	if err := model.Contact.Save(); err != nil {
		utilities.PrintDebug("Contact Save Error")
		utilities.PrintDebug(err)
		return c.JSON(http.StatusNoContent, err)
	}

	// Set the contactID in the user
	model.User.SetContactID(model.Contact.GetID())

	//Save the user model
	if err := model.User.Save(); err != nil {
		utilities.PrintDebug("User Save Error")
		utilities.PrintDebug(err)
		return c.JSON(http.StatusNoContent, err)
	}

	utilities.PrintDebug("Create User Success")
	utilities.PrintDebug(model)

	// //users[u.ID] = u
	// //seq++
	return c.JSON(http.StatusCreated, model.User)
}

// GetUser - Get a user
func GetUser(c echo.Context) error {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	utilities.PrintDebug("GetUser", id)

	userModel, err := models.GetUser(id)

	if err != nil {
		return c.JSON(http.StatusNoContent, err)
	}

	return c.JSON(http.StatusOK, userModel)
}

// GetUserList - List all users
func GetUserList(c echo.Context) error {
	log.Print("GetUserList")
	return c.JSON(http.StatusOK, models.GetUserList())
}

// UpdateUser - Update a user
func UpdateUser(c echo.Context) error {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	utilities.PrintDebug("Update User ID", id)

	//var u models.UserModel
	var u models.UserModelTest

	utilities.PrintDebug("ModelTest", u)

	if err := c.Bind(&u); err != nil {
		return err
	}
	//id, _ := strconv.Atoi(c.Param("id"))
	//id := c.Param("id")

	model, err := models.GetUser(id)
	//model.Password = u.Password

	utilities.PrintDebug("ModelTest", u, err)
	return c.JSON(http.StatusOK, model)
}

// DeleteUser - Remove a user
func DeleteUser(c echo.Context) error {
	//id, _ := strconv.Atoi(c.Param("id"))
	id := c.Param("id")
	models.DeleteUser(id)
	return c.NoContent(http.StatusNoContent)
}
