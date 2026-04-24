package saramax

import (
	"context"
	"time"

	"github.com/IBM/sarama"
	"github.com/KNICEX/InkFlow/pkg/logx"
)

const (
	defaultMaxRetryInterval = 60 * time.Second
	defaultInitInterval     = 2 * time.Second
)

// ConsumeWithRetry wraps sarama ConsumerGroup.Consume with automatic reconnect on failure.
// It retries with exponential backoff until the context is cancelled.
func ConsumeWithRetry(ctx context.Context, group sarama.ConsumerGroup, topics []string, handler sarama.ConsumerGroupHandler, l logx.Logger) {
	interval := defaultInitInterval
	for {
		err := group.Consume(ctx, topics, handler)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			l.Warn("consumer exited with error, will retry",
				logx.Error(err),
				logx.Any("retryAfter", interval.String()))
		} else {
			interval = defaultInitInterval
			continue
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
		}
		interval = min(interval*2, defaultMaxRetryInterval)
	}
}
