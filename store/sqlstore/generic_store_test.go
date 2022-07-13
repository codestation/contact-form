// Copyright 2022 codestation. All rights reserved.
// Use of this source code is governed by a MIT-license
// that can be found in the LICENSE file.

package sqlstore

import (
	"context"
	"errors"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
	"megpoid.dev/go/contact-form/model"
	"megpoid.dev/go/contact-form/store"
	"megpoid.dev/go/contact-form/store/clause"
	"testing"
	"time"
)

type testProfile struct {
	model.Model
	Avatar string
}

type testUser struct {
	model.Model
	Name      string
	ProfileID model.ID `goqu:"skipupdate"`
	Profile   *testProfile
}

func (t *testUser) AttachProfile(p *testProfile) {
	t.ProfileID = 0
	t.Profile = p
}

func newUser(name string, profileId model.ID) *testUser {
	u := &testUser{
		Model:     model.NewModel(),
		Name:      name,
		ProfileID: profileId,
	}
	return u
}

type userStore struct {
	*genericStore[*testUser]
	profile *profileStore
}

type profileStore struct {
	*genericStore[*testProfile]
}

func (s *userStore) Attach(ctx context.Context, results []*testUser, relation string) error {
	var err error
	switch relation {
	case "profile":
		err = attachRelation(ctx, results,
			func(m *testUser) model.ID { return m.ProfileID },
			func(m *testUser, r *testProfile) { m.AttachProfile(r) },
			s.profile.ListByIds)
	}
	return err
}

func TestStore(t *testing.T) {
	suite.Run(t, &storeSuite{})
}

type storeSuite struct {
	suite.Suite
	conn *connection
}

func (s *storeSuite) SetupTest() {
	s.conn = NewTestConnection(s.T(), true)
}

func (s *storeSuite) TearDownTest() {
	if s.conn != nil {
		s.conn.Close(s.T())
	}
}

func (s *storeSuite) TestNewStore() {
	st := NewStore[*testUser](s.conn.store)
	s.Equal("test_users", st.table)
	s.Equal([]any{"*"}, st.selectFields)
}

func (s *storeSuite) TestStoreGet() {
	st := NewStore[*testUser](s.conn.store)
	var tests = []struct {
		id  model.ID
		err error
	}{
		{1, nil},
		{0, store.ErrNotFound},
	}

	for _, test := range tests {
		s.Run("Get", func() {
			user, err := st.Get(context.Background(), test.id)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.NotZero(user.CreatedAt)
			}
		})
	}
}

func (s *storeSuite) TestStoreList() {
	st := NewStore[*testUser](s.conn.store)
	users, err := st.List(context.Background())
	if s.NoError(err) {
		s.GreaterOrEqual(len(users.Data), 0)
	}
}

func (s *storeSuite) TestStoreSave() {
	st := NewStore[*testUser](s.conn.store)
	var tests = []struct {
		name      string
		profileId model.ID
		err       error
	}{
		{"Some user", 1, nil},
		{"Some user", 1, store.ErrDuplicated}, // do not run more tests after a constraint error
	}

	for _, test := range tests {
		s.Run("Save", func() {
			user := newUser(test.name, test.profileId)
			err := st.Save(context.Background(), user)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
			}
		})
	}
}

func (s *storeSuite) TestStoreUpdate() {
	st := NewStore[*testUser](s.conn.store)
	var tests = []struct {
		id  model.ID
		err error
	}{
		{1, nil},
		{0, store.ErrNotFound},
	}

	for _, test := range tests {
		s.Run("Update", func() {
			user := newUser("John Doe", 1)
			user.ID = test.id
			user.UpdatedAt = time.Now()
			err := st.Update(context.Background(), user)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *storeSuite) TestStoreDelete() {
	st := NewStore[*testUser](s.conn.store)
	var tests = []struct {
		id  model.ID
		err error
	}{
		{1, nil},
		{0, store.ErrNotFound},
	}

	for _, test := range tests {
		s.Run("Delete", func() {
			err := st.Delete(context.Background(), test.id)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *storeSuite) TestStoreGetExternal() {
	st := NewStore[*testUser](s.conn.store)
	var tests = []struct {
		id  uuid.UUID
		err error
	}{
		{uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")), nil},
		{uuid.Must(uuid.NewV7(uuid.MillisecondPrecision)), store.ErrNotFound},
	}

	for _, test := range tests {
		s.Run("GetExternal", func() {
			user, err := st.GetByExternalID(context.Background(), test.id)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
				s.NotZero(user.ID)
				s.NotZero(user.CreatedAt)
			}
		})
	}
}

func (s *storeSuite) TestStoreDeleteExternal() {
	st := NewStore[*testUser](s.conn.store)
	var tests = []struct {
		id  uuid.UUID
		err error
	}{
		{uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")), nil},
		{uuid.Must(uuid.NewV7(uuid.MillisecondPrecision)), store.ErrNotFound},
	}

	for _, test := range tests {
		s.Run("DeleteExternal", func() {
			err := st.DeleteByExternalId(context.Background(), test.id)
			if test.err != nil {
				s.ErrorIs(err, test.err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *storeSuite) TestBackendError() {
	db := &fakeDatabase{
		Error: errors.New("not implemented"),
	}
	conn := &SqlStore{db: db}
	st := NewStore[*testUser](conn)
	ctx := context.Background()

	_, err := st.Get(ctx, 1)
	s.ErrorIs(err, store.ErrBackend)
	_, err = st.List(ctx)
	s.ErrorIs(err, store.ErrBackend)
	err = st.Save(ctx, newUser("John Doe", 1))
	s.ErrorIs(err, store.ErrBackend)
	err = st.Update(ctx, newUser("John Doe", 1))
	s.ErrorIs(err, store.ErrBackend)
	err = st.Delete(ctx, 1)
	s.ErrorIs(err, store.ErrBackend)
	_, err = st.GetByExternalID(ctx, uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")))
	s.ErrorIs(err, store.ErrBackend)
	err = st.DeleteByExternalId(ctx, uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")))
	s.ErrorIs(err, store.ErrBackend)

	db.Result = &fakeSqlResult{Error: errors.New("not implemented")}
	err = st.Update(ctx, newUser("John Doe", 1))
	s.ErrorIs(err, store.ErrBackend)
	err = st.Delete(ctx, 1)
	s.ErrorIs(err, store.ErrBackend)
	err = st.DeleteByExternalId(ctx, uuid.Must(uuid.FromString("00000000-0000-0000-0000-000000000000")))
	s.ErrorIs(err, store.ErrBackend)
}

func (s *storeSuite) TestIncludes() {
	st := userStore{
		genericStore: NewStore[*testUser](s.conn.store,
			WithIncludes[*testUser]([]string{"profile"}),
		),
	}

	st.profile = &profileStore{
		genericStore: NewStore[*testProfile](s.conn.store),
	}

	st.AttachFunc(st.Attach)

	users, err := st.List(context.Background(), clause.WithIncludes("profile"))
	if s.NoError(err) {
		s.Len(users.Data, 1)
		user := users.Data[0]
		s.Zero(user.ProfileID)
		s.NotNil(user.Profile)
		s.Equal(model.ID(1), user.Profile.ID)
	}
}
