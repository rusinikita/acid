# Database Transactions & Concurrency — Learning Plan

A structured curriculum for senior-level engineering interviews. Work through each phase in order. The goal is not memorization — it is being able to reason out loud about trade-offs.

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

**Interview questions at this level:**
- What happens if a process crashes mid-transaction?
- What is auto-commit and when does it matter?
- Can you roll back a `CREATE TABLE`?

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

> **Note:** `acid` is useful for Atomicity (watch a rollback undo everything), Isolation (anomaly sequences in Phase 3), and Consistency — run two concurrent transactions both trying to insert the same unique value and watch one get a constraint violation, proving the database never allows an invalid state. Not useful for Durability — that requires crash recovery, which `acid` doesn't do.

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

**PostgreSQL specifics to know:**
- PostgreSQL does not implement READ UNCOMMITTED — its minimum is READ COMMITTED
- PostgreSQL's REPEATABLE READ uses snapshot isolation, which also prevents most phantoms in practice
- PostgreSQL's SERIALIZABLE uses Serializable Snapshot Isolation (SSI), not traditional locking — it detects conflicts at commit time

**Decision guide — when to use each level in production:**
- **READ COMMITTED** (default): most OLTP workloads, short transactions, accept occasional non-repeatable reads
- **REPEATABLE READ**: long-running read transactions, reports and analytics, anywhere you need a consistent snapshot across multiple queries
- **SERIALIZABLE**: financial systems, inventory allocation, anywhere write skew would be a real bug — accept the throughput cost

**Interview questions at this level:**
- What is the default isolation level in PostgreSQL? In MySQL?
- REPEATABLE READ prevents phantom reads in PostgreSQL — why? Does that mean it is equivalent to SERIALIZABLE?
- When would you actually use SERIALIZABLE in production?
- What anomaly does READ COMMITTED still allow?

> **Note:** `acid` is useful here — take any anomaly sequence from Phase 3, add `SET TRANSACTION ISOLATION LEVEL REPEATABLE READ` as a step, re-run, and compare the result. Seeing the same scenario behave differently under a different level is more convincing than a table.

---

## Phase 5 — Locking

**Lock types:**
- **Shared lock (S)** — multiple readers can hold simultaneously; blocks exclusive locks
- **Exclusive lock (X)** — only one holder; blocks both shared and exclusive locks
- **Row-level vs. table-level** — most databases default to row-level; table locks appear in bulk operations or DDL

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

**Interview questions at this level:**
- When would you use `SELECT FOR UPDATE` instead of an isolation level upgrade?
- What is two-phase locking? Does PostgreSQL use it?
- Compare pessimistic and optimistic locking — when do you choose each?
- What happens if you forget `SELECT FOR UPDATE` in a read-then-write transaction?

> **Note:** `acid` is useful for `SELECT FOR UPDATE` and optimistic locking — write a TOML where one transaction locks a row and another tries to update it, and watch the second block. 2PL is a theoretical concept; `acid` shows the effects of locking but not the internal protocol.

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
- When a row is updated, the old version is kept; a new version is written
- No dirty reads by design — you only see committed versions
- Readers never block writers; writers never block readers

**VACUUM and dead tuples:**
- Old row versions are not immediately deleted — they remain as "dead tuples"
- `VACUUM` reclaims space by removing dead tuples that no active snapshot needs
- `autovacuum` runs automatically; in write-heavy tables it may need tuning
- Long-running transactions can prevent VACUUM from removing old versions → table bloat

**Snapshot isolation vs. Serializable:**
- Snapshot isolation (REPEATABLE READ in PostgreSQL) gives each transaction a consistent snapshot but still allows write skew
- SSI (SERIALIZABLE in PostgreSQL) adds conflict detection on top of snapshot isolation to catch write skew

**Interview questions at this level:**
- How does MVCC allow readers and writers to not block each other?
- What is a dead tuple? Why does VACUUM exist?
- What is a long-running transaction's impact on MVCC storage?
- What is the difference between snapshot isolation and SERIALIZABLE in PostgreSQL?

> **Note:** `acid` can show that a reader and writer run concurrently without blocking — which is the observable effect of MVCC. It can't show internals: `xmin`/`xmax` values, dead tuples, or VACUUM. For those, use `psql` and query system columns directly.

---

## Phase 8 — Durability & Write-Ahead Logging

**Write-Ahead Logging (WAL):**
- Before any data page is modified on disk, the change is written to the WAL log first
- On crash, the database replays the WAL to recover committed transactions and discard uncommitted ones
- This is how the D in ACID is guaranteed — even if the process crashes between writing the WAL and flushing the data page, recovery replays the log

**What to know for interviews:**
- WAL enables both crash recovery and replication (streaming replication sends WAL to replicas)
- `fsync` must be on for durability guarantees — disabling it makes writes faster but you can lose committed data on a crash
- `synchronous_commit` controls whether a commit waits for WAL to flush to disk

**Interview questions at this level:**
- How does a database guarantee durability after a crash?
- What is the role of WAL in replication?
- What happens if you set `fsync=off` in PostgreSQL?

> **Note:** `acid` can't help here. WAL and crash recovery are storage-layer concepts — demonstrating them requires crashing the database process or inspecting WAL files. Use PostgreSQL docs and `pg_waldump` instead.

---

## Phase 9 — Production Patterns & Diagnostics

**Diagnosing lock contention:**
- `pg_locks` — shows currently held locks
- `pg_stat_activity` — shows active queries; `wait_event_type = 'Lock'` means a query is blocked
- Long-running transactions visible in `pg_stat_activity` → often the root cause of contention

**Diagnosing slow queries caused by locks:**
- Find blocking query: join `pg_locks` with `pg_stat_activity` on `pid`
- Common cause: a forgotten open transaction (e.g., an idle connection mid-transaction)

**Idempotency and retries:**
- When a transaction is rolled back (e.g., deadlock victim), the application must retry
- Retried operations must be idempotent — running them twice produces the same result
- Use unique constraints, conditional inserts (`INSERT ... ON CONFLICT`), and version checks to enforce idempotency

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
- [ ] Can describe all five concurrency anomalies (including write skew) with scenarios
- [ ] Can state what each isolation level prevents and allows
- [ ] Knows PostgreSQL's default isolation level and its MVCC behavior
- [ ] Can walk through a lost update bug and fix it with `SELECT FOR UPDATE`
- [ ] Can explain Coffman's four conditions for deadlock
- [ ] Knows how a database detects and resolves a deadlock
- [ ] Understands pessimistic vs. optimistic locking trade-offs
- [ ] Can explain MVCC and why readers don't block writers
- [ ] Knows what WAL is and how it guarantees durability
- [ ] Can describe the Saga pattern vs. 2PC at a conceptual level
- [ ] Can diagnose a blocked query using `pg_locks` and `pg_stat_activity`
