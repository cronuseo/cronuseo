package util

type MongoDBConfig struct {
	DBName                     string `json:"db_name"`
	OrganizationCollectionName string `json:"organization_collection_name"`
}
