package emprepo_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/tern/migrate"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"nononsensecode.com/employee/pkg/domain/model/employee"
	"nononsensecode.com/employee/pkg/infrastructure/sql/pgsql/emprepo"
)

var versionTable = "schema_version"

func Test_Save_Employee(t *testing.T) {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, dbUrl)
	require.NoError(t, err)
	defer func() {
		conn.Close(ctx)
	}()

	m, err := migrate.NewMigrator(ctx, conn, versionTable)
	require.NoError(t, err)
	m.LoadMigrations("./migrations/employee")

	repo := emprepo.NewEmployeeRepo(pgUsername, pgPassword, "localhost", pgPort, pgDb)

	tests := map[string]struct {
		emp            employee.Employee
		wanted         employee.Employee
		wantedErr      string
		migrateVersion int32
	}{
		"employee gets saved successfully": {
			emp:            employee.New("kaushik", 42),
			wanted:         employee.UnmarshalFromPersistence(1, "kaushik", 42),
			wantedErr:      "",
			migrateVersion: 1,
		},
	}

	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			m.MigrateTo(ctx, tc.migrateVersion)

			got, err := repo.Save(ctx, tc.emp)
			if tc.wantedErr == "" {
				require.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, tc.wantedErr)
			}
			assert.Equal(t, tc.wanted, got)
		})
	}
}
