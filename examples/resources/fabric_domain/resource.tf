resource "fabric_domain" "parent" {
  display_name = "example parent"
}

resource "fabric_domain" "child" {
  display_name     = "example child"
  description      = "This is an example child domain"
  parent_domain_id = fabric_domain.parent.id
}

#domain update example
# resource "fabric_domain" "example_update" {
#   display_name     = "example child"
#   description      = "This is an example updated domain"
#   default_label_id = "11111111-1111-1111-1111-111111111111"
# }
