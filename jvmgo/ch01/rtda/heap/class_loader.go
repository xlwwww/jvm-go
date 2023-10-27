package heap

import "fmt"
import "jvmgo/ch01/classfile"
import "jvmgo/ch01/classpath"

/*
class names:
    - primitive types: boolean, byte, int ...
    - primitive arrays: [Z, [B, [I ...
    - non-array classes: java/lang/Object ...
    - array classes: [Ljava/lang/Object; ...
*/
type ClassLoader struct {
	cp       *classpath.Classpath
	verboseFlag bool
	classMap map[string]*Class // loaded classes
}

func NewClassLoader(cp *classpath.Classpath, verboseFlag bool) *ClassLoader {
	loader := &ClassLoader{
		cp:          cp,
		verboseFlag: verboseFlag,
		classMap:    make(map[string]*Class),
	}

	loader.loadBasicClasses()
	// 基本类型类
	loader.loadPrimitiveClasses()
	return loader
}
func (self *ClassLoader) loadPrimitiveClasses(){
	for primitiveType, _ := range primitiveTypes {
		self.loadPrimitiveClass(primitiveType)
	}
}
func (self *ClassLoader) loadPrimitiveClass(className string){
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name:        className,
		loader:      self,
		initStarted: true,
	}
	class.jClass = self.classMap["java/lang/Class"].NewObject()
	class.jClass.extra = class
	self.classMap[className] = class
}

func (self *ClassLoader) loadBasicClasses() {
	jlClassClass := self.LoadClass("java/lang/Class")
	for _, class := range self.classMap {
		if class.jClass == nil {
			class.jClass = jlClassClass.NewObject() // 创建类对象
			class.jClass.extra = class
		}
	}
}

func (self *ClassLoader) LoadClass(name string) *Class {
	if class, ok := self.classMap[name]; ok {
		// already loaded
		return class
	}
	var class *Class
	// 数组类
	if name[0] == '['{
		class = self.loadArrayClass(name)
	}else{
		class = self.loadNonArrayClass(name)
	}
	//所有加载到方法区的类都设置好了jClass字段
	if jlClassClass,ok:=self.classMap["java/lang/Class"];ok{
		class.jClass=jlClassClass.NewObject()
		class.jClass.extra=class
	}
	return class
}
func (self *ClassLoader) loadArrayClass(name string) *Class {
	class := &Class{
		accessFlags: ACC_PUBLIC, // todo
		name:        name,
		loader:      self,
		initStarted: true, // 数组类不需要初始化
		superClass:  self.LoadClass("java/lang/Object"),
		interfaces: []*Class{
			self.LoadClass("java/lang/Cloneable"),
			self.LoadClass("java/io/Serializable"),
		},
	}
	self.classMap[name] = class
	return class
}

func (self *ClassLoader) loadNonArrayClass(name string) *Class {
	data, entry := self.readClass(name)
	class := self.defineClass(data)
	link(class)

	if self.verboseFlag {
		fmt.Printf("[Loaded %s from %s]\n", name, entry)
	}

	return class
}

func (self *ClassLoader) readClass(name string) ([]byte, classpath.Entry) {
	data, entry, err := self.cp.ReadClass(name)
	if err != nil {
		panic("java.lang.ClassNotFoundException: " + name)
	}
	return data, entry
}

// jvms 5.3.5
func (self *ClassLoader) defineClass(data []byte) *Class {
	class := parseClass(data)
	class.loader = self
	resolveSuperClass(class)
	resolveInterfaces(class)
	self.classMap[class.name] = class
	return class
}

func parseClass(data []byte) *Class {
	cf, err := classfile.Parse(data)
	if err != nil {
		//panic("java.lang.ClassFormatError")
		panic(err)
	}
	return newClass(cf)
}

// jvms 5.4.3.1
func resolveSuperClass(class *Class) {
	if class.name != "java/lang/Object" {
		class.superClass = class.loader.LoadClass(class.superClassName)
	}
}
func resolveInterfaces(class *Class) {
	interfaceCount := len(class.interfaceNames)
	if interfaceCount > 0 {
		class.interfaces = make([]*Class, interfaceCount)
		for i, interfaceName := range class.interfaceNames {
			class.interfaces[i] = class.loader.LoadClass(interfaceName)
		}
	}
}

func link(class *Class) {
	verify(class)
	prepare(class)
}

func verify(class *Class) {
	// todo
}

// jvms 5.4.2
func prepare(class *Class) {
	calcInstanceFieldSlotIds(class)
	calcStaticFieldSlotIds(class)
	allocAndInitStaticVars(class)
}
// 计算实例字段个数,计算每个field的在方法表中的index
func calcInstanceFieldSlotIds(class *Class) {
	slotId := uint(0)
	// 如果有父类，字段表包含父类实例变量；实现方法：通过slot index = 父类实例变量数	
	if class.superClass != nil {
		slotId = class.superClass.instanceSlotCount
	}
	for _,field :=range class.fields{
		if !field.IsStatic(){
			field.slotId = slotId
			slotId++
			if(field.isLongOrDouble()){
				slotId++
			}
		}
	}
	class.instanceSlotCount = slotId
}


func calcStaticFieldSlotIds(class *Class) {
	slotId := uint(0)
	for _, field := range class.fields {
		if field.IsStatic() {
			field.slotId = slotId
			slotId++
			if field.isLongOrDouble() {
				slotId++
			}
		}
	}
	class.staticSlotCount = slotId
}
// 从常量池中加载常量值(final && static)
func initStaticFinalVar(class *Class,field *Field) {
	cp:=class.ConstantPool()
	// 常量池中index
	cpIndex := field.ConstValueIndex()
	// 方法表的index
	slotId := field.SlotId()
	//方法表
	vars := class.staticVars
	if cpIndex > 0 {
		switch field.Descriptor() {
		case "Z", "B", "C", "S", "I":
			val := cp.GetConstant(cpIndex).(int32)
			vars.SetInt(slotId, val)
		case "J":
			val := cp.GetConstant(cpIndex).(int64)
			vars.SetLong(slotId, val)
		case "F":
			val := cp.GetConstant(cpIndex).(float32)
			vars.SetFloat(slotId, val)
		case "D":
			val := cp.GetConstant(cpIndex).(float64)
			vars.SetDouble(slotId, val)
		case "Ljava/lang/String;":
			val:=cp.GetConstant(cpIndex).(string)
			jStr:= JString(class.Loader(),val)
			vars.SetRef(slotId,jStr)
		}
	}
}

// 给所有类变量初始化并赋值
func allocAndInitStaticVars(class *Class) {
	class.staticVars = newSlots(class.staticSlotCount)
	for _, field := range class.fields {
		if field.IsStatic() && field.IsFinal() {
			initStaticFinalVar(class, field)
		}
	}
}

