// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/struct"
	"reflect"
	"testing"
)

func TestJsonStruct_Contains(t *testing.T) {
	mapVal := NewJsonStruct(&structpb.Struct{Fields: map[string]*structpb.Value{
		"first":  {Kind: &structpb.Value_StringValue{"hello"}},
		"second": {Kind: &structpb.Value_NumberValue{1}}}})
	if !mapVal.Contains(String("first")).(Bool) {
		t.Error("Expected map to contain key 'first'", mapVal)
	}
	if mapVal.Contains(String("firs")).(Bool) {
		t.Error("Expected map contained non-existent key", mapVal)
	}
}

func TestJsonStruct_ConvertToNative_Error(t *testing.T) {
	val, err := NewJsonStruct(&structpb.Struct{}).ConvertToNative(jsonListValueType)
	if err == nil {
		t.Errorf("Unsupported type conversion succeeded. "+
			"Got '%v', expected error", val)
	}
}

func TestJsonStruct_ConvertToNative_Json(t *testing.T) {
	structVal := &structpb.Struct{Fields: map[string]*structpb.Value{
		"first":  {Kind: &structpb.Value_StringValue{"hello"}},
		"second": {Kind: &structpb.Value_NumberValue{1}}}}
	mapVal := NewJsonStruct(structVal)
	val, err := mapVal.ConvertToNative(jsonValueType)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(val.(proto.Message),
		&structpb.Value{Kind: &structpb.Value_StructValue{structVal}}) {
		t.Error("Got '%v', expected '%v'", val, structVal)
	}

	strVal, err := mapVal.ConvertToNative(jsonStructType)
	if err != nil {
		t.Error(err)
	}
	if !proto.Equal(strVal.(proto.Message), structVal) {
		t.Error("Got '%v', expected '%v'", strVal, structVal)
	}
}

func TestJsonStruct_ConvertToNative_Any(t *testing.T) {
	structVal := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"first":  {Kind: &structpb.Value_StringValue{"hello"}},
			"second": {Kind: &structpb.Value_NumberValue{1}}}}
	mapVal := NewJsonStruct(structVal)
	anyVal, err := mapVal.ConvertToNative(anyValueType)
	if err != nil {
		t.Error(err)
	}
	unpackedAny := ptypes.DynamicAny{}
	if ptypes.UnmarshalAny(anyVal.(*any.Any), &unpackedAny) != nil {
		t.Error("Failed to unmarshal any")
	}
	if !proto.Equal(unpackedAny.Message, mapVal.Value().(proto.Message)) {
		t.Error("Messages were not equal, got '%v'", unpackedAny.Message)
	}
}

func TestJsonStruct_ConvertToNative_Map(t *testing.T) {
	structVal := &structpb.Struct{Fields: map[string]*structpb.Value{
		"first":  {Kind: &structpb.Value_StringValue{"hello"}},
		"second": {Kind: &structpb.Value_StringValue{"world"}}}}
	mapVal := NewJsonStruct(structVal)
	val, err := mapVal.ConvertToNative(reflect.TypeOf(map[string]string{}))
	if err != nil {
		t.Error(err)
	}
	if val.(map[string]string)["first"] != "hello" {
		t.Error("Could not find key 'first' in map", val)
	}
}

func TestJsonStruct_ConvertToType(t *testing.T) {
	mapVal := NewJsonStruct(&structpb.Struct{Fields: map[string]*structpb.Value{
		"first":  {Kind: &structpb.Value_StringValue{"hello"}},
		"second": {Kind: &structpb.Value_NumberValue{1}}}})
	if mapVal.ConvertToType(MapType) != mapVal {
		t.Error("Map could not be converted to a map.")
	}
	if mapVal.ConvertToType(TypeType) != MapType {
		t.Error("Map did not indicate itself as map type.")
	}
	if !IsError(mapVal.ConvertToType(ListType)) {
		t.Error("Got list, expected error.")
	}
}

func TestJsonStruct_Equal(t *testing.T) {
	mapVal := NewJsonStruct(&structpb.Struct{Fields: map[string]*structpb.Value{
		"first":  {Kind: &structpb.Value_StringValue{"hello"}},
		"second": {Kind: &structpb.Value_StringValue{"1"}}}})

	otherVal := NewJsonStruct(&structpb.Struct{Fields: map[string]*structpb.Value{
		"first":  {Kind: &structpb.Value_StringValue{"hello"}},
		"second": {Kind: &structpb.Value_NumberValue{1}}}})
	if mapVal.Equal(otherVal) != False {
		t.Errorf("Got equals 'true', expected 'false' for '%v' == '%v'",
			mapVal, otherVal)
	}
	if mapVal.Equal(mapVal) != True {
		t.Error("Map was not equal to itself.")
	}
	if mapVal.Equal(NewJsonStruct(&structpb.Struct{})) != False {
		t.Error("Map with key-value pairs was equal to empty map")
	}
	if mapVal.Equal(String("")) != False {
		t.Error("Map was equal to a non-map type.")
	}
}

func TestJsonStruct_Get(t *testing.T) {
	if !IsError(NewJsonStruct(&structpb.Struct{}).Get(Int(1))) {
		t.Error("Structs may only have string keys.")
	}
}
