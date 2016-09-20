//
// cli_test.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This file is part of Roadie queue manager.
//
// Roadie Queue Manager  is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Roadie Queue Manager  is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

func TestRun_versionFlag(t *testing.T) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	cli := &CLI{outStream: outStream, errStream: errStream}
	args := strings.Split("./roadie-queue-manager -version", " ")

	status := cli.Run(args)
	if status != ExitCodeOK {
		t.Errorf("expected %d to eq %d", status, ExitCodeOK)
	}

	expected := fmt.Sprintf("roadie-queue-manager version %s", Version)
	if !strings.Contains(errStream.String(), expected) {
		t.Errorf("expected %q to eq %q", errStream.String(), expected)
	}
}

// Define constant for querying state.
const (
	getID = iota
	deleteContainer
	createContainer
)

var scriptTestExecuteScript = &QueuedScript{
	InstanceName: "instance-1",
	Image:        "roadie-base-image",
	Body: ScriptBody{
		APT: []string{
			"python-numpy",
			"python-scipy",
		},
		Source: "gs://sample-bucket/path/to/source",
		Data: []string{
			"gs://sample-bucket/path/to/data1",
			"gs://sample-bucket/path/to/data2",
		},
		Run: []string{
			"cmd-1",
			"cmd-2",
		},
		Result: "gs://sample-bucket/path/to/result",
		Upload: []string{
			"result-1",
			"result-2",
		},
	},
}

type mockDockerRequester struct {
	containerID   string
	containerName string
	state         int
	image         string
}

func (m *mockDockerRequester) GetID(name string) (res string, err error) {

	if m.state != getID {
		err = fmt.Errorf("GetID was called in wrong order: %d", m.state)
		return
	}

	res = m.containerID
	if len(m.containerID) == 0 {
		// This case simulates that there are no same name containers exists.
		// Thus, the next state should be createContainer.
		m.state = createContainer
	} else {
		// Otherwise, the existing container should be deleted.
		// Thus, the next state should be deleteContainer.
		m.state = deleteContainer
	}
	return

}

func (m *mockDockerRequester) CreateContainer(image, name string, script []byte) error {

	if m.state != createContainer {
		return fmt.Errorf("CreateContainer was called in wrong order: %d", m.state)
	}

	if m.image != image {
		return fmt.Errorf("Image of the creating container is wrong: %s (%s expected)", image, m.image)
	}

	if m.containerName != name {
		return fmt.Errorf("Container name is wrong: %s (%s expected)", name, m.containerName)
	}

	// State should be reset.
	m.state = getID
	return nil

}

func (m *mockDockerRequester) DeleteContainer(id string) error {

	if m.state != deleteContainer {
		return fmt.Errorf("DeleteContainer was called in wrong order: %d", m.state)
	}

	if m.containerID != id {
		return fmt.Errorf("Try to delete a wrong container: %s (%s expected)", id, m.containerID)
	}

	// The next state should be createContainer
	m.state = createContainer
	return nil

}

// Test executeScript calles methods in the expected order with an empty script.
func TestExecuteScriptWithEmptyScript(t *testing.T) {

	var err error
	var docker *mockDockerRequester

	// Test with an empty script file.
	fp, err := ioutil.TempFile("", "test-script-")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(fp.Name())

	docker = &mockDockerRequester{}
	err = executeScript(docker, fp.Name())
	if err != nil {
		t.Error("executeScript returns an error:", err.Error())
	}

	if _, exist := os.Stat(fp.Name()); exist == nil {
		t.Error("Script file was not deleted.")
	}

}

// Test executeScript calles methods in the expected order with an existing container.
func TestExecuteScriptWithExistingContainer(t *testing.T) {

	var err error
	var docker *mockDockerRequester

	fp, err := ioutil.TempFile("", "test-script-")
	if err != nil {
		t.Fatal(err.Error())
	}
	fp.Close()
	defer os.Remove(fp.Name())

	raw, err := yaml.Marshal(&scriptTestExecuteScript)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fp.Name(), raw, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}

	docker = &mockDockerRequester{
		containerID:   "asdfghjkl",
		image:         sampleScript.Image,
		containerName: sampleScript.InstanceName,
	}
	err = executeScript(docker, fp.Name())
	if err != nil {
		t.Error("executeScript returns an error:", err.Error())
	}

	if _, exist := os.Stat(fp.Name()); exist == nil {
		t.Error("Script file was not deleted.")
	}

}

func TestExecuteScriptWithoutExistingContainer(t *testing.T) {

	var err error
	var docker *mockDockerRequester

	fp, err := ioutil.TempFile("", "test-script-")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(fp.Name())

	raw, err := yaml.Marshal(&scriptTestExecuteScript)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile(fp.Name(), raw, 0644)
	if err != nil {
		t.Fatal(err.Error())
	}
	docker = &mockDockerRequester{
		image:         sampleScript.Image,
		containerName: sampleScript.InstanceName,
	}
	err = executeScript(docker, fp.Name())
	if err != nil {
		t.Error("executeScript returns an error:", err.Error())
	}

	if _, exist := os.Stat(fp.Name()); exist == nil {
		t.Error("Script file was not deleted.")
	}

}
