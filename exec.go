//
// exec.go
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
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jkawamoto/roadie-azure/roadie"
	"github.com/jkawamoto/roadie/script"
)

var (
	// RegexpDropboxURL defines a regular expression of a dropbox URL.
	RegexpDropboxURL = regexp.MustCompile(`dropbox://(?:www.dropbox.com/)?(sh?)/([^?]+)(?:\?[^:]+)?(:.*)?`)
	// DefaultDropboxArchive defines a default file name for archive files downloaded from dropbox.
	DefaultDropboxArchive = "dropbox.zip"
)

// ExecuteScript creates a sandbox container and runs a given script in the container.
func ExecuteScript(ctx context.Context, s *script.Script, logger *log.Logger) (err error) {

	logger.Println("Creating a Dockerfile and an entrypoint.sh")
	dockerfile, err := Dockerfile(s)
	if err != nil {
		return
	}
	entrypoint, err := Entrypoint(newEntrypointOpt(s))
	if err != nil {
		return
	}

	cli, err := roadie.NewDockerClient(logger)
	if err != nil {
		return
	}
	defer cli.Close()

	err = cli.Build(ctx, &roadie.DockerBuildOpt{
		ImageName:  s.Name,
		Dockerfile: dockerfile,
		Entrypoint: entrypoint,
	})
	if err != nil {
		return
	}

	err = cli.Start(ctx, s.Name, nil)
	if err != nil {
		return
	}

	return
}

// newEntrypointOpt creates a new set of options from a given script.
func newEntrypointOpt(s *script.Script) (opt *EntrypointOpt) {
	opt = new(EntrypointOpt)

	// Parse source section
	switch {
	case strings.HasSuffix(s.Source, ".git"):
		opt.Git = s.Source
	case strings.HasPrefix(s.Source, "gs://"):
		opt.GSFiles = append(opt.GSFiles, parseURL(s.Source))
	default:
		opt.Downloads = append(opt.Downloads, parseURL(s.Source))
	}

	// Parse data section
	for _, u := range s.Data {
		if strings.HasPrefix(u, "gs://") {
			opt.GSFiles = append(opt.GSFiles, parseURL(u))
		} else {
			opt.Downloads = append(opt.Downloads, parseURL(u))
		}
	}

	// Set running commands.
	opt.Run = s.Run

	// Parse result and upload section
	opt.Result = s.Result
	if !strings.HasSuffix(opt.Result, "/") {
		opt.Result += "/"
	}
	opt.Uploads = s.Upload

	return
}

// parseURL parses an extended URL and returns a download option.
func parseURL(u string) (opt DownloadOpt) {

	var noExpand bool
	if strings.HasPrefix(u, "dropbox://") {
		// The given URL has schema dropbox://
		m := RegexpDropboxURL.FindStringSubmatch(u)

		basename := DefaultDropboxArchive
		if m[1] == "s" {
			basename = filepath.Base(m[2])
		}

		opt.Src = fmt.Sprintf("https://www.dropbox.com/%v/%v?dl=1", m[1], m[2])
		if m[3] != "" {
			// The given URL has a renaming option.
			if strings.HasSuffix(m[3], "/") {
				opt.Dest = filepath.Join(strings.TrimPrefix(m[3], ":"), basename)

			} else {
				opt.Dest = strings.TrimPrefix(m[3], ":")
				noExpand = true
			}

		} else {
			// The given URL doesn't have a renaming option.
			opt.Dest = basename
		}

	} else {
		// The given URL has other schemae.
		lhs := strings.Index(u, ":")
		rhs := strings.LastIndex(u, ":")
		if lhs == rhs {
			// The given URL doesn't have a renaming option.
			opt.Src = u
			opt.Dest = filepath.Base(opt.Src)

		} else {
			// The given URL has a renaming option.
			opt.Src = u[:rhs]
			target := u[rhs+1:]
			if strings.HasSuffix(target, "/") {
				opt.Dest = filepath.Join(target, filepath.Base(opt.Src))
			} else {
				noExpand = true
				opt.Dest = target
			}

		}

	}

	if !noExpand {
		switch {
		case strings.HasSuffix(opt.Dest, ".zip"):
			opt.Zip = true
		case strings.HasSuffix(opt.Dest, ".tar.gz"):
			opt.TarGz = true
		case strings.HasSuffix(opt.Dest, ".tar"):
			opt.Tar = true
		}
	}

	return

}
