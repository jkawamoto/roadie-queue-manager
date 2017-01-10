//
// datastore.go
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
// along with Foobar.  If not, see <http://www.gnu.org/licenses/>.
//

package main

import (
	"context"
	"fmt"

	"google.golang.org/api/iterator"

	"github.com/jkawamoto/roadie/command/cloud"
	"github.com/jkawamoto/roadie/command/resource"

	"cloud.google.com/go/datastore"
)

// NoScript is an error type which represents there are no scripts in a queue.
var NoScript = fmt.Errorf("There are no script in the queue.")

// NextScript returns a next script in a queue.
// Founded script will be passed to a given handler. This request is in a transaction.
// If the handler returns any error, transaction will be rollbacked,
// otherwise the founded script will be deleted and the transaction will be commited.
// If there are no scripts, it returns nil without calling the handler.
func NextScript(ctx context.Context, project, queue string, handler func(*resource.Task) error) (err error) {

	// Create a client.
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		return
	}
	defer client.Close()

	// Execute requests.
	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) (err error) {

		query := datastore.NewQuery(cloud.QueueKind).Filter("Pending=", false).Limit(1)
		iter := client.Run(ctx, query)

		var res resource.Task
		key, err := iter.Next(&res)
		if err == iterator.Done {
			// There are no items in the given queue.
			return NoScript
		} else if err != nil {
			return
		}

		// Delete a recieved script from the queue.
		err = tx.Delete(key)
		if err != nil {
			return
		}

		// Call the handler.
		err = handler(&res)
		if err != nil {
			return
		}

		// Commit and return.
		return

	})

	return

}
