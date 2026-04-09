# Database Transactions Learning Path

Work through these tasks with your AI coach. Run each scenario in `acid` and discuss
the output before moving to the next task.

---

## Task 1 — What Is a Transaction?

**Goal:** Understand what BEGIN / COMMIT / ROLLBACK do.

**Explore:** Run `acid` (no arguments) and browse the built-in sequences. Press `s`
to reveal setup steps, `m` or Space to hide results before they render (quiz mode).

**Discussion questions:**
- What does "atomicity" mean? Can a transaction be half-applied?
- What happens to other connections if a transaction is never committed?
- What are the four ACID properties?

---

## Task 2 — Lost Update

**Scenario:** `sequences/lost_update.toml`

```
acid run -f sequences/lost_update.toml
```

**Before running:** Two transactions both read a balance of 1000 and both write back
900. What is the final balance? What did you expect?

**After running:** What did you observe?

**Key concept:** Lost Update — the second writer silently overwrites the first.

**Experiment:** Edit the TOML. Change one transaction's UPDATE to use arithmetic:

```toml
sql = "update accounts set balance = balance - 100 where id = 1"
```

How does this differ from the read-then-write pattern? Does it fix the problem?

---

## Task 3 — Dirty Read

**Scenario:** `sequences/dirty_read.toml`

```
acid run -f sequences/dirty_read.toml
```

**Before running:** `writer` inserts a row before committing. Will `reader` see it?

**After running:** What happened? Did PostgreSQL allow the dirty read?

**Key concept:** READ COMMITTED (PostgreSQL default) prevents dirty reads.

**Discussion:** Why is a dirty read dangerous in a payment system?

---

## Task 4 — Non-Repeatable Read

**Scenario:** `sequences/non_repeatable_read.toml`

```
acid run -f sequences/non_repeatable_read.toml
```

**Before running:** `auditor` reads a price twice. Between reads, `cashier` commits
an update. Will both reads return the same value?

**After running:** Explain in your own words why the second read returned a different value.

**Key concept:** Non-Repeatable Read — the same row returns different data within one
transaction because another transaction committed a change.

**Experiment:** Add this step immediately after the `begin` step for `auditor`:

```toml
[[steps]]
sql = "set transaction isolation level repeatable read"
trx = "auditor"
```

Re-run. Does the non-repeatable read still occur?

---

## Task 5 — Phantom Read

**Scenario:** `sequences/phantom_read.toml`

```
acid run -f sequences/phantom_read.toml
```

**Before running:** `reporter` counts active users twice. Between counts, `admin`
inserts a new active user. What does each count return?

**After running:** What is a "phantom row"? How is this different from a non-repeatable read?

**Key concept:** Phantom Read — a range query returns a different number of rows because
another transaction inserted or deleted qualifying rows.

---

## Task 6 — Isolation Level Comparison

**Goal:** See isolation levels prevent anomalies.

1. Choose `non_repeatable_read.toml`.
2. Add a step immediately after the `begin` for `auditor`:

```toml
[[steps]]
sql = "set transaction isolation level repeatable read"
trx = "auditor"
```

3. Re-run. Compare the result to Task 4.

**Discussion:** What guarantees does REPEATABLE READ provide? What does SERIALIZABLE add?

---

## Task 7 — Deadlock

**Goal:** Observe a deadlock and understand how the database resolves it.

Run `acid` (standalone mode) and select the "Deadlock" sequence from the list.

**After running:**
- Which transaction was chosen as the deadlock victim?
- What error did it receive?
- What happened to the surviving transaction?

**Key concept:** The database detects circular lock dependencies and aborts one
transaction to break the cycle.

---

## Task 8 — Design Your Own Scenario

**Goal:** Write a TOML file from scratch.

Ideas:
- `SELECT ... FOR UPDATE` blocking a concurrent writer
- Two transactions transferring money between accounts (safe vs. unsafe pattern)
- A unique-constraint violation in a concurrent context

Steps:
1. Create `sequences/my_scenario.toml`
2. Define the table in `setup = true` steps
3. Use at least two named transactions (`trx = "tx1"` and `trx = "tx2"`)
4. Run it: `acid run -f sequences/my_scenario.toml`
5. Explain to your coach what you expected and what happened

---

## Completion Checklist

- [ ] I can explain the four ACID properties in my own words
- [ ] I know what isolation level PostgreSQL uses by default
- [ ] I can explain the difference between a dirty read, non-repeatable read, and phantom read
- [ ] I know what a lost update is and two ways to prevent it
- [ ] I understand what a deadlock is and how the database recovers from it
- [ ] I wrote at least one TOML scenario file myself
