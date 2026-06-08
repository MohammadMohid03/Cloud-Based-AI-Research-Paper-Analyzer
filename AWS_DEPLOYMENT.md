# ☁️ AWS Deployment Guide — AI-Powered Research Paper Analyzer

> This guide walks through deploying the entire application to AWS using production-grade services: RDS, S3, EC2 / Elastic Beanstalk, CloudFront, and Amazon Bedrock.

---

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Prerequisites](#prerequisites)
- [Step 1: Set Up VPC and Security Groups](#step-1-set-up-vpc-and-security-groups)
- [Step 2: Create RDS PostgreSQL Instance](#step-2-create-rds-postgresql-instance)
- [Step 3: Create S3 Bucket for File Storage](#step-3-create-s3-bucket-for-file-storage)
- [Step 4: Set Up IAM Roles and Policies](#step-4-set-up-iam-roles-and-policies)
- [Step 5: Configure Amazon Bedrock Access](#step-5-configure-amazon-bedrock-access)
- [Step 6: Deploy Backend](#step-6-deploy-backend)
  - [Option A: EC2 with Docker](#option-a-ec2-with-docker)
  - [Option B: Elastic Beanstalk](#option-b-elastic-beanstalk)
- [Step 7: Deploy Frontend to S3 + CloudFront](#step-7-deploy-frontend-to-s3--cloudfront)
- [Step 8: Set Up CloudWatch Monitoring](#step-8-set-up-cloudwatch-monitoring)
- [Step 9: Optional — SNS Notifications](#step-9-optional--sns-notifications)
- [Step 10: Optional — API Gateway + Lambda](#step-10-optional--api-gateway--lambda)
- [Cost Estimation](#cost-estimation)
- [Security Best Practices](#security-best-practices)
- [Troubleshooting](#troubleshooting)

---

## Architecture Overview

```
                         ┌───────────────────────────────┐
                         │         Route 53              │
                         │   (DNS: app.example.com)      │
                         └──────────┬────────────────────┘
                                    │
                    ┌───────────────┴───────────────┐
                    │                               │
             ┌──────┴──────┐                ┌──────┴──────┐
             │ CloudFront  │                │     ALB     │
             │ (Frontend)  │                │  (Backend)  │
             │ CDN         │                │             │
             └──────┬──────┘                └──────┬──────┘
                    │                               │
             ┌──────┴──────┐                ┌──────┴──────┐
             │  S3 Bucket  │                │   EC2 / EB  │
             │  (Static    │                │  (Go API)   │
             │   Assets)   │                │             │
             └─────────────┘                └──────┬──────┘
                                                   │
                                    ┌──────────────┼──────────────┐
                                    │              │              │
                             ┌──────┴──────┐┌─────┴─────┐┌──────┴──────┐
                             │    RDS      ││    S3     ││  Bedrock   │
                             │ PostgreSQL  ││ (Uploads) ││  (Claude)  │
                             │ (Database)  ││           ││            │
                             └─────────────┘└───────────┘└────────────┘

                             ┌─────────────┐┌───────────┐
                             │ CloudWatch  ││   SNS     │
                             │ (Monitoring)││ (Alerts)  │
                             └─────────────┘└───────────┘
```

---

## Prerequisites

Before starting, ensure you have:

- [ ] **AWS Account** with administrative access
- [ ] **AWS CLI v2** installed and configured (`aws configure`)
- [ ] **Docker** installed locally (for building images)
- [ ] **Node.js 20+** (for building the frontend)
- [ ] **A registered domain** (optional, for custom domain setup)
- [ ] **Git** with the project repository cloned

```bash
# Verify AWS CLI
aws --version
aws sts get-caller-identity

# Set your preferred region
export AWS_REGION=us-east-1
```

---

## Step 1: Set Up VPC and Security Groups

### 1.1 Create a VPC (or use the default VPC)

```bash
# Create a new VPC
aws ec2 create-vpc \
  --cidr-block 10.0.0.0/16 \
  --tag-specifications 'ResourceType=vpc,Tags=[{Key=Name,Value=rpa-vpc}]'
```

### 1.2 Create Subnets

```bash
# Public subnet (for EC2/ALB)
aws ec2 create-subnet \
  --vpc-id <vpc-id> \
  --cidr-block 10.0.1.0/24 \
  --availability-zone us-east-1a \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=rpa-public-1a}]'

# Public subnet (second AZ for ALB)
aws ec2 create-subnet \
  --vpc-id <vpc-id> \
  --cidr-block 10.0.2.0/24 \
  --availability-zone us-east-1b \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=rpa-public-1b}]'

# Private subnet (for RDS)
aws ec2 create-subnet \
  --vpc-id <vpc-id> \
  --cidr-block 10.0.3.0/24 \
  --availability-zone us-east-1a \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=rpa-private-1a}]'

# Private subnet (second AZ for RDS Multi-AZ)
aws ec2 create-subnet \
  --vpc-id <vpc-id> \
  --cidr-block 10.0.4.0/24 \
  --availability-zone us-east-1b \
  --tag-specifications 'ResourceType=subnet,Tags=[{Key=Name,Value=rpa-private-1b}]'
```

### 1.3 Create Security Groups

```bash
# Backend security group (allow HTTP 8080 from ALB)
aws ec2 create-security-group \
  --group-name rpa-backend-sg \
  --description "Security group for RPA backend" \
  --vpc-id <vpc-id>

aws ec2 authorize-security-group-ingress \
  --group-id <backend-sg-id> \
  --protocol tcp --port 8080 --cidr 10.0.0.0/16

# Database security group (allow PostgreSQL 5432 from backend)
aws ec2 create-security-group \
  --group-name rpa-database-sg \
  --description "Security group for RPA database" \
  --vpc-id <vpc-id>

aws ec2 authorize-security-group-ingress \
  --group-id <db-sg-id> \
  --protocol tcp --port 5432 --source-group <backend-sg-id>
```

---

## Step 2: Create RDS PostgreSQL Instance

### 2.1 Create a DB Subnet Group

```bash
aws rds create-db-subnet-group \
  --db-subnet-group-name rpa-db-subnet-group \
  --db-subnet-group-description "Subnets for RPA PostgreSQL" \
  --subnet-ids <private-subnet-1a> <private-subnet-1b>
```

### 2.2 Launch the RDS Instance

```bash
aws rds create-db-instance \
  --db-instance-identifier rpa-postgres \
  --db-instance-class db.t3.micro \
  --engine postgres \
  --engine-version 15.4 \
  --master-username postgres \
  --master-user-password <STRONG_PASSWORD> \
  --allocated-storage 20 \
  --storage-type gp3 \
  --db-name research_paper_analyzer \
  --vpc-security-group-ids <db-sg-id> \
  --db-subnet-group-name rpa-db-subnet-group \
  --backup-retention-period 7 \
  --no-publicly-accessible \
  --storage-encrypted \
  --tags Key=Project,Value=ResearchPaperAnalyzer
```

### 2.3 Wait for the Instance & Get the Endpoint

```bash
aws rds wait db-instance-available --db-instance-identifier rpa-postgres

aws rds describe-db-instances \
  --db-instance-identifier rpa-postgres \
  --query 'DBInstances[0].Endpoint.Address' \
  --output text
```

> Save the endpoint — you'll need it as `DB_HOST` in the backend configuration.

### 2.4 Initialize the Schema

```bash
# From a bastion host or EC2 instance in the same VPC:
psql -h <rds-endpoint> -U postgres -d research_paper_analyzer -f database/schema.sql
```

---

## Step 3: Create S3 Bucket for File Storage

### 3.1 Create the Bucket

```bash
aws s3api create-bucket \
  --bucket rpa-paper-uploads-<account-id> \
  --region us-east-1

# Enable versioning (recommended)
aws s3api put-bucket-versioning \
  --bucket rpa-paper-uploads-<account-id> \
  --versioning-configuration Status=Enabled
```

### 3.2 Block Public Access

```bash
aws s3api put-public-access-block \
  --bucket rpa-paper-uploads-<account-id> \
  --public-access-block-configuration \
    BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true
```

### 3.3 Set CORS (for direct uploads, if applicable)

```bash
aws s3api put-bucket-cors \
  --bucket rpa-paper-uploads-<account-id> \
  --cors-configuration '{
    "CORSRules": [{
      "AllowedOrigins": ["https://app.example.com"],
      "AllowedMethods": ["GET", "PUT"],
      "AllowedHeaders": ["*"],
      "MaxAgeSeconds": 3600
    }]
  }'
```

### 3.4 Set Lifecycle Policy (optional — auto-cleanup)

```bash
aws s3api put-bucket-lifecycle-configuration \
  --bucket rpa-paper-uploads-<account-id> \
  --lifecycle-configuration '{
    "Rules": [{
      "ID": "CleanupOldUploads",
      "Status": "Enabled",
      "Filter": {"Prefix": "uploads/"},
      "Transitions": [{
        "Days": 90,
        "StorageClass": "GLACIER"
      }],
      "Expiration": {"Days": 365}
    }]
  }'
```

---

## Step 4: Set Up IAM Roles and Policies

### 4.1 Create an IAM Policy for the Backend

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "S3Access",
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:GetObject",
        "s3:DeleteObject",
        "s3:ListBucket"
      ],
      "Resource": [
        "arn:aws:s3:::rpa-paper-uploads-<account-id>",
        "arn:aws:s3:::rpa-paper-uploads-<account-id>/*"
      ]
    },
    {
      "Sid": "BedrockAccess",
      "Effect": "Allow",
      "Action": [
        "bedrock:InvokeModel",
        "bedrock:InvokeModelWithResponseStream"
      ],
      "Resource": "arn:aws:bedrock:us-east-1::foundation-model/anthropic.claude-3-sonnet-20240229-v1:0"
    },
    {
      "Sid": "CloudWatchLogs",
      "Effect": "Allow",
      "Action": [
        "logs:CreateLogGroup",
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Resource": "arn:aws:logs:us-east-1:*:*"
    }
  ]
}
```

```bash
# Save the above as policy.json, then:
aws iam create-policy \
  --policy-name RPABackendPolicy \
  --policy-document file://policy.json
```

### 4.2 Create an IAM Role for EC2

```bash
# Create trust policy for EC2
cat > trust-policy.json << 'EOF'
{
  "Version": "2012-10-17",
  "Statement": [{
    "Effect": "Allow",
    "Principal": {"Service": "ec2.amazonaws.com"},
    "Action": "sts:AssumeRole"
  }]
}
EOF

aws iam create-role \
  --role-name RPABackendRole \
  --assume-role-policy-document file://trust-policy.json

aws iam attach-role-policy \
  --role-name RPABackendRole \
  --policy-arn arn:aws:iam::<account-id>:policy/RPABackendPolicy

# Create instance profile
aws iam create-instance-profile \
  --instance-profile-name RPABackendProfile

aws iam add-role-to-instance-profile \
  --instance-profile-name RPABackendProfile \
  --role-name RPABackendRole
```

---

## Step 5: Configure Amazon Bedrock Access

### 5.1 Enable Model Access

1. Open the **Amazon Bedrock** console.
2. Navigate to **Model access** in the left sidebar.
3. Click **Manage model access**.
4. Enable **Anthropic → Claude 3 Sonnet**.
5. Wait for the status to show **Access granted**.

### 5.2 Verify Access via CLI

```bash
aws bedrock list-foundation-models \
  --query "modelSummaries[?modelId=='anthropic.claude-3-sonnet-20240229-v1:0'].{Name:modelName,Status:modelLifecycle.status}" \
  --output table
```

### 5.3 Test with a Sample Invocation

```bash
aws bedrock-runtime invoke-model \
  --model-id anthropic.claude-3-sonnet-20240229-v1:0 \
  --content-type application/json \
  --body '{"anthropic_version":"bedrock-2023-05-31","max_tokens":100,"messages":[{"role":"user","content":"Hello!"}]}' \
  /tmp/bedrock-test-response.json

cat /tmp/bedrock-test-response.json
```

---

## Step 6: Deploy Backend

### Option A: EC2 with Docker

#### 6A.1 Launch an EC2 Instance

```bash
aws ec2 run-instances \
  --image-id ami-0c55b159cbfafe1f0 \
  --instance-type t3.small \
  --key-name <your-key-pair> \
  --security-group-ids <backend-sg-id> \
  --subnet-id <public-subnet-1a> \
  --iam-instance-profile Name=RPABackendProfile \
  --tag-specifications 'ResourceType=instance,Tags=[{Key=Name,Value=rpa-backend}]' \
  --user-data file://user-data.sh
```

#### 6A.2 User Data Script (`user-data.sh`)

```bash
#!/bin/bash
set -e

# Install Docker
yum update -y
yum install -y docker git
systemctl start docker
systemctl enable docker
usermod -aG docker ec2-user

# Install Docker Compose
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" \
  -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

# Clone the repository
cd /home/ec2-user
git clone https://github.com/your-username/Cloud-based-AI-Powered-Research-Paper-Analyzer.git app
cd app

# Create environment file
cat > .env << 'ENVEOF'
DB_HOST=<rds-endpoint>
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=<rds-password>
DB_NAME=research_paper_analyzer
JWT_SECRET=<strong-random-secret>
AI_PROVIDER=bedrock
STORAGE_PROVIDER=s3
AWS_REGION=us-east-1
S3_BUCKET_NAME=rpa-paper-uploads-<account-id>
BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0
PORT=8080
CORS_ALLOWED_ORIGINS=https://app.example.com
GIN_MODE=release
SEED_DATA=false
ENVEOF

# Build and run backend only (database is on RDS)
cd backend
docker build -t rpa-backend .
docker run -d \
  --name rpa-backend \
  --env-file ../.env \
  -p 8080:8080 \
  --restart unless-stopped \
  rpa-backend
```

#### 6A.3 Set Up an Application Load Balancer

```bash
# Create ALB
aws elbv2 create-load-balancer \
  --name rpa-alb \
  --subnets <public-subnet-1a> <public-subnet-1b> \
  --security-groups <alb-sg-id> \
  --scheme internet-facing

# Create target group
aws elbv2 create-target-group \
  --name rpa-backend-tg \
  --protocol HTTP \
  --port 8080 \
  --vpc-id <vpc-id> \
  --health-check-path /api/health \
  --health-check-interval-seconds 30

# Register EC2 instance
aws elbv2 register-targets \
  --target-group-arn <target-group-arn> \
  --targets Id=<instance-id>

# Create listener (HTTPS — requires ACM certificate)
aws elbv2 create-listener \
  --load-balancer-arn <alb-arn> \
  --protocol HTTPS --port 443 \
  --certificates CertificateArn=<acm-cert-arn> \
  --default-actions Type=forward,TargetGroupArn=<target-group-arn>
```

---

### Option B: Elastic Beanstalk

#### 6B.1 Initialize Elastic Beanstalk

```bash
cd backend

# Install EB CLI
pip install awsebcli

# Initialize
eb init rpa-backend \
  --platform "Docker" \
  --region us-east-1

# Create environment
eb create rpa-backend-prod \
  --instance-type t3.small \
  --single \
  --envvars \
    DB_HOST=<rds-endpoint>,\
    DB_PORT=5432,\
    DB_USER=postgres,\
    DB_PASSWORD=<rds-password>,\
    DB_NAME=research_paper_analyzer,\
    JWT_SECRET=<strong-secret>,\
    AI_PROVIDER=bedrock,\
    STORAGE_PROVIDER=s3,\
    S3_BUCKET_NAME=rpa-paper-uploads-<account-id>,\
    BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0,\
    PORT=8080,\
    CORS_ALLOWED_ORIGINS=https://app.example.com,\
    GIN_MODE=release,\
    SEED_DATA=false
```

#### 6B.2 Deploy Updates

```bash
eb deploy
eb status
eb logs
```

---

## Step 7: Deploy Frontend to S3 + CloudFront

### 7.1 Build the Frontend

```bash
cd frontend

# Set the backend API URL
echo "VITE_API_URL=https://api.example.com" > .env.production

# Build
npm ci
npm run build
```

### 7.2 Create S3 Bucket for Hosting

```bash
aws s3api create-bucket \
  --bucket rpa-frontend-<account-id> \
  --region us-east-1

# Configure for static website hosting
aws s3 website s3://rpa-frontend-<account-id>/ \
  --index-document index.html \
  --error-document index.html   # SPA fallback
```

### 7.3 Upload Build Files

```bash
aws s3 sync dist/ s3://rpa-frontend-<account-id>/ \
  --delete \
  --cache-control "public, max-age=31536000, immutable" \
  --exclude "index.html"

# Upload index.html with no-cache
aws s3 cp dist/index.html s3://rpa-frontend-<account-id>/index.html \
  --cache-control "no-cache, no-store, must-revalidate"
```

### 7.4 Create CloudFront Distribution

```bash
aws cloudfront create-distribution \
  --distribution-config '{
    "CallerReference": "rpa-frontend-dist",
    "Origins": {
      "Quantity": 1,
      "Items": [{
        "Id": "S3-rpa-frontend",
        "DomainName": "rpa-frontend-<account-id>.s3.amazonaws.com",
        "S3OriginConfig": {
          "OriginAccessIdentity": ""
        }
      }]
    },
    "DefaultCacheBehavior": {
      "TargetOriginId": "S3-rpa-frontend",
      "ViewerProtocolPolicy": "redirect-to-https",
      "AllowedMethods": {
        "Quantity": 2,
        "Items": ["GET", "HEAD"]
      },
      "ForwardedValues": {
        "QueryString": false,
        "Cookies": {"Forward": "none"}
      },
      "MinTTL": 0,
      "DefaultTTL": 86400,
      "MaxTTL": 31536000
    },
    "DefaultRootObject": "index.html",
    "CustomErrorResponses": {
      "Quantity": 1,
      "Items": [{
        "ErrorCode": 404,
        "ResponsePagePath": "/index.html",
        "ResponseCode": "200",
        "ErrorCachingMinTTL": 300
      }]
    },
    "Enabled": true,
    "Comment": "RPA Frontend Distribution"
  }'
```

### 7.5 Invalidate Cache After Deployment

```bash
aws cloudfront create-invalidation \
  --distribution-id <distribution-id> \
  --paths "/*"
```

---

## Step 8: Set Up CloudWatch Monitoring

### 8.1 Create a Dashboard

```bash
aws cloudwatch put-dashboard \
  --dashboard-name RPAMonitoring \
  --dashboard-body '{
    "widgets": [
      {
        "type": "metric",
        "properties": {
          "title": "Backend CPU",
          "metrics": [["AWS/EC2", "CPUUtilization", "InstanceId", "<instance-id>"]],
          "period": 300
        }
      },
      {
        "type": "metric",
        "properties": {
          "title": "RDS Connections",
          "metrics": [["AWS/RDS", "DatabaseConnections", "DBInstanceIdentifier", "rpa-postgres"]],
          "period": 300
        }
      },
      {
        "type": "metric",
        "properties": {
          "title": "ALB Request Count",
          "metrics": [["AWS/ApplicationELB", "RequestCount", "LoadBalancer", "<alb-id>"]],
          "period": 60,
          "stat": "Sum"
        }
      }
    ]
  }'
```

### 8.2 Create Alarms

```bash
# High CPU alarm
aws cloudwatch put-metric-alarm \
  --alarm-name rpa-backend-high-cpu \
  --metric-name CPUUtilization \
  --namespace AWS/EC2 \
  --statistic Average \
  --period 300 \
  --threshold 80 \
  --comparison-operator GreaterThanThreshold \
  --evaluation-periods 2 \
  --dimensions Name=InstanceId,Value=<instance-id> \
  --alarm-actions <sns-topic-arn>

# RDS storage alarm
aws cloudwatch put-metric-alarm \
  --alarm-name rpa-rds-low-storage \
  --metric-name FreeStorageSpace \
  --namespace AWS/RDS \
  --statistic Average \
  --period 300 \
  --threshold 2000000000 \
  --comparison-operator LessThanThreshold \
  --evaluation-periods 1 \
  --dimensions Name=DBInstanceIdentifier,Value=rpa-postgres \
  --alarm-actions <sns-topic-arn>
```

---

## Step 9: Optional — SNS Notifications

### 9.1 Create an SNS Topic

```bash
aws sns create-topic --name rpa-alerts

# Subscribe your email
aws sns subscribe \
  --topic-arn arn:aws:sns:us-east-1:<account-id>:rpa-alerts \
  --protocol email \
  --notification-endpoint your-email@example.com
```

> Confirm the subscription via the email you receive.

### 9.2 Connect to CloudWatch Alarms

Use the SNS topic ARN as the `--alarm-actions` parameter in the CloudWatch alarm commands above.

---

## Step 10: Optional — API Gateway + Lambda

For serverless PDF processing, you can offload heavy analysis to Lambda:

### 10.1 Create a Lambda Function

```bash
# Package the PDF processor
cd backend
GOOS=linux GOARCH=amd64 go build -o bootstrap cmd/lambda/main.go
zip function.zip bootstrap

aws lambda create-function \
  --function-name rpa-pdf-processor \
  --runtime provided.al2023 \
  --handler bootstrap \
  --zip-file fileb://function.zip \
  --role arn:aws:iam::<account-id>:role/RPABackendRole \
  --timeout 120 \
  --memory-size 512 \
  --environment Variables='{
    DB_HOST=<rds-endpoint>,
    DB_USER=postgres,
    DB_PASSWORD=<password>,
    DB_NAME=research_paper_analyzer,
    BEDROCK_MODEL_ID=anthropic.claude-3-sonnet-20240229-v1:0
  }'
```

### 10.2 Trigger from S3 Upload

```bash
aws s3api put-bucket-notification-configuration \
  --bucket rpa-paper-uploads-<account-id> \
  --notification-configuration '{
    "LambdaFunctionConfigurations": [{
      "LambdaFunctionArn": "arn:aws:lambda:us-east-1:<account-id>:function:rpa-pdf-processor",
      "Events": ["s3:ObjectCreated:*"],
      "Filter": {
        "Key": {
          "FilterRules": [{"Name": "suffix", "Value": ".pdf"}]
        }
      }
    }]
  }'
```

---

## Cost Estimation

Estimated monthly costs for a low-traffic deployment:

| Service              | Spec                    | Estimated Cost (USD/month) |
| -------------------- | ----------------------- | -------------------------: |
| EC2 (backend)        | t3.small (on-demand)    |                     ~$15   |
| RDS PostgreSQL       | db.t3.micro, 20 GB      |                     ~$15   |
| S3 (uploads)         | 10 GB storage            |                      ~$1   |
| S3 (frontend)        | Static hosting           |                      ~$1   |
| CloudFront           | 10 GB transfer           |                      ~$1   |
| Amazon Bedrock       | ~1000 requests/month     |                     ~$10   |
| CloudWatch           | Basic monitoring         |                      ~$0   |
| **Total**            |                         |                   **~$43** |

> **Cost Optimization Tips:**
> - Use **EC2 Reserved Instances** or **Savings Plans** for 40-60% savings.
> - Use **t3.micro** for development/staging.
> - Enable **S3 Intelligent Tiering** for automatic storage optimization.
> - Set up **budget alerts** via AWS Budgets.

---

## Security Best Practices

### ✅ Checklist

- [ ] **Never commit secrets** — Use `.env` files or AWS Secrets Manager.
- [ ] **Enable MFA** on the root AWS account.
- [ ] **Use IAM roles** (not access keys) for EC2 instances.
- [ ] **Encrypt at rest** — RDS encryption, S3 default encryption.
- [ ] **Encrypt in transit** — HTTPS everywhere via ACM + CloudFront.
- [ ] **Restrict security groups** — Allow only required ports from known sources.
- [ ] **Enable VPC Flow Logs** for network monitoring.
- [ ] **Rotate credentials** regularly (RDS password, JWT secret).
- [ ] **Enable RDS automated backups** (7-day retention minimum).
- [ ] **Enable S3 versioning** to protect against accidental deletion.
- [ ] **Use WAF** on CloudFront and ALB for DDoS protection.
- [ ] **Set up AWS Config** rules for compliance monitoring.

### Secrets Manager (Recommended)

Instead of environment variables, use AWS Secrets Manager:

```bash
aws secretsmanager create-secret \
  --name rpa/backend/config \
  --secret-string '{
    "DB_PASSWORD": "<strong-password>",
    "JWT_SECRET": "<random-secret>"
  }'
```

Then retrieve secrets in your application at startup using the AWS SDK.

---

## Troubleshooting

### Backend won't start

| Symptom                           | Solution                                                        |
| --------------------------------- | --------------------------------------------------------------- |
| `connection refused` to database  | Check security group allows 5432 from backend SG                |
| `permission denied` for Bedrock   | Verify IAM role has `bedrock:InvokeModel` permission            |
| `no such host` for RDS endpoint   | Ensure EC2 is in the same VPC as RDS                            |
| Container exits immediately       | Check logs: `docker logs rpa-backend`                           |

### Frontend shows blank page

| Symptom                           | Solution                                                        |
| --------------------------------- | --------------------------------------------------------------- |
| 403 from CloudFront               | Check S3 bucket policy allows CloudFront OAI                    |
| API calls fail (CORS)             | Add CloudFront domain to `ALLOWED_ORIGINS`                      |
| Routes return 404                 | Ensure CloudFront custom error response returns `/index.html`   |

### Database issues

| Symptom                           | Solution                                                        |
| --------------------------------- | --------------------------------------------------------------- |
| `too many connections`            | Increase `max_connections` in RDS parameter group               |
| Slow queries                      | Check indexes; enable `pg_stat_statements`                      |
| Can't connect from local machine  | RDS is in private subnet — use a bastion host or VPN            |

### Bedrock / AI issues

| Symptom                           | Solution                                                        |
| --------------------------------- | --------------------------------------------------------------- |
| `Access denied` on InvokeModel    | Enable model access in Bedrock console; check IAM policy        |
| Timeout on large papers           | Increase Lambda/backend timeout; chunk the paper text            |
| Model returns empty response      | Check prompt format matches Anthropic's API spec                 |

---

<p align="center">
  <em>For architecture details, see <a href="ARCHITECTURE.md">ARCHITECTURE.md</a> · For local setup, see <a href="README.md">README.md</a></em>
</p>
