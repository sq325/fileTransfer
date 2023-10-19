package endpoint

import (
	"context"
	"strings"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func InstrumentingMiddleware(m CommonMetrics) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (resp any, err error) {
			req := request.(GetRequest)
			method, ok := ctx.Value("method").(string)
			if !ok {
				method = "unknow"
			}

			uri, ok := ctx.Value("uri").(string)
			if !ok {
				uri = "unknow"
			}

			defer func() {
				if err != nil {
					m.ReqErrC.With("method", method, "uri", uri).Add(1)
				}
				m.ReqC.With("method", method, "uri", uri).Add(1)
			}()

			defer func(begin time.Time) {
				m.ReqL.With("method", method, "uri", uri).Observe(float64(time.Since(begin).Microseconds()))
			}(time.Now())

			return next(ctx, req)
		}
	}
}

func CardBackInstrumentingMiddleware(m CardBackMetrics) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (resp any, err error) {
			req := request.(GetRequest)
			var rFileName string
			{
				rPath := req.RemoteFilePath
				if rPath != "" {
					for i := len(rPath) - 1; i >= 0; i-- {
						if rPath[i] == '/' {
							rFileName = rPath[i+1:]
							break
						}
					}
				}
			}
			host, user, date, bank := getInfo(rFileName)

			defer func() {
				if err != nil {
					return
				}
				m.FileAllC.With("host", host).Add(1)
				m.FileMainC.With("host", host, "user", user, "date", date, "bank", bank).Add(1)
			}()

			return next(ctx, req)

		}
	}
}

type CommonMetrics struct {
	ReqC    metrics.Counter
	ReqErrC metrics.Counter
	ReqL    metrics.Histogram
}

type CardBackMetrics struct {
	FileAllC  metrics.Counter
	FileMainC metrics.Counter
}

func NewMetrics() CommonMetrics {
	requestCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "fileTransfer_request_count",
		Help: "Number of requests received.",
	}, []string{"uri", "method"})
	requestErrCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "fileTransfer_request_error_count",
		Help: "Number of erroneous requests received.",
	}, []string{"uri", "method"})
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Name: "fileTransfer_request_latency_microseconds",
		Help: "Total duration of requests in microseconds.",
	}, []string{"uri", "method"})
	return CommonMetrics{
		ReqC:    requestCounter,
		ReqErrC: requestErrCounter,
		ReqL:    requestLatency,
	}
}

func NewCardBackMetrics() CardBackMetrics {
	fileAllCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "file_transfer_all_count",
		Help: "Number of tranfered files",
	}, []string{"host"})
	fileMainCounter := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Name: "file_transfer_main_count",
		Help: "Number of tranfered files",
	}, []string{"host", "user", "date", "bank"})

	return CardBackMetrics{
		FileAllC:  fileAllCounter,
		FileMainC: fileMainCounter,
	}
}

func getInfo(fileName string) (host, user, date, bank string) {
	s := strings.Split(fileName, "__")
	if len(s) == 4 { // 新fileback脚本
		host = s[0]
		{
			path := s[1]
			left := strings.IndexByte(path, '_')
			right := 0
			for i := left + 1; i > 0 && i < len(path); i++ {
				if path[i] == '_' {
					right = i
					break
				}
			}
			if right == 0 {
				user = ""
			} else {
				user = path[left+1 : right]
			}
		}
		date = s[2]
		bank = s[3]
	}

	if len(s) == 3 { // 老fileback脚本
		host = s[0]
		date = s[1]
		user = s[2][:strings.IndexByte(s[2], '_')]
		bank = ""
	}
	return
}
