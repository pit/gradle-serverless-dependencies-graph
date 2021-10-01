resource "aws_dynamodb_table" "storage" {
  name         = "${var.name_prefix}-storage"
  billing_mode = "PAY_PER_REQUEST"
  # read_capacity  = 20
  # write_capacity = 20

  hash_key = "Id"
  # range_key = "Version"

  attribute {
    name = "Id"
    type = "S"
  }

  attribute {
    name = "Dependency"
    type = "S"
  }

  attribute {
    name = "Version"
    type = "S"
  }

  attribute {
    name = "Repo"
    type = "S"
  }

  attribute {
    name = "Ref"
    type = "S"
  }

  global_secondary_index {
    name            = "Repository"
    hash_key        = "Repo"
    range_key       = "Ref"
    projection_type = "ALL"
  }

  global_secondary_index {
    name            = "Dependency"
    hash_key        = "Dependency"
    range_key       = "Version"
    projection_type = "ALL"
  }

  tags = {
    Name = "${var.name_prefix}-storage"
  }
}

resource "aws_dynamodb_table" "repositories" {
  name         = "${var.name_prefix}-repositories"
  billing_mode = "PAY_PER_REQUEST"
  # read_capacity  = 20
  # write_capacity = 20

  hash_key  = "Parent"
  range_key = "Child"

  attribute {
    name = "Parent"
    type = "S"
  }

  attribute {
    name = "Child"
    type = "S"
  }

  tags = {
    Name = "${var.name_prefix}-repositories"
  }
}

resource "aws_dynamodb_table" "dependencies" {
  name         = "${var.name_prefix}-dependencies"
  billing_mode = "PAY_PER_REQUEST"
  # read_capacity  = 20
  # write_capacity = 20

  hash_key  = "Parent"
  range_key = "Child"

  attribute {
    name = "Parent"
    type = "S"
  }

  attribute {
    name = "Child"
    type = "S"
  }

  tags = {
    Name = "${var.name_prefix}-dependencies"
  }
}
