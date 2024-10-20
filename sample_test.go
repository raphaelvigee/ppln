package ppln

import (
	"context"
	"fmt"
	"github.com/dlsniper/debugger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	counter := NewFuncNode1x0(func(v1 string) {
		chRes <- v1
	})

	Pipeline1(source, counter)

	go source.Run()

	select {
	case res := <-chRes:
		assert.Equal(t, "hello", res)
	case <-time.After(time.Second):
		require.Fail(t, "did not receive message")
	}
}

func Test1(t *testing.T) {
	ch := make(chan struct{})
	source := NewFuncNode0x1(func() string {
		return "3456789"
	})
	counter := NewFuncNode1x2(func(v1 string) (int, error) {
		return len(v1), nil
	})
	printer := NewFuncNode1x0(func(v1 int) {
		fmt.Println(v1)
	})
	complicatedPrinter := NewFuncNode2x0(func(v1, v2 int) {
		fmt.Println(v1, v2)
		close(ch)
	})
	errorSink := NewFuncNode1x0(func(err error) {
		fmt.Println(err)
	})

	Pipeline1(source, counter)
	Pipeline1(Take1(counter), printer)
	Pipeline2(Take1(counter), Take1(counter), complicatedPrinter)

	Pipeline1(Take2(counter), errorSink)

	<-ch
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

//func Test3(t *testing.T) {
//	ch := make(chan struct{})
//	source := NewFuncStreamNode0x1(func(emit func(v string)) {
//		for {
//			emit("dfghjk")
//		}
//	})
//	a1 := NewFuncNode1x1(func(v1 string) int {
//		return len(v1)
//	})
//	a2 := NewFuncNode1x1(func(v1 string) int {
//		return len(v1)
//	})
//
//	b2 := NewFuncNode1x2(func(v1 int) (int, error) {
//		return 1, nil
//	})
//
//	c2 := NewFuncNode1x1(func(v1 int) int {
//		return 2
//	})
//
//	join := NewFuncNode2x0(func(v1 int, v2 int) {
//
//	})
//
//	NewDAG(
//		Pipeline1(source, a1),
//		Pipeline1(source, a2),
//		Pipeline1(a2, b2),
//		Pipeline1(Take1(b2), c2),
//		Pipeline2(a1, c2, join),
//	)
//
//	<-ch
//}

//func TestStream1(t *testing.T) {
//	printer := NewFuncStreamNode1x3(func(i1 string, f func(v int), f2 func(v float32), f3 func(v bool)) {
//		f2(123)
//		f2(124)
//		f2(125)
//	})
//
//	counter := (StreamNode1x2[string, int, error])(nil)
//	printer := (StreamNode1x0[int])(nil)
//	complicatedPrinter := (Node2x0[int, int])(nil)
//
//	NewDAG(
//		Pipeline1(Take1(counter), printer),
//		Pipeline2(Take1(counter), Take1(counter), complicatedPrinter),
//	)
//}

//func Test4(t *testing.T) {
//	entrypoint := FuncNode0x1[string]{Func: func() string {
//		return ""
//	}}
//
//	parser := FuncNode1x0[string]{Func: func(string) {
//		// parsing
//
//		dag.Run(procesor, event)
//	}}
//
//	procesor := FuncStreamNode3x0[event, int]{Func: func(emit func(int)) {
//		emit(t)
//	}}
//
//	store := FuncNode1x0[int]{Func: func(v int) {
//
//	}}
//
//	dag := Dag(
//		Pipeline1(entrypoint, parser),
//		Pipeline1(parser, procesor),
//		Pipeline1(procesor, store),
//	)
//
//	dag.Run(source, "werjk")
//}
//
//func Test2(t *testing.T) {
//	source := (Node0x1[string])(nil)
//	counter1 := (Node1x2[string, int, error])(nil)
//	counter2 := (Node1x2[string, int, error])(nil)
//	summer := (Node2x1[int, int, int])(nil)
//	printer := (Node1x0[int])(nil)
//
//	errLogger := (Node1x0[error])(nil)
//
//	dag := Dag(
//		Pipeline1(source, counter1),
//		Pipeline1(source, counter2),
//
//		Pipeline1(Take2(counter1), errLogger),
//		Pipeline1(Take2(counter2), errLogger),
//
//		Pipeline2(Take1(counter1), Take1(counter2), summer),
//		Pipeline1(summer, printer),
//	)
//
//	dag.Run(source, "werjk")
//}
//
//func Test3(t *testing.T) {
//	var total atomic.Int64
//
//	counter := FuncNode1x1[string, int]{Func: func(v1 string) int {
//		return strings.Count(v1, " ")
//	}}
//
//	accumulator := FuncNode1x0[int]{Func: func(c int) {
//		total.Add(int64(c))
//	}}
//
//	ticker := FuncNode0x0[struct{}]{Func: func(emit func(int)) {
//		for {
//			emit(totalGetter)
//
//			time.Sleep(time.Second)
//		}
//	}}
//
//	totalGetter := FuncNode0x1[int]{Func: func() int {
//		return int(total.Load())
//	}}
//
//	printer := (Node1x0[int])(nil)
//
//	dag := Dag(
//		Pipeline1(counter, accumulator),
//
//		Pipeline1(totalGetter, ticker, printer),
//	)
//
//	{
//		//pubsub
//
//		dag.RunAsync(counter, "2345678io")
//	}
//
//	s := State()
//
//	dag.RunAsync(s, counter1, "2345678io")
//	dag.RunAsync(s, counter2, "2345678io")
//}
