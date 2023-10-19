package transport

import (
	"context"
	"encoding/json"
	"fileTransfer/pkg/endpoint"
	"net/http"
)


func DecodeDateTimeRequest(ctx context.Context, r *http.Request) (any, error) {
	var request endpoint.DateTimeRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func EncodeDateTimeResponse(ctx context.Context, w http.ResponseWriter, response any) error {
	return json.NewEncoder(w).Encode(response)
}
