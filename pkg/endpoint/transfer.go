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
		if req.User == "" {
			req.User, req.Passwd = patrol()
		}

		err := t.Get(req.Ip, req.User, req.Passwd, req.RemoteFilePath, req.LocalPath)
		if err != nil {
			return GetResponse{"", err.Error()}, nil
		}
		return GetResponse{"OK", ""}, nil
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
		if req.User == "" {
			req.User, req.Passwd = patrol()
		}
		list, err := t.List(req.Ip, req.User, req.Passwd, req.RemoteFilePath)
		if err != nil {
			return ListResponse{list, err.Error()}, nil
		}
		return ListResponse{list, ""}, nil
	}
}

func patrol() (user, passwd string) {
	return "patrol", "agent@9753"
}
