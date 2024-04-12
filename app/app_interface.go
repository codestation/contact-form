// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"megpoid.dev/go/contact-form/config"
	"megpoid.dev/go/contact-form/model"
)

type IApp interface {
	SaveContact(ctx context.Context, req *model.ContactRequest) (*model.Contact, error)
	HealthCheck(ctx context.Context) *model.HealthCheckResult
	Srv() *Server
	Config() *config.Config
}
