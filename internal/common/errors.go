package common

import "errors"

// ErrChannelClosed is returned when the channel that the function is uses is closed
// this is considered a normal operating error (will trigger errgroup to stop the context and shutdown).
var ErrChannelClosed = errors.New("channel closed")

// ErrContextDone is returned when the context is cancelled
// this is considered a normal operating error (will trigger errgroup to stop the context and shutdown).
var ErrContextDone = errors.New("context done")
