# FizzBuzz Server Terraform Deployment

This directory contains Terraform configuration to deploy an EC2 instance for the FizzBuzz server.

## Prerequisites

1. [Terraform](https://www.terraform.io/downloads.html) installed (v1.0.0+)
2. AWS CLI installed and configured with appropriate credentials
3. SSH key pair for EC2 access
4. Docker image of the FizzBuzz application published to DockerHub

## Configuration

Before deploying, update the `terraform.tfvars` file with your specific configuration:

```hcl
aws_region         = "us-east-1"  # Your preferred AWS region
vpc_cidr           = "10.0.0.0/16"  # VPC CIDR block
public_subnet_cidr = "10.0.1.0/24"  # Public subnet CIDR block
instance_type      = "t2.micro"  # EC2 instance type
ami_id             = "ami-0230bd60aa48260c6"  # Amazon Linux 2 AMI ID (region-specific)
public_key_path    = "~/.ssh/id_rsa.pub"  # Path to your public SSH key
docker_image       = "yourusername/fizzbuzz-server:latest"  # Your DockerHub image
```

## Deployment Steps

1. Initialize Terraform:
   ```
   terraform init
   ```

2. Preview the changes:
   ```
   terraform plan
   ```

3. Apply the configuration:
   ```
   terraform apply
   ```

4. When prompted, type `yes` to confirm the deployment.

5. After successful deployment, Terraform will output:
   - The public IP address of the EC2 instance (Elastic IP)
   - The URL to access the FizzBuzz API

## Manual Setup After Deployment

After the EC2 instance is deployed, you'll need to manually set up Docker and run the container:

1. SSH into the EC2 instance:
   ```
   ssh ec2-user@<ELASTIC_IP>
   ```

2. Install Docker:
   ```
   sudo yum update -y
   sudo amazon-linux-extras install docker -y
   sudo service docker start
   sudo systemctl enable docker
   sudo usermod -a -G docker ec2-user
   ```

3. Log out and log back in for the group changes to take effect:
   ```
   exit
   ssh ec2-user@<ELASTIC_IP>
   ```

4. Pull and run the Docker container:
   ```
   docker pull yourusername/fizzbuzz-server:latest
   docker run -d --name fizzbuzz-server --network host -e PORT=8080 --restart unless-stopped yourusername/fizzbuzz-server:latest
   ```

## Accessing the Application

The FizzBuzz API will be available at:
```
http://<ELASTIC_IP>:8080/fizzbuzz
```

Example API call:
```
http://<ELASTIC_IP>:8080/fizzbuzz?int1=3&int2=5&limit=15&str1=fizz&str2=buzz
```

## Cleanup

To destroy all resources created by Terraform:
```
terraform destroy
```

When prompted, type `yes` to confirm.