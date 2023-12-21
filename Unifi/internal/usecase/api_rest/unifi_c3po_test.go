package api_rest

import (
	"github.com/deniskaponchik/GoSoft/Unifi/internal/entity"
	"net/http"
	"testing"
	"time"
)

func TestGetUserLogin(t *testing.T) {
	unifiC3po := &UnifiC3po{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		url: "http://login:password@c3po.corp.tele2.ru/sccm/api/info/", //УКАЗАТЬ логин на время тестирования
	}
	client1 := &entity.Client{
		Hostname: "WSNS-TROFIMOV2", //УКАЗАТЬ ip на время тестирования
	}
	client2 := &entity.Client{
		Hostname: "", //УКАЗАТЬ ip на время тестирования
	}
	client3 := &entity.Client{
		Hostname: "mock", //УКАЗАТЬ ip на время тестирования
	}
	//url1 := unifiC3po.url + "pc/" + client1.Hostname
	//t.Logf(url)

	err := unifiC3po.GetUserLogin(client1)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("WSNS-TROFIMOV2 : ")
	}
	err = unifiC3po.GetUserLogin(client2)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("EMPTY : ")
	}
	err = unifiC3po.GetUserLogin(client3)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("MOCK : ")
	}

}
