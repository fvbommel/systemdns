package systemdns

import (
	"sync"
)

var (
	// Used to look up system DNS resolvers only once.
	systemDnsOnce sync.Once

	// Cached list of system DNS resolvers, or the error to return instead.
	systemDnsResolvers []string
	systemDnsError     error
)

// GetSystemDNS returns system DNS resolvers (IPv4 + IPv6).
func GetSystemDNS() ([]string, error) {
	systemDnsOnce.Do(func() {
		systemDnsResolvers, systemDnsError = getSystemDNS()
	})

	return systemDnsResolvers, systemDnsError
}
