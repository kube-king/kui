package logger

// 定义默认的常量
const (
	defaultBaseDirectoryName  = "logs"      // 日志根目录
	defaultInfoDirectoryName  = "info"      // info日志目录
	defaultWarnDirectoryName  = "warn"      // warn日志目录
	defaultErrorDirectoryName = "error"     // error日志目录
	defaultInfoFileName       = "info.log"  // info日志文件
	defaultWarnFileName       = "warn.log"  // warn日志文件
	defaultErrorFileName      = "error.log" // error日志文件
	defaultLogFileMaxSize     = 128         // 日志文件大小，单位：MB
	defaultLogFileMaxBackups  = 30          // 日志文件保留个数 多于180个文件后，清理比价旧的日志
	defaultLogFileMaxAge      = 1           // 日志文件一天一切隔
	defaultLogFileCompress    = false       // 日志文件是否压缩
	defaultLogPrintTag        = true        // true:在终端和文件同时输出日志; false:只在文件输出日志
)

// Config 配置文件结构体
type Config struct {
	BaseDirectoryName  string
	InfoDirectoryName  string
	WarnDirectoryName  string
	ErrorDirectoryName string
	InfoFileName       string
	WarnFileName       string
	ErrorFileName      string
	LogFileMaxSize     int
	LogFileMaxBackups  int
	LogFileMaxAge      int
	LogFileCompress    bool
	LogPrintTag        bool
}

// Option 定义配置选项函数
type Option func(*Config)

// SetBaseDirectoryName 自定义日志根目录
func SetBaseDirectoryName(name string) Option {
	return func(c *Config) {
		c.BaseDirectoryName = name
	}
}

// SetInfoDirectoryName 自定义info日志目录
func SetInfoDirectoryName(name string) Option {
	return func(c *Config) {
		c.InfoDirectoryName = name
	}
}

// SetWarnDirectoryName 自定义warn日志目录
func SetWarnDirectoryName(name string) Option {
	return func(c *Config) {
		c.WarnDirectoryName = name
	}
}

// SetErrorDirectoryName 自定义error日志目录
func SetErrorDirectoryName(name string) Option {
	return func(c *Config) {
		c.ErrorDirectoryName = name
	}
}

// SetInfoFileName 自定义info文件名
func SetInfoFileName(name string) Option {
	return func(c *Config) {
		c.InfoFileName = name
	}
}

// SetWarnFileName 自定义warn文件名
func SetWarnFileName(name string) Option {
	return func(c *Config) {
		c.WarnFileName = name
	}
}

// SetErrorFileName 自定义error文件名
func SetErrorFileName(name string) Option {
	return func(c *Config) {
		c.ErrorFileName = name
	}
}

// SetLogFileMaxSize 自定义日志文件大小
func SetLogFileMaxSize(size int) Option {
	return func(c *Config) {
		c.LogFileMaxSize = size
	}
}

// SetLogFileMaxBackups 自定义日志文件保留个数
func SetLogFileMaxBackups(size int) Option {
	return func(c *Config) {
		c.LogFileMaxBackups = size
	}
}

// SetLogFileMaxAge 自定义日志文件切隔间隔
func SetLogFileMaxAge(size int) Option {
	return func(c *Config) {
		c.LogFileMaxAge = size
	}
}

// SetLogFileCompress 自定义日志文件是否压缩
func SetLogFileCompress(compress bool) Option {
	return func(c *Config) {
		c.LogFileCompress = compress
	}
}

// SetLogPrintTag 自定义日志输出标记位 true:在终端和文件同时输出日志; false:只在文件输出日志
func SetLogPrintTag(tag bool) Option {
	return func(c *Config) {
		c.LogPrintTag = tag
	}
}

// NewConfig 应用函数选项配置
func NewConfig(opts ...Option) Config {
	// 初始化默认值
	defaultConfig := Config{
		BaseDirectoryName:  defaultBaseDirectoryName,
		InfoDirectoryName:  defaultInfoDirectoryName,
		WarnDirectoryName:  defaultWarnDirectoryName,
		ErrorDirectoryName: defaultErrorDirectoryName,
		InfoFileName:       defaultInfoFileName,
		WarnFileName:       defaultWarnFileName,
		ErrorFileName:      defaultErrorFileName,
		LogFileMaxSize:     defaultLogFileMaxSize,
		LogFileMaxBackups:  defaultLogFileMaxBackups,
		LogFileMaxAge:      defaultLogFileMaxAge,
		LogFileCompress:    defaultLogFileCompress,
		LogPrintTag:        defaultLogPrintTag,
	}

	// 依次调用opts函数列表中的函数，为结构体成员赋值
	for _, opt := range opts {
		opt(&defaultConfig)
	}

	return defaultConfig
}
