package main

import (
	"fmt"
)

// 判断给定的字符串是否可读（不包含特殊字符）。
func isReadable(content []byte) bool {
	for i, c := range content {
		if (c < ' ' || c > '~') && c != '\n' && c != '\r' && c != '\t' {
			fmt.Println(i, string(c))
			return false
		}
	}
	return true
}

var content = []byte(`{"errcode":0,"errmsg":"success","data":{"list":[{"id":0,"user_id":19939006,"name":"王路明","avatar":"","hospital":"西安医学院第二附属医院","section":"口腔科"},{"id":0,"user_id":30613895,"name":"王路","avatar":"","hospital":"东北国际医院","section":"内分泌科"},{"id":0,"user_id":31213716,"name":"王路遥","avatar":"","hospital":"中国人民解放军第四五五医院","section":"妇产科"},{"id":0,"user_id":32691603,"name":"王路阳","avatar":"","hospital":"浙江省妇保医院","section":"麻醉科"},{"id":0,"user_id":38473242,"name":"王路平","avatar":"","hospital":"梨树县第一人民医院","section":"心脏内科"},{"id":0,"user_id":45329383,"name":"王路","avatar":"","hospital":"天津市肿瘤医院","section":"肿瘤科"},{"id":0,"user_id":50758607,"name":"王路强","avatar":"","hospital":"大名县人民医院","section":"康复理疗科"},{"id":0,"user_id":50863161,"name":"王路","avatar":"","hospital":"","section":"关节外科"},{"id":0,"user_id":51281085,"name":"#超神王路飞","avatar":"","hospital":"","section":""},{"id":0,"user_id":51335589,"name":"王路","avatar":"","hospital":"","section":""},{"id":0,"user_id":51474984,"name":"王路路","avatar":"","hospital":"","section":""},{"id":0,"user_id":52556199,"name":"王路飞","avatar":"","hospital":"北京天健医院","section":"内科"},{"id":0,"user_id":52690956,"name":"王路飞","avatar":"","hospital":"","section":""},{"id":0,"user_id":52785006,"name":"王路","avatar":"","hospital":"","section":""},{"id":0,"user_id":52948836,"name":"王路","avatar":"","hospital":"","section":""},{"id":0,"user_id":53751485,"name":"帝王路","avatar":"","hospital":"","section":""},{"id":0,"user_id":58559709,"name":"王路平","avatar":"","hospital":"大连大学附属中山医院","section":"内科"},{"id":0,"user_id":58564649,"name":"王路","avatar":"","hospital":"常州同济男科医院","section":"泌尿外科"},{"id":0,"user_id":58728517,"name":"王路路","avatar":"","hospital":"滨州医学院附属医院","section":"重症医学科"},{"id":0,"user_id":59216892,"name":"王路","avatar":"","hospital":"新乡医学院第三附属医院","section":"内科"},{"id":0,"user_id":66361603,"name":"王路尧","avatar":"","hospital":"南京医科大学附二院","section":"消化内科"},{"id":0,"user_id":68108639,"name":"王路","avatar":"","hospital":"","section":"内分泌科"},{"id":0,"user_id":69108221,"name":"王路杰","avatar":"","hospital":"东部战区总医院","section":"骨科"},{"id":0,"user_id":69302157,"name":"王路娥","avatar":"","hospital":"河北省沧州中西医结合医院","section":"急诊医学科"}],"hasMore":false}}

`)

func main() {
	fmt.Println(string(content[70:85]))
	fmt.Println(isReadable(content))
}
