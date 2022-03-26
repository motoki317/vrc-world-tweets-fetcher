package utils

import (
	"net/url"
	"regexp"
)

const (
	uuidFormat          = "[\\da-f]{8}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{4}-[\\da-f]{12}"
	vrchatWorldIDFormat = "wrld_" + uuidFormat
)

type worldIDExtractFunc func(u *url.URL) VRChatWorldID

func generateWorldExtractFunctions() []worldIDExtractFunc {
	path1 := regexp.MustCompile("^/home/world/(" + vrchatWorldIDFormat + ")$")
	vrchatWorldIDFormatRe := regexp.MustCompile("^" + vrchatWorldIDFormat + "$")

	return []worldIDExtractFunc{
		func(u *url.URL) VRChatWorldID {
			if u.Scheme == "https" && u.Host == "vrchat.com" {
				matches := path1.FindStringSubmatch(u.Path)
				if len(matches) == 2 {
					return VRChatWorldID(matches[1])
				}
			}
			return ""
		},
		func(u *url.URL) VRChatWorldID {
			if u.Scheme == "https" && u.Host == "vrchat.com" && u.Path == "/home/launch" {
				if wid := u.Query().Get("worldId"); vrchatWorldIDFormatRe.MatchString(wid) {
					return VRChatWorldID(wid)
				}
			}
			return ""
		},
	}
}

var vrcWorldURLFormats = generateWorldExtractFunctions()

func ExtractVRChatWorldIDs(urls []string) []VRChatWorldID {
	var ret []VRChatWorldID
	for _, urlStr := range urls {
		u, err := url.Parse(urlStr)
		if err != nil {
			continue
		}
		for _, format := range vrcWorldURLFormats {
			if wid := format(u); wid != "" {
				ret = append(ret, wid)
				break
			}
		}
	}
	return ret
}
