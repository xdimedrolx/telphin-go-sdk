package telphin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	gock "gopkg.in/h2non/gock.v1"
)

func TestCreateClient(t *testing.T) {
	client, err := NewClient("clientId", "secret", Host)
	assert.NoError(t, err)
	assert.NotNil(t, client)

	client, err = NewClient("", "", Host)
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestGetAccessTokenWhenCredentialsIsValid(t *testing.T) {
	gock.Clean()

	gock.New(Host).
		Post("/oauth/token").
		MatchType("url").
		BodyString(`client_id=client&client_secret=secret&grant_type=client_credentials`).
		Reply(200).
		JSON(getTestData(t, "post_token.json"))

	client, _ := NewClient("client", "secret", Host)
	token, err := client.GetAccessToken()

	if assert.NoError(t, err) {
		assert.NotNil(t, token)
		assert.Equal(t, "59sBiAq3LOTp40CGQAoc6EqQpjwSTdr0", token.Token)
	}
}

func TestGetAccessTokenWhenCredentialsIsInvalid(t *testing.T) {
	gock.Clean()

	gock.New(Host).
		Post("/oauth/token").
		MatchType("url").
		Reply(400)

	client, _ := NewClient("client", "secret1asd", Host)
	token, err := client.GetAccessToken()

	assert.Error(t, err)
	assert.Empty(t, token.Token)
}
