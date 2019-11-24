package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
	"setuphelper/api/models"
)

// type (
// 	//Contact - Structure of a persons contact details
// 	//Contact struct {
// )

var (
	//contacts   = map[int]*Contact{}
	contacts   = map[int]*models.ContactModel{}
	contactseq = 1
)

//----------
// Handlers
//----------

// CreateContact - Create a contact
func CreateContact(c echo.Context) error {
	return nil
}

// GetContact - Get a contact
func GetContact(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	return c.JSON(http.StatusOK, contacts[id])
}

// GetContactList - List all contact
func GetContactList(c echo.Context) error {
	log.Print("GetContactList")
	return c.JSON(http.StatusOK, contacts)
}

// UpdateContact - Update a contact
func UpdateContact(c echo.Context) error {
	u := new(models.ContactModel)
	if err := c.Bind(u); err != nil {
		return err
	}
	id, _ := strconv.Atoi(c.Param("id"))
	// name := u.FirstName
	// contacts[id].FirstName = name

	return c.JSON(http.StatusOK, contacts[id])
}

// DeleteContact - Remove a contact
func DeleteContact(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	delete(contacts, id)
	return c.NoContent(http.StatusNoContent)
}
