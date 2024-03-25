package certificate

// TLSConfig TLS 配置
type TLSConfig struct {
	Subject            // 证书参数
	KeyOutPath  string // key生成路径
	CertOutPath string // cert生成路径
}

// TLSContext TLS证书内容
type TLSContext struct {
	Cert string
	Key  string
}

// Subject TLS证书标识
type Subject struct {
	Country            string // 国家
	Province           string // 城市
	Organization       string // 区域
	OrganizationalUnit string // 区域单位
	CommonName         string // 域名
}
