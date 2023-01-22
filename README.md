# Go-DDQuery

### Go语言的Bilibili的DD成分查询

- [x] 通过UID查询
- [x] 通过用户名查询
- [x] 显示所有大航海数据
- [x] 显示直播中的主播
- [x] 通过`HTTP`接口查询, 支持其他`BOT`调用

![](./pic/208259.png)

### 示例代碼:
1.[`在控制臺使用`](main/main.go#L19-L57)

2.[`作爲HttpApi(使用Gin)`](main/main.go#59-L101)

Http請求示例:
```http://127.0.0.1:8964/query?name=嘉然今天吃什么```


```http://127.0.0.1:8964/query?uid=672328094```
