// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
)

func IsJSON(content string) bool {
	return json.Valid([]byte(content))
}

func JSONNormalize(content *string) error {
	return jsonTransform(content, false)
}

func JSONNormalizePretty(content *string) error {
	return jsonTransform(content, true)
}

func jsonTransform(content *string, prettyJSON bool) error { //revive:disable-line:flag-parameter
	var result any
	if err := json.Unmarshal([]byte(*content), &result); err != nil {
		return err
	}

	var resultRaw []byte
	var err error

	if prettyJSON {
		resultRaw, err = json.MarshalIndent(result, "", "  ")
	} else {
		resultRaw, err = json.Marshal(result)
	}

	if err != nil {
		return err
	}

	*content = strings.TrimSpace(string(resultRaw))

	return nil
}

func JSONBase64GzipEncode(content *string) error {
	if err := JSONNormalize(content); err != nil {
		return err
	}

	if err := Base64GzipEncode(content); err != nil {
		return err
	}

	return nil
}

func JSONBase64GzipDecode(content string) (any, error) {
	contentPtr := &content
	if err := Base64GzipDecode(contentPtr); err != nil {
		return nil, err
	}

	var result any

	err := json.Unmarshal([]byte(*contentPtr), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Base64Decode(content *string) error {
	if content == nil {
		return nil
	}

	payload, err := base64ToByte(*content)
	if err != nil {
		return err
	}

	*content = strings.TrimSpace(string(payload))

	return nil
}

func Base64Encode(content *string) error {
	if content == nil {
		return nil
	}

	payload := byteToBase64([]byte(*content))
	*content = strings.TrimSpace(payload)

	return nil
}

func Base64GzipEncode(content *string) error {
	if content == nil {
		return nil
	}

	var bu bytes.Buffer
	gz := gzip.NewWriter(&bu)

	if _, err := gz.Write([]byte(*content)); err != nil {
		return err
	}

	if err := gz.Flush(); err != nil {
		return err
	}

	if err := gz.Close(); err != nil {
		return err
	}

	*content = byteToBase64(bu.Bytes())

	return nil
}

func Base64GzipDecode(content *string) error {
	if content == nil {
		return nil
	}

	data, err := base64ToByte(*content)
	if err != nil {
		return nil
	}

	rd, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil
	}
	defer rd.Close()

	dataDecomp, err := io.ReadAll(rd)
	if err != nil {
		return nil
	}

	*content = string(dataDecomp)

	return nil
}

func base64ToByte(content string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func byteToBase64(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}
