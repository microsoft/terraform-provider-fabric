// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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
	"github.com/ohler55/ojg/jp"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/params"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

const (
	ParameterTypeTextReplace     string = "TextReplace"
	ParameterTypeJSONPathReplace string = "JsonPathReplace"
)

func PossibleParameterTypeValues() []string {
	return []string{
		ParameterTypeTextReplace,
		ParameterTypeJSONPathReplace,
	}
}

const (
	ProcessingModeGoTemplate string = "GoTemplate"
	ProcessingModeParameters string = "Parameters"
	ProcessingModeNone       string = "None"
)

func PossibleProcessingModeValues() []string {
	return []string{
		ProcessingModeGoTemplate,
		ProcessingModeParameters,
		ProcessingModeNone,
	}
}

const (
	TokensDelimiterCurlyBraces string = "{{}}"
	TokensDelimiterAngles      string = "<<>>"
	TokensDelimiterAt          string = "@{}@"
	TokensDelimiterUnderscore  string = "____"
)

func PossibleTokensDelimiterValues() []string {
	return []string{
		TokensDelimiterCurlyBraces,
		TokensDelimiterAngles,
		TokensDelimiterAt,
		TokensDelimiterUnderscore,
	}
}

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
//
//nolint:gocognit
func SourceFileToPayload(
	srcPath string,
	processingMode string,
	tokens map[string]string,
	parameters []*params.ParametersModel,
	tokensDelimiter string,
) (string, string, diag.Diagnostics) { //revive:disable-line:confusing-results
	var diags diag.Diagnostics

	content, err := os.ReadFile(srcPath)
	if err != nil {
		diags.AddError(common.ErrorFileReadHeader, err.Error())

		return "", "", diags
	}

	var contentB64, contentSha256, contentStr string

	if utf8.Valid(content) { //nolint:nestif
		switch strings.ToLower(processingMode) {
		case strings.ToLower(ProcessingModeGoTemplate):
			tmplFuncs, err := getTmplFuncs()
			if err != nil {
				diags.AddError("Template functions error", err.Error())

				return "", "", diags
			}

			var tmpl *template.Template

			if tokensDelimiter == TokensDelimiterCurlyBraces {
				tmpl, err = template.New("tmpl").Funcs(tmplFuncs).ParseFiles(srcPath)
			} else {
				leftDelim := tokensDelimiter[:2]
				rightDelim := tokensDelimiter[len(tokensDelimiter)-2:]
				tmpl, err = template.New("tmpl").Delims(leftDelim, rightDelim).Funcs(tmplFuncs).ParseFiles(srcPath)
			}

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

			contentStr = contentBuf.String()
		case strings.ToLower(ProcessingModeParameters):
			contentStr = string(content)

			for _, param := range parameters {
				switch strings.ToLower(param.Type.ValueString()) {
				case strings.ToLower(ParameterTypeTextReplace):
					contentStr = strings.ReplaceAll(contentStr, param.Find.ValueString(), param.Value.ValueString())
				case strings.ToLower(ParameterTypeJSONPathReplace):
					if IsJSON(contentStr) {
						contentStr, diags = processJSONPathReplacement(contentStr, param)
						if diags.HasError() {
							return "", "", diags
						}
					}
				default:
					diags.AddError("Unsupported parameter type", "Invalid parameter type: "+param.Type.ValueString())

					return "", "", diags
				}
			}
		default:
			contentStr = string(content)
		}

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

// processJSONPathReplacement handles JSON path replacement for a parameter.
func processJSONPathReplacement(contentStr string, param *params.ParametersModel) (string, diag.Diagnostics) {
	var diags diag.Diagnostics

	jpExpression, err := jp.ParseString(param.Find.ValueString())
	if err != nil {
		diags.AddError("JSONPath expression", err.Error())

		return "", diags
	}

	var contentJSON any

	err = json.Unmarshal([]byte(contentStr), &contentJSON)
	if err != nil {
		diags.AddError("JSON unmarshal", err.Error())

		return "", diags
	}

	jpIter := jpExpression.Get(contentJSON)

	if len(jpIter) > 0 {
		err := jpExpression.Set(contentJSON, param.Value.ValueString())
		if err != nil {
			diags.AddError("JSONPath set", err.Error())

			return "", diags
		}

		content, err := json.Marshal(contentJSON)
		if err != nil {
			diags.AddError("JSON marshal", err.Error())

			return "", diags
		}

		return string(content), diags
	}

	return contentStr, diags
}
