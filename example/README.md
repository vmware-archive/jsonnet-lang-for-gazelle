# Example

run:

```
$ bazel run //:gazelle
$ bazel build //:demo_to_json
```

Once run, it will perform some changes to the checked out BUILD.bazel files; please don't commit those changes
as the current checked in state is a suitable base for trying out what would happen if you'd run gazelle.
