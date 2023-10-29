package hjwzw

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
				Title:         "恐怖修仙世界",
				Writer:        "龍蛇枝",
				Type:          "仙俠",
				UpdateDate:    "2021-04-06",
				UpdateChapter: "完本感言",
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
				{URL: "/Book/Read/37656,16491126", Title: "第1章 黑暗恐懼"}, {URL: "/Book/Read/37656,16491127", Title: "第2章 陰鬼"},
				{URL: "/Book/Read/37656,16491128", Title: "第3章 小燈符"}, {URL: "/Book/Read/37656,16491129", Title: "第4章 怪夢"},
				{URL: "/Book/Read/37656,16491130", Title: "第5章 束發日"}, {URL: "/Book/Read/37656,16491131", Title: "第6章 血色字數"},
				{URL: "/Book/Read/37656,16491132", Title: "第7章 壽命天定"}, {URL: "/Book/Read/37656,16491133", Title: "第8章 短命種的義務"},
				{URL: "/Book/Read/37656,17898429", Title: "第998章 陰隱線"}, {URL: "/Book/Read/37656,17900924", Title: "第999章 蝕日與千譎"},
				{URL: "/Book/Read/37656,17900925", Title: "第1000章 想當黃雀?"}, {URL: "/Book/Read/37656,17900926", Title: "第1001章 我咒你"},
				{URL: "/Book/Read/37656,17903108", Title: "第1002章 遲來的人"}, {URL: "/Book/Read/37656,17903109", Title: "第1003章 吞食"},
				{URL: "/Book/Read/37656,17903110", Title: "第1004章 收獲與筆記"}, {URL: "/Book/Read/37656,17903111", Title: "第1005章 吃骨頭"},
				{URL: "/Book/Read/37656,20061604", Title: "第1998章 回歸"}, {URL: "/Book/Read/37656,20061706", Title: "第1999章 名字"},
				{URL: "/Book/Read/37656,20061707", Title: "第2000章 攤牌了"}, {URL: "/Book/Read/37656,20065889", Title: "第2001章 小時間"},
				{URL: "/Book/Read/37656,20065890", Title: "第2002章 大道選擇"}, {URL: "/Book/Read/37656,20065916", Title: "第2003章 重返主星界"},
				{URL: "/Book/Read/37656,20065920", Title: "第2004章 云元子"}, {URL: "/Book/Read/37656,20070318", Title: "第2005章 造神宗的邀請"},
				{URL: "/Book/Read/37656,20202511", Title: "第2066章 再進階"}, {URL: "/Book/Read/37656,20202512", Title: "第2067章 蘇醒"},
				{URL: "/Book/Read/37656,20206411", Title: "第2068章 各自對手"}, {URL: "/Book/Read/37656,20206412", Title: "第2069章 釣竿"},
				{URL: "/Book/Read/37656,20209393", Title: "第2070章 譎元紀"}, {URL: "/Book/Read/37656,20209394", Title: "第2071章 唯有超脫"},
				{URL: "/Book/Read/37656,20209409", Title: "第2072章 第三刀"}, {URL: "/Book/Read/37656,20209874", Title: "第2073章 最終"},
				{URL: "/Book/Read/37656,20212949", Title: "完本感言"},
			},
			wantError: nil,
		},
		{
			name: "happy flow",
			body: `<data>
				<div id="tbchapterlist"><table><tbody><tr>
					<td><a href="chapter url 1">chapter name 1</a></td>
					<td><a href="chapter url 2">chapter name 2</a></td>
					<td><a href="chapter url 3">chapter name 3</a></td>
					<td><a href="chapter url 4">chapter name 4</a></td>
				</tr></tbody></table></div>
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
				<div id="tbchapterlist"><table><tbody><tr>
					<td><a href="chapter url 1">chapter name 1</a></td>
					<td><a href="">chapter name 2</a></td>
					<td><a href="chapter url 3">chapter name 3</a></td>
					<td><a href="chapter url 4">chapter name 4</a></td>
				</tr></tbody></table></div>
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
				<div id="tbchapterlist"><table><tbody><tr>
					<td><a href="chapter url 1">chapter name 1</a></td>
					<td><a href="chapter url 2">chapter name 2</a></td>
					<td><a href="chapter url 3"></a></td>
					<td><a href="chapter url 4">chapter name 4</a></td>
				</tr></tbody></table></div>
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
				Title: "第1章 黑暗恐懼",
				Body:  "請記住本站域名: 黃金屋\n    第1章 黑暗恐懼\n 周凡勉力眨了一下眼睛，黃泥與干草混合搭建的房墻上只有一扇小窗，窗內投進一束光，屋頂的天窗也有白光泄進來，塵埃在光線內微微蕩漾。\n 但屋內大多數地方一片昏暗，黑得看不見任何東西。\n 周凡依然覺得頭腦昏昏沉沉的，他來到這個世界三天了，但還是有些搞不清狀況。\n 他只知道這具身體的名字同樣叫周凡，晚上‘父母’才會工作回來，而他之所以躺在床上，他模模糊糊從那對夫婦口中得知是因為他傷了腦袋。\n 也幸好是如此，周凡在第一次醒來后能推托說自己什么都忘記了，否則面對前身的‘父母’，他沒有任何關于前身的記憶，他真的不知道說什么才好。\n 為了避免懷疑，這三天周凡甚至不敢多說話，晚上都是靜靜地聽著‘父母’說話，但可惜的是那個老農打扮的父親是一個沉默寡言的男人，因此父母兩人的對話寥寥可數，周凡暫時無法從中得到太多有用的信息。\n 周凡掙扎著坐了起來，這一坐起，他的臉上露出痛苦之色，他伸出左手捂住額頭，腦袋就像被針刺了一樣。\n 他的手溫很低，有股冰冷順著額頭蔓延，使得那針刺的痛感減輕了不少。\n 又過了一會，腦袋的痛感幾乎微不可察。\n 周凡的手順著額頭而上，摸了摸沒有任何發絲的光腦袋，漸漸手觸碰到后腦勺一道一指長的疤痕。\n 他看不見傷疤，但是由觸感中能感覺到傷疤比頭發大上一些，若不是認真撫摸，還無法發現。\n 怎么受傷的？\n 周凡還不清楚，但要不是受傷，他的靈魂也到不了這個身軀內，他應該已經死了的。\n 周凡放下手，掀開床前黃葛布織造的深黃色帳幔，沒有帳幔的遮擋，視線變得清晰了一些，透過微弱的光線，他看著那些簡陋的木家具微微皺眉。\n 這更讓周凡確定自己處于一個比較貧窮的環境，只是屋內實在太暗了，晚上回來他看到‘父母’似乎點燃的是油燈。\n 不過周凡又不太敢肯定，這幾天他躺在床上昏頭昏腦的，意識都是迷迷糊糊的，白天很少有清醒的時候，大多數醒了一會，又睡了過去。\n 直到今天才好了很多，周凡看著屋內光線照不到的地方，那些地方暗得就似一團渲染開的墨水，他的腦袋開始一陣陣發麻。\n 他在害怕，就好像黑幕中有什么可怕的東西在窺視著他，會突然竄出來傷害他一樣。\n 這種恐懼沒有任何道理可言，周凡苦笑了一聲，他惜命，不過因為職業的原因從來就不是膽小的人，但身體卻有著這樣的反應，難道是重生到這具身體帶來的副作用？\n 這三天來，周凡嘗試了不少次，只要他默默注視著屋內的黑暗，就會有這樣的感覺。\n 又或者這是昏暗的環境影響他心情而導致的，周凡晃了晃頭，他沒有多想下去，而是嘗試著站了起來。\n 周凡的雙腿有些發軟，嘗試了好幾次才站了起來，他向前踏出一步，卻差點栽倒在地上，好不容易維持平衡，又繼續向前，走起來歪歪扭扭的，就像喝醉了酒一樣。\n 待周凡越過內屋門檻，來到屋子正門時已經滿頭大汗。\n 他透過微弱的光線，輕輕拉一下兩扇木門，門沒有鎖，一下子就被拉開。\n 外面耀眼的光一下子照進來，周凡瞇了瞇眼才適應了這強烈的光線。\n 一碧如洗的晴空，一排排的黃泥房子，隱約中還攜著雞鳴犬吠之聲。\n 借著明亮的光，周凡低頭看清了自己身上的衣服，這是一件褐色的短窄粗衣，現代社會恐怕做工再粗糙的衣服也不會有這么粗糙。\n 周凡站得有些累，他干脆一屁股坐在門檻上。\n 現在是白天，村里顯得有些安靜，他足足坐了一小時，才會有幾個人在他門前經過，那幾個人大都穿著短褐粗衣，手上拿著鋤頭之類的農具，他們見了周凡，有的臉色木然，有的只是對周凡笑笑，周凡回以笑容。\n 但周凡在那些人走了之后，他只是嘆了嘆氣，因為那些人的穿著打扮已經告訴了他一個早已經有所猜測的事實：他已經不是處在現代世界，而是到了一個古代世界。\n 不過周凡沒有很焦慮，前世妹妹和奶奶都死了之后，他報仇后在那世界就再也沒有任何的牽掛，對他來說，離開那沒有依戀的世界也算不了什么大事。\n 只是這是什么朝代？\n 歷史知識貧乏的周凡有些難以判斷。\n 自己以后該怎么辦呢？\n 周凡思緒煩亂想了好一會，他眼皮子開始直打架，他又覺得疲憊了。\n 周凡扶著木門框站了起來，把門關上，門一關上，就像從光明走到了黑暗之中，凝視著黑暗，那種讓他感到顫栗的感覺又從心底深處浮現了出來。\n 周凡盡量讓自己看著有微光的地方，那感覺才消退。\n 他摸黑又躺回床上，眼睛瞄著的是天窗上面的那束白光，他心里在想那種恐懼感是怎么回事？\n 黑暗中什么都不會有，他為什么會感到害怕呢？這實在太奇怪了……\n 周凡緩緩閉上了眼睛，閉眼同樣是一片黑暗，但他卻不會感到害怕，否則他連睡覺都不用睡了。\n 原本就感到疲憊的周凡很快就沉沉睡去。\n 待到再次醒來的時候，周凡取出‘父母’為他準備的食物，那是好像飯團一樣的東西，只是顏色卻是黃色的，有些像粟米。\n ‘父母’早上就說過今天中午不會回來，讓周凡自己起來吃飯。\n 周凡慢慢吃著飯團，這飯團的谷米帶著細小的谷殼，很難吞咽，只能盡量嚼碎才能吞下去。\n 很難吃的食物，但周凡沒有嫌棄，小的時候，僅靠奶奶養家，家里很窮，偶爾會餓肚子，那時起他就知道食物的珍貴，長大后就從來不敢做浪費食物的事情。\n 吃完后周凡覺得自己的精神好了很多，他又起來走動了一會，打開門發現已是黃昏時分，天邊的云彩被夕陽染得就像紅火焰一樣。\n 周凡站著看了一會，前幾天‘父母’都是勞作到夜色降臨才會回來，他又關上了門，四周黑漆漆的，他放棄了尋找油燈的想法，就算找到了，沒有打火工具也沒用。\n 現在的他什么事都做不了，只能又躺下來休息。\n 屋內也越來越暗，天窗上的光已微不可見，那種恐懼的感覺再度蔓延，有一瞬間，周凡甚至覺得自己看不見的臉色很為蒼白。\n 他不敢再睜眼，而是閉上了眼睛，只有閉眼才讓他沒有那么畏懼。\n “會是黑暗恐懼癥嗎？”\n 周凡眉頭微皺，腦海里浮現出這個想法，他曾經聽過這種病，這是一種心理疾病，怕黑，只要待在黑暗中就會產生緊張害怕等恐慌情緒。\n 只是他以前壓根就沒有這樣的毛病，會是前身的原因嗎？\n 但心理應該是思想主導才對的，現在這具身體里面可是他的靈魂，前身早已經死了，為什么還會感到害怕？\n 就在這時，他聽到‘吱呀’一聲傳來。\n 那是木門被推開的聲音，是他們回來了嗎？",
			},
			wantError: nil,
		},
		{
			name: "happy flow",
			body: `<data>
				<table><tbody><tr><td>
					<h1>chapter name</h1>
					<div></div><div></div><div></div><div></div>
					<div><p>chapter content</p></div>
				</td></tr></tbody></table>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "chapter name", Body: "chapter content",
			},
			wantError: nil,
		},
		{
			name: "title not found",
			body: `<data>
				<table><tbody><tr><td>
					<div></div><div></div><div></div><div></div><div></div>
					<div><p>chapter content</p></div>
				</td></tr></tbody></table>
			</data>`,
			want: &vendor.ChapterInfo{
				Title: "", Body: "chapter content",
			},
			wantError: vendor.ErrChapterTitleNotFound,
		},
		{
			name: "body not found",
			body: `<data>
				<table><tbody><tr><td>
					<h1>chapter name</h1>
					<div></div><div></div><div></div><div></div>
					<div></div>
				</td></tr></tbody></table>
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
			body: "黃金屋",
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
