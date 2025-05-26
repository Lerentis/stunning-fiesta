variable "region" {
    description = "The AWS region to deploy resources in"
    type        = string
    default     = "eu-central-1"
}

variable "aws_tags" {
  type = map(string)
    default = {
        topic = "{{ Values.topicName }}"
        environment = "{{ Values.environment }}"
        heritage = "{{ Values.infrastructureRepository }}"
        managed_by = "terraform"
    }
}