resource "aws_vpc" "vpc" {
    cidr_block = "10.0.0.0/16"
}

resource "aws_subnet" "public" {
    count = 3
    vpc_id = aws_vpc.vpc.id
    cidr_block = cidrsubnet(
      aws_vpc.vpc.cidr_block, 8, count.index
    )
}

resource "aws_s3_bucket" "bucket" {
    bucket = "hello-world-abc12345"
    acl    = "private"
}

output "bucket_arn" {
    value = aws_s3_bucket.bucket.arn
}

output "public_subnet_ids" {
    value = aws_subnet.public.*.id
}

output "vpc_id" {
    value = aws_vpc.vpc.id
}
