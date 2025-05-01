# Test Finder

Simple utility to find tests within Go files, to power a VSCode extension to
bring better Run functionality to test files to mimic GoLand behavior.

Expect extremely limited support, as this is just a hodgepodge tool for myself
in the few hours I get to code for myself. With that, know that it might not
follow best practices with tooling.

## Args

Point to a test file.

## Flags

| Flag      | Default | Description                                 |
| --------- | ------- | ------------------------------------------- |
| subtests  | `true`  | Emits sub tests to output                   |
| fulltests | `true`  | Emits full tests to output                  |
| ignorelit | `false` | Ignores `t.Run` seeded with string literals |

## Output

For a given file:

- Test function
- Sub Test
- Line

```json
[
  {
    "testFunc": "TestSplitMapTest",
    "subTest": "simple",
    "line": 16
  },
  {
    "testFunc": "TestSplitMapTest",
    "subTest": "wrong sep",
    "line": 21
  },
  {
    "testFunc": "TestSplitMapTest",
    "subTest": "no sep",
    "line": 26
  },
  {
    "testFunc": "TestSplitMapTest",
    "subTest": "trailing sep",
    "line": 31
  }
]
```

# Contributions

More than happy to accept contributions, however since this just a simple tool
to do one job do not expect much. I am more than happy to review suggestions or
refactor points of failures, or even add edge cases, but since this was made in
an evening do not expect much.

If all else, feel free to create a fork.

# License

```
Copyright 2025 Deadly Surgeon

Licensed under the Apache License, Version 2.0 (the "License");
you may not use these files except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
