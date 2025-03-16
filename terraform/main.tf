provider "aws" {
  region = var.aws_region
}

# Create a VPC
resource "aws_vpc" "fizzbuzz_vpc" {
  cidr_block           = var.vpc_cidr
  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Name = "fizzbuzz-vpc"
  }
}

# Create an Internet Gateway
resource "aws_internet_gateway" "fizzbuzz_igw" {
  vpc_id = aws_vpc.fizzbuzz_vpc.id

  tags = {
    Name = "fizzbuzz-igw"
  }
}

# Create a public subnet
resource "aws_subnet" "fizzbuzz_public_subnet" {
  vpc_id                  = aws_vpc.fizzbuzz_vpc.id
  cidr_block              = var.public_subnet_cidr
  map_public_ip_on_launch = true
  availability_zone       = "${var.aws_region}a"

  tags = {
    Name = "fizzbuzz-public-subnet"
  }
}

# Create a route table
resource "aws_route_table" "fizzbuzz_public_rt" {
  vpc_id = aws_vpc.fizzbuzz_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.fizzbuzz_igw.id
  }

  tags = {
    Name = "fizzbuzz-public-rt"
  }
}

# Associate the route table with the public subnet
resource "aws_route_table_association" "fizzbuzz_public_rt_assoc" {
  subnet_id      = aws_subnet.fizzbuzz_public_subnet.id
  route_table_id = aws_route_table.fizzbuzz_public_rt.id
}

# Create a security group for the EC2 instance
resource "aws_security_group" "fizzbuzz_sg" {
  name        = "fizzbuzz-sg"
  description = "Security group for FizzBuzz server"
  vpc_id      = aws_vpc.fizzbuzz_vpc.id

  # Allow HTTP access on port 8080
  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "HTTP for FizzBuzz API"
  }

  # Allow SSH access
  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
    description = "SSH access"
  }

  # Allow all outbound traffic
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "fizzbuzz-sg"
  }
}

# Create an EC2 key pair
resource "aws_key_pair" "fizzbuzz_key_pair" {
  key_name   = "fizzbuzz-key"
  public_key = file(var.public_key_path)
}

# Create an EC2 instance
resource "aws_instance" "fizzbuzz_server" {
  ami                    = var.ami_id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.fizzbuzz_key_pair.key_name
  vpc_security_group_ids = [aws_security_group.fizzbuzz_sg.id]
  subnet_id              = aws_subnet.fizzbuzz_public_subnet.id

  root_block_device {
    volume_size = 8
    volume_type = "gp2"
  }

  tags = {
    Name        = "fizzbuzz-server"
    DockerImage = var.docker_image
  }
}

# Create an Elastic IP for the EC2 instance
resource "aws_eip" "fizzbuzz_eip" {
  instance = aws_instance.fizzbuzz_server.id
  domain   = "vpc"

  tags = {
    Name = "fizzbuzz-eip"
  }
}

# Output the public IP of the EC2 instance
output "fizzbuzz_server_public_ip" {
  value = aws_eip.fizzbuzz_eip.public_ip
}

# Output the URL to access the FizzBuzz API
output "fizzbuzz_api_url" {
  value = "http://${aws_eip.fizzbuzz_eip.public_ip}:8080/fizzbuzz"
}