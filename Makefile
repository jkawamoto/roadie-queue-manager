#
# Makefile
#
# Copyright (c) 2016 Junpei Kawamoto
#
# This file is part of Roadie queue manager.
#
# Roadie Queue Manager is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Roadie Queue Manager is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Foobar.  If not, see <http:#www.gnu.org/licenses/>.
#
VERSION = snapshot

default: build

.PHONY: build
build:
	goxc -d=pkg -pv=$(VERSION) -os="linux"

.PHONY: release
release:
	ghr -u jkawamoto  v$(VERSION) pkg/$(VERSION)

.PHONY: get-deps
get-deps:
	go get -d -t -v .

.PHONY: test
test:
	go test -v ./...
