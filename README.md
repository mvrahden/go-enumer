# go-enumer <!-- omit in toc -->

`go-enumer` is a tool to generate Go code upgrading specifically configured Go constants to enums.
It adds useful common methods to these types, such as validation and various (de-)serialization.
It is an opinionated remake of the long existing [enumer](https://github.com/dmarkham/enumer) package and therefore behaves different in practically all aspects.

`go-enumer` is intended to be:

- strict by default (principle of least surprise)
- flexible via opt-in

> E.g. it prevents `undefined` (or empty) values from being (de-)serialized by default.
> But you can always choose to configure your enums to permit (de-)serialization from such values.

## Table of Contents <!-- omit in toc -->

1. [Why `go-enumer`?](#why-go-enumer)
2. [What's it all about?](#whats-it-all-about)
   1. [Conventions](#conventions)
   2. [Defaults and Zero-Values](#defaults-and-zero-values)
   3. [Comment directive `//go:enum`](#comment-directive-goenum)
   4. [Validation](#validation)
   5. [Supported features](#supported-features)
      1. [The "undefined" feature](#the-undefined-feature)
      2. [Other supported features](#other-supported-features)
3. [Simple Block Spec](#simple-block-spec)
   1. [String Case Transformations](#string-case-transformations)
   2. [Handling of Name Prefixes](#handling-of-name-prefixes)
   3. [Alternative values](#alternative-values)
4. [Filebased Spec](#filebased-spec)
   1. [CSV-File sources](#csv-file-sources)
5. [Generated functions and methods](#generated-functions-and-methods)
6. [Configuration Options](#configuration-options)
7. [Caveats](#caveats)
8. [Inspiring projects](#inspiring-projects)

## Why `go-enumer`?

First and foremost enums are a common way of handling sets of distinct and customizable (and eventually human-readable) values.
Many projects have use cases where enums might be useful, but Go has no built-in support for enums.
Therefore this project aims to fill the gap, by making use of Go's features for type declaration and code generation.

In regards to the aforementioned ancestor project [enumer](https://github.com/dmarkham/enumer), I found some unfortunate obstacles like the configuration, strict case handling, validation and feature extensibility.
Its scope, ambitions and constraints were clearly different, since it tried to be a "drop-in" replacement for the known "Stringer" utility.

Starting from scratch, I decided it's best that `go-enumer` breaks with those restrictions while borrowing some of the neat ideas.

## What's it all about?

`go-enumer` understands an enum as individual type with their sets of custom-defined, distinct values - let's refer to such sets of enum values as the **enum's spec** from now on.
With `go-enumer` there are two possible ways to define enum specs.
The first one we'll refer from now on as the *simple block spec*, which is straight forward and a suitable choice for small and uncomplicated enum sets.
The second one we'll refer to as the *filebased spec*, which is much more advanced in its options and useful when you're reaching the *simple block spec's* limits.

### Conventions

To define an enum spec, one first needs to define its distinct type – the **enum type**.
To define an Enum type, you must follow the listed rules.
Failing to do so will either lead to an unsuccessful detection or to a code generation error.

- Every enum type must derive from an unsigned integer (`uint`, `uint8`, ...) type.
- Every enum type must be marked with a comment directive `//go:enum` ([read here for more](#comment-directive-goenum)).
- Every enum type must either have a *simple block spec* or a *filebased spec* associated.

Not only the enum type, but also its spec must follow a set of basic rules.
By convention, enum sets must:

- start with **index** of `1` and in those cases with a **default value** at `0`.
- consist of continuous linear increments of at most `1`.
- contain unique **values**.

Due to the nature of enums being distinct values, the majority of enum sequences don't require a default value and thus start at `1`.
This is demonstrated in the following snippet with an enum type of a *simple block spec*.

```go
//go:enum -transform=upper
type Color uint

const (
  ColorRed Color = iota + 1
  ColorBlue
  ColorGreen
)
/* yields:
enum type: Color
enum spec:
Value "RED"   at Index 1
Value "BLUE"  at Index 2
Value "GREEN" at Index 3
is invalid    at Index 0 and 4,...,max
*/
```

### Defaults and Zero-Values

`go-enumer` gives a special semantic meaning to specs starting with the value `0`.
The **[Zero Value](https://go.dev/tour/basics/12)** of native integer types in Go is `0`.
For enums it likewise means, that an enum of value `0` is by definition the zero value.
If the enum spec does not define a **default** (starts at `0`) an enum field is in an invalid state when its assigned value is zero.
In some situations you may find yourself in the need of such a **default** value.
Depending on whether your spec needs a default value, you will chose your sequence to start at value `0` with the your default value being `0`.
(You may have multiple defaults, see section for [alternative values](#alternative-values))

```go
//go:enum
type UserRole uint

const (
  UserRoleAnonymous UserRole = iota // <- this is your default
  UserRoleStandard
  UserRoleAdmin
  UserRoleUnknown = UserRoleAnonymous // <- this is your alternative default
)
```

The previous snippet defines an enum spec which can be deserialized from the following set of values `["Anonymous","Standard","Admin","Unknown"]` after successful code generation.
If you want your default value to be robust against deserialization from `undefined` (resp. zero values), then rest assured you do not need to do anything.
`go-enum` will naturally fail any attempt of unmarshalling from empty strings or `nil` if it was not explicitly instructed to do otherwise.
In these cases the returned error will be of type `ErrNoValidEnum` which is part of the generated file and can be used via Go's unwrapping mechanism `errors.Is(err, mypkg.ErrNoValidEnum)`.

If you need to also enable deserialization for your **default** enum value from a zero values please check out the section for ["undefined"-value](#the-undefined-feature).

### Comment directive `//go:enum`

`go-enumer` only needs one single `//go:generate` directive per package to screen the entire package thanks to the introduction of `//go:enum` comment directive.
It acts as a marker of all enums.
Now `go-enumer` can determine all enums of a package in a single pass, reducing redundant scans and therefore making the code generation (and your CI process) much more efficient.

But the `//go:enum` directive does not exclusively serve as a marker to detect enums, it also enables config **mixins**.
By supplying it with a finegrained configuration on an enum type level, we have the ability to overwrite the global `generate` configuration.

The following example will generate `json` interfaces for practically all enums it will detect, except for this specific `Greeting` enum, where it will generate both `json` and `yaml` interfaces.

> The comment directive currently offers support for all configuration options, which are available on a global configuration level.

```go
package mypackage

//go:generate go run github.com/mvrahden/go-enumer -serializers=json

/* ... */

//go:enum -serializers=json,yaml
type Greeting uint

const (
  GreetingWorld Greeting = iota
  GreetingMars
)
```

### Validation

Each enum type will support two validation Methods `Validate() error` and `IsValid() bool`.

By defining enum specs on a linear scale, the validation complexity is reduced to a simple check of whether or not the value is within the extent of the scale.
To ensure (de-)serialization success the type validation will be performed on every (de-)serialization operation.

### Supported features

Supported features are targeted with the `-support=arg1,arg2,...` flag and can be used globally via `go:generate` or as a mixin via `go:enum`.

#### The "undefined" feature

> how to use? `-support=undefined`

`go-enumer` returns with an error upon (de-)serialization of any value which is not explicitly defined in your set of enums with an `error` – this rule covers empty values as well.
Therefore those sequences, that neither contain a default value nor an empty value must start with the integer value of `1`.

The `undefined` feature is an opt-in, which enables the (de-)serialization of zero values.
In case you applied the `undefined` feature, lookups with an empty value will resolve as an "undefined" constant, representing an empty string.
Any undefined value upon deserialization will be considered **valid** now and pass the [Validation](#validation).
However it will not be represented in the list of possible values upon serialization.

In the case of `undefined` values, there are 4 cases that we can distinguish:

| # | Has Default Value | Supports Undefined | Enum Type can deserialize from undefined |
|:-:|:-----------------:|:------------------:|:----------------------------------------:|
| 1 |        no         |         no         |                    no                    |
| 2 |        yes        |         no         |                    no                    |
| 3 |        yes        |        yes         |                   yes                    |
| 4 |        no         |        yes         |                   yes                    |

> Read the table as follows (e.g. for row `4`): "If my enum *has NO default* and it *DOES support undefined*
> then it *CAN deserialize from undefined*"

#### Other supported features

> how to use? `-support=ent`

With `ent` a method will be generated to return all valid Value strings. This allows you to use your enum type with the ent framework.

## Simple Block Spec

The simple block spec is a very primitive and intuitive way to generate enums.
Simply define a `const`-block with the enums you want your enum spec to contain.
The enum values will be derived from the constant names (less their type prefix).
To allow some degree of customization while keeping your constant names readable
you can apply [string case transformations](#string-case-transformations) to the values.

### String Case Transformations

`go-enumer` supports a range of string case transformations.
These transformations are a feature exclusive to the *simple block spec*.
You may configure a global transformation via the `generate` command and you may mixin deviating transformations case by case via `go:enum` comment directive.
Here is an example with a `kebab` transformation.

```go
//go:enum -transform=kebab
type MyType uint

const (
  MyTypeFoo MyType = iota + 1
  MyTypeBar
  MyTypeFooBar
  // ...
)
// yields --> ["foo", "bar", "foo-bar"]
```

Please take examplary transformations from the following table:

| Transformation Name | Resulting enum values     |
|---------------------|---------------------------|
| `noop` (default)    | ["Foo", "Bar", "FooBar"]  |
| `camel`             | ["foo", "bar", "fooBar"]  |
| `pascal`            | ["Foo", "Bar", "FooBar"]  |
| `kebab`             | ["foo", "bar", "foo-bar"] |
| `snake`             | ["foo", "bar", "foo_bar"] |
| `lower`             | ["foo", "bar", "foobar"]  |
| `upper`             | ["FOO", "BAR", "FOOBAR"]  |
| `upper-kebab`       | ["FOO", "BAR", "FOO-BAR"] |
| `upper-snake`       | ["FOO", "BAR", "FOO_BAR"] |
| `whitespace`        | ["Foo", "Bar", "Foo Bar"] |

**Note**: If you are missing a transformation please raise an issue and/or open a Pull Request.

### Handling of Name Prefixes

It is good practice with `go-enumer` to prefix enum constant names with their corresponding type name and `go-enumer` will automatically detect these prefixes and strip them off their values.
Meaning: it will turn `GreetingMars` (with enum type `Greeting`) into e.g. `"MARS"` (assuming an `upper` transformation was applied).
It does not support arbitrary prefixes and we do not encourage that due to a resulting noisyness of code.
This rule is in place to keep your source code concise.

### Alternative values

Enum based on the *simple block spec* can contain indeces (enum IDs) that are assigned multiple times – such values resemble **Alternative values**.
These alternative values are shadowed by the dominant value, which is always the constant which was assigned first to the index in the block.

## Filebased Spec

The *filebased spec* allows code generation for the enum values from a file source.
The following sections describe the usage of the *filebased spec*

### CSV-File sources

`go-enumer` can extract your enum definitions from CSV file sources if you target `-from=path/to.csv` in your enum's comment directive.
Filebased specs are taken as given and will not undergo any string case transformation.

Filebased specs allow you also to augment your enums with additional data columns.
`go-enumer` can parse data and add typed Getter-funcs based on a column annotation syntax.
It supports Go's built-in data types via the following syntax `<datatype>(your-column-name)`, e.g. `uint(area-in-square-meter)` or `float64(tolerance)`.
If there's no explicit type annotated, `go-enumer` will assume a basic `string` type as a fallback.

Have a look at [the Booking, Color or Project examples](examples/README.md) for further info.

## Generated functions and methods

When `go-enumer` is applied to a type, it will generate:

- The following basic methods/functions:

  - Method `String()`: returns the string representation of the enum value. This makes the enum conform
    to the `Stringer` interface, so whenever you print an enum value, you'll get the string name instead of its numeric counterpart.
  - Function `<EnumType>FromString(raw string)`: returns the enum value from its string representation. This is useful
    when you need to read enum values from command line arguments, from a configuration file, or
    from a REST API request... In short, from those places where using the real enum value (an integer) would
    be almost meaningless or hard to trace or use by a human. `raw` string is case sensitive.
  - Function `<EnumType>FromStringIgnoreCase(raw string)`: we can not always guarantee the case matching because some systems out of our reach
    are insensitive to exact case matching. In these situations `<EnumType>FromStringIgnoreCase(raw string)` comes in handy.
    It acts the same as `<EnumType>FromString(raw string)` with the little difference of `raw` being case insensitive.
  - Function `<EnumType>Values()`: returns a slice with all the numeric values of the enum, ignoring any alternative values.
  - Function `<EnumType>Strings()`: returns a slice with all the string representations of the enum.
  - Method `IsValid()`: returns true if the current value is a value of the defined enum set.
  - Method `Validate()`: returns a wrapped error `ErrNoValidEnum` if the current value is not a valid value of the defined enum set.
    It is being used upon serialization and deserialization, allowing for detecting enum errors via `errors.Is(err, ErrNoValidEnum)`.

- The flag `serializers` in addition with any of the following values, additional methods for serialization are added.
  Valid values are:

  - `binary` makes the enum conform to the `encoding.BinaryMarshaler` and `encoding.BinaryUnmarshaler` interfaces.
  - `bson` makes the enum conform to the `bson.MarshalBSONValue` and `bson.UnmarshalBSONValue` interfaces.
  - `graphql` makes the enum conform to the `graphql.Marshaler` and `graphql.Unmarshaler` interfaces.
  - `json` makes the enum conform to the `json.Marshaler` and `json.Unmarshaler` interfaces.
  - `sql` makes the enum conform to the `sql.Scanner` and `sql.Valuer` interfaces.
  - `text` makes the enum conform to the `encoding.TextMarshaler` and `encoding.TextUnmarshaler` interfaces.
    **Note:** If you use your enum values as keys in a map and you encode the map as *JSON*,
    you need this flag set to true to properly convert the map keys to json (strings). If not, the numeric values will be used instead
  - `yaml` makes the enum conform to the `gopkg.in/yaml.v2.Marshaler` and `gopkg.in/yaml.v2.Unmarshaler` interfaces.
  - `yaml.v3` makes the enum conform to the `gopkg.in/yaml.v3.Marshaler` and `gopkg.in/yaml.v3.Unmarshaler` interfaces.
    **Note:** Supplying both yaml values (`yaml` and `yaml.v3`) will fail due to interface incompatibility.

## Configuration Options

You can add:

- transformation with `transform` option, e.g. `transform=kebab`.
- serializers with `serializers` option, e.g. `serializers=json,sql,...`.
- supported features via `support` option, e.g. `support=undefined,ent`
  - `undefined`, see ["undefined"-value](#the-undefined-feature)
  - `ignore-case`, adds support for case-insensitive lookup
  - `ent`, adds interface support for [entgo.io](https://github.com/ent/ent)

## Caveats

Following is a list of known issues:

- **Skipped calls to `UnmarshalJSON`**:  
  When using `encoding/json` the enum unmarshaling will not be triggered if there's no `key`-`value` pair for the enum within the JSON payload.
  This will lead to a zero value enum instead of a deserialized enum.
  If no subsequent validation is performed and no default value is defined this will cause a failing validation upon subsequent serialization.

## Inspiring projects

- [enumer](https://github.com/dmarkham/enumer)
