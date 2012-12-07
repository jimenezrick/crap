package network

import "errors"

type request struct {
	Val string `json:"request,omitempty"`
	Key string `json:"key,omitempty"`
	Size int64 `json:"size,omitempty"`
	Sync bool `json:"sync,omitempty"`
}

type response struct {
	Val  string `json:"response"`
	Info string `json:"info,omitempty"`
}

func responseError(res response) error {
	if res.Info != "" {
		return errors.New("network: " + res.Info)
	}
	return errors.New("network: " + res.Val + " response received")
}
