package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var SECRET = []byte("secret")

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "admin" && password == "admin123!" {
		// Create token
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "John Snow"
		claims["role"] = "admin"
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

		// Generate encoded token and send it as response
		t, err := token.SignedString(SECRET)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func HealthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]bool{
		"alive": true,
	})
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	role := claims["role"].(string)
	return c.String(http.StatusOK, "Welcome"+name+"-"+role+"!")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	e.Static("/static", "static");

	// Unauthenticated route
	e.GET("/", accessible)
	//e.GET("/health-check", HealthCheckHandler)
	e.GET("/health-check", HealthCheckHandler)

	// Restricted group
	r := e.Group("/restricted")
	r.Use(middleware.JWT(SECRET))
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":8086"))
}
