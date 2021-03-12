package pkg

import (
	"encoding/json"
	"os"
	"testing"
)

func toString(s interface{}) string {
	b, _ := json.MarshalIndent(s, "", "  ")
	return string(b)
}

func TestMapper_Parse(t *testing.T) {
	type sample struct{}
	var s sample
	// non pointer
	m := newMapper(true)
	if err := m.Parse(s); err == nil {
		t.Fatalf("expect error returned")
	}
	// nil
	var s1 *sample
	if err := m.Parse(s1); err == nil {
		t.Fatalf("expect error returned")
	}
	// non struct
	var s2 string
	if err := m.Parse(&s2); err == nil {
		t.Fatalf("expect error returned")
	}

	// no error
	if err := m.Parse(&s); err != nil {
		t.Fatalf("expect no error, got: %s", err)
	}
}

// 非导出字段测试
func TestMapperUnexported(t *testing.T) {
	type Config struct {
		url string `env:"URL,https://ningzi.club"`
	}
	var c Config
	var c1 Config
	_ = MustMapConfig(&c)
	_ = MapConfig(&c1)

	if c.url != "" {
		t.Fatalf("unexported failed should be ignore")
	}
	if c1.url == "" {
		t.Fatalf("unexported failed should be ignore")
	}
}

// 内嵌结构体测试
func TestMapperEmbedded(t *testing.T) {
	type Mysql struct {
		Url string `env:"URL,some url"`
	}
	type Config struct {
		Ptr    *Mysql
		Struct Mysql
		Mysql
	}

	var c Config
	_ = MustMapConfig(&c)
	if c.Struct.Url == "" {
		t.Fatal("embedded failed")
	}
	if c.Ptr.Url == "" {
		t.Fatal("embedded failed")
	}
}

func TestMapConfig(t *testing.T) {
	defer os.Clearenv()
	type config struct {
		A uint8 `env:"A"`
		B int8  `env:"B,10"`
	}
	var s config
	err := MapConfig(&s)
	if err != nil {
		t.Fatalf("unexpected err: %s", err)
	}
	if s.A != 0 {
		t.Fatalf("value error: %d", s.A)
	}
	if s.B != 10 {
		t.Fatalf("value error")
	}
}

func TestMustMapConfig(t *testing.T) {
	defer os.Clearenv()
	type empty struct{}
	err := MustMapConfig(empty{})
	if err == nil {
		t.Fatalf("non-pointer not acceptable")
	}

	var b *empty
	err = MustMapConfig(b)
	if err == nil {
		t.Fatalf("nil not acceptable")
	}

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

// **************************** Benchmark ****************************

func BenchmarkMustMapConfig(b *testing.B) {
	type Embed struct {
		U [5]int64 `env:"U,1,2,3,4,5"`
	}
	type Embed1 struct {
		V string `env:"hello again"`
	}
	type foo struct {
		A int8     `env:"A,-1"`
		B int16    `env:"B,-1"`
		C int32    `env:"C,-1"`
		D int64    `env:"C,-1"`
		E int      `env:"C,-1"`
		F uint8    `env:"F,1"`
		G uint16   `env:"G,1"`
		H uint32   `env:"H,1"`
		I uint64   `env:"I,1"`
		J uint     `env:"J,1"`
		K float32  `env:"K,0.123"`
		L float64  `env:"L,0.123"`
		M bool     `env:"M,True"`
		N string   `env:"N,hello"`
		O []string `env:"O,hello,world"`
		P [2]int64 `env:"P,1,2"`
		Q struct {
			R uint `env:"R,1"`
		}
		S *struct {
			T []bool `env:"T,T,F,0,1"`
		}
		Embed
		*Embed1
	}

	var s foo
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MustMapConfig(&s)
	}
}

func BenchmarkMapConfig(b *testing.B) {
	type Embed struct {
		U [5]int64 `env:"U,1,2,3,4,5"`
	}
	type Embed1 struct {
		V string `env:"hello again"`
	}
	type foo struct {
		A int8     `env:"A,-1"`
		B int16    `env:"B,-1"`
		C int32    `env:"C,-1"`
		D int64    `env:"C,-1"`
		E int      `env:"C,-1"`
		F uint8    `env:"F,1"`
		G uint16   `env:"G,1"`
		H uint32   `env:"H,1"`
		I uint64   `env:"I,1"`
		J uint     `env:"J,1"`
		K float32  `env:"K,0.123"`
		L float64  `env:"L,0.123"`
		M bool     `env:"M,True"`
		N string   `env:"N,hello"`
		O []string `env:"O,hello,world"`
		P [2]int64 `env:"P,1,2"`
		Q struct {
			R uint `env:"R,1"`
		}
		S *struct {
			T []bool `env:"T,T,F,0,1"`
		}
		Embed
		*Embed1
	}

	var s foo
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MapConfig(&s)
	}
}
