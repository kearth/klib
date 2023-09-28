package kcomp

// ComponentType 组件类型
type ComponentType string

const (
	CoreComponent      ComponentType = "core"       // 核心组件
	UserComponent      ComponentType = "user"       // 用户组件-允许自定义
	ThirdPartComponent ComponentType = "third-part" // 第三方组件
)

// String
func (ct ComponentType) String() string {
	return string(ct)
}

// IComponent
type IComponent interface {
	// 注册组件
	Register(name string) error
	// 组件类型
	Type() ComponentType
}
