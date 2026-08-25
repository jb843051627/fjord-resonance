package sqlite

import "context"

func batchOperationContext(ctx context.Context) context.Context { return ctx }

func (s *Store) PrepareBatchContext(ctx context.Context) context.Context {
	return batchOperationContext(ctx)
}
