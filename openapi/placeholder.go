// Package openapi 定义 App 和 Admin 接口文档的生成入口。
// 在 Handler 的 Swag 注释中用 @x-module-app true 或 @x-module-admin true 指定所属模块。
// 两端共用 Handler 时，各写一块注释，用空行分隔，每块注明对应模块和完整的 @Router 路径及方法。
//
// 在仓库根目录执行 go generate ./openapi，生成器按模块标记分别生成 app_docs.go 和 admin_docs.go。
// bootstrap/service.go 匿名导入本包，生成文件的 init 函数将两份文档以 app、admin 为名注册到 Swag。
//
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --dir .. --generalInfo openapi/app.go --output . --outputTypes go --instanceName app --parseExtension module-app --parseInternal --parseFuncBody
//go:generate go run github.com/swaggo/swag/cmd/swag@v1.16.6 init --dir .. --generalInfo openapi/admin.go --output . --outputTypes go --instanceName admin --parseExtension module-admin --parseInternal --parseFuncBody
package openapi
