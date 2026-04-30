# gogeneric

简体中文 | [English](README-EN.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/liu-houliang/gogeneric/v2.svg)](https://pkg.go.dev/github.com/liu-houliang/gogeneric/v2)
[![CI](https://github.com/liu-houliang/gogeneric/actions/workflows/ci.yml/badge.svg)](https://github.com/liu-houliang/gogeneric/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://go.dev/)

Go 泛型工具库——提供标准库（`slices`/`maps`）刻意不包含的函数式编程原语，以及常见的业务层工具类型。

**go 版本 >= 1.21**

## 安装

```bash
go get -u github.com/liu-houliang/gogeneric/v2
```

## 功能一览

### 🔧 函数式切片工具（`generic.go`）

这些函数是标准库**不提供**的，也是本库的核心价值所在。

```go
// Filter：保留满足条件的元素
evens := gogeneric.Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })
// → [2 4]

// Map：变换每个元素
strs := gogeneric.Map([]int{1, 2, 3}, func(n int) string { return fmt.Sprintf("%d", n) })
// → ["1" "2" "3"]

// Reduce：折叠聚合
sum := gogeneric.Reduce([]int{1, 2, 3, 4, 5}, 0, func(acc, n int) int { return acc + n })
// → 15

// ForEach：遍历每个元素并执行操作（无返回值）
gogeneric.ForEach([]string{"a", "b"}, func(s string) { fmt.Println(s) })

// Chunk：切片分块
gogeneric.Chunk([]int{1, 2, 3, 4, 5}, 2)
// → [[1 2] [3 4] [5]]

// Flatten：二维切片展平
gogeneric.Flatten([][]int{{1, 2}, {3, 4}, {5}})
// → [1 2 3 4 5]

// Unique：去重（保留首次出现顺序）
gogeneric.Unique([]int{1, 2, 2, 3, 1})
// → [1 2 3]
```

### 🧩 泛型集合类型（`set.go`）

```go
s := gogeneric.NewSet(1, 2, 3)
s.Add(4)
s.Contains(4)          // true
s.Remove(1)

a := gogeneric.NewSet(1, 2, 3)
b := gogeneric.NewSet(2, 3, 4)
a.Union(b)             // {1,2,3,4}
a.Intersection(b)      // {2,3}
a.Difference(b)        // {1}
a.IsSubset(b)          // false
```

### 🎯 Optional 类型（`optional.go`）

Go 缺少 Option/Maybe 类型，这里提供了一个类型安全的实现：

```go
opt := gogeneric.Some("hello")
opt.IsPresent()              // true
opt.GetOrElse("world")       // "hello"

empty := gogeneric.None[string]()
empty.GetOrElse("world")     // "world"
_, err := empty.GetOrError("not found") // err != nil

// 链式变换
length := gogeneric.MapOptional(gogeneric.Some("hello"),
    func(s string) int { return len(s) })
length.GetOrElse(0)          // 5
```

### 🔧 结构体工具（`generic.go`）

```go
type Student struct { Name string; Age uint }
students := []Student{{"alice", 20}, {"bob", 25}}

// 结构体切片转 Map
m := gogeneric.StructSliceToMap[string]("Name", students)
// → map["alice":{alice 20} "bob":{bob 25}]

// 提取某个字段组成切片
names := gogeneric.GetFieldArray[string]("Name", students)
// → ["alice" "bob"]

// 三元运算符
result := gogeneric.Cond(score > 60, "pass", "fail")
```

### 🌐 HTTP 工具（`util.go`）

所有 HTTP 函数返回 `*HTTPResponse[T]`，包含响应数据、状态码和完整响应头：

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// GET — T 可以是 struct、slice、map、pointer，任何 json.Unmarshal 支持的类型
res, err := gogeneric.HttpGetJson[[]User](ctx, "https://api.example.com/users", nil)
res.Data        // []User  — 反序列化后的 body
res.StatusCode  // 200
res.Header.Get("X-RateLimit-Remaining")  // 响应头

// POST JSON body
res, err := gogeneric.HttpPostJson[User](ctx, url, CreateReq{Name: "alice"}, nil)

// PUT / PATCH JSON body
res, err := gogeneric.HttpPutJson[User](ctx, url+"/1", UpdateReq{Name: "bob"}, nil)
res, err := gogeneric.HttpPatchJson[User](ctx, url+"/1", map[string]string{"name": "bob"}, nil)

// DELETE（支持 204 空 body）
res, err := gogeneric.HttpDeleteJson[DeletedResp](ctx, url+"/1", nil)

// POST 表单（OAuth、老接口）
form := url.Values{"grant_type": {"password"}, "username": {"alice"}}
res, err := gogeneric.HttpPostForm[TokenResp](ctx, "https://auth.example.com/token", form, nil)

// 文件上传 multipart/form-data
files := []gogeneric.File{
    {FieldName: "avatar", FileName: "photo.jpg", Reader: f},
}
res, err := gogeneric.HttpPostMultipart[UploadResp](ctx, url, files, map[string]string{"userId": "1"}, nil)
```

### ⚠️ 已废弃函数（兼容保留）

以下函数已被 Go 1.21+ 标准库覆盖，建议迁移：

| 本库函数 | 标准库替代 |
|---------|-----------|
| `Min` / `Max` | 内置 `min()` / `max()` |
| `SortSlice` | `slices.SortFunc` |
| `SliceIndex` | `slices.Index` |
| `CompareSlice` | `slices.Equal` |
| `MapKeys` | `maps.Keys` |
| `MapToSlice` | `maps.Values` |

## 运行测试

```bash
git clone https://github.com/liu-houliang/gogeneric.git
cd gogeneric

go test -v .                    # 全部测试（详细输出）
go test -run TestFilter .       # 只跑 Filter 相关测试
go test -run TestHttp .         # 只跑 HTTP 相关测试
go test -bench=. -benchtime=3s . # 跑 benchmark
```

## 注意事项

- `StructSliceToMap` 和 `GetFieldArray` 使用了反射，请确保字段名和类型参数匹配正确。
- `Optional.MustGet()` 在值不存在时会 panic，生产代码推荐使用 `Get()` 或 `GetOrElse()`。
- 所有类型均非并发安全，如需并发访问请自行加锁。

## 项目负责人

[@liu-houliang](https://github.com/liu-houliang)