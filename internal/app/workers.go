/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package app

import (
	"context"
	"sync"
)

type worker func(ctx context.Context, a *App)

func (a *App) initWorkers() []worker {
	workers := []worker{
		serveGRPCServer,
		serveSchedulers,
	}

	return workers
}

// runWorkers run workers.
func (a *App) runWorkers() {
	workers := a.initWorkers()

	wg := new(sync.WaitGroup)
	wg.Add(len(workers))

	for _, work := range workers {
		go func(ctx context.Context, work func(context.Context, *App), app *App) {
			work(ctx, app)
			wg.Done()
		}(a.ctx, work, a)
	}

	wg.Wait()
}
