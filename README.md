# database — порты для работы с БД

`Transactor`, `TxTransactor`, `HealthChecker`, `DB`, `ErrNotFound`. Единый паттерн транзакций через контекст для pgxdb/gormdb/sqldb.

```go
tx.WithinTx(ctx, func(ctx context.Context) error { return repo.Create(ctx, u) })
```

## Установка

```bash
go get github.com/SaveliiYam/database@v0.1.0
```

Требуется Go ≥ 1.25. Лицензия: см. репозиторий.
