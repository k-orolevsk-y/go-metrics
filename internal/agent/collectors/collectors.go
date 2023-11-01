package collectors

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/collectors/alternative"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/collectors/gopsutil"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/collectors/runtime"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/config"
	"github.com/k-orolevsk-y/go-metricts-tpl/internal/agent/metrics"
	"github.com/k-orolevsk-y/go-metricts-tpl/pkg/logger"
)

type Collector struct {
	collectors map[string]collector

	mx      sync.Mutex
	metrics []metrics.Metric

	log logger.Logger
}

type workerResult struct {
	Err error
	Res []metrics.Metric
}

type collector interface {
	Collect() error
	GetResults() []metrics.Metric
}

func NewCollector(log logger.Logger) *Collector {
	alternativeCollector := alternative.NewAlternativeCollector()
	gopsutilCollector := gopsutil.NewGopsutilCollector()
	runtimeCollector := runtime.NewRuntimeCollector()

	return &Collector{
		collectors: map[string]collector{
			"alternative": alternativeCollector,
			"gopsutil":    gopsutilCollector,
			"runtime":     runtimeCollector,
		},

		log: log,
	}
}

func (c *Collector) GetMetrics() []metrics.Metric {
	c.mx.Lock()
	defer c.mx.Unlock()

	return c.metrics
}

func (c *Collector) Run(ctx context.Context) {
	collectors := []string{"alternative", "gopsutil", "runtime"}

	jobs := make(chan string, len(collectors))
	results := make(chan workerResult, len(collectors))

	defer close(jobs)
	defer close(results)

	for w := 1; w <= config.Config.RateLimit; w++ {
		go c.worker(jobs, results)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.runWorkers(collectors, jobs, results)
			time.Sleep(time.Second * time.Duration(config.Config.PollInterval))
		}
	}
}

func (c *Collector) runWorkers(collectors []string, jobs chan<- string, results <-chan workerResult) {
	for _, col := range collectors {
		jobs <- col
	}

	c.mx.Lock()
	c.metrics = make([]metrics.Metric, 0)

	for r := 1; r <= len(collectors); r++ {
		result := <-results

		if result.Err != nil {
			c.log.Errorf("Error collect metrics: %s", result.Err)
			continue
		}

		c.metrics = append(c.metrics, result.Res...)
	}

	c.mx.Unlock()
}

func (c *Collector) worker(jobs <-chan string, results chan<- workerResult) {
	for job := range jobs {
		col, ok := c.collectors[job]
		if !ok {
			results <- workerResult{
				Err: fmt.Errorf("invalid collector"),
			}
			continue
		}

		if err := col.Collect(); err != nil {
			results <- workerResult{
				Err: err,
			}
			continue
		}

		results <- workerResult{
			Res: col.GetResults(),
		}
	}
}
