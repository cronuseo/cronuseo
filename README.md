<p align="center">
<img alt="Screenshot 2022-11-05 at 14 41 11" src="https://user-images.githubusercontent.com/43197743/205511091-bed0ace7-705d-4372-8980-872bcd71a200.png">
</p>

<h1 align="center">Authorization as a Service</h1>
<h2 align="center">Let's make auth-z easy!</h2>

## What is cronuseo ?

Cronuseo is an open-source authorization solution that allows developers to easily integrate permissions and access control into their products within minutes.

> Example: A developer can call the cronuseo and get a clear answer if User A has the permissions to create Resource B.

cronuseo is based on modern, open-source foundation which includes Open Policy Agent (OPA), Zanzibar.

## Main features:

* Role-based Access Control (RBAC)

## Get started

> Note : cronuseo still in the experimental stage. Only tested in the local environemnt.

* ``` git clone https://github.com/shashimalcse/cronuseo```
* ``` docker compose up --build```
* ``` docker compose -f docker-compose-database.yml up```
* ``` make run-console```
* Heading to http://localhost:3000 to continue your c6o journey.

## cronuseo SDKs for applications
use these sdks to check permisisons for the user.
* python - https://pypi.org/project/cronuseosdk
* nodejs - https://www.npmjs.com/package/cronuseosdk
* golang - https://github.com/shashimalcse/cronuseogosdk

> cronuseo uses environment variables for configuration, along with .env file support.

