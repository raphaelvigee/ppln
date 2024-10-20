package ppln

import (
	"context"
	"fmt"
	"github.com/dlsniper/debugger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSanity(t *testing.T) {
	debugger.SetLabels(func() []string {
		return []string{"where", t.Name()}
	})

	chRes := make(chan string)

	source := NewFuncNode0x1(func() string {
		return "hello"
	})
	sink := NewFuncNode1x0(func(v1 string) {
		chRes <- v1
	})

	Pipeline1(source, sink)

	go source.Run()

	select {
	case res := <-chRes:
		assert.Equal(t, "hello", res)
	case <-time.After(time.Second):
		require.Fail(t, "did not receive message")
	}
}

func TestFanOut(t *testing.T) {
	debugger.SetLabels(func() []string {
		return []string{"where", t.Name()}
	})

	var received atomic.Int64
	var wg sync.WaitGroup
	wg.Add(2)

	source := NewFuncNode0x1(func() string {
		return "hello"
	})
	sink1 := NewFuncNode1x0(func(v1 string) {
		defer wg.Done()

		received.Add(int64(len(v1)))
	})
	sink2 := NewFuncNode1x0(func(v1 string) {
		defer wg.Done()

		received.Add(int64(len(v1)))
	})

	Pipeline1(source, sink1)
	Pipeline1(source, sink2)

	go source.Run()

	wg.Wait()

	assert.Equal(t, int64(10), received.Load())
}

func TestSanityMultipleInOut(t *testing.T) {
	debugger.SetLabels(func() []string {
		return []string{"where", t.Name()}
	})

	chRes := make(chan string)

	source := NewFuncStreamNode3x2(func(v string, delay1, delay2 time.Duration, emit1 func(string), emit2 func(string)) {
		var g errgroup.Group
		g.Go(func() error {
			time.Sleep(delay1)

			emit1("1:" + v)

			return nil
		})

		g.Go(func() error {
			time.Sleep(delay2)

			emit2("2:" + v)

			return nil
		})

		_ = g.Wait()
	})
	sink := NewFuncNode2x0(func(v1, v2 string) {
		chRes <- v1 + " " + v2
	})

	Pipeline2(Take1(source), Take2(source), sink)

	go source.Run("hello", 0, time.Second)
	go source.Run("world", time.Second, 0)

	res1 := <-chRes
	res2 := <-chRes

	res := []string{res1, res2}
	slices.Sort(res)

	assert.Equal(t, "1:hello 1:hello", res[0])
	assert.Equal(t, "2:world 2:world", res[1])
}

func Test5(t *testing.T) {
	type PSMessage struct {
		Data []byte
		ID   string
	}

	chErr := make(chan error)
	chRes := make(chan int)

	source := NewFuncStreamNode0x1(func(emit1 func(PSMessage)) {
		for {
			emit1(PSMessage{})
		}
	})
	root := NewFuncNode1x2(func(v1 PSMessage) (context.Context, []byte) {
		return context.Background(), v1.Data
	})

	enrich1 := NewFuncNode2x2(func(ctx context.Context, data []byte) (int, error) {
		return len(data), nil
	})
	enrich2 := NewFuncNode2x2(func(ctx context.Context, data []byte) (int, error) {
		return len(data), nil
	})

	store := NewFuncNode3x1(func(ctx context.Context, v1, v2 int) error {
		chRes <- v1 + v2
		return nil
	})

	errorSink := NewFuncNode1x0(func(err error) {
		fmt.Println(err)
		chErr <- err
	})

	ctx := Take1(root)

	Pipeline1(source, root)
	Pipeline2(ctx, Take2(root), enrich1)
	Pipeline2(ctx, Take2(root), enrich2)

	Pipeline3(ctx, Take1(enrich1), Take1(enrich2), store)

	// Error handling
	Pipeline1(Take2(enrich1), errorSink)
	Pipeline1(Take2(enrich2), errorSink)
	Pipeline1(store, errorSink)

	go source.Run()

	select {
	case err := <-chErr:
		require.NoError(t, err)
	case res := <-chRes:
		assert.Equal(t, 1, res)
	}
}
