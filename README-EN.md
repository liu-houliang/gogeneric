# gogeneric
some go generic code as utils

go version >= 1.18

## How to use 
```
go get -u github.com/liu-houliang/gogeneric

max:= gogeneric.Max(3,4)
```

## Dictionary
```
generic.go       some go generic functions
generic_test.go  go test code
```

## Functions list
``` go
func Min[T Number](base, comp T) T

func Max[T Number](base, comp T) T

func Cond[T any](condition bool, a, b T) T 

func SortSlice[T any](slice []T, less func(T, T) bool) 

func StructSliceToMap[Key comparable, Value any](key string, slice []Value) map[Key]Value 

func GetFieldArray[Field any, T any](fieldName string, slice []T) []Field 
...
```



## How to test 
```
git clone https://github.com/liu-houliang/gogeneric.git
cd gogeneric/test
go test -v 
```

## Warning
some functions use reflect code , these functions should be tested before using in production environment.

## Project Manager
[@liu-houliang](https://github.com/liu-houliang)