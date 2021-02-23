// 该脚本工具用于将普通YAML修改为支持SideCar的YAML配置

package main

import (
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

// 注意当前操作环境
var env = "prod"

func main() {
	var (
		path = "/Users/john/Workspace/fundations/k8s"
	)
	if env != "prod" {
		path = "/Users/john/Workspace/dcm-test/k8s"
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

		match, err := gregex.MatchString(`([#\s]*initContainers:[\s\S]+#[\s#]+应用.+[\s#]+)`, content)
		if err != nil {
			panic(err)
		}
		if len(match) > 1 {
			replaceContent := "\n"
			for _, v := range gstr.Split(match[1], "\n") {
				if gstr.Trim(v) == "" {
					continue
				}
				if gstr.Trim(v) == "containers:" {
					replaceContent += v
				} else {
					replaceContent += "#" + v
				}
				replaceContent += "\n"
			}
			replaceContent += "      "
			gfile.PutContents(file, gstr.Replace(content, match[1], replaceContent))
		}
	}
}
