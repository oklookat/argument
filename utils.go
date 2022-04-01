package argument

import (
	"strings"
)

const (
	fullArgPrefix        = "--"
	shortArgPrefix       = "-"
	singleValueDelimiter = "="
)

// is short or full arg?
func isArg(val string) bool {
	return isFullArg(val) || isShortArg(val)
}

// is val strarts with "--"?
func isFullArg(val string) bool {
	return strings.HasPrefix(val, fullArgPrefix)
}

// is val strarts with "-"?
func isShortArg(val string) bool {
	return !isFullArg(val) && strings.HasPrefix(val, shortArgPrefix)
}

// is arg contains "=" sign?
func isArgWithOneValue(val string) bool {
	return strings.Contains(val, singleValueDelimiter)
}

// get argument name without preifx and value.
func getArgName(val string) string {
	// trim prefix.
	if isFullArg(val) {
		val = strings.TrimPrefix(val, fullArgPrefix)
	} else if isShortArg(val) {
		val = strings.TrimPrefix(val, shortArgPrefix)
	} else {
		return val
	}

	// if arg like "create-user=admin"
	// we need split name and value like [create-user, admin]
	if isArgWithOneValue(val) {
		var name, _ = splitArgSingle(val)
		return name
	}

	return val
}

// split arg name and value by delim sign.
//
// dirtyArg: arg like "--hello-world=1234"
//
// returns:
//
// name like "--hello-world"
//
// val like "1234"
func splitArgSingle(dirtyArg string) (name, val string) {
	name = ""
	val = ""

	// split by first delim sign.
	var sliced = strings.SplitN(dirtyArg, singleValueDelimiter, 2)

	// name.
	if len(sliced) > 0 {
		name = sliced[0]
	}

	// delim sign.
	if len(sliced) > 1 {
		val = sliced[1]
	}

	return
}

// get many values of arg.
//
// like: "--username hello world --another-arg"
//
// returns: [hello, world]
func getManyValues(i int, osArgs []string) []string {
	// if (maybe) many values.
	var values []string
	for iLocal := i + 1; iLocal < len(osArgs); iLocal++ {
		var localOsArg = osArgs[iLocal]

		// if we found new arg, stop.
		if isArg(localOsArg) {
			break
		}

		// add arg and update counter to skip.
		values = append(values, localOsArg)
	}
	return values
}

// match args, run callback if matched.
func matchArgs(osArgs []string, userArgs map[int]*argument) {

	// to avoid args overwriting.
	var disableProcessing = make(map[int]struct{}, 0)

	// for os args.
	for i, osArg := range osArgs {

		// go next it no arg.
		if !isArg(osArg) {
			continue
		}

		var osArgName = getArgName(osArg)

		// compare program argument with user arguments.
		for j, userArg := range userArgs {

			// avoid args overwriting.
			if _, ok := disableProcessing[j]; ok {
				continue
			}

			// compare.
			var isSame = userArg.fullName == osArgName || userArg.shortName == osArgName
			if !isSame {
				continue
			}

			// if single value.
			if isArgWithOneValue(osArg) {
				var _, val = splitArgSingle(osArg)
				userArg.callback([]string{val})
				disableProcessing[j] = struct{}{}
				break
			}

			// if (maybe) many values.
			var values = getManyValues(i, osArgs)
			userArg.callback(values)
			disableProcessing[j] = struct{}{}
		}
	}

}
