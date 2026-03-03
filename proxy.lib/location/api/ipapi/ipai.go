package ipapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	base "github.com/lancelotXie/proxy/proxy.lib/location/base"

	"github.com/pkg/errors"
)

// Resolv : 从 ip-api.com 网站的api 得到 ip 对应的 所在地
func Resolv(ip string) (location base.Location, err error) {
	url := "http://ip-api.com/json/" + ip
	mapBytes, err := getMap(url)
	if err != nil {
		return
	}
	location, err = resolvMap(mapBytes)
	return
}

// 从网站获取 IP所在地信息，并将信息 转换为 map[string]interface{}
func getMap(url string) (out map[string]interface{}, err error) {
	res, err := http.Get(url)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	defer res.Body.Close()
	info, err := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(info, &out)

	return
}

// 从 map 里得到 IP 对应的所在地，再转换成 base.Location 类型返回
func resolvMap(resMap map[string]interface{}) (location base.Location, err error) {
	_loc, ok := resMap["country"]
	if !ok {
		return
	}
	switch string(_loc.(string)) {
	case base.China.String():
		location = base.China
	case base.Japan.String():
		location = base.Japan
	default:
		location = base.Others
	}
	return
}
