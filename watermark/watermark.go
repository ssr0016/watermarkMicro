package watermark

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/velotiotech/watermark-service/api/v1/pb/watermark"
	"github.com/velotiotech/watermark-service/pkg/watermark"
	"github.com/velotiotech/watermark-service/pkg/watermark/endpoints"
	"github.com/velotiotech/watermark-service/pkg/watermark/transport"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/oklog/oklog/pkg/group"
	"google.golang.org/grpc"
)

const (
	defaultHTTPort = "8080"
	defaultGRPCPort = "9090"
)

func main () {
	var (
		logger log.Logger
		httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHTTPPort))
		grpcAddr = net.JoinsHostPort("localhost", envString("GRPC_PORT", defaultGRPCPort))
	)
}

logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
logger = log.With(logger, "ts", log.DefaultTimestampUTC)

var (
	service = watermark.NewService()
	eps = endpoints.NewEndpointSet(service)
	httpHandler = transport.NewHTTPHandler(eps)
	grpcServer = transport.NewGRPCServer(eps)
)

var g group.group
{
	// The HTTP listener mounts the Go kit HTTP handler we created.
	httpListener, err := net.Listen("tcp", httpAddr)
	if err != nil {
		logger.Log("transport", "HTTP". "during", "Listen", "err", err)
		os.Exit(1)
	}
	g.Add(func()) error {
		logger.log("transport", "HTTP", "addr", httAddr)
		return http.Serve(httpListerner, httpHandler)
	}, func(error) {
		httpListener.Close()
	}
}
{
	// The gRPC listener mounts the Go kit gRPC server we created.
	grpcListener, err := net.Listen("tcp",grpcAdddr)
	if err != nil {
		logger.log("transport", "gRPC", "during", "Listen", "err", err)
		os.Exit(1)
	}
	g.Add(func() error{
		if err != nil {
			logger.("transport", "gRPC", "during", "Listen", "err", err)
			// we add the Go Kit gRPC Interceptor to our gRPC service as it is used by
			// the here demonstrated zipkin tracing middleware.
			baseServer := grpc.NewServer(grpc.UnaryInterceptor(kitgrpc.Interceptor))
			pb.RegisterWatermarkServer(baseServer, grpcServer)
			return baseServer.Serve(grpcLister)
		}, func(error){
			grpcListener.Close()
		}
	})
}
{
	// This function just sits and waits for ctrl-C.
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select{
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}), func(error){
		close(cancelInterrupt)
	}
	logger.Log("exit", g.Run())
}

	func envString(evn, fallback string) string {
		e := os.Getenv(env)
		if e == ""{
			return fallback
		}
		return e
	}



