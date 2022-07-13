// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"megpoid.dev/go/contact-form/model"
)

func (a *App) HealthCheck(ctx context.Context) *model.HealthCheckResult {
	err := a.Srv().Store.HealthCheck().HealthCheck(ctx)
	return &model.HealthCheckResult{
		Ping: err,
	}
}
