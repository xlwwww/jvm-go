package heap

import "unicode/utf16"
// key:go字符串 value:java字符串
var internedStrings = map[string]*Object{}
//根据Go字符串返回相应的Java字符串实例，如果java字符串已经在池中，则直接返回
func JString(loader *ClassLoader,goStr string) *Object{
	if internedStr,ok :=internedStrings[goStr];ok{
		return internedStr
	}
	chars := stringToUtf16(goStr) // 把go字符串转换成java字符数组
	jChars:=&Object{
		class:loader.LoadClass("[C"),// 数组类
		data:chars,
	}
	jStr := loader.LoadClass("java/lang/String").NewObject() // string类
	jStr.SetRefVar("value", "[C", jChars) // 获取jstr的value字段并赋值为jchars
	internedStrings[goStr] = jStr
	return jStr
}

func stringToUtf16(s string) []uint16 {
	runes:=[]rune(s) // 强转utf-32
	return utf16.Encode(runes) // 编码utf-16
}
func GoString(jstr *Object) string{
	charArr:=jstr.GetRefVar("value","[C")
	return utf16ToString(charArr.Chars())
}
// utf16 -> utf8
func utf16ToString(s []uint16) string {
	runes := utf16.Decode(s) // func Decode(s []uint16) []rune
	return string(runes)
}
// todo
func InternString(jStr *Object) *Object {
	goStr := GoString(jStr)
	if internedStr, ok := internedStrings[goStr]; ok {
		return internedStr
	}

	internedStrings[goStr] = jStr
	return jStr
}
