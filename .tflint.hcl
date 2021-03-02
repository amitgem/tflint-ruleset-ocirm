plugin "ocirm" {
  enabled = true

  rule "oci_required_tags" {
    enabled = false
    required_tags "oci_core_instance" {
      tag "Applications.CostCenter" {
        values = ["AA", "BB" ]
      }
      tag "HumanResource.Provider" {
        values = ["XX", "YY"]
        }
      }
  }
}