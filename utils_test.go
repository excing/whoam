package main

import (
	"testing"
)

func TestGen(t *testing.T) {
	t.Log(RandNdigMbitString(4))
	t.Log(RandNdigMbitString(4, 16))
	t.Log(RandNdigMbitString(4, 36))
	t.Log(RandNdigMbitString(48))
	t.Log(RandNdigMbitString(48, 16))
	t.Log(RandNdigMbitString(48, 36))
	t.Log(RandNdigMbitString(48, 26, 36))
	t.Log(RandNdigMbitString(108))
}

func TestGenFor(t *testing.T) {
	for i := 0; i < 100; i++ {
		TestGen(t)
	}
}

func TestGenID(t *testing.T) {
	t.Log(New4bitID())
	t.Log(New4BitID())
	t.Log(New16bitID())
	t.Log(New16BitID())
	t.Log(New32bitID())
	t.Log(New32BitID())
	t.Log(New64BitID())
}

func TestNew128BitID(t *testing.T) {
	for i := 128; 0 < i; i-- {
		t.Log(New128BitID())
	}
}
