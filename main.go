package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/crgimenes/goconfig"
)

type config struct {
	OutputFile  string `cfg:"o" cfgRequired:"true" cfgHelper:"output file"`
	PathList    string `cfg:"path" cfgRequired:"true" cfgHelper:"path list"`
	PackageName string `cfg:"pkg" cfgHelper:"package name" cfgDefault:"bin2go"`
}

var files = make(map[string]bool)

func processPath(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Printf("%q: %v\n", path, err)
		return err
	}
	if info.IsDir() {
		return nil
	}
	files[path] = true
	return nil
}

func writeToFile(filename string, payload []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(payload)
	if err != nil {
		return err
	}
	return file.Sync()
}

func main() {
	cfg := config{}
	goconfig.PrefixEnv = "bin2go"
	err := goconfig.Parse(&cfg)
	if err != nil {
		fmt.Println(err)
		return
	}

	pathList := strings.Split(cfg.PathList, ":")
	for _, path := range pathList {
		err = filepath.Walk(path, processPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	var f []string
	for k, _ := range files {
		f = append(f, k)
	}
	sort.Strings(f)

	var bff []byte
	file := bytes.NewBuffer(bff)

	file.WriteString("package " + cfg.PackageName + "\n\n")
	file.WriteString("import \"os\"\n\n")

	b := "var b = [][]byte{\n"
	file.WriteString(b)

	var fileList []string

	n := ""
	for _, v := range f {
		fileList = append(fileList, filepath.Base(v))

		if n != "" {
			file.WriteString(n)
		}
		n = ",\n"
		file.WriteString("{\n")
		data, err := ioutil.ReadFile(v)
		if err != nil {
			log.Fatal(err)
		}
		c := ""
		i := 0
		for _, d := range data {
			if c != "" {
				file.WriteString(c)
			}
			if i > 7 {
				file.WriteString("\n")
				i = 0
			}
			i++
			s := fmt.Sprintf("0x%02X", d)
			file.WriteString(s)
			c = ", "
		}
		file.WriteString(",\n}")
	}
	file.WriteString(",\n}")

	file.WriteString(`
func GetBytes(name string) ([]byte, error) {
switch name {
`)

	for k, v := range fileList {
		m := fmt.Sprintf("case %q:\nreturn b[%v], nil\n", v, k)
		file.WriteString(m)
	}

	file.WriteString(`default:
return nil, os.ErrNotExist
}
}

var listFiles = []string{
`)
	// create file list (todo: remove and use file tree)
	for _, v := range fileList {
		m := fmt.Sprintf("%q,\n", v)
		file.WriteString(m)
	}
	file.WriteString("}\n")
	out, err := format.Source(file.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	writeToFile(cfg.OutputFile, out)
}
