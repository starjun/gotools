package gotools

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	mapset "github.com/deckarep/golang-set"
	"github.com/malfunkt/iprange"
	"github.com/tidwall/pretty"
)

const _VERSION = "v0.0.3"

var (
// 一些变量 ...
)

func init() {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
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

// ZipFiles
// @author: [piexlmax](https://github.com/piexlmax)
// @function: ZipFiles
// @description: 压缩文件
// @param: filename string, files []string, oldform, newform string
// @return: error
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

// FileMove
// @author: [songzhibin97](https://github.com/songzhibin97)
// @function: FileMove
// @description: 文件移动供外部调用
// @param: src string, dst string(src: 源位置,绝对路径or相对路径, dst: 目标位置,绝对路径or相对路径,必须为文件夹)
// @return: err error
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

// PathExists
// @author: [piexlmax](https://github.com/piexlmax)
// @function: PathExists
// @description: 文件目录是否存在
// @param: path string
// @return: bool, error
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

// Printf_Color 彩色打印传入 对象
func Printf_Color(value interface{}) {
	jjjj, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Printf("%s\n", pretty.Pretty(jjjj))
	fmt.Printf("%s\n", pretty.Color(pretty.Pretty(jjjj), pretty.TerminalStyle))
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

// Md5Sum md5哈希函数
// s  字符串
// return： MD5哈希字符串
func Md5Sum(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetUUID 获取uuid
// 返回uuid字符串（36位）
func GetUUID() string {
	return uuid.NewV4().String()
}

// StrInSlice 判断字符串是否在列表中
// return：bool值，true - 存在 false - 不存在
func StrInSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

// IntInSlice 判断整型是否在列表中
// return：bool值，true - 存在 false - 不存在
func IntInSlice(i int, slice []int) bool {
	for _, v := range slice {
		if v == i {
			return true
		}
	}
	return false
}

// Strval 任意类型转字符串
// 入参：interface{}
// return：字符串
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

// GetIpList str ips [192.168.1.1/24,10.1.1.1-200] 转 net.IP
func GetIpList(ips string) ([]net.IP, error) {
	addressList, err := iprange.ParseList(ips)
	if err != nil {
		return nil, err
	}
	list := addressList.Expand()

	return list, err
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

// StructDiffValue 比较两个结构体
func StructDiffValue(old, new interface{}, excludeFields ...string) map[string]interface{} {
	if reflect.TypeOf(old).Kind() != reflect.Struct || reflect.TypeOf(new).Kind() != reflect.Struct {
		log.Println(reflect.TypeOf(old).Kind(), reflect.TypeOf(new).Kind())
		return nil
	}

	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)
	oldType := reflect.TypeOf(old)

	result := make(map[string]interface{})

	for i := 0; i < oldVal.NumField(); i++ {

		fieldName := oldType.Field(i).Name
		if StrInSlice(fieldName, excludeFields) {
			continue
		}

		switch oldVal.Field(i).Kind() {
		case reflect.String:
			if oldVal.Field(i).String() != newVal.Field(i).String() {
				result[fieldName] = map[string]interface{}{
					"old": oldVal.Field(i).String(),
					"new": newVal.Field(i).String(),
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if oldVal.Field(i).Int() != newVal.Field(i).Int() {
				result[fieldName] = map[string]interface{}{
					"old": oldVal.Field(i).Int(),
					"new": newVal.Field(i).Int(),
				}
			}
		case reflect.Float64, reflect.Float32:
			if oldVal.Field(i).Float() != newVal.Field(i).Float() {
				result[fieldName] = map[string]interface{}{
					"old": oldVal.Field(i).Float(),
					"new": newVal.Field(i).Float(),
				}
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if oldVal.Field(i).Uint() != newVal.Field(i).Uint() {
				result[fieldName] = map[string]interface{}{
					"old": oldVal.Field(i).Uint(),
					"new": newVal.Field(i).Uint(),
				}
			}
		case reflect.Bool:
			if oldVal.Field(i).Bool() != newVal.Field(i).Bool() {
				result[fieldName] = map[string]interface{}{
					"old": oldVal.Field(i).Bool(),
					"new": newVal.Field(i).Bool(),
				}
			}
		default:
			//TODO
		}
	}

	return result
}

// GetFileLine 获取文件名和行号
func GetFileLine(v ...int) (string, int) {
	skip := 1
	if len(v) == 1 {
		skip = v[0]
	}
	_, file, line, _ := runtime.Caller(skip)

	return file, line
}

// CreateDateDir 创建目录
func CreateDateDir(Path, dateFolderName string) string {
	folderPath := filepath.Join(Path, dateFolderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.MkdirAll(folderPath, os.ModePerm)
	}
	return folderPath
}

// FloatFormat 格式化浮点数
// f  浮点数
// v  保留小数位数
func FloatFormat(f float64, v ...int) float64 {
	n := 2
	if len(v) > 0 {
		n = v[0]
	}
	fmtStr := strings.ReplaceAll(`%.{{num}}f`, "{{num}}", strconv.Itoa(n))
	fmt.Println("fmtStr = ", fmtStr)
	r, err := strconv.ParseFloat(fmt.Sprintf(fmtStr, f), 64)
	if err != nil {
		log.Println(err)
	}

	return r
}

func Random(min, max int) int {
	return rand.Intn(max-min) + min
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

// PrintColor 优雅打印对象
func PrintColor(value interface{}) {
	j, err := json.Marshal(value)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", pretty.Color(pretty.Pretty(j), pretty.TerminalStyle))
}

// StringPrefixList 开头列表检查
func StringPrefixList(str string, listStr []string) (b bool) {
	b = false
	for _, v := range listStr {
		tmpRe := strings.HasPrefix(str, v)
		if tmpRe {
			b = true
			return
		}
	}
	return
}

// StringSuffixList 结尾列表检查
func StringSuffixList(str string, listStr []string) (b bool) {
	b = false
	for _, v := range listStr {
		tmpRe := strings.HasSuffix(str, v)
		if tmpRe {
			b = true

			return
		}
	}

	return
}

func GetDigit(s string) string {
	reg, _ := regexp.Compile(`[0-9.]+`)
	ret := string(reg.Find([]byte(s)))
	if strings.TrimSpace(ret) == "" {
		ret = "0"
	}
	return ret
}
