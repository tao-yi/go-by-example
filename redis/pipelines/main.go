package main

import (
	"context"
	"errors"
	"fmt"
	"go-by-example/redis/client"
	"time"

	"github.com/go-redis/redis/v8"
)

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	rdb := client.Client()
	ctx := context.Background()

	pipe := rdb.Pipeline()
	// start a pipeline to execute multiple commands
	incr := pipe.Incr(ctx, "pipeline_counter")
	pipe.Expire(ctx, "pipeline_counter", time.Hour)
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, cmd := range cmds {
		fmt.Println(cmd.Err())
	}

	// The value is available only after Exec is called.
	fmt.Println(incr.Val())

	// alternatively, you can use Pipelined which calls Exec when the function exits
	_, err = rdb.Pipelined(ctx, func(p redis.Pipeliner) error {
		incr = pipe.Incr(ctx, "pipelined_counter")
		pipe.Expire(ctx, "pipelined_counter", time.Hour)
		return nil
	})
	exitOnError(err)

	// The value is available only after the pipeline is executed
	fmt.Println(incr.Val())

	_, err = rdb.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < 100; i++ {
			pipe.Get(ctx, fmt.Sprintf("key%d", i))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

// Redis transactions use optimistic locking.
const maxRetries = 1000

// Increment transactionally increments the key using GET and SET commands.
func increment(rdb *redis.Client, ctx context.Context, key string) error {
	// Transactional function.
	txf := func(tx *redis.Tx) error {
		// Get the current value or zero.
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		// Actual operation (local in optimistic lock).
		n++

		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n, 0)
			return nil
		})
		return err
	}

	// Retry if the key has been changed.
	for i := 0; i < maxRetries; i++ {
		err := rdb.Watch(ctx, txf, key)
		if err == nil {
			// Success.
			return nil
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		// Return any other error.
		return err
	}

	return errors.New("increment reached maximum number of retries")
}
