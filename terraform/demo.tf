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

resource "aws_instance" "demo" {
  ami                         = data.aws_ami.amazon-linux-2.id
  instance_type               = "t3.medium"
  associate_public_ip_address = true

  count = 5
  tags = {
    Name = "hyparview-demo-${count.index}"
  }


  provisioner "remote-exec" {

  }
}


