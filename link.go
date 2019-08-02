package shortener

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo"
)

type LinkStore interface {
	GetLink(id Identifier, ip_address string, user_id int) (Link, error)
	SaveLink(data *LinkPost, accountID int) error
	DeleteLink(id Identifier) error
	GetLinks(limit, offset int, orderBy string) ([]LinkElement, error)
}

type Link = string

type Identifier = string

type LinkPost struct {
	ID     string `json:"id"`
	Target string `json:"target"`
}

type GetLinksParams struct {
	limit   int
	offset  int
	orderBy string
}

type LinkModel struct {
	Store LinkStore
}

func (l *LinkModel) GetLink(c echo.Context) error {
	id := c.Param("id")

	claims := getTokenData(c)

	if id == "" {
		return ErrIDMissing
	}

	link, err := l.Store.GetLink(id, c.RealIP(), claims.AccountID)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, link)
}

func (l *LinkModel) GetLinks(c echo.Context) error {
	orderBy := "visit_count"
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))

	if err != nil {
		offset = 0
	}

	links, err := l.Store.GetLinks(limit, offset, orderBy)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, links)
}

func (l *LinkModel) DeleteLink(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return ErrIDMissing
	}

	err := l.Store.DeleteLink(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("Link with id %q was succesfully deleted", id))
}

func (l *LinkModel) PostLink(c echo.Context) error {
	data := new(LinkPost)

	err := c.Bind(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrDataMissing.Error()+" "+err.Error())
	}

	claims := getTokenData(c)

	err = l.Store.SaveLink(data, claims.AccountID)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, data)
}

func getTokenData(c echo.Context) *jwtClaims {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtClaims)

	return claims
}
