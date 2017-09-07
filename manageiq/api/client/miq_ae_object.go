package client

import (
	"fmt"
	"strings"
)

type MiqAeObject struct {
	workspace *Workspace
	name      string
	attrs     attributes
}

func (miq_object *MiqAeObject) printValue(value interface{}) {
	switch v := value.(type) {
	case string:
		fmt.Println(v)
	case int32, int64:
		fmt.Println(v)
	case float64:
		fmt.Println(v)
	case bool:
		fmt.Println(v)
	case VMDB_Object:
		fmt.Println(v.GetAttribute("href"))
	default:
		fmt.Println("unknown")
	}
	return
}

func (obj *MiqAeObject) GetAttribute(attribute_name string) interface{} {
	value := obj.attrs[attribute_name]
	original, ok := value.(string)
	if ok && strings.HasPrefix(original, "vmdb_reference::") {
		n := strings.Split(original, "::")[1]
		vmdb := NewVMDB_Object(obj.workspace.connection, n)
		vmdb.Fetch()
		return vmdb
	} else {
		return value
	}
}

func (obj *MiqAeObject) SetAttribute(attribute_name string, attribute_value interface{}) {
	obj.attrs[attribute_name] = attribute_value
	obj.workspace.GetOutputObject(obj.name)[attribute_name] = attribute_value
	return
}

func (obj *MiqAeObject) GetAttributeList() []string {
	var keys []string

	for key, _ := range obj.attrs {
		keys = append(keys, key)
	}
	return keys
}
