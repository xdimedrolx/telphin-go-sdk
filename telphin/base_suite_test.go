package telphin

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type BaseSuite struct {
	suite.Suite
	client *Client
}

func (s *BaseSuite) SetupTest() {
	gock.New(Host).
		Post("/oauth/token").
		MatchType("url").
		BodyString(`client_id=client&client_secret=secret&grant_type=client_credentials`).
		Persist().
		Reply(200).
		JSON(getTestData(s.T(), "post_token.json"))

	s.client.GetAccessToken()
}

func (s *BaseSuite) SetupSuite() {
	s.client, _ = NewClient("clientId", "secret", Host)
	s.client.SetLogger(WrapLogrus(logrus.New()))
}

func (s *BaseSuite) TearDownTest() {
	gock.Clean()
}
