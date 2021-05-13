package ewelink

import (
	"net/url"
	"strconv"
	"time"
)

func createDeviceQuery(session *Session) *url.Values {
	query := &url.Values{}
	query.Add("lang", session.User.Language)
	query.Add("getTags", "1")
	query.Add("version", session.Application.Version)
	query.Add("ts", strconv.FormatInt(time.Now().Unix(), 10))
	query.Add("appid", session.User.AppID)

	return query
}
