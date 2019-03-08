package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"
)

type Video struct {
	Vurl  string `json:"v_url"`
	Title string `json:"title"`
}

type Data struct {
	Total int     `json:"total"`
	Items []Video `json:"items"`
}

type Resp struct {
	ErrNo  int    `json:"errno"`
	ErrMsg string `json:"errmsg"`
	Data   Data
}

var success int
var fail int

func init() {
	success = 0
	fail = 0
}

func main() {

	page := 1
	count := 20

	for {

		api := fmt.Sprintf("https://vod.gate.panda.tv/api/hostvideos?token={TOKEN}&hostid={HOSTID}&pageno=%d&pagenum=%d&__plat=pc_web&_=%d", page, count, time.Now().UnixNano())

		body, err := httpGet(api, nil)

		if err != nil {
			panic(err)
		}

		var resp Resp

		err = json.Unmarshal(body, &resp)

		if err != nil {
			panic(err)
		}

		if resp.ErrNo != 0 {
			fmt.Errorf("%s\n", resp.ErrMsg)
			break
		}

		if len(resp.Data.Items) <= 0 {
			fmt.Println("完成了！")
			break
		}

		for _, video := range resp.Data.Items {
			fmt.Println("Starting...", video.Title, video.Vurl)

			download(video.Vurl, video.Title)
		}

		page++
	}

	select {}
}

func download(u string, title string) {

	c := fmt.Sprintf("./m3u8.sh %s ./video/%s", u, title)
	cmd := exec.Command("sh", "-c", c)
	_, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error(), title, "失败")
		success++
	} else {
		fmt.Println(title, "成功")
		fail++
	}
}

func httpGet(api string, param map[string]interface{}) ([]byte, error) {

	queryStr, err := build(param)

	if err != nil {
		return nil, err
	}

	apiInfo, err := url.Parse(api)

	if err != nil {
		return nil, err
	}

	if apiInfo.RawQuery == "" {
		api = fmt.Sprintf("%s?%s", api, queryStr)
	} else {
		api = fmt.Sprintf("%s&%s", api, queryStr)
	}

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	resp, err := http.Get(api)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func build(raw map[string]interface{}) (string, error) {

	p := make(map[string]string, 0)

	for k, v := range raw {

		switch vv := v.(type) {
		case []interface{}:

			parseNormal(p, vv, []string{k})

			break
		case map[string]interface{}:

			parseKeyValue(p, vv, []string{k})

			break
		default:

			p[k] = fmt.Sprintf("%s", vv)

			break
		}
	}

	data := url.Values{}

	for k, v := range p {
		data.Add(k, v)
	}

	return data.Encode(), nil
}

func parseKeyValue(p map[string]string, raw map[string]interface{}, keys []string) {

	for k, v := range raw {
		switch vv := v.(type) {
		case []interface{}:

			tmpKeys := append(keys, k)

			parseNormal(p, vv, tmpKeys)

			break
		case map[string]interface{}:

			tmpKeys := append(keys, k)

			parseKeyValue(p, vv, tmpKeys)

			break
		default:

			//keys = append(keys, k)

			var tmp []string

			for m, n := range keys {
				if m > 0 {
					n = fmt.Sprintf("[%s]", n)
				}

				tmp = append(tmp, n)
			}

			kStr := strings.Join(tmp, "")

			p[fmt.Sprintf("%s[%s]", kStr, k)] = fmt.Sprintf("%s", vv)

			break
		}
	}
}

func parseNormal(p map[string]string, raw []interface{}, keys []string) {

	for k, v := range raw {
		switch vv := v.(type) {
		case []interface{}:

			tmpKeys := append(keys, fmt.Sprintf("%d", k))

			parseNormal(p, vv, tmpKeys)

			break
		case map[string]interface{}:

			tmpKeys := append(keys, fmt.Sprintf("%d", k))

			parseKeyValue(p, vv, tmpKeys)

			break
		default:

			//keys = append(keys, fmt.Sprintf("%d", k))

			var tmp []string

			for m, n := range keys {
				if m > 0 {
					n = fmt.Sprintf("[%s]", n)
				}

				tmp = append(tmp, n)
			}

			kStr := strings.Join(tmp, "")

			p[fmt.Sprintf("%s[%d]", kStr, k)] = fmt.Sprintf("%s", vv)

			break
		}
	}
}
