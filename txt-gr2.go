package main

import (
	"flag"
	"fmt"

	//"os"
	"strconv"

	//"log"
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	textrank "github.com/DavidBelicza/TextRank"

	//"fmt"
	//"os"

	//"strings"

	//"golang.org/x/net/html/charset"
	//pdf/vendor/ out of all imports

	//"pdf/vendor/github.com/muesli/clusters"

	"github.com/dominikbraun/graph"
	"github.com/dominikbraun/graph/draw"
	"github.com/go-gota/gota/dataframe"
	"github.com/muesli/clusters"
	"golang.org/x/net/html/charset"

	log "github.com/sirupsen/logrus"

	"github.com/ledongthuc/pdf" //BEST
	//"github.com/rsc/pdf"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"

	//"github.com/go-gota/gota/series"

	_ "github.com/PuerkitoBio/goquery"
	//"github.com/dslipak/pdf"
	"encoding/json"

	_ "github.com/lib/pq"
)

func init() {
	// To get your free API key for metered license, sign up on: https://cloud.unidoc.io
	// Make sure to be using UniPDF v3.19.1 or newer for Metered API key support.
	// content, err := ioutil.ReadFile("C:/Users/neumu/pdffer/pdfs/10302054.pdf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	license.SetMeteredKey("3cb6fa920eff69f1d218f0f0c6ce2fa66f57bbc1abbf9e9b981ac67f153776e8")
	// if err != nil {
	// 	fmt.Printf("ERROR: Failed to set metered key: %v\n", err)
	// 	fmt.Printf("Make sure to get a valid key from https://cloud.unidoc.io\n")
	// 	panic(err)
	// }
}

func outputPdfText(inputPath string, file string) (error, map[int]string) {
	var pagestxt = make(map[int]string)

	//////////////////////////.PDF RAUSCUTTEN      strings.Replace(n.Text, ":", " ", -1)

	v, err := os.Create(file + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer v.Close()

	f, err := os.Open(inputPath)
	if err != nil {
		return err, nil
	}

	defer f.Close()

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return err, nil
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err, nil
	}
	fmt.Println(numPages)

	fmt.Printf("--------------------\n")
	fmt.Printf("PDF to text extraction:\n")
	fmt.Printf("--------------------\n")
	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err, nil
		}

		ex, err := extractor.New(page)
		if err != nil {
			return err, nil
		}

		text, err := ex.ExtractText()
		if err != nil {
			return err, nil
		}

		pagestxt[pageNum] = text
		_, err2 := v.WriteString(text)

		if err2 != nil {
			log.Fatal(err2)
		}
		fmt.Println("------------------------------")
		fmt.Printf("Page %d:\n", pageNum)
		fmt.Printf("\"%s\"\n", text)
		fmt.Println("------------------------------")
	}

	return nil, pagestxt
}

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
)

func readPdf3(path string) (error, map[int]string) {
	var pagestxt = make(map[int]string)
	v, r, err := pdf.Open(path)
	if err != nil {
		fmt.Println("HI")
		return err, nil

	}
	defer v.Close()
	//r, err := pdf.NewReader(v, 100)
	// if err != nil {
	// 	return nil, err
	// }
	totalPage := r.NumPage()
	fmt.Println(string(totalPage))
	fmt.Println("HI")
	//var textBuilder bytes.Buffer
	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}
		str, err := p.GetPlainText(nil)
		//pagestxt[pageIndex] = p.V.String()
		if err != nil {
			return err, nil
		}
		//textBuilder.WriteString(str)
		fmt.Println(str)
		pagestxt[pageIndex] = str
	}
	return err, pagestxt
}

type Reference struct {
	Name     string
	Startind int
	Endind   int
	Year     int
	Nr       string
	Writer   []string
	Cites_Nr int
}

type PaperNodes struct {
	PaperNodes []Paper `json:"paper_nodes"`
}

type PaperNodes2 struct {
	PaperNodes2 []Paper_Scraped `json:"paper_nodes2"`
}

type Paper struct {
	Cited_by string `json:"cited_by"`
	Year     int    `json:"year"`
	Name     string `json:"name"`
	Writer   string `json:"writer"`
	Cites_Nr int    `json:"cites_nr"`
	// Url      string   `json:"url"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}

type Paperv2 struct {
	Cited_by string `json:"cited_by"`
	Year     int    `json:"year"`
	Name     string `json:"name"`
	Writer   string `json:"writer"`
	Cites_Nr int    `json:"cites_nr"`
	Phrases  string `json:"phrases"`
	Tag2     string `json:"tag"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}
type Paperv2pre struct {
	Cited_by    string   `json:"cited_by"`
	Year        int      `json:"year"`
	Name        string   `json:"name"`
	Writer      string   `json:"writer"`
	Cites_Nr    int      `json:"cites_nr"`
	Phrases     string   `json:"phrases"`
	Phrases_Cit []string `json:"phrases_cit"`
	Tag2        string   `json:"tag"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}
type Paperv3pre struct {
	Cited_by    string `json:"cited_by"`
	Year        int    `json:"year"`
	Name        string `json:"name"`
	Writer      string `json:"writer"`
	Cites_Nr    int    `json:"cites_nr"`
	Phrases     string `json:"phrases"`
	Phrases_Cit string `json:"phrases_cit"`
	Tag2        string `json:"tag"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}
type AutoGenerated struct {
	Nodes []struct {
		Label      string  `json:"label"`
		X          float64 `json:"x"`
		Y          float64 `json:"y"`
		ID         string  `json:"id"`
		Attributes struct {
			Cluster string `json:"Cluster"`
		} `json:"attributes"`
		Color string  `json:"color"`
		Size  float64 `json:"size"`
	} `json:"nodes"`
	Edges []struct {
		Source     string `json:"source"`
		Target     string `json:"target"`
		ID         string `json:"id"`
		Attributes struct {
		} `json:"attributes"`
		Color string  `json:"color"`
		Size  float64 `json:"size"`
	} `json:"edges"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
	// Url      string   `json:"url"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}

type Graph3 struct {
	Attributes struct {
	} `json:"attributes"`
	Nodes []Node  `json:"nodes"`
	Edges []Edge3 `json:"edges"`
	// Url      string   `json:"url"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}

type Graph2 struct {
	Nodes []Node2 `json:"nodes"`
	Edges []Edge2 `json:"edges"`
	// Url      string   `json:"url"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}

type Node2 struct {
	Attributes Attributes2 `json:"attributes"`
	Color      string      `json:"color"`
	Key        string      `json:"key"`
	Label      string      `json:"label"`
	Size       float64     `json:"size"`
	X          float64     `json:"x"`
	Y          float64     `json:"y"`
}

type Node struct {
	Key        string     `json:"key"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Label         string  `json:"label"`
	Country       string  `json:"country"`
	Tag           string  `json:"tag"`
	Title         string  `json:"title"`
	WOSID_Wbubble string  `json:"wosid_wbubble"`
	Url           string  `json:"url"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
}

type Attributes2 struct {
	//Year string `json:"Year"`
}
type Attributes3 struct {
	Type string `json:"type"`
}

type Edge struct {
	Key        string `json:"key"`
	Source     string `json:"source"`
	Target     string `json:"target"`
	Undirected bool   `json:"undirected"`
}
type Edge3 struct {
	Key        string      `json:"key"`
	Source     string      `json:"source"`
	Target     string      `json:"target"`
	Undirected bool        `json:"undirected"`
	Attributes Attributes3 `json:"attributes"`
}

type Edge2 struct {
	Attributes Attributes2 `json:"attributes"`
	Color      string      `json:"color"`
	Key        string      `json:"key"`
	Size       float64     `json:"size"`
	Source     string      `json:"source"`
	Target     string      `json:"target"`
}

type Paper_str struct {
	// struct for unique paper table; for join between text paper data and scrape paper data
	Citer_Year   string `json:"citer_year"`
	Citer_Name   string `json:"citer_name"`
	Citer_Conc   string `json:"citer_conc"`
	Citer_Writer string `json:"citer_writer"`
	Cites_Nr     int    `json:"cites_nr"`
}

type Nde struct {
	Key           string  `json:"key"`
	Label         string  `json:"label"`
	Tag           string  `json:"tag"`
	URL           string  `json:"URL"`
	Title         string  `json:"title"`
	WOSID_Wbubble string  `json:"wosid_wbubble"`
	Cluster       string  `json:"cluster"`
	X             float64 `json:"x"`
	Y             float64 `json:"y"`
	Score         float64 `json:"score"`
}
type Clstr struct {
	Key          string `json:"key"`
	Color        string `json:"color"`
	ClusterLabel string `json:"clusterLabel"`
}
type Tag struct {
	Key   string `json:"key"`
	Image string `json:"image"`
}

type Full_Graph struct {
	Nodes    []Nde      `json:"nodes"`
	Edges    [][]string `json:"edges"`
	Clusters []Clstr    `json:"clusters"`
	Tags     []Tag      `json:"tags"`
}

type Paper_master struct {
	// struct for unique paper table; for join between text paper data and scrape paper data
	Citer_Year   int    `json:"year"`
	Citer_Name   string `json:"citer_name"`
	Citer_Writer string `json:"writer"`
	Cites_Nr     int    `json:"cites_nr"`
	Tag          string `json:"tag"`
}

type Paper_masterv2 struct {
	// struct for unique paper table; for join between text paper data and scrape paper data
	Citer_Year    int    `json:"year"`
	Citer_Name    string `json:"citer_name"`
	Citer_Writer  string `json:"writer"`
	Cites_Nr      int    `json:"cites_nr"`
	Tag           string `json:"tag"`
	WOSID_Wbubble string `json:"wosid_wbubble"`
}

//includes "" ??? pruefen, aber ist in beiden, json spezifisch string marker
//Cites_Nr fix

type Paper_Scraped struct {
	Name     string          `json:"name"`
	Year     int             `json:"year"`
	Yearinfo string          `json:"yearinfo"`
	Writer   []string        `json:"writer"`
	Cites_Nr int             `json:"cites_nr"`
	Cites    []Paper_Scraped `json:"cites"` //compares to cited by
	Url      string          `json:"url"`
	//Topics   []string `json:"topics"`
	File string `json:"file"`
}

type Article struct {
	Name string
	Year int
	//Yearinfo string
	Writer string
	// Cites_Nr int
	// Cites    []Paper
	// Url      string
	//Topics   []string

}

type Relation struct {
	Citer_Name    string
	Citer_Year    int
	Citer_Writer  string
	Citing_Name   string
	Citing_Year   int
	Citing_Writer string
	Cites_Nr      int
	Cites_Nr2     int
	//Cites_Nr      int
	//Topics
}

type Relationv2 struct {
	Citer_Name    string
	Citer_Year    int
	Citer_Writer  string
	Citing_Name   string
	Citing_Year   int
	Citing_Writer string
	Cites_Nr      int
	Cites_Nr2     int
	Phrases       string
	Tag           string
	Tag2          string
	//Cites_Nr      int
	//Topics
}

type Relationv3 struct {
	Citer_Name    string
	Citer_Year    int
	Citer_Writer  string
	Citing_Name   string
	Citing_Year   int
	Citing_Writer string
	Cites_Nr      int
	Cites_Nr2     int
	Phrases       string
	Phrases_Cit   string
	Tag           string
	Tag2          string

	//Cites_Nr      int
	//Topics
}

type Relationv4 struct {
	Citer_Name    string `json:"citer_name"`
	Citer_Year    int    `json:"citer_year"`
	Citer_Writer  string `json:"citer_writer"`
	Citing_Name   string `json:"citing_name"`
	Citing_Year   int    `json:"citing_year"`
	Citing_Writer string `json:"citing_writer"`
	Cites_Nr      int    `json:"cites_nr"`
	Cites_Nr2     int    `json:"cites_nr2"`
	Phrases       string `json:"phrases"`
	Phrases_Cit   string `json:"phrases_cit"`
	Tag           string `json:"tag"`
	Tag2          string `json:"tag2"`
	WOSID_Wbubble string `json:"wosid_wbubble"`
	//Cites_Nr      int
	//Topics
}

type Relation_wos struct {
	Citer_Name    string
	Citer_Year    int
	Citer_Writer  string
	WOSID         string
	Citing_Name   string
	Citing_Year   int
	Citing_Writer string
	Cites_Nr      int
	Cites_Nr2     int
	Phrases       string
	Phrases_Cit   []string
	Tag           string
	Tag2          string

	//Cites_Nr      int
	//Topics
}

type Relation_fin struct {
	Citer_Name    string
	Citer_Year    int
	Citer_Writer  string
	Citing_Name   string
	Citing_Year   int
	Citing_Writer string
	Cites_Nr      int
	Cites_Nr2     int
	Phrases       string
	Phrases_Cit   string
	Tag           string
	Tag2          string
	WOSID_Wbubble string

	//Cites_Nr      int
	//Topics
}

// type Paper struct {
// 	Name     string
// 	Year     int
// 	//Yearinfo string
// 	Writer   string
// 	// Cites_Nr int
// 	// Cites    []Paper
// 	// Url      string
// 	//Topics   []string

// }

// func read_json(text string) []string {

// }

func analyze_dref(text string) []string {
	//cited_papers := make([]Paper, 0)
	//s = make([]byte, 5, 5)
	//var ref []string
	var listref []string
	title := "missing"
	year := "0"
	author := "missing"
	//SEARCH
	bef, aft, bools := strings.Cut(text, "“")
	bef2, aft2, bools2 := strings.Cut(text, "\"")
	//RICHTIGE BEF AND AFTER FÜR DEN FALL-
	ifemptys := strings.Replace(bef, " ", "", -1)
	ifempty := strings.Replace(ifemptys, "\n", "", -1)
	//rex := regexp.MustCompile(`\“{1}[\w\s\:]+\”{1}`)
	//cuttitle := rex.FindAllString(text, -1)
	//fmt.Println(len(cuttitle))
	// for _, el := range cuttitle {

	// 	fmt.Println(el, "\n")
	// 	fmt.Println("\n")
	// }

	//start := strings.Index(text, "References")
	// fmt.Println("bef", bef, "\n")
	// fmt.Println("aft", aft, "\n")
	// fmt.Println("aft", bools, "\n")
	//bol := len(ifempty) > 0
	// fmt.Println("aft", bol, "\n")
	//_, aft, bools := strings.Cut(text, "References")
	if bools && len(ifempty) > 0 {
		author = bef
		rest := aft
		beff, aftt, boolss := strings.Cut(rest, "”")
		if boolss {
			title = beff
			restt := aftt

			//var listref []string FÜR AUTOR WEITER ZERTEILEN? nicht nötig;

			//fmt.Println(str1, "\n")
			//re := regexp.MustCompile(`^[a-zA-Z\,\.\s]+[a-zA-Z\"\:\-\s]+[\w\"\.\(\)\-\s]+$`)
			//BREAKER [a-zA-Z\"\s\:\-]?
			re := regexp.MustCompile(`([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})`)

			//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern

			//fmt.Printf("String contains any match: %v\n", re.MatchString(str1)) // True

			//submatchall := re.FindAllString(str1, -1)
			//submatchall := re.FindAllStringIndex(restt, -1)
			submatchalls := re.FindAllString(restt, -1)
			// fmt.Println(len(submatchalls))
			if len(submatchalls) > 0 {
				year = submatchalls[0]
			}

			//listref = append(listref, aft[sta:end])

			// fmt.Println("submatchalls", "\n")
			// for _, element := range submatchalls {

			// 	fmt.Println(element, "\n")
			// 	fmt.Println("\n")
			// }
			//var listref string
		}

	} else if len(ifempty) > 0 && bools2 {
		author = bef2
		rest := aft2
		beff, aftt, boolss := strings.Cut(rest, "\"")
		if boolss {
			title = beff
			restt := aftt

			//var listref []string FÜR AUTOR WEITER ZERTEILEN? nicht nötig;

			//fmt.Println(str1, "\n")
			//re := regexp.MustCompile(`^[a-zA-Z\,\.\s]+[a-zA-Z\"\:\-\s]+[\w\"\.\(\)\-\s]+$`)
			//BREAKER [a-zA-Z\"\s\:\-]?
			re := regexp.MustCompile(`([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})`)

			//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern

			//fmt.Printf("String yyy any match: %v\n", re.MatchString(str1)) // True

			//submatchall := re.FindAllString(str1, -1)
			//submatchall := re.FindAllStringIndex(restt, -1)
			submatchalls := re.FindAllString(restt, -1)
			// fmt.Println(len(submatchalls))
			if len(submatchalls) > 0 {
				year = submatchalls[0]
			}

			//listref = append(listref, aft[sta:end])

			// fmt.Println("submatchalls", "\n")
			// for _, element := range submatchalls {

			// 	fmt.Println(element, "\n")
			// 	fmt.Println("\n")
			// }
			//var listref string
		}

	} else if len(ifempty) > 0 {
		rexx := regexp.MustCompile(`\.{1}\,{1}`)
		//subs := rexx.FindAllString(bef, -1)
		subsind := rexx.FindAllStringIndex(bef, -1)
		if len(subsind) == 1 {
			author = bef[0:subsind[len(subsind)-1][1]]
			rest := bef[subsind[len(subsind)-1][1]+1:]
			//beff, aftt, boolss := strings.Cut(rest, "”")
			//if boolss {

			//restt := aftt

			//var listref []string FÜR AUTOR WEITER ZERTEILEN? nicht nötig;

			//fmt.Println(str1, "\n")
			//re := regexp.MustCompile(`^[a-zA-Z\,\.\s]+[a-zA-Z\"\:\-\s]+[\w\"\.\(\)\-\s]+$`)
			//BREAKER [a-zA-Z\"\s\:\-]?
			re := regexp.MustCompile(`([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})`)

			submatchalls := re.FindAllString(rest, -1)
			submatchs := re.FindAllStringIndex(rest, -1)
			// fmt.Println(len(submatchalls))
			if len(submatchalls) > 0 {
				year = submatchalls[0]
			}
			if len(submatchs) > 0 {
				//fmt.Println(rest)
				title_pre := rest[0 : submatchs[len(submatchs)-1][0]-1]
				// fmt.Println(title_pre, "\n")
				// bef, aft, bools := strings.Cut(author, "\"")
				// fmt.Println(bef, "\n")
				// fmt.Println(aft, "\n")
				// fmt.Println(bools, "\n")
				// fmt.Println(title, "\n")
				// fmt.Println(year, "\n")
				// fmt.Println(author, "\n")
				ind := strings.Index(title_pre, ".")
				// fmt.Println(ind)
				if ind != -1 {
					title = title_pre[0:ind]
				} else {
					title = author
				}

			}
		} else if len(subsind) > 1 {
			author = bef[0:subsind[0][1]]
			rest := bef[subsind[0][1]+1:]
			//beff, aftt, boolss := strings.Cut(rest, "”")
			//if boolss {

			//restt := aftt

			//var listref []string FÜR AUTOR WEITER ZERTEILEN? nicht nötig;

			//fmt.Println(str1, "\n")
			//re := regexp.MustCompile(`^[a-zA-Z\,\.\s]+[a-zA-Z\"\:\-\s]+[\w\"\.\(\)\-\s]+$`)
			//BREAKER [a-zA-Z\"\s\:\-]?
			re := regexp.MustCompile(`([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})`)

			submatchalls := re.FindAllString(rest, -1)
			submatchs := re.FindAllStringIndex(rest, -1)
			// fmt.Println(len(submatchalls))
			if len(submatchalls) > 0 {
				year = submatchalls[0]
			}
			if len(submatchs) > 0 {
				// fmt.Println(rest)
				title_pre := rest[0 : submatchs[len(submatchs)-1][0]-1]
				// fmt.Println(title_pre, "\n")
				// bef, aft, bools := strings.Cut(author, "\"")
				// fmt.Println(bef, "\n")
				// fmt.Println(aft, "\n")
				// fmt.Println(bools, "\n")
				// fmt.Println(title, "\n")
				// fmt.Println(year, "\n")
				// fmt.Println(author, "\n")
				ind := strings.Index(title_pre, ".")
				// fmt.Println(ind)
				if ind != -1 {
					title = title_pre[0:ind]
				} else {
					title = author
				}

			}
		} else {
			//DIESEN BZW OBEN DRÜBER ALS NEUEN CASE MACHEN; OHNE BEF AFT; ÜBER EINTEILUNG DURCH JAHR;
			//AUSNAHMEFALL KEIN AUTOR TRENNBAR da alles mit , getrennt , weg über and/s end index+1 bis searchind(,)

			//restt := aftt
			rf := bef
			//var listref []string FÜR AUTOR WEITER ZERTEILEN? nicht nötig;
			if !bools2 {
				rf = text
			}

			fmt.Println("HIHIER")
			fmt.Println(bools2)
			fmt.Println(rf)

			//re := regexp.MustCompile(`^[a-zA-Z\,\.\s]+[a-zA-Z\"\:\-\s]+[\w\"\.\(\)\-\s]+$`)
			//BREAKER [a-zA-Z\"\s\:\-]?
			// OR CASE MACHEN MIT 2000 jahren und 1000 jahren; or http wegreinigen;
			re := regexp.MustCompile(`(([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\.{1})|(([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\,{1})`)
			notopt := regexp.MustCompile(`([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})`) //BESSER NOCH JAHR MIT ODER MACHEN

			submatchalls := re.FindAllString(rf, -1)
			submatchs := re.FindAllStringIndex(rf, -1)
			fmt.Println("Testnow")

			fmt.Println(len(submatchalls))
			if len(submatchalls) == 0 {

				submatchalls = notopt.FindAllString(rf, -1)
				submatchs = notopt.FindAllStringIndex(rf, -1)
			}
			fmt.Println(len(submatchs))
			// fmt.Println(len(submatchalls))
			if len(submatchalls) > 0 {
				year = submatchalls[0]
			} else {
				year = "0"
			}
			if len(submatchs) > 0 {
				title_pre := rf[0 : submatchs[0][0]-1]
				fmt.Println("titlepre" + title_pre)

				if submatchs[0][1]+1 < len(rf) {
					//PUNKTGETRENNTE FAELLE WICHTIG NICHT ERSTER PUNKT SONDERN ZWEITER so auch zahl rausbekommen
					title_post := rf[submatchs[0][1]+1:]
					fmt.Println("check1")
					rx := regexp.MustCompile(`\.`)
					divind := rx.FindAllStringIndex(title_pre, -1)
					divind2 := rx.FindAllStringIndex(title_post, -1)
					if len(divind) > 0 && len(divind2) > 0 {
						fmt.Println("check2")
						////PUNKT EVT NOCH ENTFERNEN UND ZAHL NOCH ENTFERNEN AUS AUTHOR ODER SCHON VORHER AUS REF FUNC ODER TESTEN WEG
						/// MIT LBREAK ZAHL. LEERSPACE
						//VERSCHIEDENE VARAINTEN ZUSAMMENFÜHREN SPÄTER
						author = title_pre //[0:ind]
						title = title_post[0:divind2[0][0]]
						fmt.Println("author" + author)
						fmt.Println("year" + year)
						fmt.Println("title" + title)
					} else {
						title = title_pre //[0:ind]
						author = "missing"
					}
				} else {
					//DIESE VARIANTE NOCH VERBESSERN
					title = title_pre //[0:ind]
					author = "missing"
				}
				//evt noch bef check einabue CHECKEN OB JETZT PROBLEM WENN ERSTE GENOMMEN; vorher bef[0 : submatchs[len(submatchs)-1][0]-1]

				//ind := strings.Index(title_pre, ".")

			}

		}
		// else {
		// 	title = "missing"
		// 	year = "0"
		// 	author = "missing"

		// }
		//und noch default case = 0 all attributes

	} else {
		////////////////////////////EXISTING METHODS FAILED//////////////////Emergencymethod
		// re := regexp.MustCompile(`[12]{1}[890]{1}[0-9]{2}`)
		title = "missing"
		year = "0"
		author = "missing"
		// submatchalls := re.FindAllString(bef, -1)
		// submatchs := re.FindAllStringIndex(bef, -1)
		// fmt.Println(len(submatchalls))
		// if len(submatchalls) > 0 {
		// 	year = submatchalls[0]
		// }

		// title_pre := bef[0 : submatchs[len(submatchs)-1][0]-1]
		// //ind := strings.Index(title_pre, ".")
		// title = title_pre
	}
	//[0:ind]
	//CUT

	//DIVIDE
	//for _, element := range listref {
	// fmt.Println(author, "\n")
	// fmt.Println("\n")
	// fmt.Println(title, "\n")
	// fmt.Println("\n")
	// fmt.Println(year, "\n")
	// fmt.Println("\n")
	//}
	//SAVE
	if author != "missing" && year != "0" && title != "missing" { // ÜBERLEGEN OB DOCH ZULASSEN IN MANCHEN FÄLLEN; JAHR AUTOR REICHT;
		author = strings.Replace(author, "\n", " ", -1)
		author = strings.Replace(author, "DOI:", " ", -1) //VL AUHC IN TITLE
		rxx := regexp.MustCompile(`[0-9]+\.{1}`)
		nrinauthor := rxx.FindAllString(author, -1)
		for _, element := range nrinauthor {
			author = strings.Replace(author, element, " ", -1)
		}
		author = strings.TrimSpace(author)
		author = strings.TrimSuffix(author, ".")
		author = strings.TrimSuffix(author, ",")
		author = strings.TrimSpace(author)

		title = strings.Replace(title, "\n", " ", -1)
		title = strings.TrimSpace(title)
		title = strings.TrimSuffix(title, ".")
		title = strings.TrimSuffix(title, ",")
		title = strings.TrimSpace(title)

		year = strings.Replace(year, "\n", " ", -1)
		year = strings.TrimSpace(year)
		year = strings.TrimSuffix(year, ".")
		year = strings.TrimSuffix(year, ",")
		year = strings.TrimSpace(year)
		//UND DANACH NOCHMAL EXTRA TRIM VORNE ANFANG LEERZEICHEN
		listref = append(listref, author)
		listref = append(listref, title)
		listref = append(listref, year)
	} else {
		listref = nil
	}
	for _, element := range listref {
		fmt.Println("dref")
		fmt.Println(element, "\n")
		fmt.Println("\n")
	}

	return listref
	//err mit returnen später
}

func analyze_ref(text string) []string {
	//cited_papers := make([]Paper, 0)
	//s = make([]byte, 5, 5)
	//var ref []string
	var listref []string
	//var solution := "" IN REF MUSS MAN DAS ÜBER APPEND MACHEN DURCH DIE VIELEN EINZELNEN FÄLLE
	var submatchall [][]int
	var submatchalls []string
	//var helpind int
	//var helpstr string
	//SEARCH
	start := strings.Index(text, "References")

	fmt.Println(start, "\n")
	_, aft, bools := strings.Cut(text, "References")
	if !bools {
		_, aft, bools = strings.Cut(text, "REFERENCES")
	}
	if bools {
		str1 := aft
		//fmt.Println(str1, "\n")
		//re := regexp.MustCompile(`^[a-zA-Z\,\.\s]+[a-zA-Z\"\:\-\s]+[\w\"\.\(\)\-\s]+$`)
		//BREAKER [a-zA-Z\"\s\:\-]?
		re := regexp.MustCompile(`\[{1}[0-9]+\]{1}`)

		//fmt.Printf("Pattern: %v\n", re.String()) // Print Pattern

		//fmt.Printf("String contains any match: %v\n", re.MatchString(str1)) // True

		//submatchall := re.FindAllString(str1, -1)
		submatchall = re.FindAllStringIndex(str1, -1)
		//submatchalls = re.FindAllString(str1, -1)       ///////////////////TESTER
		fmt.Println(len(submatchall))
		if len(submatchall) > 0 {
			//fmt.Println(len(submatchall))
			//var listref string

			//struct bauen
			//und einzelne teile abspeichern

			for i := 0; i < len(submatchall)-1; i++ {
				// fmt.Println(i)
				refr := Reference{}
				sta := submatchall[i][1]
				refr.Startind = submatchall[i][1]

				//} else {
				end := submatchall[i+1][0]
				refr.Endind = submatchall[i+1][0]
				//}

				if len(aft[sta:end]) < 1000 {
					//fmt.Println("tester:" + aft[sta:end])
					listref = append(listref, aft[sta:end])
				} else {
					// fmt.Println("HI1")
					// fmt.Println(aft[sta:end]) TESTER///////////////////////////
					txt := aft[sta:end]
					// fmt.Println("HI1HI") //////////TESTER
					rx := regexp.MustCompile(`\.{1}\n{1}`)
					submatchs := rx.FindAllStringIndex(txt, -1)
					// fmt.Println("subs:")
					// fmt.Println(len(submatchs)) TESTER///////////

					reff := Reference{}
					//sta_1 := submatchall[len(submatchall)-1][1]
					reff.Startind = submatchall[len(submatchall)-1][1]
					end_1 := submatchs[0][0]
					reff.Endind = submatchs[0][0]
					// fmt.Println("txt:")
					// fmt.Println(txt[:end_1])
					listref = append(listref, txt[:end_1])

					for i := 0; i < len(submatchs)-2; i++ {
						fmt.Println(i)
						refr := Reference{}
						sta := submatchs[i][1]          //+ 1
						refr.Startind = submatchs[i][1] //+ 1

						//} else {
						end := submatchs[i+1][0]        //-1
						refr.Endind = submatchs[i+1][0] //- 1
						//}
						//fmt.Println("RESP:" + txt[sta:end])
						listref = append(listref, txt[sta:end])
						//list add References

					}

				}
				//list add References

			}
			//if i >= len(submatchall)-1 {
			reff := Reference{}
			sta := submatchall[len(submatchall)-1][1]
			reff.Startind = submatchall[len(submatchall)-1][1]
			end := len(aft) - 1
			reff.Endind = len(aft) - 1
			if len(aft[sta:end]) < 1000 {

				listref = append(listref, aft[sta:end])
			} else {
				// fmt.Println("HI2")
				// fmt.Println(aft[sta:end]) TESTER //////////////
				txt := aft[sta:end]
				//fmt.Println("HI1HI")
				rx := regexp.MustCompile(`\.{1}\n{1}`)
				submatchs := rx.FindAllStringIndex(txt, -1)
				//fmt.Println("subs:")
				//fmt.Println(len(submatchs))

				reff := Reference{}
				//sta_1 := submatchall[len(submatchall)-1][1]
				reff.Startind = submatchall[len(submatchall)-1][1]
				end_1 := submatchs[0][0]
				reff.Endind = submatchs[0][0]
				listref = append(listref, txt[:end_1])

				for i := 0; i < len(submatchs)-2; i++ {
					// fmt.Println(i)
					refr := Reference{}
					sta := submatchs[i][1]          //+ 1
					refr.Startind = submatchs[i][1] //+ 1

					//} else {
					end := submatchs[i+1][0]        //- 1
					refr.Endind = submatchs[i+1][0] //- 1
					//}
					// fmt.Println("RESP:" + txt[sta:end])
					listref = append(listref, txt[sta:end])
					//list add References

				}

			}

			// for _, element := range submatchall {
			// 	sta :=element[1]
			// 	end :=element[1]
			// 	listref = append(listref, aft[sta:end])
			// 	//list add References
			// 	refr.Startind = element[0]
			// 	refr.Endind = element[1]
			// }

			// for _, element := range listref { ///////////////////////TESTFUNCTION
			// 	fmt.Println(element, "\n")
			// 	fmt.Println("\n")
			// }
			// fmt.Println("submatchalls", "\n")   /////////////////////////TESTER
			// for _, element := range submatchalls {

			// 	fmt.Println(element, "\n")
			// 	fmt.Println("\n")
			// }
		} else { //\s{1}\[A-Z]{1} ///////////////////////MARKER11
			/////////////////////////////////////////FIX THIS
			/////////////////////////////////////////FEHLT HIER NICHT EIN DOPPELPUNKT?
			//////////////////////////////////// . LINEBREAK PROBIEREN;
			rex := regexp.MustCompile(`[0-9]+\s{1}\w{1}`) ////////LINEBREAK UND ZAHL LINEBREAK UND ZAHL : PRBOBIEREN ETC;

			//FALLS DAS BEDEUTET DAS JA AUCH; DASS HIER DAS ERSTE WORT MITENTHALTEN IST; BEI REX2 darunter fängt man ja am ende an;
			/// TEST OBERHALB MIT ODER; ABER EVT SO WIEDER ZURÜCK UND DANN ÜBER . variante darunter oder über zahL: anfang getrennt ohne wort
			rex2 := regexp.MustCompile(`([0-9]+\.{1}\n{1})|(DOI\:{1}\n{1})`) /////ODER PUNKTE ENDE LINEBREAK // geht notfalls normal
			submatchall = rex.FindAllStringIndex(str1, -1)
			submatchalls = rex.FindAllString(str1, -1)
			submatchall2 := rex2.FindAllStringIndex(str1, -1)
			submatchalls2 := rex2.FindAllString(str1, -1)
			fmt.Println("EXTRACASE")
			fmt.Println(len(submatchalls))
			//fmt.Println(submatchalls[0]) HIIIER
			fmt.Println("EXTRACASE2")
			//fmt.Println(len(submatchalls2))
			fmt.Println(submatchalls2[0])
			//fmt.Println(len(submatchall))
			//var listref string
			//CHECKN DA UMGANG MIT SUBMATCHALL ANDERS ALS MIT SUBMATCHALL2 da eins am ende cuttet und eins am anfang;
			if len(submatchall2) > len(submatchall) {
				submatchall = submatchall2
				submatchalls = submatchalls2
			}

			//struct bauen
			//und einzelne teile abspeichern

			for i := 0; i < len(submatchall)-1; i++ {
				//null and index overlauf schutz einbauen!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!![0-9]{1-2}
				//ACHTUNG FANGE bei treffer bei ende des treffers an und wenn wort enthalten fehlt
				//das erste wort und wenn punkt am ende fehlt erste referenz ganz; VERBESSERN; AKTUELL DIRTY MIT . variante
				fmt.Println(i)
				refr := Reference{}
				sta := submatchall[i][1] - 1
				refr.Startind = submatchall[i][1] - 1

				//} else {
				end := submatchall[i+1][0]
				refr.Endind = submatchall[i+1][0]
				//}
				//fmt.Println("tester:" + aft[sta:end])
				listref = append(listref, aft[sta:end])
				//list add References

			}
			//if i >= len(submatchall)-1 {
			reff := Reference{}
			sta := submatchall[len(submatchall)-1][1] - 1
			reff.Startind = submatchall[len(submatchall)-1][1] - 1
			end := len(aft) - 1
			reff.Endind = len(aft) - 1
			// fmt.Println("tester:" + aft[sta:end])
			listref = append(listref, aft[sta:end])
			// for _, element := range submatchall {
			// 	sta :=element[1]
			// 	end :=element[1]
			// 	listref = append(listref, aft[sta:end])
			// 	//list add References
			// 	refr.Startind = element[0]
			// 	refr.Endind = element[1]
			// }

			fmt.Println("strtreffer", "\n")
			fmt.Println(len(submatchalls))
			// for _, element := range submatchalls {

			// 	fmt.Println(element, "\n")
			// 	fmt.Println("\n")
			// }
		}

		//end:= strings.Index(text, "References")
	}
	for _, element := range listref {
		fmt.Println("ref")
		fmt.Println(element, "\n")
		fmt.Println("\n")
	}
	//CUT

	//DIVIDE

	//SAVE

	return listref
	//err mit returnen später
}

//take database publication collection(txt file) as an input and transform to units of publications (slices of strings) containing meta information
func analyze_txtws(text string) []string {
	var listpub []string

	//var solution := "" IN REF MUSS MAN DAS ÜBER APPEND MACHEN DURCH DIE VIELEN EINZELNEN FÄLLE
	var submatchall [][]int
	var submatchalls []string
	//var helpind int
	//var helpstr string
	//SEARCH

	rex := regexp.MustCompile(`AU {1}`)

	submatchalls = rex.FindAllString(text, -1)
	submatchall = rex.FindAllStringIndex(text, -1)
	if len(submatchall) > 0 {
		fmt.Println("READY1")
		fmt.Println(len(submatchall))
		for _, su := range submatchalls {
			fmt.Println(su)
		}
		for i := 0; i < len(submatchall)-1; i++ {
			bef, _, bools := strings.Cut(text[submatchall[i][1]:submatchall[i+1][0]], "VL ")
			if bools {
				listpub = append(listpub, bef)
				fmt.Println(bef)
			}
			// rex := regexp.MustCompile(`AU{1}`)

			// submatchalls = rex.FindAllString(text, -1)
			// submatchall = rex.FindAllStringIndex(text, -1)
		}
		bef, _, bools := strings.Cut(text[submatchall[len(submatchall)-1][1]:], "VL ")
		if bools {
			listpub = append(listpub, bef)
			fmt.Println(bef)
		}
	}

	//start := strings.Index(text, "References")

	//fmt.Println(start, "\n")
	// _, aft, bools := strings.Cut(text, "References")
	// if !bools {
	// 	_, aft, bools = strings.Cut(text, "REFERENCES")
	// }
	// if bools {

	//
	return listpub
}

func analyze_pubws(text string) []Relationv4 {
	var listrel []Relationv4

	var listpub []string

	//var solution := "" IN REF MUSS MAN DAS ÜBER APPEND MACHEN DURCH DIE VIELEN EINZELNEN FÄLLE
	//var submatchall [][]int
	//var submatchalls []string
	//var helpind int
	//var helpstr string
	//SEARCH
	if len(text) > 0 {
		rel := Relationv4{}
		var cited_ref string

		rex := regexp.MustCompile(`AF {1}`)

		af_ind := rex.FindStringIndex(text)
		if len(af_ind) > 0 {
			writer1 := chg_fmt2(text[:af_ind[0]])
			rel.Citer_Writer = writer1
		}

		rex = regexp.MustCompile(`TI {1}`)

		ti_ind := rex.FindStringIndex(text)

		rex = regexp.MustCompile(`SO {1}`)

		so_ind := rex.FindStringIndex(text)
		if len(so_ind) > 0 && len(ti_ind) > 0 {
			rel.Citer_Name = text[ti_ind[1]:so_ind[0]]
		}
		rex = regexp.MustCompile(`CR {1}`)

		cr_ind := rex.FindStringIndex(text)

		rex = regexp.MustCompile(`NR {1}`)

		nr_ind := rex.FindStringIndex(text)
		if len(nr_ind) > 0 && len(cr_ind) > 0 {
			cited_ref = text[cr_ind[1]:nr_ind[0]]
		}
		// rex = regexp.MustCompile(`TC {1}`)

		rex = regexp.MustCompile(`PY {1}`)

		py_ind := rex.FindStringIndex(text)
		if len(py_ind) > 0 {
			// fmt.Println("YEAR")
			// fmt.Println(text[py_ind[1]:])
			// fmt.Println(len(text[py_ind[1]:]))
			// bef, _, boolls := strings.Cut(text[py_ind[1]:], " ")

			intVar, err := strconv.Atoi(text[py_ind[1] : py_ind[1]+4])
			if err != nil {
				rel.Citer_Year = 0
			} else {
				rel.Citer_Year = intVar
			}
		}
		if len(cited_ref) > 0 {
			fmt.Println("READY77")
			fmt.Println(len(cited_ref))

			lnbrk := regexp.MustCompile(`\n{1}`)

			lnbrk_ind := lnbrk.FindAllStringIndex(cited_ref, -1)

			if len(lnbrk_ind) > 0 {
				publ := cited_ref[:lnbrk_ind[0][0]]
				listpub = append(listpub, publ)
				for i := 1; i < len(lnbrk_ind); i++ {
					publ = cited_ref[lnbrk_ind[i-1][1]:lnbrk_ind[i][0]]
					listpub = append(listpub, publ)
					// rex := regexp.MustCompile(`AU{1}`)
					//fmt.Println("TESTACHTUNG")
					//fmt.Println(publ)
					// submatchalls = rex.FindAllString(text, -1)
					// submatchall = rex.FindAllStringIndex(text, -1)
				}
			}
			rel.Cites_Nr = len(listpub)
			if len(listpub) > 0 {
				for _, li := range listpub {
					bef, aft, bools := strings.Cut(li, ", ")
					if bools {
						writer2 := chg_fmt(bef)
						rel.Citing_Writer = writer2
						// fmt.Println("WRITER")
						// fmt.Println(bef)
						beff, afft, boolss := strings.Cut(aft, ", ")
						if boolss {
							// fmt.Println("YEAR")
							// fmt.Println(beff)
							intVar, err := strconv.Atoi(beff)
							if err != nil {
								rel.Citing_Year = 0
							} else {
								rel.Citing_Year = intVar
							}
							beff, _, _ := strings.Cut(afft, ", ")

							rel.Citing_Name = beff
							// fmt.Println("NAME")
							// fmt.Println(beff)

						}
					}
					rel.Cites_Nr2 = -100
					rel.Phrases = ""
					rel.Phrases_Cit = ""
					rel.WOSID_Wbubble = ""
					listrel = append(listrel, rel)
				}
			}

		}
	}

	return listrel
}

// func analyze_refws(text string) []string {
// 	var listpub []string

// 	return listpub
// }

func chg_fmt(text string) string {
	res := ""
	text = strings.Trim(text, " ")

	rex := regexp.MustCompile(`\s{1}`)

	spc := rex.FindStringIndex(text)
	if len(spc) > 0 {
		name := text[:spc[0]]
		surname := text[spc[0]+1:]
		res = surname + " " + name
		fmt.Print("alt:")
		fmt.Println(text)
		fmt.Print("neu:")
		fmt.Println(res)
	}

	return res
}

func chg_fmt2(text string) string {
	//res := ""
	fin_res := ""
	text = strings.Trim(text, " ")

	r_ln := regexp.MustCompile(`\n{1}`)

	var wrt_lst []string
	var res_lst []string

	lns := r_ln.FindAllStringIndex(text, -1)
	fmt.Println("TEST1")
	fmt.Println(len(lns))
	if len(lns) > 0 {
		wrt_lst = append(wrt_lst, text[:lns[0][0]])
		for i := 1; i < len(lns); i++ {
			fmt.Println("TEST2")
			wrt_lst = append(wrt_lst, text[lns[i-1][0]+1:lns[i][0]])
		}
	}

	for i := 0; i < len(wrt_lst); i++ {
		r_spc := regexp.MustCompile(`,{1}`)
		spc := r_spc.FindStringIndex(wrt_lst[i])
		fmt.Println("TEST3")
		fmt.Println(len(spc))
		var res string
		if len(spc) > 0 {
			name := wrt_lst[i][:spc[0]]
			surname := wrt_lst[i][spc[0]+1:]
			res = surname + " " + name
		} else {
			res = wrt_lst[i]
		}

		fmt.Print("alt:")
		fmt.Println(wrt_lst[i])
		fmt.Print("neu:")
		fmt.Println(res)
		res_lst = append(res_lst, res)
	}
	fin_res = strings.Join(res_lst, ", ")
	whitespaces := regexp.MustCompile(`\s+`)
	fin_res = whitespaces.ReplaceAllString(fin_res, " ")
	fmt.Print("altCONC:")
	fmt.Println(text)
	fmt.Print("neuCONC:")
	fmt.Println(fin_res)
	return fin_res

}

func transform_strint(text string) int {
	res := 0
	intVar, err := strconv.Atoi(text)
	if err == nil {
		res = intVar
	}
	return res
}

//////////////////////////////////////////////////HELPER FUNCTIONS/////////////////////////////////////////////////////
// // 	//c.Visit(path)
// // 	fullURL := ""
// // 	full := ""
// // 	betw := ""
// // 	addit := ""
// // 	url1 := "https://scholar.google.de/scholar?start="
// // 	url2 := ""
// // 	pages := nr / 10

// // 	link := get_link(path)
//	if len(link) > 0 {
// // 		url2 = link[42:]
// // 		for i := 0; i <= pages; i++ {
// // 			if i == 0 {
// // 				full = path
// // 			} else {
// // 				fullURL = url1
// // 				betw = fmt.Sprint(i) + "0"
// // 				addit = url2
// // 				full = fullURL + betw + addit
// // 				//https://scholar.google.de/scholar?cites=15502415196660838898&as_sdt=2005&sciodt=0,5&hl=de
// // 				//https://scholar.google.de/scholar?start=0&hl=de&as_sdt=2005&sciodt=0,5&cites=15502415196660838898&scipsc=
// // 				fmt.Println("ACHTUNG!!!!!!!!!!!!!!!!" + full)
// // 			}

// // 			c.Visit(full)
// // 		}
// // 	}

func main() {
	//////////////////////////////////////////////////////////WIEDER NUTZEN UNIDOC
	// err, res := outputPdfText("pdfs/489426276.pdf")

	// listpdf, err := ioutil.ReadDir("C:/Users/neumu/pdffer/txt")
	// if err != nil {
	// 	panic(err)
	// }
	// //var def []map[int]string
	// // def := make([]map[int]string, 0)

	// // for _, li := range listpdf {
	// // 	// do something with the article
	// // 	fmt.Println(li.Name())
	// // 	err, pager := outputPdfText("pdf/"+li.Name(), li.Name())
	// // 	if err != nil {
	// // 		panic(err)
	// // 	} else {
	// // 		def = append(def, pager)
	// // 	}

	// // }

	// read WoS txt file 1
	content, err := ioutil.ReadFile("orgaPLAIN.txt")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(li.Name())
	//process and divide WoS publications file into subblocks
	wos_allpub := string(content)
	publications_ws := analyze_txtws(wos_allpub)
	//MARKER!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	fmt.Println(len(publications_ws))
	var listrelation_ws []Relationv4
	i := 0
	for _, pub := range publications_ws {

		references_ws := analyze_pubws(pub)
		fmt.Println("CECKER")
		fmt.Println(len(references_ws))

		for _, rel := range references_ws {

			ln := regexp.MustCompile(`\n{1}`)

			//submatchall := ln.FindAllStringIndex(rel.Citer_Writer , -1)

			rel.Citer_Writer = ln.ReplaceAllString(rel.Citer_Writer, ",")
			rel.Tag = "Web Source - Web of Science"
			rel.Tag2 = "Web Source - Web of Science (References)"
			rel.WOSID_Wbubble = ""
			//paper_citer := Paper{} //enthält informationen zu citationpaper not distinct data
			//paper_citing := Paper{}
			listrelation_ws = append(listrelation_ws, rel)
			i = i + 1
			// do something with the article
			//listrelation_ws = append(listrelation_ws, rel)
			//relation_pre := Relation{}

			//paper_citer.Name = data2[i].Cites[j].Name
			// relation_pre.Citer_Name = re[0]
			// relation_pre.Citing_Name = re[4]

			//paper_citer.Writer = strings.Join(data2[i].Cites[j].Writer, ", ") //first if "min" 2 elem, last "and"  and important when writer stats impl
			// relation_pre.Citer_Writer = re[1]
			// relation_pre.Citing_Writer = re[5] //first if "min" 2 elem, last "and"  and important when writer stats impl

			//paper_citer.Year = data2[i].Cites[j].Year
			// relation_pre.Citer_Year = transform_strint(re[2])
			// relation_pre.Citing_Year = transform_strint(re[6])
			// relation_pre.Cites_Nr2 = transform_strint(re[7])
			// relation_pre.Cites_Nr = transform_strint(re[3])
			//paper_citer.Cited_by = "" //NOT TRACKED BECAUSE DECIDED TO INCLUDE DATA ONLY ONE CITEDBY LEVEL DEEP FROM SCHOLAR

			//paper_citer.Cites_Nr = data2[i].Cites[j].Cites_Nr
			// fmt.Println("Paper aus cited_papers citing")
			// fmt.Println(paper_citing.Cited_by)
			// fmt.Println(paper_citing.Name)
			// fmt.Println(paper_citing.Writer)
			// fmt.Println(paper_citing.Year)
			// fmt.Println("\n")
			// fmt.Println("Paper aus cited_papers citer")
			// fmt.Println(paper_citer.Cited_by)
			// fmt.Println(paper_citer.Name)
			// fmt.Println(paper_citer.Writer)
			// fmt.Println(paper_citer.Year)
			// fmt.Println("\n")
			//listpaper = append(listpaper, paper_citer)

		}

		//cited_fin = append(cited_fin, cited_papers...)

	}

	js, err := json.MarshalIndent(listrelation_ws, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("listrelation_ws.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
		fmt.Print(i)
		fmt.Print(" references")
	} //MARKER!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

	// read WoS txt file 2
	// content, err = ioutil.ReadFile("orgaPLAIN2.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//fmt.Println(li.Name())
	//process and divide WoS publications file into subblocks

	// 	// Convert []byte to string and print to screen
	// 	text := string(content)
	// 	references := analyze_ref(text)

	// 	for _, re := range references { // TEST (http|https|ftp): WEGLASSEN
	// 		rex := regexp.MustCompile(`(http|https|ftp):[\/]{2}([a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,4})(:[0-9]+)?\/?([a-zA-Z0-9\-\._\?\,\'\/\\\+&amp;%\$#\=~]*)`)

	// 		submatchall := rex.FindAllString(re, -1)
	// 		if len(submatchall) > 0 {
	// 			fmt.Println("Foundyeyy")
	// 			fmt.Println(len(submatchall))
	// 			for _, su := range submatchall {
	// 				fmt.Println(su)
	// 				re = strings.Replace(re, su, "", -1)
	// 			}
	// 		}

	// 		ddref := analyze_dref(re)
	// 		if ddref != nil {
	// 			item := Paper{}
	// 			item.Cited_by = li.Name()
	// 			item.Writer = ddref[0]
	// 			item.Name = ddref[1]
	// 			intVar, err := strconv.Atoi(ddref[2])
	// 			if err != nil {
	// 				item.Year = 0
	// 			} else {
	// 				item.Year = intVar
	// 			}
	// 			i = i + 1
	// 			// do something with the article
	// 			cited_papers = append(cited_papers, item)
	// 		}

	// 	}
	// 	cited_fin = append(cited_fin, cited_papers...)

	// }

	// js, err := json.MarshalIndent(cited_fin, "", "    ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Writing data to file")
	// if err := os.WriteFile("cited_fin.json", js, 0664); err == nil {
	// 	fmt.Println("Data written to file successfully")
	// 	fmt.Print(i)
	// 	fmt.Print(" references")
	// }
	///////////////////////////////////////////////////////////////////////////////////////////////////////ABHIER
	var listpaper []Paper
	//var listpaper2 []Paperv2
	var listpapermaster []Paper_masterv2
	var listrelation []Relationv4
	//var prejoinpdf []Paper
	var prejoinpdf2 []Paperv3pre
	var relationwos []Relation_fin

	//Load txt reference data (Paper)
	file, _ := ioutil.ReadFile("cited_phrasesfin.json")
	//data := PaperNodes{}
	var data []Paperv2pre
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(data))
	for i := 0; i < len(data); i++ {
		data_all := Paperv3pre{}
		//are already paperobject could create relation with two empty attributes and add them later;
		// CITEDBY is the pdf name minus .pdf.txt (cut substr) - after manipulation compare to Name in Scraped
		//Create paper table entry and relation table entry later by joining table and relation dataframe;
		data_all.Cited_by = strings.TrimSuffix(data[i].Cited_by, ".txt")
		data_all.Cites_Nr = data[i].Cites_Nr
		data_all.Tag2 = "Text Source - References"
		data_all.Writer = data[i].Writer
		data_all.Year = data[i].Year
		data_all.Name = data[i].Name
		//phrasescit := ""
		//strings.Join(data[i].Phrases_Cit[:3], " - ")
		if len(data[i].Phrases_Cit) > 3 {
			data_all.Phrases_Cit = strings.Join(data[i].Phrases_Cit[:3], " - ")
		} else if len(data[i].Phrases_Cit) > 1 {
			data_all.Phrases_Cit = strings.Join(data[i].Phrases_Cit, " - ")
		} else if len(data[i].Phrases_Cit) == 1 {
			data_all.Phrases_Cit = data[i].Phrases_Cit[0]
		} else {
			data_all.Phrases_Cit = ""
		}

		if len(data[i].Phrases) > 0 {
			data_all.Phrases = data[i].Phrases
		} else {
			data_all.Phrases = ""
		}

		// data_all.Writer = data[i].Writer
		// data_all.Year = data[i].Year
		//data[i].Tag2 = "Text Source - References"

		prejoinpdf2 = append(prejoinpdf2, data_all)

	}
	fmt.Println("CHECKERFIN")

	file_wos, _ := ioutil.ReadFile("cited_phraseswosfin2.json")

	fmt.Println("CHECKERFIN2")
	//data := PaperNodes{}
	var data_wos []Relation_wos
	err = json.Unmarshal([]byte(file_wos), &data_wos)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("CHECKERFIN3")
	fmt.Println(len(data_wos))
	for i := 0; i < len(data_wos); i++ {
		data_all := Relation_fin{}
		//are already paperobject could create relation with two empty attributes and add them later;
		// CITEDBY is the pdf name minus .pdf.txt (cut substr) - after manipulation compare to Name in Scraped
		//Create paper table entry and relation table entry later by joining table and relation dataframe;
		data_all.Citer_Name = data_wos[i].Citer_Name
		data_all.Cites_Nr = data_wos[i].Cites_Nr
		data_all.Citer_Writer = data_wos[i].Citer_Writer
		data_all.Tag = data_wos[i].Tag
		data_all.Citer_Year = data_wos[i].Citer_Year

		//file:///C:/Users/neumu/sigma.js/demo/public/paper_wordcloud/000207197600002.svg
		data_all.WOSID_Wbubble = "file:///C:/Users/neumu/sigma.js/demo/public/paper_wordcloud/" + data_wos[i].WOSID + ".svg"

		// intVar, err := strconv.Atoi(refs[i][1])
		// 		if err != nil {
		// 			relation.Citing_Year = 0
		// 		} else {
		// 			relation.Citing_Year = intVar
		// 		}

		data_all.Citing_Name = data_wos[i].Citing_Name
		data_all.Citing_Writer = data_wos[i].Citing_Writer
		data_all.Citing_Year = data_wos[i].Citing_Year

		data_all.Cites_Nr2 = data_wos[i].Cites_Nr2
		data_all.Tag2 = data_wos[i].Tag2

		if len(data_wos[i].Phrases_Cit) > 3 {
			data_all.Phrases_Cit = strings.Join(data_wos[i].Phrases_Cit[:3], " - ")
		} else if len(data_wos[i].Phrases_Cit) > 1 {
			data_all.Phrases_Cit = strings.Join(data_wos[i].Phrases_Cit, " - ")
		} else if len(data_wos[i].Phrases_Cit) == 1 {
			data_all.Phrases_Cit = data_wos[i].Phrases_Cit[0]
		} else {
			data_all.Phrases_Cit = ""
		}

		if len(data_wos[i].Phrases) > 0 {
			data_all.Phrases = data_wos[i].Phrases
		} else {
			data_all.Phrases = ""
		}

		// data_all.Writer = data[i].Writer
		// data_all.Year = data[i].Year
		//data[i].Tag2 = "Text Source - References"

		relationwos = append(relationwos, data_all)

	}

	//listpaper2 = append(listpaper2, data...)
	//prejoinpdf2 = append(prejoinpdf2, data_all...)

	//Load scraper cites data (Paper_Scraped)
	var data2 []Paper_Scraped
	//READ AND ADD ALL JSON PARTS OF SCRAPED GOOGLE SCHOLAR DATA
	for i := 1; i < 46; i++ {
		file_str := fmt.Sprintf("cites_final%d.json", i)
		finstr := "cites_final/" + file_str
		file2, _ := ioutil.ReadFile(finstr)
		var data2_part []Paper_Scraped
		//var data2 Paper_Scraped
		_ = json.Unmarshal([]byte(file2), &data2_part)

		fmt.Println(len(data2_part))
		data2 = append(data2, data2_part...)
	}
	fmt.Println("SCRAPED PAPER ANZAHL")
	fmt.Println(len(data2))
	// fmt.Println("SCRAPED PAPER ANZAHL mit cites multipliziert")
	// fmt.Println(len(data2))

	for i := 0; i < len(data2); i++ {
		//vollstaendige Referenz erzeugbar mit jeweils 3 attributes
		paper_base := Paper_masterv2{} //normalised paper table
		paper_base.Citer_Name = data2[i].Name
		paper_base.Citer_Writer = strings.Join(data2[i].Writer, ", ")
		paper_base.Citer_Year = data2[i].Year
		paper_base.Cites_Nr = data2[i].Cites_Nr
		paper_base.WOSID_Wbubble = ""
		cites_nnr := data2[i].Cites_Nr
		paper_base.Tag = "Text Source - Full Text"

		// if len(data2[i].Cites) > 0 {
		// 	paper_base.Cites_nr = len(data2[i].Cites) - 1
		// 	cites_nnr = len(data2[i].Cites) - 1
		// } else {
		// 	paper_base.Cites_nr = 0
		// }

		listpapermaster = append(listpapermaster, paper_base)
		//paper cited by cant be filled for cites and for citing citer info and create paper for each cited by ref; cited by will be duplicate, and name potentially but less frequ
		//Create relation table entry
		//additionally paper table entry

		//first if "min" 2 elem, last "and"  and important when writer stats impl

		for j := 0; j < len(data2[i].Cites); j++ {
			//man koennte das objekt und citing attribute oberhalb setzen
			paper_citer := Paper{} //enthält informationen zu citationpaper not distinct data
			paper_citing := Paper{}

			relation_pre := Relationv4{}
			paper_citing.Name = data2[i].Name

			paper_citing.Writer = strings.Join(data2[i].Writer, ", ") //first if "min" 2 elem, last "and"  and important when writer stats impl

			paper_citing.Year = data2[i].Year

			paper_citing.Cited_by = data2[i].Cites[j].Name
			paper_citing.Cites_Nr = data2[i].Cites_Nr
			listpaper = append(listpaper, paper_citing)

			paper_citer.Name = data2[i].Cites[j].Name
			relation_pre.Citer_Name = data2[i].Cites[j].Name
			relation_pre.Citing_Name = data2[i].Name
			relation_pre.Tag = "Web Source - Google Scholar"
			relation_pre.Tag2 = "Text Source - Full Text"
			relation_pre.WOSID_Wbubble = ""

			paper_citer.Writer = strings.Join(data2[i].Cites[j].Writer, ", ") //first if "min" 2 elem, last "and"  and important when writer stats impl
			relation_pre.Citer_Writer = strings.Join(data2[i].Cites[j].Writer, ", ")
			relation_pre.Citing_Writer = strings.Join(data2[i].Writer, ", ") //first if "min" 2 elem, last "and"  and important when writer stats impl

			paper_citer.Year = data2[i].Cites[j].Year
			relation_pre.Citer_Year = data2[i].Cites[j].Year
			relation_pre.Citing_Year = data2[i].Year
			relation_pre.Cites_Nr2 = cites_nnr
			relation_pre.Cites_Nr = data2[i].Cites[j].Cites_Nr
			paper_citer.Cited_by = "" //NOT TRACKED BECAUSE DECIDED TO INCLUDE DATA ONLY ONE CITEDBY LEVEL DEEP FROM SCHOLAR

			paper_citer.Cites_Nr = data2[i].Cites[j].Cites_Nr
			relation_pre.Phrases = ""
			relation_pre.Phrases_Cit = ""
			// fmt.Println("Paper aus cited_papers citing")
			// fmt.Println(paper_citing.Cited_by)
			// fmt.Println(paper_citing.Name)
			// fmt.Println(paper_citing.Writer)
			// fmt.Println(paper_citing.Year)
			// fmt.Println("\n")
			// fmt.Println("Paper aus cited_papers citer")
			// fmt.Println(paper_citer.Cited_by)
			// fmt.Println(paper_citer.Name)
			// fmt.Println(paper_citer.Writer)
			// fmt.Println(paper_citer.Year)
			// fmt.Println("\n")
			listpaper = append(listpaper, paper_citer)
			listrelation = append(listrelation, relation_pre)
		}

	}

	/////////////////////////////////NEU LAST UPDATE APPEND

	//listrelation = append(listrelation, listrelation_ws...)
	fmt.Println(len(listrelation_ws))
	fmt.Println(len(listrelation))
	//for _, re := range listpaper {

	// 	fmt.Println(re)
	// }
	// for _, re := range listrelation {

	// 	fmt.Println(re)
	// }
	// for _, re := range references {
	// }

	///////////////////////(evt noch extra dataframe für verbindung von scrape cited by to reference  /////////////////////////////)

	df := dataframe.LoadStructs(listpaper)
	df2 := dataframe.LoadStructs(listrelation)
	df5 := dataframe.LoadStructs(relationwos)
	df3 := dataframe.LoadStructs(listpapermaster)
	df4 := dataframe.LoadStructs(prejoinpdf2)
	df4 = df4.Rename("Cites_Nr2", "Cites_Nr")
	//df3 := dataframe.LoadStructs(data.PaperNodes)
	fmt.Println("listpaper")
	fmt.Println(df)
	fmt.Println("listrelation")
	fmt.Println(df2)
	fmt.Println("listpapermaster")
	fmt.Println(df3)
	fmt.Println("prejoinpdf")
	fmt.Println(df4)
	df4 = df4.Rename("Citer_Name", "Cited_by")
	join := df4.InnerJoin(df3, "Citer_Name")
	fmt.Println("JOINS")
	fmt.Println(join)

	//join = join.Rename("Citer_Name", "Cited_by")
	join = join.Rename("Citing_Name", "Name")
	join = join.Rename("Citing_Year", "Year")
	join = join.Rename("Citing_Writer", "Writer")
	//newdf := df2.RBind(join)
	newdf := df5 //.RBind(newdf) //add together sooner as structs
	fmt.Println("newdf")
	fmt.Println(newdf)

	///////////////////////////EINZELNTEST SCHLEIFEN FUER DATAFRAMES AUSTEHEND//////////////////////////////////

	// g := graph.New(graph.IntHash, graph.Directed(), graph.Acyclic())

	// g.AddVertex(1)
	// g.AddVertex(2)
	// g.AddVertex(3)

	// _ = g.AddEdge(1, 2)
	// _ = g.AddEdge(1, 3)

	// filen, _ := os.Create("./mygraph.gv")
	// _ = draw.DOT(g, filen)
	////////////////////////////////////////////////////////////////////////TEMPORARILY addition//////////////////////////////////////

	sorted := newdf.Arrange(
		dataframe.RevSort("Phrases"))

	// fil := sorted.Filter(
	// 	dataframe.F{"Citer_Writer", series.Eq, "Schmid"},
	// )
	// sel1 := sorted.Select([]string{"Citer_Writer"})
	// sel2 := sorted.Select([]string{"Citing_Writer"})

	// fmt.Println("TESTSORT1 ")
	// str := sorted.Elem(1, 8).String()
	// fmt.Println(str)
	// fmt.Println("anotherone6000 ")
	// str2 := sorted.Elem(6000, 8).String()
	// fmt.Println(str2)
	////////////////////////////////////////////////////////////////////////TEMPORARILY//////////////////////////////////////

	paperHash := func(c Paper_str) string {
		//AUFPASSEN FMT INT STRING MIX UMWANDLUNG
		return c.Citer_Conc
	}
	//, graph.Acyclic()
	gr := graph.New(paperHash, graph.Directed(), graph.Tree(), graph.Weighted())

	sortedrec := sorted.Records()
	var vertices_title_add []string
	var conc_add []string
	var conc_sht_add [][]string
	counterspc := 0
	counterspc2 := 0
	verticesnr := 0
	edgesnr := 0
	//var vertices_wriyear_add []string further check fin_phrases
	full_data := Graph{}
	var nds []Node
	var edgs []Edge //
	var phr_add []string
	testc := 0
	for i := 1; i <= sorted.Nrow(); i++ {
		//vertexbuildciter := false
		//vertexbuildciting := false
		//Citer vor Citing weil nachfolgend
		paper_citing := Paper_str{}
		paper_citing.Citer_Name = sortedrec[i][3]
		paper_citing.Citer_Year = sortedrec[i][4]

		// fmt.Println("TESTSORT1")
		// fmt.Println(sortedrec[i][8])
		//var doc_lst []string
		fin_phrases := ""
		if sortedrec[i][8] != "" {
			fin_phrases = sortedrec[i][8]
			// txt_phrases := strings.Split(sortedrec[i][8], ", ")
			// if len(txt_phrases) >= 3 {
			// 	txt_phrases = txt_phrases[:3]
			// } else if len(txt_phrases) == 2 {
			// 	txt_phrases = txt_phrases[:2]
			// } else if len(txt_phrases) == 1 {
			// 	fin_phrases = txt_phrases[0]
			// }

			// if len(txt_phrases) > 0 {
			// 	// n := 0
			// 	// if len(doc_lst) > 3 {
			// 	// 	n = 3
			// 	// } else {
			// 	// 	n = len(doc_lst)
			// 	// }
			// 	fin_phrases = strings.Join(txt_phrases, ", ")

			// 	if !contains(phr_add, fin_phrases) {
			// 		phr_add = append(phr_add, fin_phrases)
			// 	}

			// 	//TEST
			// 	//fin_phrases = fin_phrases + "VORSICHT"
			// 	// fmt.Println("fin_phrases")
			// 	// fmt.Println(fin_phrases)
			// }
		}
		// fmt.Println("0")
		// fmt.Println(sortedrec[i][0])
		// fmt.Println("1")
		// fmt.Println(sortedrec[i][1])
		// fmt.Println("2")
		// fmt.Println(sortedrec[i][2])
		// fmt.Println("3")
		// fmt.Println(sortedrec[i][3])
		// fmt.Println("4")
		// fmt.Println(sortedrec[i][4])
		// fmt.Println("5")
		// fmt.Println(sortedrec[i][5])
		// fmt.Println("6")
		// fmt.Println(sortedrec[i][6])
		// fmt.Println("7")
		// fmt.Println(sortedrec[i][7])
		// fmt.Println("8")
		// fmt.Println(sortedrec[i][8])
		// fmt.Println("9")
		// fmt.Println(sortedrec[i][9])
		// fmt.Println("10")
		// fmt.Println(sortedrec[i][10])
		// fmt.Println("11")

		// fmt.Println(sortedrec[i][11])
		// fmt.Println("12")
		// fmt.Println(sortedrec[i][12])
		// fmt.Println("13")
		// fmt.Println(sortedrec[i][13])

		// if herr != nil {
		// 	fmt.Println("failed2")
		// }
		var writer_raw string
		//var readable_wrt string
		if sortedrec[i][5] != "" && paper_citing.Citer_Year != "0" && len(sortedrec[i][5]) < 100 {
			r, err := charset.NewReader(strings.NewReader(sortedrec[i][5]), "ISO-8859-1")
			if err != nil {
				log.Fatal(err)
			}
			result, err := ioutil.ReadAll(r)
			if err != nil {
				log.Fatal(err)
			}
			paper_citing.Citer_Writer = string(result)
			writer_raw = string(result)
			//paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, ",", " ", -1)
			////////////////////NAMEN VERKUERZEN////////////////////////////
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, ".", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "-", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "&", "and", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "?", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "'", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "/", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "(", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "[", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, ":", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, ")", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "]", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, "$", " ", -1)
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, ",", " ", -1)

			whitespaces := regexp.MustCompile(`\s+`)
			paper_citing.Citer_Writer = whitespaces.ReplaceAllString(paper_citing.Citer_Writer, " ")
			// if strings.Count(paper_citing.Citer_Writer, " ") < 15 {
			paper_citing.Citer_Writer = strings.Trim(paper_citing.Citer_Writer, " ")
			//readable_wrt = strings.Trim(paper_citing.Citer_Writer, " ")
			paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, " ", "_", -1)
			//} else {
			//	paper_citing.Citer_Writer = ""
			//}
			paper_citing.Citer_Year = strings.Trim(paper_citing.Citer_Year, " ")

		} else {
			paper_citing.Citer_Writer = ""
		}

		//}
		// if err != nil {
		// 	panic(err)
		// }
		paper_citing.Cites_Nr, err = strconv.Atoi(sortedrec[i][7])
		//////CHECK FEHLER HIER
		if err != nil {
			panic(err)
		}

		// citingyear := sortedrec[i][4]
		// fmt.Println("citingyear :" + citingyear)
		// citingwriter := sortedrec[i][5]
		// fmt.Println("citingwriter :" + citingwriter)
		//CHECK OB DAS DAS PROBLEM IST
		// tester3, err := strconv.Atoi(sortedrec[i][6])
		// if err != nil {
		// 	panic(err)
		// }
		// paper_citing.Citer_Writer
		paper_citing.Citer_Conc = paper_citing.Citer_Writer + "_" + paper_citing.Citer_Year
		//citerconc2 := readable_wrt + paper_citing.Citer_Year

		sht := regexp.MustCompile(`(,| and |&){1}`)

		writer_sht := paper_citing.Citer_Writer
		conc_sht := paper_citing.Citer_Conc
		var bl bool
		var ind int

		submatchall := sht.FindStringIndex(writer_raw)

		if len(submatchall) > 0 {
			// fmt.Println(",| and |&")
			// fmt.Println(len(submatchall))
			writer_sht = writer_raw[:submatchall[0]]
			writer_sht = trnsform(writer_sht)
			conc_sht = writer_sht + paper_citing.Citer_Year
			// fmt.Println(paper_citing.Citer_Writer)
			// fmt.Println(writer_sht)
			// fmt.Println(conc_sht)

		}

		if paper_citing.Cites_Nr == -100 {

			bl, ind = contains2(conc_sht_add, paper_citing.Citer_Conc)

			if bl {
				// fmt.Println("citer_conc when bl")
				// fmt.Println(paper_citing.Citer_Conc)
				// fmt.Println("conc_sht_add when bl")
				// fmt.Println(conc_sht_add[ind][0])
				// fmt.Println(conc_sht_add[ind][1])
				counterspc = counterspc + 1

			} //&& !bl
			if !contains(vertices_title_add, paper_citing.Citer_Name) && !bl && !contains(conc_add, paper_citing.Citer_Conc) && paper_citing.Citer_Year != "0" && paper_citing.Citer_Writer != "" { //&& paper_citing.Cites_Nr > 5
				//paper_citing.Citer_Conc = conc_sht
				//fmt.Println("citingconc :" + paper_citing.Citer_Conc + " name " + paper_citing.Citer_Name)

				gr.AddVertex(paper_citing)
				verticesnr = verticesnr + 1
				vertices_title_add = append(vertices_title_add, paper_citing.Citer_Name)
				conc_add = append(conc_add, paper_citing.Citer_Conc)
				//var conc_fin = []string{conc_sht, paper_citing.Citer_Conc}
				//conc_sht_add = append(conc_sht_add, conc_fin)
				node_citing := Node{}
				node_citing.Key = paper_citing.Citer_Conc
				attr_citing := Attributes{}
				attr_citing.Label = paper_citing.Citer_Conc
				//attr_citing.Tag = paper_citing.Citer_Name
				attr_citing.Tag = sortedrec[i][11]
				attr_citing.Title = paper_citing.Citer_Name
				attr_citing.WOSID_Wbubble = ""
				if sortedrec[i][9] != "" { //if not needed
					//READ SENTENCE FROM CITATION IN FULL TEXT
					attr_citing.Url = sortedrec[i][9]
				} else {
					attr_citing.Url = ""
				}
				attr_citing.Country = ""
				attr_citing.X = 0.0
				attr_citing.Y = 0.0
				node_citing.Attributes = attr_citing
				nds = append(nds, node_citing)
				//vertexbuildciting = true
			}
		} else {

			//CitnesNR CHECK AUCH HIER ODER AM ENDE ERST ZEICHNEN && tester3 > 0
			if !contains(vertices_title_add, paper_citing.Citer_Name) && !contains(conc_add, paper_citing.Citer_Conc) && paper_citing.Citer_Year != "0" && paper_citing.Citer_Writer != "" { //&& paper_citing.Cites_Nr > 5
				//fmt.Println("citingconc :" + paper_citing.Citer_Conc + " name " + paper_citing.Citer_Name)
				gr.AddVertex(paper_citing)
				verticesnr = verticesnr + 1
				vertices_title_add = append(vertices_title_add, paper_citing.Citer_Name)
				conc_add = append(conc_add, paper_citing.Citer_Conc)
				var conc_fin = []string{conc_sht, paper_citing.Citer_Conc}
				conc_sht_add = append(conc_sht_add, conc_fin)
				node_citing := Node{}
				node_citing.Key = paper_citing.Citer_Conc
				attr_citing := Attributes{}
				attr_citing.Label = paper_citing.Citer_Conc
				//attr_citing.Tag = paper_citing.Citer_Name
				attr_citing.Tag = sortedrec[i][11]

				attr_citing.Title = paper_citing.Citer_Name
				attr_citing.WOSID_Wbubble = ""
				if sortedrec[i][9] != "" {
					//READ SENTENCE FROM CITATION IN FULL TEXT
					attr_citing.Url = sortedrec[i][9]
				} else {
					attr_citing.Url = ""
				}

				attr_citing.Country = ""
				attr_citing.X = 0.0
				attr_citing.Y = 0.0
				node_citing.Attributes = attr_citing
				nds = append(nds, node_citing)
				//vertexbuildciting = true
			}
		}

		paper_citer := Paper_str{}
		paper_citer.Citer_Name = sortedrec[i][0]
		paper_citer.Citer_Year = sortedrec[i][1] //err = strconv.Atoi(sortedrec[i][1])
		// if err != nil {
		// 	panic(err)

		// 	fmt.Println("failed2")

		// }

		//HIER STECKT EIN FEHLER CITER YEAR YEAR 0 HIER HART RAUSFILTERN ERSTMAL; HELFEN SOWIESO NICHT WEITER? oder für referenz; alternative überlegen, eigener graph?
		/////////////////PAUSE
		// if sortedrec[i][2] != "" && paper_citer.Citer_Year != "0" {
		// 	rr, err := charset.NewReader(strings.NewReader(sortedrec[i][2]), "latin1")
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	results, err := ioutil.ReadAll(rr)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	} //CHECK IF NOETIG READER
		var writer_raw2 string
		//var readable_wrt2 string
		if sortedrec[i][2] != "" && paper_citer.Citer_Year != "0" && len(sortedrec[i][2]) < 100 {
			rs, err := charset.NewReader(strings.NewReader(sortedrec[i][2]), "ISO-8859-1")
			if err != nil {
				log.Fatal(err)
			}
			resultt, err := ioutil.ReadAll(rs)
			if err != nil {
				log.Fatal(err)
			}
			paper_citer.Citer_Writer = string(resultt)
			writer_raw2 = string(resultt)
			//paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, ",", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, ".", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "-", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "&", "and", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "?", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "'", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "/", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "(", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "[", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, ")", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "]", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, ":", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, "$", " ", -1)
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, ",", " ", -1)
			//if strings.Count(paper_citer.Citer_Writer, " ") < 15 {

			whitespaces := regexp.MustCompile(`\s+`)
			paper_citer.Citer_Writer = whitespaces.ReplaceAllString(paper_citer.Citer_Writer, " ")
			paper_citer.Citer_Writer = strings.Trim(paper_citer.Citer_Writer, " ")
			//readable_wrt2 = strings.Trim(paper_citer.Citer_Writer, " ")
			paper_citer.Citer_Writer = strings.Replace(paper_citer.Citer_Writer, " ", "_", -1)
			//} else {
			//paper_citer.Citer_Writer = ""
			//}
			paper_citer.Citer_Year = strings.Trim(paper_citer.Citer_Year, " ")

			/////////VERKUERZEN///////// erster name et al, über komma, mindestens 1 dann only erste und et ; bei and kein komma?
		} else {
			paper_citer.Citer_Writer = ""
		}
		// 	paper_citer.Citer_Writer = string(results)
		// } else {
		// 	paper_citer.Citer_Writer = ""
		// }
		//paper_citer.Citer_Writer = sortedrec[i][2]//////////////////////////////////////////

		//paper_citer.Cites_nr = sortedrec[i][6]
		paper_citer.Cites_Nr, err = strconv.Atoi(sortedrec[i][6])
		if err != nil {
			panic(err)
		}
		// citeryear := sortedrec[i][1]
		// fmt.Println("citeryear :" + citeryear)
		// citerwriter := sortedrec[i][2]
		// fmt.Println("citerwriter :" + citerwriter)           paper_citer.Citer_Writer

		paper_citer.Citer_Conc = paper_citer.Citer_Writer + "_" + paper_citer.Citer_Year
		//citerconc := readable_wrt2 + paper_citing.Citer_Year

		var writer_sht2 string
		conc_sht2 := paper_citer.Citer_Conc
		submatch := sht.FindStringIndex(writer_raw2)

		if len(submatch) > 0 {
			// fmt.Println(",| and |&")
			// fmt.Println(len(submatch))
			writer_sht2 = writer_raw2[:submatch[0]]
			writer_sht2 = trnsform(writer_sht2)
			conc_sht2 = writer_sht2 + paper_citer.Citer_Year
			// fmt.Println(paper_citer.Citer_Writer)
			// fmt.Println(writer_sht2)
			// fmt.Println(conc_sht2)

		}

		//&& tester3 > 0
		if !contains(vertices_title_add, paper_citer.Citer_Name) && !contains(conc_add, paper_citer.Citer_Conc) && paper_citer.Citer_Year != "0" && paper_citer.Citer_Writer != "" { //year evt raus && paper_citer.Cites_Nr > 5
			//fmt.Println("citerconc :" + paper_citer.Citer_Conc + " name " + paper_citer.Citer_Name)
			gr.AddVertex(paper_citer)
			verticesnr = verticesnr + 1
			vertices_title_add = append(vertices_title_add, paper_citer.Citer_Name)
			conc_add = append(conc_add, paper_citer.Citer_Conc)
			var conc_fin = []string{conc_sht2, paper_citer.Citer_Conc}
			conc_sht_add = append(conc_sht_add, conc_fin)

			node_citer := Node{}
			node_citer.Key = paper_citer.Citer_Conc
			attr_citer := Attributes{}
			attr_citer.Label = paper_citer.Citer_Conc
			attr_citer.Country = ""
			//attr_citer.Tag = paper_citer.Citer_Name
			attr_citer.Tag = sortedrec[i][10]
			attr_citer.Title = paper_citer.Citer_Name
			attr_citer.WOSID_Wbubble = sortedrec[i][12]
			if attr_citer.WOSID_Wbubble != "" {
				testc = testc + 1
			}

			attr_citer.Url = fin_phrases //+ "VORSICHT" sortedrec[i][8]
			attr_citer.X = 0.0
			attr_citer.Y = 0.0
			node_citer.Attributes = attr_citer
			nds = append(nds, node_citer)

			//vertexbuildciter = true
		}
		//if paper_citer.Cites_Nr > 5 || paper_citing.Cites_Nr > 5 {
		//CHECKEN OB UMWANDLUNG HIER KLAPPT UND ALLE EDGES DA SIND;
		if paper_citing.Cites_Nr == -100 && bl {
			edgesnr = edgesnr + 1
			counterspc2 = counterspc2 + 1
			similar_writer := conc_sht_add[ind][1]
			// fmt.Println("SPECIALCASE:")
			// fmt.Println(edgesnr)
			// fmt.Println(paper_citing.Citer_Conc)
			// fmt.Println(similar_writer)
			// fmt.Println(ind)

			//is_substr := strings.Contains(x, conc_sht)
			_ = gr.AddEdge(similar_writer, paper_citer.Citer_Conc, graph.EdgeWeight(1))

			edge := Edge{}
			edge.Key = paper_citing.Citer_Conc + paper_citer.Citer_Conc
			edge.Source = paper_citing.Citer_Conc
			edge.Target = paper_citer.Citer_Conc
			edge.Undirected = false
			edgs = append(edgs, edge)
		} else {
			edgesnr = edgesnr + 1
			_ = gr.AddEdge(paper_citing.Citer_Conc, paper_citer.Citer_Conc, graph.EdgeWeight(1))

			edge := Edge{}
			edge.Key = paper_citing.Citer_Conc + paper_citer.Citer_Conc
			edge.Source = paper_citing.Citer_Conc
			edge.Target = paper_citer.Citer_Conc
			edge.Undirected = false
			edgs = append(edgs, edge)
		}

		//}

	}

	full_data.Nodes = nds
	full_data.Edges = edgs

	js, err = json.MarshalIndent(full_data, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("full_data_prejoin5.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
	}
	//fmt.Println("verticesnr:" + string(verticesnr))////// SO STRING UMWANELN UND PRINT WIRFT FEHLHERHAFTEN WERT AUFÜASSEN CHECKEN DA OBEN BENUTZT
	//fmt.Println("verticesnr:" + string(edgesnr))
	fmt.Println("testc:")
	fmt.Println(testc)
	//testc
	fmt.Println("verticesnr:")
	fmt.Println(verticesnr)
	fmt.Println("edgesnr:")
	fmt.Println(edgesnr)
	fmt.Println("CHECKVERTICES:")
	fmt.Println(len(nds))
	fmt.Println("CHECKEDGES:")
	fmt.Println(len(edgs))
	fmt.Println("CHECKPHRASES:")
	fmt.Println(len(phr_add))

	for i := 0; i < len(nds); i++ {
		for j := 0; j < len(nds); j++ {
			if i != j && strings.Contains(nds[j].Key, nds[i].Key) {
				fmt.Println("DETECTED")
				fmt.Println(nds[i].Key)

				fmt.Println("IN")
				fmt.Println(nds[j].Key)

				// similarity := strutil.Similarity(nds[i].Key, nds[j].Key, metrics.NewHamming())
				// fmt.Printf("%.2f\n", similarity) // Output: 0.75
			}
		}

	}

	// var listnds []Node2
	// var listedgs []Edge2

	//Load txt reference data (Paper)
	//LOAD FILE WITH ALL RECORDS WHICH ARE ENRICHED WITH DATA ON AN OPTIMIZED LAYOUT AND AND A CLUSTERING
	file2, _ := ioutil.ReadFile("coord_exp05_fin_noo.json") //coord_cluster_noo_exp05_fin_neu
	//data := PaperNodes{}
	var data3 AutoGenerated
	//var jsonIndent []byte
	//var objmap map[string]*json.RawMessage
	//var objmap map[string]interface{}
	//err = json.Unmarshal(file2, &objmap)
	err = json.Unmarshal([]byte(file2), &data3)
	if err != nil {
		log.Fatal(err)
	}
	// jsonIndent, err = json.MarshalIndent(objmap, "", "  ")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//fmt.Println(len(data3))
	// for i := 0; i < len(data3); i++ {
	// 	//are already paperobject could create relation with two empty attributes and add them later;
	// 	// CITEDBY is the pdf name minus .pdf.txt (cut substr) - after manipulation compare to Name in Scraped
	// 	//Create paper table entry and relation table entry later by joining table and relation dataframe;
	// 	//data[i].Cited_by = strings.TrimSuffix(data2[i].Cited_by, ".txt")
	// 	listnds = append(listnds, data3[i].Nodes...)
	// 	listedgs = append(listedgs, data3[i].Edges...)
	nodes := data3.Nodes

	edges := data3.Edges //.(interface{})
	fmt.Println("ATTENTION:")
	// for key, value := range nodes {
	// 	// Each value is an interface{} type, that is type asserted as a string
	//fmt.Println(len(data3))

	fmt.Println("ATTENTION:")
	// }
	// for key, value := range edges {
	// 	// Each value is an interface{} type, that is type asserted as a string
	//fmt.Println(edges)
	//}
	//edgs_str := strings.Split(fmt.Sprintf(edges), "attributes")
	// v := reflect.ValueOf(edges)
	// if v.Kind() == reflect.Map {
	// 	for _, key := range v.MapKeys() {
	// 		fmt.Println("HI7")
	// 		strct := v.MapIndex(key)
	// 		fmt.Println(key.Interface(), strct.Interface())
	// 	}
	// }

	fmt.Println("CHECKVERTICESCOPY:")
	fmt.Println(len(nodes))
	fmt.Println("CHECKEDGESCOPY:")
	fmt.Println(len(edges))
	///////////////////////////////////////////////////FOR TESTS DISABLED

	// LOAD SAVED RECORDS (ALL)
	file3, _ := ioutil.ReadFile("full_data_prejoin5.json")
	//data := PaperNodes{}
	var data4 Graph
	//var jsonIndent []byte
	//var objmap map[string]*json.RawMessage
	//var objmap map[string]interface{}
	//err = json.Unmarshal(file2, &objmap)
	err = json.Unmarshal([]byte(file3), &data4)
	if err != nil {
		log.Fatal(err)
	}
	//var dataset []gomeans.Point
	var d clusters.Observations
	var key_clusters = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29"}

	for i := 0; i < len(data4.Nodes); i++ {
		for j := 0; j < len(nodes); j++ { //////SPAETER HIER WHILE SCHLEIFE
			//rows := df.Subset([]int{0, 2})
			if data4.Nodes[i].Attributes.Label == nodes[j].Label {
				data4.Nodes[i].Attributes.X = nodes[j].X
				data4.Nodes[i].Attributes.Y = nodes[j].Y

				if contains(key_clusters, nodes[j].Attributes.Cluster) {
					data4.Nodes[i].Attributes.Country = nodes[j].Attributes.Cluster
				} else {
					//ALL CLUSTERS WHICH ARE TO SMALL TO HAVE SIGNIFICANT VALUE ARE GROUPED TOGETHER IN A DUMP CLUSTER
					data4.Nodes[i].Attributes.Country = "30"
				}

				// point := gomeans.Point{}
				// point.X = nodes[j].X
				// point.Y = nodes[j].Y
				d = append(d, clusters.Coordinates{nodes[j].X, nodes[j].Y})
				break
			}

		}
		// for i := 0; i <= sorted.Nrow(); i++ {
		// 	if len(sortedrec) > 0 {
		// 		sel2 := "NEW:" + sortedrec[i][1] + "|" + sortedrec[i][2] + "|" + sortedrec[i][4] + "|" + sortedrec[i][5]
		// 		fmt.Println(sel2)
		// 	}

	}

	//cluster := gomeans.Run(dataset, 7)
	// km := kmeans.New()
	// clusters, err := km.Partition(d, 7)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for i := 0; i < len(data4.Nodes); i++ {
	// 	for k := 0; k < len(clusters); k++ {
	// 		for j := 0; j < len(clusters[k].Observations); j++ { //////SPAETER HIER WHILE SCHLEIFE
	// 			//rows := df.Subset([]int{0, 2})
	// 			coord := clusters[k].Observations[j].Coordinates()
	// 			if data4.Nodes[i].Attributes.X == coord[0] && data4.Nodes[i].Attributes.Y == coord[1] {
	// 				data4.Nodes[i].Attributes.Country = fmt.Sprint(k)
	// 				break
	// 			}
	// 			//fmt.Println(clusters[k].Observations[j])
	// 		}
	// 	}
	// 	// for i := 0; i <= sorted.Nrow(); i++ {
	// 	// 	if len(sortedrec) > 0 {
	// 	// 		sel2 := "NEW:" + sortedrec[i][1] + "|" + sortedrec[i][2] + "|" + sortedrec[i][4] + "|" + sortedrec[i][5]
	// 	// 		fmt.Println(sel2)
	// 	// 	}

	// }

	gr3 := Graph3{}
	gr4 := Full_Graph{}

	var final_list_nds []Nde
	var final_list_clstr []Clstr
	var final_list_tags []Tag

	tag := Tag{} // MAKE LOOP OVER LIST OF TAGS
	//{ "key": "unknown", "image": "unknown.svg" }
	tag.Key = "Web Source - Google Scholar"
	tag.Image = "field.svg"
	final_list_tags = append(final_list_tags, tag)

	tag.Key = "Text Source - Full Text (available)"
	tag.Image = "concept.svg"
	final_list_tags = append(final_list_tags, tag)

	tag.Key = "Text Source - Full Text (not available)"
	tag.Image = "unknown.svg"
	final_list_tags = append(final_list_tags, tag)

	tag.Key = "Web Source - Web of Science"
	tag.Image = "person.svg"
	final_list_tags = append(final_list_tags, tag)

	tag.Key = "Web Source - Web of Science (References)"
	tag.Image = "method.svg"
	final_list_tags = append(final_list_tags, tag)

	tag.Key = "Text Source - References (available)"
	tag.Image = "list.svg"
	final_list_tags = append(final_list_tags, tag)

	tag.Key = "Text Source - References (not available)"
	tag.Image = "unknown.svg"
	final_list_tags = append(final_list_tags, tag)

	//imrpove and change to lstedg2 without duplicates

	var final_list_edgs [][]string

	var lstnds []Node
	var lstedg []Edge
	//var cluster_grouped [][]string
	m := make(map[int][]string, 10)
	//var lstedg3 []Edge3
	cnt := 0
	for i := 0; i < len(data4.Nodes); i++ {
		if data4.Nodes[i].Key != "0" && data4.Nodes[i].Attributes.Country != "" { // hier fallen ein paar Sätze raus berechtigt
			data4.Nodes[i].Key = strings.Replace(data4.Nodes[i].Key, "_", " ", -1)
			data4.Nodes[i].Attributes.Label = strings.Replace(data4.Nodes[i].Attributes.Label, "_", " ", -1)

			lstnds = append(lstnds, data4.Nodes[i])
			nde := Nde{}
			nde.Key = strings.Replace(data4.Nodes[i].Key, "_", " ", -1)
			nde.Label = strings.Replace(data4.Nodes[i].Attributes.Label, "_", " ", -1)
			nde.Cluster = data4.Nodes[i].Attributes.Country
			nde.Title = data4.Nodes[i].Attributes.Title
			nde.WOSID_Wbubble = data4.Nodes[i].Attributes.WOSID_Wbubble
			if data4.Nodes[i].Attributes.Tag == "Text Source - References" {
				if data4.Nodes[i].Attributes.Url != "" {
					nde.Tag = "Text Source - References (available)"
				} else {
					nde.Tag = "Text Source - References (not available)"
				}
			} else if data4.Nodes[i].Attributes.Tag == "Text Source - Full Text" {
				if data4.Nodes[i].Attributes.Url != "" {
					nde.Tag = "Text Source - Full Text (available)"
				} else {
					nde.Tag = "Text Source - Full Text (not available)"
				}

			} else {
				nde.Tag = data4.Nodes[i].Attributes.Tag
			}
			//"unknown"data4.Nodes[i].Attributes.Tag
			nde.URL = data4.Nodes[i].Attributes.Url
			nde.Score = 0.0
			nde.X = data4.Nodes[i].Attributes.X
			nde.Y = data4.Nodes[i].Attributes.Y

			final_list_nds = append(final_list_nds, nde)
			cluster_nr := -1
			intVar, err := strconv.Atoi(data4.Nodes[i].Attributes.Country)
			if err == nil {
				cluster_nr = intVar
			}

			if data4.Nodes[i].Attributes.Url != "" && cluster_nr != -1 {
				cnt = cnt + 1
				fmt.Println("ADDED")
				fmt.Println(data4.Nodes[i].Attributes.Country)
				fmt.Println(cluster_nr)
				m[cluster_nr] = append(m[cluster_nr], data4.Nodes[i].Attributes.Url)
			}

		}

	}
	fmt.Println("CNTCHECK")
	fmt.Println(cnt)

	var key_clusterss = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "30"}
	var colors = []string{"#007cff", "#666666", "#57a835", "#7145cd", "#579f5f", "d043c4", "#477028", "#b174cb", "#a4923a", "#5f83cc", "#db4139", "#c94c83", "#7c5d28", "#91ad2b", "#379982", "#a54a49", "#cf7435", "#dc9034", "#8900c1", "#c5b5eb", "#8dcfc7", "#e0de5b", "#DAF7A6", "#af0e24", "#864fa4", "#5398e7", "#fbff00", "#ad771e", "#ef98f7", "#98f7e4", "#6c3e81"}
	//ji := 0
	//for _, li := range key_clusterss {
	for k := 0; k < len(key_clusterss); k++ {
		clstr := Clstr{}
		clstr.Key = key_clusterss[k]
		clstr.ClusterLabel = key_clusterss[k]
		clstr.Color = colors[k]

		if len(m[k]) > 0 {

			clstr_phrases := strings.Join(m[k], ", ")
			clstr_phrases_rnkd := phr_rank(clstr_phrases)

			if len(clstr_phrases_rnkd) > 10 {
				clstr_phrases_rnkd = clstr_phrases_rnkd[:10]
			}
			// }else {
			// 	clstr_phrases_rnkd = clstr_phrases_rnkd
			// }
			//else if len(clstr_phrases_rnkd) == 2 {
			// 	clstr_phrases_rnkd = clstr_phrases_rnkd[:2]
			// } else if len(clstr_phrases_rnkd) == 1 {
			// 	clstr_phrases_rnkd = clstr_phrases_rnkd[:1]
			// }

			clstr_phrases_fin := strings.Join(clstr_phrases_rnkd, ", ")

			clstr.ClusterLabel = clstr_phrases_fin
		}

		final_list_clstr = append(final_list_clstr, clstr)

	}

	for i := 0; i < len(data4.Edges); i++ {
		if data4.Edges[i].Source != "0" && data4.Edges[i].Target != "0" {

			data4.Edges[i].Key = strings.Replace(data4.Edges[i].Key, "_", " ", -1)
			data4.Edges[i].Source = strings.Replace(data4.Edges[i].Source, "_", " ", -1)
			data4.Edges[i].Target = strings.Replace(data4.Edges[i].Target, "_", " ", -1)

			lstedg = append(lstedg, data4.Edges[i])

			var conc_edg = []string{data4.Edges[i].Source, data4.Edges[i].Target}

			final_list_edgs = append(final_list_edgs, conc_edg)
		}

	}

	gr4.Nodes = final_list_nds
	//gr4.Edges = final_list_edgs
	gr4.Clusters = final_list_clstr
	gr4.Tags = final_list_tags

	var lstedg2 []Edge3
	for i := 0; i < len(lstedg); i++ {
		srccheck := false
		tgtcheck := false
		//CHECK IF NODES IN EDGES EXIST
		for j := 0; j < len(final_list_nds); j++ {
			if lstedg[i].Source == final_list_nds[j].Key {
				srccheck = true
			}
			if lstedg[i].Target == final_list_nds[j].Key {
				tgtcheck = true
			}
			if srccheck && tgtcheck {
				break
			}
		}
		if srccheck && tgtcheck {
			dupcheck := false
			for k := 0; k < len(lstedg2); k++ {
				if lstedg2[k].Key == lstedg[i].Key {
					dupcheck = true
				}
			}
			if !dupcheck {
				//type: "arrow"
				edg3 := Edge3{}
				edg3.Key = strings.Replace(lstedg[i].Key, "_", " ", -1)
				edg3.Source = strings.Replace(lstedg[i].Source, "_", " ", -1)
				edg3.Target = strings.Replace(lstedg[i].Target, "_", " ", -1)
				edg3.Attributes.Type = "arrow"
				edg3.Undirected = lstedg[i].Undirected
				lstedg2 = append(lstedg2, edg3)
			}

		}
	}
	fmt.Println("TESTERGR")
	fmt.Println(len(lstnds))
	fmt.Println(len(lstedg2))

	gr3.Nodes = lstnds
	gr3.Edges = lstedg2

	//CALCULATE SCORES BASED ON "CITED BY" RELATIONS
	cited_by_max := 0.0
	for i := 0; i < len(gr4.Nodes); i++ {
		cited_by := 0.0
		for j := 0; j < len(lstedg2); j++ {
			if gr4.Nodes[i].Key == lstedg2[j].Source && cited_by < 100.0 {
				cited_by = cited_by + 1.0
			}

		}
		if cited_by > cited_by_max {
			cited_by_max = cited_by
		}
		gr4.Nodes[i].Score = cited_by
		fmt.Println("TESTSCORE")
		fmt.Println(gr4.Nodes[i].Key)
		fmt.Println(cited_by)
	}

	for i := 0; i < len(gr4.Nodes); i++ {
		//CHECKEN
		gr4.Nodes[i].Score = gr4.Nodes[i].Score / cited_by_max
		// if gr4.Nodes[i].Score > 0.75 {
		// 	gr4.Nodes[i].Score = 0.25
		// } else if gr4.Nodes[i].Score > 0.5 {
		// 	gr4.Nodes[i].Score = 0.175
		// } else if gr4.Nodes[i].Score > 0.25 {
		// 	gr4.Nodes[i].Score = 0.100
		// }

		fmt.Println("TESTSCOREIMPROVED")
		fmt.Println(gr4.Nodes[i].Key)
		fmt.Println(gr4.Nodes[i].Score)
		fmt.Println("TESTSCOREMAX")
		fmt.Println(cited_by_max)

	}

	var final_list_edgs2 [][]string

	for i := 0; i < len(lstedg2); i++ {
		var conc_edgs = []string{lstedg2[i].Source, lstedg2[i].Target}

		final_list_edgs2 = append(final_list_edgs2, conc_edgs)

	}

	gr4.Edges = final_list_edgs2

	js, err = json.MarshalIndent(gr4, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("full_data_demofinv3.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
	}
	//////////////////////////////////////FOR TESTS DISABLED
	//fmt.Println(len(strr))
	//1600 vertices zu erwartem 2040 minus 410 und 1040 references; minus x anteil welcher fehlerhafte writer year or join
	//for _, re := range references {
	fmt.Println("final src")
	fmt.Println(counterspc)
	fmt.Println(counterspc2)

	fmt.Println("TESTERGR")
	fmt.Println(len(lstnds))
	fmt.Println(len(lstedg2))

	//fmt.Println("final src")

	filens, _ := os.Create("./my-graph_finv3.gv")
	_ = draw.DOT(gr, filens)
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////ABHIER
	////////////////////////////////////////////////////////////////////////TEMPORARILY//////////////////////////////////////
	//fmt.Println(df3)
	//dfconc := dataframe.Concat(df)
	//fmt.Println(dfconc)j

	//Transform scraped to paper struct

	//Join Data

	//Graph Visualize Func

	//fmt.Println("TESTUNIT:" + references[0])
	//fmt.Println("TESTUNIT:" + dref[0][0])
	// fmt.Println(dref[1][0])
	// fmt.Println(dref[3][0])
	// fmt.Println(dref[4][0])

	// content, err := ioutil.ReadFile("processdrift.txt")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Convert []byte to string and print to screen
	// textx := string(content)
	// cleanContent := stopwords.CleanString(textx, "en", true)

	// f, err := os.Create("cleanContent.txt")

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer f.Close()

	// _, err2 := f.WriteString(cleanContent)

	// if err2 != nil {
	// 	log.Fatal(err2)
	// }

	// fmt.Println("done")

	// result := make(map[string]int)

	// //start := time.Now()
	// //for _, fn := range flag.Args() {
	// processFile(result, "cleanContent.txt")
	// //}

	// //defer fmt.Printf("Processing took: %v\n", time.Since(start))
	// printResult(result)
}

func phr_rank(rawText string) []string {
	//rawText := "Your long raw text, it could be a book. Lorem ipsum..."
	// TextRank object
	var lst_ranked []string
	//var lst_abbrevations []string

	// rex := regexp.MustCompile(`[(]{1}[A-Z\-\.]+[)]{1}`)

	// submatchall := rex.FindAllString(rawText, -1)
	// submatchallind := rex.FindAllStringIndex(rawText, -1)
	// if len(submatchall) > 0 {
	// 	fmt.Println("Found------------------------------------ABBRV")
	// 	fmt.Println(len(submatchall))
	// 	for _, su := range submatchall {
	// 		fmt.Println(su)
	// 		su = strings.Replace(su, "(", "", 1)
	// 		su = strings.Replace(su, ")", "", 1)
	// 		// fmt.Println("AFTERreplc")
	// 		// fmt.Println(su)
	// 		lst_abbrevations = append(lst_abbrevations, su)

	// 	}
	// }

	tr := textrank.NewTextRank()
	// Default Rule for parsing
	rule := textrank.NewDefaultRule()
	// Default Language for filtering stop words
	language := textrank.NewDefaultLanguage()
	// Default algorithm for ranking text
	algorithmDef := textrank.NewDefaultAlgorithm()

	// Add text.
	tr.Populate(rawText, language, rule)
	// Run the ranking
	tr.Ranking(algorithmDef)

	// Get all phrases by weight
	rankedPhrases := textrank.FindPhrases(tr)

	// Most important phrase
	//fmt.Println(rankedPhrases[0])
	// Second important phrase
	//fmt.Println(rankedPhrases[1])
	iterator := 0
	if len(rankedPhrases) < 100 {
		iterator = len(rankedPhrases)
	} else {
		iterator = 100
	}

	for i := 0; i < iterator; i++ {
		s := rankedPhrases[i].Left + " " + rankedPhrases[i].Right //fmt.Sprint(rankedPhrases[i])

		// for j := 0; j < len(lst_abbrevations); j++ {
		// 	if strings.Contains(s, lst_abbrevations[j]) {
		// 		lastspc := strings.LastIndex(rawText[:submatchallind[j][0]-1], " ")
		// 		s_chg := rawText[lastspc+1 : submatchallind[j][0]-1]
		// 		fmt.Println("BEFORECHG")
		// 		fmt.Println(s)
		// 		fmt.Println("AFTERECHG")
		// 		fmt.Println(s_chg)
		// 		s = s_chg

		// 	}
		// }

		lst_ranked = append(lst_ranked, s)
	}

	return lst_ranked
}

func trnsform(text string) string {
	res := string(text)
	//paper_citing.Citer_Writer = strings.Replace(paper_citing.Citer_Writer, ",", " ", -1)
	////////////////////NAMEN VERKUERZEN////////////////////////////
	res = strings.Replace(res, ".", " ", -1)
	res = strings.Replace(res, "-", " ", -1)
	res = strings.Replace(res, "&", "and", -1)
	res = strings.Replace(res, "?", " ", -1)
	res = strings.Replace(res, "'", " ", -1)
	res = strings.Replace(res, "/", " ", -1)
	res = strings.Replace(res, "(", " ", -1)
	res = strings.Replace(res, "[", " ", -1)
	res = strings.Replace(res, ":", " ", -1)
	res = strings.Replace(res, ")", " ", -1)
	res = strings.Replace(res, "]", " ", -1)
	res = strings.Replace(res, "$", " ", -1)
	res = strings.Replace(res, ",", " ", -1)

	whitespaces := regexp.MustCompile(`\s+`)
	res = whitespaces.ReplaceAllString(res, " ")
	// if strings.Count(paper_citing.Citer_Writer, " ") < 15 {
	res = strings.Trim(res, " ")
	//temporarily
	res = strings.Replace(res, " ", "_", -1)
	//} else {
	//	paper_citing.Citer_Writer = ""
	//}
	res = strings.Trim(res, " ")

	return res
}

func JsonIndent(jsontext []byte) ([]byte, error) {
	var err error
	var jsonIndent []byte
	var objmap map[string]*json.RawMessage
	err = json.Unmarshal(jsontext, &objmap)
	if err != nil {
		return jsonIndent, err
	}
	jsonIndent, err = json.MarshalIndent(objmap, "", "  ")
	return jsonIndent, err
}

// func convrtToUTF8(str string, origEncoding string) string {
// 	strBytes := []byte(str)
// 	byteReader := bytes.NewReader(strBytes)
// 	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
// 	strBytes, _ = ioutil.ReadAll(reader)
// 	return string(strBytes)
// }

func contains(s []string, str string) bool {
	//i:=0
	for _, v := range s {
		//i++
		if v == str {
			return true
		}
	}

	return false
}
func contains2(s [][]string, str string) (bool, int) {
	ji := 0
	bools := false
	counter := 0
	for i := 0; i < len(s); i++ {

		if s[i][0] == str {
			counter = counter + 1
			bools = true
			ji = i
		}
	}

	if counter > 1 {
		//no clear publication can be identified
		bools = false
		ji = 0
	}

	return bools, ji
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func processFile(result map[string]int, fn string) {
	var w string
	r, err := os.Open(fn)
	if nil != err {
		log.Warn(err)
		return
	}
	defer r.Close()

	sc := bufio.NewScanner(r)
	sc.Split(bufio.ScanWords)

	for sc.Scan() {
		w = strings.ToLower(sc.Text())
		result[w] = result[w] + 1
	}
}

func printResult(result map[string]int) {
	fmt.Printf("%-10s%s\n", "Count", "Word")
	fmt.Printf("%-10s%s\n", "-----", "----")

	for w, c := range result {
		fmt.Printf("%-10v%s\n", c, w)
	}
}

func search(docs []string, term string) []string {
	var r []string
	for _, doc := range docs {
		if strings.Contains(doc, term) {
			r = append(r, doc)
		}
	}
	return r
}
