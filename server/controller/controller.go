package controller

import (
	"database/sql"
	"fmt"
	"github.com/go-postgres-jwt-react-starter/server/config"
	"net/http"
	"time"

	"github.com/go-postgres-jwt-react-starter/server/db"
	"github.com/go-postgres-jwt-react-starter/server/errors"
	"github.com/go-postgres-jwt-react-starter/server/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("secret")

//Claims jwt claims struct
type Claims struct {
	db.User
	jwt.StandardClaims
}

// Pong tests that api is working
func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"ping": "pong"})
}


//Initiate Password reset email with reset url
func InitiatePasswordReset(c *gin.Context){
	var createReset db.CreateReset
	c.Bind(&createReset)
	if id,ok := checkAndRetrieveUserIDViaEmail(createReset); ok{
		link := fmt.Sprintf("%s/reset/%d",config.CLIENT_URL,id)
		//Reset link is returned in json response for testing purposes since no email service is integrated
		c.JSON(http.StatusOK, gin.H{"success": true, "msg": "Successfully sent reset mail to " + createReset.Email, "link":link})
	} else{
		c.JSON(http.StatusNotFound,gin.H{"success": false, "errors":"No user found for email: " + createReset.Email})
	}
}

func ResetPassword(c *gin.Context){
	var resetPassword db.ResetPassword
	c.Bind(&resetPassword)
	if ok,errStr := utils.ValidatePasswordReset(resetPassword); ok{
		password := db.CreateHashedPassword(resetPassword.Password)
		_,err := db.DB.Query(db.UpdateUserPasswordQuery,resetPassword.ID,password)
		errors.HandleErr(c, err)
		c.JSON(http.StatusOK,gin.H{"success":true, "msg": "User password reset successfully"})
	} else{
		c.JSON(http.StatusOK, gin.H{"success":false,"errors":errStr})
	}

}


//Create new user
func Create(c *gin.Context) {
	var user db.Register
	c.Bind(&user)
	exists := checkUserExists(user)

	valErr := utils.ValidateUser(user, errors.ValidationErrors)
	if exists == true {
		valErr = append(valErr, "email already exists")
	}
	fmt.Println(valErr)
	if len(valErr) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"success": false, "errors": valErr})
		return
	}
	db.HashPassword(&user)
	_, err := db.DB.Query(db.CreateUserQuery, user.Name, user.Password, user.Email)
	errors.HandleErr(c, err)
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "User created succesfully"})
}

// Session returns JSON of user info
func Session(c *gin.Context) {
	user, isAuthenticated := AuthMiddleware(c, jwtKey)
	if !isAuthenticated {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "user": user})
}

// Login controller
func Login(c *gin.Context) {
	var user db.Login
	c.Bind(&user)

	row := db.DB.QueryRow(db.LoginQuery, user.Email)

	var id int
	var name, email, password, createdAt, updatedAt string

	err := row.Scan(&id, &name, &password, &email, &createdAt, &updatedAt)

	if err == sql.ErrNoRows {
		fmt.Println(sql.ErrNoRows, "err")
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "incorrect credentials"})
		return
	}

	match := db.CheckPasswordHash(user.Password, password)
	if !match {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "msg": "incorrect credentials"})
		return
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
	errors.HandleErr(c, err)
	// c.SetCookie("token", tokenString, expirationTime, "", "*", true, false)
	http.SetCookie(c.Writer, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	fmt.Println(tokenString)
	c.JSON(http.StatusOK, gin.H{"success": true, "msg": "logged in succesfully", "user": claims.User, "token": tokenString})
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
func checkAndRetrieveUserIDViaEmail(createReset db.CreateReset) (int,bool){
	rows, err := db.DB.Query(db.CheckUserExists,createReset.Email)
	if err != nil{
		return -1,false
	}
	if !rows.Next(){
		return -1,false
	}
	var id int
	err = rows.Scan(&id)
	if err != nil{
		return -1,false
	}
	return id,true
}
