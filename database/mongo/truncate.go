package mongo

import (
	"context"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Truncator struct {
	db *mongo.Database
}

func NewTruncator(dsn string) (*Truncator, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return nil, err
	}

	dbName := strings.TrimPrefix(u.Path, "/")
	conn := client.Database(dbName, nil)

	return &Truncator{db: conn}, err
}

func (t *Truncator) Truncate(ctx context.Context) error {
	return t.db.Drop(ctx)
}
