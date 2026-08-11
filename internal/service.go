package internal

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/moovit-sp-gmbh/hcloud-cli-repl/internal/db"
)

type ContextService struct {
	q *db.Queries
}

func New() *ContextService {
	return &ContextService{
		q: db.New(db.Database()),
	}
}

func (cs *ContextService) GetLatest() (db.Context, error) {
	ctx, err := cs.q.GetLatestContext(context.Background())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return db.Context{}, nil
		}
		return ctx, err
	}
	return ctx, err
}

func (cs *ContextService) Push(arg db.Context) error {
	return cs.q.CreateContext(context.Background(), db.CreateContextParams{
		Timestamp: time.Now().Unix(),
		Jwt:       arg.Jwt,
		OrgName:   arg.OrgName,
		NewName:   arg.NewName,
	})
}

func (cs *ContextService) List(page, limit int64) ([]db.Context, error) {
	return cs.q.List(context.Background(), db.ListParams{Offset: page * limit, Limit: limit})
}

func (cs *ContextService) SetLatest(id int64) error {
	return cs.q.SetTimestamp(context.Background(), db.SetTimestampParams{ID: id, Timestamp: time.Now().Unix()})
}
