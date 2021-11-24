package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/stretchr/testify/suite"
	"goapp/packages/api"
	"goapp/packages/config"
	"goapp/packages/db"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"
)

type sessionTestSuite struct {
	suite.Suite
	port string
	db   *sql.DB
}

type sessionResponse struct {
	Success bool    `json:"success"`
	User    db.User `json:"user"`
}

func TestSessionTest(t *testing.T) {
	err := os.Setenv("ENV", "test")
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("Failed to set ENV environment variable")
	}

	config.InitConfig()
	go api.StartServer()
	time.Sleep(1 * time.Second) // wait for server to start

	conn, err := db.ConnectDB()
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("Db connection error occurred")
	}

	serverPort, _ := config.Config[config.SERVER_PORT]

	suite.Run(t, &sessionTestSuite{
		port: serverPort,
		db:   conn,
	})
}

func (s *sessionTestSuite) AfterTest(su, t string) {
	defer s.db.Close()
	_, err := s.db.Query(db.DeleteUser, "test@test.com")
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("Failed executing DELETE USER QUERY")
	}
	api.StopServer()
}

func (s *sessionTestSuite) BeforeTest(su, t string) {
	passHash := "$2a$10$rbwisPtAZBtJFDNWhLxcAeiTr/uBSd7UL24u80hhH6qlKlmlcCJJu" // hash for "password"
	_, err := s.db.Query(db.CreateUserQuery, "testName", passHash, "test@test.com")
	s.NoError(err)
}

func (s *sessionTestSuite) TestLoginAndSession() {
	values := map[string]string{"email": "test@test.com", "password": "password"}
	jsonData, err := json.Marshal(values)
	s.NoError(err)
	response, err := http.Post(fmt.Sprintf("http://localhost%s/api/login", s.port),
		"application/json", bytes.NewBuffer(jsonData))
	defer response.Body.Close()
	s.NoError(err)

	s.Equal(http.StatusOK, response.StatusCode)
	token := response.Cookies()[0].Value

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost%s/api/session", s.port), nil)
	s.NoError(err)
	req.Header.Set("Authorization", token)
	res, err := client.Do(req)
	defer res.Body.Close()
	s.NoError(err)

	s.Equal(http.StatusOK, res.StatusCode)
	body, _ := ioutil.ReadAll(res.Body)
	var resObj sessionResponse
	err = json.Unmarshal(body, &resObj)
	s.NoError(err)
	s.Equal("test@test.com", resObj.User.Email)
	s.Equal("testName", resObj.User.Name)
}
