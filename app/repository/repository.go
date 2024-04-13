// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"context"

	"go.megpoid.dev/go-skel/pkg/repo"
	"megpoid.dev/go/contact-form/app/model"
)

type HealthcheckRepo interface {
	Execute(ctx context.Context) error
}

type ContactRepo interface {
	repo.GenericStore[*model.Contact]
}
