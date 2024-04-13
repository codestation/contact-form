// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"go.megpoid.dev/go-skel/pkg/repo"
	"go.megpoid.dev/go-skel/pkg/sql"
)

type HealthcheckRepoImpl struct {
	db sql.Pinger
}

func NewHealthCheck(conn sql.Pinger) *HealthcheckRepoImpl {
	s := &HealthcheckRepoImpl{
		db: conn,
	}

	return s
}

// Execute returns an error if the database doesn't respond
func (s HealthcheckRepoImpl) Execute(ctx context.Context) error {
	if err := s.db.Ping(ctx); err != nil {
		return repo.NewRepoError(repo.ErrBackend, err)
	}
	return nil
}
