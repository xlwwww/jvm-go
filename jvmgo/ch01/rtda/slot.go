package rtda
import "jvmgo/ch01/rtda/heap"

type Slot struct {
	num int32
	ref *heap.Object
}
