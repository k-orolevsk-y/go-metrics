package models

type MetricType string

const (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

type (
	MetricsUpdate struct {
		ID    string   `json:"id" binding:"required"`
		MType string   `json:"type" binding:"required,oneof=counter gauge"`
		Delta *int64   `json:"delta,omitempty" binding:"required_if=MType counter"`
		Value *float64 `json:"value,omitempty" binding:"required_if=MType gauge"`
	}

	MetricsValue struct {
		ID    string   `json:"id" db:"name" binding:"required"`
		MType string   `json:"type" db:"mtype" binding:"required,oneof=counter gauge"`
		Delta *int64   `json:"delta,omitempty" db:"delta"`
		Value *float64 `json:"value,omitempty" db:"value"`
	}
)
