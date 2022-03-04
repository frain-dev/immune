package database

import (
	"context"

	"github.com/frain-dev/immune/database/noop"

	"github.com/frain-dev/immune"
	"github.com/frain-dev/immune/database/mongo"
)

type Truncator interface {
	Truncate(ctx context.Context) error
}

func NewTruncator(db *immune.Database) (Truncator, error) {
	switch db.Type {
	case "mongo":
		return mongo.NewTruncator(db.Dsn)
	default:
		return noop.NewTruncator()
	}
}
