package main

import "fmt"
import "strings"
import "jvmgo/ch01/classpath"
import "jvmgo/ch01/instructions/base"
import "jvmgo/ch01/rtda"
import "jvmgo/ch01/rtda/heap"

type JVM struct{
	cmd *Cmd
	classLoader *heap.ClassLoader
	mainThread *rtda.Thread
}
func newJVM(cmd *Cmd) *JVM{
	cp := classpath.Parse(cmd.XjreOption, cmd.cpOption)
	classLoader := heap.NewClassLoader(cp, cmd.verboseClassFlag)
	return &JVM{
		cmd:cmd,
		classLoader:classLoader,
		mainThread:rtda.NewThread(),
	}
}
func (self *JVM) start() {
	// 初始化
	self.initVM()
	// 执行main方法
	self.execMain()
}

func (self *JVM) initVM() {
	//加载sun.mis.VM类
	vmClass := self.classLoader.LoadClass("sun/misc/VM")
	// 执行其类初始化方法
	base.InitClass(self.mainThread, vmClass)
	interpret(self.mainThread, self.cmd.verboseInstFlag)
}

func (self *JVM) execMain() {
	className := strings.Replace(self.cmd.class, ".", "/", -1)
	mainClass := self.classLoader.LoadClass(className)
	mainMethod := mainClass.GetMainMethod()
	if mainMethod == nil {
		fmt.Printf("Main method not found in class %s\n", self.cmd.class)
		return
	}

	argsArr := self.createArgsArray()
	frame := self.mainThread.NewFrame(mainMethod)
	frame.LocalVars().SetRef(0, argsArr) // 设置到局部变量表中
	self.mainThread.PushFrame(frame)
	interpret(self.mainThread, self.cmd.verboseInstFlag)
}
// 获取命令行参数并转换为java字符串数组
func (self *JVM) createArgsArray() *heap.Object {
	stringClass := self.classLoader.LoadClass("java/lang/String")
	argsLen := uint(len(self.cmd.args))
	argsArr := stringClass.ArrayClass().NewArray(argsLen)
	jArgs := argsArr.Refs()
	for i, arg := range self.cmd.args {
		jArgs[i] = heap.JString(self.classLoader, arg)
	}
	return argsArr
}