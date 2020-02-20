provider "aws" {
  profile    = "default"
  region     = "us-east-1"
}

data "aws_ami" "amazon-linux-2" {
  most_recent = true
  owners = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm*gp2"]
  }
  filter {
    name   = "architecture"
    values = ["x86_64"]
  }
}

resource "aws_key_pair" "demo-key" {
  key_name = "demo-key"
  public_key = file("demo-key.pub")
}

resource "aws_security_group" "hyparview-demo-sg" {
  name   = "hyparview-demo-sg"

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # these rules should be limited to our vpc
  # our GRPC ports
  ingress {
    from_port   = 10000
    to_port     = 20000
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # stats
  ingress {
    from_port   = 23456
    to_port     = 23456
    protocol    = "udp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_instance" "demo" {
  ami                         = data.aws_ami.amazon-linux-2.id
  instance_type               = "t3.medium"
  associate_public_ip_address = true
  key_name                    = "demo-key"
  security_groups             = ["hyparview-demo-sg"]

  count = 5
  tags = {
    Name = "hyparview-demo-${count.index}"
  }
}


