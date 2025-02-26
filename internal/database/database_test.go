package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/database/mock"
	zerrors "github.com/zitadel/zitadel/internal/errors"
)

func TestQueryJSONObject(t *testing.T) {
	type dst struct {
		A int `json:"a,omitempty"`
	}
	const (
		query = `select $1;`
		arg   = 1
	)

	tests := []struct {
		name    string
		mock    func(*testing.T) *mock.SQLMock
		want    *dst
		wantErr error
	}{
		{
			name: "tx error",
			mock: func(t *testing.T) *mock.SQLMock {
				return mock.NewSQLMock(t, mock.ExpectBegin(sql.ErrConnDone))
			},
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "DATAB-Oath6", "Errors.Internal"),
		},
		{
			name: "no rows",
			mock: func(t *testing.T) *mock.SQLMock {
				return mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(query,
						mock.WithQueryArgs(arg),
						mock.WithQueryResult([]string{"json"}, [][]driver.Value{}),
					),
				)
			},
			wantErr: sql.ErrNoRows,
		},
		{
			name: "unmarshal error",
			mock: func(t *testing.T) *mock.SQLMock {
				return mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(query,
						mock.WithQueryArgs(arg),
						mock.WithQueryResult([]string{"json"}, [][]driver.Value{{`~~~`}}),
					),
					mock.ExpectCommit(nil),
				)
			},
			wantErr: zerrors.ThrowInternal(nil, "DATAB-Vohs6", "Errors.Internal"),
		},
		{
			name: "success",
			mock: func(t *testing.T) *mock.SQLMock {
				return mock.NewSQLMock(t,
					mock.ExpectBegin(nil),
					mock.ExpectQuery(query,
						mock.WithQueryArgs(arg),
						mock.WithQueryResult([]string{"json"}, [][]driver.Value{{`{"a":1}`}}),
					),
					mock.ExpectCommit(nil),
				)
			},
			want: &dst{A: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := tt.mock(t)
			defer mock.Assert(t)
			db := &DB{
				DB: mock.DB,
			}
			got, err := QueryJSONObject[dst](context.Background(), db, query, arg)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
