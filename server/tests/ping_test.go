package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/apex/log"
	"github.com/stretchr/testify/suite"

	"goapp/packages/api"
	"goapp/packages/config"
)

type pingTestSuite struct {
	suite.Suite
	port string
}

func TestPingTest(t *testing.T) {
	err := os.Setenv("ENV", "test")
	if err != nil {
		log.WithField("reason", err.Error()).Fatal("Failed to set ENV environment variable")
	}

	config.InitConfig()
	go api.StartServer()
	time.Sleep(1 * time.Second) // wait for server to start

	serverPort, _ := config.Config[config.SERVER_PORT]

	suite.Run(t, &pingTestSuite{
		port: serverPort,
	})
}

func (s *pingTestSuite) AfterTest(su, t string) {
	api.StopServer()
}

func (s *pingTestSuite) TestPing() {
	response, err := http.Get(fmt.Sprintf("http://localhost%s/api/ping", s.port))
	defer response.Body.Close()
	s.NoError(err)

	s.NoError(err)
	s.Equal(http.StatusOK, response.StatusCode)
	body, _ := ioutil.ReadAll(response.Body)
	bodyString := string(body[:])
	s.Equal("pong", bodyString)
}
