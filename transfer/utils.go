package transfer

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
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

// localIPv4WithNetwork return only one local network IPv4 addresses by requesting remote network.
//
// ip, port, error
func localIPv4WithNetwork() (string, string, error) {
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

// parseRangeHeader parse incoming client Range header, return start range and end rage
func parseRangeHeader(h string) ([]int64, error) {
	if strings.Compare(h, "") == 0 {
		return []int64{}, errors.New("range header could not be empty")
	}
	h = strings.TrimPrefix(h, "bytes=")
	s := strings.Split(h, "-")
	start, err := strconv.ParseInt(s[0], 10, 0)
	if err != nil {
		return []int64{}, err
	}
	end, err1 := strconv.ParseInt(s[1], 10, 0)
	if err1 != nil {
		return []int64{}, err1
	}
	return []int64{
		start,
		end,
	}, nil
}

// formatRangeHeader format server Content-Range header
func formatRangeHeader(start, end, total int64) string {
	return fmt.Sprintf("bytes %d-%d/%d", start, end, total)
}
