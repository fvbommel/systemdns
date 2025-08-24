# System DNS resolvers

Some DNS libraries like `github.com/projectdiscovery/retryabledns` require explicitly specifying the DNS resolvers to use.

If you want to use the DNS resolvers configured on the system,
it's actually a bit tricky to figure out what they are in a cross-platform way.

This package implements it in a reusable way.

## Example

<https://github.com/fvbommel/systemdns/blob/1a89e78b0e8b622196654f8f763368b1f381b092/cmd/printdns/printdns.go#L11-L18>

## Under the hood

On Windows systems, we use system calls to figure out the resolvers associated with each network adapter.

On other systems, we look for "nameserver" lines in `/etc/resolv.conf`.

In either case, `GetSystemDNS()` returns a deduplicated list of DNS resolvers.
