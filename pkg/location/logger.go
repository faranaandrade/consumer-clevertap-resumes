package location

type Logger interface {
	Error(moduleName, functionName string, err error)
}
