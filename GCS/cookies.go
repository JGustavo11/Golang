package file 

import (
	"encoding/json"
	"google.golang.org/appengine"
	"net/http"
	"fmt"
	"google.golang.org/appengine/log"
	"encoding/base64"
	
)


func makeC(res http.ResponseWriter, req *http.Request, fname string) (map[string]bool, error) {
	mss := make(map[string]bool)
	cookie, _ := req.Cookie("file-names")


	if cookie != nil {
		bs, err := base64.URLEncoding.DecodeString(cookie.Value)
		if err != nil {
			return nil, fmt.Errorf("*** ERROR ***", err)
		}
		err = json.Unmarshal(bs, &mss)
		if err != nil {
			return nil, fmt.Errorf("***ERROR ***", err)
		}
	}
	
	mss[fname] = true
	bs, err := json.Marshal(mss)
	if err != nil {
		return mss, fmt.Errorf("makeC ERROR: json.Marshal: %s", err)
	}
	b64 := base64.URLEncoding.EncodeToString(bs)

	ctx := appengine.NewContext(req)
	log.Infof(ctx, "Cookie Json: %s", string(bs))

	http.SetCookie(res, &http.Cookie{
		Name: "file-names",
		Value: b64,
		})

	return mss, nil
}
