// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package transport

type FindResponses struct {
	Results   []FindResponse `json:"results,omitempty"`
	ErrorList APIErrors      `json:"error-list,omitempty"`
}

type FindResponse struct {
	Type           string     `json:"type"`
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Entity         Entity     `json:"result,omitempty"`
	DefaultRelease ChannelMap `json:"default-release,omitempty"`
}
