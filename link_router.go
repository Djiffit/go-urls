package shortener

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func CreateLinkRouter(e *echo.Echo, store LinkStore) {
	link := LinkModel{store}
	g := e.Group("/links")

	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte(secret),
		Claims:     &jwtClaims{},
	}))

	g.POST("/", link.PostLink)
	g.GET("/", link.GetLinks)
	g.DELETE("/:id", link.DeleteLink)
	g.GET("/:id", link.GetLink)
}
