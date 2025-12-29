// Copyright 2019 James Cote
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package esign_test

import (
	"bytes"
	"encoding/xml"
	"os"
	"testing"
	"time"

	"github.com/jacobwilson41/esign"
)

func TestXML(t *testing.T) {
	_ = bytes.NewBufferString("")
	f, err := os.Open("testdata/connect.xml")
	if err != nil {
		t.Fatalf("Open Connect.xml: %v", err)
		return
	}

	var v *esign.ConnectData = &esign.ConnectData{}
	decoder := xml.NewDecoder(f)
	err = decoder.Decode(v)
	if err != nil {
		t.Fatalf("XML Decode: %v", err)
		return
	}
	if v.EnvelopeStatus.DocumentStatuses[0].Name != "Docusign1.pdf" {
		t.Errorf("invalid document name in connect XML: %s", v.EnvelopeStatus.DocumentStatuses[0].Name)
	}
	var delivered = v.EnvelopeStatus.RecipientStatuses[0].Delivered
	if !delivered.Time().Equal(time.Date(2014, 11, 11, 13, 43, 40, 887000000, time.UTC)) {
		t.Errorf("expected Delivered = 2014-11-11 13:43:40.887 +0000 UTC; got %v", delivered.Time())
	}
	if !v.EnvelopeStatus.Signed.Time().Equal(time.Date(2014, 11, 11, 13, 44, 23, 590000000, time.UTC)) {
		t.Errorf("expected Signed = 2014-11-11 13:43:45.3; got %v", v.EnvelopeStatus.Signed.Time())
	}
}
