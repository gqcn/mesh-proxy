// 该脚本工具用于将支持SideCar的YAML应用到K8S集群中

package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
	"time"
)

// 切换集群环境
var env = "prod"

func main() {
	var (
		path = "/Users/john/Workspace/fundations/k8s"
	)
	if env != "prod" {
		path = "/Users/john/Workspace/dcm-test/k8s"
	}
	if _, err := gproc.ShellExec(`kubectl config use-context ` + env); err != nil {
		panic(err)
	}
	files, _ := gfile.ScanDirFile(path, "*.yaml", true)
	for _, file := range files {
		content := gfile.GetContents(file)
		if !gstr.ContainsI(content, "Deployment") &&
			!gstr.ContainsI(content, "StatefulSet") &&
			!gstr.ContainsI(content, "DaemonSet") {
			continue
		}
		if !gstr.ContainsI(content, "mesh-proxy") {
			continue
		}
		// 只打印命令不会真正执行
		command := fmt.Sprintf(`kubectl replace -f %s`, file)
		fmt.Println(command)
		continue
		if _, err := gproc.ShellExec(command); err == nil {
			fmt.Println("success replace for:", file)
			time.Sleep(3 * time.Second)
		} else {
			fmt.Printf("!!!!!error: %s, %v\n", file, err)
		}
	}
}
