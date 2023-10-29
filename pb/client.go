package pb

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InsecureConnect(address string, port int) (*grpc.ClientConn, error) {
	address = fmt.Sprintf("%s:%d", address, port)
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Connect(address string, port int) (*grpc.ClientConn, error) {
	// TODO: implement secure connection
	return InsecureConnect(address, port)
}
