package engine

import (
	"errors"
	"fmt"
	"github.com/abusizhishen/do-once-while-concurrent/src"
	json "github.com/json-iterator/go"
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
	Match       func(str string)  (result bool,err error)
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


func (t *TagInfo) TimeMatch(value string)  (result bool,err error)  {
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
		return t.Value.Time.Equal(tim)||t.Value.Time.After(tim),nil
	case "≤":
		return t.Value.Time.Before(tim)||t.Value.Time.Equal(tim),nil
	default:
		return false,nil
	}
}

func (t *TagInfo) Float64Match(value string)  (result bool,err error)  {
	if value == tagDefaultValue {
		value = "0"
	}
	data, err := strconv.ParseFloat(value, 64)
	if err != nil {
		panic(fmt.Sprintf("t:%+v,value:%s", t, value))
	}

	switch t.Symbol {
	case "<":
		result =  data < t.Value.Float64
	case "≤":
		result =  data <= t.Value.Float64
	case "=":
		result =  data == t.Value.Float64
	case "≠":
		result =  data != t.Value.Float64
	case "≥":
		result =  data >= t.Value.Float64
	case ">":
		result =  data > t.Value.Float64
	default:
		result =  false
	}
	
	return result,err
}

func (t *TagInfo) IntMatch(value string)  (result bool,err error)  {
	if value == tagDefaultValue {
		//tool.OutputInfo(fmt.Sprintf("tag:%+v,标签%s返回值异常 %v", t, t.TypeReflect, value))
		value = "0"
	}
	data, err := strconv.Atoi(value)
	if err != nil {
		return false,err
	}
	switch t.Symbol {
	case "<":
		result = data < t.Value.Int
	case "≤":
		result = data <= t.Value.Int
	case "=":
		result =  data == t.Value.Int	
	case "≠":
		result =  data != t.Value.Int
	case "≥":
		result =  data >= t.Value.Int
	case ">":
		result =  data > t.Value.Int
	default:
		result =  false
	}
	
	return result,err
}

func (t *TagInfo) Int64Match(value string)  (result bool,err error)  {
	if value == tagDefaultValue {
		//tool.OutputInfo(fmt.Sprintf("tag:%+v,标签%s返回值异常 %v", t, t.TypeReflect, value))
		value = "0"
	}
	data, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return false,err
	}

	switch t.Symbol {
	case "<":
		result = data < t.Value.Int64
	case "≤":
		result = data <= t.Value.Int64
	case "=":
		result = data == t.Value.Int64
	case "≠":
		result = data != t.Value.Int64
	case "≥":
		result = data >= t.Value.Int64
	case ">":
		result = data > t.Value.Int64
	default:
		result = false
	}
	
	return result,err
}

func (t *TagInfo) StringMatch(value string)  (result bool,err error)  {
	switch t.Symbol {
	case "=":
		result =  t.Value.String == value
	case "≠":
		result =  t.Value.String != value
	default:
		result =  false
	}

	return result,err
}

func (t *TagInfo) StringArrMatch(value string)  (result bool,err error)  {
	data := strings.Split(strings.Trim(value, "[]"), ",")
	for k, v := range data {
		data[k] = strings.Trim(v, `"`)
	}

	switch t.Symbol {
	case "∈":
		for _, v := range data {
			if _, ok := t.Value.StringArr[v]; ok {
				result =  true
			}
		}

		result =  false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.StringArr[v]; ok {
				result =  false
			}
		}

		result =  true

	default:
		result =  false
	}

	return result,err
}

func (t *TagInfo) IntArrMatch(value string)  (result bool,err error)  {
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
				result =  true
			}
		}
		result =  false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.IntArr[v]; ok {
				result =  false
			}
		}

		result =  true

	default:
		result =  false
	}
	
	return result,err
}

func (t *TagInfo) Int64ArrMatch(value string)  (result bool,err error)  {
	strData := strings.Split(strings.Trim(value, "[]"), ",")
	var data []int64
	for _, str := range strData {
		v, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return false,err
		}
		data = append(data, v)
	}

	switch t.Symbol {
	case "∈":
		for _, v := range data {
			if _, ok := t.Value.Int64Arr[v]; ok {
				result = true
			}
		}
		result = false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.Int64Arr[v]; ok {
				result = false
			}
		}
		result = true
	default:
		result = false
	}
	
	return result,err
}

func (t *TagInfo) Float64ArrMatch(value string)  (result bool,err error)  {
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
				result =  true
			}
		}
		result =  false
	case "∉":
		for _, v := range data {
			if _, ok := t.Value.Float64Arr[v]; ok {
				result =  false
			}
		}
		result =  true
	default:
		result =  false
	}
	
	return result,err
}

func (t *TagInfo) Init() (err error){
	t.Value.Init()
	t.SetValueType()

	funs := []func()(err error){
		t.Value.ParseValue,
		t.SetMatchFun,
		t.SetFunRelation,
	}

	for _, fun := range funs {
		err = fun()
		if err != nil{
			return
		}
	}

	return
}

func (t *TagInfo) SetFunRelation() (err error){
	switch t.Relation {
	case "", "∧":
		t.RelationFun = OutAnd
	case "∨":
		t.RelationFun = func(bool2 bool) bool {
			return bool2
		}
	case "┐":
		t.RelationFun = Not
	default:
		err = fmt.Errorf("unkonwn logic relation: %s", t.Relation)
	}

	return err
}

func (t *TagInfo) SetValueType() {
	t.TypeReflect = TypeArr[t.TagType]
	t.Value.Type = t.TypeReflect
}

func (v *Value) ParseValue() (err error){
	switch v.Type {
	case "float64":
		v.Float64, _ = strconv.ParseFloat(v.Origin, 64)
	case "int":
		v.Int, _ = strconv.Atoi(v.Origin)
	case "int64":
		v.Int64, _ = strconv.ParseInt(v.Origin, 10, 64)
	case "string":
		v.String = v.Origin
	case "time":
		v.Int64, _ = strconv.ParseInt(v.Origin, 10, 64)
	case "[]float64":
		var sli = strings.Split(v.Origin, ",")
		for _, s := range sli {
			vv, err := strconv.ParseFloat(s, 64)
			if err != nil {
				panic(err)
			}
			v.Float64Arr[vv] = ""
		}

	case "[]int":
		var sli = strings.Split(v.Origin, ",")
		for _, s := range sli {
			vv, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			v.IntArr[vv] = ""
		}
	case "[]string":
		var sli = strings.Split(v.Origin, ",")
		for _, s := range sli {
			v.StringArr[s] = ""
		}
	default:
		err = fmt.Errorf("unknown type: %s", v.Type)
	}

	return err
}

func (t *TagInfo) SetMatchFun()(err error) {
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
		err = errors.New("数据类型：" + t.TypeReflect + "错误")
	}

	return err
}

///**
//广告标签结构体数组
// */
//type Node []FirstLevel

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

func (f *FirstLevel) Init() (err error){
	var funs = []func()error{
		f.SetFuncIn,
		f.SetFuncOut,
		f.CondInit,
	}

	for _, f := range funs {
		err = f()
		if err != nil{
			return
		}
	}

	return
}

func (f *FirstLevel) SetFuncIn() (err error){
	switch f.In {
	case "∧", "":
		f.FuncIn = And
	case "∨":
		f.FuncIn = Or
	default:
		//panic("逻辑关系错误 " + f.In)
		//默认与操作
		//f.FuncIn = And
		err = fmt.Errorf("逻辑关系错误:%s", f.In)
	}

	return
}

func (f *FirstLevel) SetFuncOut()(err error) {
	switch f.Out {
	case "∧", "":
		f.FuncOut = OutAnd
	case "┐":
		f.FuncOut = Not
	default:
		//panic("逻辑关系错误 " + f.Out)
//		f.FuncOut = OutAnd
		err = fmt.Errorf("逻辑关系错误:%s", f.Out)
	}

	return
}

func (f *FirstLevel) CondInit() (err error){
	f.TagsMap = make(map[string]string)
	for k, t := range f.Cond {
		err = t.Init()
		if err != nil{
			return err
		}
		f.TagsMap[t.Id] = ""
		f.Cond[k] = t
	}

	return
}

func (f *FirstLevel) Match(values map[string]string)  (result bool,err error)  {
	var resultArr = make([]bool, 0, len(f.Cond))
	for _, t := range f.Cond {
		t.Result,err = t.Match(values[t.Id])
		if err != nil{
			return false, err
		}
		resultArr = append(resultArr, t.RelationFun(t.Result))
	}

	f.ResultIn = f.FuncIn(resultArr)
	f.ResultOut = f.FuncOut(f.ResultIn)
	return f.ResultOut,nil
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
	Node     []node      	   `json:"ad_tag"`
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
			Node    []node      `json:"ad_tag"`
			TagsMap  map[string]string `json:"tags_map"`
			TagCount int               `json:"tag_count"`
			IsInit   bool              `json:"is_init"`
		}{
			Id:       a.Id,
			Dnf:      a.Dnf,
			Node:     a.Node,
			TagsMap:  a.TagsMap,
			TagCount: a.TagCount,
			IsInit:   a.IsInit,
		})
}

func (a *AdvertTag) Unmarshal(data []byte) error {
	type A struct {
		Id       string            `json:"id"`
		Dnf      string            `json:"dnf"`
		Node    []node      `json:"ad_tag"`
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
	a.Node = aa.Node
	a.TagsMap = aa.TagsMap
	a.TagCount = aa.TagCount
	a.IsInit = aa.IsInit

	return nil
}

func (a *AdvertTag) Init()(err error) {
	if a.IsInit {
		return
	}

	if !a.DoOnceSameTime.Req(1) {
		a.DoOnceSameTime.Wait(1)
		return
	}
	defer a.DoOnceSameTime.Release(1)

	a.TagsMap = make(map[string]string)
	for k, d := range a.Node {
		err = d.Init()
		for k, v := range d.TagsMap {
			a.TagsMap[k] = v
		}
		a.Node[k] = d
	}

	a.TagCount = len(a.TagsMap)
	a.IsInit = true

	return
}

/*
广告标签和用户标签匹配方法,匹配则返回true
*/
func (a *AdvertTag) Match(tagValues map[string]string) (result bool,err error) {
	for _, t := range a.Node {
		t.ResultOut,err = t.Match(tagValues)
		if err != nil{
			return false, err
		}
		if !t.ResultOut {
			return t.ResultOut,nil
		}
	}

	a.MatchResult = true
	return a.MatchResult,nil
}

func New(dnf string)(a engine,err error) {
	a = new(AdvertTag)
	err = a.Parse(dnf)
	return 
}

func (a *AdvertTag)Parse(dnf string)(err error) {
	a.Dnf = dnf
	var express = `(∧|∨|┐){0,1}(\[.*?\])`
	reg := regexp.MustCompile(express)
	s := reg.FindAllStringSubmatch(a.Dnf, -1)
	if len(s) == 0 {
		err = fmt.Errorf("dnf解析为空:%s", dnf)
		return
	}

	var tags []FirstLevel
	for _, dnf := range s {
		if len(dnf) != 3 {
			err = fmt.Errorf("invalid sub dnf :%v+", dnf)
			return
		}
		if dnf[1] == "" {
			dnf[1] = "∧"
		}

		var firstLevel FirstLevel
		firstLevel.In = "∧"
		firstLevel.Out = dnf[1]
		if dnf[2] == "[]" {
			continue
		}
		err = firstLevel.Parse(dnf[2])
		if err != nil{
			return  err
		}
		if len(firstLevel.Cond) > 0 {
			firstLevel.In = firstLevel.Cond[len(firstLevel.Cond)-1].Relation
		}
		if len(firstLevel.Cond) == 0 {
			continue
		}
		tags = append(tags, firstLevel)
	}

	a.Node = tags
	err = a.Init()
	return
}

//[(exp_first_attend_status∈{0.000000,1.000000}|4)∧(exp_attend_num={1}|2)]
func (f *FirstLevel) Parse(dnf string) (err error){
	f.Dnf = dnf
	var express = `(∧|∨|┐){0,1}\(.*?\)`
	reg := regexp.MustCompile(express)
	s := reg.FindAllString(f.Dnf, -1)
	if len(s) == 0 {
		return errors.New("dnf解析为空")
	}

	var nodes []*TagInfo
	for _, dnf := range s {
		if dnf == "()" {
			continue
		}
		var t TagInfo
		err = t.Parse(dnf)
		if err != nil{
			return err
		}

		nodes = append(nodes, &t)
	}

	f.Cond = nodes
	return
}

// (exp_first_attend_status∈{0.000000,1.000000}|4)
func (t *TagInfo) Parse(dnf string)(err error) {
	t.Dnf = dnf
	var express = `^(∧|∨|┐){0,1}\(([\w_\-\.]+)(∈|∉|=|≠|<|≤|≥|>){1}\{(.*?)}{1}\|(\d+)\)$`
	reg := regexp.MustCompile(express)
	d := reg.FindAllStringSubmatch(t.Dnf, -1)

	if len(d) != 1 {
		return fmt.Errorf("taginfo正则匹配异常:%s",dnf)
	}

	s := d[0]
	if len(s) != 6 {
		return fmt.Errorf("taginfo正则匹配异常:%s",dnf)
	}

	t.Relation = s[1]
	t.Id = s[2]
	t.Symbol = s[3]
	t.Value.Origin = s[4]
	i, err := strconv.Atoi(s[5])
	if err != nil {
		return fmt.Errorf("dnf tag value type parse failed,should be int, dnf：%s, err:%s", dnf, err)
	}
	t.TagType = i

	return
}