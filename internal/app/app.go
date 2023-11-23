/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package app

import (
	"context"

	"github.com/DirusK/utils/log"
	"github.com/DirusK/utils/validator"
	"github.com/alitto/pond"
	"github.com/nutsdb/nutsdb"
	"google.golang.org/grpc"

	"authentication-chains/internal/config"
	"authentication-chains/internal/node"
)

type (
	// Meta is additional information about application.
	Meta struct {
		Name  string
		Level uint32
	}

	// App is a main application structure.
	App struct {
		meta       Meta
		ctx        context.Context
		validator  validator.Validator
		cfg        *config.Config
		db         *nutsdb.DB
		grpcServer *grpc.Server
		logger     log.Logger
		workerPool *pond.WorkerPool
		node       *node.Node
	}
)

// New creates a new application instance.
func New(ctx context.Context, configPath string) *App {
	app := new(App)
	app.ctx = ctx

	app.initValidator()
	app.initConfig(configPath)

	app.meta = Meta{
		Name:  app.cfg.Node.Name,
		Level: app.cfg.Node.Level,
	}

	app.initLogger()
	app.initStorage()
	app.initWorkerPool(ctx)
	app.initNode()
	app.initGRPCServer()

	return app
}
