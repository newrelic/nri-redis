package main

import (
	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestEntity_RemoteEntity(t *testing.T) {
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
