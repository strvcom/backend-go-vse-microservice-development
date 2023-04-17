OAuth Flow:
- Create `OAuth App` in GitHub developer settings.
- Fill in `clientSecret` in `main.go`.
- Visit https://github.com/login/oauth/authorize?client_id=<client_id>, where `<client_id>` is `id` of your GitHub OAuth app.
- Copy the `code` value from redirected URL and fill in `code` in `main.go`.
- Run `main.go` and see logs.
