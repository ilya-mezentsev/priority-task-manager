package settings

type (
	WebSettings struct {
		Port   int    `json:"port"`
		Domain string `json:"domain"`
	}

	DBSettings struct {
		Host       string `json:"host"`
		Port       int    `json:"port"`
		User       string `json:"user"`
		Password   string `json:"password"`
		DBName     string `json:"db_name"`
		Connection struct {
			RetryCount   int `json:"retry_count"`
			RetryTimeout int `json:"retry_timeout"`
		} `json:"connection"`
		CacheLifetimeSeconds int `json:"cache_lifetime_seconds"`
	}

	QueueSettings struct {
		Name        string `json:"name"`
		MaxPriority int    `json:"max_priority"`
		Durable     bool   `json:"durable"`
		AutoDelete  bool   `json:"auto_delete"`
		Exclusive   bool   `json:"exclusive"`
		NoWait      bool   `json:"no_wait"`
		AutoAck     bool   `json:"auto_ack"`
		NoLocal     bool   `json:"no_local"`
		Consumer    string `json:"consumer"`
	}

	RabbitMQSettings struct {
		Host                     string        `json:"host"`
		Port                     int           `json:"port"`
		User                     string        `json:"user"`
		Password                 string        `json:"password"`
		PublishingTimeoutSeconds int           `json:"publishing_timeout_seconds"`
		Queue                    QueueSettings `json:"queue"`
	}

	BasicAuthSettings struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
)
