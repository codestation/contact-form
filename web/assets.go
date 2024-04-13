// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package web

import "embed"

//go:embed static
var assets embed.FS

func Assets() embed.FS {
	return assets
}
