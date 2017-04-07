package vault

import (
	"fmt"

	log "github.com/Sirupsen/logrus"

	"github.com/hashicorp/vault/api"
)

func New(token string, host string, port int) (*api.Client, error) {

	// specify vault url
	url := fmt.Sprintf("https://%s:%v", host, port)

	// create an API cpnfig
	config := &api.Config{Address: url}

	// read any environment variables
	if err := config.ReadEnvironment(); err != nil {
		log.Warn("Error reading environment variables", err)
	}

	// create a client
	client, err := api.NewClient(config)

	if err != nil {
		log.Fatal("Error creating vault client", err)
	}

	// set token
	client.SetToken(token)

	return client, nil

}
