package engine

import (
	"fmt"
	"github.com/abusizhishen/do-once-while-concurrent/src"
	json "github.com/json-iterator/go"
	"log"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	tagDefaultValue = "-"
	tagTodayValue = "0100-01-01"
)

type DnfMap map[string]interface{}

// 数据类型
//big_float,  float ：浮点型
//num：整型
//string, time：字符串
type TagInfo struct {
	Id          string `json:"id"`
	Symbol      string `json:"symbol"`
	Value       Value  `json:"value"`
	TagType     int    `json:"tag_type"`
	Match       func(str string) bool
	TypeReflect string `json:"type_reflect"`
	Dnf         string `json:"dnf"`
	Relation    string `json:"relation"`
	RelationFun func(bool2 bool) bool
	Result      bool `json:"result"`
}

func (t *TagInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id          string `json:"id"`
		Symbol      string `json:"symbol"`
		Value       Value  `json:"value"`
		TagType     int    `json:"tag_type"`
		TypeReflect string `json:"type_reflect"`
		Dnf         string `json:"dnf"`
		Relation    string `json:"relation"`
		Result      bool   `json:"result"`
	}{
		t.Id,
		t.Symbol,
		t.Value,
		t.TagType,
		t.TypeReflect,
		t.Dnf,
		t.Relation,
		t.Result,
	})
}

type T struct {
	Type        string `json:"type"`
	Id          string `json:"id"`
	Symbol      string `json:"symbol"`
	Value       Value  `json:"value"`
	TagType     int    `json:"tag_type"`
	TypeReflect string `json:"type_reflect"`
	Dnf         string `json:"dnf"`
	Relation    string `json:"relation"`
}

func (t *TagInfo) UnmarshalJSON(data []byte) error {
	var tt T
	err := json.Unmarshal(data, &tt)
	if err != nil {
		return err
	}

	t.Id = tt.Id
	t.Symbol = tt.Symbol
	t.Value = tt.Value
	t.TagType = tt.TagType
	t.TypeReflect = tt.TypeReflect
	t.Dnf = tt.Dnf
	t.Relation = tt.Relation

	return nil
}


func (t *TagInfo) TimeMatch(value string) bool {
	switch value {
	case tagDefaultValue:
		switch t.Symbol {
		case "≥":
			value = time.Time{}.Format(layout)
		case "≤":
			value = time.Now().AddDate(10, 0, 0).Format(layout)
		}
	default:
	}

	tim, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}

	switch t.Symbol {
	case "≥":
		return t.Value.Time.Equal(tim)||t.Value.Time.After(tim)
	case "≤":
		return t.Value.Time.Before(tim)||t.Value.Time.Equal(tim)
	default:
		return false
	}
}

func (t *TagInfo) Float64Match(value string) bool {
	if value == tagDefaultValue {
		value = "0"
	}
	data, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(fmt.Sprintf("t:%+v,value:%s", t, value))
	}

	switch t.Symbol {
	case "<":
		return data < t.Value.Float64
	case "≤":
		return data <= t.Value.Float64
	case "=":
		return data == t.Value.Float64
	case "≠":
		return data != t.Value.Float64
	case "≥":
		return data >= t.Value.Float64
	case ">":
		return data > t.Value.Float64
	default:
		return false
	}
}

func (t *TagInfo) IntMatch(value string) bool {
	if value == tagDefaultValue {
		//tool.OutputInfo(fmt.Sprintf("tag:%+v,标签%s返回值异常 %v", t, t.TypeReflect, value))
		value = "0"
	}
	data, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	switch t.Symbol {
	case "<":
		return data < t.Value.Int
	case "≤":
		return data <= t.Value.Int
	case "=":
		return data == t.Value.Int
	case "≠":
		return data != t.Value.Int
	case "≥":
		return data >= t.Value.Int
	case ">":
		return data > t.Value.Int
	default:
		return false
	}
}

func (t *TagInfo) Int64Match(value string) bool {
	if value == tagDefaultValue {
		//tool.OutputInfo(fmt.Sprintf("tag:%+v,标签%s返回值异常 %v", t, t.TypeReflect, value))
		value = "0"
	}
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		panic(err)
	}

	switch t.Symbol {
	case "<":
		return data < t.Value.Int64
	case "≤":
		return data <= t.Value.Int64
	case "=":
		return data == t.Value.Int64
	case "≠":
		return data != t.Value.Int64
	case "≥":
		return data >= t.Value.Int64
	case ">":
		return data > t.Value.Int64
	default:
		return false
	}
}

func (t *TagInfo) StringMatch(value string) bool {
	switch t.Symbol {
	case "=":
		return t.Value.String == value
	case "≠":
		return t.Value.String != value
	default:
		return false
	}
}

func (t *TagInfo) StringArrMatch(value string) bool {
	data := strings.Split(strings.Trim(value, "[]"), ",")
	for k, v := range data {
		data[k] = strings.Trim(v, `"`)
	}

	switch t.Symbol {
	case "∈":
		for _, v := range data {
			if _, ok := t.Value.StringArr[v]; ok {
				return true
			}
		}

		return false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.StringArr[v]; ok {
				return false
			}
		}

		return true

	default:
		return false
	}
}

func (t *TagInfo) IntArrMatch(value string) bool {
	strData := strings.Split(strings.Trim(value, "[]"), ",")
	var data []int
	for _, str := range strData {
		v, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		data = append(data, v)
	}

	switch t.Symbol {
	case "∈":
		for _, v := range data {
			if _, ok := t.Value.IntArr[v]; ok {
				return true
			}
		}
		return false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.IntArr[v]; ok {
				return false
			}
		}

		return true

	default:
		return false
	}
}

func (t *TagInfo) Int64ArrMatch(value string) bool {
	strData := strings.Split(strings.Trim(value, "[]"), ",")
	var data []int64
	for _, str := range strData {
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			panic(err)
		}
		data = append(data, v)
	}

	switch t.Symbol {
	case "∈":
		for _, v := range data {
			if _, ok := t.Value.Int64Arr[v]; ok {
				return true
			}
		}
		return false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.Int64Arr[v]; ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (t *TagInfo) Float64ArrMatch(value string) bool {
	strData := strings.Split(strings.Trim(value, "[]"), ",")
	var data []float64
	for _, str := range strData {
		v, err := strconv.ParseFloat(str, 64)
		if err != nil {
			panic(err)
		}
		data = append(data, v)
	}

	switch t.Symbol {
	case "∈":
		for _, v := range data {
			if _, ok := t.Value.Float64Arr[v]; ok {
				return true
			}
		}
		return false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.Float64Arr[v]; ok {
				return false
			}
		}
		return true
	default:
		return false
	}
}

func (t *TagInfo) Init() {
	funs := []func(){
		t.Value.Init,
		t.SetValueType,
		t.Value.ParseValue,
		t.SetMatchFun,
		t.SetFunRelation,
	}

	for _, fun := range funs {
		fun()
	}
}

func (t *TagInfo) SetFunRelation() {
	switch t.Relation {
	case "", "∧":
		t.RelationFun = OutAnd
	case "∨":
		t.RelationFun = func(bool2 bool) bool {
			return bool2
		}
	case "┐":
		t.RelationFun = Not
	}
}

func (t *TagInfo) SetValueType() {
	t.TypeReflect = TypeArr[t.TagType]
	t.Value.Type = t.TypeReflect
}

func (t *Value) ParseValue() {
	switch t.Type {
	case "float64":
		t.Float64, _ = strconv.ParseFloat(t.Origin, 64)
	case "int":
		t.Int, _ = strconv.Atoi(t.Origin)
	case "int64":
		t.Int64, _ = strconv.ParseInt(t.Origin, 10, 64)
	case "string":
		t.String = t.Origin
	case "time":
		t.Int64, _ = strconv.ParseInt(t.Origin, 10, 64)
	case "[]float64":
		var sli = strings.Split(t.Origin, ",")
		for _, s := range sli {
			v, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(err)
			}
			t.Float64Arr[v] = ""
		}

	case "[]int":
		var sli = strings.Split(t.Origin, ",")
		for _, s := range sli {
			v, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			t.IntArr[v] = ""
		}
	case "[]string":
		var sli = strings.Split(t.Origin, ",")
		for _, s := range sli {
			t.StringArr[s] = ""
		}
	default:
		log.Println("unknown  type：")
		panic(t.Type)
	}
}

func (t *TagInfo) SetMatchFun() {
	switch t.TypeReflect {
	case "time":
		t.Match = t.TimeMatch
	case "int":
		t.Match = t.IntMatch
	case "int64":
		t.Match = t.Int64Match
	case "[]int64":
		t.Match = t.Int64ArrMatch
	case "float64":
		t.Match = t.Float64Match
	case "string":
		t.Match = t.StringMatch
	case "[]string":
		t.Match = t.StringArrMatch
	case "[]int":
		t.Match = t.IntArrMatch
	case "[]float64":
		t.Match = t.Float64ArrMatch

	default:
		panic("数据类型：" + t.TypeReflect + "错误")
	}
}

///**
//广告标签结构体数组
// */
//type AdTag []FirstLevel

/**
一级标签关系
*/


type FirstLevel struct {
	Cond        []*TagInfo `json:"cond"`
	In          string     `json:"in"`
	Out         string     `json:"out"`
	TypeReflect string     `json:"type_reflect"`
	ResultIn    bool       `json:"result_in"`
	ResultOut   bool       `json:"result_out"`
	Dnf         string     `json:"dnf"`

	FuncIn  func(data []bool) bool
	FuncOut func(data bool) bool

	TagsMap map[string]string `json:"tags_map"`
}

func (f *FirstLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Cond        []*TagInfo        `json:"cond"`
		In          string            `json:"in"`
		Out         string            `json:"out"`
		TypeReflect string            `json:"type_reflect"`
		ResultIn    bool              `json:"result_in"`
		ResultOut   bool              `json:"result_out"`
		Dnf         string            `json:"dnf"`
		TagsMap     map[string]string `json:"tags_map"`
	}{
		f.Cond,
		f.In,
		f.Out,
		f.TypeReflect,
		f.ResultIn,
		f.ResultOut,
		f.Dnf,
		f.TagsMap,
	})
}

func (f *FirstLevel) Unmarshal(data []byte) error {
	type F struct {
		Cond        []TagInfo `json:"cond"`
		In          string    `json:"in"`
		Out         string    `json:"out"`
		TypeReflect string    `json:"type_reflect"`
		ResultIn    bool      `json:"result_in"`
		ResultOut   bool      `json:"result_out"`
		Dnf         string    `json:"dnf"`
	}

	var ff F
	err := json.Unmarshal(data, &ff)
	if err != nil {
		return err
	}

	f = &FirstLevel{
		f.Cond,
		f.In,
		f.Out,
		f.TypeReflect,
		f.ResultIn,
		f.ResultOut,
		f.Dnf,
		nil,
		nil,
		f.TagsMap,
	}

	return nil
}

func (f *FirstLevel) Init() {
	var funs = []func(){
		f.SetFuncIn,
		f.SetFuncOut,
		f.CondInit,
	}

	for _, f := range funs {
		f()
	}
}

func (f *FirstLevel) SetFuncIn() {
	switch f.In {
	case "∧", "":
		f.FuncIn = And
	case "∨":
		f.FuncIn = Or
	default:
		//panic("逻辑关系错误 " + f.In)
		//默认与操作
		f.FuncIn = And
	}
}

func (f *FirstLevel) SetFuncOut() {
	switch f.Out {
	case "∧", "":
		f.FuncOut = OutAnd
	case "┐":
		f.FuncOut = Not
	default:
		//panic("逻辑关系错误 " + f.Out)
		f.FuncOut = OutAnd
	}
}

func (f *FirstLevel) CondInit() {
	f.TagsMap = make(map[string]string)
	for k, t := range f.Cond {
		t.Init()
		f.TagsMap[t.Id] = ""
		f.Cond[k] = t
	}
}

func (f *FirstLevel) Match(values map[string]string) (result bool) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("match excepted: %v",e)
			result = false
		}
	}()
	var resultArr = make([]bool, 0, len(f.Cond))
	for _, t := range f.Cond {
		t.Result = t.Match(values[t.Id])
		resultArr = append(resultArr, t.RelationFun(t.Result))
	}

	f.ResultIn = f.FuncIn(resultArr)
	f.ResultOut = f.FuncOut(f.ResultIn)
	return f.ResultOut
}

func And(data []bool) bool {
	var result = true
	for _, b := range data {
		result = result && b
		if !result {
			return false
		}
	}

	return true
}

func Or(data []bool) bool {
	for _, b := range data {
		if b {
			return true
		}
	}

	return false
}

func OutAnd(data bool) bool {
	return data
}
func Not(data bool) bool {
	return !data
}

type AdvertTag struct {
	Id       string            `json:"id"`
	Dnf      string            `json:"dnf"`
	AdTag    []FirstLevel      `json:"ad_tag"`
	TagsMap  map[string]string `json:"tags_map"`
	TagCount int               `json:"tag_count"`
	IsInit   bool              `json:"is_init"`
	sync.RWMutex
	DoOnceSameTime src.DoOnce `json:"-"`
	MatchResult    bool           `json:"match_result"`
}

func (a *AdvertTag) MarshalJSON() ([]byte, error) {
	return json.Marshal(
		struct {
			Id       string            `json:"id"`
			Dnf      string            `json:"dnf"`
			AdTag    []FirstLevel      `json:"ad_tag"`
			TagsMap  map[string]string `json:"tags_map"`
			TagCount int               `json:"tag_count"`
			IsInit   bool              `json:"is_init"`
		}{
			Id:       a.Id,
			Dnf:      a.Dnf,
			AdTag:    a.AdTag,
			TagsMap:  a.TagsMap,
			TagCount: a.TagCount,
			IsInit:   a.IsInit,
		})
}

func (a *AdvertTag) Unmarshal(data []byte) error {
	type A struct {
		Id       string            `json:"id"`
		Dnf      string            `json:"dnf"`
		AdTag    []FirstLevel      `json:"ad_tag"`
		TagsMap  map[string]string `json:"tags_map"`
		TagCount int               `json:"tag_count"`
		IsInit   bool              `json:"is_init"`
	}

	var aa A
	err := json.Unmarshal(data, &aa)
	if err != nil {
		return err
	}

	a.Id = aa.Id
	a.Dnf = aa.Dnf
	a.AdTag = aa.AdTag
	a.TagsMap = aa.TagsMap
	a.TagCount = aa.TagCount
	a.IsInit = aa.IsInit

	return nil
}

func (a *AdvertTag) Init() {
	if a.IsInit {
		return
	}

	if !a.DoOnceSameTime.Req(1) {
		a.DoOnceSameTime.Wait(1)
		return
	}
	defer a.DoOnceSameTime.Release(1)

	a.TagsMap = make(map[string]string)
	for k, d := range a.AdTag {
		d.Init()
		for k, v := range d.TagsMap {
			a.TagsMap[k] = v
		}
		a.AdTag[k] = d
	}

	a.TagCount = len(a.TagsMap)
	a.IsInit = true
}

/*
广告标签和用户标签匹配方法,匹配则返回true
*/
func (a *AdvertTag) Match(tagValues map[string]string) bool {
	for _, t := range a.AdTag {
		t.ResultOut = t.Match(tagValues)
		if !t.ResultOut {
			return t.ResultOut
		}
	}

	a.MatchResult = true
	return a.MatchResult
}

func New(dnf string)(a *AdvertTag,err error) {
	a = new(AdvertTag)
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("dnf解析报错：%+v,无效的dnf:%s", e, dnf)
		}
	}()
	defer a.Init()
	a.Dnf = dnf
	var express = `(∧|∨|┐){0,1}(\[.*?\])`
	reg := regexp.MustCompile(express)
	s := reg.FindAllStringSubmatch(a.Dnf, -1)
	if len(s) == 0 {
		panic("dnf解析为空")
	}

	a.AdTag = make([]FirstLevel, len(s))
	for idx, dnf := range s {
		if len(dnf) != 3 {
			//log.Printf("dnf异常:%v+", dnf)
			continue
		}
		if dnf[1] == "" {
			dnf[1] = "∧"
		}

		var firstLevel FirstLevel
		firstLevel.In = "∧"
		firstLevel.Out = dnf[1]
		firstLevel.FromDnf(dnf[2])
		if len(firstLevel.Cond) > 0 {
			firstLevel.In = firstLevel.Cond[len(firstLevel.Cond)-1].Relation
		}
		a.AdTag[idx] = firstLevel
	}
	return
}

//[(exp_first_attend_status∈{0.000000,1.000000}|4)∧(exp_attend_num={1}|2)]
func (a *FirstLevel) FromDnf(dnf string) {
	a.Dnf = dnf
	var express = `(∧|∨|┐){0,1}\(.*?\)`
	reg := regexp.MustCompile(express)
	s := reg.FindAllString(a.Dnf, -1)
	if len(s) == 0 {
		panic("dnf解析为空")
	}

	a.Cond = make([]*TagInfo, len(s))
	for idx, dnf := range s {
		var t TagInfo
		t.FromDnf(dnf)
		a.Cond[idx] = &t
	}
}

// (exp_first_attend_status∈{0.000000,1.000000}|4)
func (t *TagInfo) FromDnf(dnf string) {
	t.Dnf = dnf
	var express = `(∧|∨|┐){0,1}\((.*)+(∈|∉|=|≠|<|≤|≥|>){1}\{(.*?)}{1}\|(\d+)\)$`
	reg := regexp.MustCompile(express)
	d := reg.FindAllStringSubmatch(t.Dnf, -1)

	if len(d) != 1 {
		panic("taginfo正则匹配异常")
	}

	s := d[0]
	if len(s) != 6 {
		panic("taginfo正则匹配异常")
	}

	t.Relation = s[1]
	t.Id = s[2]
	t.Symbol = s[3]
	t.Value.Origin = s[4]
	i, err := strconv.Atoi(s[5])
	if err != nil {
		panic(err)
	}
	t.TagType = i
}