.PHONY: apply

apply:
	cd terraform/dev && terraform apply

destroy:
	cd terraform/dev && terraform destroy

plan:
	cd terraform/dev && terraform plan

output:
	cd terraform/dev && terraform output

build-l1:
	cd lambda/data_send_to_firehose_and_sqs && go build -o main

build-l2:
	cd lambda/receive_msg_from_dlq && go build -o main

build-all: build-l1 build-l2



