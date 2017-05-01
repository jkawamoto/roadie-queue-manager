//
// docker.go
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
	"bytes"
	"text/template"

	"github.com/jkawamoto/roadie-queue-manager/assets"
	"github.com/jkawamoto/roadie/script"
)

const (
	// DefaultImage defines the default base image.
	DefaultImage = "ubuntu:latest"
)

// loadTemplate loads a template from assets and apply given options.
func loadTemplate(name string, opt interface{}) (res []byte, err error) {

	data, err := assets.Asset(name)
	if err != nil {
		return
	}

	temp, err := template.New("").Parse(string(data))
	if err != nil {
		return
	}

	buf := bytes.NewBuffer(nil)
	err = temp.Execute(buf, opt)
	if err != nil {
		return
	}

	res = buf.Bytes()
	return

}

// Dockerfile creates a new Dockerfile for a given script.
func Dockerfile(s *script.Script) (res []byte, err error) {

	if s.Image == "" {
		s.Image = DefaultImage
	}
	return loadTemplate("assets/Dockerfile", s)

}

// DownloadOpt defines a download option for a URL.
type DownloadOpt struct {
	Src   string
	Dest  string
	Zip   bool
	TarGz bool
	Tar   bool
}

// EntrypointOpt defines options to create an entrypoint.sh.
type EntrypointOpt struct {
	Git       string
	Downloads []DownloadOpt
	GSFiles   []DownloadOpt
	Run       []string
	Result    string
	Uploads   []string
}

// Entrypoint creates a new entrypoint.sh with a given set of options.
func Entrypoint(opt *EntrypointOpt) (res []byte, err error) {

	return loadTemplate("assets/entrypoint.sh", opt)

}
