package base

import "os"

// DirExists -
func DirExists(path *string) bool {
	info, err := os.Stat(*path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// FileExists -
func FileExists(path *string) bool {
	info, err := os.Stat(*path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
