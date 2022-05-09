.PHONY: apply

apply:
	cd terraform/application && terraform apply

destroy:
	cd terraform/application && terraform destroy

destroy-dev:
	cd terraform/application && terraform workspace select default && terraform destroy

destroy-staging:
	cd terraform/application && terraform workspace select staging && terraform destroy

destroy-prod:
	cd terraform/application && terraform workspace select prod && terraform destroy

destroy-all: destroy-dev destroy-staging destroy-prod

plan:
	cd terraform/application && terraform plan

output:
	cd terraform/application && terraform output

build-l1:
	cd lambda/data_send_to_firehose_and_sqs && go build -o main

build-l2:
	cd lambda/receive_msg_from_dlq && go build -o main

build-l3:
	cd lambda/kinesis_lambda_processor && go build -o main

build-l4:
	cd lambda/s3_event_for_metadata_and_checksum && go build -o main

build-all: build-l1 build-l2 build-l3 build-l4




