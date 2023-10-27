package constants

import "jvmgo/ch01/instructions/base"
import "jvmgo/ch01/rtda"

// Do nothing
type NOP struct{ base.NoOperandsInstruction }

func (self *NOP) Execute(frame *rtda.Frame) {
	// really do nothing
}
