package shortener

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const secret = "super_secret"

type UserDBStore struct {
	db *sql.DB
}

type jwtClaims struct {
	Name      string `json:"name"`
	AccountID int    `json:"id"`
	jwt.StandardClaims
}

func (u *UserDBStore) GetUser(id int) (User, error) {
	getQuery := `SELECT id, name FROM account WHERE id = $1;`
	var (
		ID   int
		Name string
	)
	err := u.db.QueryRow(getQuery, id).Scan(&ID, &Name)

	if err != nil {
		return User{}, err
	}

	return User{ID, Name}, nil
}

func (u *UserDBStore) CreateUser(name, password string) (string, error) {
	password, err := hashPassword(password)

	if err != nil {
		fmt.Print(err)
		return "", err
	}

	insertQuery := `INSERT INTO account (name, password) values ($1, $2) RETURNING id;`
	var id int
	err = u.db.QueryRow(insertQuery, name, password).Scan(&id)

	if err != nil {
		fmt.Print(err)
		return "", err
	}

	token, err := createToken(name, id)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserDBStore) GetUsers() ([]User, error) {
	getQuery := `SELECT id, name FROM account;`
	data, err := u.db.Query(getQuery)
	if err != nil {
		println(err.Error())
	}
	users := []User{}

	for data.Next() {
		var (
			ID   int
			Name string
		)
		data.Scan(&ID, &Name)

		users = append(users, User{ID, Name})
	}

	return users, nil
}

func (u *UserDBStore) Login(name, password string) (string, error) {

	getQuery := `SELECT password, id FROM account WHERE name=$1;`
	var hash string
	var id int
	err := u.db.QueryRow(getQuery, name).Scan(&hash, &id)

	if err != nil {
		return "", err
	}

	correctPw := compareHash(password, hash)

	if !correctPw {
		return "", ErrInvalidCredentials
	}

	token, err := createToken(name, id)

	if err != nil {
		return "", err
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func compareHash(challenge, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(challenge))

	return err == nil
}

func createToken(name string, id int) (string, error) {
	claims := &jwtClaims{
		name,
		id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 73).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return t, nil
}
