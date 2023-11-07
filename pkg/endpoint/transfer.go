package endpoint

import (
	"context"
	"fileTransfer/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

type GetRequest struct {
	Ip             string `json:"ip"`
	User           string `json:"user"`
	Passwd         string `json:"passwd"`
	RemoteFilePath string `json:"remoteFilePath"`
	LocalPath      string `json:"localPath"`
}

type GetResponse struct {
	V   string `json:"v"`
	Err string `json:"err"`
}

type ListRequest struct {
	Ip             string `json:"ip"`
	User           string `json:"user"`
	Passwd         string `json:"passwd"`
	RemoteFilePath string `json:"remoteFilePath"`
}

type ListResponse struct {
	V   []string `json:"v"`
	Err string   `json:"err"`
}

type HealthCheckRequest struct{}

type HealthCheckResponse struct {
	V bool `json:"v"` // status code
}

// Get
//
//	@Summary		获取远程文件
//	@Description	支持通配符
//	@Tags			GET
//	@Accept			json
//	@Produce		json
//	@Param			body	body		GetRequest	true	"远端ep和文件路径"
//	@Success		200		{object}	GetResponse
//	@Router			/get [post]
func MakeGetEndpoint(t service.Transfer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetRequest)
		err := t.Get(req.Ip, req.User, req.Passwd, req.RemoteFilePath, req.LocalPath)
		if err != nil {
			return GetResponse{req.RemoteFilePath + " failed", err.Error()}, err
		}
		return GetResponse{req.RemoteFilePath + " OK", ""}, nil
	}
}

// List
//
//	@Summary		获取远程文件列表
//	@Description	支持通配符
//	@Tags			GET
//	@Accept			json
//	@Produce		json
//	@Param			body	body		ListRequest	true	"远端ep和文件路径"
//	@Success		200		{object}	ListResponse
//	@Router			/list [post]
func MakeListEndpoint(t service.Transfer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(ListRequest)
		list, err := t.List(req.Ip, req.User, req.Passwd, req.RemoteFilePath)
		if err != nil {
			return ListResponse{list, err.Error()}, err
		}
		return ListResponse{list, ""}, nil
	}
}

func MakeHealthEndpoint(t service.Transfer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		status := t.HealthCheck()
		return HealthCheckResponse{V: status}, nil
	}
}
