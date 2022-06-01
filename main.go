package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var f *os.File
var fList *os.File
var tableName = "hs_code"
var tableName2 = "hs_code_list"
var baseUrl = "https://www.hsbianma.com/"

type HSStaticData struct {
	Id       int             `json:"id"`
	Name     string          `json:"name"`
	Children []*HSStaticData `json:"children"`
}

func main() {
	// 创建 sql 文件
	f, _ = os.Create(tableName + ".sql")
	fList, _ = os.Create(tableName2 + ".sql")

	// 该数据可以直接在 js 中获取到，未爬取了
	data := []*HSStaticData{
		{1, "第1类 - 活动物;动物产品", []*HSStaticData{
			{1, "第1章 - 活动物", nil},
			{2, "第2章 - 肉及食用杂碎", nil},
			{3, "第3章 - 鱼、甲壳动物,软体动物及其他水生无脊椎动物", nil},
			{4, "第4章 - 乳品;蛋品;天然蜂蜜;其他食用动物产品", nil},
			{5, "第5章 - 其他动物产品", nil},
		}},
		{2, "第2类 - 植物产品", []*HSStaticData{
			{6, "第6章 - 活树及其他活植物;鳞茎、根及类似品;插花及装饰用簇叶", nil},
			{7, "第7章 - 食用蔬菜、根及块茎", nil},
			{8, "第8章 - 食用水果及坚果;柑桔属水果或甜瓜的果皮", nil},
			{9, "第9章 - 咖啡、茶、马黛茶及调味香料", nil},
			{10, "第10章 - 谷物", nil},
			{11, "第11章 - 制粉工业产品;麦芽;淀粉;菊粉;面筋", nil},
			{12, "第12章 - 含油子仁及果实;杂项子仁及果实;工业用或药用植物;稻草、秸杆及饲料", nil},
			{13, "第13章 - 虫胶;树胶、树脂及其他植物液、汁", nil},
			{14, "第14章 - 编结用植物材料;其他植物产品", nil},
		}},
		{3, "第3类 - 动、植物油、脂及其分解产品;精制的食用油脂;动、植物蜡", []*HSStaticData{
			{15, "第15章 - 动、植物油、脂及其分解产品;精制的食用油脂;动、植物蜡", nil},
		}},
		{4, "第4类 - 食品;饮料、酒及醋;烟草、烟草及烟草代用品的制品", []*HSStaticData{
			{16, "第16章 - 肉、鱼、甲壳动物,软体动物及其他水生无脊椎动物的制品", nil},
			{17, "第17章 - 糖及糖食", nil},
			{18, "第18章 - 可可及可可制品", nil},
			{19, "第19章 - 谷物、粮食粉、淀粉或乳的制品;糕饼点心", nil},
			{20, "第20章 - 蔬菜、水果、坚果或植物其他部分的制品", nil},
			{21, "第21章 - 杂项食品", nil},
			{22, "第22章 - 饮料、酒及醋", nil},
			{23, "第23章 - 食品工业的残渣及废料;配制的动物饲料", nil},
			{24, "第24章 - 烟草及烟草代用品的制品", nil},
		}},
		{5, "第5类 - 矿产品", []*HSStaticData{
			{25, "第25章 - 盐;硫磺;泥土及石料;石膏料、石灰及水泥", nil},
			{26, "第26章 - 矿砂、矿渣及矿灰", nil},
			{27, "第27章 - 矿物燃料、矿物油及其蒸馏产品;沥青物质;矿物蜡", nil},
		}},
		{6, "第6类 - 化学工业及其相关工业的产品", []*HSStaticData{
			{28, "第28章 - 无机化学品;贵金属、稀土金属、放射性元素及其同位素的有机及无机化合物", nil},
			{29, "第29章 - 有机化合物", nil},
			{30, "第30章 - 药品", nil},
			{31, "第31章 - 肥料", nil},
			{32, "第32章 - 鞣料浸膏及染料浸膏;鞣酸及其衍生物;染料、颜料及其他着色料;油漆及清漆;油灰及其他类似胶粘剂;墨水、油墨", nil},
			{33, "第33章 - 精油及香膏;芳香料制品及化妆盥洗品", nil},
			{34, "第34章 - 肥皂、有机表面活性剂、洗涤剂、 润滑剂、人造蜡、调制蜡、光洁剂、蜡烛及类似品、塑型用膏、“牙科用蜡”及牙科用熟石膏制剂", nil},
			{35, "第35章 - 蛋白类物质;改性淀粉;胶;酶", nil},
			{36, "第36章 - 炸药;烟火制品;火柴;引火合金;易燃材料制品", nil},
			{37, "第37章 - 照相及电影用品", nil},
			{38, "第38章 - 杂项化学产品", nil},
		}},
		{7, "第7类 - 塑料及其制品:橡胶及其制品", []*HSStaticData{
			{39, "第39章 - 塑料及其制品", nil},
			{40, "第40章 - 橡胶及其制品", nil},
		}},
		{8, "第8类 - 生皮、皮革、毛皮及其制品;鞍具及挽具;旅行用品、手提包及类似容器;动物肠线", []*HSStaticData{
			{41, "第41章 - 生皮(毛皮除外)及皮革", nil},
			{42, "第42章 - 皮革制品;鞍具及挽具;旅行用品、手提包及类似容器;动物肠线(蚕胶丝除外)制品", nil},
			{43, "第43章 - 毛皮、人造毛皮及其制品", nil},
		}},
		{9, "第9类 - 木及木制品;木炭;软木及软木制品;稻草、秸秆、针茅或其他编结材料制品;篮筐及柳条编结品", []*HSStaticData{
			{44, "第44章 - 木及木制品;木炭", nil},
			{45, "第45章 - 软木及软木制品", nil},
			{46, "第46章 - 稻草、秸秆、针茅或其他编结材料制品:篮筐及柳条编结品", nil},
		}},
		{10, "第10类 - 木浆及其他纤维状纤维素浆;回收", []*HSStaticData{
			{47, "第47章 - 木浆及其他纤维状纤维素浆;回收(废碎)纸或纸板", nil},
			{48, "第48章 - 纸及纸板;纸浆、纸或纸板制品", nil},
			{49, "第49章 - 书籍、报纸、印刷图画及其他印刷品;手稿、打字稿及设计图纸", nil},
		}},
		{11, "第11类 - 纺织原料及纺织制品", []*HSStaticData{
			{50, "第50章 - 蚕丝", nil},
			{51, "第51章 - 羊毛、动物细毛或粗毛;马毛纱线及其机织物", nil},
			{52, "第52章 - 棉花", nil},
			{53, "第53章 - 其他植物纺织纤维;纸纱线及其机织物", nil},
			{54, "第54章 - 化学纤维长丝", nil},
			{55, "第55章 - 化学纤维短纤", nil},
			{56, "第56章 - 絮胎、毡呢及无纺织物;特种纱线;线、绳、索、缆及其制品", nil},
			{57, "第57章 - 地毯及纺织材料的其他铺地制品", nil},
			{58, "第58章 - 特种机织物;簇绒织物;花边;装饰毯;装饰带;剌绣品", nil},
			{59, "第59章 - 浸渍、涂布、包覆或层压的纺织物;工业用纺织制品", nil},
			{60, "第60章 - 针织物及钩编织物", nil},
			{61, "第61章 - 针织或钩编的服装及衣着附件", nil},
			{62, "第62章 - 非针织或非钩编的服装及衣着附件", nil},
			{63, "第63章 - 其他纺织制成品;成套物品;旧衣着及旧纺织品;碎织物", nil},
		}},
		{12, "第12类 - 鞋、帽、伞、杖、鞭及其零件;已加工的羽毛及其制品;人造花;人发制品", []*HSStaticData{
			{64, "第64章 - 鞋靴、护腿和类似品及其零件", nil},
			{65, "第65章 - 帽类及其零件", nil},
			{66, "第66章 - 雨伞、阳伞、手杖、鞭子、马鞭及其零件", nil},
			{67, "第67章 - 已加工羽毛、羽绒及其制品;人造花;人发制品", nil},
		}},
		{13, "第13类 - 石料、石膏、水泥、石棉、云母及类似材料的制品;陶瓷产品;玻璃及其制品", []*HSStaticData{
			{68, "第68章 - 石料、石膏、水泥、石棉、云母及类似材料的制品", nil},
			{69, "第69章 - 陶瓷产品", nil},
			{70, "第70章 - 玻璃及其制品", nil},
		}},
		{14, "第14类 - 天然或养殖珍珠、宝石或半宝石、贵金属、包贵金属及其制品;仿首饰;硬币", []*HSStaticData{
			{71, "第71章 - 天然或养殖珍珠、宝石或半宝石、贵金属、包贵金属及其制品;仿首饰;硬币", nil},
		}},
		{15, "第15类 - 贱金属及其制品", []*HSStaticData{
			{72, "第72章 - 钢铁", nil},
			{73, "第73章 - 钢铁制品", nil},
			{74, "第74章 - 铜及其制品", nil},
			{75, "第75章 - 镍及其制品", nil},
			{76, "第76章 - 铝及其制品", nil},
			{78, "第78章 - 铅及其制品", nil},
			{79, "第79章 - 锌及其制品", nil},
			{80, "第80章 - 锡及其制品", nil},
			{81, "第81章 - 其他贱金属、金属陶瓷及其制品", nil},
			{82, "第82章 - 贱金属工具、器具、利口器、餐匙、餐叉及其零件", nil},
			{83, "第83章 - 贱金属杂项制品", nil},
		}},
		{16, "第16类 - 机器、机械器具、电气设备及其零件;录音机及放声机、电视图像、声音的录制和重放设备及其零件", []*HSStaticData{
			{84, "第84章 - 核反应堆、锅炉、机器、机械器具及其零件", nil},
			{85, "第85章 - 电机、电气设备及其零件;录音机及放声机、电视图像、声音的录制和重放设备及其零件、附件", nil},
		}},
		{17, "第17类 - 车辆、航空器、船舶及有关运输设备", []*HSStaticData{
			{86, "第86章 - 铁道及电车道机车、车辆及其零件;铁道及电车道轨道固定装置及其零件、附件;各种机械(包括电动机械)交通信号设备", nil},
			{87, "第87章 - 车辆及其零件、附件、但铁道及电车道车辆除外", nil},
			{88, "第88章 - 航空器、航天器及其零件", nil},
			{89, "第89章 - 船舶及浮动结构体", nil},
		}},
		{18, "第18类 - 光学、照相、电影、计量、检验、医疗或外科用仪器及设备、精密仪器及设备;钟表;乐器;上述物品", []*HSStaticData{
			{90, "第90章 - 光学、照相、电影、计量、检验、医疗或外科用仪器及设备、精密仪器及设备;上述物品的零件、附件", nil},
			{91, "第91章 - 钟表及其零件", nil},
			{92, "第92章 - 乐器及其零件、附件", nil},
		}},
		{19, "第19类 - 武器、弹药及其零件、附件", []*HSStaticData{
			{93, "第93章 - 武器、弹药及其零件、附件", nil},
		}},
		{20, "第20类 - 杂项制品", []*HSStaticData{
			{94, "第94章 - 家具;寝具、褥垫、弹簧床垫、软座垫及类似的填充制品;未列名灯具及照明装置;发光标志、发光铭牌及类似品;活动房屋", nil},
			{95, "第95章 - 玩具、游戏品、运动品及其零件、附件", nil},
			{96, "第96章 - 杂 项 制 品", nil},
		}},
		{21, "第21类 - 艺术品、收藏品及古物", []*HSStaticData{
			{97, "第97章 - 艺术品、收藏品及古物", nil},
		}},
		{22, "第22类 - 特殊交易品及未分类商品", []*HSStaticData{
			{98, "第98章 - 特殊交易品及未分类商品", nil},
		}},
	}
	dataLen := len(data)
	fmt.Println("项目已启动，请耐心等待，待执行：" + strconv.Itoa(dataLen))
	for index, item := range data {
		_, _ = io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`parent_id`) values("+strconv.Itoa(item.Id)+
			",'"+item.Name+"',"+strconv.Itoa(0)+");\r\n")
		itemLen := len(item.Children)
		for cIndex, itemChildren := range item.Children {
			fmt.Println("==============> 进度：" + strconv.Itoa(index+1) + "/" + strconv.Itoa(dataLen) +
				" 子进度：" + strconv.Itoa(cIndex+1) + "/" + strconv.Itoa(itemLen))

			id := item.Id*100 + itemChildren.Id
			_, _ = io.WriteString(f, "INSERT INTO "+tableName+"(`id`,`name`,`parent_id`) values("+
				strconv.Itoa(id)+
				",'"+itemChildren.Name+"',"+strconv.Itoa(item.Id)+");\r\n")

			execPage := findList(itemChildren.Id, id, 1)
			fmt.Println(itemChildren.Name + " 已完成，共执行页数：" + strconv.Itoa(execPage))
			time.Sleep(time.Second * 30)
		}
	}
}

// 查找第一级
func findList(docInt int, id int, page int) int {
	var keywords string
	if docInt > 10 {
		keywords = strconv.Itoa(docInt)
	} else {
		keywords = "0" + strconv.Itoa(docInt)
	}
	url := baseUrl + "/search/" + strconv.Itoa(page) + "?keywords=" + keywords
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	if doc.Find(".result-grid").Length() <= 0 {
		return page - 1
	}
	doc.Find(".result-grid").Each(func(i int, selection *goquery.Selection) {
		// 每行数据
		var columnData [8]string
		selection.Find("td").Each(func(i int, column *goquery.Selection) {
			// 每行数据中的每列数据
			if i <= 5 {
				var fieldData string
				if i == 0 {
					// 判断是否过期
					index := strings.Index(column.Text(), "[过期]")
					if index > 0 {
						// 则存在
						fieldData = strings.Replace(column.Text(), "[过期]", "", -1)
					} else {
						fieldData = column.Text()
					}
					fieldData = strings.TrimSpace(fieldData)
					columnData[6] = strconv.Itoa(id)
				} else {
					fieldData = strings.TrimSpace(column.Text())
				}
				columnData[i] = fieldData
			} else {
				//fmt.Println("td1:" + strconv.Itoa(i) + " ==> " + columnData)
				// 获取详情信息
				detailUrl, exist := column.Find("a").Attr("href")
				if exist {
					columnData[7] = strconv.Quote(string(findDetail(detailUrl)))
				}
			}
		})
		//fmt.Println("INSERT INTO "+tableName2+"(`code`,`name`,`unit`,`export_tax`,`supervise`,`quarantine`)"+
		//	" values('"+columnData[0]+"','"+columnData[1]+"','"+columnData[2]+"','"+columnData[3]+"','"+columnData[4]+"','"+columnData[5]+"')")
		// 点击进入详情
		_, _ = io.WriteString(fList, "INSERT INTO "+tableName2+"(`code`,`pid`,`name`,`unit`,`export_tax`,`supervise`,`quarantine`, `detail`)"+
			" values('"+columnData[0]+"','"+columnData[6]+"','"+columnData[1]+"','"+columnData[2]+"','"+columnData[3]+
			"','"+columnData[4]+"','"+columnData[5]+"',"+columnData[7]+");\r\n")
		time.Sleep(time.Millisecond * 1)
	})

	fmt.Println("写入完成：page " + strconv.Itoa(page))
	time.Sleep(time.Second * 10)
	page++
	return findList(docInt, id, page)
}

type detail struct {
	Title string `json:"title"`
	Info  map[int]detailInfo
}

type detailInfo struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func findDetail(detailUrl string) []byte {
	url := baseUrl + detailUrl
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	detailData := make(map[int]detail)
	doc.Find("#code-info").Find("h3").Each(func(i int, selection *goquery.Selection) {
		k := selection.Text()
		data := doc.Find("#code-info").Find(".cbox").Eq(i)
		//var detailInfoData []detailInfo
		detailInfoData := make(map[int]detailInfo)
		data.Find("tr").Each(func(j int, tr *goquery.Selection) {
			keys := tr.Find("td").Eq(0).Text()
			value := tr.Find("td").Eq(1).Text()
			detailInfoData[j] = detailInfo{Key: keys, Value: value}
		})
		detailData[i] = detail{Title: k, Info: detailInfoData}
		//_, _ = io.WriteString(fList, "INSERT INTO "+tableName2+"(`code`,`pid`,`name`,`unit`,`export_tax`,`supervise`,`quarantine`)"+
		//	" values('"+columnData[0]+"','"+columnData[6]+"','"+columnData[1]+"','"+columnData[2]+"','"+columnData[3]+"','"+columnData[4]+"','"+columnData[5]+"');\r\n")
	})
	// 转成 Json
	detailStr, _ := json.Marshal(detailData)
	return detailStr
}
