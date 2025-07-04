// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"bytes"
	"os"
	"path/filepath"
	"text/template"
	"unicode/utf8"

	"github.com/go-sprout/sprout"
	sproutchecksum "github.com/go-sprout/sprout/registry/checksum"
	sproutconversion "github.com/go-sprout/sprout/registry/conversion"
	sproutencoding "github.com/go-sprout/sprout/registry/encoding"
	sproutmaps "github.com/go-sprout/sprout/registry/maps"
	sproutnumeric "github.com/go-sprout/sprout/registry/numeric"
	sproutrandom "github.com/go-sprout/sprout/registry/random"
	sproutregexp "github.com/go-sprout/sprout/registry/regexp"
	sproutsemver "github.com/go-sprout/sprout/registry/semver"
	sproutslices "github.com/go-sprout/sprout/registry/slices"
	sproutstd "github.com/go-sprout/sprout/registry/std"
	sproutstrings "github.com/go-sprout/sprout/registry/strings"
	sprouttime "github.com/go-sprout/sprout/registry/time"
	sproutuniqueid "github.com/go-sprout/sprout/registry/uniqueid"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

// getTmplFuncs initializes and returns template functions from the sprout library.
func getTmplFuncs() (template.FuncMap, error) {
	handler := sprout.New()

	err := handler.AddRegistries(
		sproutchecksum.NewRegistry(),
		sproutconversion.NewRegistry(),
		sproutencoding.NewRegistry(),
		sproutmaps.NewRegistry(),
		sproutnumeric.NewRegistry(),
		sproutrandom.NewRegistry(),
		sproutregexp.NewRegistry(),
		sproutsemver.NewRegistry(),
		sproutslices.NewRegistry(),
		sproutstd.NewRegistry(),
		sproutstrings.NewRegistry(),
		sprouttime.NewRegistry(),
		sproutuniqueid.NewRegistry(),

		// unsupported due to security concerns or not real use cases
		// sproutcrypto.NewRegistry(),
		// sproutenv.NewRegistry(),
		// sproutfilesystem.NewRegistry(),
		// sproutreflect.NewRegistry(),
	)
	if err != nil {
		return nil, err
	}

	return handler.Build(), nil
}

// SourceFileToPayload transforms a source file into a base64 encoded payload and calculates its SHA256 hash.
// It optionally processes the file as a template with the provided tokens.
func SourceFileToPayload(srcPath string, tokens map[string]string) (string, string, diag.Diagnostics) { //revive:disable-line:confusing-results
	var diags diag.Diagnostics

	content, err := os.ReadFile(srcPath)
	if err != nil {
		diags.AddError(common.ErrorFileReadHeader, err.Error())

		return "", "", diags
	}

	var contentB64, contentSha256 string

	if utf8.Valid(content) { //nolint:nestif
		tmplFuncs, err := getTmplFuncs()
		if err != nil {
			diags.AddError("Template functions error", err.Error())

			return "", "", diags
		}

		tmpl, err := template.New("tmpl").Funcs(tmplFuncs).ParseFiles(srcPath)
		if err != nil {
			diags.AddError(common.ErrorFileReadHeader, err.Error())

			return "", "", diags
		}

		// Process template with tokens if provided
		tokensData := map[string]string{}
		if len(tokens) > 0 {
			tokensData = tokens
		}

		// Execute template
		var contentBuf bytes.Buffer

		err = tmpl.ExecuteTemplate(&contentBuf, filepath.Base(srcPath), tokensData)
		if err != nil {
			diags.AddError(common.ErrorTmplParseHeader, err.Error())

			return "", "", diags
		}

		contentStr := contentBuf.String()

		// If content is JSON, normalize it
		if IsJSON(contentStr) {
			normalizedContent, err := JSONNormalize(contentStr)
			if err != nil {
				diags.AddError(common.ErrorJSONNormalizeHeader, err.Error())

				return "", "", diags
			}

			contentStr = normalizedContent
		}

		contentSha256 = utils.Sha256(contentStr)

		contentB64, err = Base64Encode(contentStr)
		if err != nil {
			diags.AddError(common.ErrorBase64EncodeHeader, err.Error())

			return "", "", diags
		}
	} else {
		// Handle binary file
		contentSha256 = utils.Sha256(content)

		contentB64, err = Base64Encode(content)
		if err != nil {
			diags.AddError(common.ErrorBase64EncodeHeader, err.Error())

			return "", "", diags
		}
	}

	return contentB64, contentSha256, nil
}

// PayloadToGzip transforms a base64 encoded content string to a gzip compressed base64 string.
// If the content is valid JSON, it uses JSON-specific encoding.
func PayloadToGzip(content string) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	if content == "" {
		return "", diags
	}

	// Decode the base64 content first
	decoded, err := Base64Decode(content)
	if err != nil {
		diags.AddError(common.ErrorBase64DecodeHeader, err.Error())

		return "", diags
	}

	// Re-encode with compression based on content type
	var encoded string
	if IsJSON(decoded) {
		encoded, err = JSONBase64GzipEncode(decoded)
		if err != nil {
			diags.AddError(common.ErrorBase64GzipEncodeHeader, err.Error())

			return "", diags
		}
	} else {
		encoded, err = Base64GzipEncode(decoded)
		if err != nil {
			diags.AddError(common.ErrorBase64GzipEncodeHeader, err.Error())

			return "", diags
		}
	}

	return encoded, diags
}
