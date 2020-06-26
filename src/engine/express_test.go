package engine

import "testing"

func TestNew(t *testing.T)  {
	var dnfs = []string{
		"┐[(age>{1}|0)]",
		"[(age!={1}|0)]",
		"[(age≠{1}|0)]",
	}

	fun := func(dnf string) {
		_,err := New(dnf)
		if err != nil{
			t.Error(err)
		}
	}

	for _,dnf := range dnfs{
		fun(dnf)
	}
}

func TestAdvertTag_Match(t *testing.T)  {
	type testCase struct {
		Dnf string
		Data map[string]string
		Result bool
	}

	testData := []testCase{
		{
			"[(age<{1}|0)]",
			map[string]string{"age":"0"},
			true,
		},
		{
			"[(age={1}|0)]",
			map[string]string{"age":"1"},
			true,
		},
		{
			"[(age>{1}|0)]",
			map[string]string{"age":"2"},
			true,
		},
		{
			"[(age≥{1}|0)]",
			map[string]string{"age":"2"},
			true,
		},
		{
			"[(age!={1}|0)]",
			map[string]string{"age":"2"},
			true,
		},
	}

	test := func(cas testCase) bool{
		engine,err := New(cas.Dnf)
		if err != nil{
			t.Fatal(err)
		}

		return engine.Match(cas.Data) == cas.Result
	}

	for _,cas := range testData{
		if !test(cas){
			t.Errorf("test failed result should be %v case: [dnf:%s, data:%v, result:%v]", !cas.Result, cas.Data,cas.Data,cas.Result)
		}
	}
}
