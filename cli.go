//
// cli.go
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
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/net/context"

	storage "google.golang.org/api/storage/v1"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota

	gcsScope = storage.DevstorageFullControlScope

	// TimeFormat: 2016-07-08T02:05:14.446Z
	TimeFormat = "2006-01-02T15:04:05"

	// ScriptDir is the directory downloaded script file are stored.
	ScriptDir = "script"
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	if flags.NArg() != 2 {
		fmt.Println("Usage: roadie-queue-manager <project-ID> <queue-name>")
		return ExitCodeError
	}

	if err := run(flags.Arg(0), flags.Arg(1)); err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}
	return ExitCodeOK
}

func run(project, queue string) (err error) {

	// Prepare the directory to store downloaded script files.
	err = os.MkdirAll(ScriptDir, 0755)
	if err != nil {
		return
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Docker accessor.
	docker := NewDocker()

	// Check a script file exists.
	// If there are script files, it measn VM was restarted during the run.
	// The script file must be restarted.
	matches, err := filepath.Glob(filepath.Join(ScriptDir, "*.yml"))
	if err != nil {
		return
	}
	for _, file := range matches {
		err = executeScript(docker, file)
		if err != nil {
			return
		}
	}

	// Start checking queue and executing each script.
	for {

		var path string

		// Obtain a next script.
		err = NextScript(ctx, project, queue, func(script *QueuedScript) error {
			// This handler stores a given script into a file
			// so that if this program will be stopped accidentaly,
			// the given script won't be lost.
			raw, err2 := script.Bytes()
			if err2 != nil {
				return err2
			}
			path = filepath.Join(ScriptDir, script.InstanceName, ".yml")
			return ioutil.WriteFile(path, raw, 0644)

		})
		if err == NoScript {
			break
		} else if err != nil {
			return err
		}

		// Execute a script.
		err = executeScript(docker, path)
		if err != nil {
			return err
		}

	}

	return nil

}

// DockerRequester defines necessary methods to access to docker.
type DockerRequester interface {
	GetID(name string) (res string, err error)
	CreateContainer(image, name string, script []byte) error
	DeleteContainer(id string) error
}

// executeScript runs a given script via a given docker interface.
// When the script ends without errors, the script file will be deleted.
func executeScript(docker DockerRequester, file string) (err error) {

	// Parse the script.
	script, err := NewQueuedScript(file)
	if err != nil {
		return
	}

	// Check a previous container exists.
	id, err := docker.GetID(script.InstanceName)
	if err != nil {
		return
	}
	if len(id) != 0 {
		// If samne name container exists, delete it.
		err = docker.DeleteContainer(id)
		if err != nil {
			return
		}
	}

	// Start a new container.
	body, err := script.ScriptBody()
	if err != nil {
		return
	}
	err = docker.CreateContainer(script.Image, script.InstanceName, body)
	if err != nil {
		return
	}

	// Delete the script file.
	return os.Remove(file)

}
