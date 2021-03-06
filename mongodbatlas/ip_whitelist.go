package mongodbatlas

import (
	"fmt"
	"log"

	ma "github.com/akshaykarle/go-mongodbatlas/mongodbatlas"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceIPWhitelist() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPWhitelistCreate,
		Read:   resourceIPWhitelistRead,
		Update: resourceIPWhitelistUpdate,
		Delete: resourceIPWhitelistDelete,

		Schema: map[string]*schema.Schema{
			"group": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"comment": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIPWhitelistCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	params := []ma.Whitelist{
		ma.Whitelist{
			CidrBlock: d.Get("cidr_block").(string),
			GroupID:   d.Get("group").(string),
			IPAddress: d.Get("ip_address").(string),
			Comment:   d.Get("comment").(string),
		},
	}

	w, _, err := client.Whitelist.Create(d.Get("group").(string), params)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB Project IP Whitelist: %s", err)
	}
	d.SetId(w[0].CidrBlock)
	log.Printf("[INFO] MongoDB Project IP Whitelist ID: %s", d.Id())

	return resourceIPWhitelistRead(d, meta)
}

func resourceIPWhitelistRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	w, _, err := client.Whitelist.Get(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error reading MongoDB Project IP Whitelist %s: %s", d.Id(), err)
	}

	d.Set("cidr_block", w.CidrBlock)
	d.Set("ip_address", w.IPAddress)
	d.Set("group", w.GroupID)
	d.Set("comment", w.Comment)

	return nil
}

func resourceIPWhitelistUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceIPWhitelistDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ma.Client)

	log.Printf("[DEBUG] MongoDB Project IP Whitelist destroy: %v", d.Id())
	_, err := client.Whitelist.Delete(d.Get("group").(string), d.Id())
	if err != nil {
		return fmt.Errorf("Error destroying MongoDB Project IP Whitelist %s: %s", d.Id(), err)
	}

	return nil
}
