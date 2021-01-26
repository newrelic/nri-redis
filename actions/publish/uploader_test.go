package main

import "testing"

// parse the configuration
func TestParseConfig(t *testing.T) {
	t.Parallel() // marks TLog as capable of running in parallel with other tests
	tests := map[string]struct {
		config string
		expected []uploadArtifactConfig
		isError bool
	}{
		"good config multiple entries": { `
			ohi:
			  - src: "foo.tar.gz"
			    dest: "/tmp"
 				arch:
				  - amd64
				  - 386
			  - src: "{integration_name}_linux_{version}_{arch}.tar.gz"
			    dest: "infrastructure_agent/binaries/linux/{arch}/"
 				arch:
				  - ppc
		`, []uploadArtifactConfig {
			{"foo.tar.gz", "/tmp", []string { "amd64", "386"}},
			{"{integration_name}_linux_{version}_{arch}.tar.gz", "infrastructure_agent/binaries/linux/{arch}/", []string { "ppc" }},
		}, false},
		"bad config: src is omitted": { `
			ohi:
			  - dest: "/tmp"
 				arch:
				  - amd64
		`, nil,true},
		"bad config dest is omitted": { `
			ohi:
			  - src: "foo.tar.gz"
 				arch:
				  - amd64
		`, nil,true},
		"good config arch is omitted": { `
			ohi:
			  - src: "foo.tar.gz"
			    dest: "/tmp"
		`, []uploadArtifactConfig {
			{"foo.tar.gz", "/tmp", nil},
		},false},
	}
	for name, _ := range tests {
		//tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			t.Log(name)
		})
	}
}
