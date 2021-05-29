package core

import "testing"

func BenchmarkConvertRemoteImageToLocal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		maintainImageTagsForSingleFile(
			"/Users/zhongrui/Library/Mobile Documents/iCloud~app~cyan~taio/Documents/Editor/notebook/notes/技术/golang/语言/02.变量&常量.md",
			"/Users/zhongrui/Library/Mobile Documents/iCloud~app~cyan~taio/Documents/Editor/notebook/attachments",
			false,
			true)
	}
}
