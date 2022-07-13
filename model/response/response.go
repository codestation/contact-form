// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"megpoid.dev/go/contact-form/model"
	"megpoid.dev/go/contact-form/store/paginator/cursor"
)

type ListResponse[T model.Modelable] struct {
	Data []T        `json:"data"`
	Meta Pagination `json:"meta"`
}

func NewListResponse[T model.Modelable](results []T, c *cursor.Cursor) *ListResponse[T] {
	return &ListResponse[T]{
		Data: results,
		Meta: Pagination{
			Items:      len(results),
			NextCursor: c.After,
			PrevCursor: c.Before,
		},
	}
}
