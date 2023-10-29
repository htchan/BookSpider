package xbiquge

import (
	"testing"

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
				Title:         "神印王座II皓月当空",
				Writer:        "唐家三少",
				Type:          "都市小说",
				UpdateDate:    "2023-08-03 10:45:03",
				UpdateChapter: "正文 第二百二十章 陷阱，绝境？",
			},
			wantError: nil,
		},
		{
			name: "happy flow",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<meta property="og:novel:update_time" content="date" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: nil,
		},
		{
			name: "title not found",
			body: `<data>
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<meta property="og:novel:update_time" content="date" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
			</data>`,
			want: &vendor.BookInfo{
				Writer: "author", Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookTitleNotFound,
		},
		{
			name: "writer not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:category" content="type" />
				<meta property="og:novel:update_time" content="date" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Type: "type",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookWriterNotFound,
		},
		{
			name: "type not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:update_time" content="date" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author",
				UpdateDate: "date", UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookTypeNotFound,
		},
		{
			name: "date not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<meta property="og:novel:latest_chapter_name" content="chapter name" />
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateChapter: "chapter name",
			},
			wantError: vendor.ErrBookDateNotFound,
		},
		{
			name: "chapter not found",
			body: `<data>
				<meta property="og:novel:book_name" content="book name" />
				<meta property="og:novel:author" content="author" />
				<meta property="og:novel:category" content="type" />
				<meta property="og:novel:update_time" content="date" />
			</data>`,
			want: &vendor.BookInfo{
				Title: "book name", Writer: "author", Type: "type",
				UpdateDate: "date",
			},
			wantError: vendor.ErrBookChapterNotFound,
		},
		{
			name:      "all fields not found",
			body:      "<data></data>",
			want:      &vendor.BookInfo{},
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
				{URL: "40007993.html", Title: "引子：皓月当空"}, {URL: "40007994.html", Title: "第一章 龙当当与龙空空"},
				{URL: "40007995.html", Title: "第二章 魔法圣殿"}, {URL: "40007996.html", Title: "第三章 光明庇护体质"},
				{URL: "40007997.html", Title: "第四章 龙家的“两大天才”"}, {URL: "40007998.html", Title: "第五章 第一堂课"},
				{URL: "40007999.html", Title: "第六章 圣殿大礼包"}, {URL: "40008000.html", Title: "第七章 灵炉，意外"},
				{URL: "40008001.html", Title: "第八章 骑士圣殿也来了"}, {URL: "40008002.html", Title: "第九章 灵炉，再见灵炉"},
				{URL: "40008003.html", Title: "第十章 绝望周末？"}, {URL: "40064650.html", Title: "第一百章 不可能完成的任务"},
				{URL: "40064651.html", Title: "第一百零一章 禁！光之礼赞"}, {URL: "40065639.html", Title: "第一百零二章 战后"},
				{URL: "40065640.html", Title: "第一百零三章 三姐弟的秘聊"}, {URL: "40066003.html", Title: "第一百零四章 要去获得坐骑了？"},
				{URL: "40066396.html", Title: "第一百零五章 骑士圣山"}, {URL: "40066397.html", Title: "第一百零六章 我不配和你不配"},
				{URL: "40067405.html", Title: "第一百零七章 五爪金龙"}, {URL: "40068006.html", Title: "第一百零八章 龙皇相邀"},
				{URL: "40068047.html", Title: "第一百零九章 龙当当的坐骑伙伴"}, {URL: "40068048.html", Title: "第一百一十章 龙空空的坐骑伙伴"},
				{URL: "40178727.html", Title: "第二百章 五层精神之海?"}, {URL: "40178728.html", Title: "第二百零一章 分身的精神边际"},
				{URL: "40178856.html", Title: "第二百零二章 龙当当的恐怖之处"}, {URL: "40178857.html", Title: "第二百零三章 我们缺功勋值吗?"},
				{URL: "40182283.html", Title: "第二百零四章 很贵的一战"}, {URL: "40182284.html", Title: "第二百零五章 龙当当的布局"},
				{URL: "40182285.html", Title: "第二百零六章 强势过关"}, {URL: "40182286.html", Title: "第二百零七章 正心殿"},
				{URL: "40182287.html", Title: "第二百零八章 守心"}, {URL: "40182288.html", Title: "第二百零九章 正式成为猎魔者"},
				{URL: "40182289.html", Title: "第二百一十章 士级猎魔团"}, {URL: "40221749.html", Title: "第二百一十一章 兑换奖励"},
				{URL: "40222892.html", Title: "第二百一十二章 套装，圣耀之心"}, {URL: "40224052.html", Title: "第二百一十四章 龙当当想要的奖励"},
				{URL: "40226906.html", Title: "第二百一十五章 幸福中翱翔"}, {URL: "40228397.html", Title: "第二百一十六章 新装备"},
				{URL: "40235866.html", Title: "第二百一十七章 突破六阶"}, {URL: "40260038.html", Title: "第二百一十八章 紧急任务"},
				{URL: "40283912.html", Title: "第二百一十九章 神女之威"}, {URL: "40290770.html", Title: "第二百二十章 陷阱，绝境？"},
			},
			wantError: nil,
		},
		{
			name: "happy flow",
			body: `<data>
				<div>
					<dd><a href="chapter url 1">chapter name 1</a></dd>
					<dd><a href="chapter url 2">chapter name 2</a></dd>
					<dd><a href="chapter url 3">chapter name 3</a></dd>
					<dd><a href="chapter url 4">chapter name 4</a></dd>
				</div>
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
				<div>
					<dd><a href="chapter url 1">chapter name 1</a></dd>
					<dd><a href="">chapter name 2</a></dd>
					<dd><a href="chapter url 3">chapter name 3</a></dd>
					<dd><a href="chapter url 4">chapter name 4</a></dd>
				</div>
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
				<div>
					<dd><a href="chapter url 1">chapter name 1</a></dd>
					<dd><a href="chapter url 2">chapter name 2</a></dd>
					<dd><a href="chapter url 3"></a></dd>
					<dd><a href="chapter url 4">chapter name 4</a></dd>
				</div>
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
			got, err := p.ParseChapterList(test.body)
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
			name: "happy flow with real data",
			body: string(testChapterBytes),
			want: &vendor.ChapterInfo{
				Title: "第二百二十章 陷阱，绝境？",
				Body:  "笔趣阁 www.xbiquge.bz，最快更新神印王座II皓月当空 ！\n\n神圣魔法是光系魔法的升华。就像冰系魔法是水系魔法的升华一样。神圣魔法更是被誉为最强大的魔法属性之一。\n\n 普通的光明魔法是克制不了亡灵生物的，这些能够在阳光下自由行动的亡灵生物本身实力极为强大。唯有神圣魔法，对其能够产生一定的克制。\n\n 一般的光系魔法师在施展光系魔法的时候，如果能够令自身魔法多一丝神圣属性，威能就会极大程度的提升。光系魔法师魔法中蕴含的神圣属性强度，往往决定者他们的天赋和潜能。\n\n 但是，像这种纯粹的神圣属性，李泓澈甚至连听都没听说过。因为他甚至感受不到这个魔法中的光明，感受到的，唯有那崇高的神圣。\n\n 如此强大的牧师吗？\n\n 也就在这个时候，空气中突然出现了一种奇异的声响，一个闪烁着奇异色彩的光球在空中划过，带着刺耳的厉啸，直奔那在光之礼赞中强行挣扎的巫妖飞去。\n\n 那光球内，蕴含着蓝、红、黄、青、金、黑六种色彩。在飞行的过程中，甚至有些不稳定的震颤着。\n\n 刚刚完成光之礼赞的凌梦露下意识的扭头看向远方。作为同时代最优秀的天之骄女，她和她之间，一直都有着莫名的比较。\n\n 空气中，突然响起一声诡异的嗡鸣。\n\n 那名巫妖看到飞向自己的光球，内心中骤然产生出一股强烈的危机感。它没有选择硬碰，几乎是瞬问身形向后飘飞。但是，受到光之礼赞影响，虽然凌梦露更多的针对的是下方的亡灵法阵，但这个魔法的神圣气息实在是太强大了，强大到对它的速度和灵魂都产生了影响。\n\n “呼呼呼！”空气中，恐怖的元素波动几乎在下一刻绽放。\n\n 那个六色光球在空中瞬问犹如凝固了一般停顿下来，紧接着，六种色彩的恐怖光芒就以它为中心爆发了。\n\n 爆发的中心点覆盖了巫妖，也覆盖了下方所有的英灵，甚至是李泓澈与那名亡者。\n\n 李泓澈脸色一变，身形暖问爆退，与此同时，一道银光从他脚下绽放，层层迭迭的银色甲胃飞速向上反卷，在那恐怖的元素风暴降临在自己身上之前，他的身体就已经被那散发着强大气息的银色甲胃完全包覆。是的，没错，他身上这甲胃，正是骑士圣殿最核心的强大力量之一，秘银基座战铠。李泓澈不只是一名猎魔者，同时也是骑士圣殿七十二位秘银基座骑士之一。\n\n 元素风暴狂躁无比的疯狂爆发，刚刚被弱化版光之礼赞洗礼的呆立原地的英灵遗骸们顿时在风暴肆虐的支离破碎。\n\n 在子桑琉荧看来，本次任务的关键就是不让亡灵生物带走这些只剩下遗体的前辈强者，而毁掉这些遗体就是对这些先辈最好的尊重。让他们不至于在死后还受到茶毒。伴随着成为猎魔者这九个月来面对的种种有关于亡灵生物的任务，她是完全支持火化所有先辈还骨的，正是因为有了众多英灵的墓葬，才让活着的人疲于奔命并且不断的在战斗中陨落，更痛苦的是一旦英灵遗骨被亡灵生物掠走，用不了多久，就会成为亡灵大军的一份子，而且必然会是强大的亡者，更加助长亡灵生物的实力。\n\n 老一代的强者，甚至是六大圣殿的主事者们对于先辈遗骨极为珍视，也正是他们的强硬，才呈现出现在这种要不断守护墓葬的不利局面，所以，遇到机会，子桑琉荧立刻就用自己的行动，来宣誓着自己的主张。\n\n 眼看着那一具具遗骨在元素风暴中破损，身体在秘银基座战铠保护中飞退的李泓澈已是看的目瞪口呆。虽然任务中没有明说要保护好这些遗骨妥善带回去，但是，她这分明就是直接针对还骨去的啊！甚至主要攻击点都不在那名巫妖身上。\n\n 遗骨在元素风暴中破碎，被激活的意思灵魂之力在光之礼赞中净化。那原本看起来十分恐怖的还骨群顷刻间连同地面的亡灵法阵被一扫而空，没有留下丝毫痕迹。\n\n 悬浮在空中的凌梦露与另一边的子桑琉荧，此时犹如双星闪耀一般绽放着夺目光彩。\n\n 但是，也就在下一刻，突然之间，一股极其强烈的压抑感突然出现了。\n\n 光之礼赞的光辉和元素风暴的爆裂渐渐消退着，但是，不知道什么时候，天空却随之暗了下来。\n\n 强烈的压抑感禁绕在已经纷纷落在地面的四支猎魔团心头，一种不妙的感觉也随之出现。\n\n 远处，先前逃跑在最前方的十几名亡者全都停了下来，缓步的向他们的方向走来，在这些亡者之中，有四名亡者走在前面，那压抑的气息，正是从它们身上释放出来的。\n\n 其中两名亡者抬手挥向空中，顿时，天空变得更加阴暗，大片、大片的乌云凝聚，遮挡住了太阳的光彩。他们身上的气息也开始如同井喷一般进发、提升，强盛无比的力量感在空气中剧烈的波动着。\n\n 这是逃亡的亡灵？还是……\n\n “陷阱。”龙当当脸色阴沉的沉声说道，此时他和伙伴们都已经下了青鹏，来到了完成一个大魔法正在恢复的凌梦露身边。\n\n 一道冰蓝色的光芒电射而出，直奔远处的亡者们飞射而去。正是月离刚刚完成的冰雷之矛。刺目的电光爆发之下，冰雷之矛速度犹如雷霆一般，几乎是瞬间降临。\n\n 但是，对方位于正中的那名亡者突然抬起手，一把就抓佳了冰雷之矛，就在冰爆术和雷霆即将爆发的时候，漆黑如墨的光芒从他掌中蔓延开来，覆盖了整根长矛，竟是让上面的所有魔法元素为之泯灭。\n\n 要知道，这可是一个威能足以达到七阶的双系混合单体攻击魔法啊！但却就是被这样轻而易举的掐灭了，对方是什么实力？\n\n 李泓澈也同样意识到了问题，这次的紧急任务，此时看起来，更像是一个钓鱼的陷阱。\n\n 十几名亡者已经缓缓的国了上来，逃跑？那肯定是来不及了。\n\n 没有半分犹豫，李泓澈激发了身上的子灵晶，发出了求援信号。但是，此处乃是旷野距离城市很远，以子灵晶的传信能力，就算是能将求援信号传回去，救援至少也要半个时辰恐怕才能赶来。\n\n 李泓澈深吸口气，眼底闪过一抹决绝之色。当他成为猎魔者的那一天，就已经做好了随时艳牲的准备。只是没想到，这一天竟然来的会这么快，他们才刚刚晋升帅级猎魔团不久，本来有着光明的前途。\n\n 所有人，向我靠拢。“李泓澈大喝一声。在他身后不远处的本团队牧师已经将一个又一个的辅助魔法增幅在他身上。\n\n 修为达到七阶红衣主教的牧师施展着洞察之眼，沉声道：“敌人有四名八阶亡者，剩余的亡者中还有六名七阶，其他是六阶。它们之前应该是用了什么特殊的装备掩盖了气息。我们麻烦了。”\n\n 十几名亡者似乎一点都不着急，此时距离他们还有近三百米，却是一步步的朝着他们走过来，但谁都知道，一旦他们试图逃离，这些亡者立刻就会展开最快的速度向他们发起冲击。\n\n 龙当当、子桑琉荧、李雨青听令。李泓澈在这个时候显得异常冷静。\n\n “在。”三名团长同时应声。\n\n 明白了吗？\n\n 李泓澈深吸口气，沉声道：“稍候，我和我的队友会发起全力政击，不惜代价的拖住它们。你们立刻骑乘青鹏撒离。不用管我们，不许回头，听听着他声音中的决绝，龙当当三人瞬问都有种全身起鸡皮疙瘩的感觉，因为他们立刻就明白，这位从一开始就有些看不上他们的圣殿骑士已经做好了要和自己团队以身殉道的准备。没有半分的迟疑他就已经做出了决定。面对四名八阶强者带领的亡灵强者们，一―二一一七猎魔团必将全军覆没，可是，他义无反顾。\n\n “李团长……”李两青忍不住急声叫道。\n\n 李泓澈突然笑了，他的面甲开启，露出了满脸释然的微笑，“你们记住，按照猎魔团的规矩。如果是多支猎魔团共同遭遇危险时，最高级别的猎魔团必须要为其他猎魔团争取撤离的机会，无畏牺牲！这是猎魔团的传统今天是我们，或许，下一次就是你们。”\n\n 说到这里，他的目光从子桑琉荧和凌梦露身上扫过，“先前因为你们年轻却小看你们了。你们比我想象中更加优秀，一定要成长起来。还有，你们之前做得很好，或许，这样做才能让我们那些死去的先辈英灵们不被打扰，更不会去破坏他们曾经用生命守护着的人类世界。”\n\n 说到这里，他的声音突然拔高起来，“在我冲出去的那一刻，你们立即撒离，不得有误。\n\n “砰――”低沉的闷响声从他身上暴起，灿烂的红色光芒在他身上点燃。\n\n 对于这样的光芒没有谁比龙当当更熟悉了，那赫然是牺牲的光彩，点燃的，是自己的灵力，更是自身的生命火焰。\n\n 灿银色的甲胃上燃烧着红色的光焰，李泓澈自身的气息疯狂的提升着，在他身边的队友们从始至终都没有多说一句什么，每个人的眼神都是坚定而决绝的。是的，这就是猎魔者以守护人类、守护联邦为己任的猎魔者们。\n\n 李泓澈最后回过头，看向龙当当，“龙团长，你我都是骑士。如果，我是说如果，未来你的实力足够达到的时候，请想办法帮我收回秘银基座战铠，现在我没办法让价带走它，因为没有它，我的力量不足以阻止住这些亡者，拜托了。”\n\n 四名团长之中，他和龙当当是骑士，秘银基座战铠对骑士圣殿来说不只是一件装备，更是骑士的象征之一。\n\n “李团长。”眼看着李泓澈已经回过身去，龙当当突然沉声喝道。\n\n 李泓澈眉头紧蹙，正当他准备回头的时候，却突然发现，在自己身边已经多了一人，与自己并肩而立，赫然正是龙当当。然后他就看到，龙当当胸口处，一尊洁白如玉的灵炉浮现而出，紧接着，又是两尊灵炉漂荡而至，分别来到他身边。\n\n 三尊灵炉瞬间重合，再骤然落在他身上，晶莹剔透的冰蓝色甲青伴随着他的身体膨账而迅速出现在他身上。他的气息也随之开始暴涨，双眸之中中，更是散发着灿烂的金色光辉。\n\n 这是什么？\n\n “你不会失去秘银基座战铠的。”龙当当的声音响起。\n\n 看着比自己足足高出一米多的沧月天使龙当当，李泓澈一时之间不禁有些顺住，下一秒才怒喝道：“胡闹，准备撤退。\n\n “我们不会走的。”子桑琉荧清冷的声音响起，幽幽的道：“这么多八阶亡灵，这是多少功勋啊！”―边说着，她也已经来到了李泓澈的另一边，她团队中的骑士同样顶在了前面。\n\n “如果我战死了，请在第一时间毁掉我的尸体，不要被亡灵生物带走。”李两青澹澹的说道。\n\n 能够成为猎魔者，就没有一个人是怕死的，他们都经历了无数考验，经历了正心殿的正心之旅，他们内心深处，都有着坚定的信仰，那就是守护！\n\n 尽管李泓澈此时心急如焚，但是，在他的内心处却有阵阵暖流涌动。\n\n 而也就在这时，似乎是感受到了猎魔者们这边的气势变化，对面的十几名亡者，突然动了。\n\n 四名八阶亡者之中，两名先前没有出手释放阴云的亡者突然暴起，在它们背后，各自出现了一对灰色的翅膀，翅膀用力一拍，下一瞬，带着狂暴的罡风，已经直奔众人而来。\n\n 八阶亡灵，同样拥有灵罡。\n\n 龙当当扭头看向李泓澈，微笑道：“李团长，我们一人一个。一说完这句话，他背后两对光翼勐然拍动，下一刻，人已经如同流星一般飞了出去。\n\n 。",
			},
			wantError: nil,
		},
		{
			name: "happy flow",
			body: `<data>
				<div class="bookname"><h1>chapter name</h1></div>
				<div id="content">chapter content</div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "chapter name", Body: "chapter content",
			},
			wantError: nil,
		},
		{
			name: "title empty",
			body: `<data>
			<div class="bookname"><h1></h1></div>
			<div id="content">chapter content</div>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "", Body: "chapter content",
			},
			wantError: vendor.ErrChapterTitleNotFound,
		},
		{
			name: "body empty",
			body: `<data>
			<div class="bookname"><h1>chapter name</h1></div>
			<div id="content"></div>
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
			body: "笔趣阁",
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
