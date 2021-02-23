package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/text/gstr"
)

var (
	contextUsed    = "test"                     // 需要执行的配置名称
	contextRecover = "test"                     // 执行后需要切换的配置
	namespaces     = g.SliceStr{"app", "infra"} // 检索的命名空间
)

// 该脚本用于搜索没有安装有mesh-proxy sidecar容器的服务。
func main() {
	content := ""
	columns := `-o=custom-columns=LABELS:.kind,NAME:.metadata.name,DATA:'.spec.template.spec.initContainers[0].image'`
	if _, err := gproc.ShellExec(fmt.Sprintf(`kubectl config use-context %s`, contextUsed)); err != nil {
		panic(err)
	}
	defer func() {
		if _, err := gproc.ShellExec(fmt.Sprintf(`kubectl config use-context %s`, contextRecover)); err != nil {
			panic(err)
		}
	}()
	for _, namespace := range namespaces {
		content = ""
		command1 := fmt.Sprintf(`kubectl get deployment  %s -n %s`, columns, namespace)
		command2 := fmt.Sprintf(`kubectl get statefulset %s -n %s`, columns, namespace)
		command3 := fmt.Sprintf(`kubectl get daemonset   %s -n %s`, columns, namespace)
		if c, err := gproc.ShellExec(command1); err == nil {
			content += c
		}
		if c, err := gproc.ShellExec(command2); err == nil {
			content += c
		}
		if c, err := gproc.ShellExec(command3); err == nil {
			content += c
		}
		for _, line := range gstr.SplitAndTrim(content, "\n") {
			array := gstr.SplitAndTrim(line, " ")
			if len(array) == 3 {
				if !gstr.Contains(array[2], "mesh-proxy") {
					if array[0] == "LABELS" {
						continue
					}
					fmt.Printf("%s\t%s/%s\n", namespace, array[0], array[1])
				}
			}
		}
	}
}
