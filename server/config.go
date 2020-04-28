package server

// Config represents a server config
type Config struct {
	ListenAddr                 string
	TLSListenAddr              string
	TLS                        *TLSConfig
	ReadAuthToken              string
	WriteAuthToken             string
	AllowPublicCreateGenerated bool // If true, WriteAuthToken is NOT required when creating an entry that does not contain a custom ID
	GeneratedIDLen             int
	Verbose                    bool
}

// TLSConfig represents a TLS configuration
type TLSConfig struct {
	KeyFile  string
	CertFile string
}
