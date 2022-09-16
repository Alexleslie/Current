package session

import (
	"Current/Gorm/geeORM/clause"
	"Current/Gorm/geeORM/dialect"
	"Current/Gorm/geeORM/schema"
	"Current/tools/logc"
	"database/sql"
	"strings"
)

/*数据库交互*/

type Session struct {
	db      *sql.DB         // 使用sql.Open()方法连接数据库成功之后返回的指针
	dialect dialect.Dialect // Go语言转换成对应数据库SQL语句的管理实例
	tx      *sql.Tx         // SQL事务执行时指针
	schema  *schema.Schema  // 数据库表格
	clause  clause.Clause   // 数据库SQL语句生成拼接
	sql     strings.Builder // 用来拼接SQL语句和SQL语句中占位符的对应值，用户调用Raw()方法即可改变这两个变量的值
	sqlVars []interface{}
}

// CommonDB is a minimal function set of db
type CommonDB interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) ClearSession() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// DB returns tx if a tx begins, otherwise return *sql.DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	//log.Info(s.db)
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec 封装：统一打印日志（包括执行的SQL语句和错误日志），执行完成后清空s.sql和s.sqlVars两个变量，这样Session可以复用
// 开启一次会话，可以执行多次多次SQL
// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.ClearSession()
	logc.Info("[Session.Exec] sql=[%+v], sqlVars=[%+v]", s.sql.String(), s.sqlVars)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		logc.Error("[Session.Exec] Exec sql error, sql=[%+v], sqlVars=[%+v], err=[%+v]", s.sql.String(), s.sqlVars, err)
	}
	return
}

// QueryRaw gets a record form db
func (s *Session) QueryRaw() *sql.Row {
	defer s.ClearSession()
	logc.Info("[Session.QueryRaw] sql=[%+v], sqlVars=[%+v]", s.sql.String(), s.sqlVars)
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRaws gets a list of records from db
func (s *Session) QueryRaws() (rows *sql.Rows, err error) {
	defer s.ClearSession()
	logc.Info("[Session.QueryRaws] sql=[%+v], sqlVars=[%+v]", s.sql.String(), s.sqlVars)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		logc.Error("[Session.QueryRaws] Query error, err=[%+v]", err)
	}
	return
}
