package ppln

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	ch := make(chan struct{})
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

	NewDAG(
		Pipeline1(Take1(counter), printer),
		Pipeline2(Take1(counter), Take1(counter), complicatedPrinter),
	)

	<-ch
}

func TestStream1(t *testing.T) {
	//printer := NewFuncStreamNode0x1(func(emit func(v int)) {
	//	for {
	//		emit(1)
	//		time.Sleep(time.Second)
	//	}
	//})

	counter := (StreamNode1x2[string, int, error])(nil)
	printer := (StreamNode1x0[int])(nil)
	complicatedPrinter := (Node2x0[int, int])(nil)

	NewDAG(
		Pipeline1(Take1(counter), printer),
		Pipeline2(Take1(counter), Take1(counter), complicatedPrinter),
	)
}

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
