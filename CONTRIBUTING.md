# Contributing

## Linter

In order to contribute to our project, you need to make sure to install [golangci-lint](https://golangci-lint.run/usage/install/#binaries).\
Run it before committing.

Then, configure git hooks:
```bash
git config core.hooksPath .githooks
```

## Types

In `golang`, it's pretty common to put vars and types at the beginning of the file.\
By the time, the file keep growing and at some point, you will see in some projects those definitions can be found in the middle of nowhere. Let's avoid that pattern and put everything in scoped type files.\
As an example for rafty scope, it will be `rafty_types.go`.\
This will make this project more consistant by providing a better clarity of what is being done.

## Unit testing

It's very common in golang to have 2 files:
```go
// dummy.go
func Foo(){}
func Bar(){}
```
and the test file:
```go
// dummy_test.go
TestFoo(xxx)
TestBar(xxx)
```

That's very great to have this pattern. The problem is when the project is getting very big and you want to execute a set of specific tests, this won't be possible. It's avoid that by scoping tests like so:
```go
// dummy_test.go
TestDummy_Foo(xxx)
TestDummy_Bar(xxx)
```
With this scoped tests definitions, we will be able to run `go test -v -race -run Dummy`.\
Privilegiate `t.Run()` when it's possible with a meaningful name that represent what is being tested.\
Don't use name with space like `t.Run("my dummy test", func(t *testing.T) {})`.\
Use this pattern `t.Run("my_dummy_test", func(t *testing.T) {})`.

It's the same for mocks except that the pattern can be slighly different. It's better to group them all by having `mock_<scope>_test.go` or `mocks_<scope>_test.go` but you can also have `<scope>_mocks_test.go`.

## Commits

Use coventional commits overwise pull request will be rejected. See [conventional commits here](https://www.conventionalcommits.org/en/v1.0.0/).

About commit, try to keep the summary short but meaningful BUT add more in the description for more explinations.\
Example: `git commit -m "feat(cluster): My meaningful summary" -m "My meaningful description"`.
