#!/bin/bash
set -e
#
#  This script implements the locking mechanism functions to not having concurrent jobs updating the S3 repos
#  (APT, YUM, ZYPP...), to avoid wrong repo metadata. It creates a lock for each repo type in DynamoDB.
#
#
function create_dynamo_table {
  # Setup DynamoDB table
  if [ -z ${DYNAMO_TABLE_NAME+x} ]; then
    echo "$DYNAMO_TABLE_NAME is unset"
    exit 1
  fi
  if aws dynamodb describe-table --table-name $DYNAMO_TABLE_NAME --region $AWS_DEFAULT_REGION >/dev/null 2>&1 ; then
    echo "===> Dynamodb lock table already exists"
  else
    aws dynamodb create-table \
            --region $AWS_DEFAULT_REGION \
            --table-name $DYNAMO_TABLE_NAME \
            --attribute-definitions AttributeName=lock-type,AttributeType=S \
            --key-schema AttributeName=lock-type,KeyType=HASH \
            --sse-specification Enabled=true \
            --provisioned-throughput ReadCapacityUnits=2,WriteCapacityUnits=1
    aws dynamodb wait table-exists --table-name $DYNAMO_TABLE_NAME --region $AWS_DEFAULT_REGION
    aws dynamodb put-item \
        --table-name $DYNAMO_TABLE_NAME \
        --item '{"lock-type": {"S": "yum"}, "locked": {"BOOL": false}, "repo": {"S": "-"}}'
    aws dynamodb put-item \
        --table-name $DYNAMO_TABLE_NAME \
        --item '{"lock-type": {"S": "apt"}, "locked": {"BOOL": false}, "repo": {"S": "-"}}'
    aws dynamodb put-item \
        --table-name $DYNAMO_TABLE_NAME \
        --item '{"lock-type": {"S": "zypp"}, "locked": {"BOOL": false}, "repo": {"S": "-"}}'
    aws dynamodb put-item \
        --table-name $DYNAMO_TABLE_NAME \
        --item '{"lock-type": {"S": "win"}, "locked": {"BOOL": false}, "repo": {"S": "-"}}'
    aws dynamodb put-item \
        --table-name $DYNAMO_TABLE_NAME \
        --item '{"lock-type": {"S": "tarball"}, "locked": {"BOOL": false}, "repo": {"S": "-"}}'
  fi
}

function wait_free_lock {
  echo "===> Wait for Lock to be released, if takes long unlock DynamoDB item manually"
  while true; do
    locked=$(aws dynamodb get-item \
       --table-name ${DYNAMO_TABLE_NAME}  \
       --key "{ \"lock-type\": {\"S\": \"${LOCK_REPO_TYPE}\"} }" \
       --projection-expression "locked" \
      | jq -r '.Item.locked.BOOL');
    if [[ $locked == "false" ]]; then
      break
    fi
    sleep 10
  done
}

function lock {
  echo "===> Locking $LOCK_REPO_TYPE"
  aws dynamodb put-item \
    --table-name $DYNAMO_TABLE_NAME \
    --item "{\"lock-type\": {\"S\": \"${LOCK_REPO_TYPE}\"}, \"locked\": {\"BOOL\": true}}, \"repo\": {\"S\": \"${REPO_FULL_NAME}\"}}"
}

function release_lock {
  echo "===> Release Lock in $LOCK_REPO_TYPE"
  aws dynamodb put-item \
    --table-name $DYNAMO_TABLE_NAME \
    --item "{\"lock-type\": {\"S\": \"${LOCK_REPO_TYPE}\"}, \"locked\": {\"BOOL\": false}}, \"repo\": {\"S\": \"-\"}}"
}
