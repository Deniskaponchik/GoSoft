package repo

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"testing"
)

func TestUnifiRepo_GetLoginPCerr(t *testing.T) {
	ur := &UnifiRepo{
		dataSourceITsup: "nil/nil",
		dataSourceGLPI:  "root:t2root@tcp(10.77.252.153:3306)/glpi_db",
		controller:      1,
	}
	client := &entity.Client{
		Hostname: "WSIR-BRONER",
	}
	err := ur.GetLoginPCerr(client)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf(client.UserLogin)
	}
}
