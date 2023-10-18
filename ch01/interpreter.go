package main

import "fmt"
import "jvmgo/ch01/classfile"
import "jvmgo/ch01/rtda"
import "jvmgo/ch01/instructions"
import "jvmgo/ch01/instructions/base"

func interpret(methodInfo *classfile.MemberInfo) {
	// 获取方法的属性
	codeAttr := methodInfo.CodeAttribute()
	maxLocals := codeAttr.MaxLocals()
	maxStack := codeAttr.MaxStack()
	bytecode := codeAttr.Code()

	thread := rtda.NewThread()
	frame := thread.NewFrame(maxLocals, maxStack)
	thread.PushFrame(frame)
	defer catchErr(frame)
	// 执行
	loop(thread,bytecode)
}

func catchErr(frame *rtda.Frame) {
	if r := recover(); r != nil {
		fmt.Printf("LocalVars:%v\n", frame.LocalVars())
		fmt.Printf("OperandStack:%v\n", frame.OperandStack())
		panic(r)
	}
}
func loop(thread *rtda.Thread, bytecode []byte){
	frame := thread.PopFrame()
	reader := &base.BytecodeReader{}
	for{
		pc := frame.NextPC() // for branch指令
		thread.SetPC(pc)
		reader.Reset(bytecode,pc)
		// 指令码
		opcode:=reader.ReadUint8()
		// 指令
		inst := instructions.NewInstruction(opcode)
		inst.FetchOperands(reader)
		frame.SetNextPC(reader.PC())// 下一条指令
		fmt.Printf("pc:%2d inst:%T %v\n", pc, inst, inst)
		inst.Execute(frame) // 如果是branch指令，会在frame中重新SetNextPC
	}
}