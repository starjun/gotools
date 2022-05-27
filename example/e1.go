package main

import (
	"fmt"

	"github.com/starjun/gotools"
)

func main() {
	a := []int{1, 2, 1, 3, 4}
	b := []int{4, 3, 5, 6}

	c := gotools.GetDifferentIntArray(a, b)

	fmt.Println(c)

	// fmt.Println(gotools.IpCidrCheck("192.168.1.1-10", "192.168.1.5"))

	//
	fmt.Println("MatchStr cidr: ", gotools.MatchStr("192.168.1.0/24", []string{"192.168.2.5"}, "cidr", false))

	//
	fmt.Println("MatchStr prefix: ", gotools.MatchStr("http://www.aop.com}", []string{"https://"}, "prefix", false))

	//
	fmt.Println("MatchStr in: ", gotools.MatchStr("23dsafasdf", []string{"==="}, "in", false))

	// =
	fmt.Println("MatchStr =: ", gotools.MatchStr("hasdf22om", []string{"*"}, "=", false))

	//
	fmt.Println("MatchStr suffix: ", gotools.MatchStr("hasdf22.com", []string{".cn"}, "suffix", false))

	//
	fmt.Println("MatchStr reg: ", gotools.MatchStr("hasdfffdf.com", []string{".*[0-9]+.*"}, "reg", false))

	maptmp1 := make(map[string]string)

	maptmp1["method"] = "GET"
	maptmp1["uri"] = "/admin/v1/get/list"
	maptmp1["args_name"] = "admin"
	maptmp1["passwd"] = "asdfasdf324dsfdsaf"

	crule1 := gotools.CRule{
		Des: "匹配用户名",
		Opt: "=",
		ReStrList: []string{
			"user",
			"admin",
			"test",
			"ftp",
		},
		Rev:        false,
		Lcon:       "and",
		MaLocation: "args_name",
	}

	fmt.Println("maptmp1 - crule1 [true] == >", gotools.MapCRuleMatch(maptmp1, crule1))

	crule2 := gotools.CRule{
		Des: "匹配uri",
		Opt: "prefix",
		ReStrList: []string{
			"/admin/v1/",
			"/admin/v2/",
			"/admin/manage/",
		},
		Rev:        false,
		Lcon:       "and",
		MaLocation: "uri",
	}

	fmt.Println("maptmp1 - crule2 [true] ==>", gotools.MapCRuleMatch(maptmp1, crule2))

	crule3 := gotools.CRule{
		Des: "匹配密码",
		Opt: "in",
		ReStrList: []string{
			"123456789",
			"adminadmin",
			"admin888",
		},
		Rev:        false,
		Lcon:       "or",
		MaLocation: "passwd",
	}

	fmt.Println("maptmp1 - crule3 [false] ==>", gotools.MapCRuleMatch(maptmp1, crule3))

	// test uri prefix and (args_name = or passwd in) = true
	//        true and (true or false)

	listRules := []gotools.CRule{}
	listRules = append(listRules, crule2)
	listRules = append(listRules, crule3)
	listRules = append(listRules, crule1)

	fmt.Println("maptmp1 - crule2 and crule3 or crule1 ==>", gotools.MapCrulesListMatch(maptmp1, listRules))

}
