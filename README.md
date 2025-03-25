# car-thingy_GO

Ez lenni repository for car-thingy_GO

This backend appliation is part of the car-thingy system. This app is written in GO, acts as the REST API for CRUD operations for every model.

## Setup
- Have an operational MySQL database

### Terraform
- Run the following commands
```terraform
terraform init

terraform apply \
    -var="container_name=car-thingy_go" \
    -var="container_version=latest" \
    -var="env=prod" \
    -var="db_username=<db username>" \
    -var="db_password=<db passwd>" \
    -var="db_ip=<db ip>" \
    -var="db_port=<db port>" \
    -var="db_name=<db name>" \
    -var="api_secret=<random secret>" \
    -auto-approve
```

### Docker only
- Run the following command
```sh
docker run --mount source=downloaded_images,target=/app/downloaded_images \
    -e "DB_IP=<db ip>" \
    -e "DB_USERNAME=<db username>" \
    -e "DB_PASSWORD=<db passwd>" \
    -e "DB_PORT=<db port>" \
    -e "DB_NAME=<db name>" \
    -e "API_SECRET=<random secret>"
    -p <exposed port>:3000 \
    --restart unless-stopped \
    --name=car-thingy_GO \
    sc4n1a471/car-thingy_go:latest
```

## Usage
- First, you need to generate an API key by calling the `POST /auth` endpoint with your secret as `x-api-key` header (only 1 active API key is allowed)
- You will need to use this key for every operation, also need to pass it to car-thingy_Python
- You can find the endpoints of the API in `openapi.yaml`
