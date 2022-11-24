package main

import ( //"log"
	//log "github.com/Sirupsen/logrus"
	//BEST
	//"github.com/rsc/pdf"
	//_ "github.com/PuerkitoBio/goquery"
	//"github.com/dslipak/pdf"
	//_ "github.com/lib/pq"
	//sss "github.com/compscidr/go-scholar/scholar"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	//"github.com/gocolly/colly" /////////////////////DAS IST DAS RICHGIGE MERKEN
	"github.com/gocolly/colly/v2"
	//sch "github.com/sotetsuk/goscholar"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type Paper struct {
	Name     string   `json:"name"`
	Year     int      `json:"year"`
	Yearinfo string   `json:"yearinfo"`
	Writer   []string `json:"writer"`
	Cites_Nr int      `json:"cites_nr"`
	Cites    []Paper  `json:"cites"`
	Url      string   `json:"url"`
	//Topics   []string `json:"topics"`
	File     string `json:"file"`
	File_url string `json:"file_url"`
}

func main() {
	//https://scholar.google.de/scholar?start=00&q=%22organizational+routines%22+hicss&hl=de&as_sdt=0,5
	//arr := search_papers("https://scholar.google.de/scholar?start=0&q=%22organizational+routines%22&hl=de&as_sdt=8,5")

	//var papers []Paper
	var papersnew []Paper

	//Load txt reference data (Paper)
	file, _ := ioutil.ReadFile("cited_papersnew2.json")
	//data := PaperNodes{}
	var data []Paper
	_ = json.Unmarshal([]byte(file), &data)

	papers := data
	// ki := 0
	// //for i := 0; i < len(papers[:5]); i++ {
	// for _, element := range papers[828:] {
	// 	// for i := 0; i < len(papers); i++ {
	// 	// 	if len(papers[i].Url) > 0 && papers[i].Cites_Nr > 0 {
	// 	if len(element.File_url) > 0 {
	// 		ki = ki + 1

	// 		fmt.Print("Current Count: ")
	// 		fmt.Println(ki)
	// 		DownloadFile(element.Name+".pdf", element.File_url)
	// 	}

	// 	// 		decodedValue, err := url.QueryUnescape(papers[i].Url)
	// 	// 		fmt.Println("URL" + papers[i].Url)
	// 	// 		fmt.Println("DECODEURL" + decodedValue)
	// 	// 		if err != nil {
	// 	// 			log.Fatal(err)
	// 	// 			return
	// }

	//cites :=
	//FAILT vermutlich wegen sehr häufiger aufrufe und neue visits in pages, oder in allowed domains, oder versuchen immer wieder ip zu wechseln;

	//papers[i].Cites = search_papers3(papers)
	//	}

	papersnew = search_papers3(papers[280:300])
	// append(papersnew, papers[i])
	// fmt.Println("URL:"arr[i].Url)
	// fmt.Println(arr[i].Cites_Nr)

	//}
	//}
	js, err := json.MarshalIndent(papersnew, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("cites_final15.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
	}
	fmt.Println(len(papersnew))

	//err mit returnen später
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create("pdf/" + filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func search_papers(path string) []Paper {
	//colly.Async(true)
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)
	cited_papers := make([]Paper, 0)
	var p []Paper

	var titles []string
	titles_nr := 0

	c.OnHTML("div#gs_res_ccl_mid", func(e *colly.HTMLElement) {
		e.ForEach("div.gs_r.gs_or.gs_scl", func(i int, h *colly.HTMLElement) {
			var title string
			fmt.Println("HI")
			item := Paper{}
			h.ForEach("h3.gs_rt", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					if contains(titles, n.Text) {
						//unkorrekter counter aber egal, da hauptsache utnerscheidbar
						titles_nr = titles_nr + 1
						title = strings.Replace(n.Text, ":", " ", -1) + string(titles_nr)
						item.Name = title
					} else {
						title = strings.Replace(n.Text, ":", " ", -1)
						titles = append(titles, title)
						item.Name = title
					}

				})
			})
			h.ForEach("div.gs_ggs.gs_fl", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {
					fmt.Println("HI")

					str := n.Attr("href")
					//DownloadFile(title+".pdf", str)
					item.File = title + ".pdf"
					item.File_url = str

				})
				//c.OnHTML("div#gs_bdy_ccl", func(e *colly.HTMLElement) {
			})
			h.ForEach("div.gs_a", func(i int, u *colly.HTMLElement) {
				//sr := u.Text
				item.Yearinfo = u.Text

				str1 := u.Text
				//extract year
				re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

				//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern

				//fmt.Printf("String contains any match: %v\n", re.MatchString(str1)) // True

				submatchall := re.FindAllString(str1, -1)
				for _, element := range submatchall {
					intVar, err := strconv.Atoi(element)
					if err != nil {
						item.Year = 0
					} else {
						item.Year = intVar
					}

				}

				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					item.Writer = append(item.Writer, n.Text)

				})
				if len(item.Writer) == 0 {
					//extract writer
					fmt.Println("writerTEXT " + u.Text)

					writer_bef, _, status := strings.Cut(u.Text, "-")
					if status {
						fmt.Println("writerBEFORE " + writer_bef)
						item.Writer = append(item.Writer, writer_bef)
					}
					//ALTERNATIVE METHODE BIS JAHRESZAHL ////////////////////
					// submatchall := re.FindAllString(u.Text, 1)
					// for _, element := range submatchall {
					// 	if len(element) > 0 {
					// 		item.Writer = append(item.Writer, element)
					// 	}
					// }

				}
			})
			h.ForEach("div.gs_fl", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {
					// 		item.Cites_Nr = 6
					if strings.Contains(n.Text, "Zitiert von:") {

						str := "https://scholar.google.de" + n.Attr("href")

						item.Url = str
						obj := n.Text[13:]
						intVarr, err := strconv.Atoi(obj)
						fmt.Println(obj)
						fmt.Println(intVarr)
						if err != nil {
							item.Cites_Nr = 0
						} else {
							item.Cites_Nr = intVarr
						}

						//item.Cites = search_citations(str, intVarr)

						// 			// item.Image = n.ChildAttr("img", "data-src")
						// 			// item.Price = n.Attr("data-price")
						// 			// item.Url = "https://jumia.com.ng" + e.Attr("href")
						// 			// item.Discount = e.ChildText("div.tag._dsct")
					}
				})
			})

			cited_papers = append(cited_papers, item)
			p = append(p, item)
		})

	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL)
	})

	// rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.SetProxyFunc(rp)

	// // Print the response
	// c.OnResponse(func(r *colly.Response) {
	// 	log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
	// 	log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	// })

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(cited_papers, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("cited_papers_singlenew.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		Delay:       10 * time.Second,
	})

	//c.Visit(path)
	fullURL := ""
	full := ""

	for i := 0; i <= 97; i++ {
		randombool := RandBool()
		if randombool {
			full = "https://de.wikipedia.org/wiki/Fu%C3%9Fball-Europameisterschaft_2021"
			i = i - 1 // since page was not visited which we want to visits
		} else {
			if i == 0 {
				full = path
			} else {
				fullURL = fmt.Sprintf("https://scholar.google.de/scholar?start=%d0", i)
				addit := "&q=%22organizational+routines%22+hicss&hl=de&as_sdt=0,5"
				full = fullURL + addit
			}
		}

		c.Visit(full)
	}
	//c.Visit(path)
	// c.Wait()

	return cited_papers
}

func RandBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 1
}

func search_papers2(path string, nr int) []Paper {
	//colly.Async(true)
	pages := int(math.Ceil(float64(nr) / 10.0))
	//starttime := time.Now()
	c := colly.NewCollector()

	c.SetRequestTimeout(120 * time.Second) //hier hoeher oder niedriger
	cited_papers := make([]Paper, 0)
	var p []Paper

	var titles []string
	titles_nr := 0

	c.OnHTML("div#gs_res_ccl_mid", func(e *colly.HTMLElement) {
		e.ForEach("div.gs_r.gs_or.gs_scl", func(i int, h *colly.HTMLElement) {
			var title string
			fmt.Println("HI")
			item := Paper{}
			h.ForEach("h3.gs_rt", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					if contains(titles, n.Text) {
						//unkorrekter counter aber egal, da hauptsache utnerscheidbar
						titles_nr = titles_nr + 1
						title = strings.Replace(n.Text, ":", " ", -1) + string(titles_nr)
						item.Name = title
					} else {
						title = strings.Replace(n.Text, ":", " ", -1)
						titles = append(titles, title)
						item.Name = title
					}

				})
			})
			h.ForEach("div.gs_ggs.gs_fl", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {
					fmt.Println("HI")

					str := n.Attr("href")
					//DownloadFile(title+".pdf", str)
					item.File = title + ".pdf"
					item.File_url = str

				})
				//c.OnHTML("div#gs_bdy_ccl", func(e *colly.HTMLElement) {
			})
			h.ForEach("div.gs_a", func(i int, u *colly.HTMLElement) {
				//sr := u.Text
				item.Yearinfo = u.Text

				str1 := u.Text
				//extract year
				re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

				//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern

				//fmt.Printf("String contains any match: %v\n", re.MatchString(str1)) // True

				submatchall := re.FindAllString(str1, -1)
				for _, element := range submatchall {
					intVar, err := strconv.Atoi(element)
					if err != nil {
						item.Year = 0
					} else {
						item.Year = intVar
					}

				}

				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					item.Writer = append(item.Writer, n.Text)

				})
				if len(item.Writer) == 0 {
					//extract writer
					fmt.Println("writerTEXT " + u.Text)

					writer_bef, _, status := strings.Cut(u.Text, "-")
					if status {
						fmt.Println("writerBEFORE " + writer_bef)
						item.Writer = append(item.Writer, writer_bef)
					}
					//ALTERNATIVE METHODE BIS JAHRESZAHL ////////////////////
					// submatchall := re.FindAllString(u.Text, 1)
					// for _, element := range submatchall {
					// 	if len(element) > 0 {
					// 		item.Writer = append(item.Writer, element)
					// 	}
					// }

				}
			})
			h.ForEach("div.gs_fl", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {
					// 		item.Cites_Nr = 6
					if strings.Contains(n.Text, "Zitiert von:") {

						str := "https://scholar.google.de" + n.Attr("href")

						item.Url = str
						obj := n.Text[13:]
						intVarr, err := strconv.Atoi(obj)
						fmt.Println(obj)
						fmt.Println(intVarr)
						if err != nil {
							item.Cites_Nr = 0
						} else {
							item.Cites_Nr = intVarr
						}

						//item.Cites = search_citations(str, intVarr)

						// 			// item.Image = n.ChildAttr("img", "data-src")
						// 			// item.Price = n.Attr("data-price")
						// 			// item.Url = "https://jumia.com.ng" + e.Attr("href")
						// 			// item.Discount = e.ChildText("div.tag._dsct")
					}
				})
			})

			cited_papers = append(cited_papers, item)
			p = append(p, item)
		})

	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL)

		// if time.Since(starttime) > 180 { //* time.Second
		// 	r.Abort()
		// }
	})

	// rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.SetProxyFunc(rp)

	// // Print the response
	// c.OnResponse(func(r *colly.Response) {
	// 	log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
	// 	log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	// })

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)

	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(cited_papers, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("cited_papersnew3.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		Delay:       20 * time.Second,
	})

	full := ""
	betw := ""
	addit := ""
	url1 := "https://scholar.google.de/scholar?start="
	//und := "&"
	url2 := ""
	///////////////////////////HIER
	//link := get_link(path)
	//if len(link) > 0 {
	if pages > 20 {
		pages = 20
	}
	//if len(path) > 0 {
	url2 = path[34:]
	for i := 0; i <= pages; i++ {
		randombool := RandBool()
		if randombool { /////!Even(i)
			full = "https://de.wikipedia.org/wiki/Wikipedia:Hauptseite"
			i = i - 1
		} else {
			if i == 0 {
				full = path
			} else {
				//fullURL = url1
				betw = fmt.Sprint(i) + "0"
				addit = url2
				full = url1 + betw + "&" + addit
				//https://scholar.google.de/scholar?cites=9875713578711629068&as_sdt=2005&sciodt=0,5&hl=de

				//https://scholar.google.de/scholar?start=&cites=9875713578711629068&as_sdt=2005&sciodt=0,5&hl=de

				//https://scholar.google.de/scholar?cites=15502415196660838898&as_sdt=2005&sciodt=0,5&hl=de
				//https://scholar.google.de/scholar?start=0&hl=de&as_sdt=2005&sciodt=0,5&cites=15502415196660838898&scipsc=
				fmt.Println("ACHTUNG!!!!!!" + full)
			}
		}

		c.Visit(full)
	}
	// } else {
	// 	c.Visit(path)
	// }

	//c.Visit(path)
	// c.Wait()

	return cited_papers
}

func search_papers3(paperlist []Paper) []Paper {
	//colly.Async(true)

	//starttime := time.Now()
	c := colly.NewCollector()
	c.SetRequestTimeout(240 * time.Second) //hier hoeher oder niedriger

	cited_papers := make([]Paper, 0)

	ci := 0
	var p []Paper

	var titles []string
	titles_nr := 0
	//act_paper := paperlist[ci]

	c.OnHTML("div#gs_res_ccl_mid", func(e *colly.HTMLElement) {
		e.ForEach("div.gs_r.gs_or.gs_scl", func(i int, h *colly.HTMLElement) {
			var title string
			fmt.Println("HI")
			item := Paper{}
			h.ForEach("h3.gs_rt", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					if contains(titles, n.Text) {
						//unkorrekter counter aber egal, da hauptsache utnerscheidbar
						titles_nr = titles_nr + 1
						title = strings.Replace(n.Text, ":", " ", -1) + string(titles_nr)
						item.Name = title
					} else {
						title = strings.Replace(n.Text, ":", " ", -1)
						titles = append(titles, title)
						item.Name = title
					}

				})
			})
			h.ForEach("div.gs_ggs.gs_fl", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {
					fmt.Println("HI")

					str := n.Attr("href")
					//DownloadFile(title+".pdf", str)
					item.File = title + ".pdf"
					item.File_url = str

				})
				//c.OnHTML("div#gs_bdy_ccl", func(e *colly.HTMLElement) {
			})
			h.ForEach("div.gs_a", func(i int, u *colly.HTMLElement) {
				//sr := u.Text
				item.Yearinfo = u.Text

				str1 := u.Text
				//extract year
				re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

				//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern

				//fmt.Printf("String contains any match: %v\n", re.MatchString(str1)) // True

				submatchall := re.FindAllString(str1, -1)
				for _, element := range submatchall {
					intVar, err := strconv.Atoi(element)
					if err != nil {
						item.Year = 0
					} else {
						item.Year = intVar
					}

				}

				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					item.Writer = append(item.Writer, n.Text)

				})
				if len(item.Writer) == 0 {
					//extract writer
					fmt.Println("writerTEXT " + u.Text)

					writer_bef, _, status := strings.Cut(u.Text, "-")
					if status {
						fmt.Println("writerBEFORE " + writer_bef)
						item.Writer = append(item.Writer, writer_bef)
					}
					//ALTERNATIVE METHODE BIS JAHRESZAHL ////////////////////
					// submatchall := re.FindAllString(u.Text, 1)
					// for _, element := range submatchall {
					// 	if len(element) > 0 {
					// 		item.Writer = append(item.Writer, element)
					// 	}
					// }

				}
			})
			h.ForEach("div.gs_fl", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {
					// 		item.Cites_Nr = 6
					if strings.Contains(n.Text, "Zitiert von:") {

						str := "https://scholar.google.de" + n.Attr("href")

						item.Url = str
						obj := n.Text[13:]
						intVarr, err := strconv.Atoi(obj)
						fmt.Println(obj)
						fmt.Println(intVarr)
						if err != nil {
							item.Cites_Nr = 0
						} else {
							item.Cites_Nr = intVarr
						}

						//item.Cites = search_citations(str, intVarr)

						// 			// item.Image = n.ChildAttr("img", "data-src")
						// 			// item.Price = n.Attr("data-price")
						// 			// item.Url = "https://jumia.com.ng" + e.Attr("href")
						// 			// item.Discount = e.ChildText("div.tag._dsct")
					}
				})
			})

			cited_papers = append(cited_papers, item)
			p = append(p, item)
		})

	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL)

		// if time.Since(starttime) > 180 { //* time.Second
		// 	r.Abort()
		// }
	})

	// rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:1337", "socks5://127.0.0.1:1338")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.SetProxyFunc(rp)

	// // Print the response
	// c.OnResponse(func(r *colly.Response) {
	// 	log.Printf("Proxy Address: %s\n", r.Request.ProxyURL)
	// 	log.Printf("%s\n", bytes.Replace(r.Body, []byte("\n"), nil, -1))
	// })

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)

	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(cited_papers, "", "    ")
		paperlist[ci].Cites = append(paperlist[ci].Cites, cited_papers...)
		cited_papers = make([]Paper, 0)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("cited_papersnewest.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		Delay:       20 * time.Second,
	})

	for ci = 0; ci < len(paperlist); ci++ {

		// if ci>0 {
		// 	paperlist[ci-1]
		// }

		full := ""
		betw := ""
		addit := ""
		url1 := "https://scholar.google.de/scholar?start="
		//und := "&"
		url2 := ""
		///////////////////////////HIER
		//link := get_link(path)
		//if len(link) > 0 {

		if len(paperlist[ci].Url) > 0 && paperlist[ci].Cites_Nr > 0 {
			pages := int(math.Ceil(float64(paperlist[ci].Cites_Nr) / 10.0))
			url2 = paperlist[ci].Url[34:]
			if pages > 5 {
				pages = 5
			}
			fmt.Println("PAGES")
			fmt.Println(pages)
			for i := 0; i <= pages; i++ {
				randombool := RandBool()
				if randombool { /////!Even(i)
					full = "https://de.wikipedia.org/wiki/Wikipedia:Hauptseite"
					i = i - 1
				} else {
					if i == 0 {
						full = paperlist[ci].Url
					} else {
						//fullURL = url1
						betw = fmt.Sprint(i) + "0"
						addit = url2
						full = url1 + betw + "&" + addit
						//https://scholar.google.de/scholar?cites=9875713578711629068&as_sdt=2005&sciodt=0,5&hl=de

						//https://scholar.google.de/scholar?start=&cites=9875713578711629068&as_sdt=2005&sciodt=0,5&hl=de

						//https://scholar.google.de/scholar?cites=15502415196660838898&as_sdt=2005&sciodt=0,5&hl=de
						//https://scholar.google.de/scholar?start=0&hl=de&as_sdt=2005&sciodt=0,5&cites=15502415196660838898&scipsc=
						fmt.Println("ACHTUNG!!!!!!" + full)
					}
				}

				c.Visit(full)
			}
		} //else {
		// 	full = full = "https://de.wikipedia.org/wiki/Wikipedia:Hauptseite"
		// }
	}
	// } else {
	// 	c.Visit(path)
	// }

	//c.Visit(path)
	// c.Wait()

	return paperlist
}

func Even(number int) bool {
	return number%2 == 0
}

func search_citations(path string, nr int) []Paper {
	//var cited_papers = make(map[int]Paper)   , nr int
	//fmt.Println("ok")
	//colly.Async(true)
	pages := nr / 10
	fmt.Println(pages)
	c := colly.NewCollector()
	c.SetRequestTimeout(120 * time.Second)
	cited_papers := make([]Paper, 0)
	var titles []string
	titles_nr := 0

	// Callbacks

	c.OnHTML("div#gs_res_ccl_mid", func(e *colly.HTMLElement) {
		e.ForEach("div.gs_r.gs_or.gs_scl", func(i int, h *colly.HTMLElement) {
			var title string

			item := Paper{}
			h.ForEach("h3.gs_rt", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					if contains(titles, n.Text) {
						//unkorrekter counter aber egal, da hauptsache utnerscheidbar
						titles_nr = titles_nr + 1
						title = strings.Replace(n.Text, ":", " ", -1) + fmt.Sprint(titles_nr)
						item.Name = title
					} else {
						title = strings.Replace(n.Text, ":", " ", -1)
						item.Name = title
						titles = append(titles, title)
					}

				})
			})
			h.ForEach("div.gs_a", func(i int, u *colly.HTMLElement) {
				//sr := u.Text
				item.Yearinfo = u.Text

				str1 := u.Text
				//extract year
				re := regexp.MustCompile(`[-]?\d[\d,]*[\.]?[\d{2}]*`)

				//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern

				//fmt.Printf("String contains any match: %v\n", re.MatchString(str1)) // True

				submatchall := re.FindAllString(str1, -1)
				for _, element := range submatchall {
					intVar, err := strconv.Atoi(element)
					if err != nil {
						item.Year = 0
					} else {
						item.Year = intVar
					}

				}

				u.ForEach("a", func(i int, n *colly.HTMLElement) {

					item.Writer = append(item.Writer, n.Text)

				})

				if len(item.Writer) == 0 {
					//extract writer
					fmt.Println("writerTEXT " + u.Text)

					writer_bef, _, status := strings.Cut(u.Text, "-")
					if status {
						fmt.Println("writerBEFORE " + writer_bef)
						item.Writer = append(item.Writer, writer_bef)
					}
					//ALTERNATIVE METHODE MIT STRING UTEXT BIS JAHRESZAHL GEFUNDEN;
					// submatchall := re.FindAllString(u.Text, 1)
					// for _, element := range submatchall {
					// 	if len(element) > 0 {
					// 		item.Writer = append(item.Writer, element)
					// 	}
					// }

				}
			})
			h.ForEach("div.gs_fl", func(i int, u *colly.HTMLElement) {
				u.ForEach("a", func(i int, n *colly.HTMLElement) {
					// 		item.Cites_Nr = 6
					if strings.Contains(n.Text, "Zitiert von:") {

						//str := "https://scholar.google.de" + n.Attr("href")

						//item.Url = "https://scholar.google.de" + n.Attr("href")

						obj := n.Text[13:]
						intVarr, err := strconv.Atoi(obj)
						fmt.Println(obj)
						fmt.Println(intVarr)
						if err != nil {
							item.Cites_Nr = 0
						} else {
							item.Cites_Nr = intVarr
						}

						//item.Cites = search_citations(str)

						// 			// item.Image = n.ChildAttr("img", "data-src")
						// 			// item.Price = n.Attr("data-price")
						// 			// item.Url = "https://jumia.com.ng" + e.Attr("href")
						// 			// item.Discount = e.ChildText("div.tag._dsct")
					}
				})
			})

			cited_papers = append(cited_papers, item)
		})

	})

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		Delay:       60 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
		fmt.Println("Visiting", r.URL)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Got a response from", r.Request.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Println("Got this error:", e)
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
		js, err := json.MarshalIndent(cited_papers, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Writing data to file")
		if err := os.WriteFile("cited_papersnew.json", js, 0664); err == nil {
			fmt.Println("Data written to file successfully")
		}

	})

	//c.Visit(path)
	///////////////////////////HIER
	//fullURL := ""
	full := ""
	betw := ""
	addit := ""
	url1 := "https://scholar.google.de/scholar?start="
	//und := "&"
	url2 := ""
	///////////////////////////HIER
	//link := get_link(path)
	//if len(link) > 0 {
	url2 = path[34:]
	for i := 0; i <= pages; i++ {
		if i == 0 {
			full = path
		} else {
			//fullURL = url1
			betw = fmt.Sprint(i) + "0"
			addit = url2
			full = url1 + betw + "&" + addit
			//https://scholar.google.de/scholar?cites=9875713578711629068&as_sdt=2005&sciodt=0,5&hl=de

			//https://scholar.google.de/scholar?start=&cites=9875713578711629068&as_sdt=2005&sciodt=0,5&hl=de

			//https://scholar.google.de/scholar?cites=15502415196660838898&as_sdt=2005&sciodt=0,5&hl=de
			//https://scholar.google.de/scholar?start=0&hl=de&as_sdt=2005&sciodt=0,5&cites=15502415196660838898&scipsc=
			fmt.Println("ACHTUNG!!!!!!" + full)
		}

		c.Visit(full)
	}
	// fullURL := ""
	// full := ""

	// for i := 0; i <= 20; i++ {
	// 	if i == 0 {
	// 		full = path
	// 	} else {
	// 		fullURL = fmt.Sprintf("https://scholar.google.de/scholar?start=%d0", i)
	// 		addit := "&q=%22organizational+routines%22+hicss&hl=de&as_sdt=0,5"
	// 		full = fullURL + addit
	// 	}

	// 	c.Visit(full)
	// }
	// } else {
	// 	c.Visit(path)
	// }
	//////////////////////////////HIER
	// for i := 0; i <= 2; i++ {
	// 	if i == 0 {
	// 		full = path
	// 	} else {
	// 		fullURL = url1
	// 		betw = fmt.Sprint(i) + "0"
	// 		addit = url2
	// 		full = fullURL + betw + addit
	// 		//https://scholar.google.de/scholar?cites=15502415196660838898&as_sdt=2005&sciodt=0,5&hl=de
	// 		//https://scholar.google.de/scholar?start=0&hl=de&as_sdt=2005&sciodt=0,5&cites=15502415196660838898&scipsc=
	// 		fmt.Println("ACHTUNG!!!!!!!!!!!!!!!!" + full)
	// 	}

	// 	c.Visit(full)
	// }

	return cited_papers

}
