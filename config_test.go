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

func TestMustMapConfig(t *testing.T) {
	type Embed struct {
		Embed []string `env:"EMBED, Embed"`
	}
	type Embed1 struct {
		Embed1 []string `env:"EMBED1, Embed1"`
	}
	type Sub struct {
		Sub bool `env:"SUB, True"`
	}
	type Sub1 struct {
		Sub1 [2]bool `env:"SUB1, T,0"`
	}
	type Foo struct {
		a int
		A int    `env:"A"`
		B string `env:"B,b"`
		C []int  `env:"C,1,2,3"`
		D Sub
		E *Sub1
		F uint `env:"-"`
		Embed1
		*Embed
	}

	type setter = func()

	var bar = func(key, value string) setter {
		return func() {
			if err := os.Setenv(key, value); err != nil {
				t.Logf("failed to set ENV: %s", err)
			}
		}
	}

	cases := []struct {
		in  Foo
		out Foo
		err bool
		fc  []func()
	}{
		// 使用默认值
		{
			in: Foo{},
			fc: []func(){bar("A", "1")},
			out: Foo{
				a: 0,
				A: 1,
				B: "b",
				C: []int{1, 2, 3},
				D: Sub{
					Sub: true,
				},
				E: &Sub1{[2]bool{true, false}},
				F: 0,
				Embed1: Embed1{
					Embed1: []string{"Embed1"},
				},
				Embed: &Embed{
					Embed: []string{"Embed"},
				},
			},
			err: false,
		},
		// 环境变量类型错误
		{
			in: Foo{},
			fc: []func(){
				bar("A", "not int"),
			},
			err: true,
		},
		{
			in: Foo{},
			fc: []func(){
				bar("SUB1", "1,2,3"), // out of range
			},
			err: true,
		},
		// 设置 tag 为 "-" 的字段的环境变量值
		{
			in: Foo{},
			fc: []func(){
				bar("A", "1"),
				bar("F", "any thing"), // 应该被忽略
			},
			err: false,
			// 全部默认值
			out: Foo{
				a: 0,
				A: 1,
				B: "b",
				C: []int{1, 2, 3},
				D: Sub{
					Sub: true,
				},
				E: &Sub1{[2]bool{true, false}},
				F: 0,
				Embed1: Embed1{
					Embed1: []string{"Embed1"},
				},
				Embed: &Embed{
					Embed: []string{"Embed"},
				},
			},
		},
		// Foo.A 字段为空
		{
			in:  Foo{},
			err: true,
		},
		// 读取环境变量替换默认值
		{
			in: Foo{},
			fc: []func(){
				bar("A", "15"),
				bar("C", "4,5,6"),
				bar("SUB", "false"),
				bar("EMBED1", "reset Embed1"),
			},
			out: Foo{
				a: 0,
				A: 15,
				B: "b",
				C: []int{4, 5, 6},
				D: Sub{
					Sub: false,
				},
				E: &Sub1{[2]bool{true, false}},
				F: 0,
				Embed1: Embed1{
					Embed1: []string{"reset Embed1"},
				},
				Embed: &Embed{
					Embed: []string{"Embed"},
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
