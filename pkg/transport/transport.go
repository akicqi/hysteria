package transport

import (
	"github.com/tobyxdd/hysteria/pkg/faketcp"
	"net"
	"time"
)

type Transport interface {
	QUICResolveUDPAddr(address string) (*net.UDPAddr, error)
	QUICListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error)
	QUICListenFakeTCP(address string) (*faketcp.TCPConn, error)
	QUICDialFakeTCP(address string) (*faketcp.TCPConn, error)

	LocalResolveIPAddr(address string) (*net.IPAddr, error)
	LocalResolveTCPAddr(address string) (*net.TCPAddr, error)
	LocalResolveUDPAddr(address string) (*net.UDPAddr, error)
	LocalDial(network, address string) (net.Conn, error)
	LocalDialTCP(laddr, raddr *net.TCPAddr) (*net.TCPConn, error)
	LocalListenTCP(laddr *net.TCPAddr) (*net.TCPListener, error)
	LocalListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error)
}

var DefaultTransport Transport = &defaultTransport{
	Timeout: 8 * time.Second,
}

var IPv6OnlyTransport Transport = &ipv6OnlyTransport{
	defaultTransport{
		Timeout: 8 * time.Second,
	},
}

type defaultTransport struct {
	Timeout time.Duration
}

func (t *defaultTransport) QUICResolveUDPAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}

func (t *defaultTransport) QUICListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error) {
	return net.ListenUDP("udp", laddr)
}

func (t *defaultTransport) QUICListenFakeTCP(address string) (*faketcp.TCPConn, error) {
	return faketcp.Listen("tcp", address)
}

func (t *defaultTransport) QUICDialFakeTCP(address string) (*faketcp.TCPConn, error) {
	return faketcp.Dial("tcp", address)
}

func (t *defaultTransport) LocalResolveIPAddr(address string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip", address)
}

func (t *defaultTransport) LocalResolveTCPAddr(address string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr("tcp", address)
}

func (t *defaultTransport) LocalResolveUDPAddr(address string) (*net.UDPAddr, error) {
	return net.ResolveUDPAddr("udp", address)
}

func (t *defaultTransport) LocalDial(network, address string) (net.Conn, error) {
	dialer := &net.Dialer{Timeout: t.Timeout}
	return dialer.Dial(network, address)
}

func (t *defaultTransport) LocalDialTCP(laddr, raddr *net.TCPAddr) (*net.TCPConn, error) {
	dialer := &net.Dialer{Timeout: t.Timeout, LocalAddr: laddr}
	conn, err := dialer.Dial("tcp", raddr.String())
	if err != nil {
		return nil, err
	}
	return conn.(*net.TCPConn), nil
}

func (t *defaultTransport) LocalListenTCP(laddr *net.TCPAddr) (*net.TCPListener, error) {
	return net.ListenTCP("tcp", laddr)
}

func (t *defaultTransport) LocalListenUDP(laddr *net.UDPAddr) (*net.UDPConn, error) {
	return net.ListenUDP("udp", laddr)
}

type ipv6OnlyTransport struct {
	defaultTransport
}

func (t *ipv6OnlyTransport) LocalResolveIPAddr(address string) (*net.IPAddr, error) {
	return net.ResolveIPAddr("ip6", address)
}
