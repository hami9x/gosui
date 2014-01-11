package native

import (
	"fmt"
	"os"
	"path"
	"runtime"
)

var (
	assetDirs []string
)

func LocalDir(relPath string) string {
	_, f, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(f), relPath)
}

func GetFile(filename string) string {
	for _, dir := range assetDirs {
		file := path.Join(dir, filename)
		if _, err := os.Stat(file); os.IsNotExist(err) {
			continue
		}
		return file
	}
	panic(fmt.Sprintf(`Cannot find asset file %v \n
		Searched in %v\n`, filename, assetDirs))
}

func AddAssetDir(dir string) {
	assetDirs = append(assetDirs, dir)
}
