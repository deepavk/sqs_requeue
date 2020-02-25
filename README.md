This script can be used on aws lambda to requeues messages from the dead letter queue back to the source queue. 

The lambda script takes 3 environment variables

destination_queue
region
source_queue

Run lambda zip_dlq_script.sh to create a zip file that can be uploaded on aws lambda
