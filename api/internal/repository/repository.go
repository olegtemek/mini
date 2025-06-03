package repository

import (
	"context"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/olegtemek/mini/internal/config"
	"github.com/olegtemek/mini/internal/models"
)

type Repository struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) *Repository {

	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("cannot connect to db", "error", err)
		os.Exit(1)
	}

	err = dbpool.Ping(ctx)
	if err != nil {
		slog.Error("cannot ping db", "error", err)
		os.Exit(1)
	}

	// migrate)
	query := `
		create table if not exists products(
			id bigserial primary key,
			title text not null
		)
	`

	_, err = dbpool.Exec(ctx, query)
	if err != nil {
		slog.Error("cannot migrate table", "error", err, "query", query)
		os.Exit(1)
	}

	return &Repository{
		db: dbpool,
	}
}

func (r *Repository) Create(ctx context.Context, title string) (res bool, err error) {
	query := `insert into products (title) values ($1)`

	_, err = r.db.Exec(ctx, query, title)
	if err != nil {
		slog.Error("cannot create product", "error", err)
		return
	}
	res = true
	return
}

func (r *Repository) GetAll(ctx context.Context) (res []*models.Product, err error) {
	query := `select id, title from products`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		slog.Error("cannot get all products", "error", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		product := &models.Product{}
		errScan := rows.Scan(&product.Id, &product.Title)
		if errScan != nil {
			err = errScan
			slog.Error("cannot parse product", "error", errScan)
			return
		}
		res = append(res, product)
	}

	return
}

func (r *Repository) Ping() (err error) {
	err = r.db.Ping(context.Background())
	return
}

func (r *Repository) Close() (err error) {
	r.db.Close()
	return
}
