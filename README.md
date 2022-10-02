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
# In db folder
sqlc generate
```

&nbsp;

---

&nbsp;
