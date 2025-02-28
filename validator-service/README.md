# Validator Service 

This document provides a guide on how to use the Validator Service, which is a web service designed to create and manage validator requests. The service is built using the Go programming language and provides RESTful API endpoints for creating validators, checking their status, and monitoring the service's health and metrics.


## Overview

The Validator Service allows users to:

Create validator requests by specifying the number of validators and a fee recipient address.

Check the status of a validator request using its unique request_id.

Monitor the health of the service.

Access Prometheus metrics for monitoring request performance.

The service uses an SQLite database to store validator requests and their associated keys. It also integrates with Prometheus to provide metrics such as request counts and response times.

## API Endpoints
### Base URL

All endpoints are accessible under the base URL:
http://localhost:8080

### Create Validators

Creates a new validator request.

Endpoint:
`POST /validators`

Request Body:

```json
{
    "num_validators": 5,
    "fee_recipient": "0x1234567890123456789012345678901234567890"
}
```

Parameters:

`num_validators (uint)`: The number of validators to create. Must be greater than 0.

`fee_recipient (string)`: A valid Ethereum address to receive fees.

Response:

```json

{
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Validator creation in progress"
}
```

Response Codes:

`200 OK`: Validator request created successfully.

`400 Bad Request`: Invalid request body, invalid number of validators, or invalid fee recipient address.

`500 Internal Server Error`: Server error during validator creation.

### Check Validator Request Status
Retrieves the status of a validator request by its request_id.

Endpoint:
`GET /validators/{request_id}`

Response:

```json
{
    "status": "successful",
    "keys": [
        "key1",
        "key2",
        "key3",
        "key4",
        "key5"
    ]
}
```

Response Codes:

`200 OK`: Validator request found and returned.

`404 Not Found`: Validator request with the specified request_id not found.

`500 Internal Server Error`: Server error while processing the request.

### Health Check
Checks the health of the service, including database connectivity.

Endpoint:
`GET /health`

Response:

```json
{
    "status": "healthy"
}
```

Response Codes:

`200 OK`: Service is healthy.

`500 Internal Server Error`: Service is unhealthy (e.g., database connection issue).

### Metrics
   Provides Prometheus metrics for monitoring the service.

Endpoint:
`GET /metrics`

Response:
Prometheus metrics in plain text format.

Example:

```plaintext
# HELP http_requests_total Total number of requests received
# TYPE http_requests_total counter
http_requests_total{endpoint="/validators"} 10
http_requests_total{endpoint="/health"} 5
# HELP http_response_time_seconds Response time distribution
# TYPE http_response_time_seconds histogram
http_response_time_seconds_bucket{endpoint="/validators",le="0.1"} 7
http_response_time_seconds_bucket{endpoint="/validators",le="0.2"} 10
...
```

## Example Usage
### Create Validators

Send a POST request to create validators:

```bash

curl -X POST http://localhost:8080/validators \
-H "Content-Type: application/json" \
-d '{
    "num_validators": 3,
    "fee_recipient": "0x1234567890123456789012345678901234567890"
}'
```

Response:

```json
{
    "request_id": "550e8400-e29b-41d4-a716-446655440000",
    "message": "Validator creation in progress"
}
```

### Check Validator Request Status

Use the request_id to check the status:

```bash
curl -X GET http://localhost:8080/validators/550e8400-e29b-41d4-a716-446655440000
```

Response:

```json
{
    "status": "successful",
        "keys": [
        "key1",
        "key2",
        "key3"
    ]
}
```

### Health Check

Check the health of the service:

```bash
curl -X GET http://localhost:8080/health
```

Response:

```json
{
    "status": "healthy"
}
```

### Access Metrics
Access Prometheus metrics:

```bash
curl -X GET http://localhost:8080/metrics
```

## Monitoring and Metrics

### Prometheus
The service integrates with Prometheus to provide the following metrics:

`http_requests_total`: Total number of HTTP requests received, grouped by endpoint.

`http_response_time_seconds`: Distribution of response times for HTTP requests, grouped by endpoint.

These metrics can be scraped by Prometheus and visualized using tools like Grafana.

## Running the Service

1. Ensure Go is installed on your system.
2. Clone the repository.
3. Navigate to the project directory.

4. Run the service:

```bash
go run main.go
```

> The service will be available at http://localhost:8080.

## Docker

Service can be run locally using Dockerfile.

### Build docker image

```bash 
docker build -t validator-service .
```

### Run docker container

```bash 
docker run -p 8080:8080 validator-service
```

## Kubernetes

Service can be deployed on Kubernetes using deployment configuration in file `k8s-deployment.yaml`

### Deploy locally example on minikube

1. Install `minikube` from [official source](https://minikube.sigs.k8s.io/docs/start/?arch=%2Fwindows%2Fx86-64%2Fstable%2F.exe+download).
2. Install `kubectl` from [official source](https://kubernetes.io/docs/tasks/tools/)
3. Install `docker` from [official source](https://docs.docker.com/engine/install/)
4. Run minikube:
```bash
minikube start
```
5. Enable minikube tunnel (it should be running in a separate CLI window):

```bash
minikube tunnel
```

6. Apply kubernetes deployment configuration (in new CLI window):
```bash
kubectl apply -f .\k8s-deployment.yaml
```
7. Check if all pods are running. 

Run:
```bash
kubectl get pods
```
Output should be similar to:
```bash
NAME                                 READY   STATUS    RESTARTS   AGE
grafana-755fd46679-7lgpk             1/1     Running   0          9s
prometheus-b9ff8b547-pz75c           1/1     Running   0          9s
validator-service-659cd49474-p4b9n   1/1     Running   0          9s
validator-service-659cd49474-r92ct   1/1     Running   0          9s
```

8. Check if all services are running:

Run:
```bash
kubectl get svc
```
Output should be similar to:
```bash
NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
grafana-service      NodePort    10.110.218.43    <none>        3000:31466/TCP   18s
kubernetes           ClusterIP   10.96.0.1        <none>        443/TCP          15s
prometheus-service   NodePort    10.109.225.148   <none>        9090:32154/TCP   18s
validator-service    ClusterIP   10.109.249.91    <none>        8080/TCP         18s

```
9. To access any service you can use `EXTERNAL-IP` from previous command

10. If the `EXTERNAL-IP` is displayed as `<none>`, use the following Minikube commands to open the service bridges:

* Run Validator Service 
```bash
minikube service validator-service --url
```

* Run Grafana Service
```bash
minikube service grafana-service --url
```

## Grafana

### Dashboard

Dashboard contains two visualization:

* Average response time per endpoint
* Total number of requests per endpoint.

> Dashboard JSON model are located in `config/grafana-dashboard.json`

### Install

1) Take grafana `url` from previous step
2) Enter in grafana dashboard using:
    * username: `admin`
    * password: `admin`
3) Go to `Connections->Data Sources`
4) Press `Add Data Sources`
5) Choose `Prometheus`
6) In field `Prometheus server URL` enter `http://prometheus-service:9090` (prometheus address inside kubernetes)
7) Go to `Dashboards`
8) In the right top corner press `New -> Import`
9) In JSON model field copy content from file `config/grafana-dashboard.json`
10) Press load
11) You should see imported dashboard

### Possible Problems

After importing you can see issues with visualisations `Datasource <some_name> was not found`.
To solve it:
1) Open visualisation
2) Change Datasource to that you added before
3) Press `Save Dashboard` and `Run Queries`
4) You should see charts.