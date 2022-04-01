package argument

import (
	"errors"
	"os"
)

var (
	ErrNoCallback = errors.New("[argument] nil callback func")
	ErrNoFullName = errors.New("[argument] no fullName provided")
	ErrSameNames  = errors.New("[argument] names cannot be same")
)

// one argument.
type argument struct {
	// full name like: "create-superuser".
	fullName string

	// short name like: "csu".
	shortName string

	// when one of names triggered, run callback.
	callback func(values []string)
}

// create new argument instance.
func New() *Instance {
	return &Instance{}
}

// argument instance.
type Instance struct {
	args map[int]*argument
}

// add argument.
//
// fullName: full arg name like "create-superuser".
//
// shortName: short arg name like "csu".
//
// cb: when one of names triggered, run this func.
func (i *Instance) Add(fullName string, shortName string, cb func(values []string)) error {
	if len(fullName) < 1 {
		return ErrNoFullName
	}
	if fullName == shortName {
		return ErrSameNames
	}
	if cb == nil {
		return ErrNoCallback
	}
	if i.args == nil {
		i.args = make(map[int]*argument, 0)
	}
	var newArg = &argument{}
	newArg.fullName = fullName
	newArg.shortName = shortName
	newArg.callback = cb
	i.args[len(i.args)+1] = newArg
	return nil
}

// start searching arguments.
func (i *Instance) Start() {
	var osArgs = os.Args
	if i.args == nil || len(osArgs) < 1 {
		return
	}
	matchArgs(osArgs, i.args)
}
