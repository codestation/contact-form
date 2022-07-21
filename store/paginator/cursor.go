// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package paginator

import "megpoid.dev/go/contact-form/store/paginator/cursor"

type MetaType string

const (
	MetaCursor MetaType = "cursor"
	MetaOffset MetaType = "offset"
	MetaNone   MetaType = "none"
)

type Page struct {
	Items        int
	Total        int
	Page         int
	ItemsPerPage int
}

type Cursor struct {
	cursor *cursor.Cursor
	offset *Page
}

func (c *Cursor) Type() MetaType {
	switch {
	case c.cursor != nil:
		return MetaCursor
	case c.offset != nil:
		return MetaOffset
	default:
		return MetaNone
	}
}

func (c *Cursor) SetCursor(cur *cursor.Cursor) {
	c.cursor = cur
}

func (c *Cursor) SetOffset(off *Page) {
	c.offset = off
}

func (c *Cursor) Cursor() *cursor.Cursor {
	return c.cursor
}

func (c *Cursor) Offset() *Page {
	return c.offset
}
