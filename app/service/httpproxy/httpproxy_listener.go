package httpproxy

import (
	"encoding/binary"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"net"
	"runtime"
	"syscall"
)

const (
	// TCP SO_ORIGINAL_DST
	SO_ORIGINAL_DST = 80
)

var (
	// remoteToOriginAddrMap is a concurrent-safe map for remote address to
	// original destination address.
	remoteToOriginAddrMap = gmap.NewStrStrMap(true)
)

// GetOriginByLocalAddr returns the original destination address by local address.
// It returns an empty string if fails.
func GetOriginByRemoteAddr(remoteAddr string) string {
	return remoteToOriginAddrMap.Get(remoteAddr)
}

// Conn implements interface of net.Conn.
type Conn struct {
	net.Conn
}

// Close closes the connection and removes the remote address to original address mapping in memory.
func (c *Conn) Close() error {
	remoteToOriginAddrMap.Remove(c.Conn.RemoteAddr().String())
	return c.Conn.Close()
}

// Listener implements interface of net.Listener.
type Listener struct {
	net.Listener
}

// Accept accepts the connection and adds the remote address to original address mapping in memory.
func (l *Listener) Accept() (c net.Conn, err error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	if origin := l.OriginAddr(conn); origin != "" {
		remoteToOriginAddrMap.Set(conn.RemoteAddr().String(), origin)
	}
	return &Conn{Conn: conn}, nil
}

// OriginAddr retrieves and returns the original address from socket, which is redirected from iptables.
// Note that this function is only available on linux platform.
func (l *Listener) OriginAddr(conn net.Conn) string {
	if runtime.GOOS != "linux" {
		return ""
	}
	tcpConnFile, err := conn.(*net.TCPConn).File()
	if err != nil {
		return ""
	}
	defer tcpConnFile.Close()
	addr, err := syscall.GetsockoptIPv6Mreq(
		int(tcpConnFile.Fd()),
		syscall.IPPROTO_IP,
		SO_ORIGINAL_DST,
	)
	if err != nil {
		g.Log().Debugf("[HTTP] retrieving original address error: %v", err)
		return ""
	}
	return net.IP(addr.Multiaddr[4:8]).String() + ":" + gconv.String(binary.BigEndian.Uint16(addr.Multiaddr[2:4]))
}
