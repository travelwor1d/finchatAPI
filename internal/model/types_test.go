package model

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringSlice(t *testing.T) {
	const data = `{"tags":["first","second","last"]}`
	slice := StringSlice{"first", "second", "last"}
	commaSep := "first,second,last"

	t.Run("json unmarshaling", func(t *testing.T) {
		s := struct {
			Tags StringSlice `json:"tags"`
		}{}
		if err := json.Unmarshal([]byte(data), &s); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(s.Tags, slice) {
			t.Errorf("does not match: %s != %s", s.Tags, slice)
		}
	})

	t.Run("json marshaling", func(t *testing.T) {
		s := struct {
			Tags StringSlice `json:"tags"`
		}{Tags: slice}
		byt, err := json.Marshal(&s)
		if err != nil {
			t.Fatal(err)
		}
		if string(byt) != data {
			t.Errorf("does not match: %s != %s", string(byt), data)
		}
	})

	t.Run("StringSlice.Scan", func(t *testing.T) {
		var s StringSlice
		if err := s.Scan(commaSep); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(s, slice) {
			t.Errorf("does not match: %s != %s", s, slice)
		}
	})

	t.Run("StringSlice.Value", func(t *testing.T) {
		s, err := slice.Value()
		if err != nil {
			t.Fatal(err)
		}
		if s != commaSep {
			t.Errorf("does not match: %s != %s", s, commaSep)
		}
	})
}
