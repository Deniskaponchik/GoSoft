package main

/*
import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


func main() {
	client := NewSOAPClient("https://soap.example.com/call", true, nil)

	req := &CreateUserRequest{
		Email:    "jdoe@example.com",
		Password: "1234567890",
	}
	res := &CreateUserResponse{}
	if err := client.Call("create_user", req, res); err != nil {
		panic(err)
	}

	// if everything went well res.ID should have its
	// value set with the one returned by the service.
	fmt.Println(res.ID)
}

type CreateUserRequest struct {
	Email    string `xml:"Email,omitempty"`
	Password string `xml:"Password,omitempty"`
}
type CreateUserResponse struct {
	ID string `xml:"ID"`
}

type AccountUser struct {
	XMLName xml.Name `xml:"http://exacttarget.com/wsdl/partnerAPI AccountUser"`

	*APIObject

	AccountUserID             int32         `xml:"AccountUserID,omitempty"`
	UserID                    string        `xml:"UserID,omitempty"`
	Password                  string        `xml:"Password,omitempty"`
	Name                      string        `xml:"Name,omitempty"`
	Email                     string        `xml:"Email,omitempty"`
	MustChangePassword        bool          `xml:"MustChangePassword,omitempty"`
	ActiveFlag                bool          `xml:"ActiveFlag,omitempty"`
	ChallengePhrase           string        `xml:"ChallengePhrase,omitempty"`
	ChallengeAnswer           string        `xml:"ChallengeAnswer,omitempty"`
	UserPermissions           []*UserAccess `xml:"UserPermissions,omitempty"`
	Delete                    int32         `xml:"Delete,omitempty"`
	LastSuccessfulLogin       time.Time     `xml:"LastSuccessfulLogin,omitempty"`
	IsAPIUser                 bool          `xml:"IsAPIUser,omitempty"`
	NotificationEmailAddress  string        `xml:"NotificationEmailAddress,omitempty"`
	IsLocked                  bool          `xml:"IsLocked,omitempty"`
	Unlock                    bool          `xml:"Unlock,omitempty"`
	BusinessUnit              int32         `xml:"BusinessUnit,omitempty"`
	DefaultBusinessUnit       int32         `xml:"DefaultBusinessUnit,omitempty"`
	DefaultApplication        string        `xml:"DefaultApplication,omitempty"`
	Locale                    *Locale       `xml:"Locale,omitempty"`
	TimeZone                  *TimeZone     `xml:"TimeZone,omitempty"`
	DefaultBusinessUnitObject *BusinessUnit `xml:"DefaultBusinessUnitObject,omitempty"`

	AssociatedBusinessUnits struct {
		BusinessUnit []*BusinessUnit `xml:"BusinessUnit,omitempty"`
	} `xml:"AssociatedBusinessUnits,omitempty"`

	Roles struct {
		Role []*Role `xml:"Role,omitempty"`
	} `xml:"Roles,omitempty"`

	LanguageLocale *Locale `xml:"LanguageLocale,omitempty"`

	SsoIdentities struct {
		SsoIdentity []*SsoIdentity `xml:"SsoIdentity,omitempty"`
	} `xml:"SsoIdentities,omitempty"`
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{
		//Header:        SoapHeader{},
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	log.Println(buffer.String())

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	if soapAction != "" {
		req.Header.Add("SOAPAction", soapAction)
	}

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: s.tls,
		},
		Dial: dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
*/
