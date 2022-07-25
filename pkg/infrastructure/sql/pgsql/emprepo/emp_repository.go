package emprepo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/multierr"
	"nononsensecode.com/employee/pkg/domain/model/employee"
	"nononsensecode.com/employee/pkg/errors"
)

const (
	PgSqlPoolConnectErr errors.Code = iota + 100
	PgSqlTxStartErr
	PgSqlTxFinishErr
	PgSqlEmpSaveErr
	PgSqlEmpNotFoundErr
	PgSqlEmpFindErr
)

type PgSqlEmployeeRepository struct {
	username string
	password string
	host     string
	port     string
	db       string
}

func NewEmployeeRepo(username, password, host, port, db string) PgSqlEmployeeRepository {
	return PgSqlEmployeeRepository{
		username: username,
		password: password,
		host:     host,
		port:     port,
		db:       db,
	}
}

func (p PgSqlEmployeeRepository) dsn() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		p.username, p.password, p.host, p.port, p.db)
}

func (p PgSqlEmployeeRepository) Save(ctx context.Context, emp employee.Employee) (saved employee.Employee, err error) {
	pool, err := pgxpool.Connect(ctx, p.dsn())
	if err != nil {
		err = errors.NewUnknownError(PgSqlPoolConnectErr, err)
		return
	}
	defer func() {
		pool.Close()
	}()

	tx, err := pool.Begin(ctx)
	if err != nil {
		err = errors.NewUnknownError(PgSqlTxStartErr, err)
		return
	}
	defer func() {
		err = FinishTx(ctx, tx, err)

		if err == nil {
			return
		}

		if _, ok := err.(errors.ApiError); ok {
			return
		}

		err = errors.NewUnknownError(PgSqlTxFinishErr, err)
	}()

	var id int64
	err = tx.QueryRow(ctx, "INSERT INTO employee (name, age) VALUES ($1, $2) RETURNING ID", emp.Name(), emp.Age()).Scan(&id)
	if err != nil {
		err = errors.NewUnknownError(PgSqlEmpSaveErr, err)
		return
	}

	saved = employee.UnmarshalFromPersistence(id, emp.Name(), emp.Age())
	return
}

func (p PgSqlEmployeeRepository) FindByID(ctx context.Context, id int64) (found employee.Employee, err error) {
	pool, err := pgxpool.Connect(ctx, p.dsn())
	if err != nil {
		err = errors.NewUnknownError(PgSqlPoolConnectErr, err)
		return
	}
	defer func() {
		pool.Close()
	}()

	tx, err := pool.Begin(ctx)
	if err != nil {
		err = errors.NewUnknownError(PgSqlTxStartErr, err)
		return
	}
	defer func() {
		err = FinishTxReadOnly(ctx, tx, err)

		if err == nil {
			return
		}

		if _, ok := err.(errors.ApiError); ok {
			return
		}

		err = errors.NewUnknownError(PgSqlTxFinishErr, err)
	}()

	var (
		name string
		age  uint8
	)
	err = tx.QueryRow(ctx, "SELECT * FROM employee WHERE id = $1", &id).Scan(&id, &name, &age)
	if err != nil {
		switch err {
		case pgx.ErrNoRows:
			err = errors.NewNotFoundError(PgSqlEmpNotFoundErr, fmt.Errorf("there is no employee found with id %d", id))
			return
		default:
			err = errors.NewUnknownError(PgSqlEmpFindErr, err)
			return
		}
	}

	found = employee.UnmarshalFromPersistence(id, name, age)
	return
}

func FinishTx(ctx context.Context, tx pgx.Tx, err error) (newErr error) {
	if err != nil {
		rErr := tx.Rollback(ctx)
		if rErr != nil {
			newErr = multierr.Combine(rErr, err)
			return
		}
		newErr = err
		return
	}

	if cErr := tx.Commit(ctx); cErr != nil {
		newErr = cErr
		return
	}
	return
}

func FinishTxReadOnly(ctx context.Context, tx pgx.Tx, err error) (newErr error) {
	if err != nil {
		rErr := tx.Rollback(ctx)
		if rErr != nil {
			newErr = multierr.Combine(rErr, err)
			return
		}
		newErr = err
		return
	}
	return
}
