package ppln

import (
	"context"
	"errors"
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

func TestSanityTake(t *testing.T) {
	debugger.SetLabels(func() []string {
		return []string{"where", t.Name()}
	})

	chRes2 := make(chan string)

	source := NewFuncNode0x2(func() (string, string) {
		return "hello", "world"
	})
	sink2 := NewFuncNode1x0(func(v1 string) {
		chRes2 <- v1
	})

	Pipeline1(Take2(source), sink2)

	go source.Run()

	res2 := <-chRes2
	assert.Equal(t, "world", res2)
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

	assert.Equal(t, "1:hello 2:hello", res[0])
	assert.Equal(t, "1:world 2:world", res[1])
}

func TestExamplePubSub(t *testing.T) {
	type PSMessage struct {
		Data []byte
		// Other fields like pubsub id
	}

	chErr := make(chan error, 1)
	chRes := make(chan int, 1)

	source := NewFuncNode0x1(func() PSMessage {
		return PSMessage{Data: []byte(`hello`)}
	})
	root := NewFuncNode1x2(func(v1 PSMessage) (context.Context, []byte) {
		return context.Background(), v1.Data
	})
	enrich1 := NewFuncNode2x2(func(ctx context.Context, data []byte) (int, error) {
		return len(data), nil
	})
	enrich2 := NewFuncNode2x2(func(ctx context.Context, data []byte) (int, error) {
		return 0, errors.New("went wrong")
	})

	store := NewFuncNode3x1(func(ctx context.Context, v1, v2 int) error {
		chRes <- v1 + v2
		return nil
	})

	errorSink := FilterNode(func(v error) bool {
		return v != nil
	})

	errorCollector := NewFuncNode1x0(func(err error) {
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

	Pipeline1(errorSink, errorCollector)

	go source.Run()

	err := <-chErr
	assert.ErrorContains(t, err, "went wrong")
	res := <-chRes
	assert.Equal(t, 5, res)
}

// uncomment later
//func TestSources(t *testing.T) {
//	debugger.SetLabels(func() []string {
//		return []string{"where", t.Name()}
//	})
//
//	chRes := make(chan string)
//
//	type Message struct {
//		Name   string
//		Type   string
//		Events []string
//	}
//
//	source := NewFuncStreamNode0x1(func(emit1 func(Message)) {
//		emit1(Message{
//			Name:   "name1",
//			Type:   "type1",
//			Events: []string{"event1.1", "event1.2"},
//		})
//	}) // , NewSource(), WithOnDone(source, ackMsg)
//	nameProducer := NewFuncNode1x1(func(v1 Message) string {
//		return v1.Name
//	}) //, WithNonBlocking(errorSink), WithQueueLength(10000)
//	typeProducer := NewFuncNode1x1(func(v1 Message) string {
//		return v1.Type
//	})
//	typeCurator := NewFuncNode1x1(func(v1 string) string {
//		return strings.ToUpper(v1)
//	})
//	eventProducer := NewFuncStreamNode1x1(func(v1 Message, emit1 func(string)) {
//		for _, e := range v1.Events {
//			emit1(e)
//		}
//	})
//	eventCurator := NewFuncNode1x1(func(v1 string) string {
//		return strings.ToUpper(v1)
//	})
//
//	eventSink := NewFuncNode3x0(func(name, typ, event string) {
//		// does stuff
//	})
//
//	//ackMsg := NewFuncNode1x0(func(v1 Message) {
//	//	// v1.ack()
//	//})
//
//	Pipeline1(source, nameProducer)
//	Pipeline1(source, typeProducer)
//	Pipeline1(source, eventProducer)
//
//	Pipeline1(typeProducer, typeCurator)
//
//	Pipeline1(eventProducer, eventCurator)
//	Pipeline3(nameProducer, typeCurator, eventCurator, eventSink)
//
//	go source.Run()
//
//	select {
//	case res := <-chRes:
//		assert.Equal(t, "hello", res)
//	case <-time.After(time.Second):
//		require.Fail(t, "did not receive message")
//	}
//}
