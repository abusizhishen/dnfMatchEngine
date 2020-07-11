package engine

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testParseCase struct {
	Dnf string
	Err error
}

var ErrDnfParse = errors.New("")

func TestNew(t *testing.T) {
	t.Run("IntGt", func(t *testing.T) {
		e,err := New("[(age>{1}|0)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("IntEq", func(t *testing.T) {
		e,err := New("[(age={1}|0)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("IntLt", func(t *testing.T) {
		e,err := New("[(age≤{1}|0)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("StringEq", func(t *testing.T) {
		e,err := New("[(name={foo}|3)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("StringNotEq", func(t *testing.T) {
		e,err := New("[(name≠{foo}|3)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("TimeGT", func(t *testing.T) {
		e,err := New("[(date>{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("TimeGTOrEq", func(t *testing.T) {
		e,err := New("[(date≥{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("TimeEq", func(t *testing.T) {
		e,err := New("[(date={2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("TimeNotEq", func(t *testing.T) {
		e,err := New("[(date≠{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("TimeLTOrEq", func(t *testing.T) {
		e,err := New("[(date≤{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("TimeLT", func(t *testing.T) {
		e,err := New("[(date<{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("InArray", func(t *testing.T) {
		e,err := New("[(size∈{small,large}|7)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("NotInArray", func(t *testing.T) {
		e,err := New("[(size∉{small,large}|7)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("EqInArray", func(t *testing.T) {
		e,err := New("[(size={small,large}|7)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("NotEqInArray", func(t *testing.T) {
		e,err := New("[(size≠{small,large}|7)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("And", func(t *testing.T) {
		e,err := New("[(age>{1}|0)^(age>{3}|0)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})	})

	t.Run("Or", func(t *testing.T) {
		e,err := New("[(age>{1}|0)∨(age>{2}|0)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})

	t.Run("Not", func(t *testing.T) {
		e,err := New("┐[(age>{1}|0)]")
		if err != nil{
			t.Error(err)
		}
		assert.IsType(t,e, &AdvertTag{})
	})
}

type testCase struct {
	Dnf    string
	Data   map[string]string
	Result bool
}

func TestAdvertTag_Match(t *testing.T) {
	t.Run("IntGtMatch", func(t *testing.T) {
		engine,err := New("[(age>{1}|0)]")
		data := map[string]string{"age":"2"}
		result,err := engine.Match(data)
		if err != nil{
			t.Error(err)
		}
		assert.True(t, result)
	})

	t.Run("IntGtNotMatch", func(t *testing.T) {
		engine,err := New("[(age>{1}|0)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"age":"1"}
		result,err := engine.Match(data)
		if err != nil{
			t.Error(err)
		}
		assert.False(t, result)
	})

	t.Run("IntEqMatch", func(t *testing.T) {
		engine,err := New("[(age={1}|0)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"age":"1"}
		result,err := engine.Match(data)
		if err != nil{
			t.Error(err)
		}
		assert.True(t, result)
	})

	t.Run("IntEqNotMatch", func(t *testing.T) {
		engine,err := New("[(age={1}|0)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"age":"3"}
		result,err := engine.Match(data)
		if err != nil{
			t.Error(err)
		}
		assert.False(t, result)
	})

	t.Run("IntLtMatch", func(t *testing.T) {
		engine,err := New("[(age≤{1}|0)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"age":"0"}
		result,err := engine.Match(data)
		if err != nil{
			t.Error(err)
		}
		assert.True(t, result)
	})

	t.Run("IntLtNotMatch", func(t *testing.T) {
		engine,err := New("[(age≤{1}|0)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"age":"2"}
		result,err := engine.Match(data)
		if err != nil{
			t.Error(err)
		}
		assert.False(t, result)
	})
	
	t.Run("StringEqMatch", func(t *testing.T) {
		e,err := New("[(name={foo}|3)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"name":"foo"}
		result,err := e.Match(data)
		if err != nil{
			t.Error(err)
		}

		assert.True(t,result)
	})

	t.Run("StringEqNotMatch", func(t *testing.T) {
		e,err := New("[(name={foo}|3)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"name":"boo"}
		result,err := e.Match(data)
		if err != nil{
			t.Error(err)
		}

		assert.False(t,result)
	})

	t.Run("TimeGTMatchGT", func(t *testing.T) {
		e,err := New("[(date>{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"date":"2020-07-01 00:00:01"}
		result,err := e.Match(data)
		if err != nil{
			t.Error(err)
		}

		assert.True(t,result)
	})

	t.Run("TimeGTMatchEQ", func(t *testing.T) {
		e,err := New("[(date>{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"date":"2020-07-01 00:00:00"}
		result,err := e.Match(data)
		if err != nil{
			t.Error(err)
		}

		assert.False(t,result)
	})

	t.Run("TimeGTNotMatchLt", func(t *testing.T) {
		e,err := New("[(date>{2020-07-01 00:00:00}|4)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"date":"2020-06-30 23:59:59"}
		result,err := e.Match(data)
		if err != nil{
			t.Error(err)
		}

		assert.False(t,result)
	})


	t.Run("InArrayMatchIn", func(t *testing.T) {
		e,err := New("[(size∈{small,large}|7)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"size":"small"}
		result,err := e.Match(data)
		if err != nil{
			t.Error(err)
		}

		assert.True(t,result)
	})

	t.Run("InArrayMatchNotIn", func(t *testing.T) {
		e,err := New("[(size∈{small,large}|7)]")
		if err != nil{
			t.Error(err)
		}
		data := map[string]string{"size":"huge"}
		result,err := e.Match(data)
		if err != nil{
			t.Error(err)
		}

		assert.False(t,result)
	})
}