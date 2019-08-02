package main

import (
	shortener "github.com/Djiffit/url-shortener"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

func main() {

	db := shortener.InitDb()
	defer db.Close()

	linkStore := shortener.CreateLinkStore(db)
	userStore := shortener.CreateUserStore(db)
	e := echo.New()
	shortener.CreateLinkRouter(e, linkStore)
	shortener.CreateUserRouter(e, userStore)

	e.Logger.Fatal(e.Start(":1323"))
}
