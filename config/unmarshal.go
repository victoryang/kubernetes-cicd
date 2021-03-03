package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

var (
	Configuration *Config
)

func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Config

	if err := unmarshal((*plain)(c)); err!=nil {
		return err
	}

	if c.EndPoint == "" {
		c.EndPoint = DefaultConfig.EndPoint
	}

	if c.Database == nil {
		c.Database = DefaultConfig.Database
	}

	if c.Ldap == nil {
		c.Ldap = DefaultConfig.Ldap
	}

	if c.Log == nil {
		c.Log = DefaultConfig.Log
	}

	if c.DebugMode == false {
		c.DebugMode = DefaultConfig.DebugMode
	}

	if c.GithubSecret == "" {
		c.GithubSecret = DefaultConfig.GithubSecret
	}

	return nil
}

func (dc *DatabaseConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain DatabaseConfig

	if err := unmarshal((*plain)(dc)); err!=nil {
		return err
	}

	if dc.Adapter == "" {
		dc.Adapter = DefaultDatabase.Adapter
	}

	if dc.Username == "" {
		dc.Username = DefaultDatabase.Username
	}

	if dc.Password == "" {
		dc.Password = DefaultDatabase.Password
	}

	return nil
}

func (l *Ldap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Ldap

	if err := unmarshal((*plain)(l)); err!=nil {
		return err
	}

	if l.Address == "" {
		l.Address = DefaultLdap.Address
	}

	if l.Password == "" {
		l.Password = DefaultLdap.Password
	}

	return nil
}

func (l *Log) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type plain Log

	if err := unmarshal((*plain)(l)); err!=nil {
		return err
	}

	if l.File == "" {
		l.File = DefaultLog.File
	}

	if l.Level == 0 {
		l.Level = DefaultLog.Level
	}

	return nil
}

// Load parses the YAML input s into a Config.
func Load(s string) (*Config, error) {
	cfg := &Config{}

	*cfg = DefaultConfig

	err := yaml.UnmarshalStrict([]byte(s), cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*Config, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg, err := Load(string(content))
	if err != nil {
		return nil, err
	}

	Configuration = cfg

	return cfg, nil
}
