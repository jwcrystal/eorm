// Copyright 2021 gotomicro
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package test 是用于辅助测试的包。仅限于内部使用
package test

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/gotomicro/ekit"
)

// SimpleStruct 包含所有 eorm 支持的类型
type SimpleStruct struct {
	Id      uint64 `eorm:"primary_key,column=int_c"`
	Bool    bool
	BoolPtr *bool

	Int    int
	IntPtr *int

	Int8    int8 `eorm:"primary_key,column=int8_c"`
	Int8Ptr *int8

	Int16    int16
	Int16Ptr *int16

	Int32    int32
	Int32Ptr *int32

	Int64    int64
	Int64Ptr *int64

	Uint    uint
	UintPtr *uint

	Uint8    uint8
	Uint8Ptr *uint8

	Uint16    uint16
	Uint16Ptr *uint16

	Uint32    uint32
	Uint32Ptr *uint32

	Uint64    uint64
	Uint64Ptr *uint64

	Float32    float32
	Float32Ptr *float32

	Float64    float64
	Float64Ptr *float64

	// Byte byte
	// BytePtr *byte
	ByteArray []byte

	String string

	// 特殊类型
	NullStringPtr *sql.NullString
	NullInt16Ptr  *sql.NullInt16
	NullInt32Ptr  *sql.NullInt32
	NullInt64Ptr  *sql.NullInt64
	NullBoolPtr   *sql.NullBool
	// NullTimePtr    *sql.NullTime
	NullFloat64Ptr *sql.NullFloat64
	JsonColumn     *JsonColumn
}

// JsonColumn 是自定义的 JSON 类型字段
// Val 字段必须是结构体指针
type JsonColumn struct {
	Val   User
	Valid bool
}

type User struct {
	Name string
}

func (j *JsonColumn) Scan(src any) error {
	if src == nil {
		return nil
	}
	var bs []byte
	switch val := src.(type) {
	case string:
		bs = []byte(val)
	case []byte:
		bs = val
	case *[]byte:
		if val == nil {
			return nil
		}
		bs = *val
	default:
		return fmt.Errorf("不合法类型 %+v", src)
	}
	if len(bs) == 0 {
		return nil
	}
	err := json.Unmarshal(bs, &j.Val)
	if err != nil {
		return err
	}
	j.Valid = true
	return nil
}

// Value 参考 sql.NullXXX 类型定义的
func (j *JsonColumn) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	bs, err := json.Marshal(j.Val)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func NewSimpleStruct(id uint64) *SimpleStruct {
	return &SimpleStruct{
		Id:             id,
		Bool:           true,
		BoolPtr:        ekit.ToPtr[bool](false),
		Int:            12,
		IntPtr:         ekit.ToPtr[int](13),
		Int8:           8,
		Int8Ptr:        ekit.ToPtr[int8](-8),
		Int16:          16,
		Int16Ptr:       ekit.ToPtr[int16](-16),
		Int32:          32,
		Int32Ptr:       ekit.ToPtr[int32](-32),
		Int64:          64,
		Int64Ptr:       ekit.ToPtr[int64](-64),
		Uint:           14,
		UintPtr:        ekit.ToPtr[uint](15),
		Uint8:          8,
		Uint8Ptr:       ekit.ToPtr[uint8](18),
		Uint16:         16,
		Uint16Ptr:      ekit.ToPtr[uint16](116),
		Uint32:         32,
		Uint32Ptr:      ekit.ToPtr[uint32](132),
		Uint64:         64,
		Uint64Ptr:      ekit.ToPtr[uint64](164),
		Float32:        3.2,
		Float32Ptr:     ekit.ToPtr[float32](-3.2),
		Float64:        6.4,
		Float64Ptr:     ekit.ToPtr[float64](-6.4),
		ByteArray:      []byte("hello"),
		String:         "world",
		NullStringPtr:  &sql.NullString{String: "null string", Valid: true},
		NullInt16Ptr:   &sql.NullInt16{Int16: 16, Valid: true},
		NullInt32Ptr:   &sql.NullInt32{Int32: 32, Valid: true},
		NullInt64Ptr:   &sql.NullInt64{Int64: 64, Valid: true},
		NullBoolPtr:    &sql.NullBool{Bool: true, Valid: true},
		NullFloat64Ptr: &sql.NullFloat64{Float64: 6.4, Valid: true},
		JsonColumn: &JsonColumn{
			Val:   User{Name: "Tom"},
			Valid: true,
		},
	}
}

type BaseEntity struct {
	CreateTime uint64
	UpdateTime uint64
}

type CombinedModel struct {
	BaseEntity
	Id        int64 `eorm:"auto_increment,primary_key"`
	FirstName string
	Age       int8
	LastName  *string
}

func NewCombinedModel(id int64) *CombinedModel {
	return &CombinedModel{
		BaseEntity: BaseEntity{
			CreateTime: 10000,
			UpdateTime: 10000,
		},
		Id:        id,
		FirstName: "Tom" + fmt.Sprintln(id),
		Age:       20,
		LastName:  ekit.ToPtr[string]("Jerry" + fmt.Sprintln(id)),
	}
}

type Order struct {
	Id        int
	UsingCol1 string
	UsingCol2 string
}

type OrderDetail struct {
	OrderId int
	ItemId  int

	UsingCol1 string
	UsingCol2 string
}
type Item struct {
	Id int
}

func NewOrderModel(id int) *Order {
	return &Order{
		Id:        id,
		UsingCol1: "col1",
		UsingCol2: "col2",
	}
}

func NewOrderDetailModel(id int) *OrderDetail {
	return &OrderDetail{
		OrderId:   id,
		ItemId:    id + 1,
		UsingCol1: "col1",
		UsingCol2: "col2",
	}
}

func NewItemModel(id int) *Item {
	return &Item{
		Id: id,
	}
}
