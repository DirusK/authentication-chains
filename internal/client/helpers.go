/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"authentication-chains/internal/types"
)

// initClient initializes a new client.
func initClient(ctx context.Context, address string) (types.NodeClient, error) {
	conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	return types.NewNodeClient(conn), nil
}
