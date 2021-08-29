/**
 * @Author: zhangwenbing@u51-inc.com
 * @Date: 2021/7/11 上午11:26
 * @Desc:
 */
package main

import (
	"fmt"
	"github.com/XiaoMi/soar/advisor"
	"github.com/gin-gonic/gin"
	"github.com/percona/go-mysql/query"
	"strings"
)

// 1. go里面没有类的概念
// 2. go有点像进阶版的C语言

// 声明一个struct
type SQLRequest struct {
	SQL string `json:"sql"` // 字段的第一个字母要大写
}

type SQLResponse struct {
	Fingerprint string `json:"fingerprint"`
	Id string `json:"id"`
	SQL string `json:"sql"`
	Suggests map[string]advisor.Rule `json:"suggests"`
}

func main() {
	fingerprint := strings.TrimSpace(query.Fingerprint("select * from test where c like '%b%'"))
	fmt.Println(fingerprint)

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Geektutu")
	})

	r.POST("/", func(context *gin.Context) {
		json := SQLRequest{}
		context.BindJSON(&json)
		if len(json.SQL) == 0 {
			context.String(400, "sql can not be blank")
			return
		}
		fingerprint := strings.TrimSpace(query.Fingerprint(json.SQL))
		id := query.Id(fingerprint)
		resp := SQLResponse{
			Fingerprint: fingerprint,
			Id:          id,
			SQL: json.SQL,
		}
		q, syntaxErr := advisor.NewQuery4Audit(json.SQL)
		if syntaxErr != nil {
			context.String(400, fmt.Sprintf("sql syntaxErr: %s", syntaxErr.Error()))
			return
		}

		heuristicSuggest := make(map[string]advisor.Rule)

		for item, rule := range advisor.HeuristicRules {
			// 去除忽略的建议检查
			okFunc := (*advisor.Query4Audit).RuleOK
			if !advisor.IsIgnoreRule(item) && &rule.Func != &okFunc {
				r := rule.Func(q)
				if r.Item == item {
					heuristicSuggest[item] = r
				}
			}
		}

		resp.Suggests = heuristicSuggest

		context.JSON(200, resp)
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}
