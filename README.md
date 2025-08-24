# System DNS resolvers

Some DNS libraries like `github.com/projectdiscovery/retryabledns` require explicitly specifying the DNS resolvers to use.

If you want to use the DNS resolvers configured on the system,
it's actually a bit tricky to figure out what they are in a cross-platform way.

This package implements it in a reusable way.

## Test

<https://github.com/fvbommel/systemdns/blob/main/cmd/printdns/printdns.go>

## Example

```go
package main

import (
    "fmt"
    "log"

    "github.com/fvbommel/systemdns"
)

func main() {
    resolvers, err := systemdns.GetSystemDNS()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Resolvers: %v\n", resolvers)
}
```

## Under the hood

On Windows systems, we use system calls to figure out the resolvers associated with each network adapter.

On other systems, we look for "nameserver" lines in `/etc/resolv.conf`.

In either case, `GetSystemDNS()` returns a deduplicated list of DNS resolvers.
