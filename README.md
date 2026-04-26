# envdiff

Compare `.env` files across environments and highlight missing or mismatched keys with structured output.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
envdiff --base .env --compare .env.production
```

### Example Output

```
MISSING IN .env.production:
  - DATABASE_URL
  - REDIS_HOST

MISMATCHED KEYS:
  - APP_ENV: "development" → "production"
  - LOG_LEVEL: "debug" → "warn"

OK: 12 keys match across both files.
```

### Flags

| Flag        | Description                          | Default  |
|-------------|--------------------------------------|----------|
| `--base`    | Base `.env` file to compare from     | `.env`   |
| `--compare` | Target `.env` file to compare against| required |
| `--format`  | Output format: `text`, `json`        | `text`   |
| `--strict`  | Exit with non-zero code on any diff  | `false`  |

---

## Why envdiff?

Managing environment variables across staging, production, and local setups is error-prone. `envdiff` gives you a fast, scriptable way to catch configuration drift before it causes issues in deployment.

---

## License

MIT © [yourusername](https://github.com/yourusername)