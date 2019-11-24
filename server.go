package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"log"
	"setuphelper/api/controllers"
	"setuphelper/api/utilities"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func login(c echo.Context) error {
	fmt.Printf("%+v\n", c)
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "jon" || password != "shhh!" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"Jon Snow",
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

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	//client := GetClient()
	//utilities.DatabaseConnect()
	utilities.Init()
	client := utilities.DatabaseObj.GetClient()

	// Check the connection
	if err := client.Ping(context.Background(), nil); err != nil {
		utilities.PrintDebug("Could not Connect to MongoDB!")

	}
	utilities.PrintDebug("Connected MongoDB! -main")

	e := echo.New()

	// Recover middleware recovers from panics anywhere in the chain,
	// prints a stack trach and send control to the HTTPErrorHandler
	e.Use(middleware.Recover())

	// Logger middleware logs the information about each HTTP request
	e.Use(middleware.Logger())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route
	e.GET("/", accessible)
	e.GET("/test", accessible)

	e.Static("/material-dashboard-react/static", "./ui-dist/static")
	e.Static("/dashboard", "./ui-dist/")

	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })

	// Restricted group
	r := e.Group("/restricted")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)

	e.POST("/users", controllers.CreateUser)
	e.GET("/users", controllers.GetUserList)
	e.GET("/users/:id", controllers.GetUser)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	e.Logger.Fatal(e.Start(":3001"))
}
