# Environment Variables Package
## Overview
The environment variables package provides a set of functions to retrieve environment variables in Go, handling various data types and providing default values if the variables are not set.

## Features
- **Type-Specific Retrieval**: Functions to retrieve environment variables as strings, booleans, byte slices, string slices, signed integers (`int`, `int32`, `int64`), unsigned integers (`uint`, `uint8`, `uint16`, `uint32`, `uint64`), floating-point numbers (`float32`, `float64`), durations (`time.Duration`), absolute URLs (`*url.URL`), and timestamps (`time.Time`).
- **Default Values**: Each function allows specifying a default value to return if the environment variable is not set.
- **Error Handling**: Automatically returns default values if parsing fails due to incorrect data type.
- **Overflow Safety**: Sized unsigned integers reject negatives and out-of-range values (e.g. `GetUint8` rejects `>255` and negatives instead of silently wrapping).
- **Custom Separators & Trimming**: `GetStringSliceWithSep` accepts any separator and trims surrounding whitespace from each item.
- **Error-Returning Variants (`GetXE`)**: Same as the silent accessors, but return `(value, error)` so callers can detect parse failures instead of silently using the default.
- **Required Variants (`MustGetX`)**: No default — panic with `ErrMissing` when the variable is unset, or with `*ParseError` when the value cannot be parsed. Use for mandatory startup config.
- **URL & Time Accessors**: `GetURL` (requires absolute URL with scheme + host), `GetTime` (RFC3339), and `GetTimeWithLayout` (custom Go time layout).

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
Get an unsigned 8-bit integer environment variable with a default value. Rejects negatives and values above 255 (returns the default instead of wrapping).
```go
value := env.GetUint8("KEY", 0)
```


#### Unsigned 16-bit Integer (`uint16`)
Get an unsigned 16-bit integer environment variable with a default value. Rejects negatives and values above 65535.
```go
value := env.GetUint16("KEY", 0)
```


#### Unsigned 32-bit Integer (`uint32`)
Get an unsigned 32-bit integer environment variable with a default value. Rejects negatives and values above `math.MaxUint32`.
```go
value := env.GetUint32("KEY", 0)
```


#### Unsigned 64-bit Integer (`uint64`)
Get an unsigned 64-bit integer environment variable with a default value. Rejects negatives.
```go
value := env.GetUint64("KEY", 0)
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


#### String Slice with Separator (`[]string`)
Get a string slice environment variable with a default value using a custom separator. Each item has surrounding whitespace trimmed.
```go
// "a, b ,c" → []string{"a", "b", "c"}
value := env.GetStringSliceWithSep("KEY", []string{}, ",")

// "a;b;c" → []string{"a", "b", "c"}
value := env.GetStringSliceWithSep("KEY", []string{}, ";")
```

#### Floating-Point Number (`float32`)
Get a 32-bit floating-point number environment variable with a default value.
```go
value := env.GetFloat32("KEY", 0.0)
```

#### Floating-Point Number (`float64`)
Get a floating-point number environment variable with a default value.
```go
value := env.GetFloat64("KEY", 0.0)
```

#### Signed 32-bit Integer (`int32`)
Get a signed 32-bit integer environment variable with a default value. Returns the default if the value overflows int32.
```go
value := env.GetInt32("KEY", 0)
```

#### Signed 64-bit Integer (`int64`)
Get a signed 64-bit integer environment variable with a default value.
```go
value := env.GetInt64("KEY", 0)
```

#### Signed Integers (`int`)
Get a signed integer environment variable with a default value.
```go
value := env.GetInt("KEY", 0)
```

#### Duration (`time.Duration`)
Get a duration environment variable with a default value. Accepts any string parseable by `time.ParseDuration` (e.g. `"5s"`, `"1h30m"`, `"500ms"`).
```go
value := env.GetDuration("TIMEOUT", 30*time.Second)
```

#### URL (`*url.URL`)
Get an absolute URL environment variable with a default value. The value must include a scheme and host; relative paths and bare paths are rejected.
```go
defaultAPI, _ := url.Parse("https://api.example.com")
api := env.GetURL("API_BASE_URL", defaultAPI)
```

#### Time (`time.Time`) — RFC3339
Get an RFC3339 timestamp environment variable with a default value.
```go
deadline := env.GetTime("DEPLOY_DEADLINE", time.Time{})
// DEPLOY_DEADLINE="2026-05-10T12:30:00Z" → time.Time
```

#### Time (`time.Time`) — Custom Layout
Get a timestamp environment variable parsed with a custom Go time layout.
```go
date := env.GetTimeWithLayout("RELEASE_DATE", time.Time{}, "2006-01-02")
// RELEASE_DATE="2026-05-10" → time.Time
```

## Error-Returning Variants

Every `GetX` accessor has a sibling `GetXE` that returns `(value, error)`:

- Variable unset or empty (for slice accessors) → `(defaultValue, nil)`.
- Variable set but parsing fails → `(defaultValue, *ParseError)`.

```go
port, err := env.GetIntE("PORT", 8080)
if err != nil {
    log.Printf("invalid PORT: %v", err)
}

// Inspect the failure
var perr *env.ParseError
if errors.As(err, &perr) {
    fmt.Println(perr.Key, perr.Value, perr.Unwrap())
}
```

## Required Variants (Must)

For mandatory startup configuration, use `MustGetX`. They take no default and panic when the variable is missing or invalid — surfacing config errors at boot instead of running with silent defaults.

```go
port := env.MustGetInt("PORT")           // panic with *fmt.wrapError if unset, *ParseError if invalid
host := env.MustGet("DATABASE_HOST")     // panic if unset
timeout := env.MustGetDuration("TIMEOUT")
```

`errors.Is(err, env.ErrMissing)` and `errors.As(err, &*env.ParseError{})` work on the recovered values.

## Contributing
Contributions are welcome. Please submit pull requests with detailed explanations of changes.