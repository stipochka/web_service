package repository

import "google.golang.org/grpc"

func NewGRPCConn(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
