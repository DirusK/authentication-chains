/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package app

import (
	"github.com/nutsdb/nutsdb"
	"google.golang.org/grpc"
)

type App struct {
	db         *nutsdb.DB
	grpcServer *grpc.Server
}

func New(db *nutsdb.DB, grpcServer *grpc.Server) App {
	return App{
		db:         db,
		grpcServer: grpcServer,
	}
}
