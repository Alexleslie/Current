package clause

import (
	"fmt"
	"strings"
)

/*
构造数据库语句
SELECT col1,col2,...
	FROM table_name
	WHERE [conditions]
	GROUP BY col1
	HAVING [conditions]
INSERT INTO table_name(col1,col2,col3,...) VALUES
	(A1,A2,A3,...)
	(B1,B2,B3,...)
	...
*/

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// 根据nums数生成("?,?,?,?")
func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

// 生成插入语句 INSERT INTO [table_name]([col1],[col2],[col3],...)
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

// 生成VALUE语句 VALUES (a1,a2,a3,...),(b1,b2,b3,...)
func _values(values ...interface{}) (string, []interface{}) {
	var bindStr string
	var sql strings.Builder
	var vars []interface{}
	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _select(values ...interface{}) (string, []interface{}) {
	//SELECT $fileds FROM $tableName
	tableName := values[0]
	fileds := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fileds, tableName), []interface{}{}
}

func _limit(values ...interface{}) (string, []interface{}) {
	//LIMIT $nums
	return "LIMIT ?", values
}

func _where(values ...interface{}) (string, []interface{}) {
	//WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s ", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	//ORDER BT $
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

func _update(values ...interface{}) (string, []interface{}) {
	//UPDATE $table_name SET $col1=?,col2=?,...
	tableName := values[0]                  //表名
	m := values[1].(map[string]interface{}) //待更新的键对值
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

func _delete(values ...interface{}) (string, []interface{}) {
	//DELETE FROM $table_name
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	//SELECT count(*) FROM $table_name...
	return _select(values[0], []string{"count(*)"})
}
