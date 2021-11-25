package models

import "fmt"

type Conf struct {
	Log struct {
		LogLevel string `json:"log_level"`
		LogPath  string `json:"log_path"`
	} `json:"log"`

	Server struct {
		HTTPServer Server `json:"http_server"`
		RPCServer  Server `json:"rpc_server"`
	} `json:"server"`

	Rpc struct {
		Cert struct {
			Client struct {
				ClientCertFile string `json:"client_cert_file"`
			} `json:"server"`
			Server struct {
				ServerCertFile string `json:"server_cert_file"`
				ServerKeyFile  string `json:"server_key_file"`
			} `json:"server"`
		} `json:"cert"`
		SubServers struct {
			ManageServer      string `json:"manage_server"`
			UserAccountServer string `json:"user_account_server"`
			SpaceServer       string `json:"space_server"`
			ConnectServer     string `json:"connect_server"`
			ImServer          string `json:"im_server"`
			FileServer        string `json:"file_server"`
			ShopServer        string `json:"shop_server"`
			LiveServer        string `json:"live_server"`

			AccountDaoServer string `json:"account_dao_server"`
			SpaceDaoServer   string `json:"space_dao_server"`
			CommentDaoServer string `json:"comment_dao_server"`
			ImDaoServer      string `json:"im_dao_server"`
			ShopDaoServer    string `json:"shop_dao_server"`
			SmsDaoServer     string `json:"sms_dao_server"`
			GameDaoServer    string `json:"game_dao_server"`
			LiveDaoServer    string `json:"live_dao_server"`
		} `json:"sub_servers"`
	} `json:"rpc"`

	Register struct {
		ETCD struct {
			Cluster []string `json:"cluster"`
		} `json:"etcd"`
		HealthyRollTime int64 `json:"healthy_roll_time"`
	} `json:"register"`

	Cache struct {
		Redis struct {
			MainCache Redis `json:"main_cache"`
			SubCache  struct {
				ReadCache Redis `json:"read_cache"`
			} `json:"sub_cache"`
		} `json:"redis"`
	} `json:"cache"`

	MessageQueue struct {
		NSQ struct {
			Cluster string `json:"cluster"`
		} `json:"nsq"`
		PublishTopics struct {
			// websocket通知
			WebsocketNotify string `json:"websocket_notify"`
			// 用户积分变动
			UserCoin string `json:"user_coin"`
			// 文件删除
			DeleteFile string `json:"delete_file"`
		} `json:"publish_topics"`
		ConsumeTopics struct {
			WebsocketNotify TopicConsume `json:"websocket_notify"`
			UserCoin        TopicConsume `json:"user_coin"`
			DeleteFile      TopicConsume `json:"delete_file"`
		} `json:"consume_topics"`
	} `json:"message_queue"`

	Database struct {
		MainDatabase Mysql `json:"main_database"`
		SubDatabase  struct {
			AccountDatabase Mysql `json:"account_database"`
			ShopDatabase    Mysql `json:"shop_database"`
			SmsDatabase     Mysql `json:"sms_database"`
			SpaceDatabase   Mysql `json:"space_database"`
			ImDatabase      Mysql `json:"im_database"`
		} `json:"sub_database"`
	}

	Endpoint        string `json:"endpoint"`
	AccessKeyId     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`

	TokenSecret string `json:"token_secret"`

	Stuns []string `json:"stuns"`
	Turn  struct {
		URLs       []string `json:"urls"`
		Username   string   `json:"username"`
		Credential string   `json:"credential"`
	} `json:"turn"`
}

type Server struct {
	Name     string `json:"name"`
	Serial   int    `json:"serial"`
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Register string `json:"register"`
}

type Redis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

type Mysql struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	DBName   string `json:"db_name"`
}

func (ms Mysql) GetMysqlUrl() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		ms.Username, ms.Password, ms.Host, ms.Port, ms.DBName)
}

type TopicConsume struct {
	Topic   string `json:"topic"`
	Channel string `json:"channel"`
}
