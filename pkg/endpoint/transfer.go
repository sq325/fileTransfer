package endpoint

import (
	"context"
	"fmt"

	"github.com/sq325/fileTransfer/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

type GetRequest struct {
	RemoteIp       string `json:"remoteIp"`
	RemoteUser     string `json:"remoteUser"`
	RemotePasswd   string `json:"remotePasswd"`
	RemoteFilePath string `json:"remoteFilePath"` // must abs path
	SrcDir         string `json:"srcDir"`
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
	RemoteIp       string `json:"remoteIp"`
	RemoteUser     string `json:"remoteUser"`
	RemotePasswd   string `json:"remotePasswd"`
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

type DownloadRequest struct {
	RemoteIp       string `json:"remoteIp"`
	RemoteUser     string `json:"remoteUser"`
	RemotePasswd   string `json:"remotePasswd"`
	RemoteFilePath string `json:"remoteFilePath"`
	ClientIp       string `json:"clientIp"`
	ClientUser     string `json:"clientUser"`
	ClientPasswd   string `json:"clientPasswd"`
	ClientDir      string `json:"clientDir"`
}

type DownloadResponse struct {
	V   string `json:"v"`
	Err string `json:"err"`
}

// Get
//
//	@Summary 从remote端拉取文件到server端
//	@Description	支持通配符
//	@Tags			GET
//	@Accept			json
//	@Produce		json
//	@Param			body	body		GetRequest	true	"remote -> server"
//	@Success		200		{object}	GetResponse
//	@Router			/get [post]
func MakeGetEndpoint(t service.Transfer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(GetRequest)
		err := t.Get(req.RemoteIp, req.RemoteUser, req.RemotePasswd, req.RemoteFilePath, req.SrcDir)
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
//	@Param			body	body		PutRequest	true	"server -> client"
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
//	@Param			body	body		ListRequest	true	"remote"
//	@Success		200		{object}	ListResponse
//	@Router			/list [post]
func MakeListEndpoint(t service.Lister) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(ListRequest)
		list, err := t.List(req.RemoteIp, req.RemoteUser, req.RemotePasswd, req.RemoteFilePath)
		if err != nil {
			return ListResponse{list, err.Error()}, err
		}
		return ListResponse{list, ""}, nil
	}
}

// Download
//
//	@Summary		从remote端下载文件到client端
//	@Description	支持通配符
//	@Tags			GET
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DownloadRequest	true	"remote -> client"
//	@Success		200		{object}	DownloadResponse
//	@Router			/download [post]
func MakeDownloadEndpoint(d service.Downloader) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DownloadRequest)
		err := d.Download(req.RemoteIp, req.RemoteUser, req.RemotePasswd, req.RemoteFilePath, req.ClientIp, req.ClientUser, req.ClientPasswd, req.ClientDir)
		if err != nil {
			return DownloadResponse{req.RemoteFilePath + " download failed", err.Error()}, err
		}
		return DownloadResponse{req.RemoteFilePath + " downlaod OK", ""}, nil
	}
}

func MakeHealthEndpoint(t service.Transfer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		status := t.HealthCheck()
		return HealthCheckResponse{V: status}, nil
	}
}
