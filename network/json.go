package network

import "errors"

type request struct {
	Val string `json:"request"`
}

type keyRequest struct {
	Val string `json:"key"`
}

type result struct {
	Val  string `json:"result"`
	Info string `json:"info,omitempty"`
}

func resultError(res result) error {
	if res.Info != "" {
		return errors.New("network: " + res.Info)
	}
	return errors.New("network: erroneous result")
}
