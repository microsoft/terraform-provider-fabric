// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0
package shortcut

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	schemaD "github.com/hashicorp/terraform-plugin-framework/datasource/schema" //revive:disable-line:import-alias-naming
	schemaR "github.com/hashicorp/terraform-plugin-framework/resource/schema"   //revive:disable-line:import-alias-naming
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	superschema "github.com/orange-cloudavenue/terraform-plugin-framework-superschema"

	"github.com/microsoft/terraform-provider-fabric/internal/framework/customtypes"
	"github.com/microsoft/terraform-provider-fabric/internal/pkg/fabricitem"
)

func itemSchema(isList bool) superschema.Schema { //revive:disable-line:flag-parameter
	var dsTimeout *superschema.DatasourceTimeoutAttribute

	if !isList {
		dsTimeout = &superschema.DatasourceTimeoutAttribute{
			Read: true,
		}
	}

	return superschema.Schema{
		Resource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewResourceMarkdownDescription(ItemTypeInfo, false),
		},
		DataSource: superschema.SchemaDetails{
			MarkdownDescription: fabricitem.NewDataSourceMarkdownDescription(ItemTypeInfo, isList),
		},
		Attributes: map[string]superschema.Attribute{
			"id": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The " + ItemTypeInfo.Name + " ID.",
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The Workspace ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Item ID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: true,
				},
			},
			"path": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: `A string representing the full path where the shortcut is created, including either "Files" or "Tables".`,
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
					Validators: []validator.String{
						stringvalidator.LengthAtMost(256),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[^/].*$`),
							"Shortcut path can't start with forward slash '/'.",
						),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
			},
			"name": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Name of the shortcut.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.RequiresReplace(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Required: !isList,
					Computed: isList,
				},
			},
			"target": superschema.SuperSingleNestedAttributeOf[targetModel]{
				Common: &schemaR.SingleNestedAttribute{
					MarkdownDescription: "An object that contains the target datasource, and it must specify exactly one of the supported destinations: OneLake, Amazon S3, ADLS Gen2, Google Cloud Storage, S3 compatible or Dataverse.",
				},
				Resource: &schemaR.SingleNestedAttribute{
					Required: true,
				},
				DataSource: &schemaD.SingleNestedAttribute{
					Computed: true,
				},
				Attributes: map[string]superschema.Attribute{
					"type": superschema.StringAttribute{
						Common: &schemaR.StringAttribute{
							MarkdownDescription: "The type object contains properties like target shortcut account type. Additional types may be added over time.",
						},
						Resource: &schemaR.StringAttribute{
							Computed: true,
						},
						DataSource: &schemaD.StringAttribute{
							Computed: true,
						},
					},
					"onelake":              onelakeSchema(),
					"adls_gen2":            adlsGen2Schema(),
					"amazon_s3":            amazonS3Schema(),
					"azure_blob_storage":   azureBlobStorageSchema(),
					"google_cloud_storage": googleCloudStorageSchema(),
					"s3_compatible":        s3CompatibleSchema(),
					"external_data_share":  externalDataShareSchema(),
					"dataverse":            dataverseSchema(),
				},
			},

			"timeouts": superschema.TimeoutAttribute{
				Resource: &superschema.ResourceTimeoutAttribute{
					Create: true,
					Read:   true,
					Delete: true,
					Update: true,
				},
				DataSource: dsTimeout,
			},
		},
	}
}

func onelakeSchema() superschema.SuperSingleNestedAttributeOf[oneLakeModel] {
	return superschema.SuperSingleNestedAttributeOf[oneLakeModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target OneLake data source.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"item_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the target in OneLake. The target can be an item of Lakehouse, KQLDatabase, or Warehouse.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"path": superschema.StringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the full path to the target folder within the Item. This path should be relative to the root of the OneLake directory structure. For example: 'Tables/myTablesFolder/someTableSubFolder'.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
					Validators: []validator.String{
						stringvalidator.LengthAtMost(256),
						stringvalidator.RegexMatches(
							regexp.MustCompile(`^[^/].*$`),
							"OneLake path can't start with forward slash '/'.",
						),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"workspace_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "The ID of the target workspace.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func adlsGen2Schema() superschema.SuperSingleNestedAttributeOf[targetDataSourceModel] {
	return superschema.SuperSingleNestedAttributeOf[targetDataSourceModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target ADLS Gen2 data source.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the connection that is bound with the shortcut. The connectionId is a unique identifier used to establish a connection between the shortcut and the target datasource. To find this connection ID, first create a cloud connection to be used by the shortcut when connecting to the ADLS data location. Open the cloud connection's Settings view and copy the connection ID; this is a GUID.",
					CustomType:          customtypes.UUIDType{},
				},

				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"location": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies the location of the target ADLS container. The URI must be in the format https://[account-name].dfs.core.windows.net where [account-name] is the name of the target ADLS account.",
					CustomType:          customtypes.URLType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"subpath": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies the container and subfolder within the ADLS account where the target folder is located. Must be of the format [container]/[subfolder] where [container] is the name of the container that holds the files and folders; [subfolder] is the name of the subfolder within the container (optional). For example: /mycontainer/mysubfolder",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func amazonS3Schema() superschema.SuperSingleNestedAttributeOf[targetDataSourceModel] {
	return superschema.SuperSingleNestedAttributeOf[targetDataSourceModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target Amazon S3 data source.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the connection that is bound with the shortcut. The connectionId is a unique identifier used to establish a connection between the shortcut and the target datasource. To find this connection ID, first create a cloud connection to be used by the shortcut when connecting to the Amazon S3 data location. Open the cloud connection's Settings view and copy the connection ID; this is a GUID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"location": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "HTTP URL that points to the target bucket in S3. The URL should be in the format https://[bucket-name].s3.[region-code].amazonaws.com, where 'bucket-name' is the name of the S3 bucket you want to point to, and 'region-code' is the code for the region where the bucket is located. For example: https://my-s3-bucket.s3.us-west-2.amazonaws.com",
					CustomType:          customtypes.URLType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"subpath": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies a target folder or subfolder within the S3 bucket.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func azureBlobStorageSchema() superschema.SuperSingleNestedAttributeOf[targetDataSourceModel] {
	return superschema.SuperSingleNestedAttributeOf[targetDataSourceModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target Google Cloud Storage data source.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the connection that is bound with the shortcut. The connectionId is a unique identifier used to establish a connection between the shortcut and the target datasource.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"location": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "HTTP URL that points to the target bucket in GCS. The URL should be in the format https://[bucket-name].storage.googleapis.com, where [bucket-name] is the name of the bucket you want to point to. For example: https://my-gcs-bucket.storage.googleapis.com",
					CustomType:          customtypes.URLType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"subpath": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies a target folder or subfolder within the GCS bucket. For example: /folder",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func googleCloudStorageSchema() superschema.SuperSingleNestedAttributeOf[targetDataSourceModel] {
	return superschema.SuperSingleNestedAttributeOf[targetDataSourceModel]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target Google Cloud Storage data source.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the connection that is bound with the shortcut. The connectionId is a unique identifier used to establish a connection between the shortcut and the target datasource.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"location": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "HTTP URL that points to the target bucket in GCS. The URL should be in the format https://[bucket-name].storage.googleapis.com, where [bucket-name] is the name of the bucket you want to point to. For example: https://my-gcs-bucket.storage.googleapis.com",
					CustomType:          customtypes.URLType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"subpath": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies a target folder or subfolder within the GCS bucket. For example: /folder",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func s3CompatibleSchema() superschema.SuperSingleNestedAttributeOf[s3Compatible] {
	return superschema.SuperSingleNestedAttributeOf[s3Compatible]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target S3 compatible data source.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the connection that is bound with the shortcut. The connectionId is a unique identifier used to establish a connection between the shortcut and the target datasource.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"location": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "HTTP URL of the S3 compatible endpoint. This endpoint must be able to receive ListBuckets S3 API calls. The URL must be in the non-bucket specific format; no bucket should be specified here. For example: https://s3endpoint.contoso.com",
					CustomType:          customtypes.URLType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"subpath": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies a target folder or subfolder within the S3 compatible bucket. For example: /folder",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"bucket": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies the target bucket within the S3 compatible location.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func dataverseSchema() superschema.SuperSingleNestedAttributeOf[dataverse] {
	return superschema.SuperSingleNestedAttributeOf[dataverse]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target Dataverse data source.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Optional: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the connection that is bound with the shortcut. The connectionId is a unique identifier used to establish a connection between the shortcut and the target datasource. To find this connection ID, first create a cloud connection to be used by the shortcut when connecting to the Dataverse data location. Open the cloud connection's Settings view and copy the connection ID; this is a GUID.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"environment_domain": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "URI that indicates the Dataverse target environment's domain name. The URI should be formatted as 'https://[orgname].crm[xx].dynamics.com', where [orgname] represents the name of your Dataverse organization.",
					CustomType:          customtypes.URLType{},
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"table_name": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies the name of the target table in Dataverse",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
			"deltalake_folder": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "Specifies the DeltaLake folder path where the target data is stored.",
				},
				Resource: &schemaR.StringAttribute{
					Required: true,
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}

func externalDataShareSchema() superschema.SuperSingleNestedAttributeOf[externalDataShare] {
	return superschema.SuperSingleNestedAttributeOf[externalDataShare]{
		Common: &schemaR.SingleNestedAttribute{
			MarkdownDescription: "An object containing the properties of the target external data share.",
		},
		Resource: &schemaR.SingleNestedAttribute{
			Computed: true,
		},
		DataSource: &schemaD.SingleNestedAttribute{
			Computed: true,
		},
		Attributes: map[string]superschema.Attribute{
			"connection_id": superschema.SuperStringAttribute{
				Common: &schemaR.StringAttribute{
					MarkdownDescription: "A string representing the connection that is bound with the shortcut. The connectionId is a unique identifier used to establish a connection between the shortcut and the target datasource.",
					CustomType:          customtypes.UUIDType{},
				},
				Resource: &schemaR.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				DataSource: &schemaD.StringAttribute{
					Computed: true,
				},
			},
		},
	}
}
