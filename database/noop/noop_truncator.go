package noop

import "context"

type Truncator struct{}

func NewTruncator() (*Truncator, error) {
	return &Truncator{}, nil
}

func (t *Truncator) Truncate(ctx context.Context) error {
	return nil
}
