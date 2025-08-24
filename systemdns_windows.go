package systemdns

import (
	"fmt"
	"net/netip"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

// getSystemDNS returns system DNS resolvers (IPv4 + IPv6).
func getSystemDNS() ([]string, error) {
	var (
		size   uint32
		flags  uint32 = windows.GAA_FLAG_SKIP_ANYCAST | windows.GAA_FLAG_SKIP_MULTICAST
		family uint32 = windows.AF_UNSPEC // both IPv4 and IPv6
	)

	// First call to get required buffer size
	err := windows.GetAdaptersAddresses(family, flags, 0, nil, &size)
	if err != windows.ERROR_BUFFER_OVERFLOW {
		return nil, fmt.Errorf("GetAdaptersAddresses size query failed: %w", err)
	}

	// Make a buffer and fill it.
	buf := make([]byte, size)
	aa := (*windows.IpAdapterAddresses)(unsafe.Pointer(&buf[0]))
	err = windows.GetAdaptersAddresses(family, flags, 0, aa, &size)
	if err != nil {
		return nil, fmt.Errorf("GetAdaptersAddresses failed: %w", err)
	}

	// Extract deduplicated list of DNS resolvers.
	seen := make(map[netip.Addr]struct{})
	var servers []string
	for ; aa != nil; aa = aa.Next {
		for dns := aa.FirstDnsServerAddress; dns != nil; dns = dns.Next {
			sockaddr := dns.Address.Sockaddr

			var ip netip.Addr
			if uintptr(dns.Address.SockaddrLength) >= unsafe.Sizeof(windows.RawSockaddrInet4{}) && sockaddr.Addr.Family == windows.AF_INET {
				raw := (*windows.RawSockaddrInet4)(unsafe.Pointer(sockaddr))
				ip = netip.AddrFrom4(raw.Addr)
			} else if uintptr(dns.Address.SockaddrLength) >= unsafe.Sizeof(windows.RawSockaddrInet6{}) && sockaddr.Addr.Family == windows.AF_INET6 {
				raw := (*windows.RawSockaddrInet6)(unsafe.Pointer(sockaddr))
				ip = netip.AddrFrom16(raw.Addr)
				if raw.Scope_id != 0 {
					ip = ip.WithZone(strconv.FormatUint(uint64(raw.Scope_id), 10))
				}
			} else {
				continue
			}

			if _, ok := seen[ip]; !ok {
				seen[ip] = struct{}{}
				servers = append(servers, ip.String())
			}
		}
	}

	return servers, nil
}
