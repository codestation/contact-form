// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package testdata

import "embed"

//go:embed seed
var assets embed.FS

func Assets() embed.FS {
	return assets
}