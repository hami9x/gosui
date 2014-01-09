package jsbuild

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"text/template"

	"go/build"
	gopherjs "github.com/neelance/gopherjs/api"
)

const defaultDistPath = "dist/web"

func GetFabricBackendJsFiles() (l []string) {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Something is wrong.")
	}
	thisDir := path.Dir(f)
	srcDir := path.Join(path.Dir(thisDir), "htmlcanvas")
	l = make([]string, 0)
	l = append(l, path.Join(srcDir, "interface.js"))
	l = append(l, path.Join(srcDir, "fabricjs/dist/fabric.js"))
	return l
}

func Build() {
	cwd, _ := os.Getwd()
	appDir := path.Dir(cwd)
	fmt.Printf("Generating html/js for %v \n", appDir)
	pkg, err := build.ImportDir(appDir, 0)
	if err != nil {
		panic(`Cannot build javascript, make sure you're building the program
		 from its true package directory.`)
	}
	outputDir := path.Join(appDir, defaultDistPath)
	jsPath := path.Join(outputDir, "gosui.js")
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		panic("Something is wrong.")
	}
	srcDir := path.Dir(f)
	htmlOPath := path.Join(outputDir, "index.html")
	gjsPkg := &gopherjs.Package{Package: pkg}
	err = gopherjs.BuildPackage(gjsPkg)
	if err != nil {
		panic(err.Error())
	}
	gopherjs.WriteCommandPackage(gjsPkg, jsPath)
	htmlOFile, err := os.Create(htmlOPath)
	//panic(err.Error())
	t, err := template.New("template.html").ParseFiles(path.Join(srcDir, "template.html"))

	err = t.Execute(htmlOFile, nil)
	if err != nil {
		fmt.Printf("%v %v", err.Error())
		return
	}
	htmlOFile.Close()
	jsFiles := GetFabricBackendJsFiles()
	for _, srcPath := range jsFiles {
		dstPath := path.Join(outputDir, path.Base(srcPath))
		dst, _ := os.Create(dstPath)
		src, err := os.Open(srcPath)
		if err != nil {
			panic(err.Error())
		}
		io.Copy(dst, src)
		fmt.Printf("copied file %v to %v.\n", srcPath, dstPath)
	}
}
