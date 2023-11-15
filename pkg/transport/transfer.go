package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/sq325/fileTransfer/pkg/endpoint"
)

// decode
func DecodeGetRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.GetRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetResponse(_ context.Context, r *http.Response) (any, error) {
	var response endpoint.GetResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	return response, err
}

func DecodeListRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.ListRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeListResponse(_ context.Context, r *http.Response) (any, error) {
	var response endpoint.ListResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	return response, err
}

func DecodeHealthCheckRequest(ctx context.Context, r *http.Request) (any, error) {
	return endpoint.HealthCheckRequest{}, nil
}

func DecodeDownloadRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeDownloadResponse(_ context.Context, r *http.Response) (any, error) {
	var response endpoint.DownloadResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	return response, err
}

func DecodePutRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.PutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodePutResponse(_ context.Context, r *http.Response) (any, error) {
	var response endpoint.PutResponse
	err := json.NewDecoder(r.Body).Decode(&response)
	return response, err
}

func EncodeGetResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeListResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeHealthCheckResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodePutResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}

func EncodeDownloadResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}
