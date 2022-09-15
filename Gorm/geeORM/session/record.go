package session

import (
	"Current/Gorm/geeORM/clause"
	"errors"
	"reflect"
)

/*记录增删查改相关的代码*/

// Insert 实现一组结构体实例的插入数据库中的表
// 首先set需要的参数，然后再整体build语句
func (s *Session) Insert(instances ...interface{}) (int64, error) {
	recordInstances := make([]interface{}, 0)
	for _, instance := range instances {
		s.CallMethod(BeforeInsert, instance)
		table := s.SetSchemaByInstance(instance).GetSchema()
		s.clause.SetSqlAndVars(clause.INSERT, table.Name, table.FieldNames) //构造SQL插入子句
		recordInstances = append(recordInstances, table.RecordObjectValues(instance))
	}

	s.clause.SetSqlAndVars(clause.VALUES, recordInstances)           // 构造SQL插入值的子句
	sql, vars := s.clause.BuildInOrder(clause.INSERT, clause.VALUES) // 按顺序将子句构造成SQL完整语句
	result, err := s.Raw(sql, vars...).Exec()                        // 执行语句
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterInsert, nil)
	return result.RowsAffected()
}

// FindAll 根据结构体找到所有符合的数据
func (s *Session) FindAll(emptyInstances interface{}) error {
	s.CallMethod(BeforeQuery, nil)
	instSlice := reflect.Indirect(reflect.ValueOf(emptyInstances))
	structType := instSlice.Type().Elem() // 获取切片的单个元素的类型
	// 实例化结构体，解析成schema
	table := s.SetSchemaByInstance(reflect.New(structType).Elem().Interface()).GetSchema()

	s.clause.SetSqlAndVars(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.BuildInOrder(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRaws()
	if err != nil {
		return err
	}

	for rows.Next() {
		// 创建一个结构体实例，用于存储数据
		destInst := reflect.New(structType).Elem()
		var instFieldsAttrs []interface{}
		for _, name := range table.FieldNames {
			instFieldsAttrs = append(instFieldsAttrs, destInst.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(instFieldsAttrs...); err != nil {
			return err
		}
		s.CallMethod(AfterQuery, destInst.Addr().Interface())
		// 将dest添加到切片destSlice中。循环直到所有的记录都添加到切片destSlice中
		instSlice.Set(reflect.Append(instSlice, destInst))
	}
	return rows.Close()
}

// Update 接收两种入参，平铺开来的键值对和map类型的键值对
// generator接收的参数是map类型的键值对，Update方法会动态判断传入参数的类型，如果不是map类型，则自动转换
// support map[string]interface{},also support kv list:"Name","Tom","Age",18,...
func (s *Session) Update(kv ...interface{}) (int64, error) {
	s.CallMethod(BeforeUpdate, nil)
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.SetSqlAndVars(clause.UPDATE, s.GetSchema().Name, m)
	sql, vars := s.clause.BuildInOrder(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterUpdate, nil)
	return result.RowsAffected()
}

// Delete 删除表里的所有数据
func (s *Session) Delete() (int64, error) {
	s.CallMethod(BeforeDelete, nil)
	s.clause.SetSqlAndVars(clause.DELETE, s.GetSchema().Name)
	sql, vars := s.clause.BuildInOrder(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	s.CallMethod(AfterDelete, nil)
	return result.RowsAffected()
}

// Count 计算表的数据项
func (s *Session) Count() (int64, error) {
	s.clause.SetSqlAndVars(clause.COUNT, s.GetSchema().Name)
	sql, vars := s.clause.BuildInOrder(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRaw()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// Limit 链式调用，简化代码的一种方法，使代码更简洁、易读。一般而言，当某个对象需要一次调用多个方法来设置其属性时，就非常适合改造为链式调用。
// 链式调用的原理也非常简单，某个对象调用某个方法后，将该对象的引用/指针返回，即可以继续调用该对象的其他方法。
func (s *Session) Limit(num int) *Session {
	s.clause.SetSqlAndVars(clause.LIMIT, num)
	return s
}

// Where adds limit condition to clause
func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.SetSqlAndVars(clause.WHERE, append(append(vars, desc), args...)...)
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.SetSqlAndVars(clause.ORDERBY, desc)
	return s
}

// First 根据传入的类型，利用反射构造切片，调用 Limit(1) 限制返回的行数，调用 Find 方法获取到查询结果
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).FindAll(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
