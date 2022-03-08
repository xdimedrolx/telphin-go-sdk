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

func (suite *RecordSuite) TestGetRecordsByPeriodSuccessful() {
	gock.New(Host).
		Get("/api/ver1.0/client/@me/record/").
		MatchParam("page", "1").
		MatchParam("per_page", "3").
		MatchParam("order", "asc").
		MatchParam("start_datetime", "2021-01-01 00:00:00").
		MatchParam("end_datetime", "2021-01-01 01:00:00").
		Reply(200).
		JSON(getTestData(suite.T(), "get_records_200.json"))

	order := OrderAsc
	query := RecordsRequest{
		Page:          1,
		PerPage:       3,
		Order:         &order,
		StartDatetime: "2021-01-01 00:00:00",
		EndDatetime:   "2021-01-01 01:00:00",
	}

	records, err := suite.client.GetRecords("@me", query)

	if suite.NoError(err) {
		suite.Len(*records, 3)
	}
}

func (suite *RecordSuite) TestDeleteRecordSuccessful() {
	gock.New(Host).
		Delete("/api/ver1.0/client/@me/record/").
		Reply(204)

	err := suite.client.DeleteRecord("@me", "609560-a55c44ab973a4641a30b85236ae57a41")

	suite.NoError(err)
}

func (suite *RecordSuite) TestDeleteRecordWhenItDoesNotFound() {
	gock.New(Host).
		Delete("/api/ver1.0/client/@me/record/").
		Reply(404).
		JSON(getTestData(suite.T(), "delete_record_404.json"))

	err := suite.client.DeleteRecord("@me", "609560-a55c44ab973a4641a30b85236ae57a41")

	suite.Error(err)
	suite.Equal(404, err.(*ErrorResponse).Code)
}
