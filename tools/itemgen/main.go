package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// type ItemConfig struct {
// 	Name              string
// 	TypeInfo          string
// 	FabricItemType    string
// 	DefinitionFormats []string
// 	DocsURL           string
// 	DisplayNameMax    int
// 	DescriptionMax    int
// 	ItemType          ItemType
// 	HasDefinition     bool
// 	HasProperties     bool
// 	HasConfig         bool
// }

type ItemConfig struct {
	Name              string
	Type              string
	TypeInfo		  string
	Names             string
	Types             string
	TypesInfo		  string
	RenameAllowed     bool
	Package	 	      string
	FabricItemType    string
	DefinitionFormats []string
	DocsURL           string
	DisplayNameMax    int
	DescriptionMax    int
	ItemType          ItemType
	HasDefinition     bool
	HasProperties     bool
	HasConfig         bool
	IsPreview         bool
	IsSPNSupported    bool
}


type ItemType int

const (
	TypeBasic ItemType = iota
	TypeDefinition
	TypeProperties
	TypeDefinitionProperties
	TypeConfigProperties
	TypeConfigDefinitionProperties
)

func (t ItemType) String() string {
	switch t {
	// case TypeBasic:
	// 	return "basic"
	case TypeDefinition:
		return "definition"
	case TypeProperties:
		return "properties"
	case TypeDefinitionProperties:
		return "definition-properties"
	case TypeConfigProperties:
		return "config-properties"
	case TypeConfigDefinitionProperties:
		return "config-definition-properties"
	default:
		return "unknown"
	}
}


func main() {
	// Parse command line flags
	itemName := flag.String("item-name", "", "Name of the new item (e.g. Data Pipeline)")
	itemsName := flag.String("items-name", "", "Name of the new item in plural form (e.g. Data Pipelines)")
	itemTypeFlag := flag.String("item-type", "", "Type of item (definition|properties|definition-properties|config-properties|config-definition-properties)")
	renameAllowed := flag.Bool("rename-allowed", true, "Is item rename allowed?")
	isPreview := flag.Bool("is-preview", false, "Is the item in preview?")
	IsSPNSupported := flag.Bool("is-spn-supported", false, "Is the item supported for SPN?")
	flag.Parse()

	// Parse item type
	var itemTypeEnum ItemType
	switch *itemTypeFlag {
	case "basic":
		itemTypeEnum = TypeBasic
	case "definition":
		itemTypeEnum = TypeDefinition
	case "properties":
		itemTypeEnum = TypeProperties
	case "definition-properties":
		itemTypeEnum = TypeDefinitionProperties
	case "config-properties":
		itemTypeEnum = TypeConfigProperties
	case "config-definition-properties":
		itemTypeEnum = TypeConfigDefinitionProperties
	default:
		fmt.Printf("Error: Invalid item type %s. Must be one of: basic, definition, properties, definition-properties, config-properties, config-definition-properties\n", *itemTypeFlag)
		os.Exit(1)
	}

	// Create item configuration
	config := ItemConfig{
		Name:              *itemName,
		Type:              strings.ToLower(strings.ReplaceAll(*itemName, " ", "_")),
		TypeInfo:		   strings.ReplaceAll(*itemName, " ", ""),
		Names:             *itemsName,
		Types:             strings.ToLower(strings.ReplaceAll(*itemsName, " ", "_")),
		TypesInfo:		   strings.ReplaceAll(*itemsName, " ", ""),
		Package: 	       strings.ToLower(strings.ReplaceAll(*itemName, " ", "")),
		RenameAllowed:	   *renameAllowed,
		DefinitionFormats: []string{"<part1>", "<part2>"},
		DocsURL:           "<docs-url>",
		DisplayNameMax:    123,
		DescriptionMax:    256,
		ItemType:          itemTypeEnum,
		HasDefinition:     itemTypeEnum == TypeDefinition || itemTypeEnum == TypeDefinitionProperties || itemTypeEnum == TypeConfigDefinitionProperties,
		HasProperties:     itemTypeEnum == TypeProperties || itemTypeEnum == TypeDefinitionProperties || itemTypeEnum == TypeConfigProperties || itemTypeEnum == TypeConfigDefinitionProperties,
		HasConfig:         itemTypeEnum == TypeConfigProperties || itemTypeEnum == TypeConfigDefinitionProperties,
		IsPreview:         *isPreview,
		IsSPNSupported:    *IsSPNSupported,
	}

	// Create the item directory
	itemDir := filepath.Join("internal", "services", config.Package)
	if err := os.MkdirAll(itemDir, 0755); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", itemDir, err)
		os.Exit(1)
	}

	// Check if directory already contains files
	files, err := os.ReadDir(itemDir)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", itemDir, err)
		os.Exit(1)
	}
	if len(files) > 0 {
		fmt.Printf("Warning: Directory %s already contains files. This may overwrite existing files.\n", itemDir)
		fmt.Print("Do you want to continue? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("Operation cancelled")
			os.Exit(0)
		}
	}

	// Generate files based on item type
	filesToGenerate := getFilesForItemType(config.Type, config.Types, itemTypeEnum)

	for _, file := range filesToGenerate {
		if err := generateFile(itemDir, filepath.Join("tools/itemgen/templates", *itemTypeFlag, file.template), file.output, config); err != nil {
			fmt.Printf("Error generating %s: %v\n", file.output, err)
			os.Exit(1)
		}
		fmt.Printf("Generated %s\n", filepath.Join(itemDir, file.output))
	}

	fmt.Printf("\nSuccessfully generated item %s in %s\n", *itemName, itemDir)
	fmt.Println("\nNext steps:")
	fmt.Println("1. Review the generated files")
	fmt.Println("2. Update the documentation URL if needed")
	fmt.Println("3. Add any service-specific logic")
	fmt.Println("4. Run the tests to verify the implementation")
}

type fileInfo struct {
	template string
	output   string
}

func getFilesForItemType(typeName string, typesName string, itemType ItemType) []fileInfo {
	files := []fileInfo{
		{"base.go.tmpl", "base.go"},
		{"base_test.go.tmpl", "base_test.go"},
		{"data_item.go.tmpl", "data_" + typeName + ".go"},
		{"data_item_test.go.tmpl", "data_" + typeName + "_test.go"},
		{"data_items.go.tmpl", "data_" + typesName + ".go"},
		{"data_items_test.go.tmpl", "data_" + typesName + "_test.go"},
		{"resource_item.go.tmpl", "resource_" + typeName + ".go"},
		{"resource_item_test.go.tmpl", "resource_" + typeName + "_test.go"},
	}


	switch itemType {
	case TypeProperties, TypeDefinitionProperties, TypeConfigProperties, TypeConfigDefinitionProperties:
		typeSpecificFiles := []fileInfo{
			{"schema_data_item.go.tmpl", "schema_data_" + typeName + ".go"},
			{"schema_resource_item.go.tmpl", "schema_resource_" + typeName + ".go"},
			{"models.go.tmpl", "models.go"},
		}
		files = append(files, typeSpecificFiles...)
	}


	return files
}

func generateFile(dir, tmplPath, outputFile string, config ItemConfig) error {
	// Check if template file exists
	if _, err := os.Stat(tmplPath); os.IsNotExist(err) {
		return fmt.Errorf("template file not found: %s", tmplPath)
	}

	// Read template
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("error parsing template %s: %v", tmplPath, err)
	}

	// Create output file
	output := filepath.Join(dir, outputFile)
	f, err := os.Create(output)
	if err != nil {
		return fmt.Errorf("error creating output file %s: %v", output, err)
	}
	defer f.Close()

	// Execute template
	if err := tmpl.Execute(f, config); err != nil {
		return fmt.Errorf("error executing template %s: %v", tmplPath, err)
	}

	return nil
}
