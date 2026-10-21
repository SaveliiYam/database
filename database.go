// Package database — порты (интерфейсы) для работы с БД.
//
// Ключевой принцип заменяемости: бизнес-логика зависит только от
// Transactor / HealthChecker / Closer из этого пакета, а не от pgx или gorm.
//
// Транзакционная модель единая для всех адаптеров (pgxdb, gormdb, sqldb):
//
//	err := tx.WithinTx(ctx, func(ctx context.Context) error {
//		// адаптер кладёт транзакцию в ctx;
//		// репозиторий достаёт её своим executor-методом:
//		//   pgx:  pool.Querier(ctx)
//		//   gorm: db.DB(ctx)
//		//   sql:  db.Executor(ctx)
//		return repo.Create(ctx, u)
//	})
//
// Специфика СУБД (типы pgx.Rows, *gorm.DB) остаётся внутри адаптера и
// репозитория — наружу торчит только нативный хендл через Pool()/DB()
// для тех мест, где без него не обойтись.
package database

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("database: not found")

func IsNotFound(err error) bool { return errors.Is(err, ErrNotFound) }

type Transactor interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type TxLevel int

const (
	LevelDefault TxLevel = iota
	LevelReadCommitted
	LevelRepeatableRead
	LevelSerializable
)

type TxOptions struct {
	Level    TxLevel
	ReadOnly bool
}

type TxTransactor interface {
	Transactor
	WithinTxOptions(ctx context.Context, opts TxOptions, fn func(ctx context.Context) error) error
}

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type Closer interface {
	Close() error
}

type DB interface {
	Transactor
	HealthChecker
	Closer
}
type Component struct {
	name string
	db   DB
}

func NewComponent(name string, db DB) *Component { return &Component{name: name, db: db} }

func (c *Component) Name() string { return c.name }

func (c *Component) Start(ctx context.Context) error {
	if err := c.db.Ping(ctx); err != nil {
		return err
	}
	<-ctx.Done()
	return nil
}

func (c *Component) Stop(context.Context) error { return c.db.Close() }
