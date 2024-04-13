// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"
	"log/slog"

	"megpoid.dev/go/contact-form/app/repository"
)

// used to validate that the implementation matches the interface
var _ Healthcheck = &HealthcheckInteractor{}

type HealthcheckInteractor struct {
	healthcheckRepo repository.HealthcheckRepo
}

func (u *HealthcheckInteractor) Execute(ctx context.Context) error {
	slog.InfoContext(ctx, "Executing healthcheck")
	return u.healthcheckRepo.Execute(ctx)
}

func NewHealthcheck(repo repository.HealthcheckRepo) *HealthcheckInteractor {
	return &HealthcheckInteractor{
		healthcheckRepo: repo,
	}
}
