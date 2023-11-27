package gotools

import "regexp"

// IsIp 是否是合法 ipv4 地址 ipv6待增加
func IsIp(ip string) (b bool) {
	if m, _ := regexp.MatchString("^(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)", ip); !m {
		return false
	}
	return true
}

// IsIpSegment 是否是ip网段格
func IsIpSegment(ip string) (b bool) {
	if m, _ := regexp.MatchString("^((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}/(3[0-2]|[1-2]\\d|\\d)$", ip); !m {
		return false
	}
	return true
}

// DomainCheck 判断是否是一个合法域名
func DomainCheck(domain string) bool {
	var match bool
	NotLine := "^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,6}"
	//支持以http://或者https://开头并且域名中间没有/的情况
	// NotLine := "(http(s)?:\\/\\/)?(www\\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\\d+)*(\\/\\w+\\.\\w+)*$"
	match, _ = regexp.MatchString(NotLine, domain)
	return match
}

func CheckDomain(domain string) bool {
	var match bool
	// 支持以http://或者https://开头并且域名中间没有/的情况
	// NotLine := "^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,6}"
	NotLine := "(http(s)?:\\/\\/)?(www\\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\\d+)*(\\/\\w+\\.\\w+)*$"
	match, _ = regexp.MatchString(NotLine, domain)
	return match
}
