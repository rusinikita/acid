# Database Transactions & Concurrency — Learning Plan

A structured curriculum for senior-level engineering interviews. Work through each phase in order. The goal is not memorization — it is being able to reason out loud about trade-offs.

Covers both **PostgreSQL** and **MySQL (InnoDB)**. Where behavior differs between the two, both are explained.

---

## Phase 1 — Transaction Fundamentals

**What to learn:**
- What a transaction is and why databases need them
- `BEGIN` / `COMMIT` / `ROLLBACK` semantics
- Auto-commit mode vs. explicit transaction boundaries
- Savepoints and nested transactions

**Key concepts to be able to explain:**
- Why partial writes are dangerous without transactions
- What "all or nothing" means in practice
- The difference between a statement failing and a transaction being rolled back

**Nested transactions:**
- True nested transactions (a transaction inside a transaction) are not supported by PostgreSQL or MySQL
- **Savepoints** are the practical substitute: mark a point within a transaction that you can roll back to without aborting the whole transaction (`SAVEPOINT name` / `ROLLBACK TO SAVEPOINT name` / `RELEASE SAVEPOINT name`)
- Application-level "nested transaction" abstractions (e.g., in ORMs) typically use savepoints under the hood

**Implicit commits in MySQL:**
Certain statements automatically commit the current transaction before executing. If you have open changes, they are committed without warning:
- All DDL: `CREATE`, `ALTER`, `DROP`, `RENAME`, `TRUNCATE TABLE`
- `LOCK TABLES`, `UNLOCK TABLES`
- `BEGIN` (starting a new transaction commits the previous one)
- Administration statements: `ANALYZE TABLE`, `OPTIMIZE TABLE`, `REPAIR TABLE`
- PostgreSQL does not have implicit commits — DDL is fully transactional

**Which operations cannot be rolled back:**
- **MySQL**: any DDL statement (see implicit commits above); also `TRUNCATE TABLE` — it is DDL in MySQL, not DML
- **PostgreSQL**: almost everything can be rolled back, including DDL; exceptions are operations outside the transaction scope (e.g., sequences, `pg_sleep`, external side effects)

**Database differences to know:**
- Both PostgreSQL and MySQL support `BEGIN` / `COMMIT` / `ROLLBACK`; MySQL also uses `START TRANSACTION` as an alias for `BEGIN`

**Interview questions at this level:**
- What happens if a process crashes mid-transaction?
- What is auto-commit and when does it matter?
- Can you roll back a `CREATE TABLE`? (answer differs by database)
- What is a savepoint? How does it differ from a nested transaction?
- In MySQL, which statements cause an implicit commit?

> **Note:** `acid` is useful here — run a sequence with an explicit rollback and watch all intermediate writes disappear. Makes "all or nothing" tangible.

---

## Phase 2 — ACID Properties

**What to learn:**
- **Atomicity** — all statements in a transaction succeed or all are rolled back
- **Consistency** — a transaction brings the database from one valid state to another; all constraints hold
- **Isolation** — concurrent transactions behave as if they run serially (to some degree)
- **Durability** — once committed, data survives crashes

**Key concepts to be able to explain:**
- A concrete real-world example for each property (e.g., bank transfer for atomicity)
- Why consistency is partly the database's job and partly the application's job
- Why isolation is the most nuanced property — it is a spectrum, not a toggle

**Interview questions at this level:**
- Give me a real example where violating atomicity causes a bug
- What does "consistency" actually mean — who is responsible for it?
- Why is full isolation expensive?

> **Note:** `acid` is useful for Atomicity (watch a rollback undo everything), Isolation (anomaly sequences in Phase 3), and Consistency — two good scenarios: (1) two transactions both insert the same unique value, one gets a constraint violation; (2) two transactions both run `UPDATE SET value = value + 1` with a `CHECK (value < 2)` constraint — the first commits, the second re-reads the updated value and the increment violates the constraint, so the database rejects it. Not useful for Durability — that requires crash recovery, which `acid` doesn't do.

---

## Phase 3 — Concurrency Anomalies

These are the bugs that happen when isolation is imperfect. Know each one well enough to describe it with a concrete scenario.

**Dirty Read**
- Transaction A reads data written by transaction B before B commits
- B then rolls back — A has read data that never existed

**Non-Repeatable Read**
- Transaction A reads a row; transaction B updates and commits it; A reads the same row again and gets a different value
- Same query, same transaction, different results

**Phantom Read**
- Transaction A runs a range query; transaction B inserts a new row matching that range and commits; A runs the same range query and gets an extra row

**Lost Update**
- Two transactions both read a value (e.g., balance = 1000), compute a new value, and write back — the second write silently overwrites the first
- Classic bug in read-then-write patterns

**Write Skew** *(asked at senior/staff level)*
- Two transactions each read overlapping data, make decisions based on what they see, and write to different rows — but the combined result violates a constraint that either transaction alone would have respected
- Example: two doctors both check "is there at least one doctor on call?" and both go off call simultaneously

**Interview questions at this level:**
- Walk me through a lost update scenario with a bank account
- What is the difference between a phantom read and a non-repeatable read?
- What is write skew? Give me an example
- Which anomalies can happen at READ COMMITTED?

> **Note:** `acid` is most useful here — built-in sequences exist for dirty read, non-repeatable read, phantom read, and lost update. Use quiz mode (hide results, predict the output, then reveal) to build real intuition. Write skew has no built-in sequence but you can write one as a TOML.

---

## Phase 4 — Isolation Levels

**The four standard levels (SQL standard):**

| Level | Dirty Read | Non-Repeatable Read | Phantom Read |
|---|---|---|---|
| READ UNCOMMITTED | possible | possible | possible |
| READ COMMITTED | prevented | possible | possible |
| REPEATABLE READ | prevented | prevented | possible |
| SERIALIZABLE | prevented | prevented | prevented |

**Database-specific behavior:**

| | PostgreSQL | MySQL (InnoDB) |
|---|---|---|
| Default level | READ COMMITTED | REPEATABLE READ |
| READ UNCOMMITTED | Not implemented (treated as READ COMMITTED) | Supported |
| REPEATABLE READ | Snapshot isolation — also prevents most phantoms | Gap locks + next-key locks prevent phantoms |
| SERIALIZABLE | Serializable Snapshot Isolation (SSI) — detects conflicts at commit time | Converts all plain SELECTs to SELECT FOR SHARE |

**MySQL gap locks and next-key locks** *(important — unique to MySQL)*:
- At REPEATABLE READ, InnoDB locks not just the rows that match a query but also the gaps between them
- This prevents other transactions from inserting new rows into a range you've already scanned, blocking phantom reads
- Gap locks can cause unexpected contention and deadlocks that do not occur in PostgreSQL — a common source of production surprises

**Decision guide — when to use each level in production:**
- **READ COMMITTED** (PostgreSQL default): most OLTP workloads, short transactions, accept occasional non-repeatable reads
- **REPEATABLE READ** (MySQL default): consistent snapshot across multiple queries; in MySQL also prevents phantoms via gap locks
- **SERIALIZABLE**: financial systems, inventory allocation, anywhere write skew would be a real bug — accept the throughput cost

**Interview questions at this level:**
- What is the default isolation level in PostgreSQL? In MySQL?
- REPEATABLE READ prevents phantom reads in PostgreSQL — why? How does MySQL prevent them differently?
- What is a gap lock? What problem does it solve and what problem can it cause?
- When would you actually use SERIALIZABLE in production?
- What anomaly does READ COMMITTED still allow?

> **Note:** `acid` is useful here — take any anomaly sequence from Phase 3, add `SET TRANSACTION ISOLATION LEVEL REPEATABLE READ` as a step, re-run, and compare the result. Seeing the same scenario behave differently under a different level is more convincing than a table.

---

## Phase 5 — Locking

**Lock types:**
- **Shared lock (S)** — multiple readers can hold simultaneously; blocks exclusive locks
- **Exclusive lock (X)** — only one holder; blocks both shared and exclusive locks
- **Row-level vs. table-level** — most databases default to row-level; table locks appear in bulk operations or DDL

**Explicit table-level locking:**

*PostgreSQL* has 8 lock modes, from least to most restrictive:

| Mode | Acquired by | Conflicts with |
|---|---|---|
| ACCESS SHARE | `SELECT` | ACCESS EXCLUSIVE only |
| ROW SHARE | `SELECT FOR UPDATE/SHARE` | EXCLUSIVE, ACCESS EXCLUSIVE |
| ROW EXCLUSIVE | `INSERT`, `UPDATE`, `DELETE` | SHARE and above |
| SHARE UPDATE EXCLUSIVE | `VACUUM`, `ANALYZE`, `CREATE INDEX CONCURRENTLY` | SHARE UPDATE EXCLUSIVE and above |
| SHARE | `CREATE INDEX` | ROW EXCLUSIVE and above |
| SHARE ROW EXCLUSIVE | `CREATE TRIGGER`, some `ALTER TABLE` | ROW EXCLUSIVE and above |
| EXCLUSIVE | rarely used explicitly | ROW SHARE and above |
| ACCESS EXCLUSIVE | `DROP TABLE`, `TRUNCATE`, `ALTER TABLE`, `REINDEX` | everything, including `SELECT` |

*MySQL* uses a simpler explicit locking model:
- `LOCK TABLES tbl READ` — shared table lock; other sessions can read but not write
- `LOCK TABLES tbl WRITE` — exclusive table lock; blocks all other access
- MySQL also has **Metadata Locks (MDL)**: automatically acquired on any table access; DDL waits for all active transactions on the table to finish before proceeding

Key things to know (both databases):
- Most application queries never need explicit `LOCK TABLE` — the database acquires the right mode automatically
- In PostgreSQL, ACCESS EXCLUSIVE (held by `ALTER TABLE`) blocks even plain `SELECT` — this is why schema migrations on live tables cause outages
- In MySQL, DDL waits for MDL, which means a single long-running transaction can block an `ALTER TABLE` which then blocks all subsequent queries on that table (lock queue)

**`SELECT FOR UPDATE` and `SELECT FOR SHARE`:**
- `SELECT FOR UPDATE` acquires an exclusive row lock — blocks other writers and other `SELECT FOR UPDATE` until the transaction commits
- `SELECT FOR SHARE` acquires a shared row lock — allows other readers, blocks writers
- Use case: read-then-write patterns where you must prevent a concurrent update between the read and the write

**Two-Phase Locking (2PL):**
- **Growing phase**: transaction acquires all locks it needs, releases none
- **Shrinking phase**: transaction releases locks, acquires no new ones
- Guarantees serializability but limits concurrency
- Strict 2PL (hold all locks until commit) is what most databases implement

**Pessimistic vs. optimistic locking:**

| | Pessimistic | Optimistic |
|---|---|---|
| Mechanism | Lock before reading | Check version at commit |
| Good when | High contention, conflicts likely | Low contention, conflicts rare |
| Cost of conflict | Waiting (blocking) | Rollback and retry |
| Implementation | `SELECT FOR UPDATE` | Version column + conditional UPDATE |

**Optimistic locking with a version column:**
```sql
-- Read
SELECT id, balance, version FROM accounts WHERE id = 1;

-- Write (only succeeds if version hasn't changed)
UPDATE accounts
SET balance = 900, version = version + 1
WHERE id = 1 AND version = <version_you_read>;

-- If 0 rows updated → someone else changed it → retry
```

**Advisory Locks:**
- Application-level named locks provided by the database — the database enforces mutual exclusion but assigns no meaning to the lock key; semantics are entirely defined by the application
- Use cases: ensure only one worker processes a job at a time, prevent concurrent cron job execution, coordinate distributed processes that share a database without a dedicated lock service

*PostgreSQL*:
- `pg_advisory_lock(key)` — session-level exclusive lock, held until explicitly released or session ends
- `pg_advisory_xact_lock(key)` — transaction-level exclusive lock, released automatically on commit/rollback
- `pg_try_advisory_lock(key)` — non-blocking variant, returns `true` if acquired, `false` if already held
- Shared variants exist: `pg_advisory_lock_shared`, `pg_advisory_xact_lock_shared`
- Key is an integer (or two integers)

*MySQL*:
- `GET_LOCK('name', timeout)` — acquires a named string lock; `timeout = 0` is non-blocking, `-1` waits forever
- `RELEASE_LOCK('name')` — explicitly releases the lock
- `IS_FREE_LOCK('name')` — check without acquiring
- Key is an arbitrary string; always session-level (no transaction-level variant)
- One important difference: in MySQL a session can only hold one named lock at a time in versions before 5.7.5; from 5.7.5 multiple locks are supported

**Interview questions at this level:**
- When would you use `SELECT FOR UPDATE` instead of an isolation level upgrade?
- What is two-phase locking? Does PostgreSQL use it?
- Compare pessimistic and optimistic locking — when do you choose each?
- What happens if you forget `SELECT FOR UPDATE` in a read-then-write transaction?
- What is an advisory lock? When would you use one instead of a row lock?
- What is the difference between a session-level and transaction-level advisory lock?

> **Note:** `acid` is useful for `SELECT FOR UPDATE`, optimistic locking, and advisory locks — write a TOML where one transaction calls `pg_advisory_xact_lock` and another tries to acquire the same key, and watch the second block. 2PL is a theoretical concept; `acid` shows the effects of locking but not the internal protocol.

---

## Phase 6 — Deadlocks

**What a deadlock is:**
Transaction A holds lock 1 and waits for lock 2. Transaction B holds lock 2 and waits for lock 1. Neither can proceed.

**Coffman's four necessary conditions** *(always asked at interviews)*:
1. **Mutual exclusion** — a resource can be held by only one transaction at a time
2. **Hold and wait** — a transaction holds resources while waiting to acquire more
3. **No preemption** — locks cannot be forcibly taken away; they must be released voluntarily
4. **Circular wait** — there is a circular chain of transactions each waiting on the next

All four must hold simultaneously for a deadlock to occur. Prevention strategies target breaking one of them.

**Detection and resolution:**
- Databases maintain a wait-for graph; a cycle means deadlock
- The database picks a victim (usually the transaction with the least work done) and rolls it back
- The surviving transaction proceeds; your application must retry the victim

**Prevention strategies:**
- Always acquire locks in a consistent order across the codebase (breaks circular wait)
- Keep transactions short — acquire locks late, release early (reduces hold-and-wait window)
- Use lower isolation levels where safe (fewer locks acquired)
- Consider `NOWAIT` or `SKIP LOCKED` when you prefer to fail fast over waiting

**Interview questions at this level:**
- Describe Coffman's conditions. Which one is easiest to break in practice?
- How does a database detect a deadlock?
- How should your application handle the "deadlock detected" error?
- Two transactions each transfer money between account A and account B but in opposite order — what happens and how do you fix it?

> **Note:** `acid` is useful here — there is a built-in deadlock sequence. You can see which transaction gets picked as the victim and observe the error. The step-by-step TUI makes the circular wait visually obvious.

---

## Phase 7 — MVCC (Multi-Version Concurrency Control)

**What it is and why it exists:**
Traditional locking means readers block writers and writers block readers. MVCC solves this by keeping multiple versions of each row — each transaction sees a snapshot of the database as of when its transaction started.

**How PostgreSQL implements it:**
- Every row has `xmin` (the transaction that created it) and `xmax` (the transaction that deleted/updated it)
- A transaction's snapshot defines which `xmin`/`xmax` values it can see
- When a row is updated, the old version is kept in the heap; a new version is written alongside it
- No dirty reads by design — you only see committed versions
- Readers never block writers; writers never block readers

**VACUUM and dead tuples (PostgreSQL):**
- Old row versions are not immediately deleted — they remain as "dead tuples" in the heap
- `VACUUM` reclaims space by removing dead tuples that no active snapshot needs
- `autovacuum` runs automatically; in write-heavy tables it may need tuning
- Long-running transactions can prevent VACUUM from removing old versions → table bloat

**How MySQL (InnoDB) implements it:**
- Old row versions are stored in the **undo log**, not in the main table heap
- Each row in the heap has a pointer to its undo log chain; a transaction follows the chain to find the version it should see
- InnoDB has a background **purge thread** that cleans up undo log entries no longer needed by any active transaction — there is no manual `VACUUM`
- Long-running transactions cause undo log growth (similar problem to PostgreSQL bloat, different mechanism)

**Snapshot isolation vs. Serializable:**
- Snapshot isolation (REPEATABLE READ) gives each transaction a consistent snapshot but still allows write skew
- PostgreSQL SSI adds conflict detection on top of snapshot isolation to catch write skew
- MySQL SERIALIZABLE prevents write skew by converting plain SELECTs to SELECT FOR SHARE, which blocks conflicting writes

**Interview questions at this level:**
- How does MVCC allow readers and writers to not block each other?
- What is a dead tuple in PostgreSQL? What is the equivalent in MySQL?
- What is a long-running transaction's impact on MVCC storage in each database?
- What is the difference between snapshot isolation and SERIALIZABLE?

> **Note:** `acid` can show that a reader and writer run concurrently without blocking — which is the observable effect of MVCC. It can't show internals: `xmin`/`xmax` values, undo log chains, dead tuples, or VACUUM. For those, use `psql`/`mysql` and query system tables directly.

---

## Phase 8 — Durability & Write-Ahead Logging

**The concept (both databases):**
- Before any data page is modified on disk, the change is written to a log first
- On crash, the database replays the log to recover committed transactions and discard uncommitted ones
- This guarantees the D in ACID — even if the process crashes mid-write, recovery restores a consistent state

**PostgreSQL — WAL (Write-Ahead Log):**
- Single unified log used for both crash recovery and replication (streaming replication ships WAL segments to replicas)
- `fsync` must be on — disabling it risks losing committed data on a crash
- `synchronous_commit` controls whether a commit waits for WAL to flush to disk (setting to `off` improves latency at the cost of potential data loss on crash)

**MySQL — Redo Log + Binary Log:**
- InnoDB uses a **redo log** for crash recovery (equivalent role to PostgreSQL's WAL)
- MySQL also has a separate **binary log (binlog)** used for replication and point-in-time recovery — it is engine-agnostic and logs logical operations, not physical page changes
- `innodb_flush_log_at_trx_commit` controls durability: `1` (default) = flush on every commit; `2` = flush once per second; `0` = no flush (fastest, least safe)
- Replication in MySQL sends binlog events to replicas, not redo log — this is architecturally different from PostgreSQL

**Interview questions at this level:**
- How does a database guarantee durability after a crash?
- What is the difference between PostgreSQL WAL and MySQL binlog?
- What does `innodb_flush_log_at_trx_commit=2` mean in MySQL? What do you risk?
- What happens if you set `fsync=off` in PostgreSQL?

> **Note:** `acid` can't help here. WAL and crash recovery are storage-layer concepts — demonstrating them requires crashing the database process or inspecting log files. Use PostgreSQL docs / `pg_waldump` or MySQL docs / `mysqlbinlog` instead.

---

## Phase 9 — Production Patterns & Diagnostics

**Diagnosing lock contention:**

*PostgreSQL:*
- `pg_locks` — shows currently held locks
- `pg_stat_activity` — shows active queries; `wait_event_type = 'Lock'` means a query is blocked
- Find blocking query: join `pg_locks` with `pg_stat_activity` on `pid`

*MySQL:*
- `SHOW ENGINE INNODB STATUS` — shows the latest deadlock, current lock waits, and transaction list
- `performance_schema.data_locks` (MySQL 8.0+) — equivalent to `pg_locks`, shows row and table locks
- `performance_schema.data_lock_waits` — shows which transaction is blocking which
- `SHOW PROCESSLIST` or `performance_schema.processlist` — equivalent to `pg_stat_activity`

Common cause in both: a forgotten open transaction (idle-in-transaction connection) holding locks

**Idempotency and retries:**
- When a transaction is rolled back (e.g., deadlock victim), the application must retry
- Retried operations must be idempotent — running them twice produces the same result
- PostgreSQL: `INSERT ... ON CONFLICT DO UPDATE/NOTHING`
- MySQL: `INSERT ... ON DUPLICATE KEY UPDATE` or `REPLACE INTO`
- Both: version column checks in UPDATE for optimistic locking

**Distributed transactions (awareness level for senior interviews):**
- **2-Phase Commit (2PC)**: coordinator asks all participants to prepare, then commits if all agree — synchronous, strong consistency, but slow and coordinator is a single point of failure
- **Saga pattern**: sequence of local transactions, each publishing an event or triggering the next step; compensating transactions handle rollbacks — async, eventual consistency, complex to reason about
- When to use which: 2PC for tight consistency requirements on small participant sets; Saga for long-running business processes across services

**Interview questions at this level:**
- A query has been running for 10 minutes and seems stuck — how do you diagnose it?
- What is an idle-in-transaction connection and why is it dangerous?
- How do you make a payment operation safe to retry?
- What is the trade-off between 2PC and the Saga pattern?

> **Note:** `acid` can demonstrate a transaction blocking another (open a `BEGIN`, leave it uncommitted, watch a second transaction wait). It can't help with `pg_locks`/`pg_stat_activity` diagnosis, idempotency patterns, or distributed transactions — use `psql` and real application code for those.

---

## Phase 10 — Indexes

**Purpose:**
- Indexes allow the database to find rows without scanning the entire table
- Trade-off: faster reads, slower writes (index must be updated on every INSERT/UPDATE/DELETE), more storage

**Clustered vs. non-clustered:**
- **Clustered index** — the table data is physically stored in index order; there can be only one per table; a lookup by clustered index key retrieves the row directly
- **Non-clustered index** — a separate structure that stores index keys and pointers (row IDs or clustered key values) to the actual row; a lookup requires following the pointer after finding the key

*InnoDB (MySQL):*
- Every InnoDB table has exactly one clustered index — by default the primary key
- If no primary key is defined, InnoDB picks the first UNIQUE NOT NULL index; if none exists, it creates a hidden 6-byte row ID clustered index
- All secondary (non-clustered) indexes store the primary key value as their row pointer — so a secondary index lookup first finds the PK, then does a second lookup into the clustered index ("double lookup")
- An InnoDB table without any user-defined index is valid but uses the hidden clustered index; typically a bad idea for performance

*PostgreSQL:*
- PostgreSQL does not have clustered indexes in the InnoDB sense — the heap is unordered by default
- `CLUSTER tbl USING idx` physically reorders the heap once, but the ordering is not maintained on future writes
- All PostgreSQL indexes are non-clustered by this definition; they store a TID (physical row location) as the pointer

**Interview questions at this level:**
- What is the difference between a clustered and a non-clustered index?
- In InnoDB, what happens if you define no primary key?
- Why does a secondary index lookup in InnoDB require two lookups?
- Can an InnoDB table exist without indexes? What are the implications?

> **Note:** `acid` can't help here — index behavior requires running queries against real tables and using EXPLAIN to observe access patterns.

---

## Phase 11 — Query Plans

**Getting the query plan:**
- *MySQL*: `EXPLAIN SELECT ...` — shows the optimizer's plan without executing; `EXPLAIN ANALYZE SELECT ...` (MySQL 8.0+) — executes and shows actual timings
- *PostgreSQL*: `EXPLAIN SELECT ...` — estimated plan; `EXPLAIN ANALYZE SELECT ...` — executes and shows actual vs. estimated rows and timings

**Key fields in MySQL EXPLAIN output:**

| Field | What to look for |
|---|---|
| `type` | Access method — best to worst: `const` → `eq_ref` → `ref` → `range` → `index` → `ALL` |
| `key` | Which index was used; `NULL` means no index used |
| `rows` | Estimated rows examined — high numbers signal a performance problem |
| `Extra` | `Using filesort` and `Using temporary` are warning signs; `Using index` means a covering index (fast) |

**What `type` values mean:**
- `ALL` — full table scan; almost always bad on large tables
- `index` — full index scan; better than `ALL` but still reads everything
- `range` — index range scan; good for `BETWEEN`, `>`, `<`, `IN`
- `ref` — index lookup for non-unique key; good
- `eq_ref` — index lookup for unique key (one row result); very good
- `const` — primary key or unique index with a constant value; optimal

**What to do when you see a bad plan:**
- `type = ALL` → add an index on the WHERE/JOIN column
- `Using filesort` → add an index that covers the ORDER BY column
- `Using temporary` → query uses a temp table (often GROUP BY or DISTINCT without an index)
- `rows` estimate much higher than actual → stale statistics; run `ANALYZE TABLE` (MySQL) or `ANALYZE` (PostgreSQL)

**Interview questions at this level:**
- How do you get a query execution plan in MySQL? In PostgreSQL?
- What does `type = ALL` in MySQL EXPLAIN mean? How do you fix it?
- What is the difference between `EXPLAIN` and `EXPLAIN ANALYZE`?
- What does `Using filesort` mean? Is it always a problem?

> **Note:** `acid` can't help here — query plan analysis requires running EXPLAIN against real tables with real data distributions.

---

## Phase 12 — SQL Fundamentals

**TRUNCATE vs. DELETE:**

| | `DELETE` | `TRUNCATE` |
|---|---|---|
| Removes rows | One by one, with WHERE support | All rows at once, no WHERE |
| Transaction-safe | Yes — can be rolled back | MySQL: No (DDL, implicit commit); PostgreSQL: Yes |
| Triggers | Fires row-level DELETE triggers | Does not fire row-level triggers |
| Auto-increment reset | No | Yes (resets counter) |
| Speed | Slower on large tables | Much faster (deallocates pages) |
| Locks | Row locks | Table lock |

**JOIN types:**
- **INNER JOIN** — returns only rows with a match in both tables; non-matching rows from either side are excluded
- **LEFT JOIN** (LEFT OUTER JOIN) — returns all rows from the left table; for rows with no match in the right table, right-side columns are NULL
- **RIGHT JOIN** (RIGHT OUTER JOIN) — returns all rows from the right table; non-matching left-side columns are NULL; equivalent to swapping table order in a LEFT JOIN
- **FULL JOIN** (FULL OUTER JOIN) — returns all rows from both tables; NULL fills the side with no match; *MySQL does not support FULL JOIN natively* — emulate with `LEFT JOIN UNION ALL RIGHT JOIN ... WHERE left.id IS NULL`
- **CROSS JOIN** — cartesian product; every row from the left paired with every row from the right; rarely intentional

**Interview questions at this level:**
- What is the difference between DELETE and TRUNCATE? Can you roll back TRUNCATE?
- What does LEFT JOIN return when there is no match on the right side?
- How do you emulate FULL JOIN in MySQL?
- Why would TRUNCATE be faster than DELETE on a large table?

> **Note:** `acid` can't help here — these are query-level concepts best practiced in a SQL shell.

---

## Phase 13 — Schema Design & DDL

**Normalization:**
- The process of structuring tables to reduce redundancy and ensure data integrity

- **1NF (First Normal Form)**: each column holds atomic (indivisible) values; no repeating groups or arrays in a column; each row is uniquely identifiable
- **2NF (Second Normal Form)**: must be in 1NF; every non-key column must depend on the *whole* primary key — eliminates partial dependencies (relevant when the PK is composite)
- **3NF (Third Normal Form)**: must be in 2NF; every non-key column must depend *directly* on the primary key, not on another non-key column — eliminates transitive dependencies
- **BCNF** (Boyce-Codd, awareness level): a stricter form of 3NF; every determinant must be a candidate key

Practical note: full normalization optimizes for write consistency; denormalization (intentionally violating NF rules) optimizes for read performance. Most production schemas are selectively denormalized.

**Instant DDL (MySQL 8.0+):**
- Traditional `ALTER TABLE` requires rebuilding the entire table — it copies all rows to a new table structure, which is slow and locks the table
- **Instant DDL** performs the schema change by updating only the table metadata — no row data is touched, making it nearly instantaneous regardless of table size
- Operations supported for instant DDL (examples): adding a column at the end of the table, changing column default values, renaming columns, adding/dropping virtual columns, changing index visibility
- Operations that still require a table rebuild: changing a column's data type, adding a column in the middle of the table (in some versions), changing character set

*PostgreSQL equivalent:* some DDL is also instant (e.g., `ALTER TABLE ... ALTER COLUMN ... SET DEFAULT`, `ALTER TABLE ... ADD COLUMN` with a non-volatile default in Postgres 11+); others require a rewrite

**Interview questions at this level:**
- Explain 1NF, 2NF, and 3NF with an example
- When would you intentionally denormalize a table?
- What is Instant DDL in MySQL? Give two examples of operations it supports
- Why does traditional `ALTER TABLE ADD COLUMN` lock the table? How does Instant DDL avoid this?

> **Note:** `acid` can't help here — schema design and DDL behavior are best explored directly in a MySQL or PostgreSQL shell.

---

## Capstone — Design a Scenario from Scratch

Given a business requirement (e.g., "transfer money between two accounts safely"), design the complete transaction strategy:

1. Identify which anomalies could occur
2. Choose the appropriate isolation level and justify it
3. Decide between pessimistic and optimistic locking
4. Define the exact SQL, including `SELECT FOR UPDATE` if needed
5. Handle the error cases: deadlock, serialization failure, constraint violation
6. Explain what happens under high concurrency

> **Note:** `acid` is the right tool here — write the scenario as a TOML, run it, and verify your predictions match reality.

---

## Interview Readiness Checklist

- [ ] Can explain all four ACID properties with concrete examples
- [ ] Knows which MySQL statements cause implicit commits and which operations cannot be rolled back
- [ ] Can explain savepoints and how they differ from true nested transactions
- [ ] Can describe all five concurrency anomalies (including write skew) with scenarios
- [ ] Can state what each isolation level prevents and allows
- [ ] Knows the default isolation level in PostgreSQL (READ COMMITTED) and MySQL (REPEATABLE READ) and why they differ
- [ ] Can explain MySQL gap locks and next-key locks and why they can cause unexpected deadlocks
- [ ] Can walk through a lost update bug and fix it with `SELECT FOR UPDATE`
- [ ] Can explain Coffman's four conditions for deadlock
- [ ] Knows how a database detects and resolves a deadlock
- [ ] Understands pessimistic vs. optimistic locking trade-offs
- [ ] Can explain the difference between clustered and non-clustered indexes; knows InnoDB's hidden clustered index behavior
- [ ] Can read a MySQL EXPLAIN output and identify full table scans, missing indexes, and filesort
- [ ] Knows the difference between DELETE and TRUNCATE (transaction safety, triggers, speed)
- [ ] Can explain all four JOIN types and knows FULL JOIN is not natively supported in MySQL
- [ ] Can explain 1NF, 2NF, and 3NF with a concrete example
- [ ] Knows what Instant DDL is in MySQL and can give examples of instant vs. rebuild operations
- [ ] Can explain MVCC and why readers don't block writers; knows how PostgreSQL (heap versioning + VACUUM) and MySQL (undo log + purge thread) differ
- [ ] Knows what WAL is in PostgreSQL and what redo log + binlog are in MySQL
- [ ] Can describe the Saga pattern vs. 2PC at a conceptual level
- [ ] Can diagnose a blocked query in PostgreSQL (`pg_locks`, `pg_stat_activity`) and MySQL (`SHOW ENGINE INNODB STATUS`, `performance_schema.data_locks`)
