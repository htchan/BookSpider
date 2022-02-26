package main

import (
	"net/http"
	"sort"
	"strings"
	"strconv"
)

func GeneralInfo(res http.ResponseWriter, req *http.Request) {
	setHeader(res)

	siteNames := make([]string, len(siteMap))
	i := 0
	for siteName := range siteMap {
		siteNames[i] = siteName
		i++
	}
	sort.Strings(siteNames)
	//TODO: add the working sites into response in future
	response(res, map[string]interface{} { "siteNames": siteNames, "availableSiteNames": []string{} })
}

func SiteInfo(res http.ResponseWriter, req *http.Request) {
	setHeader(res)

	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[4]
	site, ok := siteMap[siteName]
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	}

	response(res, site.Map())
}

func BookInfo(res http.ResponseWriter, req *http.Request) {
	setHeader(res)

	uri := strings.Split(req.URL.Path, "/")
	siteName := uri[4]
	site, ok := siteMap[siteName]
	if !ok {
		error(res, http.StatusNotFound, "site <" + siteName + "> not found")
		return
	}
	if len(uri) < 6 {
		error(res, http.StatusBadRequest, "not enough parameters")
		return
	}
	id, err := strconv.Atoi(uri[5])
	if err != nil {
		error(res, http.StatusBadRequest, "id <" + uri[5] + "> is not a number")
		return
	}

	hashCode := "";
	if len(uri) > 6 {
		hashCode = uri[6]
	}

	book := site.SearchByIdHash(id, hashCode)
	if book == nil || book.GetTitle() == "" {
		error(res, http.StatusNotFound, "book <" + strconv.Itoa(id) + ">, " +
				"hash <" + hashCode + "> in site <" + siteName + "> not found")
	} else {
		response(res, book.Map())
	}
}