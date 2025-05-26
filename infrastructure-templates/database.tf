module "database" {
  source = "git:://github.com/terraform-aws-modules/terraform-aws-rds"
  version = "3.0.0"
  name = var.db_name
  engine = var.db_engine
  engine_version = var.db_engine_version
}