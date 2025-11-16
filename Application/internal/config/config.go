package config

type Config struct {
	StorageType string // "file" or "db"
	UserFile    string
	SessionFile string
	DBURL       string
}

var AppConfig = Config{
	StorageType: "db",                             // Change to "db" when needed
	DBURL:       "postgres://user:pass@localhost", // Placeholder (not used for file mode)
}

const (
	Bold          = "\033[1m"
	StrikeThrough = "\033[9m"
	Reset         = "\033[0m"
	Cyan          = "\033[36m"
	Green         = "\033[32m"
	Yellow        = "\033[33m"
	Magenta       = "\033[35m"
	Red           = "\033[31m"
)
