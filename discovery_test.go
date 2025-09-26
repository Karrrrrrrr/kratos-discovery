package helloworld

import (
	"context"
	"fmt"
	"testing"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/middleware/circuitbreaker"
	"github.com/go-kratos/kratos/v2/middleware/metadata"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/hashicorp/consul/api"

	v1 "helloworld/api/helloworld/v1"
)

func TestName(t *testing.T) {
	consulClient, err := api.NewClient(&api.Config{
		Address: "127.0.0.1:8500",
	})
	if err != nil {
		panic(err)
	}
	registry := consul.New(consulClient)
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint("discovery:///"+"Kar"),
		grpc.WithDiscovery(registry),
		grpc.WithMiddleware(
			recovery.Recovery(),
			tracing.Client(),
			metadata.Client(),
			circuitbreaker.Client()),
	)
	if err != nil {
		panic(err)
	}

	client := v1.NewGreeterClient(conn)
	hello, err := client.SayHello(context.Background(), &v1.HelloRequest{Name: "kratos"})

	if err != nil {
		panic(err)
	}
	fmt.Println(hello.Message)
}
