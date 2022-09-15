package session

import (
	"Current/Gorm/geeORM/schema"
	"fmt"
	"github.com/golang/glog"
	"reflect"
	"strings"
)

// SetSchemaByInstance 给schema——数据库表格赋值，解析操作比较耗时，将解析结果保存到成员变量schema
// 即使schema被调用多次，如果传入结构体名称不发生辩护啊，则不会更新schema的值
func (s *Session) SetSchemaByInstance(instance interface{}) *Session {
	// nil or different model,update schema
	if s.schema == nil || reflect.TypeOf(instance) != reflect.TypeOf(s.schema.Model) {
		s.schema = schema.ParseStructInstanceToSchema(instance, s.dialect)
	}
	return s
}

func (s *Session) GetSchema() *schema.Schema {
	if s.schema == nil {
		glog.Error("[Session.GetSchema] Schema is nil")
	}
	return s.schema
}

// CreateTable 数据库表的创建
func (s *Session) CreateTable() error {
	table := s.schema
	var columns []string
	for _, filed := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", filed.Name, filed.Type, filed.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s);", table.Name, desc)).Exec()
	return err
}

// DropTable 数据库表的删除
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s;", s.GetSchema().Name)).Exec()
	return err
}

// HasTable 判断数据库表是否存在
func (s *Session) HasTable() bool {
	sql, values := s.dialect.TableExistSQL(s.GetSchema().Name)
	row := s.Raw(sql, values...).QueryRaw() //指针*sql.Row
	var tmp string
	_ = row.Scan(&tmp)               //获取查询后返回的指针值
	return tmp == s.GetSchema().Name //判断是否存在并返回
}
