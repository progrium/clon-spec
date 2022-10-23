# Command-Line Object Notation
Ergonomic JSON-compatible input syntax for CLI tools.

CLON is an argument syntax spec for defining JSON-like objects in a way that makes sense for the command-line that isn't just taking raw JSON or pointing to a file. This was shamelessly stolen from [HTTPie](https://httpie.io/docs/cli/json) with some modification.

CLON is not meant to be a tool itself, but an argument syntax that you can use when building CLI tools that need to take arbitrary JSON-like data. An example tool would be a CLI client for an RPC system with a usage format like `rpctool call <method> <args...>`. CLON makes for a reasonably human-friendly way to input the argument data. Note that internally the tool may not use this syntax to build JSON, but a compatible format like CBOR, or even just to build an in-memory structure.

### Implementations

Below are some language specific libraries for easily accepting this syntax in your command-line tools:

 * [progrium/clon-go](https://github.com/progrium/clon-go)
 * your library here

## Simple Example
In these examples, `clon` is an *imaginary tool* that takes CLON arguments and outputs the resulting structure as JSON.

```shell
$ clon name=John email=john@example.org
```
```json
{
    "name": "John",
    "email": "john@example.org"
}
```

If the first argument is in key-value form (with `=` or `:=` operator), the arguments will be used to construct an object. If the first argument is simply a value, the arguments will be used to construct an array. 

```shell
$ clon John john@example.org "Text string"
```
```json
[
  "John",
  "john@example.org",
  "Text string"
]
```

## Non-string JSON fields
Non-string JSON fields use the `:=` separator, which allows you to embed arbitrary JSON data into the resulting JSON object. Additionally, text and raw JSON files can also be embedded into fields using `=@` and `:=@`:

```shell
$ clon \
    name=John \                        # String (default)
    age:=29 \                          # Raw JSON — Number
    married:=false \                   # Raw JSON — Boolean
    hobbies:='["http", "pies"]' \      # Raw JSON — Array
    favorite:='{"tool": "HTTPie"}' \   # Raw JSON — Object
    bookmarks:=@files/data.json \      # Embed JSON file
    description=@files/text.txt        # Embed text file
```
```json
{
    "age": 29,
    "hobbies": [
        "http",
        "pies"
    ],
    "description": "John is a nice guy who likes pies.",
    "married": false,
    "name": "John",
    "favorite": {
        "tool": "HTTPie"
    },
    "bookmarks": {
        "HTTPie": "https://httpie.org",
    }
}
```

When building a top-level array, values can be prefixed with `:` to indicate a non-string or raw JSON value. This is a shorthand for the top-level array syntax described later.

```shell
$ clon John :29 :false :'{"tool": "HTTPie"}'
```
```json
[
  "John",
  29,
  false,
  {
    "tool": "HTTPie"
  }
]
```

## Nested JSON
You still use the existing data field operators (`=`/`:=`) but instead of specifying a top-level field name (like `key=value`), you specify a path declaration telling where and how to put the given value inside an object:

```shell
$ clon \
    platform[name]=HTTPie \
    platform[about][mission]='Make APIs simple and intuitive' \
    platform[about][homepage]=httpie.io \
    platform[about][stars]:=54000 \
    platform[apps][]=Terminal \
    platform[apps][]=Desktop \
    platform[apps][]=Web \
    platform[apps][]=Mobile
```
```json
{
    "platform": {
        "name": "HTTPie",
        "about": {
            "mission": "Make APIs simple and intuitive",
            "homepage": "httpie.io",
            "stars": 54000
        },
        "apps": [
            "Terminal",
            "Desktop",
            "Web",
            "Mobile"
        ]
    }
}
```

### Basic Usage

Let’s start with a simple example, and build a simple search query:

```shell
$ clon \
    category=tools \
    search[type]=id \
    search[id]:=1
```
```json
{
    "category": "tools",
    "search": {
        "id": 1,
        "type": "id"
    }
}
```

In the example above, the `search[type]` is an instruction for creating an object called `search`, and setting the `type` field of it to the given value (`"id"`).

Also note that, just as the regular syntax, you can use the `:=` operator to directly pass raw JSON values (e.g, numbers in the case above).

Building arrays is also possible, through `[]` suffix (an append operation). This creates an array in the given path (if there is not one already), and append the given value to that array.

```shell
$ clon \
    category=tools \
    search[type]=keyword \
    search[keywords][]=APIs \
    search[keywords][]=CLI
```
```json
{
    "category": "tools",
    "search": {
        "keywords": [
            "APIs",
            "CLI"
        ],
        "type": "keyword"
    }
}
```

If you want to explicitly specify the position of elements inside an array, you can simply pass the desired index as the path:

```shell
$ clon \
    category=tools \
    search[type]=keyword \
    search[keywords][1]=APIs \
    search[keywords][0]=CLI
```
```json
{
    "category": "tools",
    "search": {
        "keywords": [
            "CLIs",
            "API"
        ],
        "type": "keyword"
    }
}
```

If there are any missing indexes, they will be nullified in order to create a concrete object:

```shell
$ clon \
    category=tools \
    search[type]=platforms \
    search[platforms][]=Terminal \
    search[platforms][1]=Desktop \
    search[platforms][3]=Mobile
```
```json
{
    "category": "tools",
    "search": {
        "platforms": [
            "Terminal",
            "Desktop",
            null,
            "Mobile"
        ],
        "type": "platforms"
    }
}
```

It is also possible to embed raw JSON to a nested structure, for example:

```shell
$ clon \
  category=tools \
  search[type]=platforms \
  'search[platforms]:=["Terminal", "Desktop"]' \
  search[platforms][]=Web \
  search[platforms][]=Mobile
```
```json
{
    "category": "tools",
    "search": {
        "platforms": [
            "Terminal",
            "Desktop",
            "Web",
            "Mobile"
        ],
        "type": "platforms"
    }
}
```

And just to demonstrate all of these features together, let’s create a very deeply nested JSON object:

```shell
$ clon \
    shallow=value \                                # Shallow key-value pair
    object[key]=value \                            # Nested key-value pair
    array[]:=1 \                                   # Array — first item
    array[1]:=2 \                                  # Array — second item
    array[2]:=3 \                                  # Array — append (third item)
    very[nested][json][3][httpie][power][]=Amaze   # Nested object
```

### Advanced Usage

#### Top level arrays

To build an array instead of a regular object, you can simply do that by omitting the starting key:

```shell
$ clon \
    []:=1 \
    []:=2 \
    []:=3
```
```json
[
    1,
    2,
    3
]
```

You can also apply the nesting to the items by referencing their index:

```shell
$ clon \
    [0][type]=platform [0][name]=terminal \
    [1][type]=platform [1][name]=desktop
```
```json
[
    {
        "type": "platform",
        "name": "terminal"
    },
    {
        "type": "platform",
        "name": "desktop"
    }
]
```

#### Escaping behavior
Nested JSON syntax uses the same escaping rules as the terminal. There are 3 special characters, and 1 special token that you can escape.

To include a bracket as is, escape it with a backslash (`\`):

```shell
$ clon \
  'foo\[bar\]:=1' \
  'baz[\[]:=2' \
  'baz[\]]:=3'
```
```json
{
    "baz": {
        "[": 2,
        "]": 3
    },
    "foo[bar]": 1
}
```

If use the literal backslash character (`\`), escape it with another backslash:

```shell
$ clon \
  'backslash[\\]:=1'
```
```json
{
    "backslash": {
        "\\": 1
    }
}
```

A regular integer in a path (e.g `[10]`) means an array index; but if you want it to be treated as a string, you can escape the whole number by using a backslash (`\`) prefix.

```shell
$ clon \
  'object[\1]=stringified' \
  'object[\100]=same' \
  'array[1]=indexified'
```
```json
{
    "array": [
        null,
        "indexified"
    ],
    "object": {
        "1": "stringified",
        "100": "same"
    }
}
```

#### Guiding syntax errors

If you make a typo or forget to close a bracket, the errors SHOULD guide you to fix it. For example:

```
$ clon \
  'foo[bar]=OK' \
  'foo[baz][quux=FAIL'
Syntax Error: Expecting ']'
foo[baz][quux
             ^
```

You can follow to given instruction (adding a `]`) and repair your expression.

#### Type safety

Each container path (e.g., `x[y][z]` in `x[y][z][1]`) has a certain type, which gets defined with the first usage and can’t be changed after that. If you try to do a key-based access to an array or an index-based access to an object, you should get an error out:

```
$ clon \
  'array[]:=1' \
  'array[]:=2' \
  'array[key]:=3'
Type Error: Can't perform 'key' based access on 'array' which has a type of 'array' but this operation requires a type of 'object'.
array[key]
     ^^^^^
```

Type Safety does not apply to value overrides, for example:

```shell
$ clon \
  user[name]:=411     # Defined as an integer
  user[name]=string   # Overridden with a string
```
```json
{
    "user": {
        "name": "string"
    }
}
```

## Raw JSON

For very complex JSON structures, it may be more convenient to pass it as raw input, for example:

```shell
$ echo -n '{"hello": "world"}' | clon
```

```shell
$ clon < files/data.json
```

