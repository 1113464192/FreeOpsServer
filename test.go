package main

import (
	"fmt"
	"strings"
)

func main() {
	str := "layout.base$view.home"
	parts := strings.Split(str, "$")
	fmt.Println(parts) // 输出: [layout base$view.home]

	//	判断是否有layout的部分在view部分前面，因为要求layout一定是在view后面的，可能有多个$，如果有view部分段在layout的前面，就报错
	if len(parts) < 2 {
		fmt.Println("layout and view must be separated by $")
		return

	}

	var viewPart string
	for _, part := range parts {
		if strings.Contains(part, "view") {
			viewPart = part
			break
		}
	}
	fmt.Println(viewPart) // 输出: base$view.home
}
