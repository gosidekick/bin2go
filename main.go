package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/crgimenes/goconfig"
)

type config struct {
	OutputFile  string `cfg:"o" cfgRequired:"true" cfgHelper:"output file"`
	PathList    string `cfg:"path" cfgRequired:"true" cfgHelper:"path list"`
	PackageName string `cfg:"pkg" cfgHelper:"package name"`
}

func main() {
	cfg := config{}
	goconfig.PrefixEnv = "bin2go"
	err := goconfig.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("%q: %v\n", path, err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			fmt.Printf("> %v\n", path)
			return nil
		})
	if err != nil {
		fmt.Println(err)
		return
	}
}
