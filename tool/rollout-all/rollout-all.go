package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gproc"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"time"
)

var (
	env = "test"
)

// 该脚本用于滚动更新所有的安装有mesh-proxy sidecar容器的服务。
// !!风险较大，请谨慎使用!!
func main() {
	content := ""
	columns := `-o=custom-columns=LABELS:.kind,NAME:.metadata.name,DATA:'.spec.template.spec.initContainers[0].image'`
	if _, err := gproc.ShellExec(fmt.Sprintf(`kubectl config use-context %s`, env)); err != nil {
		panic(err)
	}
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
					command := fmt.Sprintf(`kubectl patch %s/%s -p "{\"spec\":{\"template\":{\"metadata\":{\"labels\":{\"date\":\"%s\"}}}}}" -n %s`, array[0], array[1], gtime.TimestampNanoStr(), namespace)
					if _, err := gproc.ShellExec(command); err == nil {
						fmt.Println("success patch for:", namespace, array[0]+"/"+array[1])
						time.Sleep(3 * time.Second)
						continue
					} else {
						fmt.Println("!!!!!error:", namespace, array[0]+"/"+array[1])
					}
				}
			}
		}
	}
}
