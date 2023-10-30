package metrics

type MetricType string

var (
	GaugeType   MetricType = "gauge"
	CounterType MetricType = "counter"
)

type Metric struct {
	ID    string     `json:"id"`
	MType MetricType `json:"type"`
	Delta *int64     `json:"delta,omitempty"`
	Value *float64   `json:"value,omitempty"`
}

func NewMetric(id string, mType MetricType, delta int64, value float64) Metric {
	return Metric{
		ID:    id,
		MType: mType,
		Delta: &delta,
		Value: &value,
	}
}

func (m *Metric) IsNil() bool {
	return m.ID == ""
}
