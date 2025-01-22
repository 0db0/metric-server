package models

import (
	"database/sql"
	"errors"
)

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

var (
	ErrNotFound = errors.New("metric not found")
)

type Metric struct {
	ID    int
	Name  string          `db:"name" json:"name"`
	Type  string          `db:"type" json:"type"`
	Delta sql.NullInt64   `db:"delta" json:"delta,omitempty"`
	Value sql.NullFloat64 `db:"value" json:"value,omitempty"`
}

func NewMetric(
	name string,
	mType string,
	delta *int64,
	value *float64,
) Metric {
	deltaVal := sql.NullInt64{}
	if delta != nil {
		deltaVal.Int64 = *delta
		deltaVal.Valid = true
	}

	ValueVal := sql.NullFloat64{}
	if value != nil {
		ValueVal.Float64 = *value
		ValueVal.Valid = true
	}

	return Metric{
		Name:  name,
		Type:  mType,
		Delta: deltaVal,
		Value: ValueVal,
	}
}

func (m Metric) GetValue() *float64 {
	if !m.Value.Valid {
		return nil
	}

	return &m.Value.Float64
}

func (m Metric) GetDelta() *int64 {
	if !m.Delta.Valid {
		return nil
	}

	return &m.Delta.Int64
}
