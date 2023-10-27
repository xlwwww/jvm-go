package references
import "jvmgo/ch01/instructions/base"
import "jvmgo/ch01/rtda"
import "jvmgo/ch01/rtda/heap"

type ANEW_ARRAY struct{ base.Index16Instruction }


func (self *ANEW_ARRAY) Execute(frame *rtda.Frame) {
	stack := frame.OperandStack()
	count:=stack.PopInt()
	if count < 0 {
		panic("java.lang.NegativeArraySizeException")
	}
	// 得到引用类
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()
	// 获取数组类
	arrClass := class.ArrayClass()
	// 创建引用类数组
	arr:=arrClass.NewArray(uint(count))
	//创建完成，推入操作数栈
	stack.PushRef(arr)
}

