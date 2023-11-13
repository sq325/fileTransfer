package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sq325/fileTransfer/pkg/endpoint"
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

func DecodeHealthCheckRequest(ctx context.Context, r *http.Request) (any, error) {
	return endpoint.HealthCheckRequest{}, nil
}

func EncodeHealthCheckResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodePutRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.PutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodePutResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodeDownloadRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeDownloadResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}
