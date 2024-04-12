// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHealthCheckStore(t *testing.T) {
	suite.Run(t, &healthCheckSuite{})
}

type healthCheckSuite struct {
	suite.Suite
	conn *connection
}

func (s *healthCheckSuite) SetupTest() {
	s.conn = NewTestConnection(s.T(), false)
}

func (s *healthCheckSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *healthCheckSuite) TestPing() {
	store := newSqlHealthCheckStore(s.conn.store)
	err := store.HealthCheck(context.Background())
	s.NoError(err)
}
