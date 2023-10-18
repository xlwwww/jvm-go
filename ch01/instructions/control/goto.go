package control

import "jvmgo/ch01/instructions/base"
import "jvmgo/ch01/rtda"

// Branch always
type GOTO struct{ base.BranchInstruction }

func (self *GOTO) Execute(frame *rtda.Frame) {
	base.Branch(frame, self.Offset)
}
