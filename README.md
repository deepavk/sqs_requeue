This script can be used on aws lambda to requeues messages from the dead letter queue back to the source queue. 

The lambda script takes 3 environment variables
```
1. destination_queue
2. region
3. source_queue
```
lambda zip_dlq_script.sh can be used to create a zip file that can be uploaded
