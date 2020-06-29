# Common Service
Some standard interfaces:
- StringService
- ViewService: to get all data, to get data by identity
- GenericService: to insert, update, save, delete data
- DiffService: to check differences between 2 models
- ApprService: to approve or reject the changes
- ModelBuilder, which is using IdGenerator
- UniqueValueBuilder

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
- v0.0.3: UniqueValueBuilder only
- v0.0.4: ViewService only
- v1.0.5: ViewService, GenericService, DiffService, ApprService
- v1.0.7: UniqueValueBuilder, Loader, Generator, DiffService and ApprService
- v1.1.0: DefaultUniqueValueBuilder

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
	Save(ctx context.Context, model interface{}) (int64, error)
	Delete(ctx context.Context, id interface{}) (int64, error)
}
```

#### diff_service.go
```go
type DiffService interface {
	Diff(ctx context.Context, id interface{}) (*DiffModel, error)
}
```

#### appr_service.go
```go
type ApprService interface {
	Approve(ctx context.Context, id interface{}) (StatusCode, error)
	Reject(ctx context.Context, id interface{}) (StatusCode, error)
}
```

#### unique_value_builder.go
```go
type UniqueValueBuilder interface {
	Build(ctx context.Context, model interface{}, name string) (string, error)
}
```
