package spider_my

// 基础包
import (
	"github.com/PuerkitoBio/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/context" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	//	. "github.com/henrylee2cn/pholcus/app/spider/common"    //选用
	"github.com/henrylee2cn/pholcus/logs" //信息输出

	// net包
	//	"net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	//	"encoding/json"

	// 字符串处理包
	//	"regexp"
	"strconv"
	"strings"

	// 其他包
	// "fmt"
	// "math"
	// "time"
)

func init() {
	Booking.Register()
}

//var cookies_Taobao = "mt=ci%3D-1_0; swfstore=35673; thw=cn; cna=fcr5DRDmwnQCAT2QxZSu3Db6; sloc=%E8%BE%BD%E5%AE%81; _tb_token_=XLlMHhT9BI8IzeA; ck1=; v=0; uc3=nk2=symxAo6NBazVq7cY2z0%3D&id2=UU23CgHxOwgwgA%3D%3D&vt3=F8dAT%2BCFEEyTLicOBEc%3D&lg2=U%2BGCWk%2F75gdr5Q%3D%3D; existShop=MTQzNDM1NDcyNg%3D%3D; lgc=%5Cu5C0F%5Cu7C73%5Cu7C92%5Cu559C%5Cu6B22%5Cu5927%5Cu6D77; tracknick=%5Cu5C0F%5Cu7C73%5Cu7C92%5Cu559C%5Cu6B22%5Cu5927%5Cu6D77; sg=%E6%B5%B721; cookie2=1433b814776e3b3c61f4ba3b8631a81a; cookie1=Bqbn0lh%2FkPm9D0NtnTdFiqggRYia%2FBrNeQpwLWlbyJk%3D; unb=2559173312; t=1a9b12bb535040723808836b32e53507; _cc_=WqG3DMC9EA%3D%3D; tg=5; _l_g_=Ug%3D%3D; _nk_=%5Cu5C0F%5Cu7C73%5Cu7C92%5Cu559C%5Cu6B22%5Cu5927%5Cu6D77; cookie17=UU23CgHxOwgwgA%3D%3D; mt=ci=0_1; x=e%3D1%26p%3D*%26s%3D0%26c%3D0%26f%3D0%26g%3D0%26t%3D0%26__ll%3D-1%26_ato%3D0; whl=-1%260%260%260; uc1=lltime=1434353890&cookie14=UoW0FrfFYp27FQ%3D%3D&existShop=false&cookie16=V32FPkk%2FxXMk5UvIbNtImtMfJQ%3D%3D&cookie21=U%2BGCWk%2F7p4mBoUyTltGF&tag=7&cookie15=Vq8l%2BKCLz3%2F65A%3D%3D&pas=0; isg=C08C1D752BC08A3DCDF1FE6611FA3EE1; l=Ajk53TTUeK0ZKkG8yx7w7svcyasSxC34"

type Room struct {
	Id          string
	Title       string
	Description string
	Persons     string
	Price       string
}

type Advantage struct {
	AdvantageType string
	Title         string
	Description   string
}

type Facilitie struct {
	Title  string
	Detail []string
}

type Comment struct {
	UserName   string
	Nation     string
	Age        string
	Point      string
	Title      string
	Tag        []string
	Neg        string
	Pos        string
	SubmitTime string
}

var Booking = &Spider{
	Name:        "Booking酒店数据",
	Description: "Booking酒店数据 [Auto Page] [http://www.booking.com/destination.zh-cn.html/]",
	// Pausetime: [2]uint{uint(3000), uint(1000)},
	// Keyword:   USE,
	EnableCookie: false,
	RuleTree: &RuleTree{
		Root: func(ctx *Context) {
			logs.Log.Critical("开始请求")
			ctx.AddQueue(&context.Request{
				Url:          "http://www.booking.com/destination.zh-cn.html",
				Rule:         "生成请求",
				EnableCookie: true,
			})
			logs.Log.Critical("请求加入队列")
		},

		Trunk: map[string]*Rule{
			"生成请求": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
					//						url := []string{}
					//						//						for _, loc := range loc_Taobao {
					//						urls = append(urls, "http://www.booking.com"+aid["urlBase"].(string))
					//						}

					ctx.AddQueue(&context.Request{
						Url:          "http://www.booking.com" + aid["urlBase"].(string),
						Rule:         aid["Rule"].(string),
						EnableCookie: true,
						Priority:     1,
						//							DownloaderID: 1,
					})
					//					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析" + ctx.GetUrl())
					query := ctx.GetDom()
					logs.Log.Critical(strconv.Itoa(query.Find(".flatList a").Length()))
					query.Find(".flatList a").Each(func(i int, a *goquery.Selection) {

						if strings.Trim(a.Text(), " ") != "中国" { //排除中国酒店
							hrefNation, _ := a.Attr("href")

							ctx.Aid(map[string]interface{}{
								"urlBase": hrefNation,
								"Rule":    "城市列表",
							})
						}

					})
					logs.Log.Critical("解析完成" + ctx.GetUrl())
				},
			},

			"城市列表": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
					//						urls := []string{}
					//						//						for _, loc := range loc_Taobao {
					//						urls = append(urls, "http://www.booking.com"+aid["urlBase"].(string))
					//						}

					ctx.AddQueue(&context.Request{
						Url:          "http://www.booking.com" + aid["urlBase"].(string),
						Rule:         aid["Rule"].(string),
						EnableCookie: true,
						Priority:     2,
						//							DownloaderID: 1,
					})
					//					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析 城市列表" + ctx.GetUrl())
					query := ctx.GetDom()
					query.Find(".general tbody tr td a").Each(func(i int, a *goquery.Selection) {

						hrefCity, _ := a.Attr("href")

						ctx.Aid(map[string]interface{}{
							"urlBase": hrefCity,
							"Rule":    "酒店列表",
						})
					})

				},
			},
			"酒店列表": {
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					//					urls := []string{}
					//					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
					//						urls = append(urls, "http://www.booking.com"+aid["urlBase"].(string))
					//					}
					logs.Log.Critical("开始请求 酒店详情" + aid["urlBase"].(string))
					ctx.AddQueue(&context.Request{
						Url:          "http://www.booking.com" + aid["urlBase"].(string),
						Rule:         aid["Rule"].(string),
						EnableCookie: true,
						Priority:     3,
					})
					return nil
				},
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析 酒店列表" + ctx.GetUrl())
					query := ctx.GetDom()
					//					logs.Log.Critical("酒店列表 数：" + strconv.Itoa(query.Find("[name='hotels']").Next().Next().Find("tbode tr td a").Length()))
					//					logs.Log.Critical(query.Find("[name='hotels']").Next().Next().Html())
					//					logs.Log.Critical(query.Find("[name='hotels']").Next().Next().Find("tr td a").Text())

					query.Find("[name='hotels']").Next().Next().Find("tr td a").Each(func(i int, a *goquery.Selection) {

						hrefCity, _ := a.Attr("href")
						logs.Log.Critical(hrefCity)
						ctx.Aid(map[string]interface{}{

							"urlBase": hrefCity,
							"Rule":    "酒店详情",
						})
					})

				},
			},
			"酒店详情": {
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析 酒店详情" + ctx.GetUrl())
					query := ctx.GetDom()
					body := query.Find(".right")
					hotelTitle := body.Find("#wrap-hotelpage-top")

					bolckDisplay := body.Find("#blockdisplay1")

					photos := ""
					hotelTitle.Find("#photos_distinct a").Each(func(i int, a *goquery.Selection) {
						href, _ := a.Attr("href")
						photos += strings.Replace(href, "max400", "840x460", 1) + ";"
					})

					summery := ""
					bolckDisplay.Find("#summary p").Each(func(i int, a *goquery.Selection) {
						summery += a.Text() + "\r\n"
					})

					//description
					rooms := make([]Room, 0)

					bolckDisplay.Find("#maxotel_rooms tbody tr").Each(func(i int, a *goquery.Selection) {
						room := new(Room)

						id, _ := a.Find(".ftd div").Attr("id")
						room.Id = id
						room.Title = a.Find(".ftd .room-info .show-details").Text()
						room.Description = a.Find(".ftd .room-info .room-description").Text()
						person, _ := a.Find(".occ_no_dates .jq_tooltip").Attr("title")
						room.Persons = person
						room.Price = a.Find(".lowest-price").Text()

						rooms = append(rooms, *room)
					})

					traffic := make([]string, 0)
					bolckDisplay.Find(".pub_trans .pub_trans_item").Each(func(i int, a *goquery.Selection) {
						traffic = append(traffic, a.Text())
					})

					advantages := make([]Advantage, 0)
					bolckDisplay.Find("#uspsbox div").Each(func(i int, a *goquery.Selection) {
						id, isExist := a.Attr("id")
						if isExist {
							advantage := new(Advantage)
							advantage.AdvantageType = id
							advantage.Title = a.Find("p").First().Text()
							advantage.Description = a.Find("p").First().Next().Text()

							advantages = append(advantages, *advantage)

						}
					})

					//设施

					facilities := make([]Facilitie, 0)
					bolckDisplay.Find("#hp_facilities_box .facilitiesChecklistSection").Each(func(i int, a *goquery.Selection) {
						facilitie := new(Facilitie)
						facilitie.Title = a.Find("h5").Text()
						a.Find("ul li").Each(func(j int, b *goquery.Selection) {
							facilitie.Detail = append(facilitie.Detail, b.Text())
						})
						facilities = append(facilities, *facilitie)
					})
					cards := make([]string, 0)
					//定前必须	付款方式
					bolckDisplay.Find("#b_tt_holder_2 p").Each(func(i int, a *goquery.Selection) {
						card, isExist := a.Attr("class")
						if isExist {
							cards = append(cards, card)
						}
					})

					//评论URL
					//					comment, exists := bolckDisplay.Find(".seo_reviews_block .show_all_reviews_btn").Attr("href")

					ctx.Output(map[int]interface{}{
						0:  hotelTitle.Find("#b_tt_holder_1 span").Text(),                                         //tag
						1:  hotelTitle.Find("#hp_hotel_name").Text(),                                              //name
						2:  hotelTitle.Find(" .star_track .invisible_spoken").Text(),                              //star
						3:  hotelTitle.Find(" .hp__hotel_ratings .jq_tooltip .invisible_spoken").Text(),           //tooltip
						4:  hotelTitle.Find("#hp_address_subtitle").Text(),                                        //address_subtitle
						5:  photos,                                                                                //照片集合
						6:  summery,                                                                               //summery
						7:  bolckDisplay.Find(".hotel_meta_style").Text(),                                         //酒店上线时间以及客房数
						8:  rooms,                                                                                 //房间类型
						9:  traffic,                                                                               //交通
						10: advantages,                                                                            //优势
						11: bolckDisplay.Find("#HotelFacilities .hp_facilities_score").Text(),                     //设施评分
						12: facilities,                                                                            //设施
						13: bolckDisplay.Find("#hp_policies_box #hotelPoliciesInc #checkin_policy p").Text(),      //checkin 时间
						14: bolckDisplay.Find("#hp_policies_box #hotelPoliciesInc #checkout_policy p").Text(),     //checkout 时间 cancellation_policy
						15: bolckDisplay.Find("#hp_policies_box #hotelPoliciesInc #cancellation_policy p").Text(), //取消政策 children_policy
						16: bolckDisplay.Find("#hp_policies_box #hotelPoliciesInc #children_policy p").Text(),     // 儿童加床
						17: bolckDisplay.Find("#hp_policies_box #hotelPoliciesInc #description p").Text(),         //描述
						18: cards,                                                                                 //付款方式
						19: bolckDisplay.Find("#hp_important_info_box .description").Text(),                       //预定须知
						//						20: interface{}{},  //综合评价
						//						21: []interface{}{}, //评论列表                                                                   //评论
					}, "结果")
					logs.Log.Critical("解析完毕 酒店详情" + ctx.GetUrl())
					//					if exists {
					//						ctx.AddQueue(&context.Request{
					//							Url:          "http://http://www.booking.com" + comment,
					//							Rule:         "酒店评价",
					//							Temp:         temp,
					//							Priority:     4,
					//							EnableCookie: true,
					//							//						DownloaderID: 1,
					//						})
					//					}

				},
			},

			"酒店评价": {
				ParseFunc: func(ctx *Context) {
					logs.Log.Critical("开始解析 酒店评价" + ctx.GetUrl())
					query := ctx.GetDom()
					comments := make([]interface{}, 0)
					query.Find(".review_list li").Each(func(i int, a *goquery.Selection) {
						comment := new(Comment)
						comment.SubmitTime = a.Find(".review_item_date").Text()
						comment.UserName = a.Find(".review_item_reviewer h4").Text()
						comment.Nation = a.Find(".reviewer_country").Text()
						comment.Point = a.Find(".review_item_header_score_container").Text()
						comment.Title = a.Find(".review_item_header_content").Text()
						a.Find(".review_item_info_tags li").Each(func(j int, b *goquery.Selection) {
							tag := strings.Trim(b.Text(), "•")
							comment.Tag = append(comment.Tag, strings.Trim(tag, " "))
						})
						comment.Neg = a.Find(".review_neg").Text()
						comment.Pos = a.Find(".review_pos").Text()
						comments = append(comments, *comment)
					})

					//					targetComment := comments.([]interface{})
					//					targetComment := comments.([]interface{})

					//					temp := ctx.GetTemps()
					temp := ctx.CopyTemps()
					//					if len(temp) < 21 {
					temp.Set(ctx.GetItemField(20, "结果"), query.Find("#review_list_main_score").Text())
					//					temp[ctx.IndexOutFeild(20, "结果")] = query.Find("#review_list_main_score").Text()
					scoreBreakDown := make(map[string]string)
					query.Find("#review_list_score_breakdown li").Each(func(i int, a *goquery.Selection) {
						key := a.Find(".review_score_name").Text()
						value := a.Find(".review_score_value").Text()
						scoreBreakDown[key] = value
					})
					temp.Set(ctx.GetItemField(21, "结果"), scoreBreakDown)
					//					temp[ctx.IndexOutFeild(21, "结果")] = scoreBreakDown

					//					temp[ctx.IndexOutFeild(22, "结果")] = []interface{}{}
					//					}

					//					discussAll := ctx.GetTemp(ctx.IndexOutFeild(22, "结果")).([]interface{})

					//					discussAll:=[]interface{}

					//					discussAll = append(discussAll, comments...)
					temp.Set(ctx.GetItemField(22, "结果"), comments)
					//					ctx.SetTemp(ctx.IndexOutFeild(22, "结果"), discussAll)
					//					temp[ctx.IndexOutFeild(22, "结果")] = append(temp[ctx.IndexOutFeild(22, "结果")], targetComment...)
					prePage := query.Find("#review_previous_page_link")
					if prePage != nil {
						_, exists := prePage.Attr("href")
						if exists {
							return
						}
					}
					nextPage := query.Find("#review_next_page_link")
					if nextPage != nil && nextPage.Length() > 0 {
						href, _ := nextPage.First().Attr("href")
						ctx.AddQueue(&context.Request{
							Rule:         "酒店评价",
							Url:          "http://http://www.booking.com" + href,
							Temp:         temp,
							Priority:     4,
							EnableCookie: true,
							//							DownloaderID: 1,
						})
					}

				},
			},

			"结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"标签",      //title
					"名称",      //price
					"星级",      //currentPrice
					"tooltip", //tooltip
					"地址",      //vipPrice
					"照片",      //unitPrice
					"综述",      //unit
					"上线时间客房数", //isVirtual
					"房间类型",    //ship
					"交通",      //tradeNum
					"优势",      //formatedNum
					"设施评分",    //nick
					"设施",      //sellerId
					"入住时间",    //guarantee
					"退房时间",    //itemId
					"取消政策",    //isLimitPromotion
					"儿童加床",    //loc
					"描述",      //storeLink
					"付款方式",    //href
					"预订须知",    //commend
					//					"总评分",     //source
					//					"细项评分",    //ratesum
					//					"所有评价",    //goodRate

				},
				ParseFunc: func(ctx *Context) {
					// 结果存入Response中转
					ctx.Output(ctx.CopyTemps())
				},
			},
		},
	},
}
