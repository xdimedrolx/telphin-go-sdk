package telphin

import (
	"testing"

	"github.com/stretchr/testify/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type ExtensionSuite struct {
	BaseSuite
}

func TestExtensionSuite(t *testing.T) {
	suite.Run(t, new(ExtensionSuite))
}

func (suite *ExtensionSuite) TestCreateCallbackSuccessful() {
	gock.New(Host).
		Post("/api/ver1.0/extension/303052/callback/").
		Reply(200).
		JSON(getTestData(suite.T(), "post_callback.json"))

	transfer := "12345"
	waitForPickup := 30
	soundID := 94478
	dstAni := "+79530000001"

	callback, err := suite.client.CreateCallback(303052, CallbackRequest{
		SrcNum:                 []string{"1234"},
		DstNum:                 "+79530000000",
		DstAni:                 &dstAni,
		TransferAfterSrcHangup: &transfer,
		DstAnnounceSoundID:     &soundID,
		WaitForPickup:          &waitForPickup,
	})

	if suite.NoError(err) {
		suite.NotNil(callback)
		suite.Equal("02abb9ee685a11ea94bfe17eb6eaecba", callback.CallID)
	}
}
