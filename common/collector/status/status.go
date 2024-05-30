// Package status defines an API for any implementation services that
// wish to act as a status reporting service (server).
//
// The implementation can report literally anything to the caller.
package status

import "context"

// Status interface governs the rules of interaction with an
// underlying status server repository. Data, that is passed in
// req/resp as an interface{} (DTO), is used by implementations
// to build up and then perform queries against the repository.
//
// This means, that implementations should provide DTOs for any
// req/resp.
type Status interface {
	// Read returns a status server response resp based on a request req.
	//
	// Request req can be anything from an SQL query to a NoSQL key.
	Read(ctx context.Context, req interface{}) (resp interface{}, err error)

	// Update returns an error err if request req wasn't updated in a repository.
	//
	// Request req can be anything from an SQL query to a NoSQL key.
	Update(ctx context.Context, req interface{}) (err error)

	// Delete returns an error err if request req wasn't deleted from a repository.
	//
	// Request req can be anything from an SQL query to a NoSQL key.
	Delete(ctx context.Context, req interface{}) (err error)
}
