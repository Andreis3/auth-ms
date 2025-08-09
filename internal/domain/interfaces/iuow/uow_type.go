package iuow

import "context"

type UnitOfWorkFactory func(ctx context.Context) IUnitOfWork
