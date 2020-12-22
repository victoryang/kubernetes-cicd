package config

const (
	DefaultEndpoint = ":8080"
)

type k8sCluster struct {
	Name string
	Addr string
}

var (
	Regions = []k8sCluster{
		k8sCluster{Name: "北京星光", Addr: "http://10.12.35.2:8080"},
	}
)

// The defaults applied before parsing the respective config sections.
var (
	DefaultConfig = Config {
		EndPoint: DefaultEndpoint,
		Database: &DefaultDatabase,
		Ldap: &DefaultLdap,
		Log: &DefaultLog,
	}

	DefaultDatabase = DatabaseConfig {
		Adapter: "mysql",
		Username: "root",
		Password: "123456",
	}

	DefaultLdap = Ldap {
		Address: "192.168.x.x:389",
		Password: "xxxxxxx",
	}

	DefaultLog = Log {
		File: "ci.log",
		Level: 0,
	}
)

type Config struct {
	EndPoint 	string 		`yaml:"endpoint,omitempty"`
	Database 	*DatabaseConfig 	`yaml:"database"`
	Ldap 		*Ldap 		`yaml:"ldap,omitempty"`
	Log 		*Log			`yaml:"log,omitempty"`
}

type DatabaseConfig struct {
	Adapter 	string 		`yaml:"adapter"`
	Username 	string 		`yaml:"username"`
	Password 	string 		`yaml:"password"`
}

type Ldap struct {
	Address 	string 		`yaml:"address"`
	Password 	string 		`yaml:"password"`
}

type Log struct {
	File 	string 		`yaml:"file"`
	Level 	int 		`yaml:"level"`
}