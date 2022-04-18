# `go-enumer` Examples

The examples directory contains examples with specific use cases.
Each example is self-contained and part of a wider testing setup to improve QA and limit project overhead.

## Example Use Cases

<!--
TODO:
- add all transform cases
- add ignore case example
 -->

1. `greetings`: Generate standard enum and enums with default value (zero value).
2. `pills`: Generate enums for all unsigned integer types.
3. `planets`: Generate various combinations of standard/default vs. undefined.
4. `booking`: Generate enums from CSV source.
5. `project`: A more realistic mix of enums.

> `_invalid`: Contains various invalid edge cases which are expected to produce specific user-friendly errors.
> You can happily **ignore this directory** as it is for testing puproses only.

## Tests

As mentioned above, these examples are part of a wider testing setup.
Each example directory contains a specific set of files:

- `enums.go`: contains the enum type definitions.
- `generated.go`: contains the expected generated output.
- `enums_test.go`: contains test files to assert that the `generated.go` file performs as expected.
- `config.yml`: the base configuration for the code generation **run**.
