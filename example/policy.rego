package example

import future.keywords.contains
import future.keywords.if
import future.keywords.in

default allow := false

allow if {
	some permission in user_is_granted

	input.action == permission.action
	input.resource == permission.resource
}

user_is_granted contains permission if {
	some role in input.data.user_roles
	some permission in input.data.role_permissions[role]
}