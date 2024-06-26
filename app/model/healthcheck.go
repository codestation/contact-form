// Copyright 2024 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package model

type HealthCheckResult struct {
	Ping error `json:"ping"`
}

func (h HealthCheckResult) AllOk() bool {
	return h.Ping == nil
}
