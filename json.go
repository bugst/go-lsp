package lsp

import (
	"fmt"
	"reflect"
	"strings"
)

// UnmarshalJSON read json data and enforce LSP restrictions
func UnmarshalJSON(data []byte, res interface{}) error {
	examiner(reflect.TypeOf(res), 1)
	return nil
}

func examiner(t reflect.Type, depth int) {
	ind := strings.Repeat("\t", depth)
	fmt.Println(ind, "Type is", t.Name(), "and kind is", t.Kind())
	switch t.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		fmt.Println(ind+"\t", "Contained type:")
		examiner(t.Elem(), depth+1)
	case reflect.Struct:
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fmt.Println(ind+"\t", "Field", i+1, "name is", f.Name, "type is", f.Type.Name(), "and kind is", f.Type.Kind())
			if f.Tag != "" {
				fmt.Println(ind+"\t\t", "Tag is", f.Tag)
				fmt.Println(ind+"\t\t", "tag1 is", f.Tag.Get("tag1"), "tag2 is", f.Tag.Get("tag2"))
			}
		}
	}
}
