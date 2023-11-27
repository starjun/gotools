package gotools

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/pretty"
)

const _VERSION = "v0.0.4"

var (
// 一些变量 ...
)

func init() {
	rand.NewSource(time.Now().UnixNano()) // 产生随机种子
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
		if ExistStringArray(excludeFields, fieldName) {
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

func GetDigit(s string) string {
	reg, _ := regexp.Compile(`[0-9.]+`)
	ret := string(reg.Find([]byte(s)))
	if strings.TrimSpace(ret) == "" {
		ret = "0"
	}
	return ret
}
