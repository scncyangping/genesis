

### 测试及测试覆盖率

```
GOROOT=$GOROOT GOPATH=$GOPATH $GOROOT/bin/go test -coverprofile=coverage.out  ./mysqlRepo && $GOROOT/bin/go tool cover -html=coverage.out
```