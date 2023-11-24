/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package app

import (
	"context"
	"time"

	utils "github.com/DirusK/utils/config"
	"github.com/DirusK/utils/log"
	"github.com/DirusK/utils/validator"
	"github.com/alitto/pond"
	"github.com/go-co-op/gocron"
	"github.com/nutsdb/nutsdb"
	"google.golang.org/grpc"

	"authentication-chains/internal/config"
	"authentication-chains/internal/node"
)

func (a *App) initValidator() {
	a.validator = validator.New()
}

func (a *App) initConfig(configPath string) {
	a.cfg = new(config.Config)

	utils.MustLoadFromFile(configPath, a.cfg)

	if err := a.validator.Struct(a.cfg, "invalid configuration"); err != nil {
		panic(err)
		return
	}
}

func (a *App) initLogger() {
	a.logger = log.New(log.WithConfig(a.cfg.Logger), log.WithAppName(a.meta.Name))
}

func (a *App) initStorage() {
	var err error

	a.db, err = nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(a.cfg.Storage.Directory),
	)
	if err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) initNode(ctx context.Context) {
	var err error

	a.node, err = node.New(a.cfg.Node, a.db, a.workerPool, a.logger)
	if err != nil {
		a.logger.Fatal(err)
	}

	if err = a.node.Init(ctx); err != nil {
		a.logger.Fatal(err)
	}
}

func (a *App) initWorkerPool(ctx context.Context) {
	a.workerPool = pond.New(a.cfg.WorkerPool.MaxWorkers, a.cfg.WorkerPool.MaxCapacity, pond.Context(ctx))
}

func (a *App) initGRPCServer() {
	a.grpcServer = grpc.NewServer()
}

func (a *App) initScheduler() {
	a.scheduler = gocron.NewScheduler(time.UTC)
}
