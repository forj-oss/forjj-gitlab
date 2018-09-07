package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/forj-oss/forjj-modules/cli/interface"
	"github.com/forj-oss/forjj-modules/trace"
	"github.com/kr/text"
)

// ForjCli is the Core cli for forjj command.
type ForjCli struct {
	App          clier.Applicationer                       // *kingpin.Application       // Kingpin Application object
	flags        map[string]*ForjFlag                      // Collection of Objects at Application level
	objects      map[string]*ForjObject                    // Collection of Objects that forjj will manage.
	actions      map[string]*ForjAction                    // Collection recognized actions
	list         map[string]*ForjObjectList                // Collection of object list
	cli_context  ForjCliContext                            // Context from cli parsing
	values       map[string]*ForjRecords                   // Collection of Object Values.
	filters      map[string]string                         // List of field data identification from a list.
	err          error                                     // Last error found.
	bef_ctx_hook func(*ForjCli, interface{}) (error, bool) // Last parse hook applied on cli.
	aft_ctx_hook func(*ForjCli, interface{}) (error, bool) // Last parse hook applied on cli.
	parse        bool                                      // true is parse task is done.
	cur_cmds     []clier.CmdClauser

	sel_actions map[string]*ForjAction // Selected actions
	sel_object  *ForjObject            // Selected Object
}

// GetAllActions return the list of actions and definitions defined by the application.
func (c *ForjCli) GetAllActions() map[string]*ForjAction {
	if c == nil {
		return nil
	}
	return c.actions
}

// ParseBeforeHook defines a hook that will be executed before object list and object hooks
func (c *ForjCli) ParseBeforeHook(context_hook func(*ForjCli, interface{}) (error, bool)) *ForjCli {
	if c == nil {
		return nil
	}
	c.bef_ctx_hook = context_hook
	return c
}

// ParseAfterHook defines a hook that will be executed after object list and object hooks
func (c *ForjCli) ParseAfterHook(context_hook func(*ForjCli, interface{}) (error, bool)) *ForjCli {
	if c == nil {
		return nil
	}
	c.aft_ctx_hook = context_hook
	return c
}

func (c *ForjCli) IsParsePhase() bool {
	if c == nil {
		return false
	}
	return c.parse
}

// Parse do the parse of the command line
func (c *ForjCli) Parse(args []string, context interface{}) (cmd string, err error) {
	c.parse = false
	err = c.loadContext(args, context)
	if err != nil {
		return
	}

	// Load all object extra flags/arg data
	c.parse = true
	if cmd, err = c.App.Parse(args); err != nil {
		return
	}

	err = c.loadObjectData()
	return
}

// GetParseContext return the internal parseContext object
func (c *ForjCli) GetParseContext() clier.ParseContexter {
	if c == nil {
		return nil
	}
	return c.cli_context.context
}

func (c *ForjCli) GetCurrentCommand() []clier.CmdClauser {
	return c.cur_cmds
}

func (c *ForjCli) String() (ret string) {
	ret = fmt.Sprintf("clier.Applicationer: %p\n", c.App)
	ret += fmt.Sprintf("context : %s\n", c.cli_context)
	if c.bef_ctx_hook == nil {
		ret += fmt.Sprintf("Before Hook : %t\n", false)
	} else {
		ret += fmt.Sprintf("Before Hook : %t\n", true)
	}
	if c.aft_ctx_hook == nil {
		ret += fmt.Sprintf("After Hook : %t\n", false)
	} else {
		ret += fmt.Sprintf("After Hook : %t\n", true)
	}
	ret += fmt.Sprint("Flags (map):\n")
	for key, flag := range c.flags {
		ret += fmt.Sprintf("  %s: \n", key)
		ret += text.Indent(flag.String(), "    ")
	}
	ret += fmt.Sprint("Actions (map):\n")
	for key, action := range c.actions {
		ret += fmt.Sprintf("  %s: \n", key)
		ret += text.Indent(action.String(), "    ")
	}
	ret += fmt.Sprint("Objects (map):\n")
	for key, object := range c.objects {
		ret += fmt.Sprintf("  %s: \n", key)
		ret += text.Indent(object.String(), "    ")
	}
	ret += fmt.Sprint("Objects list (map):\n")
	for key, list := range c.list {
		ret += fmt.Sprintf("  %s: (%p)\n", key, list)
		ret += text.Indent(list.Stringer(), "    ")
	}
	ret += fmt.Sprint("Values (map):\n")
	for key, value := range c.values {
		ret += fmt.Sprintf("  %s: \n", key)
		ret += text.Indent(value.String(), "    ")
	}
	return
}

func (c *ForjCli) Error() error {
	return c.clearErr()
}

type ForjListParam interface {
	IsFound() bool
	GetAll() []ForjData
	IsList() bool
}

type ForjParamCopier interface {
	CopyToFlag(clier.CmdClauser) *ForjFlag
	CopyToArg(clier.CmdClauser) *ForjArg
}

type ForjParam interface {
	Name() string
	String() string
	Type() string
	IsFound() bool
	GetBoolValue() bool
	GetStringValue() string
	GetBoolAddr() *bool
	GetStringAddr() *string
	GetContextValue(clier.ParseContexter) (interface{}, bool)
	GetValue() interface{}
	// Update default
	Default(string) ForjParam

	// Create kingpin flag
	set_cmd(clier.CmdClauser, string, string, string, *ForjOpts)
	// Set kingpin options
	set_options(*ForjOpts)

	loadFrom(clier.ParseContexter)
	IsList() bool
	isListRelated() bool
	fromList() (*ForjObjectList, string, string)
	isObjectRelated() bool
	IsFromObject(*ForjObject) bool
	getObject() *ForjObject
	CopyToFlag(clier.CmdClauser) *ForjFlag
	CopyToArg(clier.CmdClauser) *ForjArg

	// Additional public interfaces

	Copier() ForjParamCopier

	// Internal interface
	forjParam() forjParam
	forjParamList() forjParamList
	forjParamRelated() forjParamRelated
	forjParamRelatedSetter() forjParamRelatedSetter
	forjParamSetter() forjParamSetter
	forjParamUpdater() forjParamUpdater
}

type forjParamUpdater interface {
	updateContextData()
	set_ref(*ForjData)
}

type ForjKingpinParam interface {
	GetKFlag()
}

type forjParam interface {
	GetFlag() *ForjFlag
	GetArg() *ForjArg
}

type forjParamList interface {
	getInstances() []string
	createObjectDataFromParams(params map[string]ForjParam) error
}

type forjParamObject interface {
	UpdateObject() error
	getObject() *ForjObject
}

type forjParamRelated interface {
	getFieldName() string
	getInstanceName() string
	getObjectList() *ForjObjectList
	getObjectAction() *ForjObjectAction
	getObject() *ForjObject
}

type forjParamSetter interface {
	createObjectDataFromParams(map[string]ForjParam) error
}

type forjParamRelatedSetter interface {
	setList(*ForjObjectList, string, string)
	setObjectAction(*ForjObjectAction, string)
	setObjectField(*ForjObject, string)
	setObjectInstance(string)
}

// ForjParams type
const (
	// Arg : To set a ForjParam as Argument.
	Arg = "arg"
	// Flag : To set a ForjParam as Flag.
	Flag = "flag"
)

// List of ForjParams internal types
const (
	// String - Define the Param data type to string.
	String = "string"
	// Bool - Define the Param data type to bool.
	Bool = "bool"
	// List - Define a ForjObjectList data type.
	List = "list"
)

// NewForjCli : Initialize a new ForjCli object
//
// panic if app is nil.
func NewForjCli(app clier.Applicationer) (c *ForjCli) {
	if app.IsNil() {
		panic("kingpin.Application cannot be nil.")
	}
	c = new(ForjCli)
	c.objects = make(map[string]*ForjObject)
	c.actions = make(map[string]*ForjAction)
	c.flags = make(map[string]*ForjFlag)
	c.values = make(map[string]*ForjRecords)
	c.list = make(map[string]*ForjObjectList)
	c.filters = make(map[string]string)
	c.sel_actions = make(map[string]*ForjAction)
	c.App = app
	return
}

// AddFieldListCapture Add a Field list capture
func (c *ForjCli) AddFieldListCapture(key, capture string) error {
	if _, found := c.filters[key]; found {
		return fmt.Errorf("Key '%s' already exist.", key)
	}

	if _, err := regexp.Compile(capture); err != nil {
		return fmt.Errorf("Capture '%s' not created: Regexp error found: %s", capture, err)
	} else {
		parentheses_reg, _ := regexp.Compile(`\(`)
		if len(parentheses_reg.FindAllString(strings.Replace(capture, `\(`, "", -1), -1)) == 0 {
			capture = "(" + capture + ")"
		}
	}

	c.filters[key] = capture
	return nil
}

// AddAppFlag create a flag object at the application layer.
func (c *ForjCli) AddAppFlag(paramIntType, name, help string, options *ForjOpts) {
	f := new(ForjFlag)
	f.flag = c.App.Flag(name, help)
	f.name = name
	f.help = help
	f.value_type = paramIntType
	f.set_options(options)

	switch paramIntType {
	case String:
		f.flagv = f.flag.String()
	case Bool:
		f.flagv = f.flag.Bool()
	}
	c.flags[name] = f
}

func (c *ForjCli) buildCapture(selector string) string {
	for key, capture := range c.filters {
		selector = strings.Replace(selector, "#"+key, capture, -1)
	}
	return strings.Replace(selector, "##", "#", -1)
}

// getValue : Core get value code for GetBoolValue and GetStringValue
func (c *ForjCli) getValue(object, key, param_name string) (interface{}, bool, error) {
	var value *ForjRecords

	if v, found := c.values[object]; !found {
		return nil, false, fmt.Errorf("Unable to find Object '%s'", object)
	} else {
		value = v
	}

	if v, found, err := value.Get(key, param_name); found {
		return v, true, nil
	} else {
		return nil, false, err
	}
}

// newParam create a ForjFlag or ForjArg defined by `paramType`
func (c *ForjCli) newParam(paramType, name string) ForjParam {
	switch paramType {
	case Arg:
		return new(ForjArg)
	case Flag:
		return new(ForjFlag)
	case List:
		l := new(ForjFlagList)
		if v, found := c.list[name]; found {
			l.obj = v
		} else {
			gotrace.Trace("Unable to find '%s' list.", name)
		}
		return l
	}
	return nil
}

// Create the ForjAction object to attach to the object parent.
func newForjObjectAction(ar *ForjAction, obj *ForjObject, name, desc string) *ForjObjectAction {
	a := new(ForjObjectAction)
	a.action = ar
	a.name = ar.name + "_" + name
	a.cmd = ar.cmd.Command(name, fmt.Sprintf(ar.help, desc))
	a.params = make(map[string]ForjParam)
	a.plugins = make([]string, 0, 5)
	a.obj = obj
	gotrace.Trace("kingping Command '%s' added.", name)
	return a
}

func (c *ForjCli) getObject(obj_name string) (*ForjObject, error) {
	if c == nil {
		return nil, fmt.Errorf("Invalid cli object. it is nil.")
	}
	if v, found := c.objects[obj_name]; found {
		return v, nil
	}
	return nil, fmt.Errorf("Unable to find object '%s'", obj_name)
}

func (c *ForjCli) getObjectAction(obj_name, action string) (o *ForjObject, a *ForjObjectAction, err error) {
	err = nil
	if v, err := c.getObject(obj_name); err != nil {
		return nil, nil, err
	} else {
		o = v
	}

	if v, found := o.actions[action]; !found {
		return nil, nil, fmt.Errorf("Unable to find action '%s' from object '%s'", action, obj_name)
	} else {
		a = v
	}
	return
}

func (c *ForjCli) getObjectListAction(list_name, action string) (o *ForjObject, l *ForjObjectList, a *ForjObjectAction, err error) {
	if c == nil {
		return nil, nil, nil, fmt.Errorf("Invalid cli object. It is nil.")
	}
	err = nil
	if v, found := c.list[list_name]; !found {
		return nil, nil, nil, fmt.Errorf("Unable to find object list '%s'", list_name)
	} else {
		l = v
		o = l.obj
	}

	if v, found := o.actions[action]; !found {
		return nil, nil, nil, fmt.Errorf("Unable to find action '%s' from object '%s'", action, list_name)
	} else {
		a = v
	}
	return
}

func (c *ForjCli) getAction(action string) (a *ForjAction, err error) {
	err = nil
	if v, found := c.actions[action]; !found {
		return nil, fmt.Errorf("Unable to find action '%s'", action)
	} else {
		a = v
	}
	return
}
