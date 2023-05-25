# (Draft) Contribute to cronuseo

Thanks for your interest in contributing to cronuseo. We respect the time of community contributors, so it'll be great if we can go through this guide which provides the necessary contribution information before starting your work.

# Set up the dev environment

You'll need these installed to proceed:

 - docker
 - golang

### Clone and install dependencies

Clone the repo https://github.com/shashimalcse/cronuseo in the way you like, then execute the command below in the project root:

```bash
go get ./... 
```

### Set up mongodb database

Run the command below in the project root:
```
docker compose -f docker-compose-database.yml up
```

### Start dev

To start the dev, you have to replace config file (local-debug) with your config values.

Run the command below in the project root:
```
make run CONFIG_FILE=<ROOT>/config/local-debug.yml 
```

### Commit and create pull request

We require every commit to [be signed](https://docs.github.com/en/authentication/managing-commit-signature-verification/signing-commits), and both the commit message and pull request title follow conventional commits.

If the pull request remains empty content, it'll be DIRECTLY CLOSED until it matches our contributing guideline.
