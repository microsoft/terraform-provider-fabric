# Create tags first
resource "fabric_tag" "environment_tag" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  display_name = "Environment:Development"
}

resource "fabric_tag" "team_tag" {
  workspace_id = "00000000-0000-0000-0000-000000000000"
  display_name = "Team:DataEngineering"
}

# Create environment with tags
resource "fabric_environment" "example" {
  display_name = "example"
  workspace_id = "00000000-0000-0000-0000-000000000000"
  tags         = [
    fabric_tag.environment_tag.id,
    fabric_tag.team_tag.id,
  ]
}
