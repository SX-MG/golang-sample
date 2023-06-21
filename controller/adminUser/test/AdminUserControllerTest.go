package adminUser

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type AdminUserControllerTest struct {
}

func (AdminUserControllerTest) Test_List() {
	// parment := controllerUtil.Paramet{
	// 	Start:       1,
	// 	Limit:       10,
	// 	QueryString: "admin",
	// 	OrderField:  "f_userName",
	// 	OrderType:   "DESC",
	// }
	data := url.Values{"start": {"0"}, "offset": {"xxxx"}}
	body := strings.NewReader(data.Encode())
	clt := http.Client{}
	resp, err := clt.Post("http://localhost:8088/common/doUserLogout", "application/x-www-form-urlencoded", body)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err.Error())
		}
		respBody := string(content)
		fmt.Println("respBody:", respBody)
	}

}
