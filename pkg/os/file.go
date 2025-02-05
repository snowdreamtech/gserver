package os

import "os"

// IsExists Check whether the file is Exist
func IsExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir Check whether the file is Dir
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsExistDir Check whether the file is Exist Dir
func IsExistDir(path string) bool {
	if !IsExists(path) {
		return false
	}

	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile Check whether the file is File
func IsFile(path string) bool {
	return !IsDir(path)
}

// IsExistFile Check whether the file is Exist File
func IsExistFile(path string) bool {
	return IsExists(path) && !IsDir(path)
}
