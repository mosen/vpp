package vpp

// BatchRequestOpts is a structure that defines options available when fetching objects in batches (users/licenses)
type BatchRequest struct {
	BatchToken         string `json:"batchToken,omitempty"`
	SinceModifiedToken string `json:"sinceModifiedToken,omitempty"`
}

func (br *BatchRequest) HasNext() bool {
	return br.BatchToken != "" && br.SinceModifiedToken == ""
}

// BatchOption describes the signature of a closure returned if you are adding a batch option to a request
//type BatchOption func(*BatchRequest) error

// BatchToken is an argument given to GetUsers when fetching many records in batches
//func BatchToken(batchToken string) BatchOption {
//	return func(opts *BatchRequest) error {
//		if batchToken == "" {
//			return errors.New("no batch token given")
//		}
//		opts.BatchToken = batchToken
//		return nil
//	}
//}

// SinceModifiedToken is an argument given to GetUsers when fetching users that have changed since the last query
//func SinceModifiedToken(sinceModifiedToken string) BatchOption {
//	return func(opts *BatchRequest) error {
//		if sinceModifiedToken == "" {
//			return errors.New("no since modified token given")
//		}
//		opts.SinceModifiedToken = sinceModifiedToken
//		return nil
//	}
//}
