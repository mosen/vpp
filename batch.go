package vpp

import "errors"

// BatchResult represents a paged result
type BatchResult interface {
	Results() ([]interface{}, error)
	HasNext() bool
}

type batchResult struct {
	*BatchRequestOpts
}

// BatchRequestOpts is a structure that defines options available when fetching objects in batches (users/licenses)
type BatchRequestOpts struct {
	BatchToken         string `json:"batchToken,omitempty"`
	SinceModifiedToken string `json:"sinceModifiedToken,omitempty"`
}

// BatchOption describes the signature of a closure returned if you are adding a batch option to a request
type BatchOption func(*BatchRequestOpts) error

// BatchToken is an argument given to GetUsers when fetching many records in batches
func BatchToken(batchToken string) BatchOption {
	return func(opts *BatchRequestOpts) error {
		if batchToken == "" {
			return errors.New("no batch token given")
		}
		opts.BatchToken = batchToken
		return nil
	}
}

// SinceModifiedToken is an argument given to GetUsers when fetching users that have changed since the last query
func SinceModifiedToken(sinceModifiedToken string) BatchOption {
	return func(opts *BatchRequestOpts) error {
		if sinceModifiedToken == "" {
			return errors.New("no since modified token given")
		}
		opts.SinceModifiedToken = sinceModifiedToken
		return nil
	}
}
