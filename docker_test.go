//
// docker_test.go
//
// Copyright (c) 2016-2017 Junpei Kawamoto
//
// This file is part of Roadie queue manager.
//
// Roadie Queue Manager is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Roadie Queue Manager is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// Test to obtain a container ID.
func TestGetID(t *testing.T) {

	containerName := "sample-container-name"
	containerID := "ABCDEFG"

	var (
		docker *Docker
		res    string
		err    error
	)

	// Test for existing container.
	docker = &Docker{

		command: func(name string, args ...string) *exec.Cmd {

			cmd := name + " " + strings.Join(args, " ")
			expected := fmt.Sprintf("docker ps -aq -f name=%s", containerName)
			if cmd != expected {
				t.Error("Generated command isn't correct:", cmd)
			}

			return exec.Command("echo", containerID)

		},
	}

	res, err = docker.GetID(containerName)
	if err != nil {
		t.Fatal(err.Error())
	}
	if res != containerID {
		t.Error("GetID to an existing container returns a wrong container ID:", res)
	}

	// Test for non existing container.
	docker = &Docker{
		command: func(name string, args ...string) *exec.Cmd {
			return exec.Command("echo")
		},
	}

	res, err = docker.GetID(containerName)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(res) != 0 {
		t.Error("GetID to non existing container returns a wrong container ID:", res)
	}

}

// Test to create a container.
func TestCreateContainer(t *testing.T) {

	image := "jkawamoto/roadie-gcp"
	containerName := "sample-container"
	script := `sample
  script
  `

	docker := &Docker{

		command: func(name string, args ...string) *exec.Cmd {

			cmd := name + " " + strings.Join(args, " ")
			expected := fmt.Sprintf("docker run -i --name %s %s --no-shutdown", containerName, image)
			if cmd != expected {
				t.Error("Generated command isn't correct:", cmd)
			}

			return exec.Command("echo", fmt.Sprintf("\"%s\"", string(script)))

		},
	}

	if err := docker.CreateContainer(image, containerName, []byte(script)); err != nil {
		t.Error("CreateContainer returns an error:", err.Error())
	}

}

// Test to delete a container.
func TestDeleteContainer(t *testing.T) {

	containerID := "ABCDEFG"
	docker := &Docker{

		command: func(name string, args ...string) *exec.Cmd {

			cmd := name + " " + strings.Join(args, " ")
			expected := fmt.Sprintf("docker rm -f %s", containerID)
			if cmd != expected {
				t.Error("Generated command isn't correct:", cmd)
			}

			return exec.Command("echo", containerID)

		},
	}

	if err := docker.DeleteContainer(containerID); err != nil {
		t.Error("DeleteContainer returns an error:", err.Error())
	}

}
