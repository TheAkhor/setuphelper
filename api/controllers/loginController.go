package controllers

import (
	"log"
	"net/http"
	"time"

	"setuphelper/api/models"
	"setuphelper/api/utilities"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
)

type (
	LoginController struct {
	}

	JwtCustomClaims struct {
		Name  string `json:"name"`
		Admin bool   `json:"admin"`
		jwt.StandardClaims
	}
)

//----------
// Handlers
//----------
func (controller *LoginController) Login(c echo.Context) error {
	userModel := &models.UserModel{}

	if err := c.Bind(userModel); err != nil {
		log.Print("Login Bind Issue", err)
		return c.JSON(http.StatusNoContent, err)
	}

	utilities.PrintDebug("JWT Login Controller	", userModel)

	// Throws unauthorized error
	if !userModel.IsAllowedToLogin() {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &JwtCustomClaims{
		userModel.GetFullName,
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
