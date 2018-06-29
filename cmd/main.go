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
	"flag"
)

func main() {
	output := flag.String("o", "api", "output dir")

	flag.Parse()

	p, _ := os.Getwd()
	pathPackage := p

	if len(flag.Args()) > 0 {
		pathPackage = path.Join(p, flag.Args()[0])
	}

	packageFiles, err := ioutil.ReadDir(pathPackage)
	if err != nil {
		log.Fatal(err)
	}

	var entites []xgravity.Entity

	gopath := build.Default.GOPATH
	packageFullName := strings.Replace(pathPackage, gopath+"/src/", "", 1)

	for _, file := range packageFiles {
		if !file.IsDir() && path.Ext(file.Name()) == ".go" {
			b, err := ioutil.ReadFile(path.Join(pathPackage, file.Name()))
			if err != nil {
				log.Fatal(err)
			}

			e, err := xgravity.GetEntities(packageFullName, path.Join(pathPackage, file.Name()), b)
			if err != nil {
				log.Fatal(err)
			}

			entites = append(entites, e...)
		}
	}

	files, err := xgravity.CreateApiPackage(packageFullName, entites)
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir(*output, os.ModePerm)

	for name, content := range files {
		err = ioutil.WriteFile(*output+"/"+name, content, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}
}

