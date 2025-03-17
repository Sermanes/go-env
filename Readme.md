# Environment Variables Package
## Overview
The environment variables package provides a set of functions to retrieve environment variables in Go, handling various data types and providing default values if the variables are not set.

## Features
- **Type-Specific Retrieval**: Functions to retrieve environment variables as strings, unsigned integers (`uint`, `uint8`), booleans, byte slices, string slices, floating-point numbers (`float64`), and signed integers (`int64`).
- **Default Values**: Each function allows specifying a default value to return if the environment variable is not set.
- **Error Handling**: Automatically returns default values if parsing fails due to incorrect data type.

## Constants
- **Default Slice Separator**: The package uses a comma (`,`) as the default separator for string slices.
- **Bit Size for 64-bit Numbers**: Used for parsing `float64` and `int64` values.
- **Decimal Base**: Used for parsing integers.

## Usage
### Retrieving Environment Variables

#### String
Get a string environment variable with a default value.
```go
value := env.Get("KEY", "default_value")
```


#### Unsigned Integer (`uint`)
Get an unsigned integer environment variable with a default value.
```go
value := env.GetUint("KEY", 0)
```


#### Unsigned 8-bit Integer (`uint8`)
Get an unsigned 8-bit integer environment variable with a default value.
```go
value := env.GetUint8("KEY", 0)
```


#### Boolean
Get a boolean environment variable with a default value.
```go
value := env.GetBool("KEY", false)
```


#### Byte Slice (`[]byte`)
Get a byte slice environment variable with a default value.
```go
value := env.GetBytes("KEY", []byte{})
```


#### String Slice (`[]string`)
Get a string slice environment variable with a default value. Values are split using a comma (`,`).
```go
value := env.GetStringSlice("KEY", []string{})
```

#### Floating-Point Number (`float64`)
Get a floating-point number environment variable with a default value.
```go
value := env.GetFloat64("KEY", 0.0)
```

#### Signed 64-bit Integer (`int64`)
Get a signed 64-bit integer environment variable with a default value.
```go
value := env.GetInt64("KEY", 0)
```

## Contributing
Contributions are welcome. Please submit pull requests with detailed explanations of changes.