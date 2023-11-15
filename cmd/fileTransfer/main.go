package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/sq325/fileTransfer/pkg/endpoint"
	"github.com/sq325/fileTransfer/pkg/service"
	"github.com/sq325/fileTransfer/pkg/transport"

	_ "github.com/sq325/fileTransfer/docs"

	complementConsul "github.com/sq325/kitComplement/pkg/consul"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/spf13/pflag"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/sq325/kitComplement/pkg/instrumentation"
)

var (
	port   *string = pflag.StringP("port", "p", "8080", "port to listen")
	consul *string = pflag.String("consul", "", "consul endpoint, eg: 10.10.10.10:8500")

	version *bool = pflag.BoolP("version", "v", false, "show version info")
)

const (
	_service     = "fileTransfer"
	_version     = "v0.5.4"
	_versionInfo = "add dec funcs"
)

var (
	buildTime      string
	buildGoVersion string
	author         string
)

func init() {
	pflag.Parse()
}

// @title			文件传输服务
// @version		0.5.4

// @license.name	Apache 2.0
func main() {
	if *version {
		fmt.Println(_service, _version)
		fmt.Println("build time:", buildTime)
		fmt.Println("go version:", buildGoVersion)
		fmt.Println("author:", author)
		fmt.Println("version info:", _versionInfo)
		return
	}

	metrics := instrumentation.NewMetrics()
	instrumentingMiddleware := instrumentation.InstrumentingMiddleware(metrics)

	transferSvc := service.NewTransfer()
	listSvc := service.NewLister()
	downloadSvc := service.NewDownloader()
	healthSvc := func() bool {
		return transferSvc.HealthCheck()
	}

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

	mux.POST("/put",
		instrumentation.GinHandlerFunc(
			"POST",
			"/put",
			instrumentingMiddleware(endpoint.MakePutEndpoint(transferSvc)),
			transport.DecodePutRequest,
			transport.EncodePutResponse,
		),
	)

	mux.POST("/list",
		instrumentation.GinHandlerFunc(
			"POST",
			"/list",
			instrumentingMiddleware(endpoint.MakeListEndpoint(listSvc)),
			transport.DecodeListRequest,
			transport.EncodeListResponse,
		),
	)

	mux.POST("/download",
		instrumentation.GinHandlerFunc(
			"POST",
			"/download",
			instrumentingMiddleware(endpoint.MakeDownloadEndpoint(downloadSvc)),
			transport.DecodeDownloadRequest,
			transport.EncodeDownloadResponse,
		),
	)

	mux.Any("/health",
		func(c *gin.Context) {
			if healthSvc() {
				c.String(200, "%s service is healthy", _service)
				return
			}
			c.String(500, "%s service is unhealthy", _service)
		},
	)

	mux.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	mux.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// register service
	if *consul != "" {
		var r complementConsul.Registrar
		{
			c := strings.Split(strings.TrimSpace(*consul), ":")
			ip := c[0]
			p, _ := strconv.Atoi(c[1])
			consulClient := complementConsul.NewConsulClient(ip, p)
			logger := complementConsul.NewLogger()
			r = complementConsul.NewRegistrar(consulClient, logger)
		}

		var svc *complementConsul.Service
		{
			port, _ := strconv.Atoi(*port)
			svc = &complementConsul.Service{
				Name: _service,
				Port: port,
				Check: struct {
					Path     string
					Interval string
					Timeout  string
				}{
					Path:     "/health",
					Interval: "60s",
					Timeout:  "10s",
				},
			}
		}
		r.Register(svc)
		defer r.Deregister(svc)
	}

	go func() { // no blocking
		err := mux.Run(":" + *port)
		if err != nil {
			log.Println(err)
		}
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
	}()

	// sig capture
	chSig := make(chan os.Signal, 1)
	signal.Notify(chSig, syscall.SIGHUP, syscall.SIGTERM)
	for sig := range chSig { // blocking
		switch sig {
		case syscall.SIGTERM:
			return
		case syscall.SIGHUP:
			continue
		}
	}
}
