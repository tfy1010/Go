package spider_my

import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/logs"                   //信息输出
	//	"net/http"
	"strconv"
)

func init() {
	Ctrip.Register()
}

var Ctrip = &Spider{Name: "携程酒店数据",
	Description:  "携程酒店数据 [http://hotels.ctrip.com/international/landmarks/]",
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			logs.Log.Critical("开始请求")
			//			ctx.Aid(map[string]interface{}{"Rule": "生成请求"}, "生成请求")
			ctx.AddQueue(&context.Request{
				Url:  "http://hotels.ctrip.com/international/landmarks/",
				Rule: "生成请求",
				//				Header: http.Header{"Content-Type": []string{"text/html", "charset=GBK"}},
			})

			logs.Log.Critical("请求加入队列")
		},
		Trunk: map[string]*Rule{
			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					logs.Log.Critical("开始 第一个请求")
					ctx.AddQueue(&context.Request{
						Url:  "http://hotels.ctrip.com/international/landmarks/",
						Rule: aid["Rule"].(string),
						//						Header: http.Header{"Content-Type": []string{"text/html", "charset=GBK"}},
					})
					logs.Log.Critical("第一个请求已添加到队列")
					return nil
				},
				ParseFunc: func(ctx *Context) {
					//					nationUrlArr := make([]string, 0)
					logs.Log.Critical("开始解析" + ctx.GetUrl())
					query := ctx.GetDom()
					logs.Log.Critical(strconv.Itoa(query.Find("#divA").Length()))
					query.Find(".filter_detail .nation_list .nation").Each(func(i int, a *goquery.Selection) {
						nationUrl, exists := a.Find("a").Attr("href")
						//						logs.Log.Critical(strconv.FormatBool(exists), nationUrl)
						if exists {
							//							nationUrlArr = append(nationUrlArr, nationUrl)  "http://hotels.ctrip.com" +
							ctx.AddQueue(&context.Request{
								Url:      nationUrl,
								Rule:     "重点城市列表",
								Priority: 1,
							})
						}
					})
					logs.Log.Critical("解析完成" + ctx.GetUrl())
				},
			}, //生成请求
			"重点城市列表": {
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析 重点城市列表" + ctx.GetUrl())

					query := ctx.GetDom()
					logs.Log.Critical(query.Find(".more_city").Text())
					moreUrl, exists := query.Find(".more_city").Attr("href")
					if exists {
						logs.Log.Critical(moreUrl)
						ctx.AddQueue(&context.Request{
							Url:      "http://hotels.ctrip.com" + moreUrl,
							Rule:     "城市列表",
							Priority: 2,
						})
					} else {
						logs.Log.Critical(ctx.GetUrl() + "找不到URL")
					}
				},
			}, //重点城市列表
			"城市列表": {
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析 城市列表" + ctx.GetUrl())
					query := ctx.GetDom()
					logs.Log.Critical(strconv.Itoa(query.Find(".domestic ul li").Length()))
					//					logs.Log.Critical(query.Text())
					query.Find(".domestic ul li a").Each(func(i int, a *goquery.Selection) {
						cityUrl, exists := a.Attr("href")
						if exists {
							ctx.AddQueue(&context.Request{
								Url: "http://hotels.ctrip.com" + cityUrl,
								//								Header: http.Header{"Content-Type": []string{"text/html", "charset=GBK"}},
								Rule:     "酒店列表",
								Priority: 3,
							})
						}
					})
				},
			}, //城市列表
			"酒店列表": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					url := aid["Url"].(string)

					for loop := aid["loop"].([2]int); loop[0] <= loop[1]; loop[0]++ {
						url = url + "/p" + strconv.Itoa(loop[0])
						logs.Log.Critical("开始请求 酒店列表翻页" + url)
						ctx.AddQueue(&context.Request{
							Url:      url,
							Rule:     aid["Rule"].(string),
							Priority: 3,
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析 酒店列表" + ctx.GetUrl())
					query := ctx.GetDom()
					query.Find(".hotel_list .hotel_list_item .hotel_name").Each(func(i int, a *goquery.Selection) {
						hotelUrl, exists := a.Attr("d-url")
						if exists {
							ctx.AddQueue(&context.Request{
								Url: "http://hotels.ctrip.com" + hotelUrl,
								//								Header: http.Header{"Content-Type": []string{"text/html", "charset=GBK"}},
								Rule:     "酒店详情",
								Priority: 4,
							})
						}
					})
					currentPage := query.Find("#page_info .layoutfix .current").Text()
					if currentPage == "1" {
						totalStr := query.Find("#page_info .layoutfix a").Last().Text()
						total, err := strconv.Atoi(totalStr)
						logs.Log.Critical(ctx.GetUrl() + " 酒店列表 总页数" + strconv.Itoa(total))
						if err == nil && total > 1 {
							// 调用指定规则下辅助函数
							ctx.Aid(map[string]interface{}{"loop": [2]int{2, total}, "Rule": "酒店列表", "Url": ctx.GetUrl()})
						}
					}

				},
			}, //酒店列表
			"酒店详情": {
				ItemFields: []string{
					"酒店名称",
				},
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析" + ctx.GetUrl())
					query := ctx.GetDom()
					hotelName := query.Find(".htl_info .name").Text()
					logs.Log.Critical("酒店详情:" + hotelName)
					ctx.Output(map[int]interface{}{
						0: hotelName,
					})
				},
			}, //酒店详情
		}, //Trunk
	}, //RuleTree
}
