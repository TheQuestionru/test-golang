# TheQuestion - Server

# Requirements
- Go 1.8+ 
- [govendor](https://github.com/kardianos/govendor) for dependency vendoring.

# Installation

First, fork [test repo](https://github.com/TheQuestionru/test-golang) on GitHub

Clone your fork:
```bash
mkdir -p ~/theq/src/github.com/TheQuestionru
cd ~/theq/src/github.com/TheQuestionru
git clone git@fork.git
```

Set GOPATH:
```bash
export GOPATH=~/theq
```

Install the dependencies and build the project:
```bash
make 
```

# Task

Add TeamCity element type to dashboard, including tests, use `newrelic` module as an example.

Provide `Dockerfile`('s) to build & deploy code to production, and a `Dockerfile` for testing.



Any questions to: `ngalayko@thequestion.ru`
