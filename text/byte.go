package text

import "unicode/utf8"

var whiteList = [32]bool{
	'\t': true, // 制表符
	'\n': true, // 换行符
	'\v': true, // 制表符
	'\f': true, // 换页符
	'\r': true, // 回车符
}

// IsPlainText 判断给定的字节数组是否是有效的纯文本。
// 纯文本是指不包含 NUL 字符、非法的 UTF-8 字符和无法显示的控制字符（除了常见的换行符、回车符、水平制表符和换页符）。
func IsPlainText(data []byte) bool {
	for len(data) > 0 {
		r, size := utf8.DecodeRune(data)

		// 1. 校验 UTF-8 格式
		// size == 1 && r == utf8.RuneError 表示解码失败
		if r == utf8.RuneError && size == 1 {
			return false
		}

		// 2. 排除 DEL 字符
		if r == 127 {
			return false
		}

		// 3. 处理控制字符（ASCII 0-31）
		// 允许白名单：换行、回车、制表符和换页符
		if r < 32 && !whiteList[r] {
			return false
		}

		// 更新 data，继续解码下一个字符
		data = data[size:]
	}

	// 如果通过所有检查，返回 true
	return true
}
