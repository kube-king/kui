package constant

const (
	RootPath        = "./"
	DataPath        = RootPath + "data"
	BinaryPath      = RootPath + "bin"
	TemplatePath    = RootPath + "templates"
	DownloadAddress = "http://kui.in-engine.cn:30080/kui"
)

const (
	EtcdHeartbeatInterval = 500
	EtcdElectionTimeout   = 5000
	CertPath              = DataPath + "/cert"
)
