package controllers

import (
	"log"
	"net/http"

	"setuphelper/api/models"
	"setuphelper/api/utilities"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	//"github.com/labstack/echo/middleware"
)

type (
	// createUserRequestModel is used to accept the create user json
	createUserRequestModel struct {
		User    models.UserModel    `json:user`
		Contact models.ContactModel `json:"contact"`
	}

	// UserController is used for CRUD operations
	UserController struct{}
)

//----------
// Handlers
//----------

// CreateUser - Create a user
func (controller *UserController) CreateUser(c echo.Context) error {
	log.Print("Create User")

	// Create a empty userModel we can use to Decode context into (bind function)
	model := &createUserRequestModel{}

	//Decode the incoming data into the model
	if err := c.Bind(model); err != nil {
		log.Print("Bind Issue", err)
		return c.JSON(http.StatusNoContent, err)
	}

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

	return c.JSON(http.StatusCreated, model.User)
}

// GetUser - Get a user
func (controller *UserController) GetUser(c echo.Context) error {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	userModel, err := models.GetUser(id)

	utilities.PrintDebug("GetUser", id, userModel)
	if err != nil {
		return c.JSON(http.StatusNoContent, err)
	}

	return c.JSON(http.StatusOK, userModel)
}

// GetUserList - List all users
func (controller *UserController) GetUserList(c echo.Context) error {
	log.Print("GetUserList")
	return c.JSON(http.StatusOK, models.GetUserList())
}

// UpdateUser - Update a user
func (controller *UserController) UpdateUser(c echo.Context) error {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	utilities.PrintDebug("Update User ID", id)

	//var u models.UserModelTest
	var u models.UserModel

	if err := c.Bind(&u); err != nil {
		return err
	}

	model, err := models.GetUser(id)

	if err != nil {
		utilities.PrintDebug("Error")
		return c.JSON(http.StatusOK, "Not Found")
	}

	model.Password = u.Password

	utilities.PrintDebug("ModelID:", model.GetID())

	err = model.Update()

	if err != nil {
		return c.JSON(http.StatusOK, "Update User Error")
	}

	// Don't really need to do this just good for testing
	model, _ = models.GetUser(id)

	//utilities.PrintDebug("ModelTest", u, err)
	return c.JSON(http.StatusOK, model)
}

// DeleteUser - Remove a user
func (controller *UserController) DeleteUser(c echo.Context) error {

	id, _ := primitive.ObjectIDFromHex(c.Param("id"))

	var model models.UserModel
	model, _ = models.GetUser(id)
	model.Delete()

	return c.NoContent(http.StatusNoContent)
}
