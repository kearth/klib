package main

import (
	"fmt"
	"github.com/kearth/ktools/web"
)

func main() {
	jsonStr := `{"host": "http://localhost:9090","port": 9090,"analytics_file": "","static_file_version": 1,"static_dir": "E:/Project/goTest/src/","templates_dir": "E:/Project/goTest/src/templates/","serTcpSocketHost": ":12340","serTcpSocketPort": 12340,"fruits": ["apple", "peach"]}`
	fmt.Println(web.JSONFormat(jsonStr))
}
