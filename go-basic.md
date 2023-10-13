### OOP
#### 封装
##### 定义类
```go
// 借助结构体实现
type Student struct {
    id uint
    name string
    male bool
    score float64
}
// go不存在构造函数，自定义公共方法
func NewStudent(id uint, name string, male bool, score float64) *Student {
    return &Student{id, name, male, score}
}
```
##### 定义成员方法
```go
type Student struct {
    id    uint
    name  string
    score float64
}

func NewStudent(id uint, name string, score float64) *Student {
    return &Student{id: id, name: name, score: score}
}

func NewStudentV2(id uint, name string, score float64) Student {
    return Student{id: id, name: name, score: score}
}

func (s Student) GetName() string {
    return s.name
}

func (s *Student) SetName(name string) {
    s.name = name
}

func main() {
    s := NewStudent(1, "学院君", 100)
    s.SetName("学院君1号")   // ok 正常调用指针方法
    fmt.Println(s.GetName()) // ok 指针调用值方法自动解引用: (*s).GetName()

    s2 := NewStudentV2(2, "学院君", 90)
    s2.SetName("学院君2号")   // ok s2 是可寻址的左值，所以实际调用: (&s2).SetName("学院君2号")
    fmt.Println(s2.GetName()) // ok 正常调用值方法

    NewStudent(3, "学院君", 80).SetName("学院君3号")   // ok 正常调用指针方法
    NewStudentV2(4, "学院君", 99).SetName("学院君4号") // err 值类型调用指针方法，左值非可寻址
}
```
- 值方法可以通过指针和值类型实例调用，指针类型实例调用值方法时会自动解引用
- 指针方法只能通过指针类型实例调用，但有一个例外，如果某个值是可寻址的（或者说左值），那么编译器会在值类型实例调用指针方法时自动插入取地址符，使得在此情形下看起来像指针方法也可以通过值来调用。
- 左值：出现在等号左边的值
- 可寻址：通过编译器直接寻址的，将a变为&a


#### 继承
```go
type Animal struct {
    Name string
}
// Dog具有Animal的所有方法和状态
type Dog struct {
    Animal 
}
```