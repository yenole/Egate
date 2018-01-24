//+build serialize

package serialize

import (
	"egate/elog"
	"egate/gate"
	"errors"
	"github.com/jinzhu/gorm"
	"reflect"
	"strings"
	"sync"
	"time"
)

var (
	Serialize              *serialize
	ERR_SERIALIZE_NOT_INIT = errors.New("serialize not init!")
)

type serialize struct {
	*gorm.DB
	tick  *time.Ticker
	mutex sync.RWMutex
	modes map[gate.IDer]map[string]interface{}
}

func Update(m gate.IDer, cls ...string) error {
	if Serialize != nil {
		Serialize.mutex.Lock()
		defer Serialize.mutex.Unlock()
		if _, ok := Serialize.modes[m]; !ok {
			Serialize.modes[m] = map[string]interface{}{}
		}
		for _, v := range cls {
			fv := reflect.ValueOf(m).Elem().FieldByName(v)
			if fv.CanInterface() {
				Serialize.modes[m][gorm.ToDBName(v)] = fv.Interface()
			} else if fv.CanAddr() {
				Serialize.modes[m][gorm.ToDBName(v)] = fv.Addr().Interface()
			}

		}
		return nil
	}
	return ERR_SERIALIZE_NOT_INIT
}

func (s *serialize) models() map[gate.IDer]map[string]interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	defer func() {
		s.modes = map[gate.IDer]map[string]interface{}{}
	}()
	return s.modes
}

func (s *serialize) work() {
	for {
		<-s.tick.C
		func() {
			defer func() {
				if err := recover(); err != nil {
					elog.Fatal("%v", err)
				}
			}()
			for m, cls := range s.models() {
				s.Model(m).Updates(cls)
			}
		}()
	}
}

func DSN(dsn string, dt time.Duration) *serialize {
	dsns := strings.Split(dsn, "://")
	if len(dsns) > 1 {
		db, err := gorm.Open(dsns[0], dsns[1])
		if err != nil {
			panic(err)
		}
		Serialize = &serialize{DB: db, modes: map[gate.IDer]map[string]interface{}{}}
		if dt > 0 {
			Serialize.tick = time.NewTicker(dt)
			go Serialize.work()
		}
	}
	return Serialize
}
