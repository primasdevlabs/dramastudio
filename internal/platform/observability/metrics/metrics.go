package metrics

type Collector interface {
	IncCounter(name string, value float64)
}
