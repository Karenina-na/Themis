package sync

import (
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
)

//
// ChangeToFollow
// @Description: change to follow
// @return       E error
//
func ChangeToFollow() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("ChangeToFollow-sync", util.Strval(r))
		}
	}()
	Bean.RoutinePool.CreateWork(func() (E error) {
		err := Follow()
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})
	return nil
}

//
// ChangeToCandidate
// @Description: change to candidate
// @return       E error
//
func ChangeToCandidate() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("ChangeToCandidate-sync", util.Strval(r))
		}
	}()
	Bean.RoutinePool.CreateWork(func() (E error) {
		err := Candidate()
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})
	return nil
}

//
// ChangeToLeader
// @Description: change to leader
// @return       E error
//
func ChangeToLeader() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("ChangeToLeader-sync", util.Strval(r))
		}
	}()
	Bean.RoutinePool.CreateWork(func() (E error) {
		err := Leader()
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})
	return nil
}
