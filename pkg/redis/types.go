package redis

// Config holds Redis connection settings
type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}
