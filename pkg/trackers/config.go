package trackers

// Config application's configuration file contents, the global behaviour.
type Config struct {
	Transmission TransmissionConfig `yaml:"transmission"`
	Persistence  PersistenceConfig  `yaml:"persistence"`
	Probe        ProbeConfig        `yaml:"probe"`
	Nameservers  []string           `yaml:"nameservers"`
}

// TransmissionConfig tranmission configuration.
type TransmissionConfig struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// PersistenceConfig persistence configuration.
type PersistenceConfig struct {
	DbPath string `yaml:"dbPath"`
}

// ProbeConfig probe (object) configuration.
type ProbeConfig struct {
	Timeout int64 `yaml:"timeout"`
}
