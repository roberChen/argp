package argp

import (
	"flag"
	"fmt"
	"reflect"
	"strings"
)

const (
	argpFiledName = "arg"
)

// parse all arguments from args to option struct
func parseOptions(flagset **flag.FlagSet, option any, args []string) ([]string, error) {
	optptr := reflect.ValueOf(option)
	opttype := optptr.Type()
	if opttype.Kind() != reflect.Pointer {
		return nil, fmt.Errorf("options not passed as pointer")
	}
	optval := reflect.Indirect(optptr)
	if optval.Kind() != reflect.Struct {
		return nil, fmt.Errorf("options is not pointer to struct but %s", optval.Kind())
	}
	*flagset = flag.NewFlagSet(optval.Type().Name(), flag.ExitOnError) // the flag set for option
	flagSet := *flagset
	// bind flags for fields
	for i := 0; i < optval.NumField(); i++ {
		name, describe := optionInfo(optval.Type().Field(i).Tag.Get(argpFiledName))
		field := optval.Field(i)
		switch field.Kind() {
		case reflect.Int:
			flagSet.IntVar((*int)(field.Addr().UnsafePointer()), name, int(field.Int()), describe)
		case reflect.Int64:
			flagSet.Int64Var((*int64)(field.Addr().UnsafePointer()), name, field.Int(), describe)
		case reflect.Uint:
			flagSet.UintVar((*uint)(field.Addr().UnsafePointer()), name, uint(field.Uint()), describe)
		case reflect.Uint64:
			flagSet.Uint64Var((*uint64)(field.Addr().UnsafePointer()), name, field.Uint(), describe)
		case reflect.Float64:
			flagSet.Float64Var((*float64)(field.Addr().UnsafePointer()), name, field.Float(), describe)
		case reflect.String:
			flagSet.StringVar((*string)(field.Addr().UnsafePointer()), name, field.String(), describe)
		case reflect.Bool:
			flagSet.BoolVar((*bool)(field.Addr().UnsafePointer()), name, field.Bool(), describe)
		default:
			return nil, fmt.Errorf("unexpected option type %v", field.Type())
		}
	}
	flagSet.Parse(args)
	return flagSet.Args(), nil
}

func optionInfo(tagfield string) (name, describe string) {
	splited := strings.Split(tagfield, ",")
	if len(splited) == 1 {
		return splited[0], ""
	} else if len(splited) == 2 {
		return splited[0], splited[1]
	}
	return "", ""
}
