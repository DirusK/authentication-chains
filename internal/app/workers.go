/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package app

import (
	"context"
	"sync"
)

// runWorkers run workers.
func (a *App) runWorkers() {
	workers := a.initWorkers()

	wg := new(sync.WaitGroup)
	wg.Add(len(workers))

	for _, work := range workers {
		go func(ctx context.Context, work func(context.Context, *App), t *App) {
			work(ctx, t)
			wg.Done()
		}(a.ctx, work, a)
	}

	wg.Wait()
}
