//
// storage.go
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
	"net/url"
	"os"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	storage "google.golang.org/api/storage/v1"
)

type Storage struct {
	service *storage.Service
}

func NewStorage(ctx context.Context) (res *Storage, err error) {

	// Create a client.
	client, err := google.DefaultClient(ctx, gcsScope)
	if err != nil {
		return
	}
	// Create a servicer.
	service, err := storage.New(client)
	if err != nil {
		return
	}

	res = &Storage{
		service: service,
	}
	return

}

// NextScript retrieves the oldest script file in a given Queue.
// If there are no items in the queue, it returns nil.
func (s *Storage) NextScript(queue *url.URL) (target *storage.Object, err error) {

	res, err := s.service.Objects.List(queue.Host).Prefix(queue.Path[1:]).Do()
	if err != nil {
		return
	}

	var targetTime time.Time
	for _, item := range res.Items {
		t, _ := time.Parse(TimeFormat, strings.Split(item.TimeCreated, ".")[0])

		if targetTime.IsZero() || targetTime.After(t) {
			target = item
			targetTime = t
		}

	}
	return

}

// Download downloads a given item to a given output file path.
func (s *Storage) Download(item *storage.Object, output string) (err error) {

	res, err := s.service.Objects.Get(item.Bucket, item.Name).Download()
	if err != nil {
		return
	}
	defer res.Body.Close()

	fp, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer fp.Close()

	reader := bufio.NewReader(res.Body)
	_, err = reader.WriteTo(fp)
	return

}

// Delete deletes a given item.
func (s *Storage) Delete(item *storage.Object) error {

	return s.service.Objects.Delete(item.Bucket, item.Name).Do()

}
