//
// script_test.go
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
	"io/ioutil"
	"os"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

var sampleScript = &QueuedScript{
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

// TestNewQueuedScript tests NewQueuedScript loads a script file.
func TestNewQueuedScript(t *testing.T) {

	// Prepare testing.
	var err error
	fp, err := ioutil.TempFile("", "queued-script-")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(fp.Name())
	writer := bufio.NewWriter(fp)

	raw, err := yaml.Marshal(&sampleScript)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = writer.Write(raw)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = writer.Flush()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = fp.Close()
	if err != nil {
		t.Fatal(err.Error())
	}

	// Start testing.
	res, err := NewQueuedScript(fp.Name())
	if err != nil {
		t.Error("Cannot load the script:", err.Error())
	}
	if res.InstanceName != sampleScript.InstanceName {
		t.Error("Loaded queued script isn't same as original one:", res)
	}
	if res.Image != sampleScript.Image {
		t.Error("Loaded queued script isn't same as original one:", res)
	}
	if res.Body.Result != sampleScript.Body.Result {
		t.Error("Loaded queued script isn't same as original one:", res)
	}

}

// TestScript tests generated scripts have correct values.
func TestScript(t *testing.T) {

	raw, err := sampleScript.Script()
	if err != nil {
		t.Fatal(err.Error())
	}
	script := string(raw)
	if !strings.Contains(script, sampleScript.Body.Source) {
		t.Error("Generated script isn't correct:", script)
	}
	if strings.Contains(script, sampleScript.InstanceName) {
		t.Error("Generated script isn't correct:", script)
	}

}
