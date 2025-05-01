# Test Finder

Finds the test.

## Args

Point to a test file.

## Flags

| Flag      | Default | Description                |
| --------- | ------- | -------------------------- |
| subtests  | `true`  | Emits sub tests to output  |
| fulltests | `true`  | Emits full tests to output |

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
