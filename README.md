# deepcopy
library to make deep copies in go

```Go
   newValue := deepcopy.Do(value)
```

It manages when struct have pointers to parents/children. e.g:
```Go
   type MyStruct struct {
        Parent *MyStruct
        Children []*MyStruct
   }
```

Unexported fields are not copied.