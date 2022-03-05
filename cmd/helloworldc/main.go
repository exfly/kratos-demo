package main

import (
	"context"
	"os"
	"time"

	hellowordv1 "helloworld/api/helloworld/v1"

	"helloworld/pkg/tracer"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	grpcx "google.golang.org/grpc"
)

func main() {
	logger := log.With(log.NewStdLogger(os.Stdout))

	tp, err := tracer.NewTracerProvider("GrpcClient", "http://localhost:14268/api/traces")
	if err != nil {
		logger.Log(log.LevelFatal, "msg", "NewTracerProvider failed")
	}
	_ = tp

	ctx := context.Background()

	conn, err := grpc.DialInsecure(ctx,
		grpc.WithEndpoint("127.0.0.1:9000"),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
		),
		grpc.WithTimeout(2*time.Second),
		// for tracing remote ip recording
		grpc.WithOptions(grpcx.WithStatsHandler(&tracing.ClientHandler{})),
	)
	if err != nil {
		logger.Log(log.LevelFatal, "msg", "grpc DialInsecure failed")
	}

	c := hellowordv1.NewGreeterClient(conn)

	reply, err := c.SayHello(ctx, &hellowordv1.HelloRequest{Name: "yyy"})
	if err != nil {
		logger.Log(log.LevelFatal, "msg", "grpc call SayHello failed")
	}

	logger.Log(log.LevelInfo, "msg", reply)

	err = tp.(*tracesdk.TracerProvider).ForceFlush(ctx)
	if err != nil {
		logger.Log(log.LevelFatal, "msg", "ForceFlush failed")
	}
}
