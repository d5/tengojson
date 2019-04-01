# tengojson

Scriptable JSON validation/transformation using Tengo

### Accessing 

```golang
tengojson.New().
    On(".", `func(v) { v.foo = "bar" }`).
    Run(`{}`)
```

### Transforming a value

```golang
func(v) {
    return string(v) // convert any 'v' to string value
}
```

### Adding a new key

```golang
func(v) {
    v.new_key = "value"
}
```

### Deleting a key

```golang
func(v) {
    v.new_key = undefined
}
```

### Validation

Return an `error` will stop the pipeline and return an error back to the user.

```golang
func(v) {
    if !is_string(v) { return error("must be string") } 
}
``` 
