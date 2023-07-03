package throttler

import (
    "context"
    "sync"
    "time"

    "golang.org/x/sync/errgroup"
    "golang.org/x/sync/semaphore"

    "github.com/khuenqdev/go-toolbox/types"
)

// ApplyOperationToList run concurrent operations on a list of inputs type TInput and get a list of outputs type TOutput with concurrency limitation mechanism
func ApplyOperationToList[TInput, TOutput any](items types.List[TInput], minPause, maxConcurrent int, operation func(t TInput) (TOutput, error)) (types.List[TOutput], error) {
    ctx := context.Background()
    var g errgroup.Group
    var mu sync.Mutex
    outputs := make(types.List[TOutput], len(items))
    sem := semaphore.NewWeighted(int64(maxConcurrent))

    for i, item := range items {
        i, item := i, item
        if err := sem.Acquire(ctx, 1); err != nil {
            time.Sleep(time.Duration(minPause))
            break
        }
        g.Go(func() error {
            defer sem.Release(1)

            out, err := operation(item)

            mu.Lock()
            outputs[i] = out
            mu.Unlock()

            return err
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    } else if err := ctx.Err(); err != nil {
        return nil, err
    }

    return outputs, nil
}
