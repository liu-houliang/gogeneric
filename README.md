# gogeneric
简体中文 | [English](README-EN.md)

可以直接使用的go语言泛型和工具函数

go 版本 >= 1.18


## 如何使用
```
go get -u github.com/liu-houliang/gogeneric

max:= gogeneric.Max(3,4)
```

## 目录结构
```
test   测试代码
    -- generic_test.go  
    -- util_test.go     
    
generic.go       泛型函数
util.go          工具函数
```

## 泛型函数列表
``` go
func Min[T Number](base, comp T) T

func Max[T Number](base, comp T) T

func Cond[T any](condition bool, a, b T) T 

func SortSlice[T any](slice []T, less func(T, T) bool) 

func StructSliceToMap[Key comparable, Value any](key string, slice []Value) map[Key]Value 

func GetFieldArray[Field any, T any](fieldName string, slice []T) []Field 
...
```

## 工具函数列表
``` go
func HttpGetJson[T any](apiURL string, headers map[string]string) (*T, error)

func HttpPostJson[T any](apiURL string, request any, headers map[string]string) (response T, err error)
```


## 测试方法
```
git clone https://github.com/liu-houliang/gogeneric.git
cd gogeneric/test
go test -v 
```

## 警告
部分函数使用了反射，需要足够的测试才能保证正确性


## 项目负责人
[@liu-houliang](https://github.com/liu-houliang)