// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"megpoid.dev/go/contact-form/store"
)

type fakeStore struct {
	HealthCheckError error
}

func (s fakeStore) Contact() store.ContactStore {
	panic("not implemented")
}

func (s fakeStore) HealthCheck() store.HealthCheckStore {
	return &fakeHealthCheckStore{
		Error: s.HealthCheckError,
	}
}

func (s fakeStore) WithTransaction(_ context.Context, _ func(s store.Store) error) error {
	panic("implement me")
}

type fakeHealthCheckStore struct {
	Error error
}

func (f fakeHealthCheckStore) HealthCheck(_ context.Context) error {
	return f.Error
}

func TestApp_HealthCheck(t *testing.T) {
	ss := &fakeStore{}
	srv := &Server{Store: ss}
	app := New(srv)
	result := app.HealthCheck(context.Background())
	assert.NotNil(t, result)
	assert.True(t, result.AllOk())
}

func TestApp_HealthCheckError(t *testing.T) {
	ss := &fakeStore{HealthCheckError: errors.New("an error")}
	srv := &Server{Store: ss}
	app := New(srv)
	result := app.HealthCheck(context.Background())
	assert.NotNil(t, result)
	assert.False(t, result.AllOk())
}
