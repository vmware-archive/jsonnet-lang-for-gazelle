local a = import 'db/data.jsonnet';

(import 'lib/load.libsonnet') + a {
  a+: 1,
}
