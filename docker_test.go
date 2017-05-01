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
// along with Roadie queue manager. If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jkawamoto/roadie/script"
)

func TestDockerfile(t *testing.T) {

	s := &script.Script{
		APT: []string{
			"package1",
			"package2",
		},
	}

	buf, err := Dockerfile(s)
	if err != nil {
		t.Fatal(err.Error())
	}
	res := string(buf)

	if !strings.Contains(res, fmt.Sprintf("FROM %v", DefaultImage)) {
		t.Error("Created Dockerfile doesn't use the correct base image")
	}
	if !strings.Contains(res, "apt-get install -y package1") {
		t.Error("Created Dockerfile doesn't install package1")
	}
	if !strings.Contains(res, "apt-get install -y package2") {
		t.Error("Created Dockerfile doesn't install package2")
	}

}
