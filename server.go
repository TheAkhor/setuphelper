package main

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"log"
	"setuphelper/api/controllers"
	"setuphelper/api/models"
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

	userModel := &models.UserModel{}

	if err := c.Bind(userModel); err != nil {
		log.Print("Login Bind Issue", err)
		return c.JSON(http.StatusNoContent, err)
	}

	utilities.PrintDebug("JWT Login", userModel)

	// Throws unauthorized error
	if !userModel.IsAllowedToLogin() {
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

// type CustomBinder struct{}

// func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
// 	// You may use default binder
// 	db := new(echo.DefaultBinder)
// 	utilities.PrintDebug("Custom Binder")
// 	if err = db.Bind(i, c); err != echo.ErrUnsupportedMediaType {
// 		utilities.PrintDebug("Custom Binder Err", err, err.Error())
// 		utilities.PrintDebug(err.Error())
// 		utilities.PrintDebug(c)
// 		fmt.Println(c)
// 		return
// 	}

// 	utilities.PrintDebug("Custom Binder Fail", i)

// 	// Define your custom implementation

// 	return
// }

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
	//e.Binder = &CustomBinder{}

	// Recover middleware recovers from panics anywhere in the chain,
	// prints a stack trach and send control to the HTTPErrorHandler
	e.Use(middleware.Recover())

	// Logger middleware logs the information about each HTTP request
	e.Use(middleware.Logger())

	// Login route for JWT
	//e.POST("/login", login)
	loginController := controllers.LoginController{}
	e.POST("/login", loginController.Login)

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
		Claims:     &controllers.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)

	userController := controllers.UserController{}
	r.POST("/users", userController.CreateUser)
	r.GET("/users", userController.GetUserList)
	r.GET("/users/:id", userController.GetUser)
	r.PUT("/users/:id", userController.UpdateUser)
	r.DELETE("/users/:id", userController.DeleteUser)

	//Remove these later - Just for easy postman testing
	e.POST("/users", userController.CreateUser)
	e.GET("/users", userController.GetUserList)
	e.GET("/users/:id", userController.GetUser)
	e.PUT("/users/:id", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)

	e.Logger.Fatal(e.Start(":3001"))
}
