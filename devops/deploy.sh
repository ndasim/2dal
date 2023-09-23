aws configure set aws_access_key_id > echo $AWS_ACCESS_KEY_ID
aws configure set aws_secret_access_key > echo $AWS_SECRET_ACCESS_KEY
aws configure set default.region > echo eu-central-1
aws ecr get-login-password --region eu-central-1 | sudo docker login --username AWS --password-stdin 163166970182.dkr.ecr.eu-central-1.amazonaws.com
sudo docker rm -f `sudo docker ps -q` 
sudo docker run -d -p 80:80 $ECR_REGISTRY/2dal_container:latest
```