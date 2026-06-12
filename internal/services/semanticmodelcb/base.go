// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package semanticmodelcb

import (
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/tftypeinfo"
)

var ItemTypeInfo = tftypeinfo.TFTypeInfo{ //nolint:gochecknoglobals
	Name:           "Semantic Model Connection Binding",
	Type:           "semantic_model_connection_binding",
	Names:          "Semantic Model Connection Bindings",
	Types:          "semantic_model_connection_bindings",
	DocsURL:        "https://learn.microsoft.com/rest/api/fabric/semanticmodel/items/bind-semantic-model-connection",
	IsPreview:      true,
	IsSPNSupported: true,
}
