package main

import (
	"flag"
	"github.com/aagun1234/rabbit-mtcp-ws/client"
	"github.com/aagun1234/rabbit-mtcp-ws/logger"
	"github.com/aagun1234/rabbit-mtcp-ws/server"
	"github.com/aagun1234/rabbit-mtcp-ws/tunnel"
	"log"
	"strings"
)

var Version = "1.0.9ws"//"No version information"

const (
	ClientMode = iota
	ServerMode
	DefaultPassword = "PASSWORD"
)

/*
type Config struct {
	Mode             string        `yaml:"mode"`              // 运行模式: "client" 或 "server"
	ConfigFile       string        `yaml:"-"`                 // 配置文件路径 (不写入YAML)
	LogLevel         string        `yaml:"log_level"`         // 日志级别: "debug", "info", "warn", "error"`

	// Client 模式配置
	ClientListenAddr string        `yaml:"client_listen_addr"` // 客户端侦听的本地TCP地址 (例如: "127.0.0.1:1080")
	StatusServer     string        `yaml:"status_server"`      // 状态服务侦听的本地TCP地址 (例如: "127.0.0.1:8010")
	StatusACL     string           `yaml:"status_acl"`      // 状态服务ACL
	ServerURLs       []string      `yaml:"server_urls"`        // 服务端WebSocket URL列表 (例如: ["ws://server1:8081/tunnel", "wss://server2:8082/tunnel"])
	DialTimeout      time.Duration `yaml:"dial_timeout"`       // 拨号超时时间
	PingInterval     time.Duration `yaml:"ping_interval"`      // Ping消息发送间隔
	PongTimeout      time.Duration `yaml:"pong_timeout"`       // 等待Pong消息的超时时间
	ReconnectDelay   time.Duration `yaml:"reconnect_delay"`    // 重连间隔
	MaxRetries       int           `yaml:"max_retries"`        // 连接重试最大次数
	ClientTLSCAFile  string        `yaml:"client_tls_ca_file"` // 客户端用于验证服务端证书的CA文件路径
	ClientTLSCertFile string       `yaml:"client_tls_cert_file"` // 客户端证书文件路径 (可选，用于双向认证)
	ClientTLSKeyFile string        `yaml:"client_tls_key_file"`  // 客户端密钥文件路径 (可选，用于双向认证)
	ClientInsecureSkipVerify bool `yaml:"client_insecure_skip_verify"` // 客户端是否跳过服务端证书验证

	// Server 模式配置
	ServerListenAddrs []string        `yaml:"server_listen_addrs"` // 服务端侦听的WebSocket地址 (例如: "0.0.0.0:8081")
	TargetAddr     string           `yaml:"target_addr"` // 目标服务
	TargetConnectTimeout time.Duration `yaml:"target_connect_timeout"` // 服务端连接目标地址的超时时间
	ServerTLSCertFile string       `yaml:"server_tls_cert_file"` // 服务端证书文件路径
	ServerTLSKeyFile string        `yaml:"server_tls_key_file"`  // 服务端密钥文件路径

	// 共享配置
	AuthKey          string        `yaml:"auth_key"`           // 认证密钥
	Password         string        `yaml:"password"`           // 认证密钥
	EnableEncryption bool          `yaml:"enable_encryption"`  // 是否启用加密
	Connections      int 		   `yaml:"connections"`        //客户端发起的连接数
	SessionTimeout   time.Duration `yaml:"session_timeout"`    // 会话空闲超时时间
	SequenceTimeout  time.Duration `yaml:"sequence_timeout"`   // 消息排序缓冲区中消息的超时时间
	WSStatsInterval  time.Duration `yaml:"ws_stats_interval"`  // WebSocket统计输出间隔
}

// NewDefaultConfig 返回一个默认配置实例
func NewDefaultConfig() *Config {
	return &Config{
		Mode:                 "client",
		LogLevel:             "info",
		ClientListenAddr:     "127.0.0.1:1080",
		StatusServer:         "127.0.0.1:8010",
		StatusACL:            "",
		ServerURLs:           []string{"ws://127.0.0.1:8081/tunnel"},
		DialTimeout:          5 * time.Second,
		PingInterval:         10 * time.Second,
		PongTimeout:          30 * time.Second,
		ReconnectDelay:       5 * time.Second,
		MaxRetries:           5,
		ClientTLSCAFile:      "",
		ClientTLSCertFile:    "",
		ClientTLSKeyFile:     "",
		ClientInsecureSkipVerify: false,
		ServerListenAddrs:    []string{"0.0.0.0:8081"},
		Password:             "PASSWORD",
		TargetAddr:           "",
		TargetConnectTimeout: 10 * time.Second,
		ServerTLSCertFile:    "",
		ServerTLSKeyFile:     "",
		AuthKey:              "your_secret_auth_key", // 默认认证密钥，请修改
		EnableEncryption:     true,
		Connections:          1,
		SessionTimeout:       5 * time.Minute,
		SequenceTimeout:      1 * time.Second,
		WSStatsInterval:      10 * time.Second, // 默认10秒输出一次WS统计
	}
}

// LoadConfig 从命令行参数和YAML文件加载配置
// 优先级：默认值 -> YAML文件 -> 命令行参数
func LoadConfig() (*Config, error) {
	cfg := NewDefaultConfig() // 1. 加载默认配置

	// 创建一个临时的FlagSet来解析命令行参数
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	var (
		configFileArg          string
		modeArg                string
		logLevelArg            string
		clientListenAddrArg    string
		statusServerArg        string
		statusACLArg           string
		serverURLsArg          string // 命令行参数的 server-urls 作为字符串处理
		dialTimeoutArg         time.Duration
		pingIntervalArg        time.Duration
		pongTimeoutArg         time.Duration
		reconnectDelayArg      time.Duration
		maxRetriesArg          int
		clientTLSCAFileArg     string
		clientTLSCertFileArg   string
		clientTLSKeyFileArg    string
		clientInsecureSkipVerifyArg bool
		serverListenAddrArg    string
		passwordArg              string
		targetAddrArg          string
		targetConnectTimeoutArg time.Duration
		serverTLSCertFileArg   string
		serverTLSKeyFileArg    string
		authKeyArg             string
		enableEncryptionArg    bool
		connectionsArg         int
		sessionTimeoutArg      time.Duration
		sequenceTimeoutArg     time.Duration
		wsStatsIntervalArg     time.Duration
	)

	// 定义所有命令行参数，将它们绑定到临时变量
	// 注意：这里将flag的默认值设为空字符串/0/false，以便通过flagsSeen map判断是否被设置
	fs.StringVar(&configFileArg, "c", "", "Path to configuration file (YAML)")
	fs.StringVar(&modeArg, "mode", "", "Run mode: 'client' or 'server'")
	fs.StringVar(&logLevelArg, "log-level", "", "Log level: 'debug', 'info', 'warn', 'error'")
	fs.StringVar(&clientListenAddrArg, "client-listen", "", "Client local TCP listen address")
	fs.StringVar(&statusServerArg, "status-server", "", "Sataus server listen address")
	fs.StringVar(&statusACLArg, "status-acl", "", "Status server ACL")
	fs.StringVar(&serverURLsArg, "server-urls", "", "Comma-separated list of server WebSocket URLs")
	fs.DurationVar(&dialTimeoutArg, "dial-timeout", 0, "Dial timeout for WebSocket connections")
	fs.DurationVar(&pingIntervalArg, "ping-interval", 0, "Interval for sending ping messages")
	fs.DurationVar(&pongTimeoutArg, "pong-timeout", 0, "Timeout for receiving pong messages")
	fs.DurationVar(&reconnectDelayArg, "reconnect-delay", 0, "Delay before reconnecting")
	fs.IntVar(&maxRetriesArg, "max-retries", 0, "Max connection retry attempts")
	fs.StringVar(&clientTLSCAFileArg, "client-tls-ca", "", "Client TLS CA file path")
	fs.StringVar(&clientTLSCertFileArg, "client-tls-cert", "", "Client TLS certificate file path")
	fs.StringVar(&clientTLSKeyFileArg, "client-tls-key", "", "Client TLS key file path")
	fs.BoolVar(&clientInsecureSkipVerifyArg, "client-insecure-skip-verify", false, "Client skip server certificate verification")
	fs.StringVar(&serverListenAddrArg, "server-listen", "", "Comma-separated list of Server WebSocket listen address")
	fs.DurationVar(&targetConnectTimeoutArg, "target-connect-timeout", 0, "Timeout for server connecting to target")
	fs.StringVar(&serverTLSCertFileArg, "server-tls-cert", "", "Server TLS certificate file path")
	fs.StringVar(&serverTLSKeyFileArg, "server-tls-key", "", "Server TLS key file path")
	fs.StringVar(&passwordArg, "password", "PASSWORD", "Password")
	fs.StringVar(&targetAddrArg, "target-addr", "", "target service address")
	fs.StringVar(&authKeyArg, "auth-key", "", "Authentication key for WebSocket connections")
	fs.BoolVar(&enableEncryptionArg, "enable-encryption", false, "Enable AES256 encryption")
	fs.IntVar(&connectionsArg, "connections", 1, "concurrent connections by client")
	fs.IntVar(&paddingMinArg, "padding-min", 0, "Min random padding length")
	fs.IntVar(&paddingMaxArg, "padding-max", 0, "Max random padding length")
	fs.DurationVar(&sessionTimeoutArg, "session-timeout", 0, "Session idle timeout")
	fs.DurationVar(&sequenceTimeoutArg, "sequence-timeout", 0, "Sequence buffer message timeout")
	fs.DurationVar(&wsStatsIntervalArg, "ws-stats-interval", 0, "Interval for logging WebSocket stats")

	fs.Parse(os.Args[1:]) // 2. 解析命令行参数，它们会填充到上面的临时变量中

	// 记录哪些命令行参数被显式设置了
	flagsSeen := make(map[string]bool)
	fs.Visit(func(f *flag.Flag) {
		flagsSeen[f.Name] = true
	})

	// 3. 读取YAML文件并合并 (YAML覆盖默认值)
	configFilePath := configFileArg
	if configFilePath == "" {
		// 尝试在默认路径下查找名为 "config.yaml" 的文件
		homeDir, _ := os.UserHomeDir()
		possiblePaths := []string{
			"config.yaml", // 当前目录
			fmt.Sprintf("%s/.go-tunnel/config.yaml", homeDir), // 用户主目录
			"/etc/go-tunnel/config.yaml",                      // Linux 系统常见配置目录
		}
		for _, p := range possiblePaths {
			if _, err := os.Stat(p); err == nil {
				configFilePath = p
				break
			}
		}
	}

	if configFilePath != "" {
		fileContent, err := os.ReadFile(configFilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file %s: %w", configFilePath, err)
		}
		// Unmarshal directly into `cfg` (它已经包含默认值)
		if err := yaml.Unmarshal(fileContent, cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file %s: %w", configFilePath, err)
		}
	} else {
		fmt.Printf("Warning: No config file found. Using defaults and command-line flags.\n")
	}

	// 4. 命令行参数覆盖YAML文件和默认值
	// 检查每个参数是否在命令行中被显式设置了，如果是，则用命令行值覆盖
	if flagsSeen["mode"] { cfg.Mode = modeArg }
	if flagsSeen["log-level"] { cfg.LogLevel = logLevelArg }
	if flagsSeen["client-listen"] { cfg.ClientListenAddr = clientListenAddrArg }
	if flagsSeen["status-server"] { cfg.StatusServer = statusServerArg }
	if flagsSeen["status-acl"] { cfg.StatusACL = statusACLArg }
	if flagsSeen["server-urls"] { cfg.ServerURLs = splitAndTrim(serverURLsArg, ",") } // 命令行参数需要特殊处理
	if flagsSeen["dial-timeout"] { cfg.DialTimeout = dialTimeoutArg }
	if flagsSeen["ping-interval"] { cfg.PingInterval = pingIntervalArg }
	if flagsSeen["pong-timeout"] { cfg.PongTimeout = pongTimeoutArg }
	if flagsSeen["reconnect-delay"] { cfg.ReconnectDelay = reconnectDelayArg }
	if flagsSeen["max-retries"] { cfg.MaxRetries = maxRetriesArg }
	if flagsSeen["client-tls-ca"] { cfg.ClientTLSCAFile = clientTLSCAFileArg }
	if flagsSeen["client-tls-cert"] { cfg.ClientTLSCertFile = clientTLSCertFileArg }
	if flagsSeen["client-tls-key"] { cfg.ClientTLSKeyFile = clientTLSKeyFileArg }
	if flagsSeen["client-insecure-skip-verify"] { cfg.ClientInsecureSkipVerify = clientInsecureSkipVerifyArg }
	if flagsSeen["server-listen"] { cfg.ServerListenAddrs = splitAndTrim(serverListenAddrArg, ",") }
	if flagsSeen["target-connect-timeout"] { cfg.TargetConnectTimeout = targetConnectTimeoutArg }
	if flagsSeen["server-tls-cert"] { cfg.ServerTLSCertFile = serverTLSCertFileArg }
	if flagsSeen["server-tls-key"] { cfg.ServerTLSKeyFile = serverTLSKeyFileArg }
	if flagsSeen["ws-path"] { cfg.WSPath = wsPathArg }
	if flagsSeen["target-addr"] { cfg.TargetAddr = targetAddrArg }
	if flagsSeen["auth-key"] { cfg.AuthKey = authKeyArg }
	if flagsSeen["enable-encryption"] { cfg.EnableEncryption = enableEncryptionArg }
	if flagsSeen["connections"] { cfg.Connections = connectionsArg }
	if flagsSeen["padding-min"] { cfg.PaddingMin = paddingMinArg }
	if flagsSeen["padding-max"] { cfg.PaddingMax = paddingMaxArg }
	if flagsSeen["session-timeout"] { cfg.SessionTimeout = sessionTimeoutArg }
	if flagsSeen["sequence-timeout"] { cfg.SequenceTimeout = sequenceTimeoutArg }
	if flagsSeen["ws-stats-interval"] { cfg.WSStatsInterval = wsStatsIntervalArg }

	return cfg, nil
}


*/

func parseFlags() (pass bool, mode int, password string, addr []string, listen string, dest, authkey, cafile, keyfile, crtfile string, tunnelN int, verbose int) {
	var modeString string
	var rabbitaddr string
	var printVersion bool
	flag.StringVar(&modeString, "mode", "c", "running mode(s or c)")
	flag.StringVar(&password, "password", DefaultPassword, "password")
	flag.StringVar(&rabbitaddr, "rabbit-addr", ":443", "listen(server mode) or remote(client mode) address used by rabbit-tcp, eg: 192.168.1.10:22222,192.168.1.11:22223,192.168.1.12:22224")
	flag.StringVar(&listen, "listen", "", "[Client Only] listen address, eg: 127.0.0.1:2333")
	flag.StringVar(&dest, "dest", "", "[Client Only] destination address, eg: shadowsocks server address")
	flag.StringVar(&authkey, "authkey", "", "websocket authkey, eg: mysecret")
	flag.StringVar(&cafile, "tlsca", "", "[Client Only] TLS CA file path, eg: /root/client.ca")
	flag.StringVar(&keyfile, "tlskey", "", "[Server Only] TLS key file path, eg: /root/server.key")
	flag.StringVar(&crtfile, "tlscrt", "", "[Server Only] TLS crt file path, eg: /root/server.crt")
	flag.IntVar(&tunnelN, "tunnelN", 4, "[Client Only] number of tunnels to use in rabbit-tcp")
	flag.IntVar(&verbose, "verbose", 2, "verbose level(0~5)")
	flag.BoolVar(&printVersion, "version", false, "show version")
	flag.Parse()

	pass = true

	// version
	if printVersion {
		log.Println("Rabbit TCP (https://github.com/aagun1234/rabbit-mtcp-ws/)")
		log.Printf("Version: %s.\n", Version)
		pass = false
		return
	}

	// mode
	modeString = strings.ToLower(modeString)
	if modeString == "c" || modeString == "client" {
		mode = ClientMode
	} else if modeString == "s" || modeString == "server" {
		mode = ServerMode
	} else {
		log.Printf("Unsupported mode %s.\n", modeString)
		pass = false
		return
	}

	// password
	if password == "" {
		log.Println("Password must be specified.")
		pass = false
		return
	}
	if password == DefaultPassword {
		log.Println("Password must be changed instead of default password.")
		pass = false
		return
	}

	// listen, dest, tunnelN
	if mode == ClientMode {
		if listen == "" {
			log.Println("Listen address must be specified in client mode.")
			pass = false
		}
		if dest == "" {
			log.Println("Destination address must be specified in client mode.")
			pass = false
		}
		if tunnelN == 0 {
			log.Println("Tunnel number must be positive.")
			pass = false
		}
	}
	
	addr = strings.Split(rabbitaddr, ",")
	return
}

func main() {
	pass, mode, password, addr, listen, dest, authkey, _, keyfile, crtfile, tunnelN, verbose := parseFlags()
	if !pass {
		return
	}
	cipher, _ := tunnel.NewAEADCipher("CHACHA20-IETF-POLY1305", nil, password)
	logger.LEVEL = verbose
	if mode == ClientMode {
		c := client.NewClient(tunnelN, addr, cipher, authkey)
		c.ServeForward(listen, dest)
	} else {
		s := server.NewServer(cipher, authkey, keyfile, crtfile)
		s.Serve(addr)
	}
}

