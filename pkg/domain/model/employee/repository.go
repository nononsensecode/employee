package employee

import "context"

type Repository interface {
	Save(ctx context.Context, emp Employee) (saved Employee, err error)
	FindByID(ctx context.Context, id int64) (emp Employee, err error)
}
