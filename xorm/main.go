package main

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/morya/utils/log"
)

var engine *xorm.Engine

type User struct {
	ID      int64     `xorm:"pk autoincr"`
	Name    string    `xorm:"UNIQUE NOT NULL"`
	Created time.Time `xorm:"-"`
}

func main() {
	var err error

	log.SetFlags(log.LstdFlags)
	engine, err = xorm.NewEngine("mysql", "morya:$Perl234@/coins?charset=utf8")
	if err != nil {
		log.InfoErrorf(err, "link failed")
		return
	}

	log.Infof("link mysql succ")
    engine.SetMapper(core.GonicMapper{})
	err = engine.Sync2(new(User))
	log.Infof("sync result %v", err)
}
