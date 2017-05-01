//
// metadata.go
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
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/context/ctxhttp"
)

const (
	// ProjectIDMetadataURL defines a metadata URL of the project ID.
	ProjectIDMetadataURL = "http://metadata.google.internal/computeMetadata/v1/project/project-id"
	// InstanceIDMetadataURL defines a metadata URL of the instance ID.
	InstanceIDMetadataURL = "http://metadata.google.internal/computeMetadata/v1/instance/id"
	// HostnameMetadataURL defines a matadata URL of the hostname.
	HostnameMetadataURL = "http://metadata.google.internal/computeMetadata/v1/instance/hostname"
	// ZoneMetadataURL defines a metadata URL of the zone this instance running.
	ZoneMetadataURL = "http://metadata.google.internal/computeMetadata/v1/instance/zone"
)

func getMetadata(ctx context.Context, url string) (str string, err error) {

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("Metadata-Flavor", "Google")
	res, err := ctxhttp.Do(ctx, nil, req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	str = string(data)
	return

}

// ProjectID returns the project ID this instance belonging to.
func ProjectID(ctx context.Context) (string, error) {
	return getMetadata(ctx, ProjectIDMetadataURL)
}

// InstanceID returns the ID of this instance.
func InstanceID(ctx context.Context) (string, error) {
	return getMetadata(ctx, InstanceIDMetadataURL)
}

// Hostname returns the host name of this instance.
func Hostname(ctx context.Context) (string, error) {
	return getMetadata(ctx, HostnameMetadataURL)
}

// Zone returns the zone name this instance running in.
func Zone(ctx context.Context) (zone string, err error) {
	zone, err = getMetadata(ctx, ZoneMetadataURL)
	if err == nil {
		zone = zone[strings.LastIndex(zone, "/")+1:]
	}
	return
}
