package server

// Config represents a server config
type Config struct {
	ListenAddr     string
	TLSListenAddr  string
	TLS            *TLSConfig
	ReadAuthToken  string
	WriteAuthToken string
	GeneratedIDLen int
	Verbose        bool
}

// TLSConfig represents a TLS configuration
type TLSConfig struct {
	KeyFile  string
	CertFile string
}
