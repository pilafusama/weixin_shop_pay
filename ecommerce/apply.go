package ecommerce

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"weixin_shop_pay"
	"weixin_shop_pay/tools"
)

// EcommerceApply 二级商户进件
func (c *Ecommerce) EcommerceApply(p *weixin_shop_pay.EcommerceApplyParams) (*ApplyResp, error) {

	// 请求参数
	dataJsonByte, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	// 发起请求
	urlPath := "v3/ecommerce/applyments"
	req, err := http.NewRequest("POST", "https://api.mch.weixin.qq.com/"+urlPath, bytes.NewBuffer(dataJsonByte))
	if err != nil {
		return nil, err
	}

	// 读取私钥文件
	keyByte, err := ioutil.ReadFile(c.Config.KeyPath)
	if err != nil {
		return nil, err
	}

	// 签名
	signature, err := tools.Signature(urlPath, string(dataJsonByte), string(keyByte), c.Config.SpMchID, c.Config.SerialNo)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "WECHATPAY2-SHA256-RSA2048 "+signature)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// 解析返回参数
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Println(string(respData))
	var output ApplyResp
	err = json.Unmarshal(respData, &output)
	if err != nil {
		return nil, err
	}
	return &output, nil
}

// ApplyResp 二级商户进件返回参数
type ApplyResp struct {
	ApplymentID  string `json:"applyment_id"`   // 微信支付申请单号
	OutRequestNo string `json:"out_request_no"` // 业务申请编号
}