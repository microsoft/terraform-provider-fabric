// Copyright Microsoft Corporation 2026
// SPDX-License-Identifier: MPL-2.0

package monitoring

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/mapvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabmonitoring "github.com/microsoft/fabric-sdk-go/fabric/monitoring"
	supertypes "github.com/orange-cloudavenue/terraform-plugin-framework-supertypes"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

type resourceMonitoring struct {
	fabricitem.ResourceFabricItemDefinitionProperties[monitoringPropertiesModel, fabmonitoring.Properties]
}

func NewResourceMonitoring() resource.Resource {
	propertiesSetter := func(ctx context.Context, from *fabmonitoring.Properties, to *fabricitem.ResourceFabricItemDefinitionPropertiesModel[monitoringPropertiesModel, fabmonitoring.Properties]) diag.Diagnostics {
		properties := supertypes.NewSingleNestedObjectValueOfNull[monitoringPropertiesModel](ctx)

		if from != nil {
			propertiesModel := &monitoringPropertiesModel{}

			if diags := propertiesModel.set(ctx, *from); diags.HasError() {
				return diags
			}

			if diags := properties.Set(ctx, propertiesModel); diags.HasError() {
				return diags
			}
		}

		to.Properties = properties

		return nil
	}

	itemGetter := func(ctx context.Context, fabricClient fabric.Client, model fabricitem.ResourceFabricItemDefinitionPropertiesModel[monitoringPropertiesModel, fabmonitoring.Properties], fabricItem *fabricitem.FabricItemProperties[fabmonitoring.Properties]) error {
		client := fabmonitoring.NewClientFactoryWithClient(fabricClient).NewItemsClient()

		respGet, err := client.GetMonitoring(ctx, model.WorkspaceID.ValueString(), model.ID.ValueString(), nil)
		if err != nil {
			return err
		}

		fabricItem.Set(respGet.Monitoring)

		return nil
	}

	config := fabricitem.ResourceFabricItemDefinitionProperties[monitoringPropertiesModel, fabmonitoring.Properties]{
		TypeInfo:              ItemTypeInfo,
		FabricItemType:        FabricItemType,
		NameRenameAllowed:     true,
		DisplayNameMaxLength:  128,
		DescriptionMaxLength:  256,
		DefinitionPathDocsURL: ItemDefinitionPathDocsURL,
		DefinitionPathKeysValidator: []validator.Map{
			mapvalidator.SizeAtMost(len(itemDefinitionFormats)),
			mapvalidator.KeysAre(fabricitem.DefinitionPathKeysValidator(itemDefinitionFormats)...),
		},
		DefinitionRequired:   false,
		DefinitionEmpty:      ItemDefinitionEmpty,
		DefinitionFormats:    itemDefinitionFormats,
		PropertiesAttributes: getResourceMonitoringPropertiesAttributes(),
		PropertiesSetter:     propertiesSetter,
		ItemGetter:           itemGetter,
	}

	return &resourceMonitoring{ResourceFabricItemDefinitionProperties: config}
}

func (r *resourceMonitoring) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	r.ResourceFabricItemDefinitionProperties.Schema(ctx, req, resp)
	resp.Schema.Attributes["display_name"] = schema.StringAttribute{
		MarkdownDescription: "The Monitoring display name: 2 to 128 ASCII letters, digits, underscores, periods or hyphens. Spaces and leading or trailing double underscores are not allowed.",
		Required:            true,
		Validators: []validator.String{
			stringvalidator.LengthBetween(2, 128),
			stringvalidator.RegexMatches(regexp.MustCompile(`^[A-Za-z0-9_.-]+$`), "Use only ASCII letters, digits, underscores, periods and hyphens."),
			stringvalidator.RegexMatches(regexp.MustCompile(`^[^_]|^_[^_]`), "The name must not start with a double underscore."),
			stringvalidator.RegexMatches(regexp.MustCompile(`[^_]$|[^_]_$`), "The name must not end with a double underscore."),
		},
	}
	resp.Schema.MarkdownDescription += "\n\nDefining routing does not enable collection. Use `fabric_monitoring_data_collection` to manage collection separately. Once configured, the Eventhouse destination mode and Linked parent reference cannot be changed or removed; an enabled Operations Agent cannot be disabled."
}
