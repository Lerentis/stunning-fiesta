terraform {
  backend "s3" {
    bucket         = "{{ Values.stateBucket }}"
    key            = "{{ Values.stateRegion }}/{{ Values.topicName }}/atlantis.tfstate"
    region         = "eu-central-1"
    encrypt        = true
  }
}