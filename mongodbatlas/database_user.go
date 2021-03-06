package mongodbatlas

import (
	"fmt"
	"log"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceDatabaseUserCreate,
		Read:   resourceDatabaseUserRead,
		Update: resourceDatabaseUserUpdate,
		Delete: resourceDatabaseUserDelete,

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"database": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"database": {
							Type:     schema.TypeString,
							Required: true,
						},
						"collection": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := ma.DatabaseUser{
		Username:     d.Get("username").(string),
		Password:     d.Get("password").(string),
		DatabaseName: d.Get("database").(string),
	}

	params.Roles = readRolesFromSchema(d.Get("roles").([]interface{}))

	databaseUser, _, err := client.DatabaseUsers.Create(d.Get("group").(string), &params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB DatabaseUser: %s", err)
	}
	d.SetId(databaseUser.Username)
	log.Printf("[INFO] MongoDB DatabaseUser ID: %s", d.Id())

	return resourceDatabaseUserRead(d, meta)
}

func resourceDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	c, _, err := client.DatabaseUsers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB DatabaseUser %s: %s", d.Id(), err)
	}

	d.Set("username", c.Username)
	d.Set("database", c.DatabaseName)
	rolesMap := make([]map[string]interface{}, len(c.Roles))
	for i, r := range c.Roles {
		rolesMap[i] = map[string]interface{}{
			"role":       r.RoleName,
			"database":   r.DatabaseName,
			"collection": r.CollectionName,
		}
	}
	d.Set("roles", rolesMap)

	return nil
}

func resourceDatabaseUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)
	requestUpdate := false

	c, _, err := client.DatabaseUsers.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB DatabaseUser %s: %s", d.Id(), err)
	}

	if d.HasChange("password") {
		c.Password = d.Get("password").(string)
		requestUpdate = true
	}
	if d.HasChange("roles") {
		c.Roles = readRolesFromSchema(d.Get("roles").([]interface{}))
		requestUpdate = true
	}

	if requestUpdate {
		_, _, err := client.DatabaseUsers.Update(d.Get("group").(string), d.Id(), c)
		if err != nil {
			return fmt.Errorf("Error updating MongoDB DatabaseUser %s: %s", d.Id(), err)
		}
		return resourceDatabaseUserRead(d, meta)
	}
	return nil
}

func resourceDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB DatabaseUser destroy: %v", d.Id())
	_, err := client.DatabaseUsers.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB DatabaseUser %s: %s", d.Id(), err)
	}

	return nil
}

func readRolesFromSchema(rolesMap []interface{}) (roles []ma.Role) {
	roles = make([]ma.Role, len(rolesMap))
	for i, r := range rolesMap {
		roleMap := r.(map[string]interface{})

		roles[i] = ma.Role{
			RoleName:       roleMap["name"].(string),
			DatabaseName:   roleMap["database"].(string),
			CollectionName: roleMap["collection"].(string),
		}
	}
	return roles
}
