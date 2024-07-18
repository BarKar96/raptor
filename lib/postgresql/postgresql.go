package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/barkar96/raptor/lib/actor"
	"github.com/barkar96/raptor/lib/logging"
)

var _ actor.Actor = (*PostgreSQL)(nil)

type PostgreSQL struct {
	db *sql.DB
}

func New(user, password, host, port, database string) (*PostgreSQL, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	return &PostgreSQL{db: db}, nil
}

func (p *PostgreSQL) Name() string {
	return "PostgreSQL"
}

func (p *PostgreSQL) Start(ctx context.Context) error {
	t := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-t.C:
			p.checkConnection(ctx)
		}
	}
}

func (p *PostgreSQL) checkConnection(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := p.db.PingContext(ctx); err != nil {
		logging.WithError(ctx, err, "failed to ping PostgreSQL")
	}
}

func (p *PostgreSQL) Stop(_ context.Context) {
	p.db.Close()
}
