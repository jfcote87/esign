// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package model_test

import (
	"encoding/json"
	"testing"

	"github.com/pwaterz/esign/v2/model"
)

func TestDSBool(t *testing.T) {
	type DSBoolTester struct {
		A model.DSBool `json:"a,omitempty"`
		B model.DSBool `json:"b"`
		C model.DSBool `json:"c,omitempty"`
	}
	var ex = &DSBoolTester{
		A: true,
		B: false,
	}
	b, _ := json.Marshal(ex)
	expect := `{"a":"true","b":"false"}`
	if string(b) != expect {
		t.Errorf("expected %s; got %s", expect, b)
		return
	}

	ex = nil

	json.Unmarshal([]byte(`{"c":"true","b":"false"}`), &ex)
	if ex.A || ex.B || !ex.C {
		t.Errorf("expected A: false, B: false, C: true; got %v", ex)
	}
}

func TestTabRequired(t *testing.T) {
	type TabRequiredTester struct {
		A model.TabRequired `json:"a,omitempty"`
		B model.TabRequired `json:"b,omitempty"`
		C model.TabRequired `json:"c"`
		D model.TabRequired `json:"d"`
	}
	ex := &TabRequiredTester{
		A: model.REQUIRED_DEFAULT,
		B: model.REQUIRED_FALSE,
		C: model.REQUIRED_TRUE,
		D: model.REQUIRED_FALSE,
	}
	b, _ := json.Marshal(ex)
	expect := `{"b":"false","c":"true","d":"false"}`
	if string(b) != expect {
		t.Errorf("expected %s; got %s", expect, b)
		return
	}

	ex = nil
	json.Unmarshal([]byte(`{"c":"false","b":"true","a":"false"}`), &ex)
	if ex.A.IsRequired() || !ex.B.IsRequired() || ex.C.IsRequired() || !ex.D.IsRequired() {
		t.Errorf("expected A: false, B: true, C: false, D: true; got A: %v, B: %v, C: %v, D: %v",
			ex.A.IsRequired(), ex.B.IsRequired(), ex.C.IsRequired(), ex.D.IsRequired())
	}
}
