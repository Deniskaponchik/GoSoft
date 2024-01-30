package api_rest

import (
	"github.com/deniskaponchik/GoSoft/internal/entity"
	"net/http"
	"testing"
	"time"
)

func TestGetUserLogin(t *testing.T) {
	unifiC3po := &UnifiC3po{
		client: http.Client{
			Timeout: 10 * time.Second,
		},
		url: "http://1stlinesupport:z75l!@tbtI3XDP5FQpsj@c3po.corp.tele2.ru/sccm/api/info/",
	}
	client1 := &entity.Client{
		Hostname: "NBCN-GUSIKHIN1",
	}
	client2 := &entity.Client{
		Hostname: "",
	}
	client3 := &entity.Client{
		Hostname: "mock",
	}
	//url1 := unifiC3po.url + "pc/" + client1.Hostname
	//t.Logf(url)

	//пустой логин пользователя
	err := unifiC3po.GetUserLogin(client1)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf(client1.Hostname)
		if client1.UserLogin == "" {
			t.Logf("Hostname is empty")
		}
		if client1.UserLogin == " " {
			t.Logf("Hostname is one space")
		}
	}

	//Пустое имя компьютера
	err = unifiC3po.GetUserLogin(client2)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("EMPTY : ")
	}

	//Не корректное имя компьютера
	err = unifiC3po.GetUserLogin(client3)
	if err != nil {
		t.Errorf("Incorrect result. %s", err)
	} else {
		t.Logf("MOCK : ")
	}

}
