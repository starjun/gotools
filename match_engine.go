package gotools

import (
	"fmt"
	mapset "github.com/deckarep/golang-set"
	"github.com/malfunkt/iprange"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// MatchStr -- 规则判断
// sourceStr 被匹配字符串s
// opt 匹配方式【等于、包含、前缀、后缀、cidr、正则 ...】
// []eStr 匹配字符串s
func MatchStr(_sourceStr string, e_str []string, _opt string, isNot bool) bool {
	re := false
	if isNot {
		re = true
	}
	if _opt == "=" || _opt == "" {
		for i := 0; i < len(e_str); i++ {
			if e_str[i] == "*" || e_str[i] == _sourceStr {
				return !re
			}
		}
		return re
	} else if _opt == "in" {
		for i := 0; i < len(e_str); i++ {
			if strings.Contains(_sourceStr, e_str[i]) {
				return !re
			}
		}
		return re
	} else if _opt == "u_in" {
		u_sourceStr := strings.ToUpper(_sourceStr)
		for i := 0; i < len(e_str); i++ {
			if strings.Contains(u_sourceStr, strings.ToUpper(e_str[i])) {
				return !re
			}

		}
		return re
	} else if _opt == "prefix" {
		for i := 0; i < len(e_str); i++ {
			if strings.HasPrefix(_sourceStr, e_str[i]) {
				return !re
			}
		}
		return re
	} else if _opt == "u_prefix" {
		u_sourceStr := strings.ToUpper(_sourceStr)
		for i := 0; i < len(e_str); i++ {
			if strings.HasPrefix(u_sourceStr, strings.ToUpper(e_str[i])) {
				return !re
			}
		}
		return re
	} else if _opt == "suffix" {
		for i := 0; i < len(e_str); i++ {
			if strings.HasSuffix(_sourceStr, e_str[i]) {
				return !re
			}
		}
		return re
	} else if _opt == "u_suffix" {
		u_sourceStr := strings.ToUpper(_sourceStr)
		for i := 0; i < len(e_str); i++ {
			if strings.HasSuffix(u_sourceStr, strings.ToUpper(e_str[i])) {
				return !re
			}
		}
		return re
	} else if _opt == "cidr" {
		for i := 0; i < len(e_str); i++ {
			if IpCidrCheck(_sourceStr, e_str[i]) {
				return !re
			}
		}
		return re
	} else if _opt == ">" {
		if len(e_str) == 0 {
			return re
		}
		intSource, err := strconv.Atoi(_sourceStr)
		if err != nil {
			return re
		}
		intTarget, err := strconv.Atoi(e_str[0])
		if err != nil {
			return re
		}
		if intSource > intTarget {
			return !re
		}
		return re
	} else if _opt == ">=" {
		if len(e_str) == 0 {
			return re
		}
		intSource, err := strconv.Atoi(_sourceStr)
		if err != nil {
			return re
		}
		intTarget, err := strconv.Atoi(e_str[0])
		if err != nil {
			return re
		}
		if intSource >= intTarget {
			return !re
		}
		return re
	} else if _opt == "<" {
		if len(e_str) == 0 {
			return re
		}
		intSource, err := strconv.Atoi(_sourceStr)
		if err != nil {
			return re
		}
		intTarget, err := strconv.Atoi(e_str[0])
		if err != nil {
			return re
		}
		if intSource < intTarget {
			return !re
		}
		return re
	} else if _opt == "<=" {
		if len(e_str) == 0 {
			return re
		}
		intSource, err := strconv.Atoi(_sourceStr)
		if err != nil {
			return re
		}
		intTarget, err := strconv.Atoi(e_str[0])
		if err != nil {
			return re
		}
		if intSource <= intTarget {
			return !re
		}
		return re
	} else {
		// 正则匹配
		for i := 0; i < len(e_str); i++ {
			match, _ := regexp.MatchString(e_str[i], _sourceStr)
			if match {
				return !re
			}
		}
		return re
	}
}

type CRule struct {
	Opt        string   // 匹配方式
	ReStrList  []string // 匹配字符串
	MaLocation string   // 匹配位置
	Des        string   // 规则描述
	Rev        bool     // 是否取反
	Lcon       string   // 规则连接符
}

func MapCRuleMatch(_obmap map[string]string, _crule CRule) bool {
	return MatchStr(_obmap[_crule.MaLocation], _crule.ReStrList, _crule.Opt, _crule.Rev)
}

func orListMatch(_obmap map[string]string, _crules []CRule) bool {
	for i := 0; i < len(_crules); i++ {
		if MapCRuleMatch(_obmap, _crules[i]) {
			return true
		}
	}
	return false
}

func MapCrulesListMatch(_obmap map[string]string, _crules []CRule) bool {
	var orListCrules []CRule
	cnt := len(_crules)
	for i := 0; i < cnt; i++ {
		if _crules[i].Lcon == "or" {
			orListCrules = append(orListCrules, _crules[i])
			if i == cnt {
				return orListMatch(_obmap, orListCrules)
			}
		} else {
			if len(orListCrules) == 0 {
				if MapCRuleMatch(_obmap, _crules[i]) {
					// nothing todu
				} else {
					return false
				}
			} else {
				orListCrules = append(orListCrules, _crules[i])
				if orListMatch(_obmap, orListCrules) {
					// nothing todu
				} else {
					return false
				}
				orListCrules = nil
			}
		}
	}
	return true
}

// MergeDuplicateIntArray 并集  合并两个数组并去重
func MergeDuplicateIntArray(slice []int, elems []int) []int {
	listPId := append(slice, elems...)
	t := mapset.NewSet()
	for _, i := range listPId {
		t.Add(i)
	}
	var result []int
	for i := range t.Iterator().C {
		result = append(result, i.(int))
	}
	return result
}

// MergeDuplicateStrArray 并集  合并两个字符串数组并去重
func MergeDuplicateStrArray(slice []string, elems []string) []string {
	listPId := append(slice, elems...)
	t := mapset.NewSet()
	for _, i := range listPId {
		t.Add(i)
	}
	var result []string
	for i := range t.Iterator().C {
		result = append(result, i.(string))
	}
	return result
}

// DuplicateIntArray Int数组 去重
func DuplicateIntArray(m []int) []int {
	s := make([]int, 0)
	smap := make(map[int]int)
	for _, value := range m {
		//计算map长度
		length := len(smap)
		smap[value] = 1
		//比较map长度, 如果map长度不相等， 说明key不存在
		if len(smap) != length {
			s = append(s, value)
		}
	}
	return s
}

// DuplicateStrArray 字符串数组 去重
func DuplicateStrArray(m []string) []string {
	s := make([]string, 0)
	smap := make(map[string]string)
	for _, value := range m {
		//计算map长度
		length := len(smap)
		smap[value] = value
		//比较map长度, 如果map长度不相等， 说明key不存在
		if len(smap) != length {
			s = append(s, value)
		}
	}
	return s
}

// IntersectInt 交集  Int数组取相同的元素
func IntersectInt(slice1, slice2 []int) []int {
	m := make(map[int]int)
	nn := make([]int, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// IntersectStr 交集  Str数组取相同的元素
func IntersectStr(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

// GetDifferentIntArray 差集 Int数组 取出不同元素
func GetDifferentIntArray(sourceList, sourceList2 []int) (result []int) {
	for _, src := range sourceList {
		var find bool
		for _, target := range sourceList2 {
			if src == target {
				find = true
				continue
			}
		}
		if !find {
			result = append(result, src)
		}
	}
	return
}

// GetDifferentStrArray 差集 Str数组 取出不同元素
func GetDifferentStrArray(sourceList, sourceList2 []string) (result []string) {
	for _, src := range sourceList {
		var find bool
		for _, target := range sourceList2 {
			if src == target {
				find = true
				continue
			}
		}
		if !find {
			result = append(result, src)
		}
	}
	return
}

// ExistIntArray Int数组 存在某个数字
func ExistIntArray(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ExistStringArray Str字符串数组 存在某个字符串
func ExistStringArray(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// String_Prefix_list 开头列表检查
func String_Prefix_list(_list_str []string, _str string) bool {
	for _, v := range _list_str {
		// fmt.Println(v, _str, strings.HasPrefix(_str, v))
		if strings.HasPrefix(v, _str) {
			return true
		}
	}
	return false
}

// String_Suffix_list 结尾列表检查
func String_Suffix_list(_list_str []string, _str string) bool {
	for _, v := range _list_str {
		tmp_re := strings.HasSuffix(v, _str)
		if tmp_re {
			return true
		}
	}
	return false
}

// String_Contains_list 包含列表检查
func String_Contains_list(_list_str []string, _str string) bool {
	for _, v := range _list_str {
		tmp_re := strings.Contains(v, _str)
		if tmp_re {
			return true
		}
	}
	return false
}

// GetIPList ip 地址转 net.IP
// ips = "10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24"
func GetIPList(ips string) ([]net.IP, error) {
	addrS, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	List := addrS.Expand()
	return List, nil
}

// IpCidrCheck ip cidr 范围判断
// cidrIp = "10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24"
func IpCidrCheck(cidrIp, Ip string) bool {
	iplist, err := GetIPList(cidrIp)
	if err != nil {
		return false
	}
	strList := []string{}
	for _, v := range iplist {
		strList = append(strList, v.String())
	}
	// Printf_Color(strList)
	return ExistStringArray(strList, Ip)
}

func GetPorts(selection string) ([]int, error) {
	var ports []int
	if selection == "" {
		return ports, nil
	}

	ranges := strings.Split(selection, ",")
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("Invalid port selection segment: '%s'", r)
			}

			p1, err := strconv.Atoi(parts[0])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[0])
			}

			p2, err := strconv.Atoi(parts[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", parts[1])
			}

			if p1 > p2 {
				return nil, fmt.Errorf("Invalid port range: %d-%d", p1, p2)
			}

			for i := p1; i <= p2; i++ {
				ports = append(ports, i)
			}

		} else {
			if port, err := strconv.Atoi(r); err != nil {
				return nil, fmt.Errorf("Invalid port number: '%s'", r)
			} else {
				ports = append(ports, port)
			}
		}
	}
	return ports, nil
}

// CheckPort 检验输入字符串端口
func CheckPort(portStr string) error {
	portStr = strings.TrimSpace(portStr)
	if portStr == "" {
		return fmt.Errorf("The port cannot be empty")
	}
	ranges := strings.Split(portStr, ",")
	maxPort := 65535
	for _, r := range ranges {
		r = strings.TrimSpace(r)
		if strings.Contains(r, "-") {
			parts := strings.Split(r, "-")
			if len(parts) != 2 {
				return fmt.Errorf("Invalid port selection segment: '%s'", r)
			}
			port1, err := strconv.Atoi(parts[0])
			if err != nil {
				return fmt.Errorf("Invalid port number: '%s'", parts[0])
			}
			port2, err := strconv.Atoi(parts[1])
			if err != nil {
				return fmt.Errorf("Invalid port number: '%s'", parts[1])
			}
			if port1 > port2 {
				return fmt.Errorf("Invalid port range: %d-%d", port1, port2)
			}
			if port2 > maxPort {
				return fmt.Errorf("Invalid port range: %d-%d", maxPort, port2)
			}
			continue
		}
		port, err := strconv.Atoi(r)
		if err != nil {
			return fmt.Errorf("Invalid port number: '%s'", r)
		}
		if port > maxPort {
			return fmt.Errorf("Invalid port range: %d-%d", maxPort, port)
		}
	}
	return nil
}
