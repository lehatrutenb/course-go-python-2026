package urls

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"
)

const MAX_AGGREGATE_WORKERS = 10

type ResponseCodeAggregator struct {
	client HttpClient
	lg     *zap.SugaredLogger
	eg     *errgroup.Group // Can use waitgroup, but errgroup has SetLimit.
}

func NewURLAggregator(client HttpClient) *ResponseCodeAggregator {
	zapCfg := zap.NewProductionConfig()
	zapCfg.Level.SetLevel(zapcore.WarnLevel)
	lg, _ := zapCfg.Build()

	sugar := lg.Sugar()

	eg, _ := errgroup.WithContext(context.Background())
	eg.SetLimit(MAX_AGGREGATE_WORKERS)

	return &ResponseCodeAggregator{client: client, lg: sugar, eg: eg}
}

func (u *ResponseCodeAggregator) Aggregate(ctx context.Context, urls <-chan string) (AggregationResult, error) {
	res := NewResult(u.lg)

	go func() {
		defer func() {
			res.setDone()
		}()

		aggregatenWg := sync.WaitGroup{}
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

	urlsLoop:
		for {
			select {
			case <-ctx.Done():
				return
			case url, ok := <-urls:
				if !ok {
					break urlsLoop
				}
				u.lg.Infow("aggregating url", "url", url)
				aggregatenWg.Add(1)
				u.eg.Go(func() error {
					defer aggregatenWg.Done()
					u.aggregateURL(url, res)
					return nil // Do not fail whole aggregation on one error.
				})
			}
		}

		go func() {
			aggregatenWg.Wait()
			cancel()
		}()

		<-ctx.Done()
	}()

	return res, nil
}

func (u *ResponseCodeAggregator) aggregateURL(url string, res *Result) error {
	resp, err := u.client.Get(url) // I strongly believe it is not possible to implement force cancelation without context pass (trade goroutine leak or use not Get) (but possibly lector put it to err != nil - in that case we will work correctly)
	u.lg.Infow("got response", "url", url, "status", resp.StatusCode)
	if err != nil {
		u.lg.Errorw("error getting url", "url", url, "error", err)
		return fmt.Errorf("error getting url: %v", err)
	}
	defer resp.Body.Close()

	// Ctx check will not help there.
	u.lg.Debugw("putting response", "url", url, "status", resp.StatusCode)
	res.putResponse(resp)

	return nil
}

func (u *ResponseCodeAggregator) Wait() error { // Not asked, but what if call many times and dont want to Done on every result.
	defer u.lg.Sync()
	return u.eg.Wait()
}
