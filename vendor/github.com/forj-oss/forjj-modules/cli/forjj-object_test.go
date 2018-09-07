package cli

import (
	"fmt"
	"forjj-modules/cli/kingpinMock"
	"reflect"
	"strings"
	"testing"
)

// -------------------------------
type ForjParamTester interface {
	GetFlag() *ForjFlag
	GetArg() *ForjArg
}

// -------------------------------
func flags_list(collection map[string]*ForjObjectListFlags) string {
	list := make([]string, 0, len(collection))
	for key := range collection {
		list = append(list, key)
	}
	return "'" + strings.Join(list, "' '") + "'"
}

var app = kingpinMock.New("Application")

const (
	create   = "create"
	update   = "update"
	maintain = "maintain"
)
const (
	create_help = "create-help"
	update_help = "update-help"
)

const (
	workspace  = "workspace"
	infra      = "infra"
	infra_help = "infra help"
)

const (
	w_f  = `[a-z]+[a-z0-9_-]*`
	ft_f = `[A-Za-z0-9_ !:/.-]+`
)

func TestForjCli_NewObject(t *testing.T) {
	t.Log("Expect NewObject('workspace', 'forjj workspace', true) to create a new object at App level.")

	const workspace_help = "workspace help"

	c := NewForjCli(app)
	o := c.NewObject(workspace, workspace_help, "internal")

	ot := reflect.TypeOf(o).String()
	if ot != "*cli.ForjObject" {
		t.Errorf("Expected to get ForjObject type. Got: %s", ot)
	}

	of, found := c.objects[workspace]
	if of != o {
		t.Error("Expected to get the object created registered. Is not.")
	}
	if !found {
		t.Errorf("Expected %s registered in the App layer as new object. Not found.", workspace)
	}
	if o.role != "internal" {
		t.Errorf("Expect to be an object with role '%s'. Got '%s'", "internal", o.role)
	}
	if o.name != workspace {
		t.Errorf("Expect object name to be '%s'. Got '%s'", workspace, o.name)
	}
	if o.desc != workspace_help {
		t.Errorf("Expect object help to be '%s'. Got %s", workspace_help, o.desc)
	}

	o = c.NewObject(workspace, workspace_help, "")
	if len(c.objects) != 1 {
		t.Errorf("Expect to have only one workspace object. Got %d", len(c.objects))
	}

	o = c.NewObject(infra, infra_help, "")
	if len(c.objects) != 2 {
		t.Errorf("Expect to have only one workspace object. Got %d", len(c.objects))
	}
}

func TestForjCli_GetObject(t *testing.T) {
	t.Log("Expect NewObject('workspace', 'forjj workspace', true) to create a new object at App level.")

	const workspace_help = "workspace help"

	c := NewForjCli(app)
	o := c.NewObject(workspace, workspace_help, "")

	o_found := c.GetObject(workspace)
	if o_found != o {
		t.Error("Expected any created object to be found and returned. Is not.")
	}
}

func TestForjObject_AddKey(t *testing.T) {
	t.Log("Expect ForjObject_AddKey() to add a new field key in the object.")

	// --- Setting test context ---
	const (
		docker      = "docker-exe-path"
		docker_help = "docker-exe-path-help"
	)
	c := NewForjCli(kingpinMock.New("Application"))
	o := c.NewObject(workspace, "", "")

	// --- Run the test ---
	or := o.AddKey(String, docker, docker_help, "", nil)

	// --- Start testing ---
	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}
}

func TestForjObject_AddField(t *testing.T) {
	t.Log("Expect AddField(cli.String, 'docker-exe-path', docker_exe_path_help) to add a field to workspace object.")

	const docker = "docker-exe-path"
	const docker_help = "docker-exe-path-help"
	const test = "test"

	c := NewForjCli(kingpinMock.New("Application"))
	o := c.NewObject(workspace, "", "").Single()

	or := o.AddField(String, docker, docker_help, "", nil)
	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}

	f, found := o.fields[docker]
	if !found {
		t.Errorf("Expected %s registered in the object as new field. Not found.", docker)
	}
	if f.name != docker {
		t.Errorf("Expect field name to be '%s'. Got %s", docker, f.name)
	}
	if f.help != docker_help {
		t.Errorf("Expect field help to be '%s'. Got %s", docker_help, f.help)
	}
	if f.value_type != String {
		t.Errorf("Expect field type to be '%s'. Got %s", String, f.value_type)
	}

	if  vl := len(o.cli.values) ; vl != 1 {
		t.Errorf("Expected to find the one object in object data. Found %d", vl)
	} else if v, found := o.cli.values[workspace] ; !found {
		t.Errorf("Expected to find the object '%s' in object data. Not found", workspace)
	} else if r, found2 := v.records[workspace] ; !found2 {
		t.Errorf("Expected to find object instance data '%s'. But not found.", workspace)
	} else if a, found3 := r.attrs[docker] ; !found3 {
		t.Errorf("Expected to find object instance attribute '%s'. But not found", docker)
	} else if a != nil {
		t.Errorf("Expected to find object instance attribute '%s' set to '%s'. Is not.", docker, "nil")
	}


	or = o.AddField(String, docker, "blabla", "", nil)
	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}
	if len(o.fields) > 2 {
		t.Errorf("Expected to have 2 fields with at least field '%s'. Got %d fields.", docker, len(o.fields))
	}

	f, found = o.fields[docker]
	if f.help != docker_help {
		t.Errorf("Expect field help to stay at '%s'. Got %s", docker_help, f.help)
	}

	// --------------- New Context

	opts := Opts()
	opts.Default("test")


	// --------------- Running
	or = o.AddField(String, test, "blabla", "", opts)

	// --------------- Testing
	if v, found := o.cli.values[workspace] ; !found {
		t.Errorf("Expected to find the object '%s' in object data.", workspace)
	} else if r, found2 := v.records[workspace] ; !found2 {
		t.Errorf("Expected to find object instance data '%s'. But not found.", workspace)
	} else if a, found3 := r.attrs[test] ; !found3 {
		t.Errorf("Expected to find object instance attribute '%s'. But not found", test)
	} else if a == nil {
		t.Errorf("Expected to find object instance attribute '%s' set to '%s'. Got nil.", test, test)
	} else if av, ok := a.(string) ; ok && av != test {
		t.Errorf("Expected to find object instance attribute '%s' set to '%s'. Got '%s'.", test, test, av)
	}
}

func TestForjObject_NoFields(t *testing.T) {
	t.Log("Expect ForjObject_NoFields() to register the object with no fields.")

	// --- Setting test context ---
	c := NewForjCli(kingpinMock.New("Application"))
	o := c.NewObject(workspace, "", "")

	// --- Run the test ---
	o = o.NoFields()

	// --- Start testing ---
	if o == nil {
		t.Error("Expected NoFields() to fails. but it works.")
	}
	if v, found := o.fields[no_fields]; !found {
		t.Error("Expected NoFields() to create 'no_field' record. Got nothing.")
	} else {
		if !v.key {
			t.Error("Expected NoFields() to create 'no_field' record as key. Is is not")
		}
	}

	// --- Setting test context ---
	c = NewForjCli(kingpinMock.New("Application"))
	o = c.NewObject(workspace, "", "").AddKey(String, "test", "help", "", nil)

	// --- Run the test ---
	o = o.NoFields()

	// --- Start testing ---
	if o != nil {
		t.Errorf("Expected NoFields() to work. But it fails. %s", c.GetObject(workspace).Error())
	}

	// --- Setting test context ---
	c = NewForjCli(kingpinMock.New("Application"))
	o = c.NewObject(workspace, "", "")

	// --- Run the test ---
	o = o.NoFields().AddKey(String, "test", "help", "", nil)

	// --- Start testing ---
	if o != nil {
		t.Errorf("Expected NoFields() to work. But it fails. %s", c.GetObject(workspace).Error())
	}
}

func TestForjObject_DefineActions(t *testing.T) {
	t.Log("Expect DefineActions('create') adding an action to fail if no action gets created from app.")

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject(workspace, "", "")
	or := o.DefineActions(create)
	if or != nil {
		t.Error("Expected DefineActions() to fail. Got one.")
	}

	o.AddKey(String, "test", "test help", "", nil)
	or = o.DefineActions(create)
	if or != o {
		t.Error("Expected to get the object 'object' updated. Is not.")
	}

	_, found := o.actions[create]
	if found {
		t.Errorf("Expected %s registered in the object as inexistent. Found it.", create)
	}
}

func TestForjObject_DefineActions2(t *testing.T) {
	t.Log("Expect actions to be added to the object.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	o := c.NewObject(workspace, "", "").AddKey(String, "test", "test help", "", nil)
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	c.NewActions(create, create_help, "create %s", true)
	if o.DefineActions(create) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	f, found := o.actions[create]
	if !found {
		t.Errorf("Expected %s registered in the object actions. Not found.", create)
	}
	if f.action == nil {
		t.Errorf("Expected action to refer to global action '%s'. Got nil", create)
	}
	if f.action.name != create {
		t.Errorf("Expected action name to refer to global action '%s'. Got %s", create, f.action.name)
	}

	cmd := app.GetCommand(create, workspace)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to be created in kingpin. Not found.", workspace)
	}
	if cmd.FullCommand() != workspace {
		t.Errorf("Expected Command to be '%s' in kingpin. Got '%s'", workspace, cmd.FullCommand())
	}

}

func TestForjObject_DefineActions3(t *testing.T) {
	t.Log("Expect double actions to be added to the object.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", false)
	o := c.NewObject(workspace, "", "").AddKey(String, "test", "test help", "", nil).
		DefineActions(create, update)

	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	// Check in cli
	cmd := app.GetCommand(create, workspace)
	cmd_in_cli := o.actions[create].cmd
	if cmd_in_cli == nil {
		t.Errorf("Expected Command '%s' to be found in kingpin. Is nil.", update)
	}
	if cmd != cmd_in_cli {
		t.Errorf("Expected Command '%s' to be found identical in kingpin. Is not.", create)
	}
	if c.objects[workspace].actions[create].cmd.FullCommand() != workspace {
		t.Errorf("Expected Command '%s' associated to object '%s' to be named '%s'. Got '%s'",
			create, workspace, workspace, c.objects[workspace].actions[create].action.cmd.FullCommand())
	}

	// Check in kingpin
	cmd = app.GetCommand(create)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to exist. Not found.", create)
	}
	if cmd.FullCommand() != create {
		t.Errorf("Expected '%s' has an command named '%s'", create, create)
	}

	cmd = app.GetCommand(create, workspace)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to be created under '%s'. Not found.", workspace, create)
	}
	if cmd.FullCommand() != workspace {
		t.Errorf("Expected '%s/%s' has an command named '%s'", create, workspace, workspace)
	}

	cmd = app.GetCommand(update)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to exist. Not found.", update)
	}
	if cmd.FullCommand() != update {
		t.Errorf("Expected '%s' has an command named '%s'", update, update)
	}

	cmd = app.GetCommand(update, workspace)
	if cmd == nil {
		t.Errorf("Expected Command '%s' to be created under '%s'. Not found.", workspace, update)
	}
	if cmd.FullCommand() != workspace {
		t.Errorf("Expected '%s/%s' has an command named '%s'", update, workspace, workspace)
	}
}

func TestForjObject_OnActions(t *testing.T) {
	t.Log("Expect OnAction() to select wanted action.")
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.NewActions(maintain, "", "maintain %s", false)
	o := c.NewObject(workspace, "", "").AddKey(String, "test", "test help", "", nil).
		DefineActions(create, update).
		OnActions(create)

	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}
	if len(o.actions) != 2 {
		t.Errorf("Expected 2 actions in object '%s'. Got '%d'", workspace, len(o.actions))
	}
	if len(o.sel_actions) != 1 {
		t.Errorf("Expected 1 selected action. Got '%d'", len(o.sel_actions))
	}

	a, found := o.sel_actions[create]
	if !found {
		t.Errorf("expected '%s' selected. Got nothing", create)
	}
	if a.action.name != create {
		t.Errorf("expected '%s' selected. Got '%s'", create, a.action.name)
	}

	o.OnActions(update)
	a, found = o.sel_actions[update]
	if !found {
		t.Errorf("expected '%s' selected. Got nothing", update)
	}
	if a.action.name != update {
		t.Errorf("expected '%s' selected. Got '%s'", update, a.action.name)
	}

	o.OnActions()
	if len(o.sel_actions) != 2 {
		t.Errorf("Expected 2 selected action. Got '%d'", len(o.sel_actions))
	}
}

func TestForjObject_AddFlag(t *testing.T) {
	t.Log("Expect AddFlag() to be added to selected actions.")

	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.NewActions(maintain, "", "maintain %s", false)

	const Path = "path"

	o := c.NewObject(workspace, "", "").
		AddKey(String, Path, "path help", "", nil).
		DefineActions(create, update).
		OnActions(create)
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	// --- Run the test ---
	or := o.AddFlag(Path, nil)
	// --- Start testing ---
	// check in cli.
	if or != o {
		t.Error("Expected to get the object updated. Is not.")
	}
	if of, found := o.actions[create].params[Path]; !found {
		t.Errorf("Expected to get parameter '%s' added to object action '%s'", Path, o.actions[create].name)
	} else {
		if of.(ForjParamTester).GetFlag().name != Path {
			t.Errorf("Expected Flag to be named '%s'", Path)
		}

		// check in kingpin
		f := app.GetFlag(create, workspace, Path)
		if f == nil {
			t.Errorf("Expected flag '%s' to be added to kingpin '%s' command. Got '%s'.",
				Path, workspace, app.ListOf(create, workspace))
			return
		}
		if f.GetName() != Path {
			t.Errorf("Expected flag name to be '%s'. Got '%s'", Path, f.GetName())
		}
		if of.(ForjParamTester).GetFlag().flag != f {
			t.Errorf("Expected kingpin flag '%s' to be stored in object action '%s'. Is not.",
				Path, o.actions[create].name)
		}
	}
}

func TestForjObject_ParseHook(t *testing.T) {
	t.Log("Expect ForjObject_ParseHook() to store the func provided.")

	const workspace_help = "workspace help"

	// --- Setting test context ---
	c := NewForjCli(app)

	var o *ForjObject
	// --- Run the test ---
	o_ret := o.ParseHook(func(_ *ForjObject, _ *ForjCli, _ interface{}) (error, bool) {
		return fmt.Errorf("This function is OK."), false
	})

	// --- Start testing ---
	if o_ret != nil {
		t.Error("Expected ParseHook() to return nil. But got one.")
	}

	// --- Setting test context ---
	o = c.NewObject(workspace, workspace_help, "")

	// --- Run the test ---
	o_ret = o.ParseHook(func(_ *ForjObject, _ *ForjCli, _ interface{}) (error, bool) {
		return fmt.Errorf("This function is OK."), false
	})

	// --- Start testing ---
	if o != o_ret {
		t.Error("Expected ParseHook() to return the object updated. Is not.")
	}
	if o.context_hook == nil {
		t.Error("Expected to have a hook stored. Got nil.")
	}
	if err, _ := o.context_hook(nil, nil, nil); fmt.Sprintf("%s", err) != "This function is OK." {
		t.Errorf("Expected to get the function stored to return what we want. Got '%s'", err)
	}
}

func TestForjObject_AddFlagsFromObjectAction(t *testing.T) {
	t.Log("Expect AddFlagsFromObjectAction() to be added to selected actions.")

	// --- Set context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)

	if o := c.NewObject(workspace, "", "").
		AddKey(String, "test", "test help", "", nil).
		DefineActions(update).
		OnActions(update).
		AddFlag("test", nil); o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	const (
		test        = "test"
		test2       = "test2"
		test_help   = "test help"
		another_obj = "another-obj"
	)

	infra_obj := c.NewObject(infra, "", "").
		AddKey(String, test2, test_help, "", nil).
		DefineActions(update).
		OnActions().
		AddFlag(test2, nil)

	// --- Running the test ---
	o := infra_obj.AddFlagsFromObjectAction(workspace, update)

	// --- Start Testing ---
	if o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}
	if o != infra_obj {
		t.Error("Expected to get the object updated. Is not.")
	}

	// Checking in cli
	expected_name := test
	param, found := o.actions[update].params[expected_name]
	if !found {
		t.Errorf("Expected flag '%s' added as in object action.params", test)
		return
	}

	f_cli := param.(forjParam).GetFlag()
	if f_cli == nil {
		t.Errorf("Expected to get a Flag from the object action '%s-%s'. Not found or is not a flag.",
			workspace, update)
	}

	// Checking in kingpin
	f := app.GetFlag(update, infra, test)
	if f == nil {
		t.Error("Expected to get flags from workspace added to another object action. Not found.")
		return
	}

	if f.GetName() != test {
		t.Errorf("Expected to get '%s' as flag name. Got '%s'", expected_name, f.GetName())
	}

	// Update context
	c.NewObject(another_obj, "", "").NoFields().DefineActions(update).OnActions()

	// Run test
	o = c.GetObject(another_obj).AddFlagsFromObjectAction(infra, update)

	// Start testing
	if o == nil {
		t.Errorf("Expected AddFlagsFromObjectAction() to NOT fail. %s", c.GetObject(another_obj).Error())
		return
	}

	// Checking in cli
	param, found = o.actions[update].params[test]
	if found {
		t.Errorf("Expected flag '%s' NOT added as in object action.params", test)
	}

	param, found = o.actions[update].params[test2]
	if !found {
		t.Errorf("Expected flag '%s' added as in object action.params", test2)
	}

	f_cli = param.(forjParam).GetFlag()
	if f_cli == nil {
		t.Errorf("Expected to get a Flag from the object action '%s-%s'. Not found or is not a flag.",
			workspace, update)
	}

	// Checking in kingpin
	f = app.GetFlag(update, another_obj, test2)
	if f == nil {
		t.Error("Expected to get flags from workspace added to another object action. Not found.")
		return
	}

	if f.GetName() != test2 {
		t.Errorf("Expected to get '%s' as flag name. Got '%s'", test2, f.GetName())
	}
}

func TestForjObject_AddFlagsFromObjectListActions(t *testing.T) {
	t.Log("Expect AddFlagFromObjectListActions() to be added to an object action as Flag.")

	const test = "test"

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if o := c.NewObject(workspace, "", "").
		AddKey(String, test, "test help", "", nil).
		DefineActions(update).
		OnActions(update).
		AddFlag(test, nil).
		CreateList("to_create", ",", "test", "").
		AddActions(update); o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	infra_obj := c.NewObject(infra, "", "").NoFields().
		DefineActions(update).
		OnActions()

	if infra_obj == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	// Checking in cli
	o := infra_obj.AddFlagsFromObjectListActions(workspace, "to_create", update)
	if o != infra_obj {
		t.Error("Expected to get the object updated. Is not.")
	}

	expected_name := update + "-" + workspace + "s"
	if _, found := c.objects[infra].actions[update].params[expected_name]; !found {
		t.Errorf("Expected to get a new Flag '%s'related to the Objectlist added. Not found.", expected_name)
	}

	// Checking flags list ref
	expected_ref_name := update + " " + infra_obj.name + " --" + expected_name
	if _, found := c.objects[workspace].list["to_create"].flags_list[expected_ref_name]; !found {
		t.Errorf("Expected to get a reference to the created flag '%s'. has %s.",
			expected_ref_name, flags_list(c.objects[workspace].list["to_create"].flags_list))
		return
	}

	fl := c.objects[workspace].list["to_create"].flags_list[expected_ref_name]
	if fl == nil {
		t.Errorf("Expected to have a reference created '%s'. Got nil.", expected_ref_name)
	}
	if fl.objList == nil {
		t.Error("Expected to reference to the list. Got nil")
		return
	}
	if fl.objList.name != "to_create" || fl.objList.obj.name != workspace {
		t.Errorf("Expected to reference to the list '%s %s'. Got ref to '%s %s'",
			workspace, "to_create", fl.objList.obj.name, fl.objList.name)
	}
	if !fl.multi_actions {
		t.Error("Expected to get multiple actions ref. Got single.")
	}
	if fl.objectAction == nil {
		t.Errorf("Expected to reference to the updated object action '%s'. Got nil", create)
		return
	}
	if fl.objectAction.name != update+"_"+infra {
		t.Errorf("Expected to reference to the updated action '%s'. Got '%s'", update+"_"+infra, fl.objectAction.name)
	}

	// Checking in kingpin
	flag := app.GetFlag(update, infra, expected_name)
	if flag == nil {
		t.Errorf("Expected to get a Flag in kingpin called '%s'. Got '%s'",
			update+"-"+workspace+"s", app.ListOf(update, infra))
	}
}

func TestForjObject_AddFlagFromObjectListAction(t *testing.T) {
	t.Log("Expect AddFlagFromObjectListActions() to be added to an object action as Flag.")

	// --- Setting test context ---
	const test = "test"

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if o := c.NewObject(workspace, "", "").
		AddKey(String, test, "test help", "", nil).
		DefineActions(update).
		OnActions(update).
		AddFlag(test, nil).
		CreateList("to_create", ",", "test", "").
		AddActions(update); o == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	infra_obj := c.NewObject(infra, "", "").NoFields().
		DefineActions(update).
		OnActions()

	// --- Run the test ---
	o := infra_obj.AddFlagFromObjectListAction(workspace, "to_create", update)

	// --- Start testing ---
	if o == nil {
		t.Error("Expected AddFlagFromObjectListAction() to return the object updated. Got nil")
		return
	}
	if o != infra_obj {
		t.Error("Expected to get the object updated. Is not.")
	}

	// Checking in cli
	expected_name := workspace + "s"
	if _, found := c.objects[infra].actions[update].params[expected_name]; !found {
		t.Errorf("Expected to get a new Flag '%s'related to the Objectlist added. Not found.", expected_name)
	}

	// Checking flags list ref
	expected_ref_name := infra_obj.name + " --" + expected_name
	if _, found := c.objects[workspace].list["to_create"].flags_list[expected_ref_name]; !found {
		t.Errorf("Expected to get a reference to the created flag '%s'. has %s.",
			expected_ref_name, flags_list(c.objects[workspace].list["to_create"].flags_list))
		return
	}

	fl := c.objects[workspace].list["to_create"].flags_list[expected_ref_name]
	if fl == nil {
		t.Errorf("Expected to have a reference created '%s'. Got nil", expected_ref_name)
	}
	if fl.objList == nil {
		t.Error("Expected to reference to the list. Got nil")
		return
	}
	if fl.objList.name != "to_create" || fl.objList.obj.name != workspace {
		t.Errorf("Expected to reference to the list '%s %s'. Got ref to '%s %s'",
			workspace, "to_create", fl.objList.obj.name, fl.objList.name)
	}
	if fl.multi_actions {
		t.Error("Expected to get single action ref. Got multiple.")
	}
	if fl.objectAction == nil {
		t.Errorf("Expected to reference to the updated object action '%s'. Got nil", create)
		return
	}
	if fl.objectAction.name != update+"_"+infra {
		t.Errorf("Expected to reference to the updated object action '%s'. Got '%s'", update+"_"+infra, fl.objectAction.name)
	}

	// Checking in kingpin
	flag := app.GetFlag(update, infra, expected_name)
	if flag == nil {
		t.Errorf("Expected to get a Flag in kingpin called '%s'. Got '%s'",
			update+"-"+workspace+"s", app.ListOf(update, infra))
	}

}

func TestForjObject_SetParamOptions(t *testing.T) {
	t.Log("Expect ForjObject_SetParamOptions() to update existing flags anywhere we have the flag set.")

	const (
		test             = "test"
		test_help        = "test help"
		key              = "key"
		key_help         = "key help"
		key_value        = "key-value"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value       = "flag value"
		myapp            = "app"
		apps             = "apps"
		app_help         = "app help"
		instance         = "instance"
		instance_help    = "instance help"
		driver           = "driver"
		driver_help      = "driver help"
		driver_type      = "driver_type"
		driver_type_help = "driver_type help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag2_value      = "flag2 value"
		myinstance       = "myapp"
	)
	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if c.NewObject(test, test_help, "").
		AddKey(String, key, key_help, "#w", nil).
		AddField(String, flag, flag_help, "#w", nil).
		DefineActions(update).OnActions().
		AddFlag(key, Opts().Required()).
		AddFlag(flag, nil) == nil {
		t.Error(c.GetObject(test).Error())
	}

	if c.NewObject(myapp, app_help, "").
		AddKey(String, instance, instance_help, "#w", nil).
		AddField(String, driver, driver_help, "#w", nil).
		AddField(String, driver_type, driver_type_help, "#w", nil).
		AddField(String, flag2, flag2_help, "#w", nil).
		ParseHook(func(_ *ForjObject, c *ForjCli, _ interface{}) (err error, updated bool) {
		ret, found, _, err := c.GetStringValue(myapp, myinstance, flag2)
		if found {
			t.Error("Expected GetStringValue() to NOT find the context value. Got one.")
		}
		if ret != "" {
			t.Errorf("Expected GetStringValue() to return '' from context. Got '%s'", ret)
		}

		ret, found, _, err = c.GetStringValue(test, key_value, flag)
		if !found {
			t.Errorf("Expected GetStringValue() to find the context value. Got none. %s", err)
		}
		if ret != flag_value {
			t.Errorf("Expected GetStringValue() to return '%s' from context. Got '%s'", flag_value, ret)
		}

		ret, found, _, err = c.GetStringValue(test, "", flag)
		if !found {
			t.Errorf("Expected GetStringValue() to find the context value. Got none. %s", err)
		}
		if ret != flag_value {
			t.Errorf("Expected GetStringValue() to return '%s' from context. Got '%s'", flag_value, ret)
		}
		return nil, false
	}).
		DefineActions(create).OnActions().
		AddFlag(driver_type, nil).
		AddFlag(driver, nil).
		AddFlag(instance, Opts().Required()).
		AddFlag(flag2, nil).
		CreateList("to_create", ",", "driver_type:driver[:instance]", app_help).
		AddValidateHandler(func(l *ForjListData) (err error) {
		if v, found := l.Data[instance]; !found || v == "" {
			l.Data[instance] = l.Data[driver]
		}
		return nil
	}) == nil {
		t.Error(c.GetObject(myapp).Error())
	}

	c.GetObject(test).AddFlagFromObjectListAction(myapp, "to_create", create)

	context := []string{"cmd:" + update, "cmd:" + test, key, key_value, flag, flag_value,
		apps, "mytype:mydriver", "mydriver-flag2", flag2_value}

	if _, err := c.Parse(context, nil); err != nil {
		t.Errorf("Expected Parse() to work successfully. Got '%s'", err)
	}

	// --- Run the test ---
	c.GetObject(myapp).SetParamOptions(flag2, Opts().Default("myDefaultDriver"))

	// --- Start testing ---
	// Testing in kingpin
	f := app.GetFlag(create, myapp, flag2)
	if !f.IsDefault("myDefaultDriver") {
		t.Errorf("'%s %s %s' Flag default is not to '%s'", create, myapp, flag2, "myDefaultDriver")
	}

	f = app.GetFlag(update, test, "mydriver-flag2")
	if f == nil {
		t.Error("Expected kingpin Flag to exit. Not found.")
		return
	}
	if !f.IsDefault("myDefaultDriver") {
		t.Errorf("'%s %s %s' Flag default is not to '%s'", update, test, "mydriver-flag2", "myDefaultDriver")
	}
}

func TestForjObject_HasField(t *testing.T) {
	t.Log("Expect ForjObject_HasField() to return the existence of a field.")

	// --- Setting test context ---
	const (
		test             = "test"
		test_help        = "test help"
		key              = "key"
		key_help         = "key help"
		key_value        = "key-value"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value       = "flag value"
		myapp            = "app"
		apps             = "apps"
		app_help         = "app help"
		instance         = "instance"
		instance_help    = "instance help"
		driver           = "driver"
		driver_help      = "driver help"
		driver_type      = "driver_type"
		driver_type_help = "driver_type help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag2_value      = "flag2 value"
		myinstance       = "myapp"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if c.NewObject(test, test_help, "").
		AddKey(String, key, key_help, "#w", nil).
		AddField(String, flag, flag_help, "#w", nil) == nil {
		t.Error(c.GetObject(test).Error())
	}

	// --- Run the test ---
	res1 := c.GetObject(test).HasField(flag)
	res2 := c.GetObject(test).HasField("blabla")
	// --- Start testing ---
	if !res1 {
		t.Errorf("Expected flag '%s' to exist. HasField said 'not found'.", flag)
	}
	if res2 {
		t.Errorf("Expected flag '%s' to NOT exist. HasField said 'found it'.", "blabla")
	}
}

func TestForjObject_Single(t *testing.T) {
	t.Log("Expect ForjObject_Single() to set Single record mode.")

	// --- Setting test context ---
	const (
		test             = "test"
		test_help        = "test help"
		key              = "key"
		key_help         = "key help"
		key_value        = "key-value"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value       = "flag value"
		myapp            = "app"
		apps             = "apps"
		app_help         = "app help"
		instance         = "instance"
		instance_help    = "instance help"
		driver           = "driver"
		driver_help      = "driver help"
		driver_type      = "driver_type"
		driver_type_help = "driver_type help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag2_value      = "flag2 value"
		myinstance       = "myapp"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if c.NewObject(test, test_help, "") == nil {
		t.Error(c.GetObject(test).Error())
	}

	// --- Run the test ---
	o := c.GetObject(test).Single()
	// --- Start testing ---
	if o == nil {
		t.Errorf("Expected object to be declared as single. Got nil. %s", c.GetObject(test).Error())
	}
	if !o.single {
		t.Error("Expected object to be single. But got false.")
	}
	if field, found := o.fields[test + ".key"] ; !found {
		t.Errorf("Expected single key '%s.key' to exist. But not found.", test)
	} else {
		if ! field.key {
			t.Errorf("Expected single key '%s.key' to be a key. But is not.", test)
		}
	}

	if v, found := o.cli.values[test] ; !found {
		t.Errorf("Expected to find the object '%s' in object data.", test)
	} else if r, found2 := v.records[test] ; !found2 {
		t.Errorf("Expected to find object instance data '%s'. But not found.", test)
	} else if a, found3 := r.attrs["action"] ; !found3 {
		t.Errorf("Expected to find object instance attribute '%s'. But not found", "action")
	} else if av, ok := a.(string) ; ok && av != "setup" {
		t.Errorf("Expected to find object instance attribute '%s' set to '%s'. Got '%s'", "action", "setup", av)
	}

	if o.AddKey(String, key, key_help, "#w", nil) != nil {
		t.Error("Expected key setting to fails. But can create a key.")
	}
	if o.AddField(String, key, flag_help, "#w", nil) == nil {
		t.Error("Expected adding a field without issue. It fails.", c.GetObject(test).Error())
	}
}

func TestForjObject_SingleErrors(t *testing.T) {
	t.Log("Expect ForjObject_Single() to set Single record mode.")

	// --- Setting test context ---
	const (
		test = "test"
		test_help = "test help"
		key = "key"
		key_help = "key help"
		key_value = "key-value"
		flag = "flag"
		flag_help = "flag help"
		flag_value = "flag value"
		myapp = "app"
		apps = "apps"
		app_help = "app help"
		instance = "instance"
		instance_help = "instance help"
		driver = "driver"
		driver_help = "driver help"
		driver_type = "driver_type"
		driver_type_help = "driver_type help"
		flag2 = "flag2"
		flag2_help = "flag2 help"
		flag2_value = "flag2 value"
		myinstance = "myapp"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, "", "update %s", false)
	c.AddFieldListCapture("w", w_f)

	if c.NewObject(test, test_help, "").AddField(String, key, flag_help, "#w", nil) == nil {
		t.Errorf("Expected context fails. %s", c.GetObject(test).Error())
	}

	// --- Run the test ---
	o := c.GetObject(test).Single()
	// --- Start testing ---
	if o != nil {
		t.Error("Expected single object setup to fail. But it succeeded.")
	}
}

func TestForjObject_AddInstance_one(t *testing.T) {
		t.Log("Expect ForjObject_AddInstance() to add a new object instance.")

	// --- Setting test context ---
	const (
		test             = "test"
		test_help        = "test help"
		key              = "key"
		key_help         = "key help"
		key_value        = "key-value"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value       = "flag value"
		myapp            = "app"
		apps             = "apps"
		app_help         = "app help"
		instance         = "instance"
		instance_help    = "instance help"
		driver           = "driver"
		driver_help      = "driver help"
		driver_type      = "driver_type"
		driver_type_help = "driver_type help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag2_value      = "flag2 value"
		myinstance       = "myapp"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	obj := c.NewObject(test, test_help, "")
	if  obj == nil {
		t.Error(c.GetObject(test).Error())
		return
	}

	// --- Run the test ---
	res1 := obj.AddInstances(instance)

	// --- Start testing ---
	if res1 != obj {
		t.Error("expected object return to be the object used. But different.")
	}
	if v, found := res1.instances[instance] ; ! found {
		t.Error("Expected instance '%s' to exist. But not found.", instance)
	} else {
		if v.name != instance {
			t.Error("Expected instance '%s' to exist as '%s'. But found it as '%s'.", instance, instance, v.name)
		}
	}
}

func TestForjObject_AddInstances_more(t *testing.T) {
	t.Log("Expect ForjObject_AddInstance() to add a new object instance.")

	// --- Setting test context ---
	const (
		test             = "test"
		test_help        = "test help"
		key              = "key"
		key_help         = "key help"
		key_value        = "key-value"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value       = "flag value"
		myapp            = "app"
		apps             = "apps"
		app_help         = "app help"
		instance         = "instance"
		instance_help    = "instance help"
		driver           = "driver"
		driver_help      = "driver help"
		driver_type      = "driver_type"
		driver_type_help = "driver_type help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag2_value      = "flag2 value"
		myinstance       = "myapp"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	obj := c.NewObject(test, test_help, "")
	if  obj == nil {
		t.Error(c.GetObject(test).Error())
		return
	}

	if l := len(obj.instances) ; l != 0 {
		t.Error("expected object instances to be empty. Found %d instances.", l)
	}

	// --- Run the test ---
	res1 := obj.AddInstances()

	// --- Start testing ---
	if res1 != obj {
		t.Error("expected object return to be the object used. But different.")
	}

	if l := len(obj.instances) ; l != 0 {
		t.Error("expected object instances to be empty. Found %d instances.", l)
	}

	// --- Run the test ---
	res1 = obj.AddInstances(instance, myinstance)

	// --- Start testing ---
	if res1 != obj {
		t.Error("expected object return to be the object used. But different.")
	}

	if l := len(obj.instances) ; l != 2 {
		t.Error("expected object instances to 2 items. Found %d instances.", l)
	}

	if v, found := res1.instances[instance] ; ! found {
		t.Error("Expected instance '%s' to exist. But not found.", instance)
	} else {
		if v.name != instance {
			t.Error("Expected instance '%s' to exist as '%s'. But found it as '%s'.", instance, instance, v.name)
		}
	}

	if v, found := res1.instances[myinstance] ; ! found {
		t.Error("Expected instance '%s' to exist. But not found.", myinstance)
	} else {
		if v.name != myinstance {
			t.Error("Expected instance '%s' to exist as '%s'. But found it as '%s'.", myinstance, myinstance, v.name)
		}
	}

	// --- Run the test ---
	res1 = obj.AddInstances(myinstance)

	// --- Start testing ---
	if l := len(obj.instances) ; l != 2 {
		t.Error("expected object instances to 2 items. Found %d instances.", l)
	}

	if v, found := res1.instances[instance] ; ! found {
		t.Error("Expected instance '%s' to exist. But not found.", instance)
	} else {
		if v.name != instance {
			t.Error("Expected instance '%s' to exist as '%s'. But found it as '%s'.", instance, instance, v.name)
		}
	}

	if v, found := res1.instances[myinstance] ; ! found {
		t.Error("Expected instance '%s' to exist. But not found.", myinstance)
	} else {
		if v.name != myinstance {
			t.Error("Expected instance '%s' to exist as '%s'. But found it as '%s'.", myinstance, myinstance, v.name)
		}
	}
}
