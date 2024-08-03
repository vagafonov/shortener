package utils

import "net"

// IsIPInSubnet Check if user ip in subnet.
func IsIPInSubnet(clientIP string, trustedSubnet string) bool {
	clientIPAddr := net.ParseIP(clientIP)

	// Разделяем адрес подсети и маску
	_, subnet, err := net.ParseCIDR(trustedSubnet)
	if err != nil {
		return false
	}

	// Проверяем, принадлежит ли IP-адрес клиента подсети
	return subnet.Contains(clientIPAddr)
}
