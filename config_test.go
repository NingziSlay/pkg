package components

import (
	"encoding/json"
	"os"
	"testing"
)

func toString(s interface{}) string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

type embed struct {
	Embed []string `env:"EMBED, embed"`
}
type embed1 struct {
	Embed1 []string `env:"EMBED1, embed1"`
}
type sub struct {
	Sub bool `env:"SUB, True"`
}
type sub1 struct {
	Sub1 [2]bool `env:"SUB1, T,0"`
}
type foo struct {
	a int
	A int    `env:"A"`
	B string `env:"B,b"`
	C []int  `env:"C,1,2,3"`
	D sub
	E *sub1
	F uint `env:"-"`
	embed1
	*embed
}

func TestMustMapConfig(t *testing.T) {

	type setter = func()

	var bar = func(key, value string) setter {
		return func() {
			if err := os.Setenv(key, value); err != nil {
				t.Logf("failed to set ENV: %s", err)
			}
		}
	}

	cases := []struct {
		in  foo
		out foo
		err bool
		fc  []func()
	}{
		// 使用默认值
		{
			in: foo{},
			fc: []func(){bar("A", "1")},
			out: foo{
				a: 0,
				A: 1,
				B: "b",
				C: []int{1, 2, 3},
				D: sub{
					Sub: true,
				},
				E: &sub1{[2]bool{true, false}},
				F: 0,
				embed1: embed1{
					Embed1: []string{"embed1"},
				},
				embed: &embed{
					Embed: []string{"embed"},
				},
			},
			err: false,
		},
		// 环境变量类型错误
		{
			in: foo{},
			fc: []func(){
				bar("A", "not int"),
			},
			err: true,
		},
		{
			in: foo{},
			fc: []func(){
				bar("SUB1", "1,2,3"), // out of range
			},
			err: true,
		},
		// 设置 tag 为 "-" 的字段的环境变量值
		{
			in: foo{},
			fc: []func(){
				bar("A", "1"),
				bar("F", "any thing"), // 应该被忽略
			},
			err: false,
			// 全部默认值
			out: foo{
				a: 0,
				A: 1,
				B: "b",
				C: []int{1, 2, 3},
				D: sub{
					Sub: true,
				},
				E: &sub1{[2]bool{true, false}},
				F: 0,
				embed1: embed1{
					Embed1: []string{"embed1"},
				},
				embed: &embed{
					Embed: []string{"embed"},
				},
			},
		},
		// foo.A 字段为空
		{
			in:  foo{},
			err: true,
		},
		// 读取环境变量替换默认值
		{
			in: foo{},
			fc: []func(){
				bar("A", "15"),
				bar("C", "4,5,6"),
				bar("SUB", "false"),
				bar("EMBED1", "reset embed1"),
			},
			out: foo{
				a: 0,
				A: 15,
				B: "b",
				C: []int{4, 5, 6},
				D: sub{
					Sub: false,
				},
				E: &sub1{[2]bool{true, false}},
				F: 0,
				embed1: embed1{
					Embed1: []string{"reset embed1"},
				},
				embed: &embed{
					Embed: []string{"embed"},
				},
			},
			err: false,
		},
	}

	for _, c := range cases {
		os.Clearenv()
		for _, f := range c.fc {
			f()
		}
		err := MustMapConfig(&c.in)
		if c.err {
			if err == nil {
				t.Fatalf("expect error but got nothing")
			}
		} else {
			if err != nil {
				t.Fatalf("unexpect err: %s", err)
			}
			in := toString(c.in)
			out := toString(c.out)
			if in != out {
				t.Fatalf("expect: %s\n got: %s", out, in)
			}
		}
	}
}

func BenchmarkMustMapConfig(b *testing.B) {

}
