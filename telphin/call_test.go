package telphin

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type CallSuite struct {
	BaseSuite
}

func TestCallSuite(t *testing.T) {
	suite.Run(t, new(CallSuite))
}

func (suite *CallSuite) TestGetCallHistory() {
	suite.Run("call found", func() {
		gock.New(Host).
			Get("/api/ver1.0/client/@me/call_history/caa308e3333c4a19aaed9c0254b9c61b").
			Reply(200).
			JSON(getTestData(suite.T(), "get_call_history_200_1.json"))

		history, err := suite.client.GetCallHistory("@me", "caa308e3333c4a19aaed9c0254b9c61b")

		if suite.NoError(err) {
			suite.NotNil(history)
			suite.Equal("caa308e3333c4a19aaed9c0254b9c61b", strings.ToLower(history.UUID))
		}
	})

	suite.Run("call not found", func() {
		gock.New(Host).
			Get("/api/ver1.0/client/@me/call_history/ce9308e3295c4a19aaed5c0254c9c612").
			Reply(404).
			JSON(getTestData(suite.T(), "error_404.json"))

		history, err := suite.client.GetCallHistory("@me", "ce9308e3295c4a19aaed5c0254c9c612")

		if suite.Error(err) {
			suite.Nil(history)
			suite.IsType(404, err.(*ErrorResponse).Code)
			suite.IsType("Not found", err.(*ErrorResponse).Cause)
		}
	})

	suite.Run("error", func() {
		gock.New(Host).
			Get("/api/ver1.0/client/@me/call_history/ce9308e3295c4a19aaed5c0254c9c611").
			Reply(429).
			JSON(getTestData(suite.T(), "get_call_history_429.json"))

		history, err := suite.client.GetCallHistory("@me", "ce9308e3295c4a19aaed5c0254c9c611")

		if suite.Error(err) {
			suite.Nil(history)
			suite.IsType(429, err.(*ErrorResponse).Code)
			suite.IsType("Stats request already exists, concurrent request not allowed", err.(*ErrorResponse).Cause)
		}
	})
}
