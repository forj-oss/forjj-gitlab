package cli

import (
	"fmt"
	"forjj-modules/cli/kingpinMock"
	"reflect"
	"testing"
)

func check_object_exist(c *ForjCli, o_name, o_key, flag, value, atAction string, isDefault bool) error {
	if _, found := c.values[o_name]; !found {
		return fmt.Errorf("Expected object '%s' to exist in values. Not found.", o_name)
	}
	if _, found := c.values[o_name].records[o_key]; !found {
		return fmt.Errorf("Expected object '%s', record '%s' to exist in values. Not found.", o_name, o_key)
	}
	if v, found := c.values[o_name].records[o_key].attrs["action"]; !found {
		return fmt.Errorf("Expected object '%s', record '%s' to have the 'action' attribute. Not found.", o_name, o_key)
	} else {
		if d, ok := v.(string); !ok {
			return fmt.Errorf("Expected object '%s', record '%s' to have the 'action' attribute as string. Not found.", o_name, o_key)
		} else {
			if d != atAction {
				return fmt.Errorf("Expected object '%s', record '%s' to have the 'action' attribute value set to '%s'."+
					" Got '%s'.", o_name, o_key, atAction, d)
			}
		}
	}
	if v, found := c.values[o_name].records[o_key].attrs[flag]; !found {
		return fmt.Errorf("Expected record '%s-%s' to have '%s = %s' in values. Not found.",
			o_name, o_key, flag, value)
	} else {
		switch v.(type) {
		case *string:
			if !isDefault {
				return fmt.Errorf("Expected value to NOT come from a default value (*string). But got default value addr.")
			}
			if *v.(*string) != value {
				return fmt.Errorf("Expected key value '%s-%s-%s' to be set to '%s' (default). Got '%s'.",
					o_name, o_key, flag, value, *v.(*string))
			}
		case string:
			if isDefault {
				return fmt.Errorf("Expected value to come from a default value (*string). But got '%s'.", v)
			}
			if v.(string) != value {
				return fmt.Errorf("Expected key value '%s-%s-%s' to be set to '%s'. Got '%s'",
					o_name, o_key, flag, value, v)
			}
		}
	}
	return nil
}

func check_list_action_params(c *ForjCli, list, action, param_name string) (ForjParam, bool, error) {
	if _, found := c.list[list]; !found {
		return nil, false, fmt.Errorf("Expected list to have '%s' entry. Got none.", list)
	}
	if _, found := c.list[list].actions[action]; !found {
		return nil, false, fmt.Errorf("Expected list '%s'to have '%s' action entry. Got none.", list, action)
	}
	v, found := c.list[list].actions[action].params[param_name]
	return v, found, nil
}

func check_action_params(c *ForjCli, action, param_name string) (ForjParam, bool, error) {
	if _, found := c.actions[action]; !found {
		return nil, false, fmt.Errorf("Expected to have '%s' action. Got none.", action)
	}
	v, found := c.actions[action].params[param_name]
	return v, found, nil
}

func TestForjCli_loadContext(t *testing.T) {
	t.Log("Expect LoadContext() to report the context with values.")

	// --- Setting test context ---
	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "create %s", true)

	if c.NewObject(test, test_help, "").AddKey(String, flag, flag_help, "", nil).DefineActions(update) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	app.NewContext().SetContext(update, test).SetContextValue(flag, flag_value)
	// --- Run the test ---
	err := c.loadContext([]string{}, nil)

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected LoadContext() to not fail. Got '%s'", err)
	}
	if c.cur_cmds == nil {
		t.Error("Expected LoadContext() to return the last context command. Got none.")
	}
	if len(c.cur_cmds) != 2 {
		t.Errorf("Expected to have '%d' context commands. Got '%d'", 2, len(c.cur_cmds))
	}
}

func TestForjCli_identifyObjects(t *testing.T) {
	t.Log("Expect ForjCli_identifyObjects() to identify and store context reference to action, object and object list.")

	// --- Setting test context ---
	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "create %s", true)

	if c.NewObject(test, test_help, "").
		AddKey(String, flag, flag_help, "", nil).
		DefineActions(update).OnActions().
		AddFlag(flag, nil) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	context := app.NewContext()
	if context.SetContext(update, test) == nil {
		t.Error("Expected context with SetContext() to set context. But fails.")
	}
	if _, err := context.SetContextValue(flag, flag_value); err != nil {
		t.Errorf("Expected context with SetContextValue() to set values. But fails. %s", err)
	}

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) != 2 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}

	// --- Run the test ---
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Start testing ---
	if c.cli_context.action == nil {
		t.Error("Expected action to be identified. Got nil.")
		return
	}
	if c.cli_context.action != c.objects[test].actions[update].action {
		t.Errorf("Expected Action to be '%s'. Got '%s.", update, c.cli_context.action.name)
	}
	if c.cli_context.object == nil {
		t.Error("Expected object to be identified. Got nil.")
		return
	}
	if c.cli_context.object != c.objects[test] {
		t.Errorf("Expected Object to be '%s'. Got '%s.", test, c.cli_context.object.name)
	}
	if c.cli_context.list != nil {
		t.Errorf("Expected object to be nil. Got '%s'.", c.cli_context.list.name)
		return
	}

	// ------------------------------------------------------------------
	// --- Updating test context ---
	const (
		flag2       = "flag2"
		flag2_help  = "flag2 help"
		flag2_value = "flag2 value"
	)
	c.OnActions(create).AddFlag(String, flag2, flag2_help, nil)

	if ctxt, err := app.NewContext().SetContext(create).SetContextValue(flag2, flag2_value); err != nil {
		t.Errorf("Expect context to work. But fails. %s", err)
		return
	} else {
		context = ctxt
	}
	context, _ = app.NewContext().SetContext(create).SetContextValue(flag2, flag2_value)
	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds = context.SelectedCommands()
	if len(cmds) != 1 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}

	// --- Run the test ---
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Start testing ---
	if c.cli_context.action == nil {
		t.Error("Expected action to be identified. Got nil.")
		return
	}
	if c.cli_context.action != c.actions[create] {
		t.Errorf("Expected Action to be '%s'. Got '%s.", create, c.cli_context.action.name)
	}
	if c.cli_context.object != nil {
		t.Errorf("Expected object to be nil. Got '%s'.", c.cli_context.object.name)
		return
	}
	if c.cli_context.list != nil {
		t.Errorf("Expected object to be nil. Got '%s'.", c.cli_context.list.name)
		return
	}

	// ------------------------------------------------------------------
	// --- Updating test context ---
	const (
		repo               = "repo"
		repos              = "repos"
		reposlist_value    = "myinstance:myname,otherinstance"
		repo_help          = "repo help"
		reponame           = "name"
		reponame_help      = "repo name help"
		repo_instance      = "repo_instance"
		repo_instance_help = "repo instance help"
	)

	c.AddFieldListCapture("w", w_f)

	o := c.NewObject(repo, repo_help, "").
		AddKey(String, repo_instance, repo_instance_help, "#w", nil).
		AddField(String, reponame, reponame_help, "#w", nil).
		DefineActions(create).OnActions().
		AddFlag(repo_instance, nil).
		AddFlag(reponame, nil).
		CreateList("list", ",", "repo_instance[:name]", repo_help).
		AddActions(create)

	if o == nil {
		t.Errorf("Expected context failed to work with error:\n%s", c.GetObject(repo).Error())
		return
	}
	if ctxt, err := app.NewContext().SetContext(create, repos).SetContextValue(repos, reposlist_value); err != nil {
		t.Errorf("Expected context to work. But fails. %s", err)
		return
	} else {
		context = ctxt
	}
	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds = context.SelectedCommands()
	if len(cmds) != 2 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}

	// --- Run the test ---
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Start testing ---
	if c.cli_context.action == nil {
		t.Error("Expected action to be identified. Got nil.")
		return
	}
	if v, found := c.objects[repo].actions[create]; found {
		if c.cli_context.action != v.action {
			t.Errorf("Expected Action to be '%s'. Got '%s.", create, c.cli_context.action.name)
		}
	} else {
		t.Errorf("Expected Action '%s' to exist in Object '%s'. Got Nil.", create, repo)
	}

	if c.cli_context.object == nil {
		t.Error("Expected object to be set. Got Nil.")
		return
	}
	if c.cli_context.object != c.objects[repo] {
		t.Errorf("Expected Object to be '%s'. Got '%s.", repo, c.cli_context.object.name)
	}
	if c.cli_context.list == nil {
		t.Error("Expected object to be set. Got Nil.")
		return
	}
	if c.cli_context.list != c.objects[repo].list["list"] {
		t.Errorf("Expected Object to be '%s'. Got '%s.", repo, c.cli_context.object.name)
	}
}

// TestForjCli_loadListData_contextObject :
// check if <app> update test --flag "flag value"
// => creates an unique object 'test' record with key and data set.
func TestForjCli_loadListData_contextObject(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test       = "test"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag value"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, "").
		AddKey(String, flag, flag_help, "#w", nil).
		DefineActions(update).
		OnActions().
		AddFlag(flag, nil) == nil {
		t.Errorf("Expected Context Object declaration to work. %s", c.GetObject(workspace).Error())
		return
	}

	context, err := app.NewContext().SetContext(update, test).SetContextValue(flag, flag_value)
	if err != nil {
		t.Errorf("Expected context to work. But fails. %s", err)
		return
	}

	if ctxt, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	} else {
		c.cli_context.context = ctxt
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err = c.loadListData(nil, context)

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value, flag, flag_value, update, false); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextAction :
// check if <app> update --tests "flag_key"
// => creates an unique object 'test' record with key and data set.
func TestForjCli_loadListData_contextAction(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test       = "test"
		tests      = "tests"
		test_help  = "test help"
		flag       = "flag"
		flag_help  = "flag help"
		flag_value = "flag_key"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, "test object help", "").
		AddKey(String, flag, flag_help, "#w", nil).
		// <app> create test --flag <data>
		DefineActions(update).OnActions().
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "flag", "test object help").
		// <app> create tests "flag_key"
		AddActions(update) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> create --tests "flag_key"
	if c.OnActions(create).AddActionFlagFromObjectListAction(test, "to_update", update) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.Error())
	}

	if ctx, err := app.NewContext().SetContext(create).SetContextValue(tests, flag_value); err != nil {
		t.Errorf("Expected context to work. But fails. %s", err)
		return
	} else {
		c.cli_context.context = ctx
	}

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := c.cli_context.context.SelectedCommands()
	if len(cmds) != 1 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, c.cli_context.context)

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value, flag, flag_value, update, false); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextObjectList:
// check if <app> update tests "flag value,other"
// => creates 2 objects 'test' records with key and data set.
func TestForjCli_loadListData_contextObjectList(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test        = "test"
		tests       = "tests"
		test_help   = "test help"
		flag        = "flag"
		flag_help   = "flag help"
		flag_value1 = "flag_value"
		flag_value2 = "other"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	//	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, "").
		AddKey(String, flag, flag_help, "#w", nil).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "flag", test_help).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	context, err := app.NewContext().SetContext(create, tests).SetContextValue(tests, flag_value1+","+flag_value2)
	if err != nil {
		t.Errorf("Expected context to work. But fails. %s", err)
		return
	}
	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err = c.loadListData(nil, context)

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value1, flag, flag_value1, create, false); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, test, flag_value2, flag, flag_value2, create, false); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextMultipleObjectList :
// check if <app> update --tests "flag value, other" --apps "type:driver:name"
// => creates 2 different object 'test' and 'app' records with key and data set.
func TestForjCli_loadListData_contextMultipleObjectList(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test             = "test"
		tests            = "tests"
		test_help        = "test help"
		flag             = "flag"
		flag_help        = "flag help"
		flag_value1      = "flag-value"
		flag_value2      = "other"
		myapp            = "app"
		apps             = "apps"
		myapp_help       = "app help"
		instance         = "instance"
		instance_help    = "instance_help"
		driver_type      = "type"
		driver_type_help = "type help"
		driver           = "driver"
		driver_help      = "driver help"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	// <app> create
	c.NewActions(create, create_help, "create %s", true)
	// <app> update
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, "").
		AddKey(String, flag, flag_help, "#w", nil).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).

		// create list
		CreateList("to_update", ",", "flag", test_help).
		// <app> create tests <data>
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	if c.NewObject(myapp, myapp_help, "").
		AddKey(String, instance, instance_help, "#w", nil).
		AddField(String, driver_type, driver_type_help, "#w", nil).
		AddField(String, driver, driver_help, "#w", nil).
		// <app> create app --instance <instance1> --type <type> --driver <driver>
		DefineActions(create).OnActions().
		AddFlag(instance, nil).
		AddFlag(driver_type, nil).
		AddFlag(driver, nil).

		// create list
		CreateList("to_update", ",", "instance[:driver[:type]]", myapp_help).
		// <app> create apps <data>
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(myapp).Error())
	}

	// <app> update --tests <data>
	c.OnActions(update).
		AddActionFlagFromObjectListAction(test, "to_update", create).
		// <app> update --apps <data>
		AddActionFlagFromObjectListAction(myapp, "to_update", create)

	context := app.NewContext().SetContext(update)
	c.cli_context.context = context

	if _, err := context.SetContextValue(tests, flag_value1); err != nil {
		t.Errorf("Expected context to work. Unable to add '%s' context value. %s", tests, err)
	}
	if _, err := context.SetContextValue(apps, "type:driver:name"); err != nil {
		t.Errorf("Expected context to work. Unable to add '%s' context value. %s", apps, err)
	}

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, context)

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value1, flag, flag_value1, create, false); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, myapp, "type", instance, "type", create, false); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, myapp, "type", driver, "driver", create, false); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, myapp, "type", driver_type, "name", create, false); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_loadListData_contextObjectData :
// check if <app> create test --flag "flag value" --flag2 "value"
// => creates 1 object 'test' record with key and all data set.
func TestForjCli_loadListData_contextObjectData(t *testing.T) {
	t.Log("Expect ForjCli_loadListData() to create object list instances.")

	// --- Setting test context ---
	const (
		test        = "test"
		tests       = "tests"
		test_help   = "test help"
		flag        = "flag"
		flag_help   = "flag help"
		flag_value1 = "flag value"
		flag2       = "flag2"
		flag2_help  = "flag2 help"
		flag_value2 = "other"
		cmd         = "cmd:"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, "").
		AddKey(String, flag, flag_help, "#w", nil).
		AddField(String, flag2, flag2_help, "#w", nil).
		// <app> create test --flag <data> --flag2 <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).
		AddFlag(flag2, nil) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	// <app> update --tests "flag_key"
	c.OnActions(update).AddActionFlagFromObjectListAction(test, "to_update", create)

	context := app.NewContext()

	if ctxt, err := c.App.ParseContext([]string{cmd + create, cmd + test, flag, flag_value1, flag2, flag_value2}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	} else {
		c.cli_context.context = ctxt
	}

	cmds := context.SelectedCommands()
	if len(cmds) == 0 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 2, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	// --- Run the test ---
	err := c.loadListData(nil, context)

	// --- Start testing ---
	// check in cli.
	if err != nil {
		t.Errorf("Expected loadListData to return successfully. But got an error. %s", err)
		return
	}
	if err := check_object_exist(c, test, flag_value1, flag, flag_value1, create, false); err != nil {
		t.Errorf("%s", err)
	}
	if err := check_object_exist(c, test, flag_value1, flag2, flag_value2, create, false); err != nil {
		t.Errorf("%s", err)
	}
}

// TestForjCli_addInstanceFlags:
// check if <app> update --tests "name1,name2" --name1-flag "value" --name2-flag "value2" --apps "test:blabla"
// => creates 1 object 'test' record with key and all data set.
func TestForjCli_addInstanceFlags(t *testing.T) {
	t.Log("Expect ForjCli_LoadContext_withMoreFlags() to create object list instances.")

	// --- Setting test context ---
	const (
		test             = "test"
		tests            = "tests"
		test_help        = "test help"
		flag             = "flag"
		flag_help        = "flag help"
		flag2            = "flag2"
		flag2_help       = "flag2 help"
		flag_value1      = "value"
		flag_value2      = "value2"
		myapp            = "app"
		apps             = "apps"
		myapp_help       = "app help"
		instance         = "instance"
		instance_help    = "instance_help"
		driver_type      = "type"
		driver_type_help = "type help"
		driver           = "driver"
		driver_help      = "driver help"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.AddFieldListCapture("w", w_f)

	c.NewActions(create, create_help, "create %s", true)
	c.NewActions(update, update_help, "update %s", true)

	if c.NewObject(test, test_help, "").
		AddKey(String, flag, flag_help, "#w", nil).
		AddField(String, flag2, flag2_help, "#w", nil).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(flag, nil).
		AddFlag(flag2, nil).

		// create list
		CreateList("to_update", ",", "flag", test_help).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(test).Error())
	}

	if c.NewObject(myapp, myapp_help, "").
		AddKey(String, instance, instance_help, "#w", nil).
		AddField(String, driver_type, driver_type_help, "#w", nil).
		AddField(String, driver, driver_help, "#w", nil).
		// <app> create test --flag <data>
		DefineActions(create).OnActions().
		AddFlag(instance, nil).AddFlag(driver_type, nil).AddFlag(driver, nil).

		// create list
		CreateList("to_update", ",", "type[:driver[:instance]]", myapp_help).
		// <app> create tests "flag_key"
		AddActions(create) == nil {
		t.Errorf("Expected context to work. Got '%s'", c.GetObject(myapp).Error())
	}

	// <app> update --tests "flag_key"
	c.OnActions(update).
		AddActionFlagFromObjectListAction(test, "to_update", create).
		// <app> update --apps "type:driver"
		AddActionFlagFromObjectListAction(myapp, "to_update", create)

	context := app.NewContext()
	c.cli_context.context = context
	if context.SetContext(update) == nil {
		t.Error("Expected SetContext() to work. It fails")
	}
	if _, err := context.SetContextValue(tests, "name1,name2"); err != nil {
		t.Errorf("Expected SetContextValue(tests) to work. It fails. %s", err)
	}
	if _, err := context.SetContextValue(apps, "test:blabla:instance"); err != nil {
		t.Errorf("Expected SetContext(apps) to work. It fails. %s", err)
	}

	if _, err := c.App.ParseContext([]string{}); err != nil {
		t.Errorf("Expected context with ParseContext() to work. Got '%s'", err)
	}

	cmds := context.SelectedCommands()
	if len(cmds) != 1 {
		t.Errorf("Expected context with SelectedCommands() to have '%d' commands. Got '%d'", 1, len(cmds))
		return
	}
	// Ensure objects are identified properly.
	c.identifyObjects(cmds[len(cmds)-1])

	if err := c.loadListData(nil, context); err != nil {
		t.Errorf("Expected loadListData() to work. got '%s'", err)
		return
	}

	// --- Run the test ---
	c.addInstanceFlags()

	// --- Start testing ---
	// checking in cli
	// <app> create tests blabla --blabla-...
	v, found, err := check_list_action_params(c, test+"_to_update", create, "name1-flag")
	if err != nil {
		t.Error(err)
	}
	if found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "name1-flag")
	}

	v, found, err = check_list_action_params(c, test+"_to_update", create, "name1-flag2")
	if err != nil {
		t.Error(err)
	}
	if !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name1-flag2")
	} else {
		if f, ok := v.(*ForjFlag); ok {
			if f.list == nil {
				t.Errorf("Expected '%s' to be attached to the list. Got Nil", "name1-flag2")
			}
			if f.instance_name != "name1" {
				t.Errorf("Expected '%s' to be attached to instance_name '%s'. got '%s'",
					"name1-flag2", "name1", f.instance_name)
			}
			if f.field_name != "flag2" {
				t.Errorf("Expected '%s' to be attached to instance_name '%s'. got '%s'",
					"name1-flag2", "flag2", f.field_name)
			}
		} else {
			t.Errorf("Expected '%s' to be '%s' type. got '%s'",
				"name1-flag2", "*ForjFlag", reflect.TypeOf(v))

		}
	}

	v, found, err = check_list_action_params(c, test+"_to_update", create, "name2-flag2")
	if !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name2-flag2")
	}

	v, found, err = check_list_action_params(c, myapp+"_to_update", create, "instance-"+instance)
	if found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "instance-"+instance)
	}

	v, found, err = check_list_action_params(c, myapp+"_to_update", create, "instance-"+driver)
	if found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "instance-"+driver)
	}

	v, found, err = check_list_action_params(c, myapp+"_to_update", create, "instance-"+driver_type)
	if found {
		t.Errorf("Expected instance flag '%s' to NOT exist. But found it.", "instance-"+driver_type)
	}

	// <add> update --tests blabla,bloblo --blabla-... --bloblo-... --apps blabla:blibli ...
	v, found, err = check_action_params(c, update, "name1-flag2")
	if !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name1-flag2")
	}

	v, found, err = check_action_params(c, update, "name2-flag2")
	if !found {
		t.Errorf("Expected instance flag '%s' to exist. But not found.", "name2-flag2")
	}

	// checking in kingpin
	// <app> create tests blabla --blabla-...
	if app.GetFlag(create, tests, "name1-flag2") == nil {
		t.Errorf("Expected instance flag '%s' to exist in kingpin. But not found.", "name1-flag")
	}
	if app.GetFlag(create, tests, "instance-"+driver_type) != nil {
		t.Errorf("Expected instance flag '%s' to NOT exist in kingpin. But found it.", "instance-"+driver_type)
	}

	// <add> update --tests blabla,bloblo --blabla-... --bloblo-... --apps blabla:blibli ...
	if app.GetFlag(update, "name1-flag2") == nil {
		t.Errorf("Expected instance flag '%s' to exist in kingpin. But not found.", "name1-flag2")
	}
	if app.GetFlag(update, "name2-flag2") == nil {
		t.Errorf("Expected instance flag '%s' to exist in kingpin. But not found it.", "name2-flag2")
	}
	// At context time, instance created flags are not parsed. It will be at next Parse time.
}

func TestForjCli_contextHook(t *testing.T) {
	t.Log("Expect ForjCli_contextHook() to manipulate cli/objects.")

	// --- Setting test context ---
	const (
		test  = "test"
		test2 = "test2"
		test3 = "test3"
		key   = "flag_key"
		field = "field"
	)

	app := kingpinMock.New("Application")
	c := NewForjCli(app)
	// --- Run the test ---
	err, updated := c.contextHook(nil)

	// --- Start testing ---
	if o := c.GetObject(test); o != nil {
		t.Errorf("Expected contextHook() to do nothing. But found the '%s' object.", test)
	}
	if updated {
		t.Errorf("Expected contextHook() to return updated to false. Got '%t'", updated)
	}

	// Update the context
	c.ParseBeforeHook(func(c *ForjCli, _ interface{}) (error, bool) {
		if c == nil {
			return nil, false
		}
		if c.GetObject(test) == nil {
			c.NewObject(test, "", "")
			return nil, true
		}
		return fmt.Errorf("Found object '%s'.", test), false
	})

	// --- Run the test ---
	err, updated = c.contextHook(nil)

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected contextHook() to return no error. Got '%s'", err)
	}
	if o := c.GetObject(test); o == nil {
		t.Errorf("Expected contextHook() to create the ' %s' object. Not found.", test)
	}
	if !updated {
		t.Errorf("Expected contextHook() to return updated to true. Got '%t'", updated)
	}

	// --- Run another test ---
	err, updated = c.contextHook(nil)

	// --- Start testing ---
	if err == nil {
		t.Error("Expected contextHook() to return an error. Got none")
	}
	if fmt.Sprintf("%s", err) != "Found object 'test'." {
		t.Errorf("Expected contextHook() to return a specific error. Got '%s'", err)
	}

	// --- Update the context ---
	c.ParseBeforeHook(nil).
		GetObject(test).ParseHook(func(o *ForjObject, c *ForjCli, _ interface{}) (error, bool) {
		if c == nil || o == nil {
			return nil, false
		}
		if c.GetObject(test2) == nil {
			if o2 := c.NewObject(test2, "", "") ; o2 != nil {
				if o2.AddKey(String, key, "flag help", "", nil) == nil {
					return o2.Error(), false
				}
			} else {
				return c.Error(), false
			}
			return nil, true
		}
		return fmt.Errorf("Found object '%s'.", test2), false
	})

	// --- Run the test ---
	err, updated = c.contextHook(nil)

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected contextHook() to return no error. Got '%s'", err)
	}
	if !updated {
		t.Errorf("Expected contextHook() to return updated to true. Got '%t'", updated)
	}
	if o := c.GetObject(test2); o == nil {
		t.Errorf("Expected contextHook() to create the ' %s' object. Not found.", test2)
	} else {
		if len(o.fields) != 1 {
			t.Errorf("Expected object '%s' to have %d field. Got '%d'.", test2, 1, len(o.fields))
		}
	}

	// --- Update the context ---
	// Cleanup test object hook
	c.GetObject(test).ParseHook(nil)
	// Create test3 object from before hook
	c.ParseBeforeHook(func(c *ForjCli, _ interface{}) (error, bool) {
		if c == nil {
			return nil, false
		}
		o := c.GetObject(test3)
		if o != nil {
			return fmt.Errorf("Found object '%s'.", test3), false
		}
		o = c.NewObject(test3, "", "")
		// Then add a key to the object hook
		o.ParseHook(func(o *ForjObject, c *ForjCli, _ interface{}) (error, bool) {
			if c == nil || o == nil {
				return nil, false
			}
			if o.AddKey(String, key, "flag help", "", nil) == nil {
				return o.Error(), false
			}
			return nil, true
		})
		return nil, true
	})
	// Then add a field from after hook
	c.ParseAfterHook(func(c *ForjCli, _ interface{}) (error, bool) {
		if c == nil {
			return nil, false
		}
		o := c.GetObject(test3)
		if o == nil {
			return fmt.Errorf("Object '%s' not found.", test3), false
		}
		if o.AddField(String, field, "", "#v", nil) != nil {
			return o.Error(), false
		}
		return nil, true
	})

	// --- Run the test ---
	err, updated = c.contextHook(nil)

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected contextHook() to return no error. Got '%s'", err)
	}
	if !updated {
		t.Errorf("Expected contextHook() to return updated to true. Got '%t'", updated)
	}
	if o := c.GetObject(test3); o == nil {
		t.Errorf("Expected contextHook() to create the ' %s' object. Not found.", test3)
	} else {
		if len(o.fields) != 2 {
			t.Errorf("Expected object '%s' to have %d fields. Got '%d'.", test3, 2, len(o.fields))
		}
	}
}

func TestForjCli_Parse_WithDefaultsContext(t *testing.T) {
	t.Log("Expect ForjCli_Parse_WithDefaultsContext() to create objects with all defaults set.")

	const (
		c_test           = "test"
		c_test_help      = "test help"
		c_flag           = "flag"
		c_flag_help      = "flag help"
		c_flag_value     = "flag-value"
		c_flag2          = "flag2"
		c_flag2_help     = "flag2 help"
		c_myDefaultValue = "My default value"
		c_cmd            = "cmd:"
	)
	// --- Setting test context ---
	app := kingpinMock.New("Application")
	c := NewForjCli(app)

	c.NewActions(create, create_help, "", false)
	c.NewObject(c_test, c_test_help, "").
		AddKey(String, c_flag, c_flag_help, "", nil).
		AddField(String, c_flag2, c_flag2_help, "", nil).
		DefineActions(create).OnActions().
		AddFlag(c_flag, Opts().Required()).
		AddFlag(c_flag2, Opts().Default(c_myDefaultValue))

	// --- Run the test ---
	_, err := c.Parse([]string{c_cmd + create, c_cmd + c_test, c_flag, c_flag_value}, nil)

	// --- Start testing ---
	if err != nil {
		t.Errorf("Expected Parse to work. Got '%s'", err)
	}
	if err := check_object_exist(c, c_test, c_flag_value, c_flag2, c_myDefaultValue, create, true); err != nil {
		t.Errorf("%s", err)
	}
}
