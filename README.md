# AVA-Stream API

Lightweight and super fast functional stream processing library inspired by [Java Streams API](https://docs.oracle.com/javase/8/docs/api/java/util/stream/Stream.html) and [Lodash](https://lodash.com).

## Table of contents

- [Requirements](#requirements)
- [Usage examples](#usage-examples)
- [Limitations](#limitations)
- [Performance](#performance)
- [Completion status](#completion-status)
- [Extra credits](#extra-credits)

## Requirements

- Go 1.18 or higher

This library makes intensive usage of [Type Parameters (generics)](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md) so it is not compatible with any Go version lower than 1.18.

## Usage examples

For more details about the API, please check the `stream/*_test.go` files.
