package heap
import"jvmgo/ch01/classfile"

type ExceptionTable []*ExceptionHandler

type ExceptionHandler struct{
	startPc int
	endPc int
	handlerPc int
	catchType *ClassRef
}

func newExceptionTable(entries []*classfile.ExceptionTableEntry,cp *ConstantPool) ExceptionTable{
	table := make([]*ExceptionHandler,len(entries))
	for i,entry :=range entries{
		table[i]=&ExceptionHandler {
			startPc:int(entry.StartPc()),
			endPc:int(entry.EndPc()),
			handlerPc:int(entry.HandlerPc()),
			catchType: getCatchType(uint(entry.CatchType()),cp),
		}
	}
	return table
}
func getCatchType(index uint, cp *ConstantPool) *ClassRef {
	if index == 0 {
		return nil // catch all
	}
	return cp.GetConstant(index).(*ClassRef)
}
// 遍历查找
func (self ExceptionTable) findExceptionHandler(exClass *Class,pc int) *ExceptionHandler{
	for _,handler:=range self{
		// 命中
		if pc>=handler.startPc && pc<=handler.endPc{
			if handler.catchType == nil{
				return handler // catch-all
			}
			catchClass:=handler.catchType.ResolvedClass()
			if catchClass == exClass{
				return handler
			}
		}
	}
	return nil
}