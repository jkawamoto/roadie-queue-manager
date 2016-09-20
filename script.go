//
// script.go
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
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ScriptBody defines the roadie script body.
type ScriptBody struct {
	// List of apt packages to be installed.
	APT []string `yaml:"apt,omitempty"`
	// URL to the source code.
	Source string `yaml:"source,omitempty"`
	// List of URLs to be downloaded as data files.
	Data []string `yaml:"data,omitempty"`
	// List of commands to be run.
	Run []string `yaml:"run,omitempty"`
	// URL where the computational results will be stored.
	Result string `yaml:"result,omitempty"`
	// List of glob pattern, files matches of one of them are uploaded as resuts.
	Upload []string `yaml:"upload,omitempty"`
}

// QueuedScript defines a data structure of script file enqueued.
type QueuedScript struct {
	// InstanceName to be created.
	InstanceName string `yaml:"instance-name,omitempty"`
	// Image name to be used to create the instance.
	Image string `yaml:"image,omitempty"`
	// The script body.
	Body ScriptBody `yaml:"body,omitempty"`
}

// NewQueuedScript reads a file associated with a given file name and returns
// a QueuedScript.
func NewQueuedScript(file string) (script *QueuedScript, err error) {

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	script = &QueuedScript{}
	err = yaml.Unmarshal(raw, script)
	return

}

// ScriptBody returns the script body of this queued script.
func (q *QueuedScript) ScriptBody() ([]byte, error) {

	return yaml.Marshal(q.Body)

}

// Bytes returns the YAML presentation of this object.
func (q *QueuedScript) Bytes() ([]byte, error) {

	return yaml.Marshal(q)

}
