package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
)

var (
	contextUsed    = "test" // 需要执行的配置名称
	contextRecover = "test" // 执行后需要切换的配置
)

// 该脚本用于展示所有安装有mesh-proxy sidecar容器的服务。
// 由于接下来往往会使用patch进行滚动更新，因此这里输出了patch命令。
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
	namespaces := g.SliceStr{"app", "infra"}
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
				if gstr.Contains(array[2], "mesh-proxy") {
					command := fmt.Sprintf(
						`kubectl patch -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"%s\"}}}}}" -n %s %s %s/%s`,
						gtime.TimestampNanoStr(), namespace, "\t", array[0], array[1],
					)
					fmt.Println(command)
				}
			}
		}
	}
}
