resource "fabric_domain_role_assignments" "example" {
  domain_id = "00000000-0000-0000-0000-000000000000"
  role      = "Admins"
  principals = [
    {
      id   = "11111111-1111-1111-1111-111111111111"
      type = "User"
    },
    {
      id   = "22222222-2222-2222-2222-222222222222"
      type = "Group"
    }
  ]
}
