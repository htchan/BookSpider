package uukanshu

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	vendor "github.com/htchan/BookSpider/internal/vendorservice"
	"github.com/stretchr/testify/assert"
)

func TestParser_ParseBook(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		want      *vendor.BookInfo
		wantError error
	}{
		{
			name: "happy flow with real data",
			body: string(testBookBytes),
			want: &vendor.BookInfo{
				Title:         "从零开始",
				Writer:        "雷云风暴",
				Type:          "网游竞技小说",
				UpdateDate:    strconv.Itoa(time.Now().Year()-7) + "-01-01",
				UpdateChapter: "第二十三卷 第6章 放飞希望（完结篇）",
			},
			wantError: nil,
		},
		{
			name: "happy flow with date in month format",
			body: `<data>
				<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
					<h1><a title="book name最新章节"></a></h1>
					<h2><a>author</a></h2>
					<div class="shijian">5月</div>
				</dd></dl></div>
				<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
				<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate:    strconv.Itoa(time.Now().AddDate(0, -5, 0).Year()) + "-" + fmt.Sprintf("%02d", time.Now().AddDate(0, -5, 0).Month()) + "-01",
				UpdateChapter: "chapter name",
			},
			wantError: nil,
		},
		{
			name: "happy flow with date in day format",
			body: `<data>
				<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
					<h1><a title="book name最新章节"></a></h1>
					<h2><a>author</a></h2>
					<div class="shijian">5日</div>
				</dd></dl></div>
				<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
				<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate: time.Now().AddDate(0, 0, -5).Format(time.DateOnly), UpdateChapter: "chapter name",
			},
			wantError: nil,
		},
		{
			name: "title not found",
			body: `<data>
				<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
					<h2><a>author</a></h2>
					<div class="shijian">d a	t
					e</div>
				</dd></dl></div>
				<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
				<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Writer: "author", Type: "type",
				UpdateDate: time.Now().Format(time.DateOnly), UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookTitleNotFound,
		},
		{
			name: "writer not found",
			body: `<data>
				<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
					<h1><a title="book name最新章节"></a></h1>
					<h2><a></a></h2>
					<div class="shijian">d a	t
					e</div>
				</dd></dl></div>
				<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
				<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Type: "type",
				UpdateDate: time.Now().Format(time.DateOnly), UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookWriterNotFound,
		},
		{
			name: "type not found",
			body: `<data>
			<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
				<h1><a title="book name最新章节"></a></h1>
				<h2><a>author</a></h2>
				<div class="shijian">d a	t
				e</div>
			</dd></dl></div>
			<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author",
				UpdateDate: time.Now().Format(time.DateOnly), UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookTypeNotFound,
		},
		{
			name: "date not found",
			body: `<data>
			<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
				<h1><a title="book name最新章节"></a></h1>
				<h2><a>author</a></h2>
			</dd></dl></div>
			<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
			<div class="zhangjie"><ul id="chapterList"><li><a>chapter name</a></li></ul></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate: time.Now().Format(time.DateOnly), UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookDateNotFound,
		},
		{
			name: "chapter not found",
			body: `<data>
			<div class="xiaoshuo_content"><dl class="jieshao"><dd class="jieshao_content">
				<h1><a title="book name最新章节"></a></h1>
				<h2><a>author</a></h2>
				<div class="shijian">d a	t
				e</div>
			</dd></dl></div>
			<div class="weizhi"><div class="path"><a></a><a>type</a></div></div>
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate: time.Now().Format(time.DateOnly),
			},
			wantError: vendor.ErrBookChapterNotFound,
		},
		{
			name:      "all fields not found",
			body:      "<data></data>",
			want:      &vendor.BookInfo{UpdateDate: time.Now().Format(time.DateOnly)},
			wantError: vendor.ErrFieldsNotFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got, err := p.ParseBook(test.body)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestParser_ParseChapterList(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		want      vendor.ChapterList
		wantError error
	}{
		{
			name: "happy flow with real data",
			body: string(testChapterListBytes),
			want: vendor.ChapterList{

				{URL: "https://www.uukanshu.com/b/1248/51419.html", Title: "第1章 血拼"}, {URL: "https://www.uukanshu.com/b/1248/51420.html", Title: "第2章 建号（上）"},
				{URL: "https://www.uukanshu.com/b/1248/51421.html", Title: "第3章 建号（下）"}, {URL: "https://www.uukanshu.com/b/1248/51422.html", Title: "第4章 职业认证"},
				{URL: "https://www.uukanshu.com/b/1248/51423.html", Title: "第5章 辅助职业"}, {URL: "https://www.uukanshu.com/b/1248/51425.html", Title: "第1章 黄金圣龙"},
				{URL: "https://www.uukanshu.com/b/1248/51426.html", Title: "第2章 迷失之城"}, {URL: "https://www.uukanshu.com/b/1248/51427.html", Title: "第3章 公司公告"},
				{URL: "https://www.uukanshu.com/b/1248/51428.html", Title: "第4章 NPC朋友"}, {URL: "https://www.uukanshu.com/b/1248/51429.html", Title: "第5章 购物"},
				{URL: "https://www.uukanshu.com/b/1248/51430.html", Title: "第6章 练级！练级！"}, {URL: "https://www.uukanshu.com/b/1248/51431.html", Title: "第7章 追赶"},
				{URL: "https://www.uukanshu.com/b/1248/51432.html", Title: "第8章 被追杀"}, {URL: "https://www.uukanshu.com/b/1248/51433.html", Title: "第9章 星辰之戒"},
				{URL: "https://www.uukanshu.com/b/1248/51434.html", Title: "第10章 幻影"}, {URL: "https://www.uukanshu.com/b/1248/51435.html", Title: "第11章 沼泽"},
				{URL: "https://www.uukanshu.com/b/1248/51436.html", Title: "第12章 穿越死亡线"}, {URL: "https://www.uukanshu.com/b/1248/51437.html", Title: "第13章 阿伟之死"},
				{URL: "https://www.uukanshu.com/b/1248/51438.html", Title: "第14章 奇怪的任务"}, {URL: "https://www.uukanshu.com/b/1248/51439.html", Title: "第15章 打劫来的神器"},
				{URL: "https://www.uukanshu.com/b/1248/51440.html", Title: "第16章 卖煤炭也能发财"}, {URL: "https://www.uukanshu.com/b/1248/51441.html", Title: "第17章 披风任务"},
				{URL: "https://www.uukanshu.com/b/1248/51442.html", Title: "第18章 隐形披风"}, {URL: "https://www.uukanshu.com/b/1248/51846.html", Title: "第1章 入侵从这里开始"},
				{URL: "https://www.uukanshu.com/b/1248/51847.html", Title: "第2章 虚惊1场"}, {URL: "https://www.uukanshu.com/b/1248/51848.html", Title: "第3章 墨玉的线索"},
				{URL: "https://www.uukanshu.com/b/1248/51849.html", Title: "第4章 坑王之王"}, {URL: "https://www.uukanshu.com/b/1248/51850.html", Title: "第5章 迷之BOSS"},
				{URL: "https://www.uukanshu.com/b/1248/51851.html", Title: "第6章 2败俱伤"}, {URL: "https://www.uukanshu.com/b/1248/51852.html", Title: "第7章 自然灾害"},
				{URL: "https://www.uukanshu.com/b/1248/51853.html", Title: "第8章 伪装"}, {URL: "https://www.uukanshu.com/b/1248/51854.html", Title: "第9章 非常规作战"},
				{URL: "https://www.uukanshu.com/b/1248/51855.html", Title: "第10章 围点打援"}, {URL: "https://www.uukanshu.com/b/1248/51856.html", Title: "第11章 发洪水啦！"},
				{URL: "https://www.uukanshu.com/b/1248/51857.html", Title: "第12章 速冻10小时"}, {URL: "https://www.uukanshu.com/b/1248/51858.html", Title: "第13章 全面镇压"},
				{URL: "https://www.uukanshu.com/b/1248/51859.html", Title: "第14章 垂直800米"}, {URL: "https://www.uukanshu.com/b/1248/51860.html", Title: "第15章 吉祥如意"},
				{URL: "https://www.uukanshu.com/b/1248/51861.html", Title: "第16章 可爱宝贝"}, {URL: "https://www.uukanshu.com/b/1248/51862.html", Title: "第17章 初闻大联盟"},
				{URL: "https://www.uukanshu.com/b/1248/51863.html", Title: "第18章 奇怪的任务"}, {URL: "https://www.uukanshu.com/b/1248/51864.html", Title: "第19章 暗门"},
				{URL: "https://www.uukanshu.com/b/1248/51865.html", Title: "第20章 机兵洞穴"}, {URL: "https://www.uukanshu.com/b/1248/51866.html", Title: "第21章 魔偶师"},
				{URL: "https://www.uukanshu.com/b/1248/51867.html", Title: "第22章 技术实力"}, {URL: "https://www.uukanshu.com/b/1248/51868.html", Title: "第23章 艰难大联盟"},
				{URL: "https://www.uukanshu.com/b/1248/51869.html", Title: "第24章 整人项目"}, {URL: "https://www.uukanshu.com/b/1248/51870.html", Title: "第25章 过路费"},
				{URL: "https://www.uukanshu.com/b/1248/51871.html", Title: "第26章 铁索桥"}, {URL: "https://www.uukanshu.com/b/1248/51872.html", Title: "第27章 赌博"},
				{URL: "https://www.uukanshu.com/b/1248/51873.html", Title: "第28章 迷宫追逐战"}, {URL: "https://www.uukanshu.com/b/1248/51874.html", Title: "第29章 淘汰赛"},
				{URL: "https://www.uukanshu.com/b/1248/51875.html", Title: "第30章 洗牌"}, {URL: "https://www.uukanshu.com/b/1248/51876.html", Title: "第31章 混乱"},
				{URL: "https://www.uukanshu.com/b/1248/51877.html", Title: "第32章 合作"}, {URL: "https://www.uukanshu.com/b/1248/51878.html", Title: "第33章 反差"},
				{URL: "https://www.uukanshu.com/b/1248/51879.html", Title: "第34章 扩张准备"}, {URL: "https://www.uukanshu.com/b/1248/51880.html", Title: "第35章 圈地行动"},
				{URL: "https://www.uukanshu.com/b/1248/51881.html", Title: "第36章 诅咒信"}, {URL: "https://www.uukanshu.com/b/1248/51882.html", Title: "第37章 老熟人"},
				{URL: "https://www.uukanshu.com/b/1248/51883.html", Title: "第38章 正规的业余部队"}, {URL: "https://www.uukanshu.com/b/1248/51884.html", Title: "第39章 暗流"},
				{URL: "https://www.uukanshu.com/b/1248/51885.html", Title: "第40章 此消彼长"}, {URL: "https://www.uukanshu.com/b/1248/51886.html", Title: "第41章 就绪"},
				{URL: "https://www.uukanshu.com/b/1248/51887.html", Title: "第42章 恐吓"}, {URL: "https://www.uukanshu.com/b/1248/51888.html", Title: "第43章 钢铁之躯"},
				{URL: "https://www.uukanshu.com/b/1248/51889.html", Title: "第44章 耀日陷落"}, {URL: "https://www.uukanshu.com/b/1248/51890.html", Title: "第45章 蚂蚁吃大象"},
				{URL: "https://www.uukanshu.com/b/1248/51891.html", Title: "第46章 霸主"}, {URL: "https://www.uukanshu.com/b/1248/51892.html", Title: "第47章 暴利"},
				{URL: "https://www.uukanshu.com/b/1248/51893.html", Title: "第48章 情报"}, {URL: "https://www.uukanshu.com/b/1248/51894.html", Title: "第49章 风暴战役之进攻与防守"},
				{URL: "https://www.uukanshu.com/b/1248/51895.html", Title: "第50章 风暴战役之阻截"}, {URL: "https://www.uukanshu.com/b/1248/51896.html", Title: "第51章 风暴战役之侵入"},
				{URL: "https://www.uukanshu.com/b/1248/221461.html", Title: "第2章 召集令"}, {URL: "https://www.uukanshu.com/b/1248/221462.html", Title: "第3章(二十三卷)终极福利(完结倒数四)"},
				{URL: "https://www.uukanshu.com/b/1248/221463.html", Title: "第24卷 第1章 计划通过"}, {URL: "https://www.uukanshu.com/b/1248/221464.html", Title: "第4章(二十三卷)大集合(倒数三)"},
				{URL: "https://www.uukanshu.com/b/1248/221465.html", Title: "第5章(二十三卷)打包装船"}, {URL: "https://www.uukanshu.com/b/1248/221467.html", Title: "第二十三卷 第6章 放飞希望（完结篇）"},
			},
			wantError: nil,
		},
		{
			name: "happy flow",
			body: `<data>
				<div class="zhangjie"><ul id="chapterList">
					<li><a href="chapter url 4">chapter name 4</a></li>
					<li><a href="chapter url 3">chapter name 3</a></li>
					<li><a href="chapter url 2">chapter name 2</a></li>
					<li><a href="chapter url 1">chapter name 1</a></li>
				</ul></div>
			</data>`,
			want: vendor.ChapterList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "chapter url 2", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: "chapter name 3"},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantError: nil,
		},
		{
			name: "2nd chapter missing href",
			body: `<data>
				<div class="zhangjie"><ul id="chapterList">
					<li><a href="chapter url 4">chapter name 4</a></li>
					<li><a href="chapter url 3">chapter name 3</a></li>
					<li><a href="">chapter name 2</a></li>
					<li><a href="chapter url 1">chapter name 1</a></li>
				</ul></div>
			</data>`,
			want: vendor.ChapterList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: "chapter name 3"},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantError: vendor.ErrChapterListUrlNotFound,
		},
		{
			name: "3nd chapter missing title",
			body: `<data>
				<div class="zhangjie"><ul id="chapterList">
					<li><a href="chapter url 4">chapter name 4</a></li>
					<li><a href="chapter url 3"></a></li>
					<li><a href="chapter url 2">chapter name 2</a></li>
					<li><a href="chapter url 1">chapter name 1</a></li>
				</ul></div>
			</data>`,
			want: vendor.ChapterList{
				{URL: "chapter url 1", Title: "chapter name 1"},
				{URL: "chapter url 2", Title: "chapter name 2"},
				{URL: "chapter url 3", Title: ""},
				{URL: "chapter url 4", Title: "chapter name 4"},
			},
			wantError: vendor.ErrChapterListTitleNotFound,
		},
		{
			name:      "no chapters found",
			body:      `<data></data>`,
			want:      nil,
			wantError: vendor.ErrChapterListEmpty,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got, err := p.ParseChapterList("", test.body)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestParser_ParseChapter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		body      string
		want      *vendor.ChapterInfo
		wantError error
	}{
		{
			name: "happy flow",
			body: string(testChapterBytes),
			want: &vendor.ChapterInfo{
				Title: "第二十三卷 第6章 放飞希望（完结篇）",
				Body:  "(adsbygoogle = window.adsbygoogle || []).push({});\n\n            \n            免费大结局＆完本感言\n\n————————————————————————————————————————————————\n\n全部人员装船完成之后，没有鲜花也没有掌声。△，撤离注定需要在默默无闻之中完成。中央那边倒是来了几位大员，但都是悄悄来的，不敢声张。发射的日期特意选了一个风雨交加的夜晚，目的就是要用雷雨掩盖发射时的各种声光现象。\n\n移民船并不是那种普通的小飞机，所以雷电和风雨并不影响它的起飞。\n\n来送行的大员之中有几个人看着飞船点火，在一阵山摇地动之中升空，忍不住眼泪就下来了。直到目送飞船消失在雨云之中后，我才转向他们说道：“各位别太伤心了，又不是见不着了。虽然现实中不能再见，但在游戏世界里大家还是可以见到家人的吗。”\n\n这些来送行的人之所以这么伤心，就是因为他们的妻女都在船上。虽然理智上知道送走他们可以让他们获得绝对的安全，但情感上，一想到家人都被送走了，就自己还在地球上，那种巨大的情感冲击一般人还真扛不住。好在送走的人不是马上开始星际移民，他们还要在太空中等待后续船队陆续发射，所以游戏网络依然还能用，大家还是可以见面。\n\n虽然这些人被提前发射了，但是他们真正开始星际航行将需要等到所有移民船都升空之后才能开始，毕竟我们龙族没有那么多护航舰船，所以船队肯定不能飞开走，大家必须要聚集成团，这样遇到危险的话就会发挥群落效应，至少能保证一部分人脱离危险。\n\n随着第一批飞船升空。我不但没有闲下来，反而是更加的忙了。\n\n政府第一批移民船发射之后，世界各国的移民船也开始陆续升空。这种事情瞒得住老百姓，瞒不住别的国家，所以大家都在发射，何况我们本来也没打算偷跑。\n\n但是。这发射次数多了，暴露也就成了必然，只是事情比我们想象的要早了很多。\n\n原本我们的想法是，大概六个月之后才会有人开始大规模的意识到政府正在撤离民众，但实际上，在我把冰霜玫瑰盟那帮人送走之后两个月，就已经开始出现了这种现象。\n\n最先出问题的是印度。他们的一艘飞船在发射的时候发生故障，居然刚升空十几秒就失去动力。装有3万多人的移民船从五六百米的高空猛然砸落地面，装有休眠人员的休眠舱四散在几平方公里的地面上。其中居然还有少数人没有死。\n\n\n(adsbygoogle = window.adsbygoogle || []).push({});\n\n毕竟飞船很大，先落地的地方起到了缓冲作用，后落地的位置人员冲击较小，加上休眠舱本身有防撞功能，因此还是有不少人幸存。\n\n也亏了印度的技术不过关，造不出太大的移民船，所以这艘船只有3万人的运载量。这要是我们龙缘的飞船，掉一艘下去就十几万人完蛋了。那绝对是超级惨剧。\n\n不过，虽然实际上只死了两万多人。还有七八千人活了下来，但是这个事情却是彻底曝光了。\n\n本来，这个事情要是发生在别的国家，封锁消息的话，至少还能再拖上个把月，但印度的政府执行能力真心让人蛋、疼。尽管世界各国联合起来动用各自的新闻控制能力管制新闻扩散。但奈何印度这边几乎是完全没有任何作为，不到十天全印度都知道这个事情了，然后因为有十多亿印度人知道了这个事情，别国再怎么封锁还是有消息泄露了出来。最后就变成了大家都知道了这个事情。\n\n各国政府一开始还想封锁来着，但是眼看彻底封不住了。于是就改变策略，直接开始和稀泥，各种乱七八糟的信息把新闻搅成一锅粥，反正就是让大家都搞不清楚具体情况。不过，各种暗流还是开始涌动了起来。\n\n最先出问题的是韩国。国内发生游行示威，民众要求知道真相，最后和军警发生大规模冲突，死了好几十人才被控制住。\n\n韩国高层询问移民计划联合指挥中心，是否可以公布事实真相，但是各国现在也是没个主意，最后只能让他们先拖着。\n\n又过了1个多月，事情终于再次爆发，这次又是印度，而且情况奇葩到我们这些别国的管理着听到之后都感觉有些蒙圈。\n\n印度那边竟然有一个移民船基地被一个家族史的部族武装给占领了。这种天方夜谭一样的消息差点没把我们震晕过去。更奇葩的是，那个家族中学历最高的人就是一个海归留美博士。这位按说也是文化人，但他去美国学的是计算机编程，星际移民船这么高大尚的玩意他当然是玩不转的。但是，他至少搞清楚了这东西是个移民船，而且从飞船上的电脑中弄到了整个移民计划。\n\n这个计划不是印度一家的计划，因为我们现在是全世界联合行动，移民计划都是一起制定的，各自分头完成。印度的计划其实就是大家的计划，所以现在就变成我们大家的计划全都曝光了。\n\n这个家族的人搞清楚了这基地里的几艘飞船用途之后就开始招兵买马，以许诺登船资格为条件，扩充武装力量，而且向全世界范围招收可以操纵这个移民船的人来帮他们开船。\n\n本来这个事是印度人自己搞出来的，别的国家都是无辜躺枪，被猪队友坑了。但是，还没等我们找印度抱怨，对方却是先找上我们了。\n\n基于之前的一连串奇葩事件，印度这边找我们的原因果然依旧奇葩。\n\n事情还是那个部族武装那里的问题。对方在世界上招人开飞船，然后真的来了好多飞行员和船舶方面的人才，因为大家都不知道这个东西到底更像一艘船还是一架飞机，反正两边沾边的人都来了不少。\n\n听说有可以登船带家眷一起走的机会，大家当然都很积极。但是，用了半个多月。这些人来到印度，然后研究了十几天，最终发现——这玩意他娘的没法开。\n\n飞船的驾驶舱里面根本没有任何的控制按钮或者手柄之类的东西，这里只有三座带有透明椭圆形舱盖的休眠舱。那个留美博士分析研究之后确认了这三个休眠舱就是驾驶系统。\n\n他们一开始还以为这东西是有某种验证，所以除了指定人员，别人开不了。但是，在陆续赶到的那些技术人才的帮助下，他们用了半个月，终于明白了一件事情。这个事情就是，这艘飞船的控制方式居然是神经系统直连控制。而且，整艘飞船庞大的管理网络，居然只有三个人控制。这也就意味着，每个驾驶员都要承担极为可怕的数据洪流。\n\n为了活命，这帮人也冒险测试过。强行启动的结果就是。躺在控制仓里面的人会在启动瞬间就开始七窍流血，十秒之内不把他弄出来就会彻底脑死亡。\n\n分析完这些情况后他们得出结论，要驾驶这些飞船就需要找到大脑超级发达的人，可以承受这种庞大的数据冲击。但是，他们不管怎么找都找不到这样的人，因为这个东西需要的运算量超过人类极限太多太多，就算有爱因斯坦那样的脑子也绝对扛不住。\n\n其实发生这种情况的原因很简单，因为移民计划最初的时候各国就已经决定过了。我们龙族将是全人类的领航员。也就是说，那个控制仓其实是给我们龙族准备的。\n\n这东西对我们龙族来说。计算负荷并不算很大，一个人就可以控制一整艘船，但是会比较疲劳。用三个控制仓主要是为了分担一下压力，这样大家都比较轻松。但是，轻松是相对我们来说的，人类的大脑是绝对撑不住的。所以，敢于躺进这个控制仓，而且还启动它的人，唯一的下场就是脑死亡。\n\n这些人发现自己发射无望之后就沉寂了几天，然后突然开始威胁印度政府。因为他们想到了一个情况。\n\n既然印度政府建造这样的飞船，那么印度政府肯定是有准备飞行员的，他们不可能弄出自己开不了的飞船来。所以，这些人最后开始威胁政府，内容就是，基地里现在又8艘飞船，他们只要一半，让印度政府提供12名驾驶员，然后剩下的飞船和基地一起还给政府。如果印度政府不同意，他们就炸掉飞船和基地然后全体家族武装一起攻击其他的基地，大家玉石俱焚。\n\n放在别的国家，这种威胁肯定不会答应。小国的话估计会扯皮，然后打上几次再看情况，大国的话估计直接二话不说就武装突击了，而且一般国家根本就不会发生之前那种被人占领发射基地的事情。但是，印度政府不但答应了对方的要求，而且居然还厚着脸皮跑来找我们要龙族驾驶员。\n\n这种白痴要求我们会答应才见鬼呢。我们的驾驶员过去算什么情况？俘虏？人质？我们可丢不起那人。\n\n最终我们当然是拒绝了要求，结果印度政府又干了一件奇葩的事情。他们不敢动自己国家的一个小小的部族武装，居然威胁我们，如果我们不答应，他们因此遭受损失了，他们就会向我们发射核弹，让我们的飞船也陪葬。\n\n这种奇葩的思维方式让联合国几个大流氓们都感觉自己脑袋有点不够用，话说对方这是怎么想出来这种“解决方式”的呢？\n\n无奈也没用，反正人家就这么说了，我们还的想办法。虽说我们有拦截系统，核武对我们不一定有用，但这种时候，世界已经够乱的了，大家都不想再添麻烦。所以，最终印度人的麻烦变成了我们的麻烦。之前解决美国危机的时候出动的联合国武装小组不得不再次聚集，然后跑去印度打部族武装。\n\n最后的结果当然是压倒性的。当初美国人那边的叛徒是准备已久，而且自身就有超强的科技力量，还能调动美国的军队，自然难搞。但是印度这边的部族武装最强的武器就是不知道什么年代的步兵战车，貌似还只能当固定炮台用。我们真的很想知道他们当初怎么用这些东西打赢政府军的。\n\n反正最后基本上就是一群未来战士虐杀土著的故事，我们和德国方面的突击队用了两个小时清理光了基地内的部族武装人员。法国和日本的特遣队扫荡了外围的村落，虽然比我们多用了俩小时。但时间其实都耽搁在路上了，交战时间一共没用到十分钟。\n\n美国方面这次也派出了自己的特战队，但是没有来参加扫荡，而是被拆分之后安排到了印度的各个移民船基地之中。他们将驻扎在那里，直到飞船全部起飞，免得再出这种让人蛋|疼的事情。\n\n虽然印度这边的问题解决了。但是因为这一番折腾，各国的消息就彻底控制不住了，尤其是印度自己国内基本上是一片混乱。\n\n有了一个成功的案例，于是印度这边各地的人开始效仿那个被灭的部族，开始攻占附近的移民船基地。而直到此时大家才发现，印度这边的移民船基地的位置，居然是人尽皆知的，反正当地人都可以准确的找到这些基地。别的国家那些基地就算暴露，外面的人最多也就知道这个地方有个秘密基地。在建造什么东西，但具体内容都不知道。印度这边倒好，是个人就知道那是移民船发射基地。\n\n美国人在这些基地都丢下了战斗小组，但是战斗小组完全挡不住那些老百姓，因为人太多，比他们的子弹数量还要多，所以他们最后无奈的只能自己跑路。美国人可没有和阵地共存亡的习惯，再说这是印度人的阵地。他们撑死了能算是外援。\n\n不过，跑出来的美国人至少让我们知道了为什么印度的第一个基地会被抢了。原来是因为基地里面的守卫就是附近村子上的人。所以，基地其实不应该算是被占领，而是叛变了。\n\n反正这场大规模暴动的最后结果就是印度人用了一个月损失了三分之一的飞船基地，而且，因为这次是无组织的暴动，和之前那种有一定组织性的情况不一样。所以在占领了基地之后发生更奇葩的事情。\n\n当地人居然开始拆飞船上的东西，感觉什么东西有用就拆回家自己用。比如说把飞船上的冷却管拆回家修理家里的自来水管，还有把飞船上的电线拆回去自己用或者卖到废品回收点。反正印度人用让我们目瞪口呆的速度把那些基地里的飞船全都给拆了个七零八落。\n\n印度目前是世界第一人口大国，而他们的工业能力又极端的弱，虽然计算机软件方面发展不错。但计算机又飞不起来。移民船需要的是重工业基础。所以，原本印度就是移民率偏低的国家之一，最初公布的移民比例就比其他国家低几十个百分点，现在更惨了。剩下的飞船只能带走最多8亿人，起码一半的国民都运不走了。更要命的是，剩下的那三分之二的基地也未必就是安全的，鬼知道还会发生什么奇葩的事情。毕竟从开始逐批次发射飞船到现在，磕磕绊绊已经5个月了，就只有一艘飞船坠毁，就是印度的。剩下的那些飞船就算没有再被抢夺破坏，能全飞起来的可能性也不高。\n\n印度的混乱很快就开始蔓延，一些落后国家，尤其是那种自己飞船数量严重不足的国家，各种矛盾变的尖锐了起来。\n\n第一批飞船发射6个月后，非洲发生大规模暴乱。维和部队没有出动维持秩序，反而是撤回了各自国内。现在的情况是大家都只能先顾自己，根本没空管落后国家。\n\n非洲和南美的情况急剧恶化，然后接着就是大量的难民私自前往欧洲和北美，当然亚洲也有，但是数量不高。主要是因为非洲和南美这俩暴乱区域距离亚洲都比较远。亚洲也就是一些小国家逐渐开始出现了混乱。\n\n第一批飞船发射8个月后，意大利和法国都分别发生了大规模屠杀难民的情况。这个主要是大批难民涌入，欧洲国家先开始想遣返，但是人手不够。后来想，反正碎着人员不断的发射，欧洲本地人越来越少，大部分区域都空出来了，所以他们就在欧洲划分了一些安全区，让难民在里面待着。\n\n但是，难民跑到欧洲是想要登船。不是为了住在欧洲，所以他们开始想要离开这个区域，并最终发生了大规模冲突。\n\n难民之中有人还带了武器过来，欧洲国家的军警只能开枪，而一旦开火就刹不住了，最后变成了大屠杀。难民死伤惨重。\n\n对这个事情美国和我国也什么都没说，现在不同往日，互相攻击没任何好处，大家都只能装作看不见。\n\n欧盟这边也开始破罐子破摔，反正已经开始屠杀难民，干脆彻底放开。让军舰攻击一切靠近海岸线的非法移民船，并且在主要发射区附近建立隔离带，进入其中的难民会被无条件射杀。\n\n美国的情况和欧盟没啥区别，只不过难民多来自南美地区。\n\n第一批飞船发射10个月后。朝鲜突然对韩国动武，而且上来就是核武器开路。我国的杀手卫星虽然挂掉了，但是还有一艘太空战舰在呢。新星号成功拦截了核武器，但是朝鲜大规模入侵韩国，美国人没空管，韩国人自己挡不住，最后向我们求援。这个事情搞的我们也是非常的为难。\n\n政治上来说朝鲜和我们关系更近一点，但是现在这情况。我们也不可能带朝鲜一起玩，韩国人自己有移民船。虽然带不走全部人口，但至少不用我们操心，所以现在理论上我们还是希望抱住韩国，因为这样移民总人口会增加。而如果让朝鲜入侵下去，韩国人战败前肯定会紧急发射一部分飞船，剩下的估计会被摧毁掉。这肯定不是好事。\n\n最后大家商量了一下，决定支持韩国，但是没有直接动手对付朝鲜，而是先表态，然后劝说了一番。反正就是让韩国匀出10万人的在运载能力给朝鲜，然后我们中国这边也帮忙支援朝鲜20万人的运载量，希望朝鲜就此放弃攻击韩国的计划。\n\n朝鲜答应了这个计划，一个月后答应他们的30万人运载量发射，朝鲜有30万人升空，但是这30万人刚走，朝鲜又开始进攻了。\n\n这次真的是把我们搞郁闷了，更糟糕的是全世界都开始出问题了。\n\n欧洲地区的某个基地不知道是被恐怖袭击还是怎么回事，发生了爆炸。爆炸摧毁了基地的顶棚，然后掉下来的建筑材料砸坏了两艘8万人的移民船，而且修复需要3个多月，也可能永远修不好了，毕竟现在的社会形势大家也都知道，实在是不能用平常的标准去衡量。\n\n这一情况不知道怎么搞得就泄露了出去，然后本来家里有粮心中不慌的欧洲人也坐不住了，本来他们以为是100%移民，现在却发现可能有8万人走不掉。大家都不希望成为这8万人中的一个，所以各地开始示威游行，然后冲突升级，再然后发生了一个超级严重的情况。\n\n因为民众在示威游行的时候和军警冲突演变成暴动，最后一处政府大楼被攻陷，在混乱中，一辆出逃的车辆发生侧翻后被暴民抓住了里面的人，然后那些人就被活活打死了。\n\n但是，这些人不知道的是，他们打死的人，就是美国人支援给他们，去修复那两艘飞船的工程师。也就是说，现在那两艘飞船是真的飞不起来了。\n\n在欧洲混乱的同时，日本地区也开始出问题，不过不是人的问题，而是地震导致一艘飞船的发射架倾斜。因为当时正在发射的过程中发生地震，所以飞船被倾斜扭曲的发射架拖住，没有飞起来反而砸到了旁边的飞船。虽然没有爆炸，但是两艘船都完蛋了，并且还堵住了发射口，这个基地里剩下的三艘船也飞不起来了。\n\n这一重大事件果然导致日本内部出现混乱，也开始发生暴动，不过总算没搞出更大的问题。美国那边也差不多，暴乱一直就没停过，最近连续损失了三艘大型飞船，起码50万人的运力完蛋了。\n\n各国征服当年玩了命的掩盖事实，就是怕发生这种情况，结果因为印度的猪队友行为，现在还是发展到这个状态了。要说唯一能让各国首脑稍微欣慰一点的就是，他们至少已经送了不少人上去了。就算下面彻底完蛋，起码人类的种子已经安全了。\n\n第一批飞船发射第12个月。全世界都和世界大战差不多了，不同的是大家都在对付自己的邻居和自己国家内部的人民暴动。\n\n我国目前已经彻底放弃了西部地区，将人口撤离到东部地区集中了起来。移民的真相已经彻底解密，反正也瞒不住了。干脆放开，这样大家还能心平气和的谈一谈。\n\n虽然我国人民的服从性较好，但是架不住邻居不好啊！以前是干不过我们，不敢上，但是现在生死存亡，反正不上也是死。所以周围的南亚小国纷纷向我国境内入侵。\n\n对于这种情况其实一号首长早有预料，所以我们在发射的时候，都是优先发射国境线附近的基地的飞船，现在距离最后时刻已经不远，边境大部分地区的飞船基地其实都已经彻底空掉了，人员都撤到了中部地区，人口也基本都移动了过去。\n\n现在我国人口基本上就是集中在北京、上海、南京、重庆这四个城市周围地区，其他地方几乎都快成无人区了。\n\n这一年以来我们已经陆陆续续送走了9亿多人，剩下的人口集中起来之后也在等待最后的发射时间。不过。我们目前拥有的运载量只有1亿左右，而剩下的人显然远不止这个数。不过，稍微让我们安心的是，有大约两亿左右的运输能力可以在两个月内完成，也就是说两个月后我们就有3亿的运载量了。当然，我们不会等到全部造好在一起发射，而是会在这段时间内陆续发射。\n\n印度那边现在已经彻底无政府主义了，他们的飞船总共就发射了大概不到3亿人口。剩下的飞船全都因为各种原因无法发射，而且印度目前好像还有大约1亿多的可用运载量。但全都分布在印度各地的部族武装和暴民手里，反正都浪费了。\n\n对于那部分飞船，一号首长的意思是让我们看看能不能利用一下。现在印度政府已经彻底瘫痪，原因不是被暴民袭击，而是他们把政府成员全都发射升空了。一群政府官员全都升空了，难道指望下面的人自制吗？但是不管怎么说。他们已经走了，骂他们也没用，但是那些船说不定还能用。虽然印度的飞船不爱靠谱，但是检查一下的话，应该能挑出一些能用的。再说就算我们不要。拿来做人情分配给周边国家也能减少一些我们的压力。只要让那些周边国家看到一些希望，他们就不会再和我们玩命了。\n\n各国就这样在战火中度过了两个月。期间我们真的跑了一趟印度，还别说，真的抢回来大约5千万人的运载量，不过我们只保留了两千万运载量的飞船，剩下的被分给了朝鲜、泰国、缅甸、越南之类的国家。这些国家得到这些飞船之后果然停止了对我们的攻击，让我们也松了口气。\n\n其实5千万运载量我们打算全要的，那3千万之所以送人，实在是因为质量有问题，起飞成功率不到60%，全很利弊之后还是送人划算一点。\n\n第一批飞船发射15个月，距离最后安全线还有5个月时间。不过此时各主要国家都已经完成了主要撤离任务，各国已经开始陆续撤离政府工作人员和军队了，动作最快的就是日本，已经全员发射完成，虽然最后因为各种原因全国只有1。6亿人升空，但起码人家最先完成了。\n\n主要经济体方面。\n\n欧盟磕磕绊绊的也算是大体上完成了任务，但是最后丢下了差不多5千万人口，不是没有运力，而是实在乱的没法管了，不得已只能丢下这帮人让他们占领欧洲去吧。\n\n俄罗斯方面一直不声不响的自己玩自己的，对各国的支援力度最小，不过他们的混乱情况也最小，最后大部分国民安然发射，只有预定的一些老人或者罪犯之类的被丢下来了。\n\n印度方面继续无政府状态，老百姓天天和过年一样，也不知道高兴个啥，反正我们打算把地球让给他们了，其他的我们也管不了。作为抢了他们5千万运力的感谢。我们帮助印度地区一个比较稳定的地方武装完成了一艘飞船的发射，总算让印度地区又送走了3万多人。\n\n美国地区撤离基本完成，地面上剩下的不是军队就是警察，反正已经是最后收尾阶段了。\n\n韩国比日本速度慢点，但也发射了大部分飞船，不过因为朝鲜捣乱。他们被迫引爆了两个基地，损失80万人的运载量，而且战死不少人，不过总算没有伤筋动骨，大部分人都升空了。\n\n我国重庆地区最先完成发射任务，预定人员撤离之后当地的政府部门人员也大部分撤离。军队和少量政府人员移动到南京地区和我们龙缘的人汇合。\n\n上海地区紧跟着重庆地区完成发射，然后全体政府人员和军队一起撤离。北京地区是第三个，一号首长本来是打算留下来不走的，但是因为现在我们的运输力量膨胀的厉害，既然有空位置，他也没打算非要留下和地球共存亡。之前表态要留下不过是怕他先走了，下面就乱套了，想要留下坐镇，但现在人都撤光了。他留下就毫无意义了。\n\n北京地区的普通民众发射完成后，政府人员大部分撤离，少量人员和剩余的一部分动力较强的部队则是移动到南京。\n\n我们这边等北京的人员到了之后也开始撤离。民众都提前发射完成，集团人员也多数都已经离开，剩下的基本上不是龙族就是武装人员。\n\n我们这些人和北京以及重庆地区过来的军队汇合后并不是留在南京，而是向着第四特区移动。\n\n用了五天，我们全员移动到第四特区，一路山看着千里无人烟的祖国大地。不少人都哭了，不过大家心里只是感慨。不算多伤心。至少我们把绝大部分国民都送上轨道了，这个成绩已经是各国之中数一数二的了，何况我们人口那么多。\n\n打刀第四特区之后安排军队和最后的政府人员一起进入第四特区，然后看着这艘山岳一般的飞船起飞，而我和少数龙族精锐则是等到他们起飞之后又四散开来。\n\n我们的任务很简单，对全国各地进行一次遥感扫描。看看还有没有被遗漏的人员。至于我们的撤离飞船则是神龙号，他将是我国最后撤离的一艘船。\n\n巡视工作用了整整半个月，即便是高空扫描也累得够呛。不过成果也不是没有。半路上捡到了两个七八岁大的小孩，还有一对情侣，然后就是俩老头和二十多个中年人。都是因为各种原因意外被丢下的。这种全国性质的打撤离，只丢了这么几个人，实在是非常不错了。\n\n最后确认了一遍没有遗漏之后我们终于登上了神龙号。龙缘南京基地在一阵爆炸声中被掀飞了半个顶棚，然后神龙号伴随着飞溅的碎石摇动着飞上天空直插云霄。\n\n进入同步轨道之后的神龙号开始进入已经完成编队的大型移民舰队，然后带领着整个地球上最终发射上来的人向着深空推进。\n\n尽管舰队中最终只有45亿人口，和原本的地球人口比几乎只有一半多点，但至少人类的种子已经上路了。但是，和那些躺在休眠舱中的人不同的是，我们龙族还要带着这45亿人披荆斩棘，一路抵达人类的新家。\n\n朝圣者的信息告诉我们，前路并不平坦，但是为了生存，我们别无选择。但愿我们可以顺利到达那希望之所。\n\n躺在神龙号专用控制室中，通过心灵网络，我对全体龙族领航员宣布：“全舰队起锚，让我们放飞希望。”\n\n=========================================================================\n\n《从零开始》到这里就算是彻底完结了。UU看书 ｗww.ｕｕkａnｓhu.coｍ 写了11年真的对不起大家！有哪位没能坚持到我完本就挂了的，可以托梦让家里人去坟头烧书了，阴曹地府那边的版权没卖出去，只能各位活着的帮他们代购了。\n\n以上玩笑，下面说点正经的。关于《从零开始》是否有续集的问题……其实我是有一套创意预案的，但会不会真的动笔写暂时还不能确定。\n\n风暴我这里刚接了巨人的合同，所以下一本书将是定制《征途》，对你没看错，就是那个网游界的常青树《征途》。顺便报个料，《征途》可能很快会出电影和电视剧版，反正有这么个计划，啥时候各位能看到就不知道了。\n\n最后，再次感谢大家11年来的支持，也希望大家继续捧场，支持下风暴的新书《征途》。（ps：这次肯定不会写11年了。）(未完待续。)u",
			},
			wantError: nil,
		},
		{
			name: "happy flow",
			body: `<data>
				<div class="zhengwen_box"><div class="box_left"><div class="w_main">
					<div class="h1title"><h1 id="timu">chapter name</h1></div>
					<div class="contentbox"><div id="contentbox">chapter content</div></div>
				</div></div></div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "chapter name", Body: "chapter content",
			},
			wantError: nil,
		},
		{
			name: "title not found",
			body: `<data>
				<div class="zhengwen_box"><div class="box_left"><div class="w_main">
					<div class="contentbox"><div id="contentbox">chapter content</div></div>
				</div></div></div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "", Body: "chapter content",
			},
			wantError: vendor.ErrChapterTitleNotFound,
		},
		{
			name: "body not found",
			body: `<data>
				<div class="zhengwen_box"><div class="box_left"><div class="w_main">
					<div class="h1title"><h1 id="timu">chapter name</h1></div>
				</div></div></div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "chapter name", Body: "",
			},
			wantError: vendor.ErrChapterContentNotFound,
		},
		{
			name:      "all fields not found",
			body:      "<data></data>",
			want:      &vendor.ChapterInfo{},
			wantError: vendor.ErrFieldsNotFound,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got, err := p.ParseChapter(test.body)
			assert.Equal(t, test.want, got)
			assert.ErrorIs(t, err, test.wantError)
		})
	}
}

func TestParser_IsAvailable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		body string
		want bool
	}{
		{
			name: "return true",
			body: "UU看书",
			want: true,
		},
		{
			name: "return false",
			body: "",
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got := p.IsAvailable(test.body)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestParser_FindMissingIds(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ids  []int
		want []int
	}{
		{
			name: "no missing ids",
			ids:  []int{4, 2, 3, 1, 5},
			want: nil,
		},
		{
			name: "some id is missing",
			ids:  []int{3, 5, 1},
			want: []int{2, 4},
		},
		{
			name: "input ids contains negative",
			ids:  []int{3, -1},
			want: []int{1, 2},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			p := VendorService{}
			got := p.FindMissingIds(test.ids)
			assert.Equal(t, test.want, got)
		})
	}

}
