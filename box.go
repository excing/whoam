package main

import (
	"encoding/binary"
	"encoding/json"
	"math"

	"github.com/coocood/freecache"
)

// Box a freecache inheritance structure
type Box struct {
	freecache.Cache
	defaultTimeout int
}

// NewBox return new box
func NewBox(size int, defaultTimeout int) *Box {
	return &Box{*freecache.NewCache(size), defaultTimeout}
}

// DelString return true, if delete the key fails, return false
func (box *Box) DelString(key string) bool {
	return box.Del([]byte(key))
}

// SetVal set key-value
func (box *Box) SetVal(key string, val interface{}, timeout ...int) error {
	bytes, err := json.Marshal(val)

	if err != nil {
		return err
	}

	return box.setVal([]byte(key), bytes, timeout...)
}

// Val parse the value of the key into dst
func (box *Box) Val(key string, dst interface{}) error {
	bytes, err := box.Get([]byte(key))
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &dst)

	return err
}

// SetValI set int key-value
func (box *Box) SetValI(key int, val interface{}, timeout ...int) error {
	bytes, err := json.Marshal(val)

	if err != nil {
		return err
	}

	return box.setVal(IntToBytes(key), bytes, timeout...)
}

// ValI parse the value of the int key into dst
func (box *Box) ValI(key int, dst interface{}) error {
	bytes, err := box.Get(IntToBytes(key))
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &dst)

	return err
}

// SetStringVal set key-string value
func (box *Box) SetStringVal(key string, val string, timeout ...int) error {
	return box.setVal([]byte(key), []byte(val), timeout...)
}

// StringVal return string value
func (box *Box) StringVal(key string) (string, error) {
	bytes, err := box.Get([]byte(key))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetStringValI set int key-string value
func (box *Box) SetStringValI(key int, val string, timeout ...int) error {
	return box.setVal(IntToBytes(key), []byte(val), timeout...)
}

// StringValI return string value
func (box *Box) StringValI(key int) (string, error) {
	bytes, err := box.Get(IntToBytes(key))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetByteVal set key-byte value
func (box *Box) SetByteVal(key string, val byte, timeout ...int) error {
	bytes := make([]byte, 1)
	bytes[0] = val
	return box.setVal([]byte(key), bytes, timeout...)
}

// ByteVal return byte value
func (box *Box) ByteVal(key string) (byte, error) {
	bytes, err := box.Get([]byte(key))
	if err != nil {
		return 0, err
	}

	return bytes[0], nil
}

// SetByteValI set int key-byte value
func (box *Box) SetByteValI(key int, val byte, timeout ...int) error {
	bytes := make([]byte, 1)
	bytes[0] = val
	return box.setVal(IntToBytes(key), bytes, timeout...)
}

// ByteValI return byte value
func (box *Box) ByteValI(key int) (byte, error) {
	bytes, err := box.Get(IntToBytes(key))
	if err != nil {
		return 0, err
	}

	return bytes[0], nil
}

// SetUint64Val set key-uint64 value
func (box *Box) SetUint64Val(key string, val uint64, timeout ...int) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes[:], val)
	return box.setVal([]byte(key), bytes, timeout...)
}

// Uint64Val return uint64 value
func (box *Box) Uint64Val(key string) (uint64, error) {
	bytes, err := box.Get([]byte(key))
	if err != nil {
		return 0, err
	}

	x := binary.BigEndian.Uint64(bytes)

	return x, nil
}

// SetUint64ValI set int key-uint64 value
func (box *Box) SetUint64ValI(key int, val uint64, timeout ...int) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes[:], val)
	return box.setVal(IntToBytes(key), bytes, timeout...)
}

// Uint64ValI return uint64 value
func (box *Box) Uint64ValI(key int) (uint64, error) {
	bytes, err := box.Get(IntToBytes(key))
	if err != nil {
		return 0, err
	}

	x := binary.BigEndian.Uint64(bytes)

	return x, nil
}

// SetInt64Val set key-int64 value
func (box *Box) SetInt64Val(key string, val int64, timeout ...int) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes[:], uint64(val))
	return box.setVal([]byte(key), bytes, timeout...)
}

// Int64Val return int64 value
func (box *Box) Int64Val(key string) (int64, error) {
	bytes, err := box.Get([]byte(key))
	if err != nil {
		return 0, err
	}

	x := binary.BigEndian.Uint64(bytes)

	return int64(x), nil
}

// SetIntVal set key-int value
func (box *Box) SetIntVal(key string, val int, timeout ...int) error {
	return box.SetInt64Val(key, int64(val), timeout...)
}

// IntVal return int value
func (box *Box) IntVal(key string) (int, error) {
	val, err := box.Int64Val(key)
	return int(val), err
}

// SetInt64ValI set int key-int64 value
func (box *Box) SetInt64ValI(key int, val int64, timeout ...int) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes[:], uint64(val))
	return box.setVal(IntToBytes(key), bytes, timeout...)
}

// Int64ValI return int64 value
func (box *Box) Int64ValI(key int) (int64, error) {
	bytes, err := box.Get(IntToBytes(key))
	if err != nil {
		return 0, err
	}

	x := binary.BigEndian.Uint64(bytes)

	return int64(x), nil
}

// SetIntValI set int key-int value
func (box *Box) SetIntValI(key int, val int, timeout ...int) error {
	return box.SetInt64ValI(key, int64(val), timeout...)
}

// IntValI return int value
func (box *Box) IntValI(key int) (int, error) {
	val, err := box.Int64ValI(key)
	return int(val), err
}

// SetFloat64Val set key-float64 value
func (box *Box) SetFloat64Val(key string, val float64, timeout ...int) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes[:], math.Float64bits(val))
	return box.setVal([]byte(key), bytes, timeout...)
}

// Float64Val return float64 value
func (box *Box) Float64Val(key string) (float64, error) {
	bytes, err := box.Get([]byte(key))
	if err != nil {
		return 0, err
	}

	x := math.Float64frombits(binary.BigEndian.Uint64(bytes))

	return x, nil
}

// SetFloat64ValI set int key-float64 value
func (box *Box) SetFloat64ValI(key int, val float64, timeout ...int) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes[:], math.Float64bits(val))
	return box.setVal(IntToBytes(key), bytes, timeout...)
}

// Float64ValI return float64 value
func (box *Box) Float64ValI(key int) (float64, error) {
	bytes, err := box.Get(IntToBytes(key))
	if err != nil {
		return 0, err
	}

	x := math.Float64frombits(binary.BigEndian.Uint64(bytes))

	return x, nil
}

// SetBoolVal set key-bool value
func (box *Box) SetBoolVal(key string, val bool, timeout ...int) error {
	bytes := make([]byte, 1)
	if val {
		bytes[0] = 1
	} else {
		bytes[0] = 0
	}
	return box.setVal([]byte(key), bytes, timeout...)
}

// BoolVal return bool value
func (box *Box) BoolVal(key string) (bool, error) {
	bytes, err := box.Get([]byte(key))
	if err != nil {
		return false, err
	}

	x := bytes[0] != 0

	return x, nil
}

// SetBoolValI set int key-bool value
func (box *Box) SetBoolValI(key int, val bool, timeout ...int) error {
	bytes := make([]byte, 1)
	if val {
		bytes[0] = 1
	} else {
		bytes[0] = 0
	}
	return box.setVal(IntToBytes(key), bytes, timeout...)
}

// BoolValI return bool value
func (box *Box) BoolValI(key int) (bool, error) {
	bytes, err := box.Get(IntToBytes(key))
	if err != nil {
		return false, err
	}

	x := bytes[0] != 0

	return x, nil
}

func (box *Box) setVal(key []byte, val []byte, timeout ...int) error {
	if 0 == len(timeout) {
		return box.Set([]byte(key), val, box.defaultTimeout)
	}

	return box.Set([]byte(key), val, timeout[0])
}

// IntToBytes return bytes from int
func IntToBytes(i int) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs[:], uint64(i))
	return bs
}
