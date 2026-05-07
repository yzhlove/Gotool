package main

import (
	"context"
	"fmt"
	"log"

	"golang.design/x/clipboard"
)

func main() {
	// 1. 初始化剪贴板
	err := clipboard.Init()
	if err != nil {
		log.Fatalf("初始化剪贴板失败: %v", err)
	}

	fmt.Println("开始监听剪贴板... (按 Ctrl+C 退出)")

	// 2. 监听剪贴板内容变化 (返回一个通道)
	// clipboard.FmtText 表示监听纯文本，也支持 FmtImage 等格式
	changeChan := clipboard.Watch(context.Background(), clipboard.FmtText)

	// 3. 循环读取变更
	for data := range changeChan {
		fmt.Printf("监听到剪贴板更新: %s\n", string(data))
	}
}
