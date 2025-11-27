## Date: 2025-11-24
### Password Validation

[OWASP](https://owasp.org/) no longer recommends strict composition rules like:
“must contain 1 uppercase”
“must contain 1 lowercase”
“must contain 1 number”
“must contain 1 special character”
Why?

Because composition rules do NOT significantly improve security and they make passwords harder for users to remember!


## Date: 2025-11-25
```go 
value := x.(string)
```

This tells go that we expect x to be a string. 

```go
str, ok := x.(string)
```

For failure case we can make use of two variables.

### Key Concepts Mastered
- **Clean Architecture**: Separation of handlers → services → repositories
- **Dependency Injection**: Interface-based design for testability
- **JWT Authentication**: Token-based auth with middleware
- **Docker Multi-stage Builds**: Production vs development environments
- **Testing Strategy**: Unit tests with mocks, integration tests with dockertest
- **Database Design**: GORM with soft deletes, migrations, foreign keys
