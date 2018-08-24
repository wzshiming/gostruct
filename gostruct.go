package gostruct

import (
	"bytes"
	"fmt"
	"go/format"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/wzshiming/namecase"
)

// GenStruct
func GenStruct(name string, i interface{}) string {
	g := NewGen(name)
	g.Add(name, reflect.ValueOf(i))
	f := g.Gen()
	d, err := format.Source([]byte(f))
	if err != nil {
		return f
	}
	return string(d)
}

// NewGen
func NewGen(name string) *Gen {
	return &Gen{
		Name:  name,
		Types: map[string]*bytes.Buffer{},
	}
}

type Gen struct {
	Name  string
	Types map[string]*bytes.Buffer
}

// Gen
func (g *Gen) Gen() string {
	buf := []string{}
	for name, v := range g.Types {
		buf = append(buf, fmt.Sprintf("// %s generated\ntype %s %s", name, name, v.String()))
	}
	sort.Strings(buf)
	return strings.Join(buf, "\n")
}

func (g *Gen) isStar(name string) string {
	if strings.HasPrefix(name, "*") {
		return name
	}
	if _, ok := g.Types[name]; ok {
		return "*" + name
	}
	return name
}

// Add
func (g *Gen) Add(name string, val reflect.Value) string {
	switch kind := val.Kind(); kind {
	case reflect.Float32, reflect.Float64:
		return "json.Number"
	case reflect.Bool,
		reflect.Complex64, reflect.Complex128,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return strings.ToLower(kind.String())
	case reflect.String:
		v := val.String()
		if _, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return "time.Time"
		}
		return "string"
	case reflect.Slice:
		if val.Len() == 0 {
			return fmt.Sprintf("[]%s", g.Add(name, reflect.New(val.Type().Elem()).Elem()))
		}
		return fmt.Sprintf("[]%s", g.Add(name, val.Index(0)))
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", val.Len(), g.Add(name, val.Elem()))
	case reflect.Ptr, reflect.Interface:
		return g.isStar(g.Add(name, val.Elem()))
	case reflect.Map:
		mk := val.MapKeys()
		if len(mk) == 0 {
			return "json.RawMessage"
		}
		valueSlice(mk).Sort()
		named := map[string]int{}
		buf := bytes.NewBuffer(nil)
		buf.WriteString("struct {\n")
		for _, k := range mk {
			v := val.MapIndex(k)
			name := k.String()
			newName := namecase.ToUpperHumpInitialisms(name)
			if named[newName] != 0 {
				newName = fmt.Sprintf("%s%d", newName, named[newName]+1)
			}
			named[newName]++
			newTypeName := namecase.ToUpperHumpInitialisms(g.Name + " " + name)
			newTypeName = g.Add(newTypeName, v)
			buf.WriteString(fmt.Sprintf("%s %s `json:\"%s,omitempty\"`\n", newName, newTypeName, name))
		}
		buf.WriteString("}")
		g.Types[name] = buf
		return g.isStar(name)
	case reflect.Struct:
		buf := bytes.NewBuffer(nil)
		buf.WriteString("struct {\n")
		typ := val.Type()
		num := typ.NumField()
		for i := 0; i != num; i++ {
			t := typ.Field(i)
			v := val.Field(i)
			buf.WriteString(fmt.Sprintf("%s %s %s\n", t.Name, g.Add(t.Name, v), string(t.Tag)))
		}
		buf.WriteString("}")
		g.Types[name] = buf
		return g.isStar(name)
	case reflect.Chan, reflect.Func:
		// No action
		return "json.RawMessage"
	default:
		return "json.RawMessage"
	}
}
