# Taskflow CLI

Taskflow is a simple CLI tool for running shell commands on a schedule using cron expressions or human-readable intervals.

### Commands

#### `taskflow exec`

Run shell commands directly.

```bash
taskflow exec "echo Hello"
```

Run multiple commands:

```bash
taskflow exec "echo Task 1" "echo Task 2"
```

#### `taskflow scheduler`

Schedule commands with cron or human-readable intervals.

```bash
taskflow scheduler --cron "* * * * *" --cmd "echo Hello"
```

#### Using a File for Jobs

You can also schedule jobs from a file:

```bash
taskflow scheduler --file "jobs.txt"
```

Where `jobs.txt` contains cron expressions like:

```txt
*/5 * * * * echo "Every 5 minutes"
```

### Cron Expression Format

Cron expressions are in the standard format:

```
* * * * *   (minute, hour, day, month, weekday)
```
