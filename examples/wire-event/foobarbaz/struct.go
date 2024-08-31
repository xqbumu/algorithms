package foobarbaz

import "github.com/google/wire"

type FooBar struct {
	MyFoo Foo
	MyBar Bar
}

func ProvideFooBar() FooBar {
	return FooBar{
		MyFoo: Foo{},
		MyBar: Bar{},
	}
}

var StructSet = wire.NewSet(
	ProvideFoo,
	ProvideBar,
	wire.Struct(new(FooBar), "MyFoo", "MyBar"),
)
