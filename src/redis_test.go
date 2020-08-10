package main

import (
	"testing"

	"github.com/newrelic/infra-integrations-sdk/integration"
	"github.com/stretchr/testify/assert"
)

func TestEntity_Metadata(t *testing.T) {
	tests := map[string]struct {
		args         *argumentList
		wantMetadata *integration.EntityMetadata
	}{
		"LocalEntity": {
			&argumentList{RemoteMonitoring: false},
			nil,
		},
		"RemoteEntity": {
			&argumentList{Hostname: "localhost", Port: 8080, RemoteMonitoring: true},
			&integration.EntityMetadata{Name: "localhost:8080", Namespace: entityRemoteType},
		},
		"RemoteEntityUnixSocket": {
			&argumentList{Hostname: "localhost", Port: 8080, RemoteMonitoring: true, UnixSocketPath: "/socket/path"},
			&integration.EntityMetadata{Name: "localhost:8080", Namespace: entityRemoteType},
		},
		"RemoteEntityUnixSocketAndUseUnixSocket": {
			&argumentList{Hostname: "localhost", Port: 8080, RemoteMonitoring: true, UnixSocketPath: "/socket/path", UseUnixSocket: true},
			&integration.EntityMetadata{Name: "localhost:/socket/path", Namespace: entityRemoteType},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			i, err := integration.New(name, integrationVersion)
			assert.NoError(t, err)

			e, err := entity(i, tc.args)
			assert.NoError(t, err)

			if tc.wantMetadata == nil {
				assert.Nil(t, e.Metadata)
			} else {
				assert.Equal(t, tc.wantMetadata.Name, e.Metadata.Name)
				assert.Equal(t, tc.wantMetadata.Namespace, e.Metadata.Namespace)
			}
		})
	}
}
