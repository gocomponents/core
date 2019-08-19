package grpc

import (
	"context"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gocomponents/core/proto"
	"google.golang.org/grpc"
	"net"
	"os"
	"runtime"
	"time"
)

var logCh = make(chan *proto.Log, 500)

var logServer,appName,hostIp string

type logToml struct {
	Version string
	Server  string
}

func init()  {
	var tomlPath string
	if runtime.GOOS == `windows` {
		tomlPath = "e:/glog/log.toml"
	} else {
		tomlPath = "/config/log.toml"
	}
	var logConfig *logToml
	_, err := toml.DecodeFile(tomlPath, &logConfig)
	if err != nil{
		panic(err)
	}

	logServer=logConfig.Server
	appName= os.Getenv("APP_NAME")
	hostIp= os.Getenv("HOST_IP")
	if ""==hostIp {
		address, err := net.InterfaceAddrs()

		if err != nil {
			fmt.Println(err)
			return
		}

		for _, address := range address {
			if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					hostIp= ipNet.IP.String()
					break
				}
			}
		}
	}

	go consume()
}

func consume()  {
	conn, err := grpc.Dial(logServer, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	client := proto.NewLogStashClient(conn)

	for {
		log, ok := <-logCh
		if ok {
			go func(log *proto.Log) {
				_,err:=client.Send(context.Background(), log)
				if nil!=err {
					fmt.Println(err)
				}
			}(log)
		}
	}
}


func Info(module,message,traceId string,execTime int32)  {
	if ""==appName {
		fmt.Println(module,message,traceId,execTime)
		return
	}

	log:=proto.Log{
		App:        appName,
		Module:     module,
		Level:      proto.Log_Info,
		TraceId:    traceId,
		Message:    message,
		Exception:  "",
		UserIp:     hostIp,
		ExecTime:   execTime,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	logCh<-&log
}

func Error(module,message,exception,traceId string)  {
	if ""==appName {
		fmt.Println(module,message,exception,traceId)
		return
	}
	log:=proto.Log{
		App:        appName,
		Module:     module,
		Level:      proto.Log_Error,
		TraceId:    traceId,
		Message:    message,
		Exception:  exception,
		UserIp:     hostIp,
		ExecTime:   0,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	logCh<-&log
}


func Warn(module,message,traceId string,execTime int32)  {
	if ""==appName {
		fmt.Println(module,message,traceId)
		return
	}
	log:=proto.Log{
		App:        appName,
		Module:     module,
		Level:      proto.Log_Warning,
		TraceId:    traceId,
		Message:    message,
		Exception:  "",
		UserIp:     hostIp,
		ExecTime:   execTime,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	logCh<-&log
}


func Debug(module,message,traceId string,execTime int32)  {
	if ""==appName {
		fmt.Println(module,message,traceId)
		return
	}
	log:=proto.Log{
		App:        appName,
		Module:     module,
		Level:      proto.Log_Debug,
		TraceId:    traceId,
		Message:    message,
		Exception:  "",
		UserIp:     hostIp,
		ExecTime:   execTime,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
	}
	logCh<-&log
}