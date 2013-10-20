package  mdns-pub

import (
  "github.com/davecheney/mdns"
  "flag"
  "log"
  "fmt"
  "net"
)

func mustPublish(rr string) {
  if err := mdns.Publish(rr); err != nil {
    log.Fatal(`Unable to publish record "%s": %v`, rr, err)
  }
}

func usage() {
  log.Fatal("Usage: mdns-pub [-t] <address> <name> [service] [port]")
  flag.PrintDefaults()
  os.Exit(2)
}

var expiry int
func init() {
  flag.IntVar( &expiry, "e", 60, "Set the record timeout in seconds"
}

func reverseaddr(addr string) (arpa string, err error) {
  ip := net.ParseIP(addr)
  if ip == nil {
    return "", "unrecognized address"
  }

  if ip.To4() != nil {
    return itoa(int(ip[15])) + "." + itoa(int(ip[14])) + "." + itoa(int(ip[13])) + "." +
      itoa(int(ip[12])) + ".in-addr.arpa.", nil
  }
  // Must be IPv6
  buf := make([]byte, 0, len(ip)*4+len("ip6.arpa."))
  // Add it, in reverse, to the buffer
  for i := len(ip) - 1; i >= 0; i-- {
    v := ip[i]
    buf = append(buf, hexDigit[v&0xF])
    buf = append(buf, '.')
    buf = append(buf, hexDigit[v>>4])
    buf = append(buf, '.')
  }
  // Append "ip6.arpa." and return (buf already has the final .)
  buf = append(buf, "ip6.arpa."...)
  return string(buf), nil
}

func main() {
  flag.Usage = usage
  flag.Parse()

  args := flag.Args()
  if len(args) < 2 {
    fmt.Println("Must specify at least Adress and Name")
    usage()
    os.Exit(1)
  }

  if len(args) == 3 {
    fmt.Println("Must provide a port for service registry")
    usage()
    os.Exit
  }

  address, name := args[0], args[1]
  if ! net.ParseIP(address) {
    fmt.Println("Adress is not invalid: ", address)
  }

  mustPublish(name + ". " + timeout + " IN A " + address)
  mustPublish(reverseaddr(address) + " " + timeout + " IN PTR " + name + ".")


  /* A simple example. Publish an A record for my router at 192.168.1.254.
  mustPublish("router.local. 60 IN A 192.168.1.254")
  mustPublish("254.1.168.192.in-addr.arpa. 60 IN PTR router.local.")

  // A more compilcated example. Publish a SVR record for ssh running on port
  // 22 for my home NAS.

  // Publish an A record as before
  mustPublish("stora.local. 60 IN A 192.168.1.200")
  mustPublish("200.1.168.192.in-addr.arpa. 60 IN PTR stora.local.")

  // Publish a PTR record for the _ssh._tcp DNS-SD type
  mustPublish("_ssh._tcp.local. 60 IN PTR stora._ssh._tcp.local.")

  // Publish a SRV record tying the _ssh._tcp record to an A record and a port.
  mustPublish("stora._ssh._tcp.local. 60 IN SRV 0 0 22 stora.local.")

  // Most mDNS browsing tools expect a TXT record for the service even if there
  // are not records defined by RFC 2782.
  mustPublish(`stora._ssh._tcp.local. 60 IN TXT ""`)

  // Bind this service into the list of registered services for dns-sd.
  mustPublish("_service._dns-sd._udp.local. 60 IN PTR _ssh._tcp.local.")
*/
  select {}
}
