package cli

import (
	"fmt"
	"github.com/kr/text"
	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
)

// ForjActionRef To define an action reference
type ForjAction struct {
	help          string                      // String which will 'printf' the object name as %s
	name          string                      // Action Name
	cmd           clier.CmdClauser            // Action used at action level
	params        map[string]ForjParam        // Collection of Arguments/Flags
	internal_only bool                        // True if this action cannot be enhanced by plugins
	to_refresh    map[string]*ForjContextTime // List of Object to refresh with context flags
}

func (a *ForjAction) String() string {
	ret := fmt.Sprintf("Action (%p):\n", a)
	ret += fmt.Sprintf("  name: '%s'\n", a.name)
	ret += fmt.Sprintf("  help: '%s'\n", a.help)
	ret += fmt.Sprintf("  internal_only: '%b'\n", a.internal_only)
	ret += fmt.Sprintf("  cmd: '%p'\n", a.cmd)
	ret += fmt.Sprintf("  params: %d\n", len(a.params))
	for key, param := range a.params {
		ret += fmt.Sprintf("    %s:\n", key)
		ret += text.Indent(param.String(), "      ")
	}
	return ret
}

func (a *ForjAction) GetBoolAddr(param_name string) *bool {
	if v, found := a.params[param_name]; found {
		return v.GetBoolAddr()
	}
	return nil
}

func (a *ForjAction) GetStringAddr(param_name string) *string {
	if v, found := a.params[param_name]; found {
		return v.GetStringAddr()
	}
	return nil
}

// ForjContextTime. Structure used at context time to add more flags from Objectlist instances or
// Add Object fields Flags from instances.
type ForjContextTime struct {
	objects_list *ForjObjectList       // List of Object list flags added - Used to add detailed flags
	action       *ForjObjectAction     // Action to refresh with ObjectList detailed flags.
	fields       map[string]*ForjField // List of fields added to the action
}

// Internal function to add ObjectList context information to an action.
func (c *ForjContextTime) addObjectListContext(o_object_list *ForjObjectList, o_action *ForjObjectAction) *ForjContextTime {
	if c == nil {
		c = new(ForjContextTime)
	}
	c.action = o_action
	c.objects_list = o_object_list
	return c
}

// Internal function to add ObjectField context information to an action.
func (c *ForjContextTime) addObjectFieldContext(o_field *ForjField) *ForjContextTime {
	if c == nil {
		c = new(ForjContextTime)
	}
	if c.fields == nil {
		c.fields = make(map[string]*ForjField, 0)
	}
	if _, found := c.fields[o_field.name]; !found {
		c.fields[o_field.name] = o_field
	}
	return c
}

// AddActionFlagFromObjectListAction add one ObjectList action to the selected list of actions (OnActions).
//
// Ex:<app> update --tests "flag_key"
// The collection of object flag can be added at parse time.
// ex: <app> update --tests "key1,key2" --key1-flag <data> --key2-flag <data>
func (c *ForjCli) AddActionFlagFromObjectListAction(obj_name, obj_list, obj_action string) *ForjCli {
	if c == nil {
		return nil
	}
	o_object, o_object_list, o_action, err := c.getObjectListAction(obj_name+"_"+obj_list, obj_action)

	if err != nil {
		c.err = fmt.Errorf("Unable to find object '%s' action '%s'. %s. Adding flags into selected actions ignored.",
			obj_name+"_"+obj_list, obj_action, err)
		return nil
	}

	for _, action := range c.sel_actions {
		action_name := action.name
		if action_name == o_action.name {
			c.err = fmt.Errorf("Unable to add '%s' Action flag to itself.", action_name)
			return nil
		}

		var action *ForjAction

		if a, found := c.actions[action_name]; !found {
			c.err = fmt.Errorf("Unable to find action '%s'. Adding object list action %s '%s-%s' as flag ignored.",
				action_name, obj_action, obj_name, obj_list)
			return nil
		} else {
			action = a
		}

		d_flag := new(ForjFlagList)

		new_object_name := obj_name + "s"
		d_flag.obj = o_object_list

		help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
		d_flag.set_cmd(action.cmd, String, new_object_name, help, nil)
		d_flag.action = o_action.action.name
		action.params[new_object_name] = d_flag

		// Need to add all others object fields not managed by the list, but At context time.
		action_context := action.to_refresh[obj_name]
		action.to_refresh[obj_name] = action_context.addObjectListContext(o_object_list, o_action)

		// Add reference to the Object list for context instance flags creation.
		flags_ref := new(ForjObjectListFlags)
		flags_ref.params = make(map[string]ForjParam)
		flags_ref.multi_actions = false
		flags_ref.objList = o_object_list
		flags_ref.action = action
		gotrace.Trace("Adding reference '%s'", action_name+" --"+new_object_name)
		o_object_list.flags_list[action_name+" --"+new_object_name] = flags_ref

	}
	return c
}

// AddActionFlagsFromObjectListActions add one ObjectList action to the selected list of actions (OnActions).
// Ex: <app> update --add-tests "flag_key" --remove-tests "test,test2"
func (c *ForjCli) AddActionFlagsFromObjectListActions(obj_name, obj_list string, obj_actions ...string) *ForjCli {
	if c == nil {
		return nil
	}
	for _, action := range c.sel_actions {
		action_name := action.name
		for _, obj_action := range obj_actions {
			o_object, o_object_list, o_action, err := c.getObjectListAction(obj_name+"_"+obj_list, obj_action)

			if err != nil {
				c.err = fmt.Errorf("Unable to find object '%s' action '%s'. %s. Adding flags into selected actions ignored.",
					obj_name+"_"+obj_list, obj_action, err)
				return nil
			}

			if action_name == o_action.name {
				c.err = fmt.Errorf("Unable to add '%s' Action flag to itself.", action_name)
				return nil
			}

			var action *ForjAction

			if a, found := c.actions[action_name]; !found {
				c.err = fmt.Errorf("Unable to find action '%s'. Adding object list action %s '%s-%s' as flag ignored.",
					action_name, obj_action, obj_name, obj_list)
				return nil
			} else {
				action = a
			}

			new_obj_name := obj_action + "-" + obj_name + "s"
			d_flag := new(ForjFlagList)
			d_flag.obj = o_object_list
			help := fmt.Sprintf("%s one or more %s", obj_action, o_object.desc)
			d_flag.set_cmd(action.cmd, String, new_obj_name, help, nil)
			d_flag.action = o_action.action.name
			action.params[new_obj_name] = d_flag

			// Need to add all others object fields not managed by the list, but At context time.
			action_context := action.to_refresh[obj_name]
			action.to_refresh[obj_name] = action_context.addObjectListContext(o_object_list, o_action)

			// Add reference to the Object list for context instance flags creation.
			flags_ref := new(ForjObjectListFlags)
			flags_ref.params = make(map[string]ForjParam)
			flags_ref.multi_actions = true
			flags_ref.objList = o_object_list
			flags_ref.action = action
			o_object_list.flags_list[action_name+" --"+new_obj_name] = flags_ref
		}
	}

	return c
}

// AddActionFlagsFromObjectAction create all flags defined on an object action to selected action.
func (c *ForjCli) AddActionFlagsFromObjectAction(obj_name, obj_action string) *ForjCli {
	if c == nil {
		return nil
	}
	o, o_action, _ := c.getObjectAction(obj_name, obj_action)
	if o.fields == nil {
		c.setErr("Unable to add flags from object action '%s-%s' that has no flags declared.", obj_name, obj_action)
		return nil
	}
	for _, action := range c.sel_actions {
		for fname := range o.fields {
			if p, found := o_action.params[fname]; found {
				d_flag := p.Copier().CopyToFlag(action.cmd)
				d_flag.setObjectAction(o_action, fname)
				if o.single {
					d_flag.setObjectInstance(o.name)
				}
				action.params[fname] = d_flag
				o.fields[fname].inActions[action.name] = d_flag
			}
		}
	}
	return c
}

// AddActionFlagFromObjectAction create one flag defined on an object action to selected action.
func (c *ForjCli) AddActionFlagFromObjectAction(obj_name, obj_action, param_name string) *ForjCli {
	if c == nil {
		return nil
	}
	o, o_action, _ := c.getObjectAction(obj_name, obj_action)
	for _, action := range c.sel_actions {
		if _, found := o.fields[param_name]; found {
			if p, found := o_action.params[param_name]; found {
				d_flag := p.Copier().CopyToFlag(action.cmd)
				d_flag.setObjectAction(o_action, param_name)
				if o.single {
					d_flag.setObjectInstance(o.name)
				}
				action.params[param_name] = d_flag
				o.fields[param_name].inActions[action.name] = d_flag
			}
		}
	}
	return c
}

// AddActionFlagFromObjectField declare one flag from an object field to selected action.
// One or more instance flags can be created as soon as object instances are loaded from
// addInstanceFlags() function.
func (c *ForjCli) AddActionFlagFromObjectField(param_name string, options *ForjOpts) *ForjCli {
	if c == nil {
		return nil
	}
	if c.sel_object == nil {
		c.setErr("Object is not selected. Use WithObject or WithObjectInstance functions.")
		return nil
	}
	o := c.sel_object
	instance_name := o.sel_instance

	var field *ForjField
	if instance_name != "" {
		if v, found := o.instances[instance_name].additional_fields[param_name]; found {
			field = v
		}
	} else {
		if v, found := o.fields[param_name]; found {
			field = v
		}
	}
	if field == nil {
		c.setErr("Field '%s' not found in ")
		return nil
	}

	for _, action := range c.sel_actions {
		if o.single {
			d_flag := new(ForjFlag)

			d_flag.setObjectField(o, param_name)
			d_flag.set_cmd(action.cmd, field.value_type, field.name, field.help, options)
			action.params[param_name] = d_flag
			field.inActions[action.name] = d_flag
			gotrace.Trace("Single object '%s' Flag '%s' added to action '%s'.", o.name, field.name, action.name)
		} else {
			if instance_name != "" {
				d_flag := new(ForjFlag)

				d_flag.setObjectField(o, param_name)
				d_flag.setObjectInstance(instance_name)
				d_flag.set_cmd(action.cmd, field.value_type, field.name,
					"Flag for instance "+instance_name+". "+field.help, options)
				action.params[param_name] = d_flag
				field.inActions[action.name] = d_flag
				gotrace.Trace("object instance '%s-%s' Flag '%s' added to action '%s'.",
					o.name, instance_name, field.name, action.name)
			} else {
				// Flags must be added by addInstanceFlags later. == DEAD CODE
				gotrace.Trace("Object '%s' Flag '%s' added to action '%s' context.", o.name, field.name, action.name)
				o_context := action.to_refresh[o.name]
				action.to_refresh[o.name] = o_context.addObjectFieldContext(field)
				gotrace.Trace("Object '%s' Flag '%s' added to action '%s' context.", o.name, field.name, action.name)
			}
		}
	}
	return c
}

// AddArg Add an arg on selected actions
// You could get values With 2 different possible functions (where <Type> can be cli.String or cli.Bool)
// - Get<Type>Value() to get the typed value
// - Get<Type>Addr() to get the pointer to the value
func (c *ForjCli) AddArg(value_type, name, help string, options *ForjOpts) *ForjCli {
	return c.addFlag(func() ForjParam {
		return new(ForjArg)
	}, value_type, name, help, options)
}

// AddFlag Add an flag on selected actions
// You could get values With 2 different possible functions (where <Type> can be cli.String or cli.Bool)
// - Get<Type>Value() to get the typed value
// - Get<Type>Addr() to get the pointer to the value
func (c *ForjCli) AddFlag(value_type, name, help string, options *ForjOpts) *ForjCli {
	return c.addFlag(func() ForjParam {
		return new(ForjFlag)
	}, value_type, name, help, options)
}

func (c *ForjCli) addFlag(newParam func() ForjParam, value_type, name, help string, options *ForjOpts) *ForjCli {
	if c == nil {
		return nil
	}
	for _, action := range c.sel_actions {
		p := newParam()

		p.set_cmd(action.cmd, value_type, name, help, options)

		action.params[name] = p
	}

	return c
}

// NewActions create the list of referenced valid actions supported. kingpin layer created.
// It add them to the kingpin application layer.
//
// name     : Name of the action to add
// help     : Generic help to add to the action.
// for_forjj: True if the action is protected against plugins features.
func (c *ForjCli) NewActions(name, act_help, compose_help string, for_forjj bool) (r *ForjAction) {
	if c == nil {
		return nil
	}
	r = new(ForjAction)
	r.cmd = c.App.Command(name, act_help)
	r.help = compose_help
	r.internal_only = for_forjj
	r.params = make(map[string]ForjParam)
	r.to_refresh = make(map[string]*ForjContextTime)
	r.name = name
	c.actions[name] = r
	return
}

func (c *ForjCli) GetAction(name string) *ForjAction {
	if v, found := c.actions[name]; found {
		return v
	}
	return nil
}

// OnActions Do a selection of action to apply more functionality
func (c *ForjCli) OnActions(actions ...string) *ForjCli {
	if c == nil {
		return nil
	}
	if len(actions) == 0 {
		c.sel_actions = c.actions
		return c
	}
	c.sel_actions = make(map[string]*ForjAction)

	for _, action := range actions {
		if v, err := c.getAction(action); err == nil {
			c.sel_actions[action] = v
		}
	}
	return c
}

// setErr - Set an error flag to the cli and none exists.
func (c *ForjCli) setErr(format string, a ...interface{}) {
	if c.err != nil {
		return
	}
	c.err = fmt.Errorf(format, a...)
}

// cleanErr - Cleanup cli error flag.
func (c *ForjCli) clearErr() error {
	err := c.err
	c.err = nil
	return err
}

func (c *ForjCli) WithObjectInstance(object_name, instance_name string) *ForjCli {
	if c == nil {
		return nil
	}
	if object_name == "" || instance_name == "" {
		c.setErr("object name AND instance name required. Got empty strings.")
		return nil
	}

	o := c.GetObject(object_name)
	if o == nil {
		c.setErr("Object '%s' not found.", object_name)
		return nil
	}
	if _, found := o.instances[instance_name]; !found {
		o.setErr("Instance '%s' is not found.", instance_name)
		return nil
	}
	o.sel_instance = instance_name
	c.sel_object = o
	return c
}

func (a *ForjAction) search_object_param(object, key, param_name string) (p ForjParam) {
	for _, param := range a.params {
		if fl, pi, pn := param.fromList(); fl == nil {
			if fl.obj.name != object || pi != key || pn != param_name {
				continue
			}
			return param
		}
	}
	return p
}
