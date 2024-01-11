package kafka

type QueueSetup struct {
	TopicMain    string   // Tópico principal (siempre es requerido), ejemplo: "atreel-message-placed"
	GroupID      string   // Identificador de grupo de kafka
	Brokers      []string // Colección de brokers de kafka
	Timeout      int      // Tiempo de espera máximo por un mensaje de la cola
	RequestDelay int      // Tiempo de inactividad después de un error o el consumo del timeout de la cola
}
