package log

import "os"

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil { // path exists
		return true, nil
	} else if os.IsNotExist(err) { // error is 'not exist'
		return false, nil
	}
	return false, err // other error
}

func createDirsIfNotExist(dirs ...string) (err error) {
	for _, dir := range dirs {
		if existing, pathExistsErr := pathExists(dir); !existing && pathExistsErr == nil {
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return
			}
		}
	}
	return
}
