include {
  path = find_in_parent_folders()
}

terraform {
  source = "../..//."
}

inputs = {
  region = "us-east-2"

  availability_zones = ["us-east-2b", "us-east-2c"]

  environment = "test"

  name = "example"
}
