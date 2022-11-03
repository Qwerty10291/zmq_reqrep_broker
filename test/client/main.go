package main

import (
	"fmt"
	"log"
	"time"

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
	client, err := app.NewClient("math", "request", ctx)
	if err != nil{
		panic(err)
	}
	client.Connect()

	socket := client.GetSocket()
	for {
		_, err := socket.Send("hello", 0)
		if err != nil{
			panic(err)
		}
		msg, err := socket.Recv(0)
		if err != nil{
			panic(err)
		}
		fmt.Println(msg)
		time.Sleep(time.Second)
	}
}