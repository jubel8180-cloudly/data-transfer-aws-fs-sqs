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
	cd lambda/data_send_to_firehose_sqs && go build -o main




