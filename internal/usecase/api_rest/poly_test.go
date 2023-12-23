package api_rest

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"net/http"
	"testing"
	"time"
)

func TestApiSafeRestart(t *testing.T) {
	polyWebApi := &PolyWebAPI{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		polyUserName: "", //УКАЗАТЬ логин на время тестирования
		polyPassword: "", //УКАЗАТЬ пароль на время тестирования
	}

	polyStruct := &entity.PolyStruct{
		IP: "", //УКАЗАТЬ ip на время тестирования
	}
	url := "http://" + polyStruct.IP + "/api/v1/mgmt/lineInfo"
	t.Logf(url)

	err := polyWebApi.ApiSafeRestart(*polyStruct)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("Успех")
	}

}
