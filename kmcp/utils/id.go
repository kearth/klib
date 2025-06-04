package utils

import "github.com/gogf/gf/v2/util/guid"

// GenID 生成一个唯一的 ID
func GenID() string {
	return guid.S([]byte("kmcp"))
}
