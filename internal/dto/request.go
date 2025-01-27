package dto

type CollectDto struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type ValueDto struct {
	ID    string `json:"id"`
	MType string `json:"type"`
}

func NewValueDto(id string, mType string) ValueDto {
	return ValueDto{
		ID:    id,
		MType: mType,
	}
}
