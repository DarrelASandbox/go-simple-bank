<details>
  <summary>Table of Contents</summary>
  <ul>
    <li><a href="#about-the-project">About The Project</a></li>
    <li><a href="#sql">SQL</a></li>
    <li><a href="#deadlock">Deadlock</a></li>
    <li><a href="#transaction-isolation-level">Transaction Isolation Level</a></li>
    <li><a href="#github-actions">GitHub Actions</a></li>
    <li><a href="#mockgen">mockgen</a></li>
    <li><a href="#migration">Migration</a></li>
    <li><a href="#paseto">PASETO</a></li>
    <li><a href="#docker">Docker</a></li>
    <li><a href="#aws-ecr-with-github-actions">AWS ECR with GitHub Actions</a></li>
    <li><a href="#aws-rds">AWS RDS</a></li>
    <li><a href="#aws-secrets-manager">AWS Secrets Manager</a></li>
    <li><a href="#aws-eks">AWS EKS</a></li>
    <li><a href="#kubectl--k9s">kubectl & k9s</a></li>
    <li><a href="#k8s-with-aws-eks">k8s with AWS EKS</a></li>
    <li><a href="#aws-route-53">AWS Route 53</a></li>
    <li><a href="#k8s-ingress">k8s Ingress</a></li>
    <li><a href="#k8s-cert-manager--lets-encrypt">k8s cert-manager & Let's Encrypt</a></li>
    <li><a href="#k8s-with-github-actions">k8s with GitHub Actions</a></li>
    <li><a href="#dbdocsio">dbdocs.io</a></li>
    <li><a href="#grpc">gRPC</a></li>
  </ul>
</details>

&nbsp;

## About The Project

- Backend Master Class [Golang + Postgres + Kubernetes + gRPC]
- Learn everything about backend web development: Golang, Postgres, Gin, gRPC, Docker, Kubernetes, AWS, GitHub Actions
- [Original Repo - simple_bank](https://github.com/techschool/simplebank)
- [techschool](https://github.com/techschool)

&nbsp;

---

&nbsp;

## SQL

- [dbdiagram.io](https://dbdiagram.io/home)
- **CRUD**
  - **SQL**
    - Very fast & straightforward
    - Manual mapping SQL fields to variables
  - **GORM**
    - CRUD functions already implemented
    - Run slowly on high load
  - **SQLX**
    - Quite fast & easy to use
    - Fields mapping via query text & struct tags
  - [**SQLC**](https://sqlc.dev/)
    - Very fast & easy to use
    - Automatic code generation
    - Catch SQL query errors before generating codes
    - Full support Postgres. MySQL is experimental

```sh
# Login using root
psql simple_bank -U root
# Check connection info
\conninfo

# In db folder
sqlc generate
```

```sql
-- Set idle session limit as superuser
ALTER system SET idle_in_transaction_session_timeout='5min';
-- Disable idle session limit as superuser
ALTER system SET idle_in_transaction_session_timeout=0;
```

&nbsp;

---

&nbsp;

## Deadlock

- **store_test.go: Logs for deadlock bug**
- Refer to commit `45de7cb5a930e3bcdddae513d968b4943327983e`

```
=== RUN   TestTransferTx
>> before: 807 9
tx 2 create transfer
tx 2 create entry 1
tx 2 create entry 2
tx 2 get account 1
tx 2 update account 1
tx 2 get account 2
tx 2 update account 2
tx 1 create transfer
>> tx: 797 19
tx 1 create entry 1
tx 1 create entry 2
tx 1 get account 1
tx 1 update account 1
tx 1 get account 2
tx 1 update account 2
>> tx: 787 29
>> after: 787 29
--- PASS: TestTransferTx (0.21s)
PASS
```

```sql
BEGIN;

INSERT INTO transfers (from_account_id, to_account_id, amount) VALUES (1, 2, 10) RETURNING *;

INSERT INTO entries (account_id, amount) VALUES (1, -10) RETURNING *;
INSERT INTO entries (account_id, amount) VALUES (2, 10) RETURNING *;

SELECT * FROM accounts WHERE id = 1 FOR UPDATE;
UPDATE accounts SET balance = 90 WHERE id = 1 RETURNING *;

SELECT * FROM accounts WHERE id = 2 FOR UPDATE;
UPDATE accounts SET balance = 110 WHERE id = 2 RETURNING *;

ROLLBACK
```

- [postgresql - Lock_Monitoring](https://wiki.postgresql.org/wiki/Lock_Monitoring)

```sql
-- edited code snippets from link above
-- see ShareLock under mode
SELECT a.application_name,
       l.relation::regclass,
       l.transactionid,
       l.mode,
	     l.locktype,
       l.GRANTED,
       a.usename,
       a.query,
       a.pid
FROM pg_stat_activity a
JOIN pg_locks l ON l.pid = a.pid
ORDER BY a.pid;
```

&nbsp;

---

&nbsp;

- Run the 4 queries in the order below to simulate deadlock

```sql
-- Tx1: transfer $10 from account 1 to account 2
BEGIN;

-- 1
UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;
-- 3
UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;

ROLLBACK;

-- Tx2: transfer $10 from account 2 to account 1
BEGIN;

-- 2
UPDATE accounts SET balance = balance - 10 WHERE id = 2 RETURNING *;
-- 4
UPDATE accounts SET balance = balance - 10 WHERE id = 1 RETURNING *;

ROLLBACK;
```

&nbsp;

---

&nbsp;

## Transaction Isolation Level

- **Read Phenomena**
  - **Dirty Read:** A transaction reads data written by other concurrent uncommitted transaction
  - **Non-Repeatable Read:** A transaction reads the same row twice and sees different value because it has been modified by other **committed** transaction
  - **Phantom Read:** A transaction re-executes a query to **find rows** that satisfy a condition and sees a **different set** of rows, due to changes by other **committed** transaction
  - **Serialization Anomaly:** The result of a **group** of concurrent **committed transactions** is **impossible to achieve** if we try to run them **sequentially** in any order without overlapping
- **4 Standard Isolation Levels**
  - **Read Uncommitted:** Can see data written by uncommitted transaction
  - **Read Committed:** Only see data written by committed transaction
  - **Repeatable Read:** Same read query always returns same result
  - **Serializable:** Can achieve same result if execute transactions serially in some order instead of concurrently

|                       | Read Uncommitted | Read Committed | Repeatable Read | Serializable |
| :-------------------: | :--------------: | :------------: | :-------------: | :----------: |
|      Dirty Read       |        ✅        |       ❌       |       ❌        |      ❌      |
|  Non-Repeatable Read  |        ✅        |       ✅       |       ❌        |      ❌      |
|     Phantom Read      |        ✅        |       ✅       |       ❌        |      ❌      |
| Serialization Anomaly |        ✅        |       ✅       |       ✅        |      ❌      |

|       MySQL        |        Postgres        |
| :----------------: | :--------------------: |
| 4 Isolation Levels |   3 Isolation Levels   |
| Locking Mechanism  | Dependencies Detection |
|  Repeatable Read   |     Read Committed     |

- High Level Isolation Methods
  - Retry Mechanism: There might be errors, timeout or deadlock
  - Read documentation: Each database engine might implement isolation level differently

&nbsp;

---

&nbsp;

## GitHub Actions

1. Actions -> Configure in Go -> Copy code to vscode
2. [GitHub - Creating PostgreSQL service containers](https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers)
3. [GitHub - migrate CLI usage](https://github.com/golang-migrate/migrate#cli-usage)
4. [GitHub - migrate releases](https://github.com/golang-migrate/migrate/releases)
   1. [migrate.linux-amd64.tar.gz](https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz)

&nbsp;

---

&nbsp;

## mockgen

- [Go Package - gomock](https://pkg.go.dev/github.com/golang/mock/gomock)

```sh
# Go 1.16+
go install github.com/golang/mock/mockgen@v1.6.0
mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/DarrelASandbox/go-simple-bank/db/sqlc Store
```

&nbsp;

---

&nbsp;

## Migration

- [Go Package - migrate: CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

```sh
migrate create -ext sql -dir db/migration -seq add_users
```

&nbsp;

---

&nbsp;

## PASETO

- Platform-Agnostic Security Tokens
- Only 2 most recent PASTO versions are accepted
- **Non-trivial Forgery**
  - No more "alg" header or "none" algorithm
  - Everything is authenticated
  - Encrypted payload for local use &lt;symmetric key&gt;

&nbsp;

---

&nbsp;

## Docker

- [sh-compatible wait-for](https://github.com/Eficode/wait-for)
  - Download from release and place it in root folder (`wait-for.sh`)

```sh
docker build -t simple_bank:latest .
docker run --name simple_bank -p 4000:4000 -e GIN_MODE=release simple_bank:latest
docker compose up

# Make script executable
chmod +x start.sh
chmod +x wait-for.sh
```

&nbsp;

---

&nbsp;

## AWS ECR with GitHub Actions

1. AWS ECR -> Create repository for URI
2. Instead of using the commands given by AWS (View push commands), we will use GitHub Actions
3. [GitHub Marketplace Actions](https://github.com/marketplace?category=&query=&type=actions&verification=)
   1. Search for `ECR`
   2. Refer to `deploy.yml` file
4. AWS IAM -> Add user
   1. **User name:** github-ci
   2. **Access type:** Programmatic access
5. Next -> Create group
   1. **Group name:** deployment
   2. **Filter policies:** elastic container registry
   3. Check AmazonEC2ContainerRegistryFullAccess -> Create group -> Next -> Review
6. GitHub simple_bank repo Settings -> Secrets -> New repository secret
   1. Refer to the code snippets below
   2. **Name:** `AWS_ACCESS_KEY_ID`
   3. **Value:** `AWS Access Key ID`
   4. **Name:** `AWS_SECRET_ACCESS_KEY`
   5. **Value:** `AWS Secret access key`

```yml
# From tutorial video
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v1
  with:
    aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
    aws-secret-access-key: ${{secrets.AWS_SECRET_ACCESS_KEY }}

# From deploy.yml
- name: Configure AWS credentials
  uses: aws-actions/configure-aws-credentials@v1
  with:
    role-to-assume: arn:aws:iam::123456789012:role/my-github-actions-role
```

&nbsp;

---

&nbsp;

## AWS RDS

1. AWS RDS -> Create database
2. **Choose a database creation method:** Standard create
3. **Engine options:** PostgreSQL
4. **DB instance identifier:** simple_bank
5. **Master username:** root
6. Check **Auto generate a password**
7. **DB instance class:** db.t2.micro
8. Uncheck **Storage autoscaling**
9. Use Default VPC for **Connectivity**
10. **Public access:** Yes
11. **VPC security group:**
    1. Create new
    2. **New VPC security group name:** access-postgres-anywhere
12. **Database authentication:** Password authentication
13. **Initial database name:** simple_bank
14. Create database
15. View credential details
16. Register new Server/Connection for PostgreSQL
    1. **Name:** AWS Postgres
    2. **Host:** From AWS connection details endpoint _(See point 19)_
    3. **Username:** root
    4. **Password:** From AWS RDS credential details
    5. **Database:** simple_bank
17. Check VPC Security Groups from AWS RDS
18. Via the Security group ID -> Edit inbound rules
    1. **Source:** Anywhere _(Unless ip is static)_
    2. Save rules
19. View connection details from AWS RDS -> Copy Endpoint
20. Update `Makefile` postgres url

&nbsp;

---

&nbsp;

## AWS Secrets Manager

```sh
# Generate a 32 characters string
openssl rand -hex 64 | head -c 32
```

1. AWS Secrets Manager -> Store a new secret
2. Other type of secrets
   1. **DB_DRIVER:** postgres
   2. **DB_SOURCE:** _Use AWS credentials for the postgres connection_
   3. **SERVER_ADDRESS:** 0.0.0.0:4000
   4. **TOKEN_SYMMETRIC_KEY:** _Generate a 32 characters string_
   5. **ACCESS_TOKEN_DURATION:** 15m
3. Next -> **Secret name:** simple_bank -> Next -> Next -> Store
4. [Install AWS CLI](https://aws.amazon.com/cli/)
   1. `aws configure` _(Refer to IAM profile)_
   2. `ls -l ~/.aws`
   3. `cat ~/.aws/credentials`
   4. `aws secretsmanager help`
   5. `aws secretsmanager get-secret-value --secret-id simple_bank` _Might want to try with arn as well_
   6. AWS IAM User groups -> Permissions -> Add permissions -> **Attach Policies:** SecretsManagerReadWrite
   7. `aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text`
5. `brew install jq` (To output into json format)
   1. `aws secretsmanager get-secret-value --secret-id simple_bank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env`
      1. [jq - String interpolation](<https://stedolan.github.io/jq/manual/#Stringinterpolation-(foo)>)
      2. [jq - Array/Object Value Iterator](https://stedolan.github.io/jq/manual/#Array/ObjectValueIterator:.[])
      3. [jq - Invoking jq (`--raw-output / -r`)](https://stedolan.github.io/jq/manual/#Invokingjq)
6. After GitHub Actions are completed, check AWS ECR for the newly built image
   1. [AWS CLI Command Reference: get-login-password](https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login-password.html)
   2. `aws ecr get-login-password | docker login --username AWS --password-stdin <aws_account_id>.dkr.ecr.<region>.amazonaws.com`
   3. `docker pull <AWS ECR Image URI>`

&nbsp;

---

&nbsp;

## AWS EKS

1. AWS EKS -> clusters -> Create cluster
   1. **Name:** simple_bank
   2. **Kubernetes version:** 1.20
   3. **Cluster Service Role:** AWSEKSClusterRole
      1. AWS IAM -> Create role -> `EKS - Cluster` use case -> Next (`AmazonEKSClusterPolicy`) -> Next
      2. **Role name:** AWSEKSClusterRole
2. Next -> Use default VPC
3. **Cluster endpoint access:** Public and private
4. **Amazon VPC CNI Version:** `v1.8.0-eksbuild. 1`
5. Next -> Next -> Create -> Refresh later when cluster is created
6. Add Node Group
   1. **Name:** simple_bank
   2. **Node IAM Role:** AWSEKSNodeRole
      1. AWS IAM -> Create role -> `EC2` use case -> Next
      2. Pick `AmazonEKS_CNI_Policy` & `AmazonEKSWorkerNodePolicy` & `AmazonEC2ContainerRegistryReadOnly` -> Next
      3. **Role name:** AWSEKSNodeRole
      4. Next
      5. **AMI type:** Amazon Linux 2 (AL2_x86_64)
      6. **Capacity type:** On-Demand
      7. **Instance types:** t3.small
      8. **Disk size:** 10 GiB
      9. **Node Group scaling configuration**
         1. **Minimum size:** 1 nodes
         2. **Maximum size:** 2 nodes
         3. **Desired size:** 1 nodes
      10. Next -> Disable **Allow remote access to nodes** -> Next -> Create
7. Refresh later when node is created

&nbsp;

---

&nbsp;

## kubectl & k9s

1. [Command line tool (kubectl)](https://kubernetes.io/docs/reference/kubectl/)
2. `kubectl cluster-info`
3. AWS IAM -> Users -> Groups -> Pick `deployment` group that was created before
4. Permissions -> Add permissions -> Create inline policy
   1. **Service:** EKS
   2. **Actions:** All EKS actions
   3. **Resources:** All resources
5. Review policy
6. **Name:** EKSFullAccess
7. Create policy
8. `aws eks update-kubeconfig --name simple_bank --region us-east-1`
9. `ls -l ~/.kube`
10. **Switch between clusters:** `kubectl config use-context arn:aws:eks:ap-southeast-1:560476749134:cluster/dep-aws-eks`
11. [How do I provide access to other IAM users and roles after cluster creation in Amazon EKS?](https://aws.amazon.com/premiumsupport/knowledge-center/amazon-eks-cluster-access/)
    1. AWS Profile -> My Security Credentials -> Access Keys (access key ID and secret access key) -> Create New Access Key
    2. `vi ~/.aws/credentials` to add the new credentials
    3. `export AWS_PROFILE=github` or `export AWS_PROFILE=default`
    4. Input User ARN from AWS IAM Users into `aws-auth.yaml`
    5. `kubectl apply -f eks/aws-auth.yaml`
12. [k9s](https://k9scli.io/)
    1. [Commands](https://k9scli.io/topics/commands/)

&nbsp;

---

&nbsp;

## k8s with AWS EKS

1. [Kubernetes Documentation - Concepts - Workloads - Workload Resources - Deployments](https://kubernetes.io/docs/concepts/workloads/controllers/deployment/)
2. `kubectl apply -f eks/deployment.yaml`
3. AWS EKS -> Clusters -> Compute -> `simple_bank` Node Group -> Autoscaling group name -> Edit group details
   1. **Desired capacity:** 1
   2. Update -> Activity -> Refresh
   3. **The points below are for switching instance type:**
      1. [AWS EC2 - IP addresses per network interface per instance type: t3.micro](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-eni.html)
      2. Edit -> Edit Node Group: simple_bank page -> **Cannot change eni**
      3. So we must **delete** node group
      4. Add Node Group
      5. **Name:** simple_bank
      6. **Role name:** AWSEKSNodeRole
      7. Next
      8. **AMI type:** Amazon Linux 2 (AL2_x86_64)
      9. **Capacity type:** On-Demand
      10. **Instance types:** t3.small
      11. **Disk size:** 10 GiB
      12. **Node Group scaling configuration**
          1. **Minimum size:** 0 nodes
          2. **Maximum size:** 2 nodes
          3. **Desired size:** 1 nodes
      13. Next -> Next -> Create
4. [Kubernetes Documentation - Concepts - Services - Load Balancing, and Networking - Service](https://kubernetes.io/docs/concepts/services-networking/service/)
5. `kubectl apply -f eks/service.yaml`
6. [Linux nslookup command](https://www.computerhope.com/unix/unslooku.htm)

&nbsp;

---

&nbsp;

## AWS route 53

1. [Amazon Route 53 pricing](https://aws.amazon.com/route53/pricing/)
2. AWS Route 53 -> Register domain (e.g. simple-bank.org) -> Check -> Continue
3. Fill in **Registrant Contact** form -> Continue -> Email Verification
4. Automatically Renew Domain -> Terms and Conditions -> Complete Order
5. Enable Transfer lock (Confirm in Registered domains page)
6. Hosted zones -> domain name
   1. **Name Server (NS) record**
   2. **Start of Authority (SOA) record**
7. Create record:
   1. **Record Name:** api.domain-name.com
   2. **Record Type:** Address (A) Record
   3. **Value:** Check `Alias`
      1. Alias to Network Load Balancer
      2. **Region of Load Balancer:** us-east-1
      3. **URL of Load Balancer from API Service**
8. `nslookup api.simple-bank.org`

&nbsp;

---

&nbsp;

## k8s Ingress

- [Kubernetes Documentation - Concepts - Services, Load Balancing, and Networking - Ingress #hostname-wildcards](https://kubernetes.io/docs/concepts/services-networking/ingress/#hostname-wildcards)
- [Kubernetes Documentation - Concepts - Services, Load Balancing, and Networking - Ingress Controllers](https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/)
- [Ingress NGINX Controller](https://github.com/kubernetes/ingress-nginx/blob/main/README.md#readme)
- [Kubernetes Documentation - Concepts - Services, Load Balancing, and Networking - Ingress #ingress-class](https://kubernetes.io/docs/concepts/services-networking/ingress/#ingress-class)

1. `kubectl apply -f eks/service.yaml` (Change from LoadBalancer to ClusterIP because we do not want to expose it outside world)
2. `kubectl apply -f eks/ingress.yaml`
3. `kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.5.1/deploy/static/provider/aws/deploy.yaml`
4. Copy `ingress-nginx` namespace's address from k9s -> AWS Route 53 Hosted zones page -> Check `api.simple-bank.org` Record name -> Edit record
   1. Replace the URL of Load Balancer with the copied address
   2. Save
5. `nslookup api.simple-bank.org`
6. `kubectl apply -f eks/ingress.yaml` (Apply #ingress-class settings)

&nbsp;

---

&nbsp;

## k8s cert-manager & Let's Encrypt

- [Kubernetes Documentation - Concepts - Services, Load Balancing, and Networking - Ingress #tls](https://kubernetes.io/docs/concepts/services-networking/ingress/#tls)
- [cert-manager](https://cert-manager.io/docs/)
- [Let's Encrypt](https://letsencrypt.org/docs/)
- [cert-manager - Automated Certificate Management Environment (ACME)](https://cert-manager.io/docs/configuration/acme/)
- `kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.10.0/cert-manager.yaml`
- Check `cert-manager` namespace
- `kubectl apply -f eks/issuer.yaml`
- Search for `clusterissuer` in k9s
- Search for `secrets` in k9s
- Search for `certificates` in k9s
- Check `all+` namespace then search for `ingress`
- Describe simple-bank-ingress

&nbsp;

---

&nbsp;

## k8s with GitHub Actions

- [GitHub Marketplace - Kubectl tool installer](https://github.com/marketplace/actions/kubectl-tool-installer)
  - Find latest version using this URL: `storage.googleapis.com/kubernetes-release/release/stable.txt`

&nbsp;

---

&nbsp;

## dbdocs.io

- [dbdocs.io](https://dbdocs.io/docs)
  - Refer to `db.dbml`
  - Login -> `dbdocs build doc/db.dbml`

&nbsp;

---

&nbsp;

## gRPC

- **Remote Procedure Call Framework**
  - The client can execute a remote procedure on the server
  - The remote interaction code is handled by gRPC
  - The API & data structure code is automatically generated
  - Support multiple programming languages
- **How it works?**
  1. **Define API & data structure**
     - The RPC and its request/ response structure are defined using protobuf
  2. **Generate gRPC stubs**
     - Generate codes for the server and client in the language of your choice
  3. **Implement the server**
     - Implement the RPC handler on the server side
  4. **Use the client**
     - Use the generated client stubs to call the RPC on the serser
- **Why GRPC?**
  - **High performance:** HTTP/2: binary framing, multiplexing, header compression, bidirectional communication
  - **Strong API contract:** Server & client share the same protobuf RPC definition with strongly typed data
  - **Automatic code generation:** Codes that serialize/ deserialize data, or transfer data between client & server are automatically generated
- **4 Types of GRPC**
  1. Unary gRPC
  2. Client streaming gRPC
  3. Server streaming gRPC
  4. Bidirectional streaming gRPC
- [**gRPC Gateway:**](https://github.com/grpc-ecosystem/grpc-gateway) Serve both gRPC and HTTP requests at once
  - A plugin of protobuf compiler
  - Generate proxy codes from protobuf
    - In-process translation: only for unary
    - Separate proxy server: both unary and streaming
  - Write code once, serve both gRPC and HTTP requests
  - `protoc-gen-grpc-gateway --help`
  - [gRPC-Gateway - Using proto names in JSON](https://grpc-ecosystem.github.io/grpc-gateway/docs/mapping/customizing_your_gateway/#using-proto-names-in-json)
- [gRPC - Introduction to gRPC](https://grpc.io/docs/what-is-grpc/introduction/)
- [gRPC - Docs - Languages - Go - Quick start](https://grpc.io/docs/languages/go/quickstart/)
- `brew install protobuf`
- `protoc --version`

```json
  "protoc": {
    "options": ["--proto_path=protos/v3"]
  }
```

- **`api` folder: REST api - HTTP request**
- **`api_rpc` folder: gRPC api - grpc message**
- [Why do we need to register reflection service on gRPC server](https://stackoverflow.com/questions/41424630/why-do-we-need-to-register-reflection-service-on-grpc-server)
- [List of TCP and UDP port numbers](https://en.wikipedia.org/wiki/List_of_TCP_and_UDP_port_numbers)

```sh
evans --host localhost --port=5000 -r repl
show service
call CreateUser
exit
```

&nbsp;

---

&nbsp;

## Swagger

- [Swagger](https://swagger.io/)
- [Swagger UI](https://github.com/swagger-api/swagger-ui)

```sh
# From source `~/Downloads/grpc-gateway/protoc-gen-openapiv2/options`
# Copy required files
cp *.proto ~/Projects/go-simple-bank/proto/protoc-gen-openapiv2/options/

# From `~/Downloads/swagger-ui`
ls -l dist
cp -r dist/* ~/Projects/go-simple-bank/doc/swagger/
# Change url in `swagger-initializer.js` file to `simple_bank.swagger.json`
```
