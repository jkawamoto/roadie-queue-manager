//
// datastore.go
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
	"fmt"

	"cloud.google.com/go/datastore"
	"golang.org/x/net/context"
)

// KindFormat is a format string to generate a kind to a queue.
const KindFormat = "roadie-queue-%s"

// ScriptHandler is a function type which handles a script.
type ScriptHandler func(*QueuedScript) error

// NoScript is an error type which represents there are no scripts in a queue.
var NoScript = fmt.Errorf("There are no script in the queue.")

// NextScript returns a next script in a queue.
// Founded script will be passed to a given handler. This request is in a transaction.
// If the handler returns any error, transaction will be rollbacked,
// otherwise the founded script will be deleted and the transaction will be commited.
// If there are no scripts, it returns nil without calling the handler.
func NextScript(ctx context.Context, project, queue string, handler ScriptHandler) (err error) {

	// Create the kind associated with a given queue.
	kind := fmt.Sprintf(KindFormat, queue)

	// Create a client.
	client, err := datastore.NewClient(ctx, project)
	if err != nil {
		return
	}
	defer client.Close()

	// Create a Transaction.
	trans, err := client.NewTransaction(ctx)
	if err != nil {
		return
	}

	// Execute requests.
	var res QueuedScript
	query := datastore.NewQuery(kind).Transaction(trans).Limit(1)
	iter := client.Run(ctx, query)

	key, err := iter.Next(&res)
	if err == datastore.Done {
		// There are no items in the given queue.
		trans.Commit()
		return nil
	}

	// Delete a recieved script from the queue.
	err = trans.Delete(key)
	if err != nil {
		trans.Rollback()
		return
	}

	// Call the handler.
	err = handler(&res)
	if err != nil {
		trans.Rollback()
		return
	}

	// Commit and return.
	_, err = trans.Commit()
	return

}
