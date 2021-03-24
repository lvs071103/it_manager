package tools

import (
	"os"
)

func WriteFile(content, fileName string) bool {
	f, cErr := os.Create(fileName)
	defer func() {
		_ = f.Close()
	}()
	_, wErr := f.WriteString(content)

	if cErr == nil && wErr == nil {
		return true
	} else {
		return false
	}
}
