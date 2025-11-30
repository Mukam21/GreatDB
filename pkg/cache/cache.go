package cache

import (
	"time"

	"github.com/ccoveille/go-safecast"
	"github.com/dustin/go-humanize"
	"github.com/rs/zerolog"
)

type KeyString interface {
	comparable
	KeyString() string
}

type StringKey string

func (sk StringKey) KeyString() string {
	return string(sk)
}

type Config struct {
	NumCounters int64

	MaxCost int64

	DefaultTTL time.Duration
}

func (c *Config) MarshalZerologObject(e *zerolog.Event) {
	maxCost, _ := safecast.Convert[uint64](c.MaxCost)
	e.
		Str("maxCost", humanize.IBytes(maxCost)).
		Int64("numCounters", c.NumCounters).
		Dur("defaultTTL", c.DefaultTTL)
}

type Cache[K KeyString, V any] interface {
	Get(key K) (V, bool)

	Set(key K, entry V, cost int64) bool

	Wait()

	Close()

	GetMetrics() Metrics

	zerolog.LogObjectMarshaler
}

type Metrics interface {
	Hits() uint64

	Misses() uint64

	CostAdded() uint64

	CostEvicted() uint64
}

func NoopCache[K KeyString, V any]() Cache[K, V] { return &noopCache[K, V]{} }

type noopCache[K KeyString, V any] struct{}

var _ Cache[StringKey, any] = (*noopCache[StringKey, any])(nil)

func (no *noopCache[K, V]) Get(_ K) (V, bool)          { return *new(V), false }
func (no *noopCache[K, V]) Set(_ K, _ V, _ int64) bool { return false }
func (no *noopCache[K, V]) Wait()                      {}
func (no *noopCache[K, V]) Close()                     {}
func (no *noopCache[K, V]) GetMetrics() Metrics        { return &noopMetrics{} }
func (no *noopCache[K, V]) MarshalZerologObject(e *zerolog.Event) {
	e.Bool("enabled", false)
}

type noopMetrics struct{}

var _ Metrics = (*noopMetrics)(nil)

func (no *noopMetrics) Hits() uint64        { return 0 }
func (no *noopMetrics) Misses() uint64      { return 0 }
func (no *noopMetrics) CostAdded() uint64   { return 0 }
func (no *noopMetrics) CostEvicted() uint64 { return 0 }
