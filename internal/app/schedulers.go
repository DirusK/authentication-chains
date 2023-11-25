/*
 * Copyright © 2023 Khruslov Dmytro khruslov.work@gmail.com
 */

package app

import (
	"context"
)

// serveSchedulers method as starting point for running of all schedulers.
func serveSchedulers(ctx context.Context, app *App) {
	if app.cfg.Schedulers.Sync.Enabled {
		app.scheduler.Every(app.cfg.Schedulers.Sync.Interval)
		if !app.cfg.Schedulers.Sync.StartImmediately {
			app.scheduler.WaitForSchedule()
		}

		app.scheduler.Do(app.node.Sync)
	}

	app.scheduler.StartAsync()

	<-ctx.Done()
	app.scheduler.Stop()
}
