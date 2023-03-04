package argp

import (
	"reflect"
	"strings"
	"testing"
)

type Options struct {
	IntOptions   int     `arg:"i"`
	Int64Option  int64   `arg:"int64,specify int64 option"`
	UIntOptions  uint    `arg:"u"`
	UInt64Option uint64  `arg:"uint64,specify int64 option"`
	FloatOption  float64 `arg:"f,specify float/64 option"`
	BoolOption   bool    `arg:"b,specify bool option"`
	StringOption string  `arg:"name,specify name"`
}

type SubCmdOpt struct {
	A int `arg:"a,argument a"`
}

func TestArgp(t *testing.T) {
	cases := []struct {
		defaultValues Options
		expected      Options
		args          []string
	}{
		{
			defaultValues: Options{},
			args:          strings.Split("-i -12 -int64 -64 -u 12 -uint64 64 -f 12.34 -b -name pike", " "),
			expected: Options{
				IntOptions:   -12,
				Int64Option:  -64,
				UIntOptions:  12,
				UInt64Option: 64,
				FloatOption:  12.34,
				BoolOption:   true,
				StringOption: "pike",
			},
		},
		{
			defaultValues: Options{
				IntOptions: -12,
			},
			args: strings.Split("-int64 -64 -u 12 -uint64 64 -f 12.34 -b -name pike", " "),
			expected: Options{
				IntOptions:   -12,
				Int64Option:  -64,
				UIntOptions:  12,
				UInt64Option: 64,
				FloatOption:  12.34,
				BoolOption:   true,
				StringOption: "pike",
			},
		},
	}

	for _, cc := range cases {
		var cmd = &Cmd[Options]{
			CmdName: "root",
			Usage:   "the root command",
			Options: &cc.defaultValues,
			Command: func(opts *Options) error {
				if !reflect.DeepEqual(cc.defaultValues, cc.expected) {
					t.Errorf("not equal: \nexpected: %+v, got\n %+v", cc.expected, cc.defaultValues)
				}
				return nil
			},
		}
		cmd.Run(cc.args)
	}
}

func TestSubCmd(t *testing.T) {
	cases := []struct {
		rootVals                        Options
		subValues                       SubCmdOpt
		args                            []string
		expectedRoot                    Options
		expectedSub                     SubCmdOpt
		expectedRootRun, expectedSubRun bool
	}{
		{
			rootVals:        Options{},
			subValues:       SubCmdOpt{},
			args:            strings.Split("-i 12 sub -a 1", " "),
			expectedRoot:    Options{IntOptions: 12},
			expectedSub:     SubCmdOpt{A: 1},
			expectedRootRun: false, expectedSubRun: true,
		},
		{
			rootVals:        Options{},
			subValues:       SubCmdOpt{},
			args:            strings.Split("-i 12", " "),
			expectedRoot:    Options{IntOptions: 12},
			expectedSub:     SubCmdOpt{A: 1},
			expectedRootRun: true, expectedSubRun: false,
		},
		{
			rootVals: Options{
				StringOption: "pike",
			},
			subValues:       SubCmdOpt{A: 1},
			args:            strings.Split("-i 12 sub", " "),
			expectedRoot:    Options{IntOptions: 12, StringOption: "pike"},
			expectedSub:     SubCmdOpt{A: 1},
			expectedRootRun: true, expectedSubRun: false,
		},
		{
			rootVals: Options{
				StringOption: "pike",
			},
			subValues:       SubCmdOpt{A: 1},
			args:            strings.Split("sub -i 12", " "),
			expectedRoot:    Options{IntOptions: 12, StringOption: "pike"},
			expectedSub:     SubCmdOpt{A: 1},
			expectedRootRun: true, expectedSubRun: false,
		},
	}

	for _, cc := range cases {
		var cmd = &Cmd[Options]{
			CmdName: "root",
			Usage:   "the root command",
			Options: &cc.rootVals,
			Command: func(opts *Options) error {
				if !cc.expectedRootRun {
					t.Errorf("root cmd is not expected to run")
				}
				if !reflect.DeepEqual(cc.rootVals, cc.expectedRoot) {
					t.Errorf("not equal: \nexpected: %+v, got\n%+v", cc.expectedRoot, cc.rootVals)
				}
				return nil
			},
		}
		cmd.AddSubCmd(&Cmd[SubCmdOpt]{
			CmdName: "sub",
			Usage:   "sub command",
			Options: &cc.subValues,
			Command: func(opts *SubCmdOpt) error {
				if !cc.expectedSubRun {
					t.Errorf("sub cmd is not expected to run")
				}
				if !reflect.DeepEqual(cc.rootVals, cc.expectedRoot) {
					t.Errorf("not equal: \nexpected: %+v, got\n%+v", cc.expectedRoot, cc.rootVals)
				}
				if !reflect.DeepEqual(cc.subValues, cc.expectedSub) {
					t.Errorf("not equal: \nexpected: %+v, got\n%+v", cc.expectedSub, cc.subValues)
				}
				return nil
			},
		})
		cmd.Run(cc.args)

	}
}
