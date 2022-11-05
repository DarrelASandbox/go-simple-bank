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
  </ul>
</details>

&nbsp;

## About The Project

- Backend Master Class [Golang + Postgres + Kubernetes + gRPC]
- Learn everything about backend web development: Golang, Postgres, Gin, gRPC, Docker, Kubernetes, AWS, GitHub Actions
- [Original Repo - simplebank](https://github.com/techschool/simplebank)
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
psql simplebank -U root
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
docker build -t simplebank:latest .
docker run --name simplebank -p 4000:4000 -e GIN_MODE=release simplebank:latest
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
6. GitHub simplebank repo Settings -> Secrets -> New repository secret
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
4. **DB instance identifier:** simplebank
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
13. **Initial database name:** simplebank
14. Create database
15. View credential details
16. Register new Server/Connection for PostgreSQL
    1. **Name:** AWS Postgres
    2. **Host:** From AWS connection details endpoint _(See point 19)_
    3. **Username:** root
    4. **Password:** From AWS RDS credential details
    5. **Database:** simplebank
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
3. Next -> **Secret name:** simplebank -> Next -> Next -> Store
4. [Install AWS CLI](https://aws.amazon.com/cli/)
   1. `aws configure` _(Refer to IAM profile)_
   2. `ls -l ~/.aws`
   3. `cat ~/.aws/credentials`
   4. `aws secretsmanager help`
   5. `aws secretsmanager get-secret-value --secret-id simplebank` _Might want to try with arn as well_
   6. AWS IAM User groups -> Permissions -> Add permissions -> **Attach Policies:** SecretsManagerReadWrite
   7. `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text`
5. `brew install jq` (To output into json format)
   1. `aws secretsmanager get-secret-value --secret-id simplebank --query SecretString --output text | jq -r 'to_entries|map("\(.key)=\(.value)")|.[]' > app.env`
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
   1. **Name:** simplebank
   2. **Kubernetes version:** 1.20
   3. **Cluster Service Role:** AWSEKSClusterRole
      1. AWS IAM -> Create role -> `EKS - Cluster` use case -> Next (`AmazonEKSClusterPolicy`) -> Next
      2. **Role name:** AWSEKSClusterRole
2. Next -> Use default VPC
3. **Cluster endpoint access:** Public and private
4. **Amazon VPC CNI Version:** `v1.8.0-eksbuild. 1`
5. Next -> Next -> Create -> Refresh later when cluster is created
6. Add Node Group
   1. **Name:** simplebank
   2. **Node IAM Role:** AWSEKSNodeRole
      1. AWS IAM -> Create role -> `EC2` use case -> Next
      2. Pick `AmazonEKS_CNI_Policy` & `AmazonEKSWorkerNodePolicy` & `AmazonEC2ContainerRegistryReadOnly` -> Next
      3. **Role name:** AWSEKSNodeRole
   3. Next
   4. **AMI type:** Amazon Linux 2 (AL2_x86_64)
   5. **Capacity type:** On-Demand
   6. **Instance types:** t3.micro
   7. **Disk size:** 10 GiB
   8. **Node Group scaling configuration**
      1. **Minimum size:** 1 nodes
      2. **Maximum size:** 2 nodes
      3. **Desired size:** 1 nodes
   9. Next -> Disable **Allow remote access to nodes** -> Next -> Create
   10. Refresh later when node is created

&nbsp;

---

&nbsp;
