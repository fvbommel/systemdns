//go:build !windows

package systemdns

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// getSystemDNS returns system DNS resolvers (IPv4 + IPv6).
func getSystemDNS() ([]string, error) {
	file, err := os.Open("/etc/resolv.conf")
	if err != nil {
		return nil, fmt.Errorf("failed to open /etc/resolv.conf: %w", err)
	}
	defer file.Close()

	seen := make(map[string]struct{})
	var servers []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		fields := strings.Fields(line)
		if len(fields) >= 2 && fields[0] == "nameserver" {
			ns := fields[1]
			if _, ok := seen[ns]; !ok {
				seen[ns] = struct{}{}
				servers = append(servers, ns)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading /etc/resolv.conf: %w", err)
	}

	return servers, nil
}
