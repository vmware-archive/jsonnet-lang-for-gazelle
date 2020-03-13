(import 'a.jsonnet') +
(import 'b.jsonnet') +
(import 'c.json') {
    text: importstr 'data.json',
    // d should not be rendered as dep in the BUILD.bazel file
    // d: import 'd.jsonnet',
} +
(import 'e.jsonnet')
