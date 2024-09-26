// Copyright (c) Microsoft Corporation
// SPDX-License-Identifier: MPL-2.0

package testhelp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	supertypes "github.com/FrangipaneTeam/terraform-plugin-framework-supertypes"
	"github.com/hashicorp/go-azure-sdk/sdk/auth"
	"github.com/hashicorp/go-azure-sdk/sdk/environments"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/manicminer/hamilton/msgraph"
	"github.com/microsoft/fabric-sdk-go/fabric"
	fabadmin "github.com/microsoft/fabric-sdk-go/fabric/admin"
	fabcore "github.com/microsoft/fabric-sdk-go/fabric/core"

	"github.com/microsoft/terraform-provider-fabric/internal/pkg/transforms"
)

type wellKnownItemModel struct {
	ID          *string `json:"id"`
	DisplayName *string `json:"displayName"`
	Description *string `json:"description"`
}

type wellKnownPrincipalModel struct {
	ID   *string `json:"id"`
	Type *string `json:"type"`
}

type wellKnown struct {
	Workspace          wellKnownItemModel      `json:"workspace"`
	Dashboard          wellKnownItemModel      `json:"dashboard"`
	Datamart           wellKnownItemModel      `json:"datamart"`
	DataPipeline       wellKnownItemModel      `json:"dataPipeline"`
	Environment        wellKnownItemModel      `json:"environment"`
	Eventhouse         wellKnownItemModel      `json:"eventhouse"`
	Eventstream        wellKnownItemModel      `json:"eventstream"`
	KQLDatabase        wellKnownItemModel      `json:"kqlDatabase"`
	KQLQueryset        wellKnownItemModel      `json:"kqlQueryset"`
	Lakehouse          wellKnownItemModel      `json:"lakehouse"`
	MLExperiment       wellKnownItemModel      `json:"mlExperiment"`
	MLModel            wellKnownItemModel      `json:"mlModel"`
	Notebook           wellKnownItemModel      `json:"notebook"`
	Report             wellKnownItemModel      `json:"report"`
	SemanticModel      wellKnownItemModel      `json:"semanticModel"`
	SparkJobDefinition wellKnownItemModel      `json:"sparkJobDefinition"`
	Warehouse          wellKnownItemModel      `json:"warehouse"`
	Capacity           wellKnownItemModel      `json:"capacity"`
	DomainParent       wellKnownItemModel      `json:"domainParent"`
	DomainChild        wellKnownItemModel      `json:"domainChild"`
	Principal          wellKnownPrincipalModel `json:"principal"`
	Group              wellKnownPrincipalModel `json:"group"`
}

const (
	wellKnownEnvKey                      = "FABRIC_TESTACC_WELLKNOWN"
	wellKnownShouldCreateResourcesEnvKey = "FABRIC_TESTACC_WELLKNOWN_CREATE_RESOURCES"
	wellKnownCapacityNameEnvKey          = "FABRIC_TESTACC_WELLKNOWN_CAPACITY_NAME"
)

var (
	wellKnownFilePath = getFixtureFilePath(".wellknown.json")
	wellKnownData     *wellKnown
)

// WellKnownData returns the well-known IDs for the Fabric resources.
// If the FABRIC_TESTACC_WELLKNOWN environment variable is set, it will try to parse the values from a json string.
// If the parsing fails, it will return an error.
// If the environment variable is not set, it will try to read it from the file "fixtures/.wellknown.json".
// If no data is found, and the environment variable FABRIC_TESTACC_WELLKNOWN_CREATE_RESOURCES is set to "true", it will create the resources and write the values to the file.
// When creating the resources, it will use the first capacity it can find.
func WellKnown() wellKnown { //revive:disable-line:unexported-return
	if wellKnownData != nil {
		return *wellKnownData
	}

	if !IsWellKnownDataAvailable() {
		panicOnError(fmt.Errorf("well-known resources file %s does not exist", wellKnownFilePath))
	}

	var wk wellKnown

	if wellKnownJSON, ok := os.LookupEnv(wellKnownEnvKey); ok {
		err := json.Unmarshal([]byte(wellKnownJSON), &wk)
		panicOnError(err)

		return wk
	}

	// read the file into a string
	wellKnownJSONBytes, err := os.ReadFile(wellKnownFilePath)
	panicOnError(err)

	// parse the json string
	err = json.Unmarshal(wellKnownJSONBytes, &wk)
	panicOnError(err)

	wellKnownData = &wk

	return wk
}

func IsWellKnownDataAvailable() bool {
	wellKnownJSON := os.Getenv(wellKnownEnvKey)

	// if the environment variable is set, we don't need to check the file
	if wellKnownJSON != "" {
		return true
	}

	// check if the file exists
	_, err := os.Stat(wellKnownFilePath)

	return !os.IsNotExist(err)
}

func ShouldCreateWellKnownResources() bool {
	shouldCreateResources := os.Getenv(wellKnownShouldCreateResourcesEnvKey)

	return strings.EqualFold(shouldCreateResources, "true")
}

// CreateWellKnownResources creates the well-known resources and writes the values to the file "fixtures/.wellknown.json".
// It will use the first capacity it can find.
// It will use those values for the principal.
// IMPORTANT: this function is still a work-in-progress and the SemanticModel and Report resource creation is not yet implemented.
func CreateWellKnownResources() { //nolint:maintidx
	ctx := context.Background()
	values := wellKnown{}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	panicOnError(err)

	fabricClientOpts := &policy.ClientOptions{}
	fabricClientOpts.Retry.MaxRetryDelay = -1
	fabClient, err := fabric.NewClient(cred, nil, fabricClientOpts)
	panicOnError(err)

	coreCF := fabcore.NewClientFactoryWithClient(*fabClient)
	adminCF := fabadmin.NewClientFactoryWithClient(*fabClient)

	capacities, err := coreCF.NewCapacitiesClient().ListCapacities(ctx, nil)
	panicOnError(err)

	if len(capacities) == 0 {
		panicOnError(fabcore.ErrCapacity.CapacityNotFound)
	}

	// try to get the capacity by name if environment variable is set
	if capacityName, ok := os.LookupEnv(wellKnownCapacityNameEnvKey); ok {
		for _, capacity := range capacities {
			if *capacity.DisplayName == capacityName {
				values.Capacity.ID = capacity.ID
				values.Capacity.DisplayName = capacity.DisplayName
				values.Capacity.Description = capacity.SKU

				break
			}
		}
	}

	// if the capacity is still not set, use the first one
	if values.Capacity.ID == nil {
		values.Capacity.ID = capacities[0].ID
		values.Capacity.DisplayName = capacities[0].DisplayName
		values.Capacity.Description = capacities[0].SKU
	}

	log.Printf("Using 'Capacity' (DisplayName: %s, ID: %s)\n", *values.Capacity.DisplayName, *values.Capacity.ID)

	// Create parent  domain
	domainClient := adminCF.NewDomainsClient()

	domainResp, err := domainClient.CreateDomain(ctx, fabadmin.CreateDomainRequest{
		DisplayName: to.Ptr(RandomName()),
		Description: to.Ptr(RandomName()),
	}, nil)
	panicOnError(err)
	log.Printf("Created parent 'Domain' (DisplayName: %s, ID: %s)\n", *domainResp.DisplayName, *domainResp.ID)

	values.DomainParent.ID = domainResp.ID
	values.DomainParent.DisplayName = domainResp.DisplayName
	values.DomainParent.Description = domainResp.Description

	// Create child domain
	domainResp, err = domainClient.CreateDomain(ctx, fabadmin.CreateDomainRequest{
		DisplayName:    to.Ptr(RandomName()),
		Description:    to.Ptr(RandomName()),
		ParentDomainID: domainResp.ID,
	}, nil)
	panicOnError(err)
	log.Printf("Created child 'Domain' (DisplayName: %s, ID: %s)\n", *domainResp.DisplayName, *domainResp.ID)

	values.DomainChild.ID = domainResp.ID
	values.DomainChild.DisplayName = domainResp.DisplayName
	values.DomainChild.Description = domainResp.Description

	// Create a workspace
	workspaceResp, err := coreCF.NewWorkspacesClient().CreateWorkspace(ctx, fabcore.CreateWorkspaceRequest{
		DisplayName: to.Ptr(RandomName()),
		Description: to.Ptr(RandomName()),
		CapacityID:  values.Capacity.ID,
	}, nil)
	panicOnError(err)
	log.Printf("Created 'Workspace' (DisplayName: %s, ID: %s)\n", *workspaceResp.DisplayName, *workspaceResp.ID)

	values.Workspace.ID = workspaceResp.ID
	values.Workspace.DisplayName = workspaceResp.DisplayName
	values.Workspace.Description = workspaceResp.Description

	itemsClient := coreCF.NewItemsClient()

	itemTypes := []string{"DataPipeline", "Environment", "Eventhouse", "Eventstream", "Lakehouse", "MLExperiment", "MLModel", "Notebook", "SparkJobDefinition", "Warehouse"}

	// Create items
	for _, itemType := range itemTypes {
		itemResp, err := itemsClient.CreateItem(ctx, *values.Workspace.ID, fabcore.CreateItemRequest{
			DisplayName: to.Ptr(RandomName()),
			Description: to.Ptr(RandomName()),
			Type:        to.Ptr(fabcore.ItemType(itemType)),
		}, nil)
		panicOnError(err)
		log.Printf("Created '%s' (DisplayName: %s, ID: %s)\n", itemType, *itemResp.DisplayName, *itemResp.ID)

		switch itemType {
		case "DataPipeline":
			values.DataPipeline.ID = itemResp.ID
			values.DataPipeline.DisplayName = itemResp.DisplayName
			values.DataPipeline.Description = itemResp.Description
		case "Environment":
			values.Environment.ID = itemResp.ID
			values.Environment.DisplayName = itemResp.DisplayName
			values.Environment.Description = itemResp.Description
		case "Eventhouse":
			values.Eventhouse.ID = itemResp.ID
			values.Eventhouse.DisplayName = itemResp.DisplayName
			values.Eventhouse.Description = itemResp.Description
		case "Eventstream":
			values.Eventstream.ID = itemResp.ID
			values.Eventstream.DisplayName = itemResp.DisplayName
			values.Eventstream.Description = itemResp.Description
		case "KQLQueryset":
			values.KQLQueryset.ID = itemResp.ID
			values.KQLQueryset.DisplayName = itemResp.DisplayName
			values.KQLQueryset.Description = itemResp.Description
		case "Lakehouse":
			values.Lakehouse.ID = itemResp.ID
			values.Lakehouse.DisplayName = itemResp.DisplayName
			values.Lakehouse.Description = itemResp.Description

			log.Printf("!!! Please go to Lakehouse and manually run 'Start with sample data' to populate the data")
		case "MLExperiment":
			values.MLExperiment.ID = itemResp.ID
			values.MLExperiment.DisplayName = itemResp.DisplayName
			values.MLExperiment.Description = itemResp.Description
		case "MLModel":
			values.MLModel.ID = itemResp.ID
			values.MLModel.DisplayName = itemResp.DisplayName
			values.MLModel.Description = itemResp.Description
		case "Notebook":
			values.Notebook.ID = itemResp.ID
			values.Notebook.DisplayName = itemResp.DisplayName
			values.Notebook.Description = itemResp.Description
		case "SparkJobDefinition":
			values.SparkJobDefinition.ID = itemResp.ID
			values.SparkJobDefinition.DisplayName = itemResp.DisplayName
			values.SparkJobDefinition.Description = itemResp.Description
		case "Warehouse":
			values.Warehouse.ID = itemResp.ID
			values.Warehouse.DisplayName = itemResp.DisplayName
			values.Warehouse.Description = itemResp.Description
		}
	}

	// no Create API
	log.Printf("!!! Please create a KQLQueryset manually, and update details in the well-known file")

	values.KQLQueryset.ID = to.Ptr("00000000-0000-0000-0000-000000000000")
	values.KQLQueryset.DisplayName = to.Ptr("test")
	values.KQLQueryset.Description = to.Ptr("test")

	log.Printf("!!! Please create a Datamart manually, and update details in the well-known file")

	values.Datamart.ID = to.Ptr("00000000-0000-0000-0000-000000000000")
	values.Datamart.DisplayName = to.Ptr("test")
	values.Datamart.Description = to.Ptr("")

	log.Printf("!!! Please create a Dashboard manually, and update details in the well-known file")

	values.Dashboard.ID = to.Ptr("00000000-0000-0000-0000-000000000000")
	values.Dashboard.DisplayName = to.Ptr("test")
	values.Dashboard.Description = to.Ptr("")

	itemTypesDef := map[string][]map[string]string{
		"SemanticModel": {
			{"definition.pbism": "semantic_model_tmsl/definition.pbism"},
			{"model.bim": "semantic_model_tmsl/model.bim.tmpl"},
		},
		"Report": {
			{"definition.pbir": "report_pbir_legacy/definition.pbir.tmpl"},
			{"report.json": "report_pbir_legacy/report.json"},
			{"StaticResources/SharedResources/BaseThemes/CY24SU06.json": "report_pbir_legacy/StaticResources/SharedResources/BaseThemes/CY24SU06.json"},
		},
	}

	tokensSemanticModel := map[string][]map[string]string{
		"model.bim": {
			{"ColumnName": "TestAcc"},
		},
	}
	itemRespSemanticModel := createItemDef(ctx, itemsClient, *values.Workspace.ID, "SemanticModel", itemTypesDef["SemanticModel"], &tokensSemanticModel)
	log.Printf("Created '%s' (DisplayName: %s, ID: %s)\n", "SemanticModel", *itemRespSemanticModel.DisplayName, *itemRespSemanticModel.ID)
	values.SemanticModel.ID = itemRespSemanticModel.ID
	values.SemanticModel.DisplayName = itemRespSemanticModel.DisplayName
	values.SemanticModel.Description = itemRespSemanticModel.Description

	tokensReport := map[string][]map[string]string{
		"definition.pbir": {
			{"SemanticModelID": *values.SemanticModel.ID},
		},
	}
	itemRespReport := createItemDef(ctx, itemsClient, *values.Workspace.ID, "Report", itemTypesDef["Report"], &tokensReport)
	log.Printf("Created '%s' (DisplayName: %s, ID: %s)\n", "Report", *itemRespReport.DisplayName, *itemRespReport.ID)
	values.Report.ID = itemRespReport.ID
	values.Report.DisplayName = itemRespReport.DisplayName
	values.Report.Description = itemRespReport.Description

	creationPayloadKQLDatabase := map[string]any{
		"databaseType":           "ReadWrite",
		"parentEventhouseItemId": *values.Eventhouse.ID,
	}
	itemRespKQLDatabase := createItemCreationPayload(ctx, itemsClient, *values.Workspace.ID, "KQLDatabase", creationPayloadKQLDatabase)
	log.Printf("Created '%s' (DisplayName: %s, ID: %s)\n", "KQLDatabase", *itemRespKQLDatabase.DisplayName, *itemRespKQLDatabase.ID)
	values.KQLDatabase.ID = itemRespKQLDatabase.ID
	values.KQLDatabase.DisplayName = itemRespKQLDatabase.DisplayName
	values.KQLDatabase.Description = itemRespKQLDatabase.Description
	// Init MS Graph client
	env := environments.AzurePublic()
	credentials := auth.Credentials{
		Environment:                       *env,
		EnableAuthenticatingUsingAzureCLI: true,
	}

	authorizer, err := auth.NewAuthorizerFromCredentials(ctx, credentials, env.MicrosoftGraph)
	panicOnError(err)

	appClient := msgraph.NewApplicationsClient()
	appClient.BaseClient.Authorizer = authorizer

	app, _, err := appClient.Create(ctx, msgraph.Application{
		DisplayName: to.Ptr(RandomName()),
	})
	panicOnError(err)
	log.Printf("Created Entra App (DisplayName: %s, AppID: %s, ObjectID: %s)\n", *app.DisplayName, *app.AppId, *app.Id)

	spClient := msgraph.NewServicePrincipalsClient()
	spClient.BaseClient.Authorizer = authorizer

	sp, _, err := spClient.Create(ctx, msgraph.ServicePrincipal{
		AppId: app.AppId,
	})
	panicOnError(err)
	log.Printf("Created Entra Service Principal (DisplayName: %s, AppID: %s, ObjectID: %s)\n", *sp.DisplayName, *sp.AppId, *sp.Id)

	values.Principal.ID = sp.Id
	values.Principal.Type = to.Ptr("ServicePrincipal")

	groupsClient := msgraph.NewGroupsClient()
	groupsClient.BaseClient.Authorizer = authorizer

	groupName := RandomName(10)
	group, _, err := groupsClient.Create(ctx, msgraph.Group{
		DisplayName:     &groupName,
		MailNickname:    &groupName,
		MailEnabled:     to.Ptr(false),
		SecurityEnabled: to.Ptr(true),
	})
	panicOnError(err)
	log.Printf("Created Entra Group (DisplayName: %s, ObjectID: %s)\n", *group.DisplayName, *group.Id)

	values.Group.ID = group.Id
	values.Group.Type = to.Ptr("Group")

	// write the values to the file, pretty printed
	wellKnownJSONBytes, err := json.MarshalIndent(values, "", "	")
	panicOnError(err)

	err = os.WriteFile(wellKnownFilePath, wellKnownJSONBytes, 0o600)
	panicOnError(err)
}

func createItemDef(ctx context.Context, itemsClient *fabcore.ItemsClient, workspaceID, itemType string, itemDefParts []map[string]string, tokens *map[string][]map[string]string) fabcore.ItemsClientCreateItemResponse {
	var defParts []fabcore.ItemDefinitionPart

	for _, itemDefPart := range itemDefParts {
		for defPath, defSource := range itemDefPart {
			defParts = append(defParts, fabcore.ItemDefinitionPart{
				Path:        to.Ptr(defPath),
				Payload:     to.Ptr(sourceToPayload(ctx, defPath, defSource, tokens)),
				PayloadType: to.Ptr(fabcore.PayloadTypeInlineBase64),
			})
		}
	}

	itemResp, err := itemsClient.CreateItem(ctx, workspaceID, fabcore.CreateItemRequest{
		DisplayName: to.Ptr(RandomName()),
		Description: to.Ptr(RandomName()),
		Type:        to.Ptr(fabcore.ItemType(itemType)),
		Definition: &fabcore.ItemDefinition{
			Parts: defParts,
		},
	}, nil)
	panicOnError(err)

	return itemResp
}

func createItemCreationPayload(ctx context.Context, itemsClient *fabcore.ItemsClient, workspaceID, itemType string, creationPayload map[string]any) fabcore.ItemsClientCreateItemResponse {
	itemResp, err := itemsClient.CreateItem(ctx, workspaceID, fabcore.CreateItemRequest{
		DisplayName:     to.Ptr(RandomName()),
		Description:     to.Ptr(RandomName()),
		Type:            to.Ptr(fabcore.ItemType(itemType)),
		CreationPayload: creationPayload,
	}, nil)
	panicOnError(err)

	return itemResp
}

func sourceToPayload(ctx context.Context, fabricPath, sourcePath string, tokens *map[string][]map[string]string) string {
	fixturePath := getFixtureFilePath(sourcePath)

	var tokensObj supertypes.MapValueOf[string]

	if tokens != nil && (*tokens)[fabricPath] != nil {
		for _, replacements := range (*tokens)[fabricPath] {
			tokensObj, _ = supertypes.NewMapValueOfMap(ctx, replacements)

			break
		}
	}

	payload, _, _ := transforms.SourceFileToPayload(ctx, types.StringValue(fixturePath), tokensObj)

	return *payload
}

func getFixtureFilePath(sourcePath string) string {
	_, filename, _, _ := runtime.Caller(0) //nolint:dogsled

	return filepath.Join(filepath.Dir(filename), "fixtures", sourcePath)
}

func panicOnError(err error) {
	if err != nil {
		panic(err) // lintignore:R009
	}
}
