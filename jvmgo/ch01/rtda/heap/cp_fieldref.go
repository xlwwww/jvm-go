package heap

import "jvmgo/ch01/classfile"

type FieldRef struct {
	MemberRef
	field *Field
}

func newFieldRef(cp *ConstantPool, refInfo *classfile.ConstantFieldrefInfo) *FieldRef {
	ref := &FieldRef{}
	ref.cp = cp
	ref.copyMemberRefInfo(&refInfo.ConstantMemberrefInfo)
	return ref
}

func (self *FieldRef) ResolvedField() *Field {
	if self.field == nil {
		self.resolveFieldRef()
	}
	return self.field
}

// jvms 5.4.3.2
// 类D想通过字段符号引用访问类C的某个字段
func (self *FieldRef) resolveFieldRef() {
	// 当前类
	d := self.cp.class
	// 字段所属类
	c := self.ResolvedClass()
	field := lookupField(c, self.name, self.descriptor)

	if field == nil {
		panic("java.lang.NoSuchFieldError")
	}
	// 检查类d是否具有权限
	if !field.isAccessibleTo(d) {
		panic("java.lang.IllegalAccessError")
	}

	self.field = field
}

func lookupField(c *Class, name, descriptor string) *Field {
	// 从类c的字段表找
	for _, field := range c.fields {
		if field.name == name && field.descriptor == descriptor {
			return field
		}
	}
	// 从类c的接口找
	for _, iface := range c.interfaces {
		if field := lookupField(iface, name, descriptor); field != nil {
			return field
		}
	}
	// 父类找
	if c.superClass != nil {
		return lookupField(c.superClass, name, descriptor)
	}

	return nil
}
