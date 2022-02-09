package stdgolibs

import (
	pkg "net"

	"reflect"
)

func init() {
	registerValues("net", map[string]reflect.Value{
		// Functions
		"ResolveTCPAddr":     reflect.ValueOf(pkg.ResolveTCPAddr),
		"DialTCP":            reflect.ValueOf(pkg.DialTCP),
		"ListenTCP":          reflect.ValueOf(pkg.ListenTCP),
		"ResolveUDPAddr":     reflect.ValueOf(pkg.ResolveUDPAddr),
		"DialUDP":            reflect.ValueOf(pkg.DialUDP),
		"ListenUDP":          reflect.ValueOf(pkg.ListenUDP),
		"ListenMulticastUDP": reflect.ValueOf(pkg.ListenMulticastUDP),
		"Interfaces":         reflect.ValueOf(pkg.Interfaces),
		"InterfaceAddrs":     reflect.ValueOf(pkg.InterfaceAddrs),
		"InterfaceByIndex":   reflect.ValueOf(pkg.InterfaceByIndex),
		"InterfaceByName":    reflect.ValueOf(pkg.InterfaceByName),
		"Pipe":               reflect.ValueOf(pkg.Pipe),
		"Dial":               reflect.ValueOf(pkg.Dial),
		"DialTimeout":        reflect.ValueOf(pkg.DialTimeout),
		"Listen":             reflect.ValueOf(pkg.Listen),
		"ListenPacket":       reflect.ValueOf(pkg.ListenPacket),
		"FileConn":           reflect.ValueOf(pkg.FileConn),
		"FileListener":       reflect.ValueOf(pkg.FileListener),
		"FilePacketConn":     reflect.ValueOf(pkg.FilePacketConn),
		"LookupHost":         reflect.ValueOf(pkg.LookupHost),
		"LookupIP":           reflect.ValueOf(pkg.LookupIP),
		"LookupPort":         reflect.ValueOf(pkg.LookupPort),
		"LookupCNAME":        reflect.ValueOf(pkg.LookupCNAME),
		"LookupSRV":          reflect.ValueOf(pkg.LookupSRV),
		"LookupMX":           reflect.ValueOf(pkg.LookupMX),
		"LookupNS":           reflect.ValueOf(pkg.LookupNS),
		"LookupTXT":          reflect.ValueOf(pkg.LookupTXT),
		"LookupAddr":         reflect.ValueOf(pkg.LookupAddr),
		"ResolveUnixAddr":    reflect.ValueOf(pkg.ResolveUnixAddr),
		"DialUnix":           reflect.ValueOf(pkg.DialUnix),
		"ListenUnix":         reflect.ValueOf(pkg.ListenUnix),
		"ListenUnixgram":     reflect.ValueOf(pkg.ListenUnixgram),
		"IPv4":               reflect.ValueOf(pkg.IPv4),
		"IPv4Mask":           reflect.ValueOf(pkg.IPv4Mask),
		"CIDRMask":           reflect.ValueOf(pkg.CIDRMask),
		"ParseIP":            reflect.ValueOf(pkg.ParseIP),
		"ParseCIDR":          reflect.ValueOf(pkg.ParseCIDR),
		"ParseMAC":           reflect.ValueOf(pkg.ParseMAC),
		"ResolveIPAddr":      reflect.ValueOf(pkg.ResolveIPAddr),
		"DialIP":             reflect.ValueOf(pkg.DialIP),
		"ListenIP":           reflect.ValueOf(pkg.ListenIP),
		"SplitHostPort":      reflect.ValueOf(pkg.SplitHostPort),
		"JoinHostPort":       reflect.ValueOf(pkg.JoinHostPort),

		// Consts

		"FlagUp":           reflect.ValueOf(pkg.FlagUp),
		"FlagBroadcast":    reflect.ValueOf(pkg.FlagBroadcast),
		"FlagLoopback":     reflect.ValueOf(pkg.FlagLoopback),
		"FlagPointToPoint": reflect.ValueOf(pkg.FlagPointToPoint),
		"FlagMulticast":    reflect.ValueOf(pkg.FlagMulticast),
		"IPv4len":          reflect.ValueOf(pkg.IPv4len),
		"IPv6len":          reflect.ValueOf(pkg.IPv6len),

		// Variables

		"ErrWriteToConnected":        reflect.ValueOf(&pkg.ErrWriteToConnected),
		"ErrClosed":                  reflect.ValueOf(&pkg.ErrClosed),
		"DefaultResolver":            reflect.ValueOf(&pkg.DefaultResolver),
		"IPv4bcast":                  reflect.ValueOf(&pkg.IPv4bcast),
		"IPv4allsys":                 reflect.ValueOf(&pkg.IPv4allsys),
		"IPv4allrouter":              reflect.ValueOf(&pkg.IPv4allrouter),
		"IPv4zero":                   reflect.ValueOf(&pkg.IPv4zero),
		"IPv6zero":                   reflect.ValueOf(&pkg.IPv6zero),
		"IPv6unspecified":            reflect.ValueOf(&pkg.IPv6unspecified),
		"IPv6loopback":               reflect.ValueOf(&pkg.IPv6loopback),
		"IPv6interfacelocalallnodes": reflect.ValueOf(&pkg.IPv6interfacelocalallnodes),
		"IPv6linklocalallnodes":      reflect.ValueOf(&pkg.IPv6linklocalallnodes),
		"IPv6linklocalallrouters":    reflect.ValueOf(&pkg.IPv6linklocalallrouters),
	})
	registerTypes("net", map[string]reflect.Type{
		// Non interfaces

		"TCPAddr":             reflect.TypeOf((*pkg.TCPAddr)(nil)).Elem(),
		"TCPConn":             reflect.TypeOf((*pkg.TCPConn)(nil)).Elem(),
		"TCPListener":         reflect.TypeOf((*pkg.TCPListener)(nil)).Elem(),
		"UDPAddr":             reflect.TypeOf((*pkg.UDPAddr)(nil)).Elem(),
		"UDPConn":             reflect.TypeOf((*pkg.UDPConn)(nil)).Elem(),
		"Interface":           reflect.TypeOf((*pkg.Interface)(nil)).Elem(),
		"Flags":               reflect.TypeOf((*pkg.Flags)(nil)).Elem(),
		"OpError":             reflect.TypeOf((*pkg.OpError)(nil)).Elem(),
		"ParseError":          reflect.TypeOf((*pkg.ParseError)(nil)).Elem(),
		"AddrError":           reflect.TypeOf((*pkg.AddrError)(nil)).Elem(),
		"UnknownNetworkError": reflect.TypeOf((*pkg.UnknownNetworkError)(nil)).Elem(),
		"InvalidAddrError":    reflect.TypeOf((*pkg.InvalidAddrError)(nil)).Elem(),
		"DNSConfigError":      reflect.TypeOf((*pkg.DNSConfigError)(nil)).Elem(),
		"DNSError":            reflect.TypeOf((*pkg.DNSError)(nil)).Elem(),
		"Buffers":             reflect.TypeOf((*pkg.Buffers)(nil)).Elem(),
		"Dialer":              reflect.TypeOf((*pkg.Dialer)(nil)).Elem(),
		"ListenConfig":        reflect.TypeOf((*pkg.ListenConfig)(nil)).Elem(),
		"SRV":                 reflect.TypeOf((*pkg.SRV)(nil)).Elem(),
		"MX":                  reflect.TypeOf((*pkg.MX)(nil)).Elem(),
		"NS":                  reflect.TypeOf((*pkg.NS)(nil)).Elem(),
		"Resolver":            reflect.TypeOf((*pkg.Resolver)(nil)).Elem(),
		"UnixAddr":            reflect.TypeOf((*pkg.UnixAddr)(nil)).Elem(),
		"UnixConn":            reflect.TypeOf((*pkg.UnixConn)(nil)).Elem(),
		"UnixListener":        reflect.TypeOf((*pkg.UnixListener)(nil)).Elem(),
		"IP":                  reflect.TypeOf((*pkg.IP)(nil)).Elem(),
		"IPMask":              reflect.TypeOf((*pkg.IPMask)(nil)).Elem(),
		"IPNet":               reflect.TypeOf((*pkg.IPNet)(nil)).Elem(),
		"HardwareAddr":        reflect.TypeOf((*pkg.HardwareAddr)(nil)).Elem(),
		"IPAddr":              reflect.TypeOf((*pkg.IPAddr)(nil)).Elem(),
		"IPConn":              reflect.TypeOf((*pkg.IPConn)(nil)).Elem(),
	})
}
