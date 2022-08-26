package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const cookie = "first_visit_datetime_pc=2022-01-25+15%3A13%3A33; p_ab_id=3; p_ab_id_2=2; p_ab_d_id=840906738; yuid_b=JIaGdGc; _fbp=fb.1.1647667910426.139149109; a_type=0; b_type=1; login_ever=yes; adr_id=TZY97OfgrdMWC6QlIvF3gES1lhFPJBrdzPCo7lOxbqicHpUN; c_type=19; __utmv=235335808.|2=login%20ever=yes=1^3=plan=normal=1^5=gender=male=1^6=user_id=79555544=1^9=p_ab_id=3=1^10=p_ab_id_2=2=1^11=lang=zh=1; cto_bundle=CNGuS19jTGNnbGFmelhOTVBMeUElMkZhbnFpOVBOeUMzciUyQm9RdlhnTjA4UURhbElvZVRXRW0wbE5xb0F6NmJSSEljNjZOT0tPMzJMdjFNRFQ4QmQ4UEtiUHJKS2hQVWxFeiUyQmRWSUMwelglMkJlWUR5dG9IdjZCdUh5byUyRkslMkZTZDNuWjhzOWNSWWdvdk1VcGhhMEc0QXlzZjNVWnloOVElM0QlM0Q; __utmc=235335808; _im_vid=01GA9CSP6P42HMGFHMKC16AZ3D; _im_uid.3929=i.pt-aBQfeTAeFsCPEPrSsPQ; QSI_S_ZN_5hF4My7Ad6VNNAi=v:0:0; __utmz=235335808.1660360088.4.3.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); tag_view_ranking=0Sds1vVNKR~i8u6Dgt7ao~Hry6GxyqEm~0xsDLqCEW6~_EOd7bsGyl~ziiAzr_h04~TqiZfKmSCg~DjBjqF155l~oCqKGRNl20~sfVA-LlZEm~eVxus64GZU~qiO14cZMBI~RVRPe90CVr~oCR2Pbz1ly~FySY6ZVB78~_giyO1uU9O~7y6oYXppQq~x_jB0UM4fe~2-ZLcTJsOe~pzzjRSV6ZO~TOd0tpUry5~5qLN1_M1lU~S9rSdAbx_2~oPZuHRRfMy~8oxa62fI2y~RRr-UOsU3F~RlNSEdcf6r~cneBYSwz3B~LxsZ2n85kN~8c5TFFwEPe~6o-mWtVX-B~Ekc63fSc_I~-wt3kJbLCW~x5q_KUbrlK~gpglyfLkWs~AteUadc2QM~63aLATR0Cy~YCF2322lfw~vH5NhSLEru~yOqOtdektt~1hJHysZey1~yJInVP0PK6~hcEdDGw-eP~xufWQ15ZA3; __utma=235335808.1969022270.1657207716.1660360088.1660529111.5; _gid=GA1.2.793944985.1660529115; __utmb=235335808.9.10.1660529111; _ga_75BBYNYN9J=GS1.1.1660529112.5.1.1660531360.0; _ga=GA1.2.1969022270.1657207716"

var urlID []string

func fetchPage(num int) {
	orginURL := "https://www.pixiv.net/ajax/search/artworks/%E5%8F%A4%E6%98%8E%E5%9C%B0%E3%81%93%E3%81%84%E3%81%97?word=%E5%8F%A4%E6%98%8E%E5%9C%B0%E3%81%93%E3%81%84%E3%81%97&order=date_d&mode=all&p=" + strconv.Itoa(num) + "&s_mode=s_tag_full&type=all&lang=zh"
	proxy, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	req, _ := http.NewRequest("GET", orginURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Http get err:", err)
		os.Exit(1)
	}
	if resp.StatusCode != 200 {
		fmt.Println("Http status code:", resp.StatusCode)
		os.Exit(1)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	value := gjson.Get(string(body), "body.illustManga.data")
	value.ForEach(func(_, value gjson.Result) bool {
		valueID := value.Get("id")
		if valueID.String() != "" {
			urlID = append(urlID, valueID.String())
		}
		return true
	})
}

//*[@id="root"]/div[2]/div/div[3]/div/div[6]/div/section[2]/div[2]/ul/li[1]/div/div[1]/div/a/div[1]/img
func GetImageURL(id string) string {
	imageURL := "https://www.pixiv.net/artworks/" + id
	fmt.Printf(imageURL)
	proxy, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	req, _ := http.NewRequest("GET", imageURL, nil)
	req.Header.Set("Referer", "www.pixiv.net")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("爬取image失败")
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	reg := regexp.MustCompile(`(?:"original":")(.*?)(?:")\}`)
	str := reg.FindStringSubmatch(string(body))
	if strings.HasPrefix(str[1], "http") {
		fmt.Println("成功抓取图片链接")
		return str[1]
	}
	return ""
}
func DownLoadImage(val string, ch chan bool) {
	imageURL := GetImageURL(val)
	fmt.Println(imageURL)
	proxy, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxy),
		},
	}
	req, _ := http.NewRequest("GET", imageURL, nil)
	req.Header.Set("Referer", "http://pixiv.net")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36")
	req.Header.Add("Cookie", cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("爬取image失败")
		return
	}
	fmt.Println("爬取image成功")
	defer resp.Body.Close()
	downloadPath := fmt.Sprintf("image/%s.png", val)
	body, err := ioutil.ReadAll(resp.Body)
	file, err := os.OpenFile(downloadPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("打开文件失败", err)
		return
	}
	defer file.Close()
	file.Write(body)
	if err != nil {
		fmt.Println("写入文件失败")
		return
	}
	fmt.Println("写入文件成功")
	ch <- true
}
func main() {
	start := time.Now()
	fetchPage(1)
	ch := make(chan bool)
	for _, val := range urlID {
		go DownLoadImage(val, ch)
	}
	for range urlID {
		<-ch
	}
	//wg.Wait()
	end := time.Since(start)
	fmt.Println("耗时时间", end, "爬取图片", len(urlID))
}
