package postgres

import (
	"context"
	"github.com/jmoiron/sqlx"
	"metric-server/internal/models"
)

type Storage struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Storage {
	return &Storage{db: db}
}

func (s Storage) Save(ctx context.Context, metric models.Metric) error {
	upsertSql := `
		INSERT INTO metrics (name, type, delta, value)
		VALUES (:name, :type, :delta, :value)
		ON CONFLICT (name, type)
		DO UPDATE SET name=EXCLUDED.name, type=EXCLUDED.type, delta=EXCLUDED.delta, value=EXCLUDED.value`
	_, err := s.db.NamedExec(upsertSql, metric)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) SaveMany(ctx context.Context, metrics []models.Metric) error {
	upsertSql := `
		INSERT INTO metrics (name, type, delta, value)
		VALUES (:name, :type, :delta, :value)
		ON CONFLICT (name, type)
		DO UPDATE SET name=EXCLUDED.name, type=EXCLUDED.type, delta=EXCLUDED.delta, value=EXCLUDED.value`
	_, err := s.db.NamedExec(upsertSql, metrics)
	if err != nil {
		return err
	}

	return nil
}

func (s Storage) FindByNameAndType(ctx context.Context, name string, metricType string) (models.Metric, error) {
	metric := models.Metric{}

	sqlSelect := "SELECT m.id, m.name, m.type, m.delta, m.value FROM metrics m WHERE m.name = $1 and m.type = $2"
	if err := s.db.GetContext(ctx, &metric, sqlSelect, name, metricType); err != nil {
		return models.Metric{}, err

	}

	return metric, nil
}

func (s Storage) FindAllByType(ctx context.Context, metricType string) ([]models.Metric, error) {
	metrics := []models.Metric{}

	sql := "SELECT m.id, m.name, m.type, m.delta, m.value FROM metrics m WHERE m.type = $1"
	err := s.db.SelectContext(ctx, &metrics, sql, metricType)

	if err != nil {
		return nil, err
	}

	return metrics, nil
}
