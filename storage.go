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
