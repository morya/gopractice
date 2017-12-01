package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"regexp"
	"time"

	"github.com/go-redis/redis"
	"github.com/morya/utils/log"
	"golang.org/x/net/context"
)

type JsonValue map[string]interface{}
type ChanEvent chan *Counter

type App struct {
	sock      net.Conn
	chanEvent ChanEvent

	reg *regexp.Regexp

	rconn  *redis.Client
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApp(redisaddr, rediskey string) *App {
	rconn := redis.NewClient(&redis.Options{
		Addr: redisaddr,
		DB:   0,
	})

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		chanEvent: make(ChanEvent, 10),

		rconn:  rconn,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) isExiting() bool {
	select {
	case <-a.ctx.Done():
		return true
	default:
		break
	}

	return false
}

func (a *App) Run() {
	go a.reporter()

	defer func() {
		log.Info("main proc exiting")
		a.cancel()
	}()

	for {
		if a.isExiting() {
			log.Info("finish")
			return
		}

		msg, err := a.pollEvent()
		if err != nil {
			log.Info("nothing to process")
			continue
		}

		log.Infof("got event %v", msg)
		err = a.procEvent(msg)
		if err != nil {
			return
		}
	}
}

func (a *App) init(unixaddress string, reg string) bool {
	if !a.buildConn(unixaddress) {
		return false
	}

	if !a.buildRegexp(reg) {
		return false
	}
	return true
}

func (a *App) buildRegexp(reg string) bool {
	var reobj = regexp.MustCompile(reg)
	a.reg = reobj

	return true
}

// poll things from redis
func (a *App) pollEvent() (string, error) {
	value, err := a.rconn.BLPop(time.Second*2, *flagRedisKey).Result()
	if err != nil {
		return "", err
	}

	return value[1], nil
}

func (a *App) procMsgLine(msg string) error {
	/* matches 结果
	* 0: 全部匹配值
	* 1: first group
	 */
	matches := a.reg.FindStringSubmatch(msg)
	if len(matches) > 0 {
		c, err := fromLogStr(matches[1])
		if err != nil {
			return err
		}

		a.chanEvent <- c
	}
	return nil
}

// send decoded data to inside channel
func (a *App) procEvent(evt string) error {
	var fbevt = make(JsonValue)
	err := json.Unmarshal([]byte(evt), &fbevt)
	if err != nil {
		log.InfoError(err, "decode original filebeat event failed")
		return err
	}

	msg := fmt.Sprintf("%v", fbevt["message"])
	log.Infof("msg = %v", msg)
	a.procMsgLine(msg)
	return nil
}

func (app *App) sendLog(c *Counter) error {
	var err error
	var data []byte

	var ts = time.Unix(c.UnixTimeStamp/1000, (c.UnixTimeStamp%1000)*1000)
	c.Timestamp = ts.Format(time.RFC3339Nano)

	data, err = json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = app.sock.Write(data)
	if err != nil {
		return err
	}
	io.WriteString(app.sock, "\n") // this is kind important!
	return nil
}

func (a *App) reporter() {
	var c *Counter
	for {
		select {
		case <-a.ctx.Done():
			return

		case c = <-a.chanEvent:
			break
		}

		a.sendLog(c)
	}
}
