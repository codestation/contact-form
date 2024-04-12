// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

import (
	"math"

	"megpoid.dev/go/contact-form/model"
	"megpoid.dev/go/contact-form/store/paginator"
)

type ListResponse[T model.Modelable] struct {
	Data []T  `json:"data"`
	Meta Meta `json:"meta,omitempty"`
}

func NewListResponse[T model.Modelable](results []T, c *paginator.Cursor) *ListResponse[T] {
	switch c.Type() {
	case paginator.MetaCursor:
		cur := c.Cursor()
		return &ListResponse[T]{
			Data: results,
			Meta: Meta{
				Items:      len(results),
				NextCursor: cur.After,
				PrevCursor: cur.Before,
			},
		}
	case paginator.MetaOffset:
		off := c.Offset()
		return &ListResponse[T]{
			Data: results,
			Meta: Meta{
				Items:          len(results),
				TotalRecords:   model.NewType(off.Total),
				CurrentPage:    model.NewType(off.Page),
				MaxPage:        model.NewType(int(math.Ceil(float64(off.Total) / float64(off.ItemsPerPage)))),
				RecordsPerPage: model.NewType(off.ItemsPerPage),
			},
		}
	default:
		return &ListResponse[T]{
			Data: results,
		}
	}
}
