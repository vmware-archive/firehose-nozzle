package api

import (
	"errors"

	"github.com/cloudfoundry-community/go-cfclient"
)

type apiTokenFetcher struct {
	clientConfig *cfclient.Config
}

func NewAPITokenFetcher(apiUrl string, username string, password string, sslSkipVerify bool) *apiTokenFetcher {
	client := &cfclient.Config{
		ApiAddress:        apiUrl,
		Username:          username,
		Password:          password,
		SkipSslValidation: sslSkipVerify,
	}

	return &apiTokenFetcher{
		clientConfig: client,
	}
}

func (api *apiTokenFetcher) FetchAuthToken() (string, error) {
	client := cfclient.NewClient(api.clientConfig)
	//client.Endpoint.DopplerEndpoint
	token := client.GetToken()
	if token == "" {
		return "", errors.New("Error fetching token")
	}
	return token, nil
}
