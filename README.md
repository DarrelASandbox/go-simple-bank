## About The Project

- Backend Master Class [Golang + Postgres + Kubernetes + gRPC]
- Learn everything about backend web development: Golang, Postgres, Gin, gRPC, Docker, Kubernetes, AWS, GitHub Actions
- [Original Repo - simplebank](https://github.com/techschool/simplebank)
- [techschool](https://github.com/techschool)

&nbsp;

---

&nbsp;

## Notes

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
# Login using u1
psql simplebank -U u1
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

- [GitHub - Creating PostgreSQL service containers](https://docs.github.com/en/actions/using-containerized-services/creating-postgresql-service-containers)

1. Actions -> Configure in Go -> Copy code to vscode

&nbsp;

---

&nbsp;
