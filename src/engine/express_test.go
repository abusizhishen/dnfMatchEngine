package engine

import (
	"errors"
	"testing"
)

type testParseCase struct {
	Dnf string
	Err error
}

var ErrDnfParse = errors.New("")

var testParseCases = []testParseCase{
	{
		Dnf: "┐[(age>{1}|0)]",
		Err: nil,
	},
	{
		"[(age!={1}|0)]",
		ErrDnfParse,
	},
	{
		"[(age≠{1}|0)]",
		nil,
	},
	{
		"[()]",
		nil,
	},
	{
		"[]",
		nil,
	},
	//{
	//	"()",
	//	nil,
	//},
}

func TestNew(t *testing.T) {
	fun := func(cas testParseCase) {
		_, err := New(cas.Dnf)
		//log.Println(cas.Dnf, engine)

		if err == nil && cas.Err == nil{
			//pass
			return
		}

		if err != nil && cas.Err != nil{
			//pass
			return
		}

		t.Log(err)
		t.Errorf("dnf parse failed dnf: %s result should be: %v",cas.Dnf,cas.Err == nil)
	}

	for _, dnf := range testParseCases {
		fun(dnf)
	}
}

type testCase struct {
	Dnf    string
	Data   map[string]string
	Result bool
}

var testMatchData = []testCase{
	{
		"[(age<{1}|0)]",
		map[string]string{"age": "0"},
		true,
	},
	{
		"[(age={1}|0)]",
		map[string]string{"age": "1"},
		true,
	},
	{
		"[(age>{1}|0)]",
		map[string]string{"age": "2"},
		true,
	},
	{
		"[(age≥{1}|0)]",
		map[string]string{"age": "2"},
		true,
	},
	{
		"[(age≠{1}|0)]",
		map[string]string{"age": "2"},
		true,
	},
	{
		"[]",
		map[string]string{"age": "2"},
		true,
	},
	{
		"[(age≠{1}|0)]",
		map[string]string{"name": "lisa"},
		false,
	},
}

func TestAdvertTag_Match(t *testing.T) {
	test := func(cas testCase) bool {
		engine, err := New(cas.Dnf)
		if err != nil {
			t.Fatal(err)
		}

		result,err :=  engine.Match(cas.Data)
		if err != nil{
			t.Logf("case failed:%+v, err:%v",cas,err)
		}
		return result == cas.Result
	}

	for _, cas := range testMatchData {
		if !test(cas) {
			t.Errorf("test failed result should be %v case: [dnf:%s, data:%v, result:%v]", !cas.Result, cas.Dnf, cas.Data, cas.Result)
		}
	}
}
