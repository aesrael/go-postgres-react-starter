package controller

import (
	"database/sql"
	"net/http"
	"time"

	"goapp/packages/config"

	"goapp/packages/db"
	"goapp/packages/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	db.User
	jwt.StandardClaims
}

func Pong(c *fiber.Ctx) error {
	return c.SendString("pong")
}

func CreateUser(c *fiber.Ctx, dbConn *sql.DB) error {
	user := db.User{}

	if err := c.BodyParser(user); err != nil {
		return err
	}

	if user.UserExists(dbConn) {
		return c.Status(400).JSON(&fiber.Map{"success": false, "errors": []string{"email already exists"}})
	}

	if errs := utils.ValidateUser(user); len(errs) > 0 {
		return c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"success": false, "errors": errs})
	}

	user.HashPassword()
	_, err := dbConn.Query(db.CreateUserQuery, user.Name, user.Password, user.Email)
	if err != nil {
		return err
	}
	return c.JSON(&fiber.Map{"success": true, "msg": "User created succesfully"})
}

func Session(c *fiber.Ctx, dbConn *sql.DB) error {
	user := db.User{}

	if err := dbConn.QueryRow(db.GetUserQuery, user.Email).
		Scan(&user.Name, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return err
	}
	return c.JSON(&fiber.Map{"success": true, "user": user})
}

func Login(c *fiber.Ctx, dbConn *sql.DB) error {
	loginUser := db.User{}

	if err := c.BodyParser(loginUser); err != nil {
		return err
	}

	user := db.User{}
	if err := dbConn.QueryRow(db.GetUserQuery, user.Email).
		Scan(&user.Name, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"success": false, "errors": []string{"Incorrect credentials"}})
		}
	}

	match := utils.ComparePassword(user.Password, loginUser.Password)
	if !match {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"success": false, "errors": []string{"Incorrect credentials"}})
	}

	//expiration time of the token ->30 mins
	expirationTime := time.Now().Add(30 * time.Minute)

	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var jwtKey = []byte(config.Config["JWT_KEY"])
	tokenValue, err := token.SignedString(jwtKey)

	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   tokenValue,
		Expires: expirationTime,
		Domain:  "*",
	})

	return c.JSON(&fiber.Map{"success": true, "msg": "logged in succesfully", "user": claims.User, "token": tokenValue})
}
