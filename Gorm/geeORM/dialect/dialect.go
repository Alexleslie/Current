package dialect

import (
	"reflect"
)

// 管理多个数据库语言
var dialectsMap = map[string]Dialect{}

// Dialect 将Go语言转换为数据库语言的接口
type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    // 将Go语言的类型转换为该数据库的数据类型
	TableExistSQL(tableName string) (string, []interface{}) // 返回某个表是否存在的SQL语句，参数是表名（tableName）
}

// RegisterDialect 注册dialect实例
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect 获取dialect实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
