// profile
package test

import (
	"hope/http/request"
	"hope/model"
	"regexp"
	"strconv"
)

//将解析表达式都编译好，避免每次都编译
var ageRe = regexp.MustCompile(`"basicInfo":["[\s\S]*?","([\s\S]*?)岁",`)
var heightRe = regexp.MustCompile(`"basicInfo":["[\s\S]*?","[\s\S]*?","[\s\S]*?","([\s\S]*?)cm",`)
var weightRe = regexp.MustCompile(`"basicInfo":["[\s\S]*?","[\s\S]*?","[\s\S]*?","[\s\S]*?","([[\s\S]*?)kg",`)
var nameRe = regexp.MustCompile(`"nickname":"([\s\S]*?)",`)
var genderRe = regexp.MustCompile(`"genderString":"([\s\S]*?)",`)
var incomeRe = regexp.MustCompile(`"basicInfo":[\s\S]*?"月收入:([\s\S]*?)",`)
var marriageRe = regexp.MustCompile(`"marriageString":"([\s\S]*?)",`)
var educationRe = regexp.MustCompile(`"educationString":"([\s\S]*?)",`)
var occupationRe = regexp.MustCompile(`"basicInfo":[\s\S]*?月收入:[\s\S]*?","([\s\S]*?)","[\s\S]*?"]`)
var hukouRe = regexp.MustCompile(`"basicInfo":[\s\S]*?"工作地:([\s\S]*?)",`)

//用户信息解析
func ParseProfile(contents []byte) *request.ParseResult {
	profile := &model.Profile{}
	age, err := strconv.Atoi(extractString(contents, ageRe))
	if err == nil {
		profile.Age = age
	}
	height, err := strconv.Atoi(extractString(contents, heightRe))
	if err == nil {
		profile.Height = height
	}
	weight, err := strconv.Atoi(extractString(contents, weightRe))
	if err == nil {
		profile.Weight = weight
	}
	profile.Name = extractString(contents, nameRe)
	profile.Gender = extractString(contents, genderRe)
	profile.Income = extractString(contents, incomeRe)
	profile.Marriage = extractString(contents, marriageRe)
	profile.Education = extractString(contents, educationRe)
	profile.Occupation = extractString(contents, occupationRe)
	profile.Hukou = extractString(contents, hukouRe)
	result := &request.ParseResult{
		Items: []interface{}{profile},
	}
	return result
}

//将正则解析抽象为一个函数
func extractString(contents []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(contents)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}
