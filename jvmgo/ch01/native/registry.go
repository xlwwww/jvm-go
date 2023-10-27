package native

import "jvmgo/ch01/rtda"

type NativeMethod func(frame *rtda.Frame)
var registry = map[string]NativeMethod{}

func Register(className, methodName, methodDescriptor string, method NativeMethod) {
	key := className + "~" + methodName + "~" + methodDescriptor // 类名+方法名+方法描述符唯一确定一个方法
	registry[key] = method
}

func FindNativeMethod(className, methodName, methodDescriptor string) NativeMethod {
	key := className + "~" + methodName + "~" + methodDescriptor
	if method, ok := registry[key]; ok {
		return method
	}
	if methodDescriptor == "()V" && methodName == "registerNatives" {
		return emptyNativeMethod
	}
	return nil
}
func emptyNativeMethod(frame *rtda.Frame) {
	// do nothing
}
