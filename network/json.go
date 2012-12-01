package network

import "errors"

type request struct {
	Val string `json:"request"`
}

type keyRequest struct {
	Val string `json:"key"`
}

type response struct {
	Val  string `json:"response"`
	Info string `json:"info,omitempty"`
}

func responseError(res response) error {
	if res.Info != "" {
		return errors.New("network: " + res.Info)
	}
	return errors.New("network: erroneous response")
}
