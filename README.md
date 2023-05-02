<h1 align="center"><img src="https://user-images.githubusercontent.com/43197743/233458042-c0b08684-87fa-471b-8f13-5b23d84ecd0a.png" alt="cronuseo - open-source authorization solution"></h1>

<p align="left">
    <a href="https://goreportcard.com/report/github.com/shashimalcse/cronuseo"><img src="https://goreportcard.com/badge/github.com/shashimalcse/cronuseo" alt="Go Report Card"></a>
</p>

## What is cronuseo ?

Cronuseo is an open-source authorization solution that allows developers to easily integrate permissions and access control into their products within minutes.

> Example: A developer can call the cronuseo and get a clear answer if User A has the permissions to create Resource B.

By using Cronuseo, it is possible to create a separate authorization service, effectively separating our policies from our code. Controlling access management centrally via a separate authorization service enables you to provide this service to any system that needs to check whether a user can or cannot access its resources. Cronuseo is based on modern, open-source foundation which includes [Open Policy Agent (OPA)](https://www.openpolicyagent.org/), [Zanzibar](https://research.google/pubs/pub48190/).

## Main features:

* Role-based Access Control (RBAC)

## Get started

> Note : cronuseo still in the experimental stage. Only tested in the local environemnt.

You can use Docker to run cronuseo locally
### Set up the local environment

* ``` curl -LJO https://raw.githubusercontent.com/shashimalcse/cronuseo/HEAD/docker-compose-db.yml | curl -LJO https://raw.githubusercontent.com/shashimalcse/cronuseo/HEAD/docker-compose.yml ```
* Prepare a [mongodb](https://hub.docker.com/_/mongo) instance ``` docker compose -f docker-compose-database.yml up```
* Make sure to update the necessary configuration in the `config/local.yml` file, and don't forget to replace the jwks endpoint with the ones provided by your own identity provider. (only tested with [asgardeo](https://wso2.com/asgardeo/))
* Start management server and check server (Policy Decision Point) ``` docker compose up --build```

## How to implement RBAC using cronuseo

In order to use RBAC, we need two types of information:
- Which `users/groups` have which `roles`
- Which `roles` have which `permissions`

Once we provide RBAC with this information, we decide how to make an authorization decision; A user/group is assigned to a role and is authorized to do what the role permits. 

For example, let us look at the following role assignments:

| User/Group (Who is performing the action) | Role (What is the user’s/group's assigned role) |
| :---:   | :---: |
| John | Director   |
| Bob | Coordinator |
| Finance Department | Accountant |

> Use user and group endpoints to create users and groups

```
curl --location --request POST 'localhost:8080/api/v1/<org_id>/user' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: <Token> \
--data-raw '{
  "first_name": "First Name",
  "last_name": "Last Name",
  "username": "Username"
}'
```

```
curl --location --request POST 'localhost:8080/api/v1/<org_id>/group' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: <Token> \
--data-raw '{
  "display_name": "Display Name",
  "identifier": "Identifier",
  "users": [
    "<user_id>"
  ]
}'
```

> Use role endpoint to assign roles to users/groups

```
curl --location --request POST 'localhost:8080/api/v1/<org_id>/role' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: <Token> \
--data-raw '{
  "display_name": "Display Name",
  "identifier": "Identifier",
  "groups": [
    "<group_id>"
  ],
  "users": [
    "<user_id>"
  ]
}'
```

And this role/permission assignment: 

| Role | Action (What are they doing) | Resource (What are they performing the action on) |
| :---:   | :---: | :---: |
| Director | Write | Budget |
| Coordinator | Read | Budget |
| Accountant | Want | Budget |

> Use resource `POST` request to create a resource with actions

```
curl --location --request POST 'localhost:8080/api/v1/<org_id>/resource' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: <Token> \
--data-raw '{
  "actions": [
    {
      "display_name": "Display Name",
      "identifier": "Identifier",
    }
  ],
  "display_name": "Display Name",
  "identifier": "Identifier",
}'
```

> Use role permission `PATCH` request to assign permission to role

```
curl --location --request PATCH 'localhost:8080/api/v1/<org_id>/role/<role_id>/permission' \
--header 'Content-Type: application/json' \
--header 'Accept: application/json' \
--header 'Authorization: <Token> \
--data-raw '{
  "added_permissions": [
    {
      "action": "Action Identifier",
      "resource": "Resource Identifier"
    }
  ]
}'
```

In this example, RBAC will make the following authorization decisions:

| User/Group | Action | Resource | Decision (Should the action be allowed, and why?)|
| :---:   | :---: | :---: | :---: |
| John | Write | Budget | `Allow` because John is in Director |
| Bob | Read | Budget | `Allow` because Bob is in Coordinator |
| Bob | Write | Budget | `Deny` because Bob is in Coordinator |
| Finance Department | Write | Budget | `Deny` because Finance Department is in Accountant |

> Use permission check endpoint or SDKs to get authorization decisions

```
curl --location --request POST 'localhost:8081/api/v1/<org_identifier>/permission/check' \
--header 'Content-Type: application/json' \
--header 'API_KEY: <API_KEY>' \
--data-raw '{
  "action": "Action Identifier",
  "resource": "Resource Identifier",
  "username": "shashimal20"
}'
```

> Response will be `true` or `false`

## cronuseo SDKs for applications
use these sdks to check permissions for the user.
* python - https://pypi.org/project/cronuseosdk
* nodejs - https://www.npmjs.com/package/cronuseosdk
* golang - https://github.com/shashimalcse/cronuseogosdk
