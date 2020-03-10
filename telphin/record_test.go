package telphin

import (
	"io/ioutil"
	"testing"

	"github.com/gobuffalo/suite"
	gock "gopkg.in/h2non/gock.v1"
)

type RecordSuite struct {
	BaseSuite
}

func TestRecordSuite(t *testing.T) {
	suite.Run(t, new(RecordSuite))
}

func (suite *RecordSuite) TestGetRecordStorageUrlLSuccessful() {
	gock.New(Host).
		Get("/api/ver1.0/client/@me/record/test-id/storage_url/").
		Reply(200).
		JSON(getTestData(suite.T(), "record_storage_url.json"))

	url, err := suite.client.GetRecordStorageUrl("@me", "test-id")

	if suite.NoError(err) {
		suite.NotNil(url)
		suite.Equal("https://storage.telphin.ru/12345", url.RecordUrl)
	}
}

func (suite *RecordSuite) TestGetRecordSuccessful() {
	fileData := getTestData(suite.T(), "record.mp3")

	gock.New(Host).
		Get("/api/ver1.0/client/@me/record/test-id").
		Reply(200).
		AddHeader("Content-Type", "audio/mpeg").
		AddHeader("Content-Disposition", "attachment; filename=\"record-2020-01-01.mp3\"").
		BodyString(string(fileData))

	record, err := suite.client.GetRecord("@me", "test-id")

	if suite.NoError(err) {
		suite.NotNil(record)
		suite.Equal("record-2020-01-01.mp3", record.Name())
		content, _ := ioutil.ReadAll(record)
		suite.Equal(fileData, content)
	}
}
