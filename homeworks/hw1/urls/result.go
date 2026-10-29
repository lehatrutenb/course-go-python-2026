package urls

import (
	"net/http"
	"sync"

	"go.uber.org/zap"
)

type Result struct {
	codes map[int]int
	mu    sync.RWMutex
	done  bool
	lg    *zap.SugaredLogger
}

func NewResult(lg *zap.SugaredLogger) *Result {
	return &Result{codes: make(map[int]int), mu: sync.RWMutex{}, done: false, lg: lg}
}

func (r *Result) GetResponsesCount(code int) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.codes[code]
}

func (r *Result) GetResult() map[int]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.codes
}

func (r *Result) Done() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.done
}

func (r *Result) setDone() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.done = true
}

func (r *Result) putResponse(resp *http.Response) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.done { // After aggreaction ctx cancellation it may get non-canceled response result.
		r.lg.Debugw("result already done, skipping response", "status", resp.StatusCode)
		return
	}
	r.codes[resp.StatusCode]++
}
