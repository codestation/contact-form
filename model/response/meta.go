// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package response

type Meta struct {
	Items          int     `json:"items"`
	NextCursor     *string `json:"next_cursor,omitempty"`
	PrevCursor     *string `json:"prev_cursor,omitempty"`
	CurrentPage    *int    `json:"current_page,omitempty"`
	MaxPage        *int    `json:"max_page,omitempty"`
	TotalRecords   *int    `json:"total,omitempty"`
	RecordsPerPage *int    `json:"records_per_page,omitempty"`
}

func (m *Meta) Next() bool {
	if m.NextCursor != nil {
		return true
	}
	if m.CurrentPage != nil && m.MaxPage != nil && *m.CurrentPage < *m.MaxPage {
		return true
	}

	return false
}
