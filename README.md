## 整合Swagger
参考文档：https://github.com/swaggo/gin-swagger
接口访问地址：http://localhost:8080/swagger/index.html
```text
GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Success 200 {string} json "{"code":"200","msg","","data":""}"
// @Router /problem-list [get]
```

## 安装JWT
```shell
go get github.com/dgrijalva/jwt-go
```