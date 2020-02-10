package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cnabio/cnab-go/driver"

	"github.com/cnabio/duffle/pkg/duffle/home"

	"github.com/cnabio/cnab-go/bundle"
	"github.com/cnabio/cnab-go/credentials"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

func CreateTestHome(t *testing.T) home.Home {
	t.Helper()
	tempDir, err := ioutil.TempDir("", "duffle")
	if err != nil {
		t.Fatal(err)
	}
	duffleHome = tempDir
	testHome := home.Home(tempDir)
	dirs := []string{
		testHome.String(),
		testHome.Bundles(),
		testHome.Logs(),
		testHome.Plugins(),
		testHome.Claims(),
		testHome.Credentials(),
	}
	if err := ensureDirectories(dirs); err != nil {
		t.Fatal(err)
	}
	if err := ensureFiles([]string{testHome.Repositories()}); err != nil {
		t.Fatal(err)
	}
	return testHome
}

func TestLoadCredentials(t *testing.T) {
	cred1 := credentials.CredentialSet{
		Name: "first",
		Credentials: []credentials.CredentialStrategy{
			{Name: "knapsack", Source: credentials.Source{Value: "cred1"}},
			{Name: "gym-bag", Source: credentials.Source{Value: "cred1"}},
		},
	}
	cred2 := credentials.CredentialSet{
		Name: "second",
		Credentials: []credentials.CredentialStrategy{
			{Name: "knapsack", Source: credentials.Source{Value: "cred2"}},
			{Name: "haversack", Source: credentials.Source{Value: "cred2"}},
		},
	}
	cred3 := credentials.CredentialSet{
		Name: "third",
		Credentials: []credentials.CredentialStrategy{
			{Name: "haversack", Source: credentials.Source{Value: "cred3"}},
		},
	}

	// The above should generate:
	// -- knapsack: cred2
	// -- havershack: cred3
	// -- gym-bag: cred1

	tmpdir, err := ioutil.TempDir("", "duffle-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	files := make([]string, 3)
	for i, c := range []credentials.CredentialSet{cred1, cred2, cred3} {
		data, err := yaml.Marshal(c)
		t.Log(string(data))
		if err != nil {
			t.Fatal(err)
		}
		fp := filepath.Join(tmpdir, c.Name+".yaml")
		if err := ioutil.WriteFile(fp, data, 0644); err != nil {
			t.Fatal(err)
		}
		files[i] = fp
	}

	bun := bundle.Bundle{
		Name: "test-load-creds",
		Credentials: map[string]bundle.Credential{
			"knapsack": {
				Location: bundle.Location{
					EnvironmentVariable: "KNAP",
				},
			},
			"haversack": {
				Location: bundle.Location{
					EnvironmentVariable: "HAVER",
				},
			},
			"gym-bag": {
				Location: bundle.Location{
					EnvironmentVariable: "GYM",
				},
			},
		},
	}

	is := assert.New(t)
	creds, err := loadCredentials(files, &bun)
	is.NoError(err)
	is.Equal("cred2", creds["knapsack"])
	is.Equal("cred3", creds["haversack"])
	is.Equal("cred1", creds["gym-bag"])
}

func TestFindCreds(t *testing.T) {
	credDir, err := ioutil.TempDir("", "credTest")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(credDir)
	is := assert.New(t)

	tests := map[string]struct {
		input            string
		expectedFilePath func() string
	}{
		"no path yaml": {
			input: "creds1",
			expectedFilePath: func() string {
				credPath := filepath.Join(credDir, "creds1.yaml")
				err = ioutil.WriteFile(credPath, []byte("test"), 0644)
				if err != nil {
					t.Fatal(err)
				}
				return credPath
			},
		},
		"no path yml": {
			input: "creds2",
			expectedFilePath: func() string {
				credPath := filepath.Join(credDir, "creds2.yml")
				err = ioutil.WriteFile(credPath, []byte("test"), 0644)
				if err != nil {
					t.Fatal(err)
				}
				return credPath
			},
		},
		"path": {
			input: "testdata/dufflehome/credentials/testing.yaml",
			expectedFilePath: func() string {
				return "testdata/dufflehome/credentials/testing.yaml"
			},
		},
	}

	for name, testCase := range tests {
		is.Equal(testCase.expectedFilePath(), findCreds(credDir, testCase.input), "Fail on test: "+name)
		os.RemoveAll(filepath.Join(credDir, "*"))
	}
}

func TestMakeOpRelocator(t *testing.T) {
	is := assert.New(t)

	tests := map[string]struct {
		relMapFile                      string
		invocationImage                 string
		expectedErrorMessage            string
		expectedInvocationImage         string
		expectedOpRelocatorErrorMessage string
	}{
		"bad relocation mapping file": {
			relMapFile:                      "no-such-file",
			invocationImage:                 "example.com/original",
			expectedErrorMessage:            "failed to read relocation mapping from no-such-file:",
			expectedInvocationImage:         "",
			expectedOpRelocatorErrorMessage: "",
		},
		"valid relocation mapping file": {
			relMapFile:                      "testdata/oprelocator/relmap.json",
			invocationImage:                 "example.com/original",
			expectedErrorMessage:            "",
			expectedInvocationImage:         "example.com/relocated",
			expectedOpRelocatorErrorMessage: "",
		},
		"omitted relocation mapping file": {
			relMapFile:                      "",
			invocationImage:                 "example.com/original",
			expectedErrorMessage:            "",
			expectedInvocationImage:         "example.com/original",
			expectedOpRelocatorErrorMessage: "",
		},
		"relocation mapping file with malformed contents": {
			relMapFile:                      "testdata/oprelocator/badrelmap.json",
			invocationImage:                 "example.com/original",
			expectedErrorMessage:            "failed to unmarshal relocation mapping:",
			expectedInvocationImage:         "example.com/original",
			expectedOpRelocatorErrorMessage: "",
		},
		"invocation image not in relocation mapping": {
			relMapFile:                      "testdata/oprelocator/relmap.json",
			invocationImage:                 "example.com/other",
			expectedErrorMessage:            "",
			expectedInvocationImage:         "example.com/other",
			expectedOpRelocatorErrorMessage: "invocation image example.com/other not present in relocation mapping map[",
		},
	}

withNextTest:
	for name, testCase := range tests {
		opRelocator, err := makeOpRelocator(testCase.relMapFile)
		if testCase.expectedErrorMessage != "" {
			is.Contains(err.Error(), testCase.expectedErrorMessage, "Failed on test: "+name)
			continue withNextTest
		}
		is.Nil(err, "Failed on test: "+name)

		op := driver.Operation{
			Files: make(map[string]string),
			Image: bundle.InvocationImage{
				BaseImage: bundle.BaseImage{
					Image: testCase.invocationImage,
				},
			},
		}
		err = opRelocator(&op)
		if testCase.expectedOpRelocatorErrorMessage != "" {
			is.Contains(err.Error(), testCase.expectedOpRelocatorErrorMessage, "Failed on test: "+name)
			continue withNextTest
		}
		is.Nil(err, "Failed on test: "+name)

		is.Equal(testCase.expectedInvocationImage, op.Image.Image, "Failed on test: "+name)

		// In the success case, a relocation mapping should be mounted if and only if a relocation mapping was specified
		_, mounted := op.Files["/cnab/app/relocation-mapping.json"]
		is.Equal(testCase.relMapFile != "", mounted, "Failed on test: "+name)
	}
}

func TestConfigureDriver(t *testing.T) {
	is := assert.New(t)

	tests := map[string]struct {
		options map[string]string
		envvars map[string]string
	}{
		"no options": {
			map[string]string{},
			map[string]string{},
		},
		"no env vars set": {
			map[string]string{
				"ONE": "first option",
				"TWO": "second option",
			},
			map[string]string{},
		},
		"some env vars set": {
			map[string]string{
				"ONE":   "first option",
				"TWO":   "second option",
				"THREE": "third option",
			},
			map[string]string{
				"TWO":   "FOO",
				"THREE": "BAR",
			},
		},
		"env var set to empty string": {
			map[string]string{
				"ONE": "first option",
			},
			map[string]string{
				"ONE": "",
			},
		},
	}

	for name, testCase := range tests {
		c := &fakeConfigurable{
			opts: testCase.options,
		}

		for k, v := range testCase.envvars {
			os.Setenv(k, v)
		}

		configureDriver(c)

		for k := range testCase.envvars {
			os.Unsetenv(k)
		}

		is.Equal(testCase.envvars, c.vals, fmt.Sprintf("Failed on test: %s", name))
	}
}

type fakeConfigurable struct {
	opts map[string]string
	vals map[string]string
}

func (c *fakeConfigurable) Config() map[string]string {
	return c.opts
}

func (c *fakeConfigurable) SetConfig(vals map[string]string) {
	c.vals = vals
}
