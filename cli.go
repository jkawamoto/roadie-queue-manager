//
// cli.go
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
// along with Roadie queue manager. If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jkawamoto/roadie/cloud/gce"
	"github.com/jkawamoto/roadie/script"

	yaml "gopkg.in/yaml.v2"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota

	// TimeFormat: 2016-07-08T02:05:14.446Z
	TimeFormat = "2006-01-02T15:04:05"

	// ScriptDir is the directory downloaded script file are stored.
	ScriptDir = "/root"
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
		fmt.Println("Usage: roadie-queue-manager <project id> <queue name>")
		return ExitCodeError
	}

	if err := run(flags.Arg(0), flags.Arg(1)); err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}
	return ExitCodeOK
}

func run(project, queue string) (err error) {
	logger := log.New(os.Stdout, "", 0)

	defer func() (err error) {

		// The context used in this function may be canceled when the following defer
		// function is called; a new background context is thereby used here.
		ctx := context.Background()
		hostname, err := Hostname(ctx)
		if err != nil {
			logger.Println("Cannot retrieve the hostname and stop this instace:", err.Error())
			return
		}
		zone, err := Zone(ctx)
		if err != nil {
			logger.Println("Cannot retrieve zone name and stop this instance:", err.Error())
			return
		}
		instanceID := strings.Split(hostname, ".")[0]

		logger.Println("Deleting instance", instanceID)
		cService := gce.NewComputeService(&gce.GcpConfig{
			Project: project,
			Zone:    zone,
		}, logger)
		err = cService.DeleteInstance(ctx, instanceID)
		if err != nil {
			logger.Println("Cannot stop instance", instanceID, ":", err.Error())
		}
		return

	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Check a script file exists.
	// If there are script files, it measn VM was restarted during the run.
	// The script file must be restarted.
	logger.Println("Checking unfinished tasks")
	matches, err := filepath.Glob(filepath.Join(ScriptDir, "*.yml"))
	if err != nil {
		logger.Println("Failed retrieving unfinished tasks:", err.Error())

	} else {

		for _, filename := range matches {
			logger.Println("Find an unfinished task", filename, "and run it again")

			var s *script.Script
			s, err = script.NewScript(filename)
			if err != nil {
				logger.Println("Cannot read", filename, "and skip it:", err.Error())
				continue
			}

			err = ExecuteScript(ctx, s, logger)
			if err != nil {
				logger.Println("Cannot finish task", filename, ":", err.Error())
				continue
			}
			os.Remove(filename)

		}

	}

	// Start checking queue and executing each script.
	logger.Println("Requesting a task from queue", queue)
	qService, err := gce.NewQueueService(ctx, &gce.GcpConfig{
		Project: project,
	}, logger)
	if err != nil {
		logger.Println("Cannot create a queue service:", err.Error())
		return
	}

	var task *gce.Task
	for {
		task, err = qService.Fetch(ctx, queue)
		if err != nil {
			logger.Println("Cannot fetch any tasks:", err.Error())
			return
		} else if task == nil {
			logger.Println("No tasks are found")
			return
		}

		logger.Println("Recieved a task", task.Name, "from queue", queue)

		// Store a given script into a file so that if this program will be stopped accidentaly,
		// the given script won't be lost.
		var raw []byte
		raw, err = yaml.Marshal(task.Script)
		if err != nil {
			logger.Println("Cannot marshal the task", task.Name, "but can continue processing:", err.Error())
		} else {
			path := filepath.Join(ScriptDir, fmt.Sprintf("%s.yml", task.Name))
			err = ioutil.WriteFile(path, raw, 0644)
			if err != nil {
				logger.Println("Cannot store the task", task.Name, "but can continue processing:", err.Error())
			}
		}

		// Execute a script.
		err = ExecuteScript(ctx, task.Script, log.New(os.Stdout, fmt.Sprintf("task-%v:", task.Name), 0))
		if err != nil {
			logger.Println("Failed to execute task", task.Name, ":", err.Error())
		}
		qService.DeleteTask(ctx, queue, task.Name)

	}

}
