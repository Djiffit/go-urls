package shortener

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

type User struct {
	AccountID int
	Name      string
}

type UserStore interface {
	GetUser(id int) (User, error)
	CreateUser(name, password string) (string, error)
	GetUsers() ([]User, error)
	Login(name, password string) (string, error)
}

type UserModel struct {
	store UserStore
}

type PostUser struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreateUserStore(db *sql.DB) UserStore {
	return &UserDBStore{db}
}

func (u *UserModel) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, ErrDataMissing)
	}

	user, err := u.store.GetUser(id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UserModel) CreateUser(c echo.Context) error {
	data := new(PostUser)

	err := c.Bind(data)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrDataMissing)
	}

	token, err := u.store.CreateUser(data.Name, data.Password)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrInvalidCredentials)
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"token": token,
	})
}

func (u *UserModel) GetUsers(c echo.Context) error {
	users, err := u.store.GetUsers()

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

func (u *UserModel) Login(c echo.Context) error {
	data := new(PostUser)

	err := c.Bind(data)

	if err != nil {
		return echo.NewHTTPError(401, err)
	}

	token, err := u.store.Login(data.Name, data.Password)

	if err != nil {
		return echo.NewHTTPError(401, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
