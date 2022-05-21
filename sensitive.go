package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/axgle/mahonia"
	"github.com/hankcs/gohanlp/hanlp"

	//thulac "github.com/unionj-cloud/thulacgo"
	thulac "github.com/unionj-cloud/thulacgo"
)

var (
	b_list      []string
	w_list      []string
	o_word_list []string //一個字array
	t_word_list []string //兩個字array

	username      []string
	id_name       []string
	phone         []string
	mail          []string
	address       []string
	vaildate      []string
	cmp_username  bool
	namelist_path string = "./Trad_Chinese_Names_Corpus_Gender.txt"
)

func ConvertToString(src string, srcCode string, tagCode string) string {

	srcCoder := mahonia.NewDecoder(srcCode)

	srcResult := srcCoder.ConvertString(src)

	tagCoder := mahonia.NewDecoder(tagCode)

	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

	result := string(cdata)

	return result

}

func Sen(url string) map[string][]string {
	username = username[:0]
	id_name = id_name[:0]
	phone = phone[:0]
	mail = mail[:0]
	address = address[:0]
	vaildate = vaildate[:0]

	random := browser.Random() //隨機UA

	// {"username": [], "id":[], "ph_no":[], "em":[], "address":[], "bir":[]}

	m := make(map[string][]string)

	resp, err := http.Get(url)
	resp.Header.Add("User-Agent", random)
	if err != nil {
		fmt.Println(err)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	bodystr := ConvertToString(string(respData), "big5", "utf-8") //網頁編碼

	// 中文正則
	pattern := "[\u9673|\u6797|\u9ec3|\u5f35|\u674e|\u738b|\u5433|\u5289|\u8521|\u694a|\u8a31|\u912d|\u8b1d|\u6d2a|\u90ed|\u90b1|\u66fe|\u5ed6|\u8cf4|\u5f90|\u5468|\u8449|\u8607|\u838a|\u5442|\u6c5f|\u4f55|\u856d|\u7f85|\u9ad8|\u6f58|\u7c21|\u6731|\u937e|\u6e38|\u5f6d|\u8a79|\u80e1|\u65bd|\u6c88|\u4f59|\u76e7|\u6881|\u8d99|\u984f|\u67ef|\u7fc1|\u9b4f|\u5b6b|\u6234|\u8303|\u65b9|\u5b8b|\u9127|\u675c|\u5085|\u4faf|\u66f9|\u859b|\u4e01|\u5353|\u962e|\u99ac|\u8463|\u6e29|\u5510|\u85cd|\u77f3|\u8523|\u53e4|\u7d00|\u59da|\u9023|\u99ae|\u6b50|\u7a0b|\u6e6f|\u9ec4|\u7530|\u5eb7|\u59dc|\u767d|\u6c6a|\u9112|\u5c24|\u5deb|\u9418|\u9ece|\u6d82|\u9f94|\u56b4|\u97d3|\u8881|\u91d1|\u7ae5|\u9678|\u590f|\u67f3|\u51c3|\u90b5|\u9322|\u4f0d|\u502a|\u6eab|\u4e8e|\u8b5a|\u99f1|\u718a|\u7518|\u79e6|\u9867|\u6bdb|\u7ae0|\u53f2|\u5b98|\u842c|\u4fde|\u96f7|\u7c98|\u9952|\u95d5|\u51cc|\u5d14|\u5c39|\u5b54|\u8f9b|\u6b66|\u8f9c|\u9676|\u6613|\u6bb5|\u9f8d|\u97cb|\u845b|\u6c60|\u5b5f|\u891a|\u6bb7|\u9ea5|\u8cc0|\u8cc8|\u83ab|\u6587|\u7ba1|\u95dc|\u5411|\u5305|\u4e18|\u6885|\u83ef|\u88f4|\u6a0a|\u623f|\u5168|\u4f58|\u5de6|\u82b1|\u9b6f|\u5b89|\u9b91|\u90dd|\u7a46|\u5857|\u90a2|\u84b2|\u6210|\u8c37|\u5e38|\u95bb|\u7df4|\u76db|\u9114|\u803f|\u8076|\u7b26|\u7533|\u795d|\u7e46|\u967d|\u89e3|\u66f2|\u5cb3|\u9f4a|\u7c43|\u55ae|\u8212|\u7562|\u55ac|\u9f8e|\u7fdf|\u725b|\u911e|\u7559|\u5b63|\u8983|\u535c|\u9805|\u51c3|\u55bb|\u5546|\u6ed5|\u7126|\u8eca|\u8cb7|\u865e|\u82d7|\u621a|\u725f|\u96f2|\u5df4|\u529b|\u827e|\u6a02|\u81e7|\u6a13|\u8cbb|\u5c48|\u5b97|\u5e78|\u885b|\u5c1a|\u9773|\u7941|\u8af6|\u6842|\u6c99|\u6b12|\u5bae|\u8def|\u5201|\u9f90|\u77bf|\u67f4|\u67cf|\u913a|\u8ac7|\u67e5|\u970d|\u968b|\u9594|\u9ad9|\u7ac7|\u677e|\u5409|\u752f|\u9072|\u5132|\u98a8|\u91cb|\u4ef2|\u5189|\u9102|\u6e5b|\u4ec7|\u6771|\u5321|\u69ae|\u4f0a|\u660c|\u5a41|\u862d|\u51b7|\u535e|\u6851|\u664f|\u5c91|\u88d8|\u59ec|\u5e2d|\u8499|\u521d|\u95bb|\u90c1|\u7c73|\u53e2|\u5055|\u660e|\u5175|\u7504|\u52de|\u97a0|\u8096|\u834a|\u666f|\u805e|\u90ce|\u595a|\u4f5f|\u8306|\u9854|\u53b2|\u5c60|\u76e4|\u6a5f|\u624d|\u7c9f|\u5c01|\u7987|\u539f|\u5e72|\u9321|\u8c9d|\u5e73|\u5b9c|\u5bb9|\u6a19|\u51bc|\u76ae|\u5ba3|\u7afa|\u84cb|\u8305|\u6556|\u9122|\u85fa|\u5bc7|\u9ee8|\u82ae|\u5357|\u6822|\u5371|\u72c4|\u70cf|\u60e0|\u79b9|\u678b|\u8fb2|\u82d1|\u675e|\u80e5|\u5ca9|\u6d66|\u5d47|\u8af8|\u9004|\u695a|\u5e2b|\u57ce|\u6b09|\u4fee|\u6708|\u620e|\u8c50|\u9e7f|\u54c0|\u5ffb|\u5f37|\u6eff|\u5019|\u5be7|\u68ee|\u5143][\u4e00-\u9fa5]{1,2}"

	re := regexp.MustCompile(`(?s)<meta.*?</head>`) //過濾網頁header
	an := regexp.MustCompile(`<.*?>`)
	regex := regexp.MustCompile(pattern)

	match := re.ReplaceAllString(bodystr, " ")
	match = an.ReplaceAllString(bodystr, " ")

	//fmt.Println(match)

	match_chinese := regex.FindAllString(match, -1)

	fileObj, err := os.Open("./new_blacklist.txt") //檔案開啟
	if err != nil {
		fmt.Println("Error Open File")
		//return
	}

	defer fileObj.Close() //檔案關閉

	scanner := bufio.NewScanner(fileObj) //blacklist檔案讀取
	for scanner.Scan() {
		b_list = append(b_list, scanner.Text())
	}

	for _, v := range match_chinese {
		//fmt.Println(v[0])
		var AllowEntry bool = true

		for _, blist := range b_list {

			find := strings.Contains(v, blist) //判斷是否存在blacklist

			if find {

				AllowEntry = false //存在黑名單 break
				break
			}

		}

		if AllowEntry {

			check := nlp_name(v)
			if check {
				han_chk := hanlp_api(v)
				if han_chk {

					chk_result := in(v, username) //判斷是否存在array
					if !chk_result {
						username = append(username, v)
					}
				}

			}

		}

	}

	//身分證比對
	re_id := regexp.MustCompile(`[a-zA-z][12][0-9]{8}`)
	match_id := re_id.FindAllString(match, -1)

	for _, id := range match_id {
		id_name = append(id_name, id)
	}

	//比對電話
	re_phone := regexp.MustCompile(`(?:0|886-?)9\d{2}(?:\d{6}|-\d{6}|-\d{3}-\d{3}?)`)
	re_hms := regexp.MustCompile(`(\d{3}-\d{8}|0[24]-\d{8}|0[24]-\d{4}-\d{4}|0[3-8]-\d{7}|0[3-8]-\d{3}-\d{4}|(?:037|082|089)-\d{6}|(?:037|082|089)-\d{3}-\d{3}|049-\d{7}|049-\d{3}-\d{4}|0836-\d{5})`)
	match_phone := re_phone.FindAllString(match, -1)
	match_hms := re_hms.FindAllString(match, -1)

	for _, p := range match_phone {
		phone = append(phone, p)
	}

	for _, hm := range match_hms {
		phone = append(phone, hm)
	}

	//比對信箱
	re_mail := regexp.MustCompile(`([a-zA-Z0-9_+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+)`)
	match_mail := re_mail.FindAllString(match, -1)

	for _, mails := range match_mail {
		mail = append(mail, mails)
	}

	//比對地址
	adrpattern := "(?:\u57fa\u9686\u5e02|\u81fa\u5317\u5e02|\u53f0\u5317\u5e02|\u5317\u5e02|\u65b0\u5317\u5e02|\u6843\u5712\u5e02|\u65b0\u7af9\u5e02|\u7af9\u5e02|\u65b0\u7af9\u7e23|\u7af9\u7e23|\u82d7\u6817\u7e23|\u81fa\u4e2d\u5e02|\u53f0\u4e2d\u5e02|\u4e2d\u5e02|\u5f70\u5316\u7e23|\u5357\u6295\u7e23|\u96f2\u6797\u7e23|\u5609\u7fa9\u5e02|\u5609\u7fa9\u7e23|\u81fa\u5357\u5e02|\u53f0\u5357\u5e02|\u9ad8\u96c4\u5e02|\u5c4f\u6771\u7e23|\u81fa\u6771\u7e23|\u53f0\u6771\u7e23|\u82b1\u84ee\u7e23|\u5b9c\u862d\u7e23|\u6f8e\u6e56\u7e23|\u91d1\u9580\u7e23|\u9023\u6c5f\u7e23).{3,20}[\u865f]"
	re_addr := regexp.MustCompile(adrpattern)
	match_addr := re_addr.FindAllString(match, -1)
	var tmp_addr string
	for _, addr := range match_addr {

		if strings.Contains(addr, "63799") {
			tmp_addr = strings.Replace(addr, "&#63799;", "路", -1)
		}
		address = append(address, tmp_addr)
	}

	//比對生日
	re_vids := regexp.MustCompile(`((?:^19\d{2}|20[01]\d|202[0-2])(?:年|/|-|)(?:0[1-9]|1[0-2]|[1-9])(?:月|/|-|)(?:0\d|3[0-1]|[1-9]|[1-2]\d))`)
	re_roc := regexp.MustCompile(`((?:^0[2-9]\d|10\d|110|^[2-9][\d])(?:年|/|-|)(?:0\d|1[0-2]|[1-9])(?:月|/|-|)(?:0\d|3[0-1]|[1-9]|[1-2]\d)$)`)

	match_vids := re_vids.FindAllString(match, -1)
	match_roc := re_roc.FindAllString(match, -1)

	for _, vid := range match_vids {

		vaildate = append(vaildate, vid)
	}
	for _, roc := range match_roc {

		vaildate = append(vaildate, roc)
	}

	m["username"] = username
	m["id"] = id_name
	m["ph_no"] = phone
	m["mail"] = mail
	m["addr"] = address
	m["vid"] = vaildate

	return m
}

func hanlp_api(v string) bool {

	HanLP := hanlp.HanLPClient(hanlp.WithAuth(""), hanlp.WithLanguage("zh"), hanlp.WithTasks("ner/pku"))
	s, _ := HanLP.Parse(v,
		hanlp.WithLanguage("zh"))
	//fmt.Println(s)
	check := strings.Contains(s, "nr")
	if check {
		return true
	} else {
		return false
	}

}

func nlp_name(v string) bool {

	lac := thulac.NewThulacgo("models", "", false, false, false, byte('_'))
	defer lac.Deinit()

	ret := lac.Seg(v)
	check := strings.Contains(ret, "np")
	if check {
		return true
	} else {
		return false
	}

}

func in(target string, str_array []string) bool {
	sort.Strings(str_array)
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}

func Unicodelen(str string) int {
	var r = []rune(str)
	return len(r)
}

func UnicodeIndex(str, substr string) int {
	result := strings.Index(str, substr)
	if result > 0 {
		prefix := []byte(str)[0:result]
		rs := []rune(string(prefix))
		result = len(rs)
	}
	return result
}

func SubString(source string, start int, end int) string {
	var unicodestr = []rune(source)
	length := len(unicodestr)
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	if start <= 0 && end >= length {
		return source
	}
	var substring = ""
	for i := start; i < end; i++ {
		substring += string(unicodestr[i])
	}
	return substring
}

func txtReader(filepath string) {
	start1 := time.Now()
	FileHandle, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		return
	}
	defer FileHandle.Close()
	lineReader := bufio.NewReader(FileHandle)

	for {
		line, _, err := lineReader.ReadLine()
		if err == io.EOF {
			break
		}
		w_list = append(w_list, string(line))
	}

	fmt.Println("檔案載入Spend : ", time.Now().Sub(start1))
}
