This script is to be used for the aws lambda function that requeues messages from the dead letter queue back to the live queue. 

The lambda takes 3 environment variables

destination_queue
region
source_queue

To create a zip file that can be uploaded for the lambda zip_dlq_script.sh can be run 