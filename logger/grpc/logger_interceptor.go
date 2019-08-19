package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocomponents/core/util"
	"google.golang.org/grpc"
	"time"
)

//gRpc服务注册Logger Interceptor
func LoggerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start:=time.Now()
	traceId:=util.GetGUID()

	defer func() {
		msg,err:=json.Marshal(req)
		if nil!=err {
			fmt.Println(err)
			return
		}
		Info(info.FullMethod,string(msg),traceId,int32(time.Since(start)))
	}()

	return handler(context.WithValue(ctx,"TraceId",traceId), req)
}