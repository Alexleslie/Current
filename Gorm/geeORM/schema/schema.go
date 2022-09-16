package schema

import (
	"Current/Gorm/geeORM/dialect"
	"Current/tools/logc"
	"go/ast"
	"reflect"
)

/*
Go转换为数据库中对应的表格类型——字段、表格，以及SQL插入语句
*/

// Field 字段，表示数据库的列
// Field represents a column of database
type Field struct {
	Name string // 列名,字段名
	Type string // 类型
	Tag  string // 约束条件，如非空、主键等
}

// Schema 表示数据库的表格
// Schema represents a table of database
type Schema struct {
	Model      interface{}       // 被映射的对象
	Name       string            // 表名
	Fields     []*Field          // 字段表
	FieldNames []string          // 字段名表，包含所有字段名（列名）
	fieldMap   map[string]*Field // 字典，通过列名和列实例映射管理表中的列，方便以后使用，无需遍历Fields
}

// GetField 通过列名获取列
func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// ParseStructInstanceToSchema 为任意结构体解析为Schema实例（数据库用）
func ParseStructInstanceToSchema(dest interface{}, d dialect.Dialect) *Schema {
	// TypeOf()和ValueOf()是reflect包分别同来返回入参的类型和值，设计的入参是一个对象的指针
	// 因此需要reflect.Indirect()获取指针指向的实例。
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(), // 获取到结构体的名称作为表名
		fieldMap: make(map[string]*Field),
	}

	// NumField()获取实例字段的个数
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)                     // 获取特定字段
		if !p.Anonymous && ast.IsExported(p.Name) { // 是否是嵌入字段&&是否是大写字母开头
			filed := &Field{
				Name: p.Name,                                              // 字段名
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))), // 转换为数据库的字段类型
			}
			if v, ok := p.Tag.Lookup("primary_key"); ok {
				filed.Tag = v
			}
			schema.Fields = append(schema.Fields, filed)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = filed
		}
	}
	return schema
}

// RecordObjectValues 根据Schema的字段值，记录所有实例所有的字段值
func (schema *Schema) RecordObjectValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	logc.Info("[Schema.RecordObjectValues] Object=[%+v],fieldValues=[%+v]", dest, fieldValues)
	return fieldValues
}
