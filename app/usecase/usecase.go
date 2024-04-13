// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package usecase

import (
	"context"

	"megpoid.dev/go/contact-form/app/model"
)

type Contact interface {
	SaveContact(ctx context.Context, req *model.ContactRequest) (*model.Contact, error)
}

type Healthcheck interface {
	Execute(ctx context.Context) error
}
