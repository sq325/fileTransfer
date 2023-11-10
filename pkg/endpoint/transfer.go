package endpoint

import (
	"context"
	"fmt"

	"github.com/sq325/fileTransfer/pkg/service"

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

type PutRequest struct {
	ClientUser   string `json:"clientUser"`
	ClientPasswd string `json:"clientPasswd"`
	ClientDir    string `json:"clientDir"` // must abs path
	SrcFilePath  string `json:"srcFilePath"`
}

type PutResponse struct {
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
//	@Summary 从remote端拉取文件到server端
//	@Description	支持通配符
//	@Tags			GET
//	@Accept			json
//	@Produce		json
//	@Param			body	body		GetRequest	true	"远端ep和文件路径"
//	@Success		200		{object}	GetResponse
//	@Router			/get [post]
func MakeGetEndpoint(t service.Transfer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		clientIp, ok := ctx.Value("clientip").(string)
		if !ok {
			return nil, fmt.Errorf("failed to get client ip")
		}
		req := request.(GetRequest)
		err := t.Get(clientIp, req.User, req.Passwd, req.RemoteFilePath, req.LocalPath)
		if err != nil {
			return GetResponse{req.RemoteFilePath + " failed", err.Error()}, err
		}
		return GetResponse{req.RemoteFilePath + " OK", ""}, nil
	}
}

// Put
//
//	@Summary		从server端put文件到client端
//	@Description	支持通配符
//	@Tags			Put
//	@Accept			json
//	@Produce		json
//	@Param			body	body		PutRequest	true	"client and server file"
//	@Success		200		{object}	PutResponse
//	@Router			/put [post]
func MakePutEndpoint(t service.Transfer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(PutRequest)
		clientIp, ok := ctx.Value("ClientIp").(string)
		if !ok {
			return nil, fmt.Errorf("failed to get client ip")
		}
		err := t.Put(clientIp, req.ClientUser, req.ClientPasswd, req.ClientDir, req.SrcFilePath)
		if err != nil {
			return PutResponse{"put" + req.SrcFilePath + " failed", err.Error()}, err
		}
		return PutResponse{"put" + req.SrcFilePath + " OK", ""}, nil
	}
}

// List
//
//	@Summary		从remote端获取文件列表
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
