package service

import (
	"go.unistack.org/micro/v3/client"
	"go.unistack.org/micro/v3/config"
)

type clientKey struct{}

// Client to call config service
func Client(c client.Client) config.Option {
	return config.SetOption(clientKey{}, c)
}

type serviceKey struct{}

// Service to which data load
func Service(s string) config.Option {
	return config.SetOption(serviceKey{}, s)
}
