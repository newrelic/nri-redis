package main

import "testing"

// parse the configuration
func TestParseConfig(t *testing.T) {
	t.Parallel() // marks TLog as capable of running in parallel with other tests
	tests := map[string]struct {
		config string
		expected []uploadArtifactSchema
		isError bool
	}{
		"good config multiple entries": { `
			ohi:
			  - Src: "foo.tar.gz"
			    Dest: "/tmp"
 				Arch:
				  - amd64
				  - 386
			  - Src: "{integration_name}_linux_{version}_{Arch}.tar.gz"
			    Dest: "infrastructure_agent/binaries/linux/{Arch}/"
 				Arch:
				  - ppc
		`, []uploadArtifactSchema{
			{"foo.tar.gz", "/tmp", []string { "amd64", "386"}},
			{"{integration_name}_linux_{version}_{Arch}.tar.gz", "infrastructure_agent/binaries/linux/{Arch}/", []string { "ppc" }},
		}, false},
		"bad config: Src is omitted": { `
			ohi:
			  - Dest: "/tmp"
 				Arch:
				  - amd64
		`, nil,true},
		"bad config Dest is omitted": { `
			ohi:
			  - Src: "foo.tar.gz"
 				Arch:
				  - amd64
		`, nil,true},
		"good config Arch is omitted": { `
			ohi:
			  - Src: "foo.tar.gz"
			    Dest: "/tmp"
		`, []uploadArtifactSchema{
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
