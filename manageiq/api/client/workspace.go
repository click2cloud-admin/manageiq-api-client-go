package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

type attributes map[string]interface{}
type ae_objects map[string]attributes

type Workspace struct {
	connection *ConnectionParameters_t
	Guid       string
	Input      struct {
		AeObjects        ae_objects `json:"workspace"`
		StateVars        attributes `json:"state_vars"`
		Current          attributes `json:"current"`
		MethodParameters attributes `json:"method_parameters"`
	}
	Output struct {
		Action   string `json:"action"`
		Resource struct {
			AeObjects attributes `json:"workspace"`
			StateVars attributes `json:"state_vars"`
		} `json:"resource"`
	}
}

func NewWorkspace(connection *ConnectionParameters_t, guid string) *Workspace {
	var workspace Workspace
	workspace.connection = connection
	workspace.Guid = guid
	workspace.Output.Action = "edit"
	workspace.Output.Resource.StateVars = make(map[string]interface{})
	workspace.Output.Resource.AeObjects = make(map[string]interface{})
	return &workspace
}

func (workspace *Workspace) Fetch() error {
	href_slug := "automate_workspaces/" + workspace.Guid
	fmt.Println("Href slug", href_slug)

	body, _ := workspace.connection.Get(href_slug)

	if err := json.Unmarshal(body, workspace); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (workspace *Workspace) Update() error {
	b, _ := json.Marshal(&workspace.Output)
	href_slug := "automate_workspaces/" + workspace.Guid
	b, err := workspace.connection.Post(href_slug, b)
	if err != nil {
		fmt.Println("Error updating workspace")
		return err
	}
	fmt.Println(string(b))
	return nil
}

func (workspace *Workspace) DumpObject(name string) {
	obj, err := workspace.GetObject(name)
	if err != nil {
		fmt.Println("Object ", name, "Not Found")
		return
	}
	fmt.Println("Attribute List", obj.GetAttributeList())
	for _, attr := range obj.GetAttributeList() {
		fmt.Println("Attribute Name", attr)
		obj.printValue(obj.GetAttribute(attr))
	}
	return
}

func (workspace *Workspace) Dump() {
	fmt.Println("Object List", workspace.GetObjectList())
	for _, object_name := range workspace.GetObjectList() {
		fmt.Println("Object:-Name-", object_name)
		workspace.DumpObject(object_name)
	}
	return
}

func (workspace *Workspace) GetOutputObject(object_name string) map[string]interface{} {
	_, err := workspace.Output.Resource.AeObjects[object_name]
	if !err {
		workspace.Output.Resource.AeObjects[object_name] = make(map[string]interface{})
	}
	return workspace.Output.Resource.AeObjects[object_name].(map[string]interface{})
}

func (workspace *Workspace) GetObject(object_name string) (*MiqAeObject, error) {
	var obj MiqAeObject
	attrs, err := workspace.Input.AeObjects[object_name]
	if !err {
		return nil, errors.New("Object " + object_name + "not found in workspace")
	}

	obj.workspace = workspace
	obj.name = object_name
	obj.attrs = attrs
	return &obj, nil
}

func (workspace *Workspace) GetStateVar(name string) interface{} {
	return workspace.Input.StateVars[name]
}

func (workspace *Workspace) SetStateVar(name string, value interface{}) {
	workspace.Input.StateVars[name] = value
	workspace.Output.Resource.StateVars[name] = value
}

func (workspace *Workspace) GetObjectList() []string {
	var keys []string

	for key, _ := range workspace.Input.AeObjects {
		keys = append(keys, key)
	}
	return keys
}

func (workspace *Workspace) StateVarExist(name string) bool {
	if _, ok := workspace.Input.StateVars[name]; ok {
		return true
	}
	return false
}

func (workspace *Workspace) GetCurrentObject() (*MiqAeObject, error) {
	object_name := "/" + workspace.Input.Current["namespace"].(string) + "/" +
		workspace.Input.Current["class"].(string) + "/" +
		workspace.Input.Current["instance"].(string)
	return workspace.GetObject(object_name)
}
