aws_region         = "us-east-1"
vpc_cidr           = "10.0.0.0/16"
public_subnet_cidr = "10.0.1.0/24"
instance_type      = "t2.micro"
ami_id             = "ami-0230bd60aa48260c6" # Amazon Linux 2 AMI in us-east-1
public_key_path    = "~/.ssh/id_rsa.pub"     # Update this to your actual public key path
docker_image       = "yourusername/fizzbuzz-server:latest" # Replace with your actual DockerHub image