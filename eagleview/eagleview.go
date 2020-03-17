package eagleview

import (
	"errors"
)

type Client struct {
	InitialRequest InitialRequest
	Token          Token
}

func NewClient(username, password, sourceId, clientSecret string) (*Client, error) {
	if username == "" {
		return nil, errors.New("username is missing")
	}

	if password == "" {
		return nil, errors.New("password is missing")
	}

	if sourceId == "" {
		return nil, errors.New("source_id is missing")
	}

	if clientSecret == "" {
		return nil, errors.New("client_secret is missing")
	}

	newClient := new(Client)
	newClient.InitialRequest.Username = username
	newClient.InitialRequest.Password = password
	newClient.InitialRequest.SourceId = sourceId
	newClient.InitialRequest.ClientSecret = clientSecret

	return newClient, nil
}
