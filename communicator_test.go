package communicator

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/turbot/go-exec-communicator/shared"
)

func TestCommunicator_new(t *testing.T) {
	ci := shared.ConnectionInfo{
		Type: "telnet",
		Host: "127.0.0.1",
	}

	if _, err := New(ci); err == nil {
		t.Fatalf("expected error with telnet")
	}

	ci.Type = "ssh"
	if _, err := New(ci); err != nil {
		t.Fatalf("err: %v", err)
	}

	ci.Type = "winrm"
	if _, err := New(ci); err != nil {
		t.Fatalf("err: %v", err)
	}
}
func TestRetryFunc(t *testing.T) {
	origMax := maxBackoffDelay
	maxBackoffDelay = time.Second
	origStart := initialBackoffDelay
	initialBackoffDelay = 10 * time.Millisecond

	defer func() {
		maxBackoffDelay = origMax
		initialBackoffDelay = origStart
	}()

	// succeed on the third try
	errs := []error{io.EOF, &net.OpError{Err: errors.New("ERROR")}, nil}
	count := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := Retry(ctx, func() error {
		if count >= len(errs) {
			return errors.New("failed to stop after nil error")
		}

		err := errs[count]
		count++

		return err
	})

	if count != 3 {
		t.Fatal("retry func should have been called 3 times")
	}

	if err != nil {
		t.Fatal(err)
	}
}

func TestRetryFuncBackoff(t *testing.T) {
	origMax := maxBackoffDelay
	maxBackoffDelay = time.Second
	origStart := initialBackoffDelay
	initialBackoffDelay = 100 * time.Millisecond

	retryTestWg = &sync.WaitGroup{}
	retryTestWg.Add(1)

	defer func() {
		maxBackoffDelay = origMax
		initialBackoffDelay = origStart
		retryTestWg = nil
	}()

	count := 0

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	Retry(ctx, func() error {
		count++
		return io.EOF
	})
	cancel()
	retryTestWg.Wait()

	if count > 4 {
		t.Fatalf("retry func failed to backoff. called %d times", count)
	}
}
