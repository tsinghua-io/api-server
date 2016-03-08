package old

import (
	"github.com/tsinghua-io/api-server/resource"
	"net/http"
)

// func (adapter *OldAdapter) PersonalInfo() (*resource.User, int) {
// 	url := "/MultiLanguage/vspace/vspace_userinfo1.jsp"
// 	doc, err := adapter.getOldResponse(url, make(map[string]string))

// 	if err != nil {
// 		glog.Warningf("Failed to get response from learning web: %s", err)
// 		return nil, http.StatusBadGateway
// 	} else {
// 		// parsing the response body
// 		docTable := doc.Find("form")
// 		infos := docTable.Find(".tr_l,.tr_l2").Map(func(i int, valueTR *goquery.Selection) (info string) {
// 			switch valueTR.Nodes[0].FirstChild.Type {
// 			case html.TextNode:
// 				info, _ = valueTR.Html()
// 			case html.ElementNode:
// 				info, _ = valueTR.Children().Attr("value")
// 			default:
// 				info, _ = valueTR.Html()
// 			}
// 			return
// 		})

// 		if len(infos) < 15 {
// 			glog.Errorf("User information parsing error: cannot parse all the informations from %s", infos)
// 			return nil, http.StatusBadGateway
// 		} else {
// 			return &resource.User{
// 				Id:     infos[0],
// 				Name:   infos[1],
// 				Type:   infos[14],
// 				Gender: infos[13],
// 				Email:  infos[6],
// 				Phone:  infos[7]}, http.StatusOK
// 		}
// 	}
// }

func (ada *OldAdapter) Profile(_ string, _ map[string]string) (user *resource.User, status int) {
	return nil, http.StatusNotImplemented
}
