// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package web

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (web *Web) InitStatic() {
	assetHandler := http.FileServer(getFileSystem())
	web.root.GET("/", echo.WrapHandler(assetHandler))
	web.root.GET("/static/*", echo.WrapHandler(http.StripPrefix("/static/", assetHandler)))
}

func getFileSystem() http.FileSystem {
	fsTree, err := fs.Sub(Assets(), "static")
	if err != nil {
		panic(err)
	}

	return http.FS(fsTree)
}
