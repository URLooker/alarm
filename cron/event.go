package cron

import (
	"encoding/json"
	"log"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/urlooker/web/api"
	"github.com/urlooker/web/model"

	"github.com/urlooker/alarm/backend"
	"github.com/urlooker/alarm/cache"
	"github.com/urlooker/alarm/g"
	"github.com/urlooker/alarm/sender"
)

func ReadEvent() {
	for {
		event, err := PopEvent()
		if err != nil {
			log.Println("error:", err)
			time.Sleep(time.Second)
			continue
		}
		if event != nil {
			mail := make([]string, 0)
			sms := make([]string, 0)
			users := getUsers(event.StrategyId)

			mailContent := sender.BuildMail(event)
			smsContent := sender.BuildSms(event)
			for _, user := range users {
				mail = append(mail, user.Email)
				sms = append(sms, user.Phone)
			}

			sender.WriteSms(sms, smsContent)
			sender.WriteMail(mail, smsContent, mailContent)
		}

		time.Sleep(1 * time.Second)
	}
}

func getUsers(sid int64) []*model.User {
	var usersResp api.UsersResponse
	var users []*model.User
	strategy, exists := cache.StrategyMap.Get(sid)
	if !exists {
		log.Printf("strategyId: %d not exists", sid)
		return users
	}

	err := backend.CallRpc("Web.GetUsersByTeam", strategy.Teams, &usersResp)
	if err != nil {
		log.Println("Web.GetUsersByTeam Error:", err)
		return users
	}

	if usersResp.Message != "" {
		log.Println("Web.GetUsersByTeam Error:", usersResp.Message)
		return users
	}
	users = usersResp.Data

	return users
}

func PopEvent() (*model.Event, error) {
	rc := g.RedisConnPool.Get()
	defer rc.Close()

	reply, err := redis.String(rc.Do("RPOP", g.Config.Alarm.QueuePattern))
	if err != nil {
		if err != redis.ErrNil {
			log.Printf("get alarm event from redis fail: %v", err)
		}
		return nil, nil
	}

	var event model.Event
	err = json.Unmarshal([]byte(reply), &event)
	if err != nil {
		log.Printf("parse alarm event fail: %v", err)
		return nil, err
	}

	if g.Config.Debug {
		log.Println("======>>>>")
		log.Println(event.String())
	}
	return &event, nil
}
