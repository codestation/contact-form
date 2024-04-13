// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package db

import "embed"

//go:embed migrations
var assets embed.FS

//go:embed seed
var seeds embed.FS

func Assets() embed.FS {
	return assets
}

func Seeds() embed.FS {
	return seeds
}
