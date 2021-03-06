package admin

import (
	"fmt"

	"github.com/kim-sardine/kfadmin/client"
)

// CreateUser create kubeflow static password user
func CreateUser(email, password string) {

	username := getUsernameFromEmail(email)

	cm := c.GetConfigMap("auth", "dex")
	originalData := cm.Data["config.yaml"]
	dc := client.UnmarshalDexConfig(originalData)
	users := dc.StaticPasswords

	uuids := make([]string, len(users)+1)
	for _, user := range users {
		uuids = append(uuids, user.UserID)
	}

	newUser := client.StaticPasswordManifest{
		Email:    email,
		Hash:     hashPassword(password),
		Username: username,
		UserID:   getUniqueUUID(uuids),
	}
	fmt.Println(newUser)

	dc.StaticPasswords = append(dc.StaticPasswords, newUser)
	cm.Data["config.yaml"] = client.MarshalDexConfig(dc)

	err := c.UpdateConfigMap("auth", "dex", cm)
	if err != nil {
		panic(err)
	}

	err = c.RestartDexDeployment()
	if err != nil {
		fmt.Println("restart failed")
		fmt.Println(err)
		fmt.Println("rollback dex")

		cm = c.GetConfigMap("auth", "dex")
		cm.Data["config.yaml"] = originalData
		err := c.UpdateConfigMap("auth", "dex", cm)
		if err != nil {
			panic(err)
		}
		fmt.Println("user creation failed.")
		return
	}
	fmt.Println("User Created")
}
