// Copyright 2016 The StudyGolang Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://studygolang.com
// Author: polaris	polaris@studygolang.com

package admin

import (
	"encoding/json"
	"net/http"

	"sander/db/nosql"
	xhttp "sander/http"
	"sander/logger"

	"github.com/labstack/echo"
	"github.com/polaris1119/goutils"
)

func parsePage(ctx echo.Context) (curPage, limit int) {
	curPage = goutils.MustInt(ctx.FormValue("page"), 1)
	limit = goutils.MustInt(ctx.FormValue("limit"), 20)
	return
}

func parseConds(ctx echo.Context, fields []string) map[string]string {
	conds := make(map[string]string)

	for _, field := range fields {
		if value := ctx.FormValue(field); value != "" {
			conds[field] = value
		}
	}

	return conds
}

// render html 输出
func render(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	return xhttp.RenderAdmin(ctx, contentTpl, data)
}

func renderQuery(ctx echo.Context, contentTpl string, data map[string]interface{}) error {
	return xhttp.RenderQuery(ctx, contentTpl, data)
}

func success(ctx echo.Context, data interface{}) error {
	result := map[string]interface{}{
		"ok":   1,
		"msg":  "操作成功",
		"data": data,
	}

	b, err := json.Marshal(result)
	if err != nil {
		return err
	}

	go func(b []byte) {
		if cacheKey := ctx.Get(nosql.CacheKey); cacheKey != nil {
			nosql.DefaultLRUCache.CompressAndAdd(cacheKey, b, nosql.NewCacheData())
		}
	}(b)

	if ctx.Response().Committed() {
		return nil
	}

	return ctx.JSONBlob(http.StatusOK, b)
}

func fail(ctx echo.Context, code int, msg string) error {
	if ctx.Response().Committed() {
		return nil
	}

	result := map[string]interface{}{
		"ok":    0,
		"error": msg,
	}

	logger.Error("operate fail:%+v", result)

	return ctx.JSON(http.StatusOK, result)
}
