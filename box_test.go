package main

import (
	"math/rand"
	"testing"
	"time"
)

func TestBox(t *testing.T) {
	var src = rand.NewSource(time.Now().UnixNano())
	var box = NewBox(3*1024, 3)
	var vals = make(map[int]interface{})

	for i := 120; 0 < i; i-- {
		var k = i
		var m = src.Int63()
		vals[k] = m
		box.SetInt64ValI(k, m)
		k = i + 1000
		var s = New64BitID()
		vals[k] = s
		box.SetStringValI(k, s)
	}

	for i := 1; i <= 120; i++ {
		var k = i
		bint, err := box.Int64ValI(k)
		vint := vals[k]
		t.Log(k, bint == vint, err)
		k = i + 1000
		bstr, err := box.StringValI(k)
		vstr := vals[k]
		t.Log(k, bstr == vstr, err)
	}
}

func TestBoxByBool(t *testing.T) {
	box := NewBox(3*1024, 1)
	box.SetBoolVal("test", false)
	ok, err := box.BoolVal("test")
	t.Log("test: ", ok, err)
	box.SetBoolVal("test", true)
	ok, err = box.BoolVal("test")
	t.Log("test: ", ok, err)
	ok, err = box.BoolVal("test1")
	t.Log("test: ", ok, err)
}

func TestBoxByInt64(t *testing.T) {
	box := NewBox(3*1024, 1)
	box.SetInt64Val("test", 3930392034232434291)
	val, err := box.Int64Val("test")
	t.Log("test: ", val, err)
	box.SetInt64Val("test", -93999391)
	val, err = box.Int64Val("test")
	t.Log("test: ", val, err)
	val, err = box.Int64Val("test1")
	t.Log("test: ", val, err)
}

func TestBoxByFloat64(t *testing.T) {
	box := NewBox(3*1024, 1)
	box.SetFloat64Val("test", -3930392034232434291.0939209329020402948)
	val, err := box.Float64Val("test")
	t.Log("test: ", val, err)
	box.SetFloat64Val("test", 93999391.939382048922028)
	val, err = box.Float64Val("test")
	t.Log("test: ", val, err)
	val, err = box.Float64Val("test1")
	t.Log("test: ", val, err)
}

func TestBoxByStruct(t *testing.T) {
	box := NewBox(3*1024, 1)
	var test = struct {
		Name string
	}{
		Name: "Tom cat",
	}
	box.SetVal("test", &test)
	test.Name = "Tom dog"
	t.Log("test: ", test)
	err := box.Val("test", &test)
	t.Log("test: ", test, err)
	var wrong string
	err = box.Val("test", &wrong)
	t.Log("test: ", wrong, err)
}
