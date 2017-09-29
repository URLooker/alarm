package g

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/toolkits/file"
)

type GlobalConfig struct {
	Debug  bool          `json:"debug"`
	Remain int           `json:"remain"`
	Rpc    *RpcConfig    `json:"rpc"`
	Web    *WebConfig    `json:"web"`
	Alarm  *AlarmConfig  `json:"alarm"`
	Queue  *QueueConfig  `json:"queue"`
	Mysql  *MysqlConfig  `json:"mysql"`
	Worker *WorkerConfig `json:"worker"`
	Smtp   *SmtpConfig   `json:"smtp"`
	Sms    string        `json:"sms"`
}

type MysqlConfig struct {
	Addr string `json:"addr"`
	Idle int    `json:"idle"`
	Max  int    `json:"max"`
}

type RpcConfig struct {
	Listen string `json:"listen"`
}

type RedisConfig struct {
	Dsn          string `json:"dsn"`
	Db           int    `json.db`
	MaxIdle      int    `json:"maxIdle"`
	ConnTimeout  int    `json:"connTimeout"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
}

type AlarmConfig struct {
	Enabled      bool         `json:"enabled"`
	MinInterval  int64        `json:"minInterval"`
	QueuePattern string       `json:"queuePattern"`
	Redis        *RedisConfig `json:"redis"`
}

type WebConfig struct {
	Addrs    []string `json:"addrs"`
	Timeout  int      `json:"timeout"`
	Interval int      `json:"interval"`
}

type QueueConfig struct {
	Mail string `json:"mail"`
	Sms  string `json:"sms"`
}

type WorkerConfig struct {
	Sms  int `json:"sms"`
	Mail int `json:"mail"`
}

type SmtpConfig struct {
	Addr     string `json:"addr"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
}

var (
	Config     *GlobalConfig
	configLock = new(sync.RWMutex)
)

func Parse(cfg string) error {
	if cfg == "" {
		return fmt.Errorf("use -c to specify configuration file")
	}

	if !file.IsExist(cfg) {
		return fmt.Errorf("configuration file %s is not exists", cfg)
	}

	configContent, err := file.ToTrimString(cfg)
	if err != nil {
		return fmt.Errorf("read configuration file %s fail %s", cfg, err.Error())
	}

	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		return fmt.Errorf("parse configuration file %s fail %s", cfg, err.Error())
	}

	configLock.Lock()
	defer configLock.Unlock()
	Config = &c

	log.Println("load configuration file", cfg, "successfully")
	return nil
}
