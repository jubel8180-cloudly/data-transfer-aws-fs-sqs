# firehose File Creation from Lambda trigger


## About <a name = "about"></a>

This project will work for sending data to firehose or SQS. Firehose will put object directly to s3 bucket. If any error happened while sending data to SQS or firehose, data will be sent to the dead letter queue. When the Dead letter queue gets any message, it writes data to Cloudwatch through the other lambda function.

### Workflow
1. Create a application load balancer
2. Trigger a aws lambda function
3. Add a flag type A and type Y in payload condition
4. If flag type A call aws kinesis firehose and write data to aws s3
5. If flag type Y send message to aws sqs
6. Otherwise send message to aws dead lock.

## AWS Configuration Setup Instruction using terraform
    - Initial Configuration
      - Please install AWS cli
        ~ url: https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html
      - Install terraform
        ~ url: https://learn.hashicorp.com/tutorials/terraform/install-cli?in=terraform/aws-get-started
    - Add credential in AWS cli
      # provide your aws access key and secret key
      $ aws configure
      - provide access key
      - provide secret key

    - Need to configure first in Terraform
      - Please go to the terraform/dev folder
      - go to the terraform.tfvars file and change some variable name
        - region = "please provide the region where you want to configure"
        - ingress_cidr_blocks = "Please provide your public ip address where from you can access application load balancer"
          - for example:
            # ingress_cidr_blocks = ["your_public_ip_1/32","your_public_ip_2/32"]
            - ingress_cidr_blocks = ["37.111.201.253/32","27.147.201.239/32"]
    - Apply terraform functionality
      - go to the terraform folder
      - then go to dev folder
      $ cd terraform/dev

      # check any configuration is write or wrong
      $ terraform plan

      # Now deploy all configuration using this apply command
      $ terraform apply
      - After applying terraform, you need to provide input.
         - s3 bucket name (name must be camel case or multiple word will be add with dash(-))

         - Next terraform will ask you about aws changes. yes or no ?
          - please provide yes

         - After that every setup will be done for this project.

      Or you can apply terraform command with bucket name like,
        # terraform apply -var="bucket_name=your-bucket-name"
        $ terraform apply -var="bucket_name=new-firehose-data-storage"

        - Next terraform will ask you about aws changes. yes or no ?
          - please provide yes

         - After that every setup will be done for this project.

     # Now you can show the all created services name using this command
     $ terraform output

## Remember [Note]:
  - We used the region “ap-south-1”. So, every AWS feature will be available in this region.
  - If you want to change the region or other service name, just go to the terraform.tfvars file in the     terraform/dev directory and change the region and others service name.


## Test the full project
- Now go to the application load balancer and copy the load balancer dns link
or you can copy load balancer dns link from the terminal after providing command [$ terraform output]
- now test the code and check the requirements
  - copy the url and provide condition in payload "flag": "A" for firehose or "flag": "Y" for sqs
  - POST: Url link like, lb-firehose-sqs-dev-1432623056.ap-south-1.elb.amazonaws.com
  - provide json data in Body


```
******** Test Firehose Data Transfer ********

- Sample Input Data for kinesis firehose
- POST: lb-firehose-sqs-dev-1432623056.ap-south-1.elb.amazonaws.com

{
  "condition": {
    "flag": "A"
  },
  "records": [
    {
      "app_id":"app_1",
      "device_id":"608d99c2-9ff2-40a1-9ab5-8ee97861ccb2",
      "request_id":"19cc1f31-2a2e-4897-a8fd-a8dddf7a8e56"
    },
    {
      "app_id":"app_2",
      "device_id":"7d4586f0-e73f-49f0-ba9d-1b66f5413580",
      "request_id":"21805093-e285-40ed-b2c8-d0c7164f7953    "
    }
    ,{
      "app_id":"app_2",
      "device_id":"11fd3383-11f8-47f9-a110-37e2e409fdc9",
      "request_id":"f86b8e64-bc5f-4486-96ba-5792dc1836fd"
    },
    {
      "app_id":"app_1",
      "device_id":"887e6f18-2850-44df-a896-183403f94a63",
      "request_id":"d24c7916-fdf9-4275-8eac-926c16188615"
    },
    {
      "app_id":"app_1",
      "device_id":"07e140fa-ebb1-485f-848e-ab8d743f5788",
      "request_id":"4ed50b8b-5b47-4a62-a63a-2df073a3e602"
    },
    {
      "app_id":"app_2",
      "device_id":"83962919-5fbf-4790-8f17-ceeef6e11b4b",
      "request_id":"f77e8c6e-4852-47ec-b52a-786f1816a617"
    }
  ]
}

******** Test Send Message to SQS ********

- POST: lb-firehose-sqs-dev-1432623056.ap-south-1.elb.amazonaws.com
- Sample Input Data for sqs

{
   "condition": {
    "flag": "Y"
    },
    "records": {
        "body": "All provider data collected successfully",
        "attributes": {
            "status": "success",
            "provider": "dynata"

        }
    }
}

```

## lambda function development process:

```
- go to the lambda directory
- different lambda function available in different directory
- go to the specific directory

# execute these go commands to build go executable file

$ go mod tidy

# Build the executable code
go build -o {a file name}
go build -o main

# optional:
  - Zip file will create automatically using terraform.
  - However, Make zip file for uploading as lambda function in aws manually.
  - zip main.zip main

```


