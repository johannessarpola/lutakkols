package pipes

import (
	"context"
	"time"
)

// DoneOrSend tries to send a element into chan or cancels if context is done
func DoneOrSend[T any](ele T, dst chan T, ctx context.Context) {
	select {
	case <-ctx.Done():
		// ctx canceled before sent successfully
	case dst <- ele:
		// sent successfully within the time out context
	}
	return
}

// TimeoutOrSend tries to send a element into chan or cancels after timeout
func TimeoutOrSend[T any](ele T, dst chan T, duration time.Duration) {
	select {
	case <-time.After(duration):
		// timeout before sent
	case dst <- ele:
		// sent successfully within the time out context
	}
	return
}
