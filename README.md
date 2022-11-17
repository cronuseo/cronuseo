<p align="center">
<img width="467" alt="Screenshot 2022-11-05 at 14 41 11" src="https://user-images.githubusercontent.com/43197743/201187748-a5af0870-8e49-4313-b01d-cc59d08f76c6.png">
</p>


## Let's make authorization easy!


![lint](https://github.com/shashimalcse/cronuseo/actions/workflows/golangci-lint.yml/badge.svg)


Lets start again!



- Organization is kind of similar to tenant
- Each an every organization has users 
- Users can add to a group that is under a organiztion
- Each user or group has a role which is under a organiztion
- Resource can be create under the organization 
- Each resource hold actions (pernissions)
- Customer can create roles and assign them to users and groups
- Each reosurce has attributes 
- Each user or group has attributes 
- We only keep control plain 

### phase 1
 - design database
 - build all entity cruds (without attribute/groups)
 - keto/openfga connect
