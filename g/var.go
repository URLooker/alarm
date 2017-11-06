package g

import "fmt"

type Event struct {
	Id          string `json:"id"`
	EventId     string `json:"event_id"`
	Status      string `json:"status"`
	Url         string `json:"url"`
	Ip          string `json:"ip"`
	EventTime   int64  `json:"event_time"`
	StrategyId  int64  `json:"strategy_id"`
	RespTime    int64  `json:"resp_time"`
	RespCode    string `json:"resp_code"`
	Result      int64  `json:"result"`
	CurrentStep int    `json:"current_step"`
}

type HistoryData struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Sms struct {
	Tos     string `json:"tos"`
	Content string `json:"content"`
}

type Mail struct {
	Tos     string `json:"tos"`
	Subject string `json:"subject"`
	Content string `json:"content"`
}

func (this *Sms) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Content:%s>",
		this.Tos,
		this.Content,
	)
}

func (this *Mail) String() string {
	return fmt.Sprintf(
		"<Tos:%s, Subject:%s, Content:%s>",
		this.Tos,
		this.Subject,
		this.Content,
	)
}
