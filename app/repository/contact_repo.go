// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package repository

import (
	"go.megpoid.dev/go-skel/pkg/repo"
	"go.megpoid.dev/go-skel/pkg/sql"
	"megpoid.dev/go/contact-form/app/model"
)

type ContactRepoImpl struct {
	*repo.GenericStoreImpl[*model.Contact]
}

func NewContact(conn sql.Executor) *ContactRepoImpl {
	s := &ContactRepoImpl{
		GenericStoreImpl: repo.NewStore[*model.Contact](conn),
	}
	return s
}
