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
      - Install terraform
    - Add credential in AWS cli
      # provide your aws access key and secret key
      $ aws configure
      - provide access key
      - provide secret key
    

    - Apply terraform functionality 
      - go to the terraform folder
      - then go to dev folder
      $ cd terraform/dev
      
      # initialized the terraform
      $ terraform init

      # check any configuration is write or wrong
      $ terraform plan

      # deploy all configuration
      $ terraform apply
      - After applying terraform, you need to provide input.
         - s3 bucket name (name must be camel case or multiple word will be add with dash(-))

     - After that every setup will be done for this project.
     


## Test the full project
- Now go to the application load balancer. 
- Copy the url link
- now test the code and check the requirements

  - copy the url and provide condition in payload "flag": "A" for firehose or "flag": "Y" for sqs
  - POST: lb-firehose-sqs-dev-264299050.ap-south-1.elb.amazonaws.com
  - provide json data in Body


```
******** Test Firehose Data Transfer ********

- Sample Input Data for kinesis firehose
- POST: lb-firehose-sqs-dev-264299050.ap-south-1.elb.amazonaws.com

{
  "condition": {
    "flag": "A"
  },
  "records": [
    {
      "title": "E sign save as template process 2017 1",
      "desc": "",
      "signers": [
        {
          "name": "Chris Ward",
          "email": "office@dragonsdesign.co.uk"
        },
        {
          "name": "Sandra Norris",
          "email": "office12@dragonsdesign.co.uk"
        }
      ]
    }
  ]
}

******** Test Send Message to SQS ********

- POST: lb-firehose-sqs-dev-264299050.ap-south-1.elb.amazonaws.com
- Sample Input Data for sqs

{
   "condition": {
    "flag": "Y"
    },
    "records": {
        messageId: '19dd0b57-b21e-4ac1-bd88-01bbb068cb78',
        receiptHandle: 'MessageReceiptHandle',
        "body": "Hello from SQS!",
        attributes: {
            ApproximateReceiveCount: '1',
            SenderId: '012345678910'
        }
    }
    
}

```

## lambda function development process: 

```
- go to the lambda directory
- different lambda function available in different directory
- go the the specific directory

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
                  