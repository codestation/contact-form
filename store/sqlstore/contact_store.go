// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"megpoid.dev/go/contact-form/model"
	"megpoid.dev/go/contact-form/store"
)

type SqlContactStore struct {
	*genericStore[*model.Contact]
}

func newSqlContactStore(sqlStore *SqlStore) store.ContactStore {
	s := &SqlContactStore{
		genericStore: NewStore[*model.Contact](sqlStore),
	}
	return s
}
