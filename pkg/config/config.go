// package config contains configuration information for `sfid` app
package config

import "runtime"

// Config keeps configuration data.
type Config struct {
	// NameSpace of UUID v5, not used if only Sha256 output is requested.
	NameSpace string

	//OutputSha256 returns Sha256 hash back
	OutputSha256 bool

	// OutputUUID returns UUIDv5 back using either supplied or default UUID
	// NameSpace.
	OutputUUID bool

	// Input is a string that might mean either a string, a file or a directory.
	Input string

	// Recursive sets recursive directory walking and creating output for each
	// file. This setting is ignored if Input is not a directory.
	Recursive bool

	// JobsNum sets provided concurrency number for running the algorithm
	JobsNum int
}

type Option func(c *Config)

func OptNameSpace(s string) Option {
	return func(c *Config) {
		c.NameSpace = s
	}
}

func OptInput(s string) Option {
	return func(c *Config) {
		c.Input = s
	}
}

func OptRecursive(b bool) Option {
	return func(c *Config) {
		c.Recursive = b
	}
}

func OptJobsNum(i int) Option {
	return func(c *Config) {
		c.JobsNum = i
	}
}

func New(opts ...Option) Config {
	res := Config{
		NameSpace: "speciesfilegroup.org",
		JobsNum:   runtime.NumCPU(),
	}

	for _, opt := range opts {
		opt(&res)
	}

	return res
}
