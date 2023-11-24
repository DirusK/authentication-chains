/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package app

import (
	"context"
	"net"

	"authentication-chains/internal/types"
)

// serveGRPCServer starts gRPC server.
func serveGRPCServer(ctx context.Context, app *App) {
	lis, err := net.Listen("tcp", app.cfg.Node.GRPC.Address)
	if err != nil {
		app.logger.Fatalf("failed to listen: %v", err)
	}

	types.RegisterNodeServer(app.grpcServer, app.node)

	// graceful shutdown listener.
	go func() {
		<-ctx.Done()
		app.grpcServer.GracefulStop()
	}()

	if err = app.grpcServer.Serve(lis); err != nil {
		app.logger.Fatalf("failed to serve gRPC server: %v", err)
	}
}
