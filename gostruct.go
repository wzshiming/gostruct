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
func GenStruct(name string, i interface{}) []byte {
	g := NewGenStruct(name)
	g.Add(name, i)
	return g.Generate()
}

// NewGenStruct
func NewGenStruct(name string) *Gen {
	return &Gen{
		Name:  name,
		Types: map[string][]byte{},
	}
}

type Gen struct {
	Name  string
	Types map[string][]byte
}

// Generate returns struct source code
func (g *Gen) Generate() []byte {
	named := make([]string, 0, len(g.Types))
	for name := range g.Types {
		named = append(named, name)
	}
	sort.Strings(named)
	buf := bytes.NewBuffer(nil)
	for _, name := range named {
		buf.Write(g.Types[name])
	}
	return buf.Bytes()
}

// AddByValue Add the struct from reflect.Value.
func (g *Gen) AddByValue(name string, val reflect.Value) {
	g.defineStruct(name, val)
}

// Add Add the struct from interface.
func (g *Gen) Add(name string, val interface{}) {
	g.AddByValue(name, reflect.ValueOf(val))
}

func (g *Gen) toStar(name string) string {
	if strings.HasPrefix(name, "*") {
		return name
	}
	if _, ok := g.Types[name]; ok {
		return "*" + name
	}
	return name
}

func (g *Gen) defineStruct(name string, val reflect.Value) string {
	switch kind := val.Kind(); kind {
	case reflect.Float64:
		// Can't distinguish between integer and floating.
		return "json.Number"
	case reflect.Bool:
		return "bool"
	case reflect.String:
		v := val.String()
		// Identify the RFC3339Nano time format
		if _, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return "time.Time"
		}
		return "string"
	case reflect.Slice:
		if val.Len() == 0 {
			// No data can be identified in a slice.
			return fmt.Sprintf("[]%s", g.defineStruct(name, reflect.New(val.Type().Elem()).Elem()))
		}
		return fmt.Sprintf("[]%s", g.defineStruct(name, val.Index(0)))
	case reflect.Ptr, reflect.Interface:
		return g.toStar(g.defineStruct(name, val.Elem()))
	case reflect.Map:
		mk := val.MapKeys()
		if len(mk) == 0 {
			return "json.RawMessage"
		}
		if _, ok := g.Types[name]; ok {
			return name
		}
		valueSlice(mk).Sort()
		named := map[string]int{}
		g.Types[name] = nil
		buf := bytes.NewBuffer(nil)
		buf.WriteString(fmt.Sprintf("\n// %s This structure is generated from data.\ntype %s ", name, name))
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
			newTypeName = g.defineStruct(newTypeName, v)
			buf.WriteString(fmt.Sprintf("%s %s `json:\"%s,omitempty\"`\n", newName, newTypeName, name))
		}
		buf.WriteString("}\n")
		g.Types[name] = formatSrc(buf.Bytes())
		return g.toStar(name)

	// From other definitions.
	case
		reflect.Float32,
		reflect.Complex64, reflect.Complex128,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.UnsafePointer:
		return kind.String()
	case reflect.Array:
		return fmt.Sprintf("[%d]%s", val.Len(), g.defineStruct(name, val.Elem()))
	case reflect.Struct:
		if _, ok := g.Types[name]; ok {
			return name
		}
		buf := bytes.NewBuffer(nil)
		buf.WriteString(fmt.Sprintf("\n// %s This structure is generated from other definitions.\ntype %s ", name, name))
		buf.WriteString("struct {\n")
		typ := val.Type()
		num := typ.NumField()
		g.Types[name] = nil
		for i := 0; i != num; i++ {
			t := typ.Field(i)
			v := val.Field(i)
			buf.WriteString(fmt.Sprintf("%s %s %s\n", t.Name, g.defineStruct(t.Name, v), string(t.Tag)))
		}
		buf.WriteString("}\n")
		g.Types[name] = formatSrc(buf.Bytes())
		return g.toStar(name)
	default:
		// No action.
		return "json.RawMessage"
	}
}

func formatSrc(src []byte) []byte {
	newSrc, err := format.Source(src)
	if err != nil {
		return src
	}
	return newSrc
}
