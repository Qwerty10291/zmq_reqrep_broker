package broker

import (zmq "github.com/pebbe/zmq4")

type ReqRepBroker struct {
	poller *zmq.Poller
	frontend *zmq.Socket
	backend *zmq.Socket
}

func NewReqRepBroker(router *zmq.Socket, dealer *zmq.Socket) *ReqRepBroker {
	poller := zmq.NewPoller()
	poller.Add(router, zmq.POLLIN)
	poller.Add(dealer, zmq.POLLIN)
	return &ReqRepBroker{
		poller:   poller,
		frontend: router,
		backend:  dealer,
	}
}

func (b *ReqRepBroker) Start(){
	for {
		sockets, _ := b.poller.Poll(-1)
		for _, socket := range sockets {
			switch s := socket.Socket; s {
			case b.frontend:
				for {
					msg, _ := s.Recv(0)
					if more, _ := s.GetRcvmore(); more {
						b.backend.Send(msg, zmq.SNDMORE)
					} else {
						b.backend.Send(msg, 0)
						break
					}
				}
			case b.backend:
				for {
					msg, _ := s.Recv(0)
					if more, _ := s.GetRcvmore(); more {
						b.frontend.Send(msg, zmq.SNDMORE)
					} else {
						b.frontend.Send(msg, 0)
						break
					}
				}
			}
		}
	}
}