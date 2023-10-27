package heap

// symbolic reference
type SymRef struct {
	cp        *ConstantPool
	className string
	class     *Class
}

func (self *SymRef) ResolvedClass() *Class {
	if self.class == nil {
		self.resolveClassRef()
	}
	return self.class
}

// jvms8 5.4.3.1
// 类D通过符号引用N引用类C的话，要解析N， 先用D的类加载器加载C
func (self *SymRef) resolveClassRef() {
	d := self.cp.class
	c := d.loader.LoadClass(self.className)
	// 检查权限，是否public获取是否处于同一package
	if !c.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	self.class = c
}
