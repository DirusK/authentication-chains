/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package helpers

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

const TagCLI = "CLI"

var Ctx = registerGracefulHandle()

func Truncate(str string, num int) string {
	if len(str) <= num {
		return str
	}
	return str[:num] + "..."
}

// registerGracefulHandle registers a graceful shutdown handler for the application.
func registerGracefulHandle() context.Context {
	gracefulCtx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		cancel()
	}()

	return gracefulCtx
}
