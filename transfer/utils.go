package transfer

import (
	"net"
)

// some range slice s and judge whether filter could shot.
//
// s is any type(interface {}), filter func has three params and need one bool returned,
// current is value of slice piece, all is s itself, index is index number of slice piece.
func some[T any](s []T, filter func(current T, all []T, index int) bool) bool {
	for i, v := range s {
		if filter(v, s, i) {
			return true
		}
	}
	return false
}

// includes range slice s and judge whether v is existed in s.
//
// Be sure that s piece is type comparable(string|int|unit|bool).
func includes[T comparable](s []T, v T) bool {
	return some(s, func(current T, all []T, index int) bool {
		return current == v
	})
}

// splice range slice s and remove item by filter func, it returns an new splice rather than effects s itself.
//
// s is any type(interface {}), filter func has two params and need one bool returned,
// current is value of slice piece, index is index number of slice piece.
func splice[T any](s []T, filter func(current T, index int) bool) []T {
	newS := make([]T, 0)
	for i, v := range s {
		if !filter(v, i) {
			newS = append(newS, v)
		}
	}
	return newS
}

// localIPv4s return all non-loopback IPv4 addresses
func localIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

// remoteIPv4 return public network IPv4 addresses.
//
// ip, port, error
func remoteIPv4() (string, string, error) {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		return "", "", err
	}
	defer conn.Close()
	ip, port, err := net.SplitHostPort(conn.LocalAddr().String())
	if err != nil {
		return "", "", err
	}
	return ip, port, nil
}
