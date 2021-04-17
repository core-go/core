# Common Service
Some standard interfaces:
- StringService
- ViewService: to get all data, to get data by identity
- GenericService: to insert, update, save, delete data

## Installation

Please make sure to initialize a Go module before installing common-go/service:

```shell
go get -u github.com/common-go/service
```

Import:

```go
import "github.com/common-go/service"
```

You can optimize the import by version:
- v0.0.2: StringService
- v0.0.4: ViewService only
- v1.0.5: ViewService, GenericService

## Details:
#### string_service.go
```go
type StringService interface {
	Load(ctx context.Context, key string, max int64) ([]string, error)
	Save(ctx context.Context, values []string) (int64, error)
	Delete(ctx context.Context, values []string) (int64, error)
}
```

#### view_service.go
```go
type ViewService interface {
	IdNames() []string
	All(ctx context.Context) (interface{}, error)
	Load(ctx context.Context, id interface{}) (interface{}, error)
	LoadAndDecode(ctx context.Context, id interface{}, result interface{}) (bool, error)
	Exist(ctx context.Context, id interface{}) (bool, error)
}
```

#### generic_service.go
```go
type GenericService interface {
	ViewService
	Insert(ctx context.Context, model interface{}) (int64, error)
	Update(ctx context.Context, model interface{}) (int64, error)
	Patch(ctx context.Context, model map[string]interface{}) (int64, error)
	Delete(ctx context.Context, id interface{}) (int64, error)
}
```
