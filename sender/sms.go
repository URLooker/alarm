package sender

import (
	"log"
	"strings"
	"time"

	"github.com/toolkits/net/httplib"

	"github.com/urlooker/alarm/g"
)

func ConsumeSms() {
	queue := g.Config.Queue.Sms
	for {
		L := PopAllSms(queue)
		if len(L) == 0 {
			time.Sleep(time.Millisecond * 200)
			continue
		}
		SendSmsList(L)
	}
}

func SendSmsList(L []*g.Sms) {
	for _, sms := range L {
		if sms.Tos == "" || sms.Content == "" {
			continue
		}

		toArr := strings.Split(sms.Tos, ",")
		log.Println("SmsCount", len(toArr))

		SmsWorkerChan <- 1
		go SendSms(sms)
	}
}

func SendSms(sms *g.Sms) {
	defer func() {
		<-SmsWorkerChan
	}()

	url := g.Config.Sms
	r := httplib.Post(url).SetTimeout(5*time.Second, 2*time.Minute)
	tos := strings.Replace(sms.Tos, "+86", "", -1)
	r.Param("tos", tos)
	r.Param("content", sms.Content)
	resp, err := r.String()
	if err != nil {
		log.Println(err)
		//SendMailToMaintainer("sender sms provider error", err.Error())
	}

	if g.Config.Debug {
		log.Println("==sms==>>>>", sms)
		log.Println("<<<<==sms==", resp)
	}

}
