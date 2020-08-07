package main

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
)

func TestEntity_LocalEntity(t *testing.T) {
	args := argumentList{
		RemoteMonitoring: false,
	}
	i, err := integration.New("test", integrationVersion)
	assert.NoError(t, err)

	e, err := entity(i, &args)
	assert.NoError(t, err)
	assert.Nil(t, e.Metadata)
}

func TestEntity_RemoteEntityPort(t *testing.T) {
	args := argumentList{
		Hostname:         "localhost",
		Port:             8080,
		RemoteMonitoring: true,
	}
	i, err := integration.New("test", integrationVersion)
	assert.NoError(t, err)

	e, err := entity(i, &args)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:8080", e.Metadata.Name)
	assert.Equal(t, entityRemoteType, e.Metadata.Namespace)

}
func TestEntity_RemoteEntityUnixSocket(t *testing.T) {
	args := argumentList{
		Hostname:         "localhost",
		Port:             8080,
		RemoteMonitoring: true,
		UnixSocketPath:   "/socket/path",
	}
	i, err := integration.New("test", integrationVersion)
	assert.NoError(t, err)

	e, err := entity(i, &args)
	assert.NoError(t, err)
	assert.Equal(t, "localhost:/socket/path", e.Metadata.Name)
	assert.Equal(t, entityRemoteType, e.Metadata.Namespace)
}
