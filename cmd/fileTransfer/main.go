package main

import (
	"fileTransfer/pkg/endpoint"
	"fileTransfer/pkg/service"
	"fileTransfer/pkg/transport"
	"fmt"
	"log"

	_ "fileTransfer/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/spf13/pflag"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/sq325/kitComplement/pkg/instrumentation"
)

var (
	port *string = pflag.StringP("port", "p", "8080", "port to listen")

	version *bool = pflag.BoolP("version", "v", false, "show version info")
)

const (
	_version     = "fileTransfer v0.2.0"
	_versionInfo = "delete datetime and extract instrumentation as individual package"
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
// @version		0.2.0

// @license.name	Apache 2.0
func main() {
	if *version {
		fmt.Println(_version)
		fmt.Println("build time:", buildTime)
		fmt.Println("go version:", buildGoVersion)
		fmt.Println("author:", author)
		fmt.Println("version info:", _versionInfo)
		return
	}

	metrics := instrumentation.NewMetrics()
	instrumentingMiddleware := instrumentation.InstrumentingMiddleware(metrics)

	transferSvc := service.NewTransfer()

	mux := gin.Default()

	mux.POST("/get",
		instrumentation.GinHandlerFunc(
			"POST",
			"/get",
			instrumentingMiddleware(
				endpoint.MakeGetEndpoint(transferSvc),
			),
			transport.DecodeGetRequest,
			transport.EncodeGetResponse,
		),
	)

	mux.POST("/list",
		instrumentation.GinHandlerFunc(
			"POST",
			"/list",
			instrumentingMiddleware(endpoint.MakeListEndpoint(transferSvc)),
			transport.DecodeListRequest,
			transport.EncodeListResponse,
		),
	)

	mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	mux.GET("/metrics", gin.WrapH(promhttp.Handler()))

	log.Fatal(mux.Run(":" + *port))

}
