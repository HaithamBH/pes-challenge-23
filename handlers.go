package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func Register(ctx echo.Context) error {
	userFrom := new(UserSignUpForm)

	// parse request into object
	err := json.NewDecoder(ctx.Request().Body).Decode(&userFrom)
	if err != nil || userFrom.Username == "" || userFrom.Email == "" || userFrom.Password == "" {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	// verif if user existance
	userTest := new(User)
	FindCollection("users", bson.M{"email": userFrom.Email}, &userTest)
	if userTest.Email == userFrom.Email {
		return ctx.JSON(http.StatusConflict, "email already exist")
	}

	// create user collection
	newUser := User{
		Name:  userFrom.Username,
		Email: userFrom.Email,
	}
	newUser.HashPassword(userFrom.Password)

	// create user in database
	userID, err := CreateCollection("users", newUser)
	if err != nil || userID == "" {
		fmt.Println(err)
		return ctx.JSON(http.StatusConflict, "can't create user")
	}

	// generate token for user
	token, err := GenerateToken(userID, newUser.Email, newUser.HashedPassword)
	if err != nil {
		return ctx.String(http.StatusInternalServerError, "generate token failed")
	}
	// parse token
	tokenObject := TokenResponse{Token: token}
	return ctx.JSON(http.StatusAccepted, tokenObject)
}

func Login(ctx echo.Context) error {
	loginForm := new(UserLoginForm)

	// parse request into object
	err := json.NewDecoder(ctx.Request().Body).Decode(&loginForm)
	if err != nil || loginForm.Email == "" || loginForm.Password == "" {
		return ctx.String(http.StatusBadRequest, "bad request")
	}

	// verify credentials
	// verify email exist
	user := new(User)
	FindCollection("users", bson.M{"email": loginForm.Email}, &user)
	if user == new(User) {
		return ctx.String(http.StatusBadRequest, "user doesn't exist")
	}
	isAuth, _ := user.PasswordMatch(loginForm.Password)

	if isAuth {
		token, err := GenerateToken(user.ID.Hex(), user.Email, user.HashedPassword)
		if err != nil {
			return ctx.String(http.StatusInternalServerError, "generate token failed")
		}
		tokenObject := TokenResponse{Token: token}
		return ctx.JSON(http.StatusAccepted, tokenObject)
	} else {
		return ctx.JSON(http.StatusUnauthorized, "credential doesn't match")
	}
}
