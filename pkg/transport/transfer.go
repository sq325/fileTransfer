package transport

import (
	"context"
	"encoding/json"
	"fileTransfer/pkg/endpoint"
	"net/http"
)

func DecodeGetRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.GetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeListRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeGetResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeListResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}
