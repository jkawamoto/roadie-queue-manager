package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// DockerRequester defines necessary methods to access to docker.
type DockerRequester interface {
	GetID(name string) (res string, err error)
	DeleteContainer(id string) error
}

// Docker defines an inteface of docker API.
type Docker struct {
}

// NewDocker returns a new docker interface.
func NewDocker() *Docker {
	return &Docker{}
}

// GetID returns the ID associated with a given container name.
func (docker *Docker) GetID(name string) (res string, err error) {

	cmd := exec.Command("docker", "ps", "-aq", "-f", fmt.Sprintf("name=%s", name))
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
func (docker *Docker) CreateContainer(image, name, script string) (err error) {

	// docker run -i --name {{.Name}} {{.Image}} < run.yml
	cmd := exec.Command("docker", "run", "-i", "--name", name, image, "--no-shutdown")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return
	}
	writer := bufio.NewWriter(stdin)

	err = cmd.Start()
	if err != nil {
		return
	}

	_, err = writer.WriteString(script)
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

	cmd := exec.Command("docker", "rm", "-f", id)
	return cmd.Run()

}
