package base

import "jvmgo/ch01/rtda"
import "jvmgo/ch01/rtda/heap"
/**
 * 定位到调用方法后，jvm要创建栈帧并推入栈顶，并传递参数
 */ 
func InvokeMethod(invokerFrame *rtda.Frame,method *heap.Method){
	thread := invokerFrame.Thread()
	newFrame := thread.NewFrame(method)
	thread.PushFrame(newFrame)
	// 传递参数
	argsSlotCount := int(method.ArgSlotCount())
	if argsSlotCount>0{
		for i:=argsSlotCount-1;i>=0;i--{
			slot := invokerFrame.OperandStack().PopSlot()
			newFrame.LocalVars().SetSlot(uint(i),slot)
		}
	}

}