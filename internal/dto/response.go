package dto

type MetricResponseDto struct {
	Name  string   `json:"name"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}
