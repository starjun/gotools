package gotools

import (
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
