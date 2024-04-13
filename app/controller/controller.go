// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package controller

import (
	"megpoid.dev/go/contact-form/oapi"
)

const (
	appName    = "forms"
	apiVersion = "v1"
)

func BaseURL() string {
	return "/apis/" + appName + "/" + apiVersion
}

var _ oapi.ServerInterface = &Controller{}

type Controller struct {
	ContactController
	HealthcheckController
}
