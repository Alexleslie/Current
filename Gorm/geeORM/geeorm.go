package geeORM

import (
	"Current/Gorm/geeORM/dialect"
	"Current/Gorm/geeORM/session"
	"database/sql"
	"fmt"
	"github.com/golang/glog"
	"strings"
)

/*用户交互*/
/*
数据库迁移，只支持字段新增和删除
新增字段：ALTER TABLE table_name ADD COLUMN col_name,col_type;
删除字段：CREATE TABLE new_table AS SELECT col1,col2,...from old_table
		DROP TABLE old_table
		ALTER TABLE new_table RENAME TO old_table;
*/

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

// TxFunc 回调函数，作为入参传递给engin.Transaction()，发生任何错误，自动回滚，没发生错误则提交
type TxFunc func(*session.Session) (interface{}, error)

// Transaction 执行事务
func (engine *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := engine.NewSession()
	// 事务开始
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil;don't change it
		} else {
			err = s.Commit() // err is nil;if Commit return error update err
		}
	}()
	return f(s)
}

// NewEngine 链接数据库，返回*sql.DB，调用db.Ping()，检查数据库能否正常连接
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		glog.Error(err)
		return
	}
	//Sent a ping to make sure the databases connection is alive.
	if err = db.Ping(); err != nil {
		glog.Error(err)
		return
	}
	//make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		glog.Errorf("[NewEngine] Dialect %s Not Found", driver)
	}
	e = &Engine{db: db, dialect: dial}
	glog.Info("[NewEngine] Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		glog.Error("Failed to close database")
	}
	glog.Info("Close database success")
}

// NewSession 通过Engine实例创建会话，与数据库进行交互
func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

// difference return a-b
// 计算前后两表字段切片的差集，新表-旧表=新增字段，旧表-新表=删除字段
func difference(a []string, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}
	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(value interface{}) error {
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.SetSchemaByInstance(value).HasTable() {
			glog.Info("[engine.Migrate] Table %s doesn't exist", s.GetSchema().Name)
			return nil, s.CreateTable()
		}
		table := s.GetSchema()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRaws()
		columns, _ := rows.Columns()
		//计算新增和删除字段
		addCols := difference(table.FieldNames, columns)
		delCols := difference(columns, table.FieldNames)
		glog.Infof("[engine.Migrate] added cols %v,deleted cols %v", addCols, delCols)

		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.Name, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		filedStr := strings.Join(table.FieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s from %s;", tmp, filedStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	})
	return err
}
