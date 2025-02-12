// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package transforms

import (
	"bytes"
	"context"
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
	"github.com/hashicorp/terraform-plugin-framework/types"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/common"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/utils"
)

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

func SourceFileToPayload(ctx context.Context, srcPath types.String, tokens supertypes.MapValueOf[string]) (*string, *string, diag.Diagnostics) {
	var diags diag.Diagnostics

	source := srcPath.ValueString()

	content, err := os.ReadFile(srcPath.ValueString())
	if err != nil {
		diags.AddError(
			common.ErrorFileReadHeader,
			err.Error(),
		)

		return nil, nil, diags
	}

	var contentSha256 string
	var contentB64 string

	if utf8.Valid(content) { //nolint:nestif
		tmplFuncs, err := getTmplFuncs()
		if err != nil {
			diags.AddError(
				"Template functions error",
				err.Error(),
			)

			return nil, nil, diags
		}

		tmpl, err := template.New("tmpl").Funcs(tmplFuncs).ParseFiles(source)
		if err != nil {
			diags.AddError(
				common.ErrorFileReadHeader,
				err.Error(),
			)

			return nil, nil, diags
		}

		tokensData := map[string]string{}

		if !tokens.IsNull() && !tokens.IsUnknown() {
			tokensData, diags = tokens.Get(ctx)
			if diags.HasError() {
				return nil, nil, diags
			}
		}

		var contentBuf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&contentBuf, filepath.Base(source), tokensData); err != nil {
			diags.AddError(
				common.ErrorTmplParseHeader,
				err.Error(),
			)

			return nil, nil, diags
		}

		content := contentBuf.String()

		if IsJSON(content) {
			if err := JSONNormalize(&content); err != nil {
				diags.AddError(
					common.ErrorJSONNormalizeHeader,
					err.Error(),
				)

				return nil, nil, diags
			}
		}

		contentSha256 = utils.Sha256(content)

		contentB64, err = Base64Encode(content)
		if err != nil {
			diags.AddError(
				common.ErrorBase64EncodeHeader,
				err.Error(),
			)

			return nil, nil, diags
		}
	} else {
		contentSha256 = utils.Sha256(content)

		contentB64, err = Base64Encode(content)
		if err != nil {
			diags.AddError(
				common.ErrorBase64EncodeHeader,
				err.Error(),
			)

			return nil, nil, diags
		}
	}

	return &contentB64, &contentSha256, nil
}

func PayloadToGzip(content *string) diag.Diagnostics {
	var diags diag.Diagnostics

	if err := Base64Decode(content); err != nil {
		diags.AddError(
			common.ErrorBase64DecodeHeader,
			err.Error(),
		)

		return diags
	}

	if IsJSON(*content) {
		if err := JSONBase64GzipEncode(content); err != nil {
			diags.AddError(
				common.ErrorBase64GzipEncodeHeader,
				err.Error(),
			)

			return diags
		}
	} else {
		if err := Base64GzipEncode(content); err != nil {
			diags.AddError(
				common.ErrorBase64GzipEncodeHeader,
				err.Error(),
			)

			return diags
		}
	}

	return nil
}
