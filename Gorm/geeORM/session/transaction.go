package session

import "github.com/golang/glog"

/*封装十五的Begin、Commit和Rollback三个接口，统一打印日志，方便定位问题*/

func (s *Session) Begin() (err error) {
	glog.Info("[Session.Begin] Transaction begin")
	//调用s.db.Begin()获得*sql.Tx对象，赋值给s.tx
	if s.tx, err = s.db.Begin(); err != nil {
		glog.Error(err)
		return
	}
	return
}

func (s *Session) Commit() (err error) {
	glog.Info("[Session.Commit] Transaction commit")
	if err = s.tx.Commit(); err != nil {
		glog.Error(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	glog.Info("[Session.Rollback] Transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		glog.Error(err)
	}
	return
}
