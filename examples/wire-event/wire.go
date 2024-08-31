//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"algorithms/examples/wire-event/foobarbaz"
	"context"

	"github.com/google/wire"
)

// InitializeEvent creates an Event. It will error if the Event is staffed with
// a grumpy greeter.
func InitializeEvent(phrase string) (Event, error) {
	panic(wire.Build(NewEvent, NewGreeter, NewMessage))
}

func InitializeBaz(ctx context.Context) (foobarbaz.Baz, error) {
	panic(wire.Build(foobarbaz.MegaSet))
}

func InjectStructFooBar() foobarbaz.FooBar {
	panic(wire.Build(foobarbaz.StructSet))
}

func InjectValueFoo() foobarbaz.Foo {
	panic(wire.Build(wire.Value(foobarbaz.Foo{X: 42})))
}

func InjectedMessage() foobarbaz.Foo {
	panic(wire.Build(
		InjectStructFooBar,
		wire.FieldsOf(new(foobarbaz.FooBar), "MyFoo"),
	))
}
