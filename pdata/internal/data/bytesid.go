// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package data // import "go.opentelemetry.io/collector/pdata/internal/data"

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
)

var escapeForwardslash = []byte(`\/`)
var forwardslash = []byte(`/`)
var escapePlus = []byte(`\+`)
var plus = []byte(`+`)

// marshalJSON converts trace id into a hex string enclosed in quotes.
// Called by Protobuf JSON deserialization.
func marshalJSON(id []byte) ([]byte, error) {
	// Plus 2 quote chars at the start and end.
	hexLen := hex.EncodedLen(len(id)) + 2

	b := make([]byte, hexLen)
	hex.Encode(b[1:hexLen-1], id)
	b[0], b[hexLen-1] = '"', '"'

	return b, nil
}

// unmarshalJSON inflates trace id from hex string, possibly enclosed in quotes.
// Called by Protobuf JSON deserialization.
func unmarshalJSON(dst []byte, src []byte) error {
	if l := len(src); l >= 2 && src[0] == '"' && src[l-1] == '"' {
		src = src[1 : l-1]
	}
	nLen := len(src)
	if nLen == 0 {
		return nil
	}

	isHex := true
	var isCustomErr error
	if len(dst) == hex.DecodedLen(nLen) {
		_, err := hex.Decode(dst, src)
		if err != nil {
			isCustomErr = fmt.Errorf("cannot unmarshal ID from string '%s': %w", string(src), err)
			isHex = false
		}
	} else {
		isCustomErr = errors.New("invalid length for ID")
		isHex = false
	}

	if isHex {
		return isCustomErr
	}

	// StdEncode the base64 decoded ID. The result is stored in dst
	_, err1 := base64.StdEncoding.Decode(dst, src)

	if err1 != nil {
		log.Printf("error while decoding src using std encoding '%s': %v", string(src), err1)
		replace := bytes.ReplaceAll(src, escapeForwardslash, forwardslash)
		replace = bytes.ReplaceAll(replace, escapePlus, plus)
		// Retry to decode
		_, err2 := base64.StdEncoding.Decode(dst, replace)
		if err2 != nil {
			return fmt.Errorf("error while decoding modified src using std encoding '%s': %w", string(src), err2)
		}
	}

	// Encode the ID bytes to hex
	hexEncodedBytes := make([]byte, hex.EncodedLen(len(dst)))
	hex.Encode(hexEncodedBytes, dst)

	// Check the length of the hex ID bytes
	nLen1 := len(hexEncodedBytes)
	if nLen1 == 0 {
		return nil
	}

	// Decode the hex encoded bytes
	hexDecodedBytes := make([]byte, hex.DecodedLen(len(hexEncodedBytes)))
	_, err := hex.Decode(hexDecodedBytes, hexEncodedBytes)
	if err != nil {
		return fmt.Errorf("cannot unmarshal ID from string '%s': %w", string(hexEncodedBytes), err)
	}

	// Original code
	/*if len(dst) != hex.DecodedLen(nLen) {
		return errors.New("invalid length for ID")
	}
	_, err := hex.Decode(dst, src)
	if err != nil {
		return fmt.Errorf("cannot unmarshal ID from string '%s': %w", string(src), err)
	}
	}*/

	return nil
}
