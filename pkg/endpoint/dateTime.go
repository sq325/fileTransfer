package endpoint

import (
	"context"
	"fileTransfer/pkg/service"

	"github.com/go-kit/kit/endpoint"
)

type DateTimeRequest struct {
	Duration string `json:"duration" example:"-5d3h22m11s"`
	Layout   string `json:"layout" example:"2006-01-02 15:04:05"`
}

type DateTimeResponse struct {
	V   string `json:"v"`
	Err string `json:"err"`
}

// DateTime
//
//	@Summary		获取制定时间
//	@Description 支持s,m,h,d,w, 例如 -5d：表示5天前。5d：表示5天后
//	@Tags			DateTime
//	@Accept			json
//	@Produce		json
//	@Param			body	body		DateTimeRequest	true	"时间间隔和格式"
//	@Success		200		{object}	DateTimeResponse
//	@Router			/datetime [post]
func MakeDateTimeEndpoint(t service.DateTimer) endpoint.Endpoint {
	return func(ctx context.Context, request any) (any, error) {
		req := request.(DateTimeRequest)
		v, err := t.DateTime(req.Duration, req.Layout)
		if err != nil {
			return DateTimeResponse{"", err.Error()}, nil
		}
		return DateTimeResponse{v, ""}, nil
	}
}
