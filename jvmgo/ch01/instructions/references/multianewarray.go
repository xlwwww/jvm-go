package references

import "jvmgo/ch01/instructions/base"
import "jvmgo/ch01/rtda"
import "jvmgo/ch01/rtda/heap"

// Create new multidimensional array
type MULTI_ANEW_ARRAY struct {
	index      uint16
	dimensions uint8
}

func (self *MULTI_ANEW_ARRAY) FetchOperands(reader *base.BytecodeReader) {
	self.index = reader.ReadUint16()
	self.dimensions = reader.ReadUint8()
}
func (self *MULTI_ANEW_ARRAY) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(uint(self.index)).(*heap.ClassRef)
	arrClass := classRef.ResolvedClass()

	stack := frame.OperandStack()
	counts := popAndCheckCounts(stack, int(self.dimensions)) // 检查每一个维度是否>=0
	arr := newMultiDimensionalArray(counts, arrClass) // 创建数组
	stack.PushRef(arr)
}

func popAndCheckCounts(stack *rtda.OperandStack, dimensions int) []int32 {
	counts := make([]int32, dimensions)
	for i := dimensions - 1; i >= 0; i-- {
		counts[i] = stack.PopInt()
		if counts[i] < 0 {
			panic("java.lang.NegativeArraySizeException")
		}
	}

	return counts
}

func newMultiDimensionalArray(counts []int32, arrClass *heap.Class) *heap.Object {
	count := uint(counts[0]) // int[3][4]，这个count是3
	arr := arrClass.NewArray(count) // 创建的是[[I 类的数组
	// 由内向外剥
	if len(counts) > 1 {
		refs := arr.Refs() // [[I 数组类的引用
		for i := range refs {
			refs[i] = newMultiDimensionalArray(counts[1:], arrClass.ComponentClass()) // 向内，获取下一层元素ComponentClass,这里是[I
		}
	}

	return arr
}
