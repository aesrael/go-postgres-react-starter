package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/stretchr/testify/suite"
	"goapp/packages/db"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"goapp/packages/api"
	"goapp/packages/config"
)

type registerTestSuite struct {
	suite.Suite
	port string
	db   *sql.DB
}

func TestRegisterTest(t *testing.T) {
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

	suite.Run(t, &registerTestSuite{
		port: serverPort,
		db:   conn,
	})
}

func (s *registerTestSuite) AfterTest(su, t string) {
	defer s.db.Close()
	_, err := s.db.Query(db.DeleteUser, "test@test.com")
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("Failed executing DELETE USER QUERY")
	}
	api.StopServer()
}

func (s *registerTestSuite) TestRegister() {
	values := map[string]string{"name": "testName", "email": "test@test.com", "password": "password"}
	jsonData, err := json.Marshal(values)
	s.NoError(err)
	response, err := http.Post(fmt.Sprintf("http://localhost%s/api/register", s.port),
		"application/json", bytes.NewBuffer(jsonData))
	defer response.Body.Close()
	s.NoError(err)

	s.Equal(http.StatusOK, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	bodyString := string(body[:])
	s.Equal("{\"success\":true}", bodyString)
	s.db.QueryRow(db.CheckUserExists, "test@test.com")
	var exists bool
	err = s.db.QueryRow(db.CheckUserExists, "test@test.com").Scan(&exists)
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("Failed executing CheckUserExists")
	}
	s.True(exists)
}
