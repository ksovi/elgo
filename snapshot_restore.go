// Copyright 2012-present Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.
// snap restore by ovi chis ovi@ovios.org http://ovios.org | Sun Dec  2 18:59:59 EST 2018

package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"github.com/olivere/elastic/uritemplates"
)

// SnapshotRestoreService restores a snapshot from a snapshot repository.
// It is documented at
// https://www.elastic.co/guide/en/elasticsearch/reference/6.4/modules-snapshots.html.
type SnapshotRestoreService struct {
	client     *Client
	repository string
	snapshot   string
	bodyString string
}

// NewSnapshotRestoreService creates a new SnapshotRestoreService.
func NewSnapshotRestoreService(client *Client) *SnapshotRestoreService {
	return &SnapshotRestoreService{
		client: client,
	}
}

// Repository is the repository name.
func (s *SnapshotRestoreService) Repository(repository string) *SnapshotRestoreService {
	s.repository = repository
	return s
}

// Snapshot is the snapshot name.
func (s *SnapshotRestoreService) Snapshot(snapshot string) *SnapshotRestoreService {
	s.snapshot = snapshot
	return s
}

func (s *SnapshotRestoreService) BodyString(body string) *SnapshotRestoreService {
	s.bodyString = body
	return s
}

// buildURL builds the URL for the operation.
func (s *SnapshotRestoreService) buildURL() (string, url.Values, error) {
	// Build URL
	path, err := uritemplates.Expand("/_snapshot/{repository}/{snapshot}/_restore", map[string]string{
		"repository": s.repository,
		"snapshot":   s.snapshot,
	})
	if err != nil {
		return "", url.Values{}, err
	}
	return path, url.Values{}, nil
}

// Validate checks if the operation is valid.
func (s *SnapshotRestoreService) Validate() error {
	var invalid []string
	if s.repository == "" {
		invalid = append(invalid, "Repository")
	}
	if s.snapshot == "" {
		invalid = append(invalid, "Snapshot")
	}
	if len(invalid) > 0 {
		return fmt.Errorf("missing required fields: %v", invalid)
	}
	return nil
}

// Do executes the operation.
func (s *SnapshotRestoreService) Do(ctx context.Context) (*SnapshotRestoreResponse, error) {
	// Check pre-conditions
	if err := s.Validate(); err != nil {
		return nil, err
	}

	// Get URL for request
	path, params, err := s.buildURL()
	if err != nil {
		return nil, err
	}
    var body interface{}
    body = s.bodyString
    
    
	// Get HTTP response
	res, err := s.client.PerformRequest(ctx, PerformRequestOptions{
		Method: "POST",
		Path:   path,
		Params: params,
        Body: body,
	})
	if err != nil {
		return nil, err
	}

	// Return operation response
	ret := new(SnapshotRestoreResponse)
	if err := json.Unmarshal(res.Body, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// SnapshotRestoreResponse is the response of SnapshotRestoreService.Do.
type SnapshotRestoreResponse struct {
	Accepted bool `json:"accepted"`
}
