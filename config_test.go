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
	if c1.url != "" {
		t.Fatalf("unexported failed should be ignore")
	}
}

// 内嵌结构体测试
func TestMapperEmbedded(t *testing.T) {
	type Mysql struct {
		Url string `env:"URL,some url"`
	}
	type Book struct {
		Name string `env:"NAME,活着"`
	}
	type Config struct {
		Ptr    *Mysql
		Struct Mysql
		Mysql
		*Book
	}

	var c Config
	_ = MustMapConfig(&c)
	if c.Struct.Url == "" {
		t.Fatal("embedded failed")
	}
	if c.Ptr.Url == "" {
		t.Fatal("embedded failed")
	}
	if c.Url == "" {
		t.Fatal("embedded failed")
	}
	if c.Name == "" {
		t.Fatal("embedded failed")
	}
}

// 测试从环境变量读取
func TestMapperENV(t *testing.T) {
	defer os.Clearenv()
	type Config struct {
		Name string `env:"LIST,hello"`
	}
	var c Config
	if err := os.Setenv("LIST", "world"); err != nil {
		panic(err)
	}
	_ = MustMapConfig(&c)
	if c.Name != "world" {
		t.Fatalf("unexpected value")
	}
}

// 测试 "-" 标签
func TestMapperSkip(t *testing.T) {
	defer os.Clearenv()
	type Config struct {
		Name string `env:"-"`
	}
	if err := os.Setenv("NAME", "should be ignored"); err != nil {
		panic(err)
	}
	var c Config
	_ = MustMapConfig(&c)
	if c.Name != "" {
		t.Fatalf("tag with `-` should be ignored")
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
