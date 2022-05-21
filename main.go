package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var (
	input    string
	keyword  string = `+("證號"+|+"姓名"+|+"生日"+|+"出生"+|+"電話"+|+"手機"+|+"護照"+|+"婚姻"+|+"家庭"+|+"教育"+|+"職業"+|+"病歷"+|+"聯絡")`
	ext      string = `+ext:doc+|+ext:docx+|+ext:xls+|+ext:xlsx+|+ext:ppt+|+ext:pptx+|+ext:pdf+|+ext:csv+|+ext:odt+|+ext:rtf+|+ext:sxw+|+ext:csv+|+ext:pps`
	next     string
	FileName string
	n_time   []string
)

func Proxy() {
	//測試用

	p := colly.NewCollector() // 在colly中使用 Collector 這類物件 來做事情

	p.OnResponse(func(r *colly.Response) { // 當Visit訪問網頁後，網頁響應(Response)時候執行的事情
		//fmt.Println(string(r.Body)) // 返回的Response物件r.Body 是[]Byte格式，要再轉成字串
	})

	p.OnRequest(func(r *colly.Request) { // 需要寫這一段 User-Agent才給爬
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	})

	p.OnHTML("#proxy_list", func(e *colly.HTMLElement) {
		dom := e.DOM
		ip_result := dom.Find("td[class=left]").Text()
		re := regexp.MustCompile(`".*?"`)
		match_ip := re.FindAllString(ip_result, -1)

		port_result := dom.Find("td[class=fport]").Text()
		fmt.Println(port_result)

		var i int = 0
		for _, ip := range match_ip {
			split_ip := strings.Trim(ip, "\"")
			fmt.Println(split_ip)
			if i == 5 {
				break
			}
			i++
		}

	})

	p.Visit("http://free-proxy.cz/zh/proxylist/country/all/http/ping/all")
}

func main() {
	n_time = n_time[:0]
	fmt.Print("請輸入搜尋")
	_, err := fmt.Scanln(&input)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	timestr := "Starting at " + time.Now().Format("2006-01-02 15:04:05")
	n_time = append(n_time, timestr)

	// csv 操作
	FileName = string(input + ".csv")
	Pwd, err := os.Getwd()
	FilePath := filepath.Join(Pwd, FileName)

	Person_File := string("PersonInfo_" + input + ".csv")
	Person_FilePath := filepath.Join(Pwd, Person_File)

	//創建csv
	file, err := os.OpenFile(FilePath, os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil && os.IsNotExist(err) {
		file, err := os.Create(FilePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		file.WriteString("\xEF\xBB\xBF")
		w := csv.NewWriter(file) //寫入文件流
		w.Write(n_time)
		title := []string{"No", "Title", "Result", "URL"}
		w.Write(title)
		w.Flush()
	}

	exist_file, err := os.OpenFile(FilePath, os.O_APPEND|os.O_RDWR, 0666)
	defer exist_file.Close()
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(exist_file)

	// csv 個人資料
	pfile, err := os.OpenFile(Person_FilePath, os.O_WRONLY, 0777)
	defer pfile.Close()
	if err != nil && os.IsNotExist(err) {
		pfile, err := os.Create(Person_FilePath)
		if err != nil {
			panic(err)
		}
		defer pfile.Close()
		pfile.WriteString("\xEF\xBB\xBF") //避免中文亂碼
		p_w := csv.NewWriter(pfile)       //寫入文件流
		p_w.Write(n_time)
		title := []string{"No", "Title", "Result", "URL"}
		p_w.Write(title)
		p_w.Flush()
	}
	person_exist_file, err := os.OpenFile(Person_FilePath, os.O_APPEND|os.O_RDWR, 0666)
	defer person_exist_file.Close()
	if err != nil {
		panic(err)
	}
	p_w := csv.NewWriter(person_exist_file)

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {

	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.75 Safari/537.36")
	})
	var num int = 0
	c.OnHTML("#main", func(e *colly.HTMLElement) {

		dom := e.DOM
		result := dom.Find("div[id=result-stats]").Text() //搜尋數量

		metaname := e.DOM.Find(".LC20lb,.MBeuO,.DKV0Md") //標題
		url_num := 0                                     //每頁第幾筆

		fmt.Println("\n\n" + result + "\n\n")
		metaname.Each(func(_ int, s *goquery.Selection) {

			fmt.Println(strconv.Itoa(url_num) + ")----------------" + s.Text() + "----------------") //title
			href := e.DOM.Find("div.znKVS > a").Eq(url_num)
			line, chk := href.Attr("href")
			if chk {
				fmt.Println(line + "\n") //cache href

				res := Sen(line)

				fmt.Printf("疑似洩漏%d姓名 %d身分證 %d電話 %d信箱 %d地址 %d生日\n", len(res["username"]), len(res["id"]), len(res["ph_no"]), len(res["mail"]), len(res["addr"]), len(res["vid"]))
				count := "符合特徵" + strconv.Itoa(len(res["username"])) + "姓名" + strconv.Itoa(len(res["id"])) + "身分證" + strconv.Itoa(len(res["ph_no"])) + "電話" + strconv.Itoa(len(res["mail"])) + "信箱" + strconv.Itoa(len(res["addr"])) + "地址" + strconv.Itoa(len(res["vid"])) + "生日"
				total := []string{strconv.Itoa(num), s.Text(), string(count), string(line)}

				w.Write(total)
				w.Flush()

				person_uname, _ := json.Marshal(res["username"])
				person_uid, _ := json.Marshal(res["id"])
				person_phone, _ := json.Marshal(res["phone"])
				person_mail, _ := json.Marshal(res["mail"])
				person_addr, _ := json.Marshal(res["addr"])
				person_vid, _ := json.Marshal(res["vid"])

				person_str_uname := string(person_uname)
				person_str_uid := string(person_uid)
				person_str_phone := string(person_phone)
				person_str_mail := string(person_mail)
				person_str_addr := string(person_addr)
				person_str_vid := string(person_vid)

				person_info := "姓名" + person_str_uname + "\n" +
					"身分證" + person_str_uid + "\n" +
					"電話" + person_str_phone + "\n" +
					"信箱" + person_str_mail + "\n" +
					"地址" + person_str_addr + "\n" +
					"生日" + person_str_vid + "\n"

				p_total := []string{strconv.Itoa(num), s.Text(), string(person_info), string(line)}

				p_w.Write(p_total)
				p_w.Flush()
				num++
			} else {
				fmt.Println("----Web cache no here!----")
			}

			url_num += 1

		})

	})

	c.OnHTML("#pnnext", func(e *colly.HTMLElement) {
		linkStr := e.Attr("href")
		next = string(linkStr)
		n_page := "https://www.google.com" + next
		e.Request.Visit(n_page)
	})

	c.Visit("https://www.google.com/search?q=site%3A" + input + keyword + ext) // Visit 要放最後
}
