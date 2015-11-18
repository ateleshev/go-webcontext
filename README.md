# go-webcontext

```
go get -u github.com/ArtemTeleshev/go-webcontext
```

```
import wc "github.com/ArtemTeleshev/go-webcontext"
```

Usage:
```
var (
  context *wc.Context
)

func init() {
  context = wc.NewContext()
  if db, err := gorm.Open("mysql", "username:password@tcp(localhost:3306)/dbname?charset=utf8"); err != nil {
    panic("Cannot connect to DB")
  }
  context.SetDB(db)
}

func main() {
  var cnt int
  db := context.GetDB()
  db.Table("test").Count(&cnt)
  log.Printf("Total rows in table 'test': %d", cnt)
}
```
