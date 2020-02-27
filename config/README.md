Contains config loading code and config files. 

1. Secrets are injected via environment variables: eg. API keys, DB passwords
2. Service config is read from config files per each environment: development, staging, production etc, determined by `env` variable. Eg. feature flags, API rate limits, DB/Cache endpoints