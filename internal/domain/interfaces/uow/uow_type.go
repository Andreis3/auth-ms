package uow

import "context"

type UnitOfWorkFactory func(ctx context.Context) UnitOfWork
