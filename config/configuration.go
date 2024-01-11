package config

type Configuration struct {
	ClevertapConfig   ClevertapConfig
	Environment       string `env:"GO_ENV,required"`
	ArtifactVersion   string `env:"VERSION,required"`
	APIEvents         string `env:"API_EVENTS,required"`
	APITimeout        int    `env:"API_TIMEOUT,default=5"`
	HTTPClient        HTTPClient
	QueueSetup        KafkaQueue
	ProcessingRetries ProcessingRetries
	GormConnection    GormConnection
}

type ClevertapConfig struct {
	APISetup        ClevertapAPISetup
	HTTPClientSetup HTTPClient
}

type ClevertapAPISetup struct {
	Host              string `env:"API_CLEVERTAP,required"`
	Timeout           int    `env:"API_TIMEOUT,required"`
	AccountID         string `env:"CLEVERTAP_ACCOUNT_ID,required"`
	ClevertapPasscode string `env:"CLEVERTAP_PASSCODE,required"`
}

type KafkaQueue struct {
	TopicMain    string   `env:"TOPIC_MAIN,required"`             // Tópico principal"
	GroupID      string   `env:"KAFKA_GROUPID,required"`          // Identificador de grupo de kafka
	Brokers      []string `env:"KAFKA_BROKERS,required"`          // Colección de brokers de kafka
	Timeout      int      `env:"QUEUE_TIMEOUT,default=10"`        // Tiempo de espera máximo por un mensaje de la cola
	RequestDelay int      `env:"QUEUE_REQUEST_DELAY,default=100"` // Tiempo de inactividad entre un error o timeout de la cola
}

type HTTPClient struct {
	APITimeout  int `env:"API_TIMEOUT,default=4"`
	MaxRetries  int `env:"MAX_RETRIES_REQUESTS,default=2"`
	MinTimeWait int `env:"MIN_RETRIES_TIME_WAIT,default=10"`
}

type ProcessingRetries struct {
	TopicRetry   string `env:"TOPIC_RETRY"`            // Tópico para reintentos de procesamiento (opcional)
	WaitForRetry int    `env:"WAIT_FOR_RETRY"`         // Cantidad de tiempo a esperar para procesar (en segundos)
	MaxRetries   int    `env:"MAX_PROCESSING_RETRIES"` // Cantidad máxima de reprocesamientos permitidos en el tópico principal
}

type GormConnection struct {
	Server             string `env:"MS_SQL_SERVER,required"`
	Database           string `env:"MS_SQL_DATABASE,required"`
	User               string `env:"MS_SQL_USER,required"`
	Password           string `env:"MS_SQL_PASSWORD,required"`
	ConnectTimeOut     int    `env:"MS_SQL_CONNECT_TIMEOUT,default=45"`
	MaxOpenConnections int    `env:"MS_SQL_MAX_OPEN_CONNECTIONS,default=5"`
	MaxIdleConnections int    `env:"MS_SQL_MAX_IDLE_CONNECTIONS,default=30"`
	ConnMaxLifetime    int    `env:"MS_SQL_CONN_MAX_LIFETIME,default=60"`
	Port               int    `env:"MS_SQL_PORT,default=1433"`
}

// https://blog.pragmatists.com/retrying-consumer-architecture-in-the-apache-kafka-939ac4cb851a
// https://eng.uber.com/reliable-reprocessing/

// EJEMPLO DE CONFIGURACIÓN DE UNA ESTRATEGIA DE RETRIES:

/*

Every 5 minutes — 2 times (1 processing, 1 retry)
Then after 30 minutes — 3 times (1 processing, 2 retry)
Then after 1 hour only one time (1 processing, 0 retry)
Then we skip the message (Dead Letter Queue)

-- processor principal:
TopicMain    = "mys-application-successful"
TopicRetry   = "mys-application-successful-backwards-retry_5m"
WaitForRetry = 0
MaxRetries   = 0

-- processor retry-0:
TopicMain    = "mys-application-successful-backwards-retry_5m"
TopicRetry   = "mys-application-successful-backwards-retry_30m"
WaitForRetry = 300
MaxRetries   = 1

-- processor retry-1:
TopicMain    = "mys-application-successful-backwards-retry_30m"
TopicRetry   = "mys-application-successful-backwards-retry_1h"
WaitForRetry = 1800
MaxRetries   = 2

-- processor retry-2:
TopicMain    = "mys-application-successful-backwards-retry_1h"
TopicRetry   = "mys-application-successful-backwards-dlq"
WaitForRetry = 3600
MaxRetries   = 0

*/
