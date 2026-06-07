package config

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ListenAddr      string
	FrontendOrigins []string
	MySQLDSN        string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	TokenSecret     string
	SessionTTL      time.Duration
	VerificationTTL time.Duration
	CodeCooldown    time.Duration
	CodeMaxAttempts int
	CodeAttemptTTL  time.Duration
	CookieName      string
	CookieSecure    bool
	SMTPHost        string
	SMTPPort        int
	SMTPUsername    string
	SMTPPassword    string
	SMTPFrom        string
}

func Load() Config {
	loadDotEnv()

	return Config{
		ListenAddr:      listenAddr(),
		FrontendOrigins: splitCSV(getenvFallback("CORS_ALLOWED_ORIGINS", getenvFallback("FRONTEND_ORIGIN", "http://localhost:5173"))),
		MySQLDSN:        os.Getenv("MYSQL_DSN"),
		RedisAddr:       getenvFallback("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:         getenvIntDefault("REDIS_DB", 0),
		TokenSecret:     getenvFallback("TOKEN_SECRET", "change-me"),
		SessionTTL:      getenvDurationDefault("SESSION_TTL_SECONDS", 7*24*time.Hour),
		VerificationTTL: getenvDurationDefault("VERIFICATION_TTL_SECONDS", 5*time.Minute),
		CodeCooldown:    getenvDurationDefault("CODE_COOLDOWN_SECONDS", time.Minute),
		CodeMaxAttempts: getenvIntDefault("CODE_MAX_ATTEMPTS", 5),
		CodeAttemptTTL:  getenvDurationDefault("CODE_ATTEMPT_TTL_SECONDS", 15*time.Minute),
		CookieName:      getenvFallback("SESSION_COOKIE_NAME", "e_director_session"),
		CookieSecure:    getenvBoolDefault("SESSION_COOKIE_SECURE", false),
		SMTPHost:        getenvFallback("SMTP_HOST", "smtp.qq.com"),
		SMTPPort:        getenvIntDefault("SMTP_PORT", 465),
		SMTPUsername:    os.Getenv("SMTP_USERNAME"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		SMTPFrom:        getenvFallback("SMTP_FROM", os.Getenv("SMTP_USERNAME")),
	}
}

func loadDotEnv() {
	for _, path := range dotenvCandidates() {
		if loadDotEnvFile(path) == nil {
			return
		}
	}
}

func dotenvCandidates() []string {
	cwd, err := os.Getwd()
	if err != nil {
		return []string{".env"}
	}

	candidates := []string{filepath.Join(cwd, ".env")}
	for dir := cwd; ; {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		candidates = append(candidates, filepath.Join(parent, ".env"))
		dir = parent
	}
	return candidates
}

func loadDotEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		if key == "" || os.Getenv(key) != "" {
			continue
		}
		os.Setenv(key, trimEnvValue(value))
	}
	return scanner.Err()
}

func trimEnvValue(value string) string {
	value = strings.TrimSpace(value)
	if len(value) >= 2 {
		quote := value[0]
		if (quote == '\'' || quote == '"') && value[len(value)-1] == quote {
			return value[1 : len(value)-1]
		}
	}
	return value
}

func listenAddr() string {
	if port := os.Getenv("HTTP_PORT"); port != "" {
		if strings.HasPrefix(port, ":") {
			return port
		}
		return ":" + port
	}
	return getenvFallback("GATEWAY_ADDR", ":8080")
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	items := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}

func getenvFallback(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvIntDefault(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func getenvBoolDefault(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed
		}
	}
	return fallback
}

func getenvDurationDefault(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return time.Duration(n) * time.Second
		}
	}
	return fallback
}
