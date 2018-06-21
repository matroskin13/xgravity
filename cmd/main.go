package main

import (
	//"github.com/matroskin13/xgravity"
	"io/ioutil"
	"log"
	"os"
	"go/build"

	"github.com/matroskin13/xgravity"
	"path"
	"strings"
)

func main() {
	p, _ := os.Getwd()
	pathPackage := p

	if len(os.Args) > 1 {
		pathPackage = path.Join(p, os.Args[1])
	}

	packageFiles, err := ioutil.ReadDir(pathPackage)
	if err != nil {
		log.Fatal(err)
	}

	var entites []xgravity.Entity

	for _, file := range packageFiles {
		if !file.IsDir() {
			b, err := ioutil.ReadFile(path.Join(pathPackage, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			e, err := xgravity.GetEntities(path.Join(pathPackage, file.Name()), b)
			if err != nil {
				log.Fatal(err)
			}

			entites = append(entites, e...)
		}
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}

	packageFullName := strings.Replace(pathPackage, gopath+"/src/", "", 1)
	files, err := xgravity.CreateApiPackage(packageFullName, entites)
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir("api", os.ModePerm)

	for name, content := range files {
		err = ioutil.WriteFile("./api/"+name, content, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

