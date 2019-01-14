package profiler

//Implementation of StockInterface
var _ ProfilerInterface = (*rootProfiler)(nil)
var _ ProfilerInterface = (*profiler)(nil)

type ProfilerInterface interface {
	Start(key string) *profiler
	End()
}
