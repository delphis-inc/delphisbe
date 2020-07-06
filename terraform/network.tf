# network.tf

# Fetch AZs in the current region
data "aws_availability_zones" "available" {
}

resource "aws_vpc" "main" {
  cidr_block = "172.17.0.0/16"
}

resource "aws_vpc" "ecs" {
  cidr_block = "172.31.0.0/16"
}

# Create var.az_count private subnets, each in a different AZ
resource "aws_subnet" "private" {
  count             = var.az_count
  cidr_block        = cidrsubnet(aws_vpc.main.cidr_block, 8, count.index)
  availability_zone = data.aws_availability_zones.available.names[count.index]
  vpc_id            = aws_vpc.main.id
}

# Create var.az_count public subnets, each in a different AZ
resource "aws_subnet" "public" {
  count                   = var.az_count
  cidr_block              = cidrsubnet(aws_vpc.main.cidr_block, 8, var.az_count + count.index)
  availability_zone       = data.aws_availability_zones.available.names[count.index]
  vpc_id                  = aws_vpc.main.id
  map_public_ip_on_launch = true
}

# Create var.az_count public subnets. hacky as these were created through the console and need to be imported in
resource "aws_subnet" "public-ecs-1" {
  cidr_block              = "172.31.48.0/20"
  availability_zone       = "us-west-2d"
  vpc_id                  = aws_vpc.ecs.id
  map_public_ip_on_launch = true
}

resource "aws_subnet" "public-ecs-2" {
  cidr_block              = "172.31.16.0/20"
  availability_zone       = "us-west-2b"
  vpc_id                  = aws_vpc.ecs.id
  map_public_ip_on_launch = true
}

resource "aws_subnet" "public-ecs-3" {
  cidr_block              = "172.31.0.0/20"
  availability_zone       = "us-west-2c"
  vpc_id                  = aws_vpc.ecs.id
  map_public_ip_on_launch = true
}

resource "aws_subnet" "public-ecs-4" {
  cidr_block              = "172.31.32.0/20"
  availability_zone       = "us-west-2a"
  vpc_id                  = aws_vpc.ecs.id
  map_public_ip_on_launch = true
}

# Internet Gateway for the public subnet
resource "aws_internet_gateway" "gw" {
  vpc_id = aws_vpc.main.id
}

# Internet Gateway for the public subnet
resource "aws_internet_gateway" "gw-ecs" {
  vpc_id = aws_vpc.ecs.id
}

# Route the public subnet traffic through the IGW
resource "aws_route" "internet_access" {
  route_table_id         = aws_vpc.main.main_route_table_id
  destination_cidr_block = "0.0.0.0/0"
  gateway_id             = aws_internet_gateway.gw.id
}

# Create a NAT gateway with an Elastic IP for each private subnet to get internet connectivity
resource "aws_eip" "gw" {
  count      = var.az_count
  vpc        = true
  depends_on = [aws_internet_gateway.gw]
}

resource "aws_nat_gateway" "gw" {
  count         = var.az_count
  subnet_id     = element(aws_subnet.public.*.id, count.index)
  allocation_id = element(aws_eip.gw.*.id, count.index)
}

# Create a new route table for the private subnets, make it route non-local traffic through the NAT gateway to the internet
resource "aws_route_table" "private" {
  count  = var.az_count
  vpc_id = aws_vpc.main.id

  route {
    cidr_block     = "0.0.0.0/0"
    nat_gateway_id = element(aws_nat_gateway.gw.*.id, count.index)
  }

  route {
    cidr_block                = "172.31.0.0/16"
    vpc_peering_connection_id = aws_vpc_peering_connection.main_to_ecs.id
  }
}

# Create a new route table for the private subnets, make it route non-local traffic through the NAT gateway to the internet
resource "aws_route_table" "private-ecs" {
  vpc_id = aws_vpc.ecs.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gw-ecs.id
  }

  route {
    cidr_block                = "172.17.0.0/16"
    vpc_peering_connection_id = aws_vpc_peering_connection.main_to_ecs.id
  }
}

# Explicitly associate the newly created route tables to the private subnets (so they don't default to the main route table)
resource "aws_route_table_association" "private" {
  count          = var.az_count
  subnet_id      = element(aws_subnet.private.*.id, count.index)
  route_table_id = element(aws_route_table.private.*.id, count.index)
}

resource "aws_vpc_peering_connection" "main_to_ecs" {
  peer_vpc_id = aws_vpc.ecs.id
  vpc_id      = aws_vpc.main.id

  tags = {
    Name = "main to ecs vpc"
  }
}

