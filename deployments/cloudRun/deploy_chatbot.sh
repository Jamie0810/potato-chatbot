#!/usr/bin/env bash
. ./deployments/cloudRun/env/chatbot.env

gcloud run deploy "${SERVICE_NAME}" \
      --image "${IMAGE}" \
      --allow-unauthenticated \
      --cpu 1000m \
      --memory 128Mi \
      --max-instances 1 \
      --platform managed \
      --region asia-east1 \
      --project "${PROJECT_NAME}" \
      --add-cloudsql-instances="${CLOUDSQL_INSTANCES}" \
      --set-env-vars  ENV="${ENV}" \
      --set-env-vars  DATABASE_DRIVER="${DATABASE_DRIVER}" \
      --set-env-vars  DATABASE_INSTANCE_NAME="${CLOUDSQL_INSTANCES}" \
      --set-env-vars  DATABASE_USER="${DATABASE_USER}" \
      --set-env-vars  DATABASE_PASSWORD="${DATABASE_PASSWORD}" \
      --set-env-vars  DATABASE_DBNAME="${DATABASE_DBNAME}" \
      --set-env-vars  DATABASE_LOG_MODE="${DATABASE_LOG_MODE}" \
      --set-env-vars  LOG_LEVEL="${LOG_LEVEL}" \
      --set-env-vars  LOG_FORMAT="${LOG_FORMAT}" \
      --set-env-vars  PUBSUB_PROJECTID="${PROJECT_NAME}" \
      --set-env-vars  PUBSUB_TOPIC_CHATBOT="${PUBSUB_TOPIC_CHATBOT}" \
      --set-env-vars  ^$$^PUBSUB_CREDENTIALS="${PUBSUB_CREDENTIALS}" \
      --args 'chatbot','--port','8080'

