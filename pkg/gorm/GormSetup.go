package gorm

type Setup struct {
	User               string
	Password           string
	Server             string
	Database           string
	Port               int
	ConnectTimeOut     int
	Environment        string
	MaxIdleConnections int
	MaxOpenConnections int
	ConnMaxLifetime    int
}
