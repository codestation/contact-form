// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"megpoid.dev/go/contact-form/model"
)

func attachRelation[T, U model.Modelable](
	ctx context.Context,
	entries []T,
	getRelationId func(m T) model.ID,
	setRelation func(m T, r U),
	listByIds func(ctx context.Context, ids []model.ID) ([]U, error),
) error {
	if len(entries) == 0 {
		return nil
	}
	// list to hold the identifiers to query
	var idList []model.ID
	// map used to keep the above list with unique items
	var uniqueMap = map[model.ID]struct{}{}

	for _, entry := range entries {
		id := getRelationId(entry)
		if _, ok := uniqueMap[id]; !ok {
			uniqueMap[id] = struct{}{}
			idList = append(idList, id)
		}
	}

	results, err := listByIds(ctx, idList)
	if err != nil {
		return err
	}

	// keep the results in a map for quicker access
	var resultMap = map[model.ID]U{}
	for _, result := range results {
		resultMap[result.GetID()] = result
	}

	for _, entry := range entries {
		if result, ok := resultMap[getRelationId(entry)]; ok {
			setRelation(entry, result)
		}
	}

	return nil
}
