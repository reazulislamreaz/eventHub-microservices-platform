package client

import (
	eventv1 "github.com/eventhub/proto/gen/event/v1"
	ticketv1 "github.com/eventhub/proto/gen/ticket/v1"
	userv1 "github.com/eventhub/proto/gen/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClients struct {
	User   userv1.UserServiceClient
	Event  eventv1.EventServiceClient
	Ticket ticketv1.TicketServiceClient
	conns  []*grpc.ClientConn
}

func NewGRPCClients(userAddr, eventAddr, ticketAddr string) (*GRPCClients, error) {
	userConn, err := dial(userAddr)
	if err != nil {
		return nil, err
	}
	eventConn, err := dial(eventAddr)
	if err != nil {
		userConn.Close()
		return nil, err
	}
	ticketConn, err := dial(ticketAddr)
	if err != nil {
		userConn.Close()
		eventConn.Close()
		return nil, err
	}

	return &GRPCClients{
		User:   userv1.NewUserServiceClient(userConn),
		Event:  eventv1.NewEventServiceClient(eventConn),
		Ticket: ticketv1.NewTicketServiceClient(ticketConn),
		conns:  []*grpc.ClientConn{userConn, eventConn, ticketConn},
	}, nil
}

func dial(addr string) (*grpc.ClientConn, error) {
	return grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func (c *GRPCClients) Close() {
	for _, conn := range c.conns {
		_ = conn.Close()
	}
}
