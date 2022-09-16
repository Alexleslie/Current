package session

import (
	"Current/tools/logc"
)

/*封装十五的Begin、Commit和Rollback三个接口，统一打印日志，方便定位问题*/

func (s *Session) Begin() (err error) {
	logc.Info("[Session.Begin] Transaction begin")
	//调用s.db.Begin()获得*sql.Tx对象，赋值给s.tx
	if s.tx, err = s.db.Begin(); err != nil {
		logc.Error("[Session.Begin] Begin error, err=[%+v]", err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	logc.Info("[Session.Commit] Transaction commit")
	if err = s.tx.Commit(); err != nil {
		logc.Error("[Session.Commit] Commit error, err=[%+v]", err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	logc.Info("[Session.Rollback] Transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		logc.Error("[Session.Rollback] Rollback error, err=[%+v]", err)
	}
	return
}
