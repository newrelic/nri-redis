package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
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
		"dst no file replacement": { "{app_name}-{arch}-{version}", "/tmp/{arch}/{app_name}/{version}/file",
			"nri-foobar", "1.2.3", "amd64",
			"nri-foobar-amd64-1.2.3","/tmp/amd64/nri-foobar/1.2.3/file"},
		"dst src replacement": { "{app_name}-{arch}-{version}", "/tmp/{arch}/{app_name}/{version}/{src}",
			"nri-foobar", "1.2.3", "amd64",
			"nri-foobar-amd64-1.2.3","/tmp/amd64/nri-foobar/1.2.3/nri-foobar-amd64-1.2.3"},
		"dst multiple replacements": { "{app_name}-{arch}-{version}", "/tmp/{arch}/{app_name}/{version}/{app_name}-{arch}-{version}",
			"nri-foobar", "1.2.3", "amd64",
			"nri-foobar-amd64-1.2.3","/tmp/amd64/nri-foobar/1.2.3/nri-foobar-amd64-1.2.3"},
		"src multiple replacements": { "{app_name}-{arch}-{version}-{app_name}-{arch}-{version}", "/tmp/{arch}/{app_name}/{version}/file",
			"nri-foobar", "1.2.3", "amd64",
			"nri-foobar-amd64-1.2.3-nri-foobar-amd64-1.2.3","/tmp/amd64/nri-foobar/1.2.3/file"},
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
		{"{app_name}-{arch}-{version}.txt", "{arch}/{app_name}/{src}", []string{"amd64"}},
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

	file, err := os.Create(path.Join(src, "nri-foobar-amd64-2.0.0.txt"))
	assert.NoError(t, err)

	_, err = file.Write([]byte("test"))
	assert.NoError(t, err)

	err = file.Close()
	assert.NoError(t, err)

	err = uploadArtifacts(cfg, schema)
	assert.NoError(t, err)

	name := path.Join(dest, "amd64/nri-foobar/nri-foobar-amd64-2.0.0.txt")
	_, err = os.Stat(name)
	assert.NoError(t, err)

	//if _,  os.IsNotExist(err) {
	//	t.Fatalf("file %s doesn't exist", name)
	//}
}