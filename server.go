package main

import (
	"context"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"setuphelper/api/controllers"
	"setuphelper/api/utilities"
)

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func main() {
	//client := GetClient()
	//utilities.DatabaseConnect()

	// command := exec.Command("npm", "start", "--prefix ./ui/")
	// command.Stdout = os.Stdout
	// command.Stderr = os.Stderr
	// if err := command.Run(); err != nil {
	// 	log.Fatal(err)
	// }

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
