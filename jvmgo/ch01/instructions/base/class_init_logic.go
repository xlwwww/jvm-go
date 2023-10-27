package base

import "jvmgo/ch01/rtda"
import "jvmgo/ch01/rtda/heap"

func InitClass(thread *rtda.Thread,class *heap.Class){
	class.StartInit() // set flag 
	scheduleClinit(thread, class) // 新建frame并推入栈
	initSuperClass(thread, class) // 递归初始化父类
}
func scheduleClinit(thread *rtda.Thread,class *heap.Class){
	clinit := class.GetClinitMethod()
	if clinit!=nil{
		frame:=thread.NewFrame(clinit)
		thread.PushFrame(frame)
	}
}

func initSuperClass(thread *rtda.Thread,class *heap.Class){
	if !class.IsInterface() {
		superClass := class.SuperClass()
		if superClass != nil && !superClass.InitStarted() {
			InitClass(thread, superClass)
		}
	}
}