package shortener

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func CreateUserRouter(e *echo.Echo, store UserStore) {
	user := UserModel{store}
	g := e.Group("/users")
	middleware := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
		Claims:     &jwtClaims{},
	})

	g.POST("/", user.CreateUser)
	g.GET("/", user.GetUsers, middleware)
	g.POST("/login", user.Login)
	g.GET("/:id", user.GetUser, middleware)
}
