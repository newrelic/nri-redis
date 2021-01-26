package main

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
	"os"
)

var (
	schema = ` 
- src: "foo.tar.gz"
  dest: "/tmp"
  arch:
    - amd64
    - 386
- src: "foo-ARCH.zip"
  dest: "infrastructure_agent/binaries/linux/{Arch}/"
  arch:
    - ppc`

	schemaNoSrc = `
- dest: /tmp
  arch:
   - amd64
`
	schemaNoDest = `
- src: foo.tar.gz
  arch:
    - amd64
`
	schemaNoArch = `
- src: foo.tar.gz
  dest: /tmp
`
)

// parse the configuration
func TestParseConfig(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		schema string
		output []uploadArtifactSchema
	}{
		"multiple entries": { schema, []uploadArtifactSchema{
			{"foo.tar.gz", "/tmp", []string { "amd64", "386"}},
			{"{integration_name}_linux_{version}_{Arch}.tar.gz", "infrastructure_agent/binaries/linux/{Arch}/", []string { "ppc" }},
		}},
		"src is omitted": { schemaNoSrc, []uploadArtifactSchema{
			{"", "/tmp", []string { "amd64" }},
		}},
		"dest is omitted": { schemaNoDest, []uploadArtifactSchema{
			{"foo.tar.gz", "", []string { "amd64" }},
		}},
		"arch is omitted": { schemaNoArch, []uploadArtifactSchema{
			{"foo.tar.gz", "/tmp", nil},
		}},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			schema, err := parseUploadSchema([]byte(tt.schema))
			assert.NoError(t, err)
			assert.EqualValues(t, tt.output, schema)
		})
	}
}

func TestReplacePlaceholders(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		srcTemplate string
		destTemplate string
		appName string
		version string
		arch string
		srcOutput string
		destOutput string
	}{
		"full replacement": { "{app_name}-{arch}-{version}", "/tmp/{arch}/{app_name}/{version}",
			"nri-foobar", "1.2.3", "amd64",
			"nri-foobar-amd64-1.2.3","/tmp/amd64/nri-foobar/1.2.3"},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			src, dest := replacePlaceholders(tt.srcTemplate, tt.destTemplate, tt.arch, tt.appName, tt.version)
			assert.EqualValues(t, tt.srcOutput, src)
			assert.EqualValues(t, tt.destOutput, dest)
		})
	}
}

func TestUploadArtifacts(t *testing.T) {
	schema := []uploadArtifactSchema{
		{"{app_name}-{arch}-{version}", "{arch}/{app_name}", []string{"amd64"}},
	}

	dest := t.TempDir()
	src := t.TempDir()
	cfg := config{
		version:              "2.0.0",
		artifactsDestFolder:  dest,
		artifactsSrcFolder:   src,
		uploadSchemaFilePath: "",
		appName:              "nri-foobar",
	}

	file, err := os.Create(path.Join(src, "nri-foobar-amd64-2.0.0"))
	assert.NoError(t, err)
	err = file.Close()
	assert.NoError(t, err)

	uploadArtifacts(cfg, schema)

	//name := path.Join(dest, "amd64/nri-foobar/nri-foobar-amd64.2.0.0")
	//if _, err := os.Stat(name); os.IsNotExist(err) {
	//	t.Fatalf("file %s doesn't exist", name)
	//}
}