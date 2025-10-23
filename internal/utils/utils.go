package utils

import "strings"

// GetStackTraceAsJSON Возвращает стек вызовов в удобном для чтения виде.
func GetStackTraceAsJSON(stackTrace []byte) string {
	stack := string(stackTrace)

	return strings.ReplaceAll(stack, `\n\t`, "\n")
}
