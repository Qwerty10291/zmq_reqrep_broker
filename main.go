package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	zmq "github.com/pebbe/zmq4"
	
	apps "github.com/Qwerty10291/golang_zmq_ipc/app"
	"github.com/Qwerty10291/golang_zmq_ipc/objects"
	"github.com/Qwerty10291/golang_zmq_ipc/server"
	"ipc_router/broker"
)

var AppName *string = new(string)
var ControllerHost *string = new(string)
var ControllerPort *int = new(int)
var ControllerResponseTimeout *int = new(int)
var FrontendServerName *string = new(string)
var BackendServerName *string = new(string)

func init() {
	flag.StringVar(AppName, "name", "", "router app name for ipc system")
	flag.StringVar(ControllerHost, "chost", "", "ipc controller host")
	flag.IntVar(ControllerPort, "cport", 0, "ipc controller port")
	flag.IntVar(ControllerResponseTimeout, "ctimeout", 0, "controller response timeout in milliseconds")
	flag.StringVar(FrontendServerName, "iname", "", "input server name")
	flag.StringVar(BackendServerName, "oname", "", "worker server name")
	flag.Parse()

	if len(*AppName) == 0 || len(*ControllerHost) == 0 || *ControllerPort == 0 || len(*FrontendServerName) == 0 || len(*BackendServerName) == 0{
		fmt.Println("all required parameters must be specified")
		os.Exit(1)
	}
}

func main() {
	ctx, err := zmq.NewContext()
	if err != nil{
		panic(err)
	}
	app := apps.NewApp(*AppName, apps.AppConfig{
		ControllerHost:            *ControllerHost,
		ControllerPort:            fmt.Sprintf("%d", *ControllerPort),
		ControllerResponseTimeout: time.Duration(*ControllerResponseTimeout) * time.Millisecond,
	}, ctx, log.Default())
	app.Init()
	
	routerServer, err := app.NewServer(*FrontendServerName, objects.ROUTER_SERVER, server.TCP, ctx)
	if err != nil{
		panic(err)
	}
	routerServer.Bind()
	dealerServer, err := app.NewServer(*BackendServerName, objects.DEALER_SERVER, server.TCP, ctx)
	if err != nil{
		panic(err)
	}
	dealerServer.Bind()
	
	msgBroker := broker.NewReqRepBroker(routerServer.GetSocket(), dealerServer.GetSocket())
	msgBroker.Start()
}