package consumer

type Setup struct {
	ArtifactVersion string
	RequestDelay    int // Tiempo de inactividad después de un error o el consumo del timeout de la cola
}
