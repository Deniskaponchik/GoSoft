package my

import (
	"context"
	unifiv1 "github.com/deniskaponchik/GoSoft/pkg/grpc/unifi/v1"
)

func (g *GrpcServ) GetClient(ctx context.Context, req *unifiv1.ClientRequest) (*unifiv1.ClientResponse, error) {

	g.slogger.Info("")
	clientHostname := req.Hostname
	g.slogger.Info(clientHostname)

	client := g.urest.GetClientForRest(clientHostname)
	if client != nil {
		g.slogger.Info("client was found in map")
		anomaliesStructSlice := client.SliceAnomalies
		//lenSlice := len(anomaliesStructSlice)
		pbAnomSlice := make([]*unifiv1.Anomaly, len(anomaliesStructSlice))

		for numb, anomStruct := range client.SliceAnomalies {
			//pbUnifiv1Anomaly := &unifiv1.Anomaly{}
			pbAnomSlice[numb] = &unifiv1.Anomaly{
				ApName:   anomStruct.ApName,
				DateHour: anomStruct.DateHour,
				AnomStr:  anomStruct.SliceAnomStr,
			}
		}
		//anoms := []*unifiv1.Anomaly{		ApName: "",	}

		return &unifiv1.ClientResponse{
			Hostname:  client.Hostname, //"взять из Unifi.Rest"
			Error:     "ошибок нет",
			Anomalies: pbAnomSlice,
			//Anomalies: []*unifiv1.Anomaly{			ApName: "",		},
		}, nil
	} else {
		return &unifiv1.ClientResponse{
			Hostname:  clientHostname,
			Error:     "клиент в базе не найден",
			Anomalies: nil,
		}, nil
	}

}
