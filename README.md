*This document is a model for how to order the sections of your README.  This format is _highly_ encouraged for structuring your README to keep documentation as consistent as possible. A good living example of this document is: https://github.com/shipt/shipt-search.*

# Tempest Template
A template project for new services using [tempest](https://github.com/shipt/tempest). 

You can create a new project using this template via the + button in the top right corner of the [kubedashian](https://kubedashian.shipt.com) landing page.


The service in this repo implements a simple CRUD API backed by a Postgres database.  See the [Developer Details](#developer-details) section for insturctions on how to run the service locally. 

## Project Structure

```
  cmd
  └── tempest-template
      └── service entrypoint and cobra-based CLI
  internal
  ├── customers
  │   ├── customer domain models and postgres db layer
  ├── test
  │   └── unit/integration test harness
  ├── webserver
  │   ├── customers
  │   │   ├── http.Handler implementing a simple CRUD api for customers
  │   ├── health
  │   │   └── http.Handler implementing a service liveness/readiness api
  │   └── a basic webserver with a couple of handlers
  └── worker
      └── a placeholder worker, t.b.d. kafka consumer example
  migrations
  ├── database migrations
  public
  └── doc
      ├── api documentation for the customers CRUD api
```

## Table of Contents

- [Tempest Template](#tempest-template)
  - [Project Structure](#project-structure)
  - [Table of Contents](#table-of-contents)
  - [Maintainers](#maintainers)
  - [Slack Channel](#slack-channel)
  - [Overview](#overview)
    - [Architecture](#architecture)
    - [Dependencies](#dependencies)
    - [Terminology](#terminology)
  - [Deployment and Configuration](#deployment-and-configuration)
    - [Metrics and Telemetry](#metrics-and-telemetry)
    - [Logging](#logging)
    - [Analytics](#analytics)
    - [Versioning](#versioning)
    - [CI/CD](#cicd)
- [Developer Details](#developer-details)
  - [Installation](#installation)
      - [Go](#go)
      - [Development Environment](#development-environment)
  - [File Linting](#file-linting)
  - [Git, Pull Requests and Reviews Process](#git-pull-requests-and-reviews-process)
  - [Testing](#testing)

## Maintainers
A list of names of people who actively maintain, review PRs, and can support this code.

## Slack Channel

The channel you'd prefer people contact your team at in case of issues. 

## Overview
A paragraph or two describing your thing in appropriate detail for a fairly wide (non-engineering) audience. Include code names here.  

### Architecture
As this applies, include information about how the system is structured. Include a diagram or a link to one if at all possible. A picture is worth a thousand words. I encourage you to use [Lucidchart](https://www.lucidchart.com) as we have an enterprise account for this. 

Ex:

![Microservice Architecture October 2019 (2)](https://user-images.githubusercontent.com/42652171/66686504-4bcaef00-ec45-11e9-9292-8427a3896a31.jpeg)

### Dependencies
List any service dependencies that this thing has - Shipt or non-Shipt. Optimally these are also covered in the diagram above. They should also be listed [here](https://github.com/shipt/TechHub/blob/master/content/services/a-guide-to-service-repos.md) under your service.

### Terminology
Ideally you have links throughout this document to terminology already described in the [Shipt Dictionary](https://github.com/shipt/TechHub/blob/master/content/culture/tech-culture/shipt-dictionary.md). If an explanation of terms is required beyond that, do so here.

## Deployment and Configuration

### Environmental Variables
As applies

### Metrics and Telemetry
As applies

### Logging
As applies

### Analytics
As applies
Ex: [Segway Analytics](https://github.com/shipt/segway##analytics)

### Versioning
Describe the versioning strategy, how to create a version or a release, and how different versions are maintained (if supporting multiple versions).

### CI/CD
- Where are the tests executed (drone, jenkins, codeship, circle ci, other hosted solution, other in house solution)? 
- When are they executed? 

# Developer Details

## Installation
- Golang 1.16.3 - [Installation](https://golang.org/doc/install)
- Docker Desktop - [Installation](https://www.docker.com/products/docker-desktop)
- redoc-cli - To install, run: `npm install -g redoc-cli`
- golangci-lint - To install, run: `brew install golangci-lint`
- golang-migrate - To install, run `brew install golang-migrate`

#### Go

If you haven't already installed the Go programming language, you'll want to do so.  Follow the instructions [here](https://golang.org/doc/install)

#### Development Environment

The first step in getting your local development environment running is to spin up local versions of external dependencies.  We use `docker-compose` to manage these, and you can bring them up by running the following command:

```
$ docker-compose up -d
```

Now you should be able to run the webserver.  You can verify that the webserver and database are live by hitting the service's readiness endpoint with `curl`.  

```
$ go run cmd/tempest-template/main.go webserver
...
$ curl -v http://localhost:8080/health/ready
Trying ::1...
TCP_NODELAY set
Connected to localhost (::1) port 8080 (#0)
GET /health/ready HTTP/1.1
Host: localhost:8080
User-Agent: curl/7.64.1
Accept: */*

HTTP/1.1 200 OK
Date: Tue, 25 Aug 2020 16:32:14 GMT
Content-Length: 0
```

VS Code users can also add the JSON blob below to `.vscode/launch.json` in the root of the repository for integrated debugging support. Anything in the `.vscode` directory will be ignored by git. With the VS Code launch configuration in place, press <kbd>F5</kbd> to start the service in the debugger. While the service is running, press <kbd>Shift</kbd>+<kbd>F5</kbd> to stop the service or <kbd>Ctrl</kbd>+<kbd>Shift</kbd>+<kbd>F5</kbd> to rebuild and restart the service.
```
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "cwd": "${workspaceFolder}",
      "program": "${workspaceFolder}/cmd/tempest-template/main.go",
      "args": [
        "webserver"
      ]
    }
  ]
}
```

## File Linting

We recommend you use `golangci-lint` for source code linting. A linting pass will be run as part of the CI pipeline on any branch push.

## Git, Pull Requests and Reviews Process

Before beginning any development task, be sure to create a Clubhouse Card in the . All development should be done in feature branches off of the `HEAD` of `main`, ideally with names that can be automatically tracked by Clubhouse such as `johncarmack/ch123456/all-the-features`. When you are satisfied with your changes, you must open a pull request and solicit the wisdom of at least one approving review prior to merging to `main`.

## Testing

From the root directory of the repository, you can run the entire test suite by executing `make test`. Always strive to _increase_ test coverage when contributing to the project.

## Migrations

If you're running `docker-compose up -d` migrations will be run for you. If you'd prefer to run migrations locally without Docker you can run them manually with the following:

Apply migrations:
```bash
make migrate
```

Create new migration:
```bash
make migration NAME=[MIGRATION_NAME]
```

See `migrations.mk` for details.

## Recommendations
* The phrase `tempest-template` is everywhere (directories, in the code, in infraspec).  Change those to be your own server's name.
* If you will be having a postgres database, you need to request a staging DB
* DBAs will not create a production postgres DB until you understand the characteristics of your service and have a general
idea of the needed database t-shirt sizing
* If you will not have a postgres DB, in addition to the code, be sure to adjust the .drone.yml, Makefile, docker-compose.yml, and infraspec.yaml files
* If you will not be using kafka, be sure to adjust the Makefile, docker-compose.yml, and infraspec.yaml files
* If you are using the '+' on the https://kubedashian.shipt.com/ dashboard to create your server, a github repo is created for you.  We suggest you modify your created github repo to
  * Update CODEOWNERS file
  * Add a file under .github directory called pull_request_template.md
  * Add shipt/Developers as a team that can access the repo under Settings/Manage access 
  * Enable webhook merge restrictions to branch main as follows:
    * Require status checks to pass before merging
    * Require branches to be up to date before merging
    * Add status checks for `captain-hook/infraspec_validation` and `continuous-integration/drone/push`
* Setup [sonarqube](https://sonarqube.gcp.shipttech.com) to scan new builds, see the [sonarqube setup guide](https://techhub.shipt.com/engineering/infrastructure/devops/sonarqube/setup/) on techhub.

## Todos for merging to main and having a staging and production env running
* Download the latest `platformctl` tool from here: https://github.com/shipt/platformctl/releases
* Install `platformctl` using the readme at that github repo
* `PSQL_URL` will be automatically injected based the CloudSQL DB defined in Infraspec
* In Drone, create `STAGING_DATABASE_URL` and `PRODUCTION_DATABASE_URL` as a setting/secret in drone UI.  This is how migrations would work then deploying to staging and prod.
* You'll need at least one push to main to get staging kubernetes all up and going.  Check slack channel sre-beholder-stag-err to ensure
  your env is up and running
* Before pushing to staging (via a merge to main), be sure to comment or remove this line in infraspec.yaml file
```yaml
  environments:
    - name: staging
      disabled: true # comment this out when you are ready to push to staging

* Before pushing to production (via creating a git tag), be sure to comment or remove this line in infraspec.yaml file
```yaml
  environments:
    - name: production
      disabled: true # comment this out when you are ready to push to production
```