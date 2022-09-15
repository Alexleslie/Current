package dialect

import (
	"fmt"
	"reflect"
	"time"
)

// 将Go语言转换为sqlite语言的具体操作
type sqlite3 struct{}

var _ Dialect = (*sqlite3)(nil)

// 包在第一次加载时，会将sqlite3的dialect自动注册到全局
func init() {
	RegisterDialect("sqlite3", &sqlite3{})
}

// DataTypeOf 将Go语言变量类型转换为sqlite3数据库对应类型
func (s *sqlite3) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "bigint"
	case reflect.Float32, reflect.Float64:
		return "real"
	case reflect.String:
		return "text"
	case reflect.Array, reflect.Slice:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("[sqlite3.DataTypeOf] Invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

// TableExistSQL 返回查询表是否存在的sqlite3语言
func (s *sqlite3) TableExistSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name =?", args
}
