package references

import "jvmgo/ch01/instructions/base"
import "jvmgo/ch01/rtda"
import "jvmgo/ch01/rtda/heap"

// Create new object
type NEW struct{ base.Index16Instruction }

func (self *NEW) Execute(frame *rtda.Frame) {
	cp := frame.Method().Class().ConstantPool()
	classRef := cp.GetConstant(self.Index).(*heap.ClassRef)
	class := classRef.ResolvedClass()
	if !class.InitStarted(){
		frame.RevertNextPC(); // 此时frame指令已经指向下一条指令，需要revert to 当前指令
		base.InitClass(frame.Thread(),class);
		return
	}
	if class.IsInterface() || class.IsAbstract() {
		panic("java.lang.InstantiationError")
	}

	ref := class.NewObject()
	frame.OperandStack().PushRef(ref)
}
