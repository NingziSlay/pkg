package components

import (
	"encoding/json"
	"os"
	"testing"
)

type sample struct {
	unexported   int      `env:"UNEXPORTED, 110"` // 非导出字段
	Int          int      `env:"INT, 1"`
	Uint         uint     `env:"UINT, 1"`
	Bool         bool     `env:"BOOL, true"`
	String       string   `env:"STRING, Mariah Carey"`
	SliceInt     []int    `env:"SLICE_INT,1,2,3"`
	SliceString  []string `env:"SLICE_STRING,Fly, Like, A, Bird"`
	SubStruct    sub
	SubStructPtr *sub1
	Embed
}

func toString(s *sample) string {
	b, _ := json.Marshal(s)
	return string(b)
}

type Embed struct {
	Array [2]bool `env:"ARRAY,0,1"`
}

type sub struct {
	SliceString []string `env:"SLICE_STRING_1,We, Belong, Together"`
}

type sub1 struct {
	String string
}

func TestMustMapConfig(t *testing.T) {
	type setter func()

	var Setter = func(env, value string) setter {
		return func() {
			_ = os.Setenv(env, value)
		}
	}

	cases := []struct {
		input  *sample
		expect *sample
		err    bool
		fn     []setter
	}{
		{
			input:  &sample{},
			expect: &sample{},
			err:    true,
		},
		{
			input: &sample{},
			expect: &sample{
				unexported:   0,
				Int:          01,
				Uint:         1,
				Bool:         true,
				String:       "hello world",
				SliceInt:     []int{1, 2, 3},
				SliceString:  []string{"Fly", "Like", "A", "Bird"},
				SubStruct:    sub{SliceString: []string{"We", "Belong", "Together"}},
				SubStructPtr: &sub1{String: "hello world"},
				Embed:        Embed{Array: [2]bool{false, true}},
			},
			err: false,
			fn: []setter{
				Setter("String", "hello world"),
			},
		},
	}
	for _, c := range cases {
		for _, f := range c.fn {
			f()
		}
		err := MustMapConfig(c.input)
		if c.err {
			if err == nil {
				t.Fatal("error should not be nil")
			}
			//}
			//} else if !c.err {
			//	if err != nil {
			//		t.Logf("error should be nil, got: %s", err)
			//	}
		} else {
			s1 := toString(c.input)
			s2 := toString(c.expect)
			if s1 != s2 {
				t.Fatalf("expect result: %s, got: %s", s2, s1)
			}
		}
	}
}
