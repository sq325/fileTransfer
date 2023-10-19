package main

import (
	"context"
	"fileTransfer/pkg/endpoint"
	"fileTransfer/pkg/service"
	"fileTransfer/pkg/transport"
	"fmt"
	"log"
	"net/http"

	_ "fileTransfer/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	kitendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/spf13/pflag"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	port *string = pflag.StringP("port", "p", "8080", "port to listen")

	version *bool = pflag.BoolP("version", "v", false, "show version info")
)

const (
	_version   = "fileTransfer v0.1.2"
	updateInfo = "/get localFilePath bugfix"
)

var (
	buildTime      string
	buildGoVersion string
	author         string
)

func init() {
	pflag.Parse()
}

// @title			运行室 文件传输服务
// @version		0.1.2
//
// @contact.name	Sun Quan
// @contact.email	xxx
// @license.name	Apache 2.0
func main() {
	if *version {
		fmt.Println(_version)
		fmt.Println("build time:", buildTime)
		fmt.Println("go version:", buildGoVersion)
		fmt.Println("author:", author)
		fmt.Println("update info:", updateInfo)
		return
	}

	metrics := endpoint.NewMetrics()
	cardMetrics := endpoint.NewCardBackMetrics()
	instrumentingMiddleware := endpoint.InstrumentingMiddleware(metrics)
	cardinstrumentingMiddleware := endpoint.CardBackInstrumentingMiddleware(cardMetrics)

	transferSvc := service.NewTransfer()
	datetimeSvc := service.NewDateTimer()

	mux := gin.Default()

	mux.POST("/get",
		ginHandlerFunc(
			"POST",
			"/get",
			instrumentingMiddleware(
				cardinstrumentingMiddleware(
					endpoint.MakeGetEndpoint(transferSvc),
				),
			),
			transport.DecodeGetRequest,
			transport.EncodeGetResponse,
		),
	)

	mux.POST("/list",
		ginHandlerFunc(
			"POST",
			"/list",
			instrumentingMiddleware(endpoint.MakeListEndpoint(transferSvc)),
			transport.DecodeListRequest,
			transport.EncodeListResponse,
		),
	)

	mux.POST("/datetime",
		ginHandlerFunc(
			"POST",
			"/datetime",
			instrumentingMiddleware(endpoint.MakeDateTimeEndpoint(datetimeSvc)),
			transport.DecodeDateTimeRequest,
			transport.EncodeDateTimeResponse,
		),
	)
	mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	mux.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Fatal(mux.Run(":" + *port))

}

func ginHandlerFunc(
	method,
	uri string,
	ep kitendpoint.Endpoint,
	dec httptransport.DecodeRequestFunc,
	enc httptransport.EncodeResponseFunc,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		opt := httptransport.ServerBefore(
			func(ctx context.Context, _ *http.Request) context.Context {
				ctx = context.WithValue(ctx, "method", method)
				ctx = context.WithValue(ctx, "uri", uri)
				return ctx
			})

		h := httptransport.NewServer(
			ep,
			dec,
			enc,
			opt,
		)
		h.ServeHTTP(c.Writer, c.Request)
	}
}
