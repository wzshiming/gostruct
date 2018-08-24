package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/wzshiming/gostruct"
	"github.com/wzshiming/namecase"
)

var f = flag.String("f", "", "json file path")
var name = flag.String("n", "Foo", "struct name")

func main() {
	flag.Parse()
	fn := *f
	var data []byte
	var err error
	if fn != "" {
		data, err = ioutil.ReadFile(fn)
		if err != nil {
			return
		}
	} else {
		data, err = ioutil.ReadAll(os.Stdin)
		if err != nil {
			return
		}
	}
	if len(data) == 0 {
		flag.Usage()
		return
	}

	var i interface{}
	json.Unmarshal(data, &i)
	n := *name
	if n == "" {
		_, n = filepath.Split(fn)
		n = strings.SplitN(n, ".", 2)[0]
		if n == "" {
			flag.Usage()
			return
		}
		n = namecase.ToUpperHumpInitialisms(n)
	}

	code := gostruct.GenStruct(n, i)
	fmt.Println(code)
}
