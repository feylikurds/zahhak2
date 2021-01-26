/*
Zahhak2, a Golang multiplayer console game.
Copyright (C) 2016 Aryo Pehlewan feylikurds@gmail.com
This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.
This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.
You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package networking

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func RemoteIP() (string, error) {
	res, e := http.Get("http://myexternalip.com/raw")
	defer res.Body.Close()

	if e == nil {
		content, err := ioutil.ReadAll(res.Body)
		ip := string(content)

		return strings.TrimSpace(ip), err
	}

	return "", e
}

func RemoteHostnames() ([]string, error) {
	res, err := http.Get("http://myexternalip.com/raw")
	defer res.Body.Close()

	if err == nil {
		content, err := ioutil.ReadAll(res.Body)
		ip := string(content)

		if err != nil {
			return nil, err
		}

		hostnames, e := net.LookupAddr(ip)

		return hostnames, e
	}

	return nil, err
}

func LocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return strings.TrimSpace(ipnet.IP.String()), nil
			}
		}
	}

	return "", errors.New("cannot find local address")
}

func LocalHostnames() ([]string, error) {
	lip, err := LocalIP()

	if err == nil {
		addr, e := net.LookupAddr(lip)

		return addr, e
	}

	return nil, err
}

func ValidIP4(ipAddress string) bool {
	ipAddress = strings.Trim(ipAddress, " ")

	re, _ := regexp.Compile(`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)
	if re.MatchString(ipAddress) {
		return true
	}
	return false
}

func IsRemotePortOpen(ipAddress string, port string) bool {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ipAddress+":"+port)

	if err != nil {
		return false
	}

	conn, err := net.DialTimeout("tcp", tcpAddr.String(), 1*time.Minute)

	if err != nil {
		return false
	}

	defer conn.Close()

	return true
}

func IsTCPPortAvailable(port int) bool {
	if port < 1 || port > 65534 {
		return false
	}
	conn, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
