package gostruct

import (
	"fmt"
	"reflect"
	"sort"
)

type valueSlice []reflect.Value

func (p valueSlice) Len() int { return len(p) }
func (p valueSlice) Less(i, j int) bool {
	return fmt.Sprint(p[i].Interface()) < fmt.Sprint(p[j].Interface())
}
func (p valueSlice) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p valueSlice) Sort()         { sort.Sort(p) }
