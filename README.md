# gogetgmail

Simple Go utility to list Gmail messages for a label and print headers/body preview.

Quick start

1. Place your Google OAuth client file in the project root and name it `credentials.json`.
   - (You already have a `credentials.json` in the repo root; keep it private.)

2. Authorize and run:

```bash
go run main.go            # defaults to INBOX
go run main.go "TEST"    # use a different label name
```

On first run the program will print an authorization URL. Open it, allow access, then paste
the authorization code into the terminal. A `token.json` file will be created for future runs.

Security

- `credentials.json` and `token.json` are listed in `.gitignore` and should never be committed.
- Keep these files out of source control and rotate credentials in Google Cloud if exposed.

Files

- `main.go` — main application logic that creates a Gmail client and lists messages.
- `mail/mail.go` — helpers: `GetLabelID`, `GetHeader`, `GetBodyText`.
- `token/token.go` — OAuth helpers: `GetClient`, token storage and refresh handling.

License

This repo is provided as-is.