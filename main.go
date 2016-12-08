package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

//物流状态：2-在途中,3-签收,4-问题件

const (
	EBusinessID = ""
	AppKey      = ""
	ReqUrl      = "http://api.kdniao.cc/Ebusiness/EbusinessOrderHandle.aspx"

	SHIP_CODE_YUNDA    = "YD"
	SHIP_CODE_SHUNFENG = "SF"
)

type RequestData struct {
	ShipperCode  string
	LogisticCode string
}

type PostParams struct {
	RequestData string
	EBusinessID string
	RequestType string
	DataSign    string
	DataType    string
}

type TraceItem struct {
	AcceptTime    string
	AcceptStation string
}

type TraceResult struct {
	Traces  []TraceItem
	Success bool
	State   string
}

func main() {
	KdnTraces(SHIP_CODE_YUNDA, "xxxx")
}

func KdnTraces(shipperCode string, logisticCode string) (traceResult *TraceResult, err error) {
	if AppKey == "" || EBusinessID == "" {
		fmt.Println("Please fill AppKey & EBusinessID")
		return nil, nil
	}

	if requestDataJson, err := json.Marshal(&RequestData{
		ShipperCode:  shipperCode,
		LogisticCode: logisticCode,
	}); err != nil {
		fmt.Printf("KdnTraces request to json error:%v\n", err)
		return nil, err
	} else {
		md5Ctx := md5.New()
		md5Ctx.Write([]byte(string(requestDataJson) + AppKey))
		b64 := base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(md5Ctx.Sum(nil))))

		resp, err := http.PostForm(ReqUrl,
			url.Values{
				"RequestData": {url.QueryEscape(string(requestDataJson))},
				"EBusinessID": {EBusinessID},
				"RequestType": {"1002"},
				"DataSign":    {url.QueryEscape(b64)},
				"DataType":    {"2"},
			})
		defer resp.Body.Close()

		if err != nil {
			fmt.Printf("KdnTraces post error:%v\n", err)
			return nil, err
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("KdnTraces read body error:%v\n", err)
				return nil, err
			}

			fmt.Println(string(body))
			// Parser body
			traceResult := TraceResult{}
			json.Unmarshal(body, &traceResult)

			fmt.Printf("Trace result:%v\n", traceResult)

			return &traceResult, nil
		}
	}
}
