package gotools

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
	"github.com/malfunkt/iprange"
	"github.com/tidwall/pretty"
)

const GOTOOLS_VERSION = "v0.0.2"

var (
// 一些变量 ...
)

//并集  合并两个数组并去重
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

//并集  合并两个字符串数组并去重
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

// Int数组 去重
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

// 字符串数组 去重
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

//交集  Int数组取相同的元素
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

//交集  Str数组取相同的元素
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

//差集 Int数组 取出不同元素
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

//差集 Str数组 取出不同元素
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

// Int数组 存在某个数字
func ExistIntArray(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Str字符串数组 存在某个字符串
func ExistStringArray(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// 是否是合法 ipv4 地址 ipv6待增加
func IsIp(ip string) (b bool) {
	if m, _ := regexp.MatchString("^(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)\\.(25[0-5]|2[0-4]\\d|[0-1]\\d{2}|[1-9]?\\d)", ip); !m {
		return false
	}
	return true
}

//是否是ip网段格
func IsIpSegment(ip string) (b bool) {
	if m, _ := regexp.MatchString("^((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}/(3[0-2]|[1-2]\\d|\\d)$", ip); !m {
		return false
	}
	return true
}

// 判断是否是一个合法域名
func DomainCheck(domain string) bool {
	var match bool
	NotLine := "^([a-zA-Z0-9]([a-zA-Z0-9\\-]{0,61}[a-zA-Z0-9])?\\.)+[a-zA-Z]{2,6}"
	//支持以http://或者https://开头并且域名中间没有/的情况
	// NotLine := "(http(s)?:\\/\\/)?(www\\.)?[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(\\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})+(:\\d+)*(\\/\\w+\\.\\w+)*$"
	match, _ = regexp.MatchString(NotLine, domain)
	return match
}

// 开头列表检查
func String_Prefix_list(_list_str []string, _str string) bool {
	for _, v := range _list_str {
		// fmt.Println(v, _str, strings.HasPrefix(_str, v))
		if strings.HasPrefix(v, _str) {
			return true
		}
	}
	return false
}

// 结尾列表检查
func String_Suffix_list(_list_str []string, _str string) bool {
	for _, v := range _list_str {
		tmp_re := strings.HasSuffix(v, _str)
		if tmp_re {
			return true
		}
	}
	return false
}

// 包含列表检查
func String_Contains_list(_list_str []string, _str string) bool {
	for _, v := range _list_str {
		tmp_re := strings.Contains(v, _str)
		if tmp_re {
			return true
		}
	}
	return false
}

// interface2sting
func Strval(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// MD5 计算
func MD5V(str []byte) string {
	h := md5.New()
	h.Write(str)
	return hex.EncodeToString(h.Sum(nil))
}

// ReadAll 读取所有文件内容
func ReadAll(filePth string) ([]byte, error) {
	// return ioutil.ReadFile(filename)
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ZipFiles
//@description: 压缩文件
//@param: filename string, files []string, oldform, newform string
//@return: error
func ZipFiles(filename string, files []string, oldForm, newForm string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		_ = newZipFile.Close()
	}()

	zipWriter := zip.NewWriter(newZipFile)
	defer func() {
		_ = zipWriter.Close()
	}()

	// 把files添加到zip中
	for _, file := range files {

		err = func(file string) error {
			zipFile, err := os.Open(file)
			if err != nil {
				return err
			}
			defer zipFile.Close()
			// 获取file的基础信息
			info, err := zipFile.Stat()
			if err != nil {
				return err
			}

			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// 使用上面的FileInforHeader() 就可以把文件保存的路径替换成我们自己想要的了，如下面
			header.Name = strings.Replace(file, oldForm, newForm, -1)

			// 优化压缩
			// 更多参考see http://golang.org/pkg/archive/zip/#pkg-constants
			header.Method = zip.Deflate

			writer, err := zipWriter.CreateHeader(header)
			if err != nil {
				return err
			}
			if _, err = io.Copy(writer, zipFile); err != nil {
				return err
			}
			return nil
		}(file)
		if err != nil {
			return err
		}
	}
	return nil
}

//@author: [songzhibin97](https://github.com/songzhibin97)
//@function: FileMove
//@description: 文件移动供外部调用
//@param: src string, dst string(src: 源位置,绝对路径or相对路径, dst: 目标位置,绝对路径or相对路径,必须为文件夹)
//@return: err error
func FileMove(src string, dst string) (err error) {
	if dst == "" {
		return nil
	}
	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	dst, err = filepath.Abs(dst)
	if err != nil {
		return err
	}
	var revoke = false
	dir := filepath.Dir(dst)
Redirect:
	_, err = os.Stat(dir)
	if err != nil {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
		if !revoke {
			revoke = true
			goto Redirect
		}
	}
	return os.Rename(src, dst)
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: PathExists
//@description: 文件目录是否存在
//@param: path string
//@return: bool, error
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// VersionOrdinal 版本号转化
func VersionOrdinal(version string) string {
	// ISO/IEC 14651:2011
	const maxByte = 1<<8 - 1
	vo := make([]byte, 0, len(version)+8)
	j := -1
	for i := 0; i < len(version); i++ {
		b := version[i]
		if '0' > b || b > '9' {
			vo = append(vo, b)
			j = -1
			continue
		}
		if j == -1 {
			vo = append(vo, 0x00)
			j = len(vo) - 1
		}
		if vo[j] == 1 && vo[j+1] == '0' {
			vo[j+1] = b
			continue
		}
		if vo[j]+1 > maxByte {
			panic("VersionOrdinal: invalid version")
		}
		vo = append(vo, b)
		vo[j]++
	}
	return string(vo)
}

// 彩色打印传入 对象
func Printf_Color(value interface{}) {
	jjjj, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("%s\n", pretty.Pretty(jjjj))
	fmt.Printf("%s\n", pretty.Color(pretty.Pretty(jjjj), pretty.TerminalStyle))
}

// ip 地址转 net.IP
// ips = "10.0.0.1, 10.0.0.5-10, 192.168.1.*, 192.168.10.0/24"
func GetIPList(ips string) ([]net.IP, error) {
	addrS, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	List := addrS.Expand()
	return List, nil
}

// ip cidr 范围判断
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

// -- 规则判断
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
	} else if _opt == "prefix" {
		for i := 0; i < len(e_str); i++ {
			if strings.HasPrefix(_sourceStr, e_str[i]) {
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
	} else if _opt == "cidr" {
		for i := 0; i < len(e_str); i++ {
			if IpCidrCheck(_sourceStr, e_str[i]) {
				return !re
			}
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
