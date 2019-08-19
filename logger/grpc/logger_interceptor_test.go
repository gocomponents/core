package grpc

import (
	"context"
	"fmt"
	"github.com/gocomponents/core/proto"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

func TestGrpcLoggerInterceptor_server(t *testing.T) {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(LoggerInterceptor))

	s := grpc.NewServer(opts...)

	proto.RegisterLogStashServer(s, &server{})

	go TestGrpcLoggerInterceptor_client(t)
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func TestGrpcLoggerInterceptor_client(t *testing.T) {
	time.Sleep(10*time.Second)
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithInsecure())
	if err != nil {
		fmt.Println(err)
	}

	defer conn.Close()
	client := proto.NewLogStashClient(conn)

	for i := 0; i < 2; i++ {
		log := proto.Log{
			App:        "test",
			Module:     "consume",
			Level:      proto.Log_Info,
			TraceId:    "123",
			Message:    "456",
			Exception:  "",
			UserIp:     "192.168.11.11",
			ExecTime:   12,
			CreateTime: time.Now().Add(time.Duration(i) * time.Millisecond).Format("2006-01-02 15:04:05"),
		}

		timer := time.Now()
		_,err:=client.Send(context.Background(), &log)
		fmt.Println(err)
		fmt.Println(time.Since(timer))

	}
}


type server struct{}

func (s *server) Send(ctx context.Context, request *proto.Log) (*proto.Response, error) {
	time.Sleep(10*time.Second)
	return &proto.Response{
		ErrorCode: 0,
		Message:   "",
	}, nil
}
