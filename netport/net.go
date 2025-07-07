package netport

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func FindListenablePort(host string, port int, limit int) (int, error) {
	for i := 0; limit <= 0 || i < limit; i++ {
		addr := net.JoinHostPort(host, strconv.Itoa(port))
		serving, err := IsTCPAddrDialable(addr, 20*time.Millisecond)
		if err != nil {
			return 0, err
		}
		if serving {
			port++
			continue
		}

		return port, nil
	}
	return 0, fmt.Errorf("no available port found")
}

// addr example:
//   - 127.0.0.1:7070
func IsTCPAddrDialable(addr string, timeout time.Duration) (bool, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false, nil
	}
	conn.Close()
	return true, nil
}

func WaitTCPAddrDialable(addr string, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return err
	}
	conn.Close()
	return nil
}
