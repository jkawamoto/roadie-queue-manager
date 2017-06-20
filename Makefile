#
# Makefile
#
# Copyright (c) 2016-2017 Junpei Kawamoto
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
VERSION := $(subst v,,$(shell git describe --abbrev=0 --tags))
.PHONY: build release get-deps test
default: build

asset: get-deps
	rm -f assets/assets.go
	go-bindata -pkg assets -o assets/assets.go -nometadata assets/*

build: asset test
	mkdir -p pkg/$(VERSION)/roadie-queue-manager_$(VERSION)_linux_amd64
	GOOS=linux GOARCH=amd64 go build -o pkg/$(VERSION)/roadie-queue-manager_$(VERSION)_linux_amd64/roadie-queue-manager
	cd pkg/$(VERSION) && tar -zcvf roadie-queue-manager_$(VERSION)_linux_amd64.tar.gz roadie-queue-manager_$(VERSION)_linux_amd64
	rm -r pkg/$(VERSION)/roadie-queue-manager_$(VERSION)_linux_amd64

release: build
	ghr -u jkawamoto v$(VERSION) pkg/$(VERSION)

get-deps:
	go get -d -t -v .
	go get -u github.com/jteeuwen/go-bindata/...

test:
	go test -v ./...
