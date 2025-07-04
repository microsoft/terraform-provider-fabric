// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

// Package transforms provides utility functions for data transformation operations
// like JSON normalization, base64 encoding/decoding, and gzip compression/decompression.
package transforms

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrEmptyContent = errors.New("empty content")
	ErrEmptyJSON    = errors.New("empty JSON content")
)

// IsJSON checks if a string contains valid JSON.
func IsJSON(content string) bool {
	return json.Valid([]byte(content))
}

// JSONNormalize normalizes JSON content to a compact format.
func JSONNormalize(content string) (string, error) {
	return jsonTransform(content, false)
}

// JSONNormalizePretty normalizes JSON content to a human-readable format with indentation.
func JSONNormalizePretty(content string) (string, error) {
	return jsonTransform(content, true)
}

// jsonTransform is a helper function that unmarshals and marshals JSON to normalize it.
// If prettyJSON is true, the output will be indented.
func jsonTransform(content string, prettyJSON bool) (string, error) { //revive:disable-line:flag-parameter
	if content == "" {
		return "", ErrEmptyJSON
	}

	var result any
	var err error

	err = json.Unmarshal([]byte(content), &result)
	if err != nil {
		return content, fmt.Errorf("invalid JSON format: %w", err)
	}

	var resultRaw []byte

	if prettyJSON {
		resultRaw, err = json.MarshalIndent(result, "", "  ")
	} else {
		resultRaw, err = json.Marshal(result)
	}

	if err != nil {
		return content, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return strings.TrimSpace(string(resultRaw)), nil
}

// JSONBase64GzipEncode normalizes, compresses, and base64 encodes JSON content.
func JSONBase64GzipEncode(content string) (string, error) {
	if content == "" {
		return "", ErrEmptyJSON
	}

	normalized, err := JSONNormalize(content)
	if err != nil {
		return "", err
	}

	encoded, err := Base64GzipEncode(normalized)
	if err != nil {
		return "", err
	}

	return encoded, nil
}

// JSONBase64GzipDecode decodes, decompresses and unmarshals base64 gzipped JSON content.
func JSONBase64GzipDecode(content string) (any, error) {
	if content == "" {
		return nil, ErrEmptyContent
	}

	decoded, err := Base64GzipDecode(content)
	if err != nil {
		return nil, err
	}

	var result any

	err = json.Unmarshal([]byte(decoded), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

// Base64Encode encodes a string or byte slice to base64.
func Base64Encode[T string | []byte](content T) (string, error) {
	switch v := any(content).(type) {
	case string:
		if v == "" {
			return "", nil
		}

		return byteToBase64([]byte(v)), nil
	case []byte:
		if len(v) == 0 {
			return "", nil
		}

		return byteToBase64(v), nil
	}

	return "", nil
}

// Base64Decode decodes a base64 encoded string.
func Base64Decode(content string) (string, error) {
	if content == "" {
		return "", ErrEmptyContent
	}

	payload, err := base64ToByte(content)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	return strings.TrimSpace(string(payload)), nil
}

// Base64GzipEncode compresses and base64 encodes content.
func Base64GzipEncode(content string) (string, error) {
	if content == "" {
		return "", nil
	}

	var buf bytes.Buffer

	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return "", fmt.Errorf("failed to create gzip writer: %w", err)
	}

	_, err = gz.Write([]byte(content))
	if err != nil {
		gz.Close() //revive:disable-line:unhandled-error

		return "", fmt.Errorf("failed to write to gzip: %w", err)
	}

	err = gz.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close gzip writer: %w", err)
	}

	return byteToBase64(buf.Bytes()), nil
}

// Base64GzipDecode decodes and decompresses base64 gzipped content.
func Base64GzipDecode(content string) (string, error) {
	if content == "" {
		return "", nil
	}

	data, err := base64ToByte(content)
	if err != nil {
		return "", fmt.Errorf("base64 decode error: %w", err)
	}

	rd, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer rd.Close()

	dataDecomp, err := io.ReadAll(rd)
	if err != nil {
		return "", fmt.Errorf("failed to read gzip data: %w", err)
	}

	return string(dataDecomp), nil
}

// base64ToByte decodes a base64 string into bytes.
func base64ToByte(content string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, err
	}

	return decoded, nil
}

// byteToBase64 encodes bytes into a base64 string.
func byteToBase64(content []byte) string {
	return base64.StdEncoding.EncodeToString(content)
}
