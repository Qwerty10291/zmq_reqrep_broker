package main

import (
	"log"

	apps "github.com/Qwerty10291/golang_zmq_ipc/app"
	zmq "github.com/pebbe/zmq4"
)


func main() {
	ctx, _ := zmq.NewContext()
	app := apps.NewApp("worker1", apps.AppConfig{
		ControllerHost:            "127.0.0.1",
		ControllerPort:            "5000",
		ControllerResponseTimeout: 0,
	}, ctx, log.Default())
	app.Init()
	server, err := app.NewClient("math", "workers", ctx)
	if err != nil{
		panic(err)
	}
	server.Connect()
	socket := server.GetSocket()
	for {
		data, err := socket.RecvBytes(0)
		if err != nil{
			panic(err)
		}
		socket.SendBytes(data, 0)
	}
}