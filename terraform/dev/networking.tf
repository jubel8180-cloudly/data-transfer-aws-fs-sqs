
# this module will create a custome vpc, subnets, gateway and routing table
module "my_vpc" {
    source = "../modules/vpc"
    vpc_cidr = "172.31.0.0/16"
    tenancy = "default"
    vpc_id = "${module.my_vpc.vpc_id}"
    subnet_cidrs = ["172.31.0.0/20","172.31.16.0/20"]
    availability_zones = ["${var.region}a","${var.region}b"]
}