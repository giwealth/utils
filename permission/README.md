# 通过注释生成ULA权限列表文件

## 接口函数示例
```
// List ECS列表
// @Permission 资源中心:云服务器:服务器列表
// @Router /api/clouds/resources/ecs [get]
func (c *ecsController) List(ctx context.Context) error {
	return nil
}

// Detail ECS详情
// @Permission 资源中心:云服务器:服务器详情
// @Router /api/clouds/resources/ecs/detail [get]
func (c *ecsController) Detail(ctx context.Context) error {
	return nil
}
```

## 调用示例
```
package main

import (
    "github.com/giwealth/utils/permission"
)

func main() {
    // ./internal 为接口go文件所在目录; ./configs 为生成的权限列表json文件保存目录
    if err := permission.Generate("./internal", "./configs"); err != nil {
        panic(err)
    }
}
```

## 使用方法
- 使用embed静态文件包
```
// configs/config.go
package configs

import (
	_ "embed" // 导入内嵌静态文件包
)

var (
	// PermissionJSON 默认权限文件
	//go:embed permission.json
	PermissionJSON []byte
)
```

- 权限列表接口函数
```
// permission list api
type Permission struct {
    Code string `json:"code"`         // 权限代码
    Desc string `json:"desc"`         // 权限描述，分组使用冒号分隔
    Type string `json:"type"`         // 权限类型, 值为: 系统或自定义
}

func PermissionList() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        var permissions []*Permission
        if err := json.Unmarshal(configs.PermissionJSON, &permissions); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }

        ctx.JSON(http.StatusOK, gin.H{"permissions": permissions})
    }
}
```