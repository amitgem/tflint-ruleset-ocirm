plugin "ocirm" {
  enabled = true

  rule "oci_required_tags" {
    enabled = true
    required_tags "oci_core_instance" {
      tag {
        name = "Applications.CostCenter"
        values = ["AA", "BB" ]
      }
      tag {
        name = "HumanResource.Provider" 
        values = ["XX", "YY"]
      }
    }
    required_tags "oci_core_compartment" {
      tag {
        name = "Applications.CostCenter"
        values = ["CC", "DD" ]
      }
    }
  }
}