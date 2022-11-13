package repositories

import (
	"github.com/shashimalcse/Cronuseo/config"
	"github.com/shashimalcse/Cronuseo/models"
)

func GetGroups(tenant_id string, groups *[]models.Group) error {
	err := config.DB.Select(groups, "SELECT * FROM tenant_group WHERE tenant_id = $1", tenant_id)
	if err != nil {
		return err
	}
	return nil
}

func GetGroup(tenant_id string, id string, group *models.Group) error {
	err := config.DB.Get(group, "SELECT * FROM tenant_group WHERE tenant_id = $1 AND group_id = $2", tenant_id, id)
	if err != nil {
		return err
	}
	return nil
}

func CreateGroup(tenant_id string, group *models.Group) error {
	stmt, err := config.DB.Prepare("INSERT INTO tenant_group(group_key,name,tenant_id) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(group.Key, group.Name, tenant_id)
	if err != nil {
		return err
	}
	return nil
}

func DeleteGroup(tenant_id string, id string) error {
	stmt, err := config.DB.Prepare("DELETE FROM tenant_group WHERE tenant_id = $1 AND group_id = $2")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(tenant_id, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGroup(group *models.Group) error {
	stmt, err := config.DB.Prepare("UPDATE tenant_group SET name = $1 WHERE tenant_id = $2 AND group_id = $3")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(group.Name, group.TenantID, group.ID)
	if err != nil {
		return err
	}
	return nil
}

func CheckGroupExistsById(id string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT group_id FROM tenant_group WHERE group_id = $1)", id).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}

func CheckGroupExistsByKey(tenant_id string, key string, exists *bool) error {
	err := config.DB.QueryRow("SELECT exists (SELECT group_key FROM tenant_group WHERE tenant_id = $1 AND group_key = $2)",
		tenant_id, key).Scan(exists)
	if err != nil {
		return err
	}
	return nil
}
