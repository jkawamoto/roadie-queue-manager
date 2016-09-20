//
// docker.go
//
// Copyright (c) 2016 Junpei Kawamoto
//
// This file is part of Roadie queue manager.
//
// Roadie is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Roadie is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// commandExecutor is a type of function takes command name and arguments,
// and returns a pointer of exec.Cmd to be run a command.
type commandExecutor func(string, ...string) *exec.Cmd

// Docker defines an inteface of docker API.
type Docker struct {
	// Command executor which will be used to run any command in this class.
	command commandExecutor
}

// NewDocker returns a new docker interface.
func NewDocker() *Docker {
	// Default command executor is `exec.Command`.
	return &Docker{
		command: exec.Command,
	}
}

// GetID returns the ID associated with a given container name.
func (docker *Docker) GetID(name string) (res string, err error) {

	// docker ps -aq -f name=<name>
	cmd := docker.command("docker", "ps", "-aq", "-f", fmt.Sprintf("name=%s", name))
	output, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer output.Close()

	err = cmd.Start()
	if err != nil {
		return
	}

	buf, err := ioutil.ReadAll(output)
	if err != nil {
		return
	}
	res = strings.TrimRight(string(buf), "\n")

	err = cmd.Wait()
	return

}

// CreateContainer creates a container based on a given image. The container will
// have a given name and receive a given script.
// This function will block when the container ends.
func (docker *Docker) CreateContainer(image, name string, script []byte) (err error) {

	// docker run -i --name <name> <image> --no-shutdown
	cmd := docker.command("docker", "run", "-i", "--name", name, image, "--no-shutdown")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	writer := bufio.NewWriter(stdin)

	err = cmd.Start()
	if err != nil {
		return
	}

	_, err = writer.Write(script)
	if err != nil {
		return
	}
	err = writer.Flush()
	if err != nil {
		return
	}
	err = stdin.Close()
	if err != nil {
		return
	}

	return cmd.Wait()

}

// DeleteContainer deletes a container associated with a given ID or name.
func (docker *Docker) DeleteContainer(id string) error {

	// docker run -f <id>
	cmd := docker.command("docker", "rm", "-f", id)
	return cmd.Run()

}
