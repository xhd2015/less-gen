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

func IsTCPAddrDialable(url string, timeout time.Duration) (bool, error) {
	conn, err := net.DialTimeout("tcp", url, timeout)
	if err != nil {
		return false, nil
	}
	conn.Close()
	return true, nil
}
