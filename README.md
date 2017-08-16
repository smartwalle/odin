Permission System base Redis

## 简介
使用 Go 语言开发的一套基于 Redis 的权限管理工具。

## 使用说明

#### 初始化
```go
import (
	"github.com/smartwalle/odin"
)

// 主要是配置 Redis 数据库信息
odin.Init(
	Redis 地址,
	Redis 访问密码, 
	Redis 数据库 Index,
	Redis 最大激活连接数量,
	Redis 最大闲置连接数量,
)
```

#### 添加权限信息
```go
odin.NewPermission(
	权限组名称,
	权限名称,
	权限标识符,
)
```
权限添加成功之后，会返回一个权限 id，该 id 为权限标识符 MD5 哈希之后的值。 

#### 获取所有的权限信息
```go
odin.GetPermissionList()
```
该方法会返回当前权限系统所维护的所有权限信息。

#### 添加角色信息
```go
odin.NewRole(
	角色组名称,
	角色名称,
	角色被赋予的权限 id 列表,
)
```
角色添加成功之后，会返回一个角色 id。

#### 获取所有的角色信息
```go
odin.GetRoleList()
```
该方法会返回当前权限系统所维护的所有角色信息。

#### 授权
```go
odin.GrantRole(
	被授权对象标识,
	角色 id 列表,
)
```
为指定对象授予指定角色。

#### 权限验证
```go
odin.Check(
	被授权对象标识,
	权限标识符,
)
```
用于验证指定对象是否拥有指定的权限。