# go-enumer [![GoDoc](https://godoc.org/github.com/mvrahden/go-enumer?status.svg)](https://godoc.org/github.com/mvrahden/go-enumer) [![Go Report Card](https://goreportcard.com/badge/github.com/mvrahden/go-enumer)](https://goreportcard.com/report/github.com/mvrahden/go-enumer) [![GitHub Release](https://img.shields.io/github/release/mvrahden/go-enumer.svg)](https://github.com/mvrahden/go-enumer/releases)[![Build Status](https://travis-ci.com/mvrahden/go-enumer.svg?branch=master)](https://travis-ci.com/mvrahden/go-enumer)

`go-enumer` is a tool to generate Go code to upgrade Go constants to enums.
It furthermore adds (of unsigned integer types) useful methods to the types, such as validation and (de-)serialization.
It is an opinionated remake of the existing [enumer](https://github.com/dmarkham/enumer) package and therefore behaves different in practically all aspects.

This remake of `go-enumer` is intended to be:

- strict by default (principle of least surprise)
- flexible via opt-in

> E.g. it prevents `undefined` values from being (de-)serialized by default,
> but allows (de-)serialization from empty/undefined values if you configure it to do so.

## What's new?

`go-enumer` is an implementation with improved handling of default values and empty values (Zero Values).
Additionally the type derivation (lookup) was improved and was given more freedom to preferences of case sensitivity.
It furthermore introduces a clear and unified usage pattern on how to implement enums, by bringing its own easy to understand and hard to misuse way of defining them.

`go-enumer` defines an enum as a set custom-defined, distinct values which are implemented by defining a range of continuously incrementing constants of any unsigned integer type or any type alias deriving from those.
By convention, enum sets must:

- be marked with a [magic comment](#magic-comment-goenum).
- be of an unsigned integer type.
- start at `1` or in defined edge cases with `0`.
- consist of continuous linear increments of `1`.

Due to the nature of enums being distinct values, the majority of enum sequences will start at `1`, as you can see in the following snippet.

```go
//go:enum
type Greeting uint

const (
  GreetingWorld Greeting = iota + 1
  GreetingMars
)
```

`go-enumer` gives a special semantic meaning to constants with the value `0`.
The **[Zero Value](https://go.dev/tour/basics/12)** of native integer types in Go is `0`.
For enums it likewise means, that an enum of value `0` is by definition the zero value or the **default** of the enum set.
In some situations you may find yourself in the need of such a **default** value.
Depending on whether your set of values needs a default value, you will chose your sequence to start at value `0` with the your default value being `0`.
All the other values should be `>= 1`.
(You may have multiple defaults, see section for [equivalent values](#equivalent-values))

```go
//go:enum
type Greeting uint

const (
  GreetingWorld Greeting = iota  // <- this is your default
  GreetingMars
)
```

If you want your default value to be robust against deserialization from zero values (like an empty string), then rest assured you do not need to do anything.
However, if you need to also enable deserialization for your **default** enum value from a zero values please check out the section for ["undefined"-value](#the-undefined-value).

If your enum set does not define a default value, but you still want to be able to deserialize empty values into an "undefined" value, please check out the section for ["undefined"-value](#the-undefined-value).
`go-enumer` can generate this "undefined" value for you, as follows and makes it available to you.

```go
const (
  GreetingUndefined Greeting = 0
)
```

### Single pass screening

Thanks to [magic comments](#magic-comment-goenumer) `go-enumer` we can now determine all enums of a package in a single pass, making the code generation much more efficient.

### Magic comment `//go:enum`

The magic comment `//go:enum` serves as a marker to detect all enums.
It also allows for a finegrained configuration on an enum type level, giving the ability to overwrite the global `generate` configuration.
The following example will generate `json` and `yaml` interfaces for the `Greeting` enum, while only generating `json` for all the other enums it can find.
Currently the magic comment supports for all configuration options, which are availble on a global configuration level.

```go
//go:generate github.com/mvrahden/go-enumer -serializers=json

//go:enum -serializers=json,yaml
type Greeting uint

const (
  GreetingWorld Greeting = iota
  GreetingMars
)
```

### Equivalent values

Enums can contain values that are assigned multiple times – such values resemble **equivalent values**.
These alternative values are shadowed by the dominant value, which is always the one which was assigned first.

### Handling of Name Prefixes

In short: You can prefix constant names with their corresponding type alias name and `go-enumer` will automatically strip that off its value – meaning: it will turn `GreetingMars` (with type alias `Greeting`) into e.g. `"MARS"` (assuming an `upper` transformation was applied). It does not allow any other prefixes. This rule is in place to keep your source code concise. Please see [here](#prefix-auto-stripping) for further information.

### Type Validation

By defining enums on a linear scale, the validation time is of constant complexity.
Type validation will be performed on every (de-)serialization.

### Supported features

Supported features are targeted with the `support` flag.

#### The "undefined" feature

`go-enumer` returns with an error upon (de-)serialization of any value which is not explicitly defined in your set of enums with an `error` – this rule covers empty values as well.
Therefore those sequences, that neither contain a default value nor an empty value must start with the integer value of `1`.

The `undefined` feature is an opt-in, which enables the (de-)serialization of zero values.
In case you applied the `undefined` feature, lookups with an empty value will resolve as an "undefined" constant, representing an empty string.
The undefined constant is considered being a **valid** member of the enum set now.
However it will not be represented in the list of possible values.
If you do not have a **default** constant (with the value `0`) `go-enumer` will generate one for you.
Your source code will now support an "undefined" value alongside your explicitly defined set of enums.
However, **unknown** values will still fail upon serialization or deserialization.

| Has Default Value | Supports Undefined | Enum Type can deserialize from undefined | Is Nullable |
|:-----------------:|:------------------:|:----------------------------------------:|:-----------:|
|        no         |         no         |                    no                    |     no      |
|        yes        |         no         |                    no                    |     no      |
|        yes        |        yes         |                   yes                    |     no      |
|        no         |        yes         |                   yes                    |     yes     |

#### Other supported features

With `ent` a method will be generated to return all valid Value strings. This allows you to use your enum type with the ent framework.

## Generated functions and methods

When `go-enumer` is applied to a type, it will generate:

- The following basic methods/functions:

  - Method `String()`: returns the string representation of the enum value. This makes the enum conform
    to the `Stringer` interface, so whenever you print an enum value, you'll get the string name instead of its numeric counterpart.
  - Function `<Type>FromString(raw string)`: returns the enum value from its string representation. This is useful
    when you need to read enum values from command line arguments, from a configuration file, or
    from a REST API request... In short, from those places where using the real enum value (an integer) would
    be almost meaningless or hard to trace or use by a human. `raw` string is case sensitive.
  - Function `<Type>FromStringIgnoreCase(raw string)`: we can not always guarantee the case matching because some systems out of our reach
    are insensitive to exact case matching. In these situations `<Type>FromStringIgnoreCase(raw string)` comes in handy.
    It acts the same as `<Type>FromString(raw string)` with the little difference of `raw` being case insensitive.
  - Function `<Type>Values()`: returns a slice with all the numeric values of the enum, ignoring any alternative values.
  - Function `<Type>Strings()`: returns a slice with all the string representations of the enum.
  - Method `IsValid()`: returns true if the current value is a value of the defined enum set.
  - Method `Validate()`: returns a wrapped error `ErrNoValidEnum` if the current value is not a valid value of the defined enum set.   
    It is being used upon serialization and deserialization, allowing for detecting enum errors via `errors.Is(err, ErrNoValidEnum)`.

- The flag `serializers` in addition with any of the following values, additional methods for serialization are added.
  Valid values are:

  - `binary` makes the enum conform to the `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler` interfaces.
  - `gql` makes the enum conform to the `graphql.Marshaler` and `graphql.Unmarshaler` interfaces.
  - `json` makes the enum conform to the `json.Marshaler` and `json.Unmarshaler` interfaces.
  - `sql` makes the enum conform to the `sql.Scanner` and `sql.Valuer` interfaces.
  - `text` makes the enum conform to the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces.
    **Note:** If you use your enum values as keys in a map and you encode the map as _JSON_,
    you need this flag set to true to properly convert the map keys to json (strings). If not, the numeric values will be used instead
  - `yaml` makes the enum conform to the `gopkg.in/yaml.v2.Marshaler` and `gopkg.in/yaml.v2.Unmarshaler` interfaces.
  - `yaml.v3` makes the enum conform to the `gopkg.in/yaml.v3.Marshaler` and `gopkg.in/yaml.v3.Unmarshaler` interfaces.
    **Note:** Supplying both yaml values (`yaml` and `yaml.v3`) will fail due to interface incompatibility.

For example, if we have an enum type called `Pill`,

> CAUTION: The following example does not use a type prefix.
> Generally it is recommended, to **always** prefix your constants for improved understandability throughout your code base
> \- especially for exported enums.

```go
//go:generate github.com/mvrahden/go-enumer -serializers=json

//go:enum
type Pill uint

const (
  Placebo Pill = iota
  Aspirin
  Ibuprofen
  Paracetamol
  Acetaminophen = Paracetamol
)
```

executing `//go:generate github.com/mvrahden/go-enumer -serializers=json` will generate a new file with four basic package functions and five methods (of which two are for JSON (de-)serialization):

```go
func PillValues() []Pill {
  //...
}

func PillStrings() []string {
  //...
}

func PillFromString(s string) (Pill, error) {
  //...
}

func PillFromStringIgnoreCase(s string) (Pill, error) {
  //...
}

func (i Pill) String() string {
  //...
}

func (i Pill) IsValid() bool {
  //...
}

func (i Pill) Validate() error {
  //...
}

func (i Pill) MarshalJSON() ([]byte, error) {
  //...
}

func (i *Pill) UnmarshalJSON(data []byte) error {
  //...
}
```

From now on, we can:

```go
// Convert any Pill value to string
var aspirinString string = Aspirin.String()
// (or use it in any place where a Stringer is accepted)
fmt.Println("I need ", Paracetamol) // Will print "I need Paracetamol"

// Convert a string with the enum name to the corresponding enum value
pill, err := PillFromString("Ibuprofen")
if err != nil {
    fmt.Println("Unrecognized pill: ", err)
    return
}
// Now pill == Ibuprofen

// Convert a string with the enum name, but degenerated string case
// to the corresponding enum value
pill, err := PillFromStringIgnoreCase("IbUpRoFeN")
if err != nil {
    fmt.Println("Unrecognized pill: ", err)
    return
}

// Get all the values of the string
allPills := PillValues()
fmt.Println(allPills) // Will print [Placebo Aspirin Ibuprofen Paracetamol]

// Check if a value belongs to the Pill enum values
var notAPill Pill = 42
if ok := notAPill.IsValid(); !ok {
  fmt.Println(notAPill, "is not a value of the Pill enum")
}
if err := notAPill.Validate(); err != nil {
  fmt.Printf("%s", err)
}
// Infer whether the error is a validation error
if _, err := json.Marshal(notAPill); errors.Is(err, ErrNoValidEnum) {
  fmt.Printf("this is not a valid enum")
}

// Marshal/unmarshal to/from json strings, either directly or automatically when
// the enum is a field of a struct
pillJSON, _ := Aspirin.MarshalJSON()
// Now pillJSON == `"Aspirin"`
```

## The string representation of the enum value

### Prefix auto-stripping

When dealing with a lot of exported enum values of different type aliases in your project it can easily happen that you run into naming collisions.
To avoid naming collisions while maintaining the same enum string, `go-enumer` automatically strips off the type alias name when determining any string values.
Consider the following example, which will generate the same string values:

```go
//go:enum
type Greeting uint

const (
  Россия Greeting = iota + 1
  中國
  日本
  한국
  ČeskáRepublika
)

//go:enum
type GreetingWithPrefix uint

const (
  GreetingWithPrefixРоссия Greeting = iota + 1
  GreetingWithPrefix中國
  GreetingWithPrefix日本
  GreetingWithPrefix한국
  GreetingWithPrefixČeskáRepublika
)
```

By default, `go-enumer` uses the same name of the enum value for generating the string representation (usually PascalCase in Go).

```go
//go:enum
type Greeting uint

 ...

fmt.Print(ČeskáRepublika) // name => "ČeskáRepublika"
```

Sometimes you need to use some other string representation format than PascalCase.
To transform the string values from PascalCase to another format, you can use the `transform` flag.

For example, the command `//go:generate github.com/mvrahden/go-enumer -transform=whitespace` would generate the following string representation:

```go
fmt.Print(ČeskáRepublika) // name => "Česká Republika"
```

**Note**: If you are missing a transformation please raise an issue and/or open a Pull Request.

### Transformers

Please take the example transformation from the following table for this example:

```go
//go:enum
type MyType uint

const (
  FooBar MyType = iota
  // ...
)
```

| Transformer Name | Example (Stringer value) |
|---|----|
| noop (default) | FooBar |
| camel  | fooBar |
| pascal | FooBar |
| kebab  | foo-bar |
| snake  | foo_bar |
| lower  | foobar |
| upper  | FOOBAR |
| upper-kebab | FOO-BAR |
| upper-snake | FOO_BAR |
| whitespace | Foo Bar |

## Configuration Options

You can add:

- transformation with `transform` option, e.g. `transform=kebab`.
- serializers with `serializers` option, e.g. `serializers=json,sql,...`.
- supported features `support` option, e.g. `support=undefined,ent`
  - `undefined`, see ["undefined"-value](#the-undefined-value)
  - `ignore-case`, adds support for case-insensitive lookup
  - `ent`, adds interface support for [entgo.io](https://github.com/ent/ent)

## Inspiring projects

- [enumer](https://github.com/dmarkham/enumer)
