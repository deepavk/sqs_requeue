AWS sqs: 

If a message fails to be consumed it is sent to the dead letter queue on SQS. An AWS lambda function can be setup to requeue messages from dead letter queue to source queue
This script can be used to requeues messages. 

The environment variables have to be setup for:
```
1. destination_queue
2. region
3. source_queue
```
lambda zip_dlq_script.sh can be used to create a zip file that can be uploaded
