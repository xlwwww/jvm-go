package main

import "fmt"
import "jvmgo/ch01/instructions"
import "jvmgo/ch01/instructions/base"
import "jvmgo/ch01/rtda"


func interpret(thread *rtda.Thread,logInst bool) {
	defer catchErr(thread)
	// 执行
	loop(thread,logInst)
}

func catchErr(thread *rtda.Thread) {
	if r := recover(); r != nil {
		logFrames(thread)
		panic(r)
	}
}
func loop(thread *rtda.Thread,logInst bool){
	reader := &base.BytecodeReader{}
	for{
		frame := thread.CurrentFrame() 
		pc := frame.NextPC() // 取出下一条指令
		thread.SetPC(pc)
		reader.Reset(frame.Method().Code(),pc)
		// method := frame.Method()
		// 指令码
		opcode:=reader.ReadUint8()
		// 指令
		inst := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		frame.SetNextPC(reader.PC())// 重置下一条指令
		if logInst{
			logInstruction(frame, inst)
		}
		inst.Execute(frame) // 如果是branch指令，会在frame中重新SetNextPC
		if thread.IsStackEmpty(){
			break
		}
	}
}
func logInstruction(frame *rtda.Frame, inst base.Instruction) {
	method := frame.Method()
	className := method.Class().Name()
	methodName := method.Name()
	pc := frame.Thread().PC()
	fmt.Printf("%v.%v() #%2d %T %v\n", className, methodName, pc, inst, inst)
}

func logFrames(thread *rtda.Thread) {
	for !thread.IsStackEmpty() {
		frame := thread.PopFrame()
		method := frame.Method()
		className := method.Class().Name()
		fmt.Printf(">> pc:%4d %v.%v%v \n",
			frame.NextPC(), className, method.Name(), method.Descriptor())
	}
}
