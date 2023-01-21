package util

import rts "github.com/ory/keto/proto/ory/keto/relation_tuples/v1alpha2"

type KetoClients struct {
	WriteClient rts.WriteServiceClient
	ReadClient  rts.ReadServiceClient
	CheckClient rts.CheckServiceClient
}
