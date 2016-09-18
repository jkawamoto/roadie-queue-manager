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
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota

	gcsScope = storage.DevstorageFullControlScope

	// TimeFormat: 2016-07-08T02:05:14.446Z
	TimeFormat = "2006-01-02T15:04:05"
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

	if flags.NArg() != 1 {
		fmt.Println("Usage: roadie-queue-manager <queue URL>")
		return ExitCodeError
	}

	if err := run(flags.Arg(0)); err != nil {
		fmt.Println(err.Error())
		return ExitCodeError
	}
	return ExitCodeOK
}

func run(queue string) error {

	queueURL, err := url.Parse(queue)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Create a client.
	client, err := google.DefaultClient(ctx, gcsScope)
	if err != nil {
		return err
	}
	// Create a servicer.
	service, err := storage.New(client)
	if err != nil {
		return err
	}

	// Search the oldest item in the queue.
	res, err := service.Objects.List(queueURL.Host).Prefix(queueURL.Path[1:]).Do()
	if err != nil {
		return err
	}

	var target *storage.Object
	var targetTime time.Time
	for _, item := range res.Items {
		t, _ := time.Parse(TimeFormat, strings.Split(item.TimeCreated, ".")[0])

		if targetTime.IsZero() || targetTime.After(t) {
			target = item
			targetTime = t
		}

	}

	// Download the oldest item.
	res2, err := service.Objects.Get(queueURL.Host, target.Name).Download()
	if err != nil {
		return err
	}
	defer res2.Body.Close()

	fp, err := ioutil.TempFile(".", "test")
	if err != nil {
		return err
	}
	defer fp.Close()

	reader := bufio.NewReader(res2.Body)
	_, err = reader.WriteTo(fp)
	if err != nil {
		return err
	}

	return nil

}
