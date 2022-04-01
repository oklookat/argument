package argument

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func Test_MatchArgs(t *testing.T) {
	type caser struct {
		num          int
		osArgs       []string
		userArgs     map[int]*argument
		gotVals      []string
		expectedVals []string
	}

	var cases = make([]caser, 0)

	var setGotVals = func(i int, vals []string) {
		if cases[i].gotVals == nil {
			cases[i].gotVals = make([]string, 0)
		}
		cases[i].gotVals = append(cases[i].gotVals, vals...)
	}

	cases = []caser{
		{
			num:          1,
			osArgs:       []string{"--username", "123", "456", "-another"},
			expectedVals: []string{"123", "456"},
			userArgs: map[int]*argument{
				1: {
					fullName:  "username",
					shortName: "u",
					callback: func(values []string) {
						setGotVals(0, values)
					},
				},
			},
		},

		{
			num:          2,
			osArgs:       []string{"--name=oklookat", "-d", "123", "456", "--e"},
			expectedVals: []string{"oklookat", "123", "456"},
			userArgs: map[int]*argument{
				1: {
					fullName:  "data",
					shortName: "d",
					callback: func(values []string) {
						setGotVals(1, values)
					},
				},
				2: {
					fullName:  "name",
					shortName: "n",
					callback: func(values []string) {
						setGotVals(1, values)
					},
				},
			},
		},

		{
			num:          3,
			osArgs:       []string{"--name", "    ", "123", "    ", "--e"},
			expectedVals: []string{"    ", "123", "    "},
			userArgs: map[int]*argument{
				1: {
					fullName:  "name",
					shortName: "n",
					callback: func(values []string) {
						setGotVals(2, values)
					},
				},
			},
		},

		{
			num:          4,
			osArgs:       []string{"--username", "123", "456", "-u", "OVERWRITE_1", "OVERWRITE_2", "-p", "root"},
			expectedVals: []string{"123", "456", "root"},
			userArgs: map[int]*argument{
				1: {
					fullName:  "username",
					shortName: "u",
					callback: func(values []string) {
						setGotVals(3, values)
					},
				},
				2: {
					fullName:  "password",
					shortName: "p",
					callback: func(values []string) {
						setGotVals(3, values)
					},
				},
			},
		},
	}

	for i, _ := range cases {
		matchArgs(cases[i].osArgs, cases[i].userArgs)
		if !reflect.DeepEqual(cases[i].expectedVals, cases[i].gotVals) {
			t.Fatalf("case %v | expected %v | got %v", cases[i].num, cases[i].expectedVals, cases[i].gotVals)
		}
	}
}

type caserValidator struct {
	val      string
	expected bool
}

func (c *caserValidator) Run(cases []caserValidator, run func(val string) bool) error {
	for _, cased := range cases {
		var result = run(cased.val)
		if result != cased.expected {
			var msg = fmt.Sprintf("val %v | expected: %v | got %v", cased.val, cased.expected, result)
			return errors.New(msg)
		}
	}
	return nil
}

func Test_IsArg(t *testing.T) {
	var cases = []caserValidator{
		{
			val:      "hello",
			expected: false,
		},
		{
			val:      "--hello",
			expected: true,
		},
		{
			val:      "----hello",
			expected: true,
		},
		{
			val:      "--hello--world",
			expected: true,
		},
		{
			val:      "-hello",
			expected: true,
		},
		{
			val:      "-hello---world",
			expected: true,
		},
	}
	var cv = caserValidator{}
	var err = cv.Run(cases, isArg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_IsFullArg(t *testing.T) {
	var cases = []caserValidator{
		{
			val:      "hello",
			expected: false,
		},
		{
			val:      "--hello",
			expected: true,
		},
		{
			val:      "----hello",
			expected: true,
		},
		{
			val:      "--hello--world",
			expected: true,
		},
		{
			val:      "-hello",
			expected: false,
		},
		{
			val:      "-hello---world",
			expected: false,
		},
	}
	var cv = caserValidator{}
	var err = cv.Run(cases, isFullArg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_IsShortArg(t *testing.T) {
	var cases = []caserValidator{
		{
			val:      "hello",
			expected: false,
		},
		{
			val:      "--hello",
			expected: false,
		},
		{
			val:      "----hello",
			expected: false,
		},
		{
			val:      "--hello--world",
			expected: false,
		},
		{
			val:      "-hello",
			expected: true,
		},
		{
			val:      "-hello---world-ok",
			expected: true,
		},
		{
			val:      "----hello---world",
			expected: false,
		},
	}
	var cv = caserValidator{}
	var err = cv.Run(cases, isShortArg)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func Test_SplitArgSingle(t *testing.T) {
	type caser struct {
		num          int
		dirty        string
		expectedName string
		expectedVal  string
	}
	var cases = []caser{
		{
			num:          0,
			dirty:        "",
			expectedName: "",
			expectedVal:  "",
		},
		{
			num:          1,
			dirty:        "-hello",
			expectedName: "-hello",
			expectedVal:  "",
		},
		{
			num:          2,
			dirty:        "--hello",
			expectedName: "--hello",
			expectedVal:  "",
		},
		{
			num:          3,
			dirty:        "--hello=1234",
			expectedName: "--hello",
			expectedVal:  "1234",
		},
		{
			num:          4,
			dirty:        "--hello-world====1234",
			expectedName: "--hello-world",
			expectedVal:  "===1234",
		},
	}
	for _, cased := range cases {
		var name, val = splitArgSingle(cased.dirty)
		if name != cased.expectedName {
			t.Fatalf("case %v | expected name %v | got %v", cased.num, cased.expectedName, name)
		}
		if val != cased.expectedVal {
			t.Fatalf("case %v | expected value %v | got %v", cased.num, cased.expectedVal, val)
		}
	}
}

func Test_GetArgName(t *testing.T) {
	type caser struct {
		num      int
		val      string
		expected string
	}
	var cases = []caser{
		{
			num:      0,
			val:      "not-an-argument",
			expected: "not-an-argument",
		},
		{
			num:      1,
			val:      "-hello",
			expected: "hello",
		},
		{
			num:      2,
			val:      "--hello",
			expected: "hello",
		},
		{
			num:      3,
			val:      "----hello--world",
			expected: "--hello--world",
		},
		{
			num:      4,
			val:      "-hello--world--",
			expected: "hello--world--",
		},
		{
			num:      5,
			val:      "--with-value=hello",
			expected: "with-value",
		},
		{
			num:      6,
			val:      "-wv=hello",
			expected: "wv",
		},
	}

	for _, cased := range cases {
		var result = getArgName(cased.val)
		if result != cased.expected {
			t.Fatalf("case %v | expected %v | got %v", cased.num, cased.expected, result)
		}
	}
}
