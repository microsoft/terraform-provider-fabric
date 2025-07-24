resource "fabric_onelake_data_access_security" "example" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  item_id = "11111111-1111-1111-1111-111111111111"
	value = [
		{
			name = "example",
			decision_rules = [
				{
					effect = "Permit"
					permission = [
						{
							attribute_name = "Path"
							attribute_value_included_in = ["*"]
						},
						{
							attribute_name = "Action"
							attribute_value_included_in = ["Read"]
						}
					]
				}
			]
			members = {
				fabric_item_members = [
					{
						item_access = ["ReadAll"]
						source_path = "00000000-0000-0000-0000-000000000000/11111111-1111-1111-1111-111111111111"
					}
				]
			}
		}
	]
}
