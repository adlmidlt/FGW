package pkg

import (
	"net"
	"os"
	"strings"
)

// InfoPC информация ПК.
type InfoPC struct {
	domain string
	ipAddr string
}

// NewInfoPC создает новый экземпляр InfoPC и инициализирует его поля.
func NewInfoPC() *InfoPC {
	hostname := hostName()
	ipAddr := ipAddress()

	return &InfoPC{domain: hostname, ipAddr: ipAddr}
}

// IPAddr возвращает IP-адрес.
func (i *InfoPC) IPAddr() string {
	return i.ipAddr
}

// HostName возвращает имя хоста (домен), связанный с экземпляром InfoPC.
func (i *InfoPC) HostName() string {
	return i.domain
}

// hostName получить текущее имя хоста.
func hostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "hostname не найден"
	}

	return hostname
}

// ipAddress выполняет поиск IP-адресов для текущего хоста, используя его имя (hostName), объединяем их в одну строку,
// разделенную символом `|`, и возвращает эту строку.
func ipAddress() string {
	ips, err := net.LookupIP(hostName())
	if err != nil {
		return "IP Addresses не найдены"
	}

	var ipAddresses strings.Builder
	for _, ip := range ips {
		ipAddresses.WriteString(ip.String())
		ipAddresses.WriteString(" | ")
	}

	allIPAddrInLine := ipAddresses.String()
	if len(allIPAddrInLine) > 0 {
		allIPAddrInLine = allIPAddrInLine[:len(allIPAddrInLine)-3]
	}

	return allIPAddrInLine
}
