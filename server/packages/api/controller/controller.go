package controller

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"goapp/packages/config"

	"goapp/packages/db"

	"goapp/packages/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte("secret")

//Claims jwt claims struct
type Claims struct {
	db.User
	jwt.StandardClaims
}

// Pong tests that api is working
func Pong(c *fiber.Ctx) error {
	return c.SendString("pong")
}

//Initiate Password reset email with reset url
func InitiatePasswordReset(c *fiber.Ctx, db *sql.DB) error {
	createReset := db.CreateReset{}

	if id, ok := checkAndRetrieveUserIDViaEmail(createReset); ok {
		link := fmt.Sprintf("%s/reset/%d", config.Config["CLIENT_URL"], id)
		//Reset link is returned in json response for testing purposes since no email service is integrated
		return c.JSON(&fiber.Map{"success": true, "msg": "Successfully sent reset mail to " + createReset.Email, "link": link})
	}
	return c.JSON(&fiber.Map{"success": false, "errors": "No user found for email: " + createReset.Email})
}

func ResetPassword(c *fiber.Ctx, db *sql.DB) error {
	var resetPassword db.ResetPassword
	if ok, errStr := utils.ValidatePasswordReset(resetPassword); !ok {
		return c.JSON(&fiber.Map{"success": false, "errors": errStr})
	}
	password := db.CreateHashedPassword(resetPassword.Password)
	_, err := db.DB.Query(db.UpdateUserPasswordQuery, resetPassword.ID, password)
	if err != nil {
		return err
	}
	return c.JSON(&fiber.Map{"success": true, "msg": "User password reset successfully"})
}

func CreateUser(c *fiber.Ctx, db *sql.DB) error {
	user := db.Register{}
	exists := checkUserExists(user)
	if exists {
		return c.Status(400).JSON(&fiber.Map{"success": false, "msg": "email already exists"})
	}
	errs := utils.ValidateUser(user)

	if len(errs) > 0 {
		c.JSON(http.StatusUnprocessableEntity, &fiber.Map{"success": false, "errors": valErr})
		return nil
	}
	db.HashPassword(&user)
	_, err := db.DB.Query(db.CreateUserQuery, user.Name, user.Password, user.Email)
	if err != nil {
		return err
	}
	c.JSON(&fiber.Map{"success": true, "msg": "User created succesfully"})
}

func Session(c *fiber.Ctx, db *sql.DB) error {
	return c.JSON(&fiber.Map{"success": true, "user": "user"})
}

func Login(c *fiber.Ctx, db *sql.DB) error {
	user := db.Login{}

	row := db.DB.QueryRow(db.LoginQuery, user.Email)

	var id int
	var name, email, password, createdAt, updatedAt string

	err := row.Scan(&id, &name, &password, &email, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		fmt.Println(sql.ErrNoRows, "err")
		c.Status(http.StatusUnauthorized).JSON(fiber.Map{"success": false, "msg": "incorrect credentials"})
		return err
	}

	match := db.CheckPasswordHash(user.Password, password)
	if !match {
		return c.Status(http.StatusUnauthorized).JSON(&fiber.Map{"success": false, "msg": "incorrect credentials"})
	}

	//expiration time of the token ->30 mins
	expirationTime := time.Now().Add(30 * time.Minute)

	// Create the JWT claims, which includes the User struct and expiry time
	claims := &Claims{

		User: db.User{
			Name: name, Email: email, CreatedAt: createdAt, UpdatedAt: updatedAt,
		},
		StandardClaims: jwt.StandardClaims{
			//expiry time, expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT token string
	tokenString, err := token.SignedString(jwtKey)

	// c.SetCookie("token", tokenString, expirationTime, "", "*", true, false)
	// http.SetCookie(c.Writer, &http.Cookie{
	// 	Name:    "token",
	// 	Value:   tokenString,
	// 	Expires: expirationTime,
	// })

	return c.JSON(&fiber.Map{"success": true, "msg": "logged in succesfully", "user": claims.User, "token": tokenString})
}

func checkUserExists(user db.Register) bool {
	rows, err := db.DB.Query(db.CheckUserExists, user.Email)
	if err != nil {
		return false
	}
	if !rows.Next() {
		return false
	}
	return true
}

//Returns -1 as ID if the user doesnt exist in the table
func checkAndRetrieveUserIDViaEmail(createReset db.CreateReset) (int, bool) {
	rows, err := db.DB.Query(db.CheckUserExists, createReset.Email)
	if err != nil {
		return -1, false
	}
	if !rows.Next() {
		return -1, false
	}
	var id int
	err = rows.Scan(&id)
	if err != nil {
		return -1, false
	}
	return id, true
}
