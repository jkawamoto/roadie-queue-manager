//
// queue.go
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

	"cloud.google.com/go/datastore"
	"google.golang.org/api/iterator"

	"github.com/jkawamoto/roadie/cloud/gce"
	"github.com/jkawamoto/roadie/script"
)

// RecieveTask recieves tasks in a queue and pass it to the handler.
// If the handler returns non nil error, this function will stop and return the
// the same error.
func RecieveTask(ctx context.Context, project, queue string, handler func(*script.Script) error) (err error) {

	// Create a client.
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		return
	}
	defer client.Close()

	// Query for taking one task in a queue.
	query := datastore.NewQuery(gce.QueueKind).Filter("QueueName=", queue).Filter("Pending=", false).Limit(1)

	// Execute requests.
	for {

		var task gce.Task
		_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) (err error) {

			iter := client.Run(ctx, query)
			key, err := iter.Next(&task)
			if err != nil {
				return
			}

			// Delete a recieved script from the queue.
			err = tx.Delete(key)
			if err != nil {
				return
			}

			// Commit and return.
			return

		})
		if err != iterator.Done {
			return nil
		} else if err != nil {
			return
		}

		// Call the handler.
		err = handler(task.Script)
		if err != nil {
			return
		}

	}

}
