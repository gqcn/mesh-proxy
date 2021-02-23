// 该脚本工具用于将普通YAML修改为支持SideCar的YAML配置

package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
)

// 注意当前操作环境
var env = "test"

var meshContent = fmt.Sprintf(`
      initContainers:
        ###############################################
        # 网关中间件iptables
        ################################################
        - name : init
          image: "loads/mesh-proxy:init"
          imagePullPolicy: "Always"
          securityContext:
            privileged: true

      containers:
      ################################################
      ## 网关中间件容器
      #################################################
      - name : mesh
        image: "loads/mesh-proxy:%s"
        imagePullPolicy: "Always"
        readinessProbe:
          initialDelaySeconds:	5
          periodSeconds: 5
          exec:
            command: ["health.sh"]
        livenessProbe:
          initialDelaySeconds:	5
          periodSeconds: 5
          exec:
            command: ["health.sh"]
        volumeMounts:
        - name      : logmesh
          mountPath : /var/log/www

      ################################################
      # 应用容器
      ################################################
      `, env)

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
		if gstr.ContainsI(content, "mesh-proxy") {
			continue
		}
		// 替换 volumns
		if !gstr.ContainsI(content, "logmesh") {
			match, _ := gregex.MatchString(`((\s+)\-\s*name\s*:\s*\w*log\w*\s*emptyDir\s*:\s*\{\})`, content)
			if len(match) > 1 {
				tab := match[2]
				content = gstr.Replace(content, match[1], fmt.Sprintf("%s%s%s%s  %s", match[1], tab, "- name    : logmesh", tab, "emptyDir: {}"))
				//gfile.PutContents(file, content)
			}
		}

		// 替换 containers
		match, err := gregex.MatchString(`(#[\s#]+应用[\s\S]+)(\s+)containers:\s*`, content)
		if err != nil {
			panic(err)
		}
		if len(match) > 2 {
			// 应用容器注释在containers标签前面
			content = gstr.Replace(content, match[0], meshContent)
			gfile.PutContents(file, content)
		} else {
			// 应用容器注释在containers标签后面
			match, err := gregex.MatchString(`((\s+)containers:\s*#[\s#]+应用.+[\s#]+)`, content)
			if err != nil {
				panic(err)
			}
			if len(match) > 2 {
				content = gstr.Replace(content, match[0], meshContent)
				gfile.PutContents(file, content)
			}
		}
	}
}
