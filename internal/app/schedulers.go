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

		if _, err := app.scheduler.Do(func() { app.node.Sync(ctx) }); err != nil {
			app.logger.Fatal(err)
		}
	}

	if app.cfg.Schedulers.Explore.Enabled {
		app.scheduler.Every(app.cfg.Schedulers.Explore.Interval)
		if !app.cfg.Schedulers.Explore.StartImmediately {
			app.scheduler.WaitForSchedule()
		}

		if _, err := app.scheduler.Do(func() { app.node.Explore(ctx) }); err != nil {
			app.logger.Fatal(err)
		}
	}

	app.scheduler.StartAsync()

	<-ctx.Done()
	app.scheduler.Stop()
}
