package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"

	//"log"
	"bufio"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	//"fmt"
	//"os"

	//"strings"

	//"golang.org/x/net/html/charset"

	textrank "github.com/DavidBelicza/TextRank"
	log "github.com/sirupsen/logrus"

	"github.com/ledongthuc/pdf" //BEST
	//"github.com/rsc/pdf"

	"github.com/jdkato/prose/v2"
	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"

	//"github.com/go-gota/gota/series"

	_ "github.com/PuerkitoBio/goquery"
	//"github.com/dslipak/pdf"

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

type Pub_sht struct {
	Citer_Writer string `json:"writer"`
	Citer_Name   string `json:"citer_name"`
	WOSID        string `json:"wosid"`
	Citer_Year   int    `json:"year"`
	Cites_Nr     int    `json:"cites_nr"`
	Tag          string `json:"tag"`

	//
	//
	//
	//

	// Cited_by    string   `json:"cited_by"`
	// Year        int      `json:"year"`
	// Name        string   `json:"name"`
	// Writer      string   `json:"writer"`
	// Cites_Nr    int      `json:"cites_nr"`
	// Phrases     string   `json:"phrases"`
	// Phrases_Cit []string `json:"phrases_cit"`

}

type Pub struct {
	PublicationType        string      `json:"Publication Type"`
	Authors                string      `json:"Authors"`
	BookAuthors            string      `json:"Book Authors"`
	BookEditors            string      `json:"Book Editors"`
	BookGroupAuthors       string      `json:"Book Group Authors"`
	AuthorFullNames        string      `json:"Author Full Names"`
	BookAuthorFullNames    string      `json:"Book Author Full Names"`
	GroupAuthors           string      `json:"Group Authors"`
	ArticleTitle           string      `json:"Article Title"`
	SourceTitle            string      `json:"Source Title"`
	BookSeriesTitle        string      `json:"Book Series Title"`
	BookSeriesSubtitle     string      `json:"Book Series Subtitle"`
	Language               string      `json:"Language"`
	DocumentType           string      `json:"Document Type"`
	ConferenceTitle        string      `json:"Conference Title"`
	ConferenceDate         string      `json:"Conference Date"`
	ConferenceLocation     string      `json:"Conference Location"`
	ConferenceSponsor      string      `json:"Conference Sponsor"`
	ConferenceHost         string      `json:"Conference Host"`
	AuthorKeywords         string      `json:"Author Keywords"`
	KeywordsPlus           string      `json:"Keywords Plus"`
	Abstract               string      `json:"Abstract"`
	Addresses              string      `json:"Addresses"`
	ReprintAddresses       string      `json:"Reprint Addresses"`
	EmailAddresses         string      `json:"Email Addresses"`
	ResearcherIds          string      `json:"Researcher Ids"`
	ORCIDs                 string      `json:"ORCIDs"`
	FundingOrgs            string      `json:"Funding Orgs"`
	FundingText            string      `json:"Funding Text"`
	CitedReferences        string      `json:"Cited References"`
	CitedReferenceCount    int         `json:"Cited Reference Count"`
	TimesCitedWoSCore      int         `json:"Times Cited, WoS Core"`
	TimesCitedAllDatabases int         `json:"Times Cited, All Databases"`
	One80DayUsageCount     int         `json:"180 Day Usage Count"`
	Since2013UsageCount    int         `json:"Since 2013 Usage Count"`
	Publisher              string      `json:"Publisher"`
	PublisherCity          string      `json:"Publisher City"`
	PublisherAddress       string      `json:"Publisher Address"`
	Issn                   string      `json:"ISSN"`
	EISSN                  string      `json:"eISSN"`
	Isbn                   string      `json:"ISBN"`
	JournalAbbreviation    string      `json:"Journal Abbreviation"`
	JournalISOAbbreviation string      `json:"Journal ISO Abbreviation"`
	PublicationDate        string      `json:"Publication Date"`
	PublicationYear        int         `json:"Publication Year"`
	Volume                 string      `json:"Volume"`
	Issue                  string      `json:"Issue"`
	PartNumber             interface{} `json:"Part Number"`
	Supplement             interface{} `json:"Supplement"`
	SpecialIssue           string      `json:"Special Issue"`
	MeetingAbstract        string      `json:"Meeting Abstract"`
	StartPage              string      `json:"Start Page"`
	EndPage                string      `json:"End Page"`
	ArticleNumber          string      `json:"Article Number"`
	Doi                    string      `json:"DOI"`
	BookDOI                string      `json:"Book DOI"`
	EarlyAccessDate        string      `json:"Early Access Date"`
	NumberOfPages          int         `json:"Number of Pages"`
	WoSCategories          string      `json:"WoS Categories"`
	ResearchAreas          string      `json:"Research Areas"`
	IDSNumber              string      `json:"IDS Number"`
	UTUniqueWOSID          string      `json:"UT (Unique WOS ID)"`
	PubmedID               interface{} `json:"Pubmed Id"`
	OpenAccessDesignations string      `json:"Open Access Designations"`
	HighlyCitedStatus      string      `json:"Highly Cited Status"`
	HotPaperStatus         string      `json:"Hot Paper Status"`
	DateOfExport           string      `json:"Date of Export"`
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
	Cited_by    string   `json:"cited_by"`
	Year        int      `json:"year"`
	Name        string   `json:"name"`
	Writer      string   `json:"writer"`
	Cites_Nr    int      `json:"cites_nr"`
	Phrases     string   `json:"phrases"`
	Phrases_Cit []string `json:"phrases_cit"`

	// Url      string   `json:"url"`
	// //Topics   []string `json:"topics"`
	// File string `json:"file"`
}

type Paper_str struct {
	// struct for unique paper table; for join between text paper data and scrape paper data
	Citer_Year   string `json:"citer_year"`
	Citer_Name   string `json:"citer_name"`
	Citer_Conc   string `json:"citer_conc"`
	Citer_Writer string `json:"citer_writer"`
	Cites_nr     int    `json:"cites_nr"`
}

type Paper_master struct {
	// struct for unique paper table; for join between text paper data and scrape paper data
	Citer_Year   int    `json:"year"`
	Citer_Name   string `json:"citer_name"`
	Citer_Writer string `json:"writer"`
	Cites_nr     int    `json:"cites_nr"`
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
	Cites_nr      int
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
	Phrases_Cit   []string
	Tag           string
	Tag2          string

	//Cites_Nr      int
	//Topics
}

type Relationv4 struct {
	Citer_Name    string   `json:"citer_name"`
	Citer_Year    int      `json:"citer_year"`
	Citer_Writer  string   `json:"citer_writer"`
	WOSID         string   `json:"wosid"`
	Citing_Name   string   `json:"citing_name"`
	Citing_Year   int      `json:"citing_year"`
	Citing_Writer string   `json:"citing_writer"`
	Cites_Nr      int      `json:"cites_nr"`
	Cites_Nr2     int      `json:"cites_nr2"`
	Phrases       string   `json:"phrases"`
	Phrases_Cit   []string `json:"phrases_cit"`
	Tag           string   `json:"tag"`
	Tag2          string   `json:"tag2"`

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
			var rest string
			//////////////////////////////////////////////ABFRAGEHIER/////////////////// TEST FURTHER- WHICH CASE
			if subsind[len(subsind)-1][1]+1 < len(bef) {
				rest = bef[subsind[len(subsind)-1][1]+1:]
			} else {
				rest = bef[subsind[len(subsind)-1][1]:]
			}

			//rest := bef[subsind[len(subsind)-1][1]+1:]
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
			if len(submatchs) > 0 && len(rest) > 0 && submatchs[len(submatchs)-1][0] > 0 {
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
			if len(submatchs) > 0 && len(rest) > 0 && submatchs[len(submatchs)-1][0] > 0 {
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

			// fmt.Println("HIHIER")
			// fmt.Println(bools2)
			// fmt.Println(rf)

			//re := regexp.MustCompile(`^[a-zA-Z\,\.\s]+[a-zA-Z\"\:\-\s]+[\w\"\.\(\)\-\s]+$`)
			//BREAKER [a-zA-Z\"\s\:\-]?
			// OR CASE MACHEN MIT 2000 jahren und 1000 jahren; or http wegreinigen;
			re := regexp.MustCompile(`(([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\.{1})|(([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\,{1})`)
			notopt := regexp.MustCompile(`([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})`) //BESSER NOCH JAHR MIT ODER MACHEN

			submatchalls := re.FindAllString(rf, -1)
			submatchs := re.FindAllStringIndex(rf, -1)
			//fmt.Println("Testnow")

			//fmt.Println(len(submatchalls))
			if len(submatchalls) == 0 {

				submatchalls = notopt.FindAllString(rf, -1)
				submatchs = notopt.FindAllStringIndex(rf, -1)
			}
			//fmt.Println(len(submatchs))
			// fmt.Println(len(submatchalls))
			if len(submatchalls) > 0 {
				year = submatchalls[0]
			} else {
				year = "0"
			}
			if len(submatchs) > 0 && len(rf) > 0 && submatchs[0][0] > 0 {
				title_pre := rf[0 : submatchs[0][0]-1]
				//fmt.Println("titlepre" + title_pre)

				if submatchs[0][1]+1 < len(rf) {
					//PUNKTGETRENNTE FAELLE WICHTIG NICHT ERSTER PUNKT SONDERN ZWEITER so auch zahl rausbekommen
					title_post := rf[submatchs[0][1]+1:]
					//fmt.Println("check1")
					rx := regexp.MustCompile(`\.`)
					divind := rx.FindAllStringIndex(title_pre, -1)
					divind2 := rx.FindAllStringIndex(title_post, -1)
					if len(divind) > 0 && len(divind2) > 0 {
						//fmt.Println("check2")
						////PUNKT EVT NOCH ENTFERNEN UND ZAHL NOCH ENTFERNEN AUS AUTHOR ODER SCHON VORHER AUS REF FUNC ODER TESTEN WEG
						/// MIT LBREAK ZAHL. LEERSPACE
						//VERSCHIEDENE VARAINTEN ZUSAMMENFÜHREN SPÄTER
						author = title_pre //[0:ind]
						title = title_post[0:divind2[0][0]]
						// fmt.Println("author" + author)
						// fmt.Println("year" + year)
						// fmt.Println("title" + title)
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
	// for _, element := range listref {
	// 	// fmt.Println("dref")
	// 	// fmt.Println(element)
	// 	//fmt.Println("\n")
	// }

	return listref
	//err mit returnen später
}

func analyze_ref2(text string) ([][]string, []string) {
	//count und citref erstmal weggelassen
	//->citref mit Zahl oder autor noch umsetzen hier;
	var listref [][]string
	//var countref []int
	//var citref []string
	var phrref []string

	// var authref []string
	// var yearref []int
	// var titleref []string

	var submatchall [][]int
	var submatchalls []string
	//var helpind int
	//var helpstr string
	//SEARCH
	//start := strings.Index(text, "References")

	_, aftintr, bls1 := strings.Cut(text, "Introduction")
	if !bls1 {
		_, aftintr, bls1 = strings.Cut(text, "INTRODUCTION")
	}
	if bls1 {
		text = aftintr
	}

	preref, aft, bools := strings.Cut(text, "References")
	if !bools {
		preref, aft, bools = strings.Cut(text, "REFERENCES")
	}

	if bools {
		if aft[0:4] == " . ." {
			prereff, aft2, _ := strings.Cut(aft, "References")
			if aft2 != "" {
				aft = aft2
				preref = prereff
			}
		}
		phrref = phr_rank(preref)

		//rex2 := regexp.MustCompile(`([0-9]+\.{1}\n{1})|(DOI\:{1}\n{1})`) /////ODER PUNKTE ENDE LINEBREAK // geht notfalls normal
		r := regexp.MustCompile(`(\({1}([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\){1}\.{1})|(\({1}([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\){1}\,{1})|(\({1}([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\){1}\s{1})`)

		submatchall = r.FindAllStringIndex(aft, -1)
		submatchalls = r.FindAllString(aft, -1)

		fmt.Println("submatchall")
		fmt.Println(len(submatchall))

		if len(submatchall) < 3 { //if year with format (year) is not matched at least 3 times (threshold for false hits), the format is assumed to be year instead
			re := regexp.MustCompile(`(([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\.{1})|(([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\,{1})|(([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})\s{1})`)
			submatchall = re.FindAllStringIndex(aft, -1)
			submatchalls = re.FindAllString(aft, -1)
		}

		if len(submatchall) > 0 {
			//FIRST REF
			firstauth := aft[:submatchall[0][0]]
			firstauth = strings.TrimPrefix(firstauth, " ")
			firstauth = strings.TrimPrefix(firstauth, ".")
			firstauth = strings.TrimSuffix(firstauth, " ")
			firstauth = strings.TrimSuffix(firstauth, ".")
			firstauth = strings.TrimSuffix(firstauth, "(")
			firstauth = strings.Replace(firstauth, "&", ",", 1)
			firstauth = strings.Replace(firstauth, "and", ",", 1)
			firstyear := strings.Replace(submatchalls[0], "(", "", 1)
			firstyear = strings.Replace(firstyear, ")", "", 1)
			firstyear = strings.Replace(firstyear, ".", "", 1)

			//REGEXP AB ZWEITEM TREFFER;
			left := regexp.MustCompile(`[0-9]{2,16}\.{1}`)
			left2 := regexp.MustCompile(`(\w{2,16}\.{1})`)
			right := regexp.MustCompile(`\.{1}`) //left frame of reference
			// title is assumed to be next to the year since its standard format of most journals and conference papers
			newaft := aft[submatchall[0][1]:]
			submatchright := right.FindAllStringIndex(newaft, -1)
			//submatchrights := right.FindAllString(newaft, -1)

			//submatchlefts := left.FindAllString(newaft, -1)
			firsttitle := ""
			if len(submatchright) > 0 {
				// newaft = strings.TrimPrefix(newaft, ".")
				firsttitle = newaft[:submatchright[0][0]]
				if len(firsttitle) < 2 && len(submatchright) > 1 {
					firsttitle = newaft[:submatchright[1][0]]
				}
			}

			firsttitle = strings.TrimPrefix(firsttitle, ")")
			firsttitle = strings.TrimPrefix(firsttitle, ".")
			firsttitle = strings.TrimPrefix(firsttitle, " ")

			fmt.Println("firstauth")
			fmt.Println(firstauth)
			fmt.Println(firstyear)
			fmt.Println(firsttitle)
			//fyear := 0

			// intVar, err := strconv.Atoi(firstyear)
			// if err != nil {
			// 	fyear = 0
			// } else {
			// 	fyear = intVar
			// }
			var firstref []string
			firstauth = author_flip2(firstauth)
			fmt.Println("aftflip")
			fmt.Println(firstauth)
			firstref = append(firstref, firstauth)
			firstref = append(firstref, firstyear)
			firstref = append(firstref, firsttitle)

			listref = append(listref, firstref)

			// it:=1
			// if len(submatchleft)>len(submatchright) {
			// 	it = len(submatchleft)   //-1
			// } else {
			// 	it = len(submatchright)  //-1
			// }

			for i := 1; i < len(submatchall)-1; i++ { //-1
				//ind_start_auth := strings.LastIndex(aft[submatchall[i-1][1]:submatchall[i][0]], ".")
				//newaft := aft[submatchall[0][1]:]
				refstr := aft[submatchall[i-1][1]:submatchall[i][0]]
				submatchleft := left.FindAllStringIndex(refstr, -1)
				submatchleft2 := left2.FindAllStringIndex(refstr, -1)
				newstr := aft[submatchall[i][1]:]
				submatchright := right.FindAllStringIndex(newstr, -1)

				fmt.Println("testright")
				fmt.Println(len(submatchright))
				ref_auth := ""
				ref_year := "0"
				ref_title := ""
				//ryear := 0
				fmt.Println("testleft")
				fmt.Println(len(submatchleft))
				fmt.Println(len(submatchleft2))

				if len(submatchleft) > 0 {
					ref_auth = refstr[submatchleft[len(submatchleft)-1][1]:]
					ref_auth = strings.TrimPrefix(ref_auth, " ")
					ref_auth = strings.TrimPrefix(ref_auth, ".")
					ref_auth = strings.TrimSuffix(ref_auth, " ")
					ref_auth = strings.TrimSuffix(ref_auth, ".")
					ref_auth = strings.TrimSuffix(ref_auth, "(")
					ref_auth = strings.Replace(ref_auth, "&", ",", 1)
					ref_auth = strings.Replace(ref_auth, "and", ",", 1)
					ref_year = strings.Replace(submatchalls[i], "(", "", 1)
					ref_year = strings.Replace(ref_year, ")", "", 1)
					ref_year = strings.Replace(ref_year, ".", "", 1)
					//newstr = strings.TrimPrefix(newstr, ".")
					if len(submatchright) > 0 {
						ref_title = newstr[:submatchright[0][0]]
						if len(ref_title) < 2 && len(submatchright) > 1 {
							ref_title = newstr[:submatchright[1][0]]
						}
					} else {
						space := regexp.MustCompile(`\s{1}`)
						submatchspace := space.FindAllStringIndex(newstr, -1)
						if len(submatchspace) > 0 {
							ref_title = newstr[:submatchspace[0][0]]
						} else {
							ref_title = newstr
						}
						// if i <len(submatchall)-2 {
						// 	ref_title= newstr[:submatchall[i+1][0]]
						// } else {

						//}

					}

					ref_title = strings.TrimPrefix(ref_title, ")")
					ref_title = strings.TrimPrefix(ref_title, ".")
					ref_title = strings.TrimPrefix(ref_title, " ")
				} else if len(submatchleft2) > 0 {
					ref_auth = refstr[submatchleft2[len(submatchleft2)-1][1]:]
					ref_auth = strings.TrimPrefix(ref_auth, " ")
					ref_auth = strings.TrimPrefix(ref_auth, ".")
					ref_auth = strings.TrimSuffix(ref_auth, " ")
					ref_auth = strings.TrimSuffix(ref_auth, ".")
					ref_auth = strings.TrimSuffix(ref_auth, "(")
					ref_auth = strings.Replace(ref_auth, "&", ",", 1)
					ref_auth = strings.Replace(ref_auth, "and", ",", 1)
					ref_year = strings.Replace(submatchalls[i], "(", "", 1)
					ref_year = strings.Replace(ref_year, ")", "", 1)
					ref_year = strings.Replace(ref_year, ".", "", 1)
					//newstr = strings.TrimPrefix(newstr, ".")
					if len(submatchright) > 0 {
						ref_title = newstr[:submatchright[0][0]]
						if len(ref_title) < 2 && len(submatchright) > 1 {
							ref_title = newstr[:submatchright[1][0]]
						}
					} else {
						space := regexp.MustCompile(`\s{1}`)
						submatchspace := space.FindAllStringIndex(newstr, -1)
						if len(submatchspace) > 0 {
							ref_title = newstr[:submatchspace[0][0]]
						} else {
							ref_title = newstr
						}
					}
					ref_title = strings.TrimPrefix(ref_title, ")")
					ref_title = strings.TrimPrefix(ref_title, ".")
					ref_title = strings.TrimPrefix(ref_title, " ")
				}

				// intVar, err := strconv.Atoi(ref_year)
				// if err != nil {
				// 	ryear = 0
				// } else {
				// 	ryear = intVar
				// }
				fmt.Println("testauth")
				fmt.Println(ref_auth)
				fmt.Println(ref_year)
				fmt.Println(ref_year)
				fmt.Println(ref_title)

				ref_auth = author_flip2(ref_auth)
				fmt.Println("aftflip")
				fmt.Println(ref_auth)

				var ref []string

				ref = append(ref, ref_auth)
				ref = append(ref, ref_year)
				ref = append(ref, ref_title)

				listref = append(listref, ref)

				// authref = append(authref, ref_auth)
				// yearref = append(yearref, ref_year)
				// titleref = append(titleref, ref_title)
			}
		}
	} else {
		phrref = phr_rank(text)
	}

	return listref, phrref //authref, yearref, titleref
	//err mit returnen später
}

func analyze_ref(text string) ([]string, []int, []string, []string) {
	//cited_papers := make([]Paper, 0)
	//s = make([]byte, 5, 5)
	//var ref []string
	var listref []string
	var countref []int
	var citref []string
	//var phrref []string
	phrref := phr_rank(text)
	//phrases

	//var solution := "" IN REF MUSS MAN DAS ÜBER APPEND MACHEN DURCH DIE VIELEN EINZELNEN FÄLLE
	var submatchall [][]int
	var submatchalls []string
	//var helpind int
	//var helpstr string
	//SEARCH
	start := strings.Index(text, "References")

	_, aftintr, bls1 := strings.Cut(text, "Introduction")
	if !bls1 {
		_, aftintr, bls1 = strings.Cut(text, "INTRODUCTION")
	}
	if bls1 {
		text = aftintr
	}
	fmt.Println(start)
	///evt noch LINEBREAK EINBAUEN
	preref, aft, bools := strings.Cut(text, "References")
	if !bools {
		preref, aft, bools = strings.Cut(text, "REFERENCES")
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

		//ZAHLENCOUNTER NOCH MACHEN UND DANN ALLES ZUSAMMEZIEHEN;
		citnrref := re.FindAllString(str1, -1)

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

					counter_ref := 0
					searchstr := citnrref[i]
					nrselect := regexp.MustCompile(`[0-9]+`)
					//slternativ cites übwe (namen finden,mit punkt , )

					nrmatch := nrselect.FindAllString(searchstr, -1)
					//SEARCH TEXT FOR CITATION WITH [NR]
					counter_ref = count_ref(preref, nrmatch[0])

					countref = append(countref, counter_ref)

					citref = append(citref, nrmatch[0])
				} else {
					// fmt.Println("HI1")
					// fmt.Println(aft[sta:end]) TESTER///////////////////////////
					txt := aft[sta:end]
					// fmt.Println("HI1HI") //////////TESTER
					rx := regexp.MustCompile(`\.{1}\n{1}`)
					submatchs := rx.FindAllStringIndex(txt, -1)
					// fmt.Println("subs:")
					// fmt.Println(len(submatchs)) TESTER///////////

					if len(submatchs) > 0 {
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

							counter_ref := 0
							searchstr := citnrref[i]
							nrselect := regexp.MustCompile(`[0-9]+`)
							//slternativ cites übwe (namen finden,mit punkt , )

							nrmatch := nrselect.FindAllString(searchstr, -1)
							//SEARCH TEXT FOR CITATION WITH [NR]
							counter_ref = count_ref(preref, nrmatch[0])

							countref = append(countref, counter_ref)

							citref = append(citref, nrmatch[0])

						}
					}

				}
				//list add References

			}
			if len(submatchall) > 0 && len(aft) > 0 {
				reff := Reference{}
				sta := submatchall[len(submatchall)-1][1]
				reff.Startind = submatchall[len(submatchall)-1][1]
				end := len(aft) - 1
				reff.Endind = len(aft) - 1
				//ab jetzt betrifft das ja nur den letzten aft[sta:end] block wenn früher über 1000 length nicht betreten - generelisieren
				if len(aft[sta:end]) < 1000 {

					listref = append(listref, aft[sta:end])
					//SAVE LAST ELEMENT OF FOUND STRINGS REFERENCES ([NR])
					counter_ref := 0
					searchstr := citnrref[len(submatchall)-1]
					nrselect := regexp.MustCompile(`[0-9]+`)
					//slternativ cites übwe (namen finden,mit punkt , )

					nrmatch := nrselect.FindAllString(searchstr, -1)
					//SEARCH TEXT FOR CITATION WITH [NR]
					counter_ref = count_ref(preref, nrmatch[0])

					countref = append(countref, counter_ref)

					citref = append(citref, nrmatch[0])

				} else {
					// fmt.Println("HI2")
					// fmt.Println(aft[sta:end]) TESTER //////////////
					txt := aft[sta:end]
					//fmt.Println("HI1HI")
					rx := regexp.MustCompile(`\.{1}\n{1}`)
					submatchs := rx.FindAllStringIndex(txt, -1)

					//////////////////////NUR WEITER IF SUBMATCHS LEN GR?ER 0/////// und dann auch counterref exakt gleich
					//fmt.Println("subs:")
					//fmt.Println(len(submatchs))

					if len(submatchs) > 0 {
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

							//SAVE LAST ELEMENT OF FOUND STRINGS REFERENCES ([NR])
							counter_ref := 0
							searchstr := citnrref[len(submatchall)-1]
							nrselect := regexp.MustCompile(`[0-9]+`)
							//slternativ cites übwe (namen finden,mit punkt , )

							nrmatch := nrselect.FindAllString(searchstr, -1)
							//SEARCH TEXT FOR CITATION WITH [NR]
							counter_ref = count_ref(preref, nrmatch[0])

							countref = append(countref, counter_ref)

							citref = append(citref, nrmatch[0])

						}

					}

				}

			}

			//SAVE COUNT
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
			//fmt.Println("EXTRACASE")
			//HIER NOCH WEITER SPEZIFIEREN WEITERE CASES
			//fmt.Println(len(submatchalls))
			//fmt.Println(submatchalls[0]) HIIIER
			//fmt.Println("EXTRACASE2")
			//fmt.Println(len(submatchalls2))
			//fmt.Println(len(submatchalls2))
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

				counter_ref := 0

				nrselect := regexp.MustCompile(`[0-9]+`)
				//slternativ cites übwe (namen finden,mit punkt , )

				nrmatch := nrselect.FindAllString(submatchalls[i], -1)
				//citesmatch durchgehen und häufigkeit zahl finden searchstr
				//ZAHLENCOUNTER NOCH MACHEN UND DANN ALLES ZUSAMMEZIEHEN;
				searchstr := nrmatch[0] // "[" + + "]"
				//submatchalls bis zum punkt oder bis zum leerzeichen und dann drum herum klammer,

				//SEARCH TEXT FOR CITATION WITH [NR]
				counter_ref = count_ref(preref, searchstr)
				countref = append(countref, counter_ref)

				citref = append(citref, searchstr)

			}
			if len(submatchall) > 0 { ///len(aft) > 0
				reff := Reference{}
				sta := submatchall[len(submatchall)-1][1] - 1
				reff.Startind = submatchall[len(submatchall)-1][1] - 1
				end := len(aft) - 1
				reff.Endind = len(aft) - 1
				// fmt.Println("tester:" + aft[sta:end])
				listref = append(listref, aft[sta:end])

				counter_ref := 0

				nrselect := regexp.MustCompile(`[0-9]+`)
				//slternativ cites übwe (namen finden,mit punkt , )

				nrmatch := nrselect.FindAllString(submatchalls[len(submatchalls)-1], -1)
				//citesmatch durchgehen und häufigkeit zahl finden searchstr
				//ZAHLENCOUNTER NOCH MACHEN UND DANN ALLES ZUSAMMEZIEHEN;
				searchstr := nrmatch[0] // "[" + + "]"
				//submatchalls bis zum punkt oder bis zum leerzeichen und dann drum herum klammer,

				//SEARCH TEXT FOR CITATION WITH [NR]
				counter_ref = count_ref(preref, searchstr)
				countref = append(countref, counter_ref)

				citref = append(citref, searchstr)

			}

			// fmt.Println("strtreffer", "\n")
			// fmt.Println(len(submatchalls))
			// for _, element := range submatchalls {

			// 	fmt.Println(element, "\n")
			// 	fmt.Println("\n")
			// }
		}

		//end:= strings.Index(text, "References")
	}
	// for _, element := range listref {
	// 	fmt.Println("ref")
	// 	fmt.Println(element, "\n")
	// 	fmt.Println("\n")
	// }
	//CUT

	//DIVIDE

	//SAVE

	return listref, countref, phrref, citref
	//err mit returnen später
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

func count_ref2(text string, ref string) int {
	preref, _, bools := strings.Cut(text, "References")
	if !bools {
		preref, _, bools = strings.Cut(text, "REFERENCES")
	}
	if bools {
		text = preref
	}
	count := 0
	//CREATES LIST OF ALL CITATIONS AND COUNTS CITATIONS OF PAPER X
	strf := fmt.Sprintf(`\[{1}[0-9\,\s]*%s{1}[0-9\,\s]*\]{1}`, regexp.QuoteMeta(ref)) //meta uberfl?//notfalls gleich meta in die mitte und jeweils zwei klammer darum
	// fmt.Println("strf")
	// fmt.Println(strf)
	// fmt.Println("ref")
	// fmt.Println(ref)
	rex := regexp.MustCompile(strf) /////ODER PUNKTE ENDE LINEBREAK // geht notfalls normal
	//submatchall = rex.FindAllStringIndex(text, -1)
	submatchs := rex.FindAllString(text, -1)
	submatchs_ind := rex.FindAllStringIndex(text, -1)
	fmt.Println(len(submatchs))
	for i := 0; i < len(submatchs); i++ {
		refstr := fmt.Sprintf(`%s{1}`, regexp.QuoteMeta(ref))
		rext := regexp.MustCompile(refstr)
		refstr_ind := rext.FindAllStringIndex(submatchs[i], -1)
		newtext := submatchs[i]
		for l := 0; l < len(refstr_ind); l++ {
			if submatchs_ind[i][1]+1 < len(text) && strings.Contains(submatchs[i], ref) {
				fmt.Println("TESTSUBMATCHS ")
				fmt.Println(submatchs[i])
				fmt.Println(len(newtext))
				//fmt.Println(conc_cites[submatchs_ind[i][1] : submatchs_ind[i][1]+1])
				// fmt.Println("TESTCOUNTSSUBSTR2: ")
				// fmt.Println(conc_cites[submatchs_ind[i][1]+1 : submatchs_ind[i][1]+2]) //:submatchs_ind[i][1]+1
				if (string(newtext[refstr_ind[l][1]]) == "," || string(newtext[refstr_ind[l][1]]) == "]") && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {

					count = count + 1

				}
			} else if strings.Contains(submatchs[i], ref) && refstr_ind[l][1]+1 == len(text) && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {
				count = count + 1

			}
		}
	}
	fmt.Println("count")
	fmt.Println(count)
	return count
}

func count_ref(text string, ref string) int {

	//CREATES LIST OF ALL CITATIONS AND COUNTS CITATIONS OF PAPER X
	nrselect2 := regexp.MustCompile(`\[{1}[0-9\,\s]*\]{1}`)
	citesmatchlist := nrselect2.FindAllString(text, -1)

	// for _, el := range citesmatchlist {
	// 	fmt.Println("CHECK IF RIGHT READ CITES  ")
	// 	fmt.Println(el)
	// }
	conc_cites := strings.Join(citesmatchlist, "")
	count := 0
	//str := "ref+`\,|\]{1}`"
	// count = strings.Count(conc_cites, ref)
	rex := regexp.MustCompile(regexp.QuoteMeta(ref)) /////ODER PUNKTE ENDE LINEBREAK // geht notfalls normal
	//submatchall = rex.FindAllStringIndex(text, -1)
	submatchs := rex.FindAllString(conc_cites, -1)
	submatchs_ind := rex.FindAllStringIndex(conc_cites, -1)
	fmt.Println(len(submatchs))
	for i := 0; i < len(submatchs); i++ {
		if submatchs_ind[i][1]+1 < len(conc_cites) {

			//fmt.Println("TESTCOUNTSSUBSTR: ")
			//fmt.Println(conc_cites[submatchs_ind[i][1] : submatchs_ind[i][1]+1])
			// fmt.Println("TESTCOUNTSSUBSTR2: ")
			// fmt.Println(conc_cites[submatchs_ind[i][1]+1 : submatchs_ind[i][1]+2]) //:submatchs_ind[i][1]+1
			if (conc_cites[submatchs_ind[i][1]:submatchs_ind[i][1]+1] == "," || conc_cites[submatchs_ind[i][1]:submatchs_ind[i][1]+1] == "]") && (conc_cites[submatchs_ind[i][0]-1:submatchs_ind[i][0]] == "," || conc_cites[submatchs_ind[i][0]-1:submatchs_ind[i][0]] == "[" || conc_cites[submatchs_ind[i][0]-1:submatchs_ind[i][0]] == " ") {
				count = count + 1
			}
		} else if submatchs_ind[i][1]+1 == len(conc_cites) && (conc_cites[submatchs_ind[i][0]-1:submatchs_ind[i][0]] == "," || conc_cites[submatchs_ind[i][0]-1:submatchs_ind[i][0]] == "[" || conc_cites[submatchs_ind[i][0]-1:submatchs_ind[i][0]] == " ") {
			count = count + 1
		}

	}
	//count = len(submatchs)
	// fmt.Println("CHECK COUNT  ")
	// fmt.Println(count)
	return count
}

func author_flip(author string) string {
	var authorslst []string
	var authorslstfinal []string

	authors := ""

	if len(author) > 0 {
		authorslst = strings.Split(author, ";")
	}

	for i := 0; i < len(authorslst); i++ {
		res := ""
		auth := strings.Trim(authorslst[i], " ")
		auth = strings.Replace(auth, " ", "", -1)
		auth = strings.Replace(auth, ",", " ", -1)
		rex := regexp.MustCompile(`\s{1}`)

		spc := rex.FindStringIndex(auth)
		if len(spc) > 0 {
			name := auth[:spc[0]]
			surname := auth[spc[0]+1:]
			res = surname + " " + name
			fmt.Print("alt:")
			fmt.Println(author)
			fmt.Print("neu:")
			fmt.Println(res)
		}
		authorslstfinal = append(authorslstfinal, res)
	}
	authors = strings.Join(authorslstfinal, ", ")

	return authors
}

func author_flip2(author string) string {
	var authorslst []string
	var authorslstfinal []string

	authors := ""

	if strings.Contains(author, ";") {
		if len(author) > 0 {
			authorslst = strings.Split(author, ";")
		}

		for i := 0; i < len(authorslst); i++ {
			res := ""
			auth := strings.Trim(authorslst[i], " ")
			auth = strings.Replace(auth, " ", "", -1)
			auth = strings.Replace(auth, ",", " ", -1)
			rex := regexp.MustCompile(`\s{1}`)

			spc := rex.FindStringIndex(auth)
			if len(spc) > 0 {
				name := auth[:spc[0]]
				surname := auth[spc[0]+1:]
				res = surname + " " + name
				fmt.Print("alt:")
				fmt.Println(author)
				fmt.Print("neu:")
				fmt.Println(res)
			}
			authorslstfinal = append(authorslstfinal, res)
		}
	} else {
		authorslst2 := strings.Split(author, ".,")
		if len(authorslst2) > 2 {
			for i := 0; i < len(authorslst2); i++ {
				res := ""
				auth := strings.Trim(authorslst2[i], " ")
				auth = strings.Replace(auth, " ", "", -1)
				auth = strings.Replace(auth, ",", " ", -1)
				rex := regexp.MustCompile(`\s{1}`)
				// rex := regexp.MustCompile(`\.{1}`)

				spc := rex.FindStringIndex(auth)
				if len(spc) > 0 {
					name := auth[:spc[0]]
					surname := auth[spc[0]+1:]
					res = surname + " " + name
					fmt.Print("alt:")
					fmt.Println(author)
					fmt.Print("neu:")
					fmt.Println(res)
				}
				authorslstfinal = append(authorslstfinal, res)
			}
		} else {
			preref, aft, bools := strings.Cut(author, ".,")
			fmt.Println(bools)
			if bools {
				fmt.Println("preref")
				fmt.Println(preref)
				fmt.Println("aft")
				fmt.Println(aft)

				res := ""
				auth := strings.Trim(preref, " ")
				auth = strings.Replace(auth, " ", "", -1)
				auth = strings.Replace(auth, ",", " ", -1)
				rex := regexp.MustCompile(`\s{1}`)
				//rex := regexp.MustCompile(`\.{1}`)

				spc := rex.FindStringIndex(auth)
				if len(spc) > 0 {
					name := auth[:spc[0]]
					surname := auth[spc[0]+1:]
					res = surname + " " + name
					fmt.Print("alt:")
					fmt.Println(author)
					fmt.Print("neu:")
					fmt.Println(res)
				}

				authorslstfinal = append(authorslstfinal, res)
				rest := ""
				if strings.Contains(aft, ",") {

					auth := strings.Trim(aft, " ")
					auth = strings.TrimPrefix(auth, ",")
					auth = strings.Replace(auth, " ", "", -1)
					auth = strings.Replace(auth, ",", " ", -1)
					rex := regexp.MustCompile(`\s{1}`)
					//rex := regexp.MustCompile(`\.{1}`)

					spc := rex.FindStringIndex(auth)
					if len(spc) > 0 {
						name := auth[:spc[0]]
						surname := auth[spc[0]+1:]
						rest = surname + " " + name
						fmt.Print("alt:")
						fmt.Println(author)
						fmt.Print("neu:")
						fmt.Println(rest)
					}
				}

				authorslstfinal = append(authorslstfinal, rest)
			} else {
				res := ""
				auth := strings.Trim(author, " ")
				auth = strings.Replace(auth, " ", "", -1)
				auth = strings.Replace(auth, ",", " ", -1)
				rex := regexp.MustCompile(`\s{1}`)
				//rex := regexp.MustCompile(`\.{1}`)

				spc := rex.FindStringIndex(auth)
				if len(spc) > 0 {
					name := auth[:spc[0]]
					surname := auth[spc[0]+1:]
					res = surname + " " + name
					fmt.Print("alt:")
					fmt.Println(author)
					fmt.Print("neu:")
					fmt.Println(res)
				}

				authorslstfinal = append(authorslstfinal, res)
			}

		}

		// // if len(aft)>0 {
		// authorslstfinal = append(authorslstfinal, aft)
		// 		authorslst = strings.Split(aft, ",")

		// }
	}

	authors = strings.Join(authorslstfinal, ", ")

	return authors
}

func main() {
	//////////////////////////////////////////////////////////WIEDER NUTZEN UNIDOC
	// err, res := outputPdfText("pdfs/489426276.pdf")
	// if err != nil {
	// 	log.Fatal(err)

	//READ CSV TO STRUCT

	// pubsFile, err := os.OpenFile("webofscience.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }
	// defer pubsFile.Close()

	// pubs := []*AutoGenerated{}

	// if err := gocsv.UnmarshalFile(pubsFile, &pubs); err != nil { // Load clients from file
	// 	panic(err)
	// }

	listpdf, err := ioutil.ReadDir("C:/Users/neumu/pdffer3_wos/paper_fulltext")
	if err != nil {
		panic(err)
	}
	var lsttxt []string
	for _, li := range listpdf {
		// content, err := ioutil.ReadFile("paper_fulltext/" + li.Name())
		// if err != nil {
		// 	log.Fatal(err)
		// }
		fmt.Println(li.Name())
		txtfile := strings.TrimSuffix(li.Name(), ".txt")
		fmt.Println("ADDTHIS")
		fmt.Println(txtfile)
		lsttxt = append(lsttxt, txtfile)

	}

	// fmt.Println("TESTER")
	// fmt.Println(len(pubs))
	// for _, pub := range pubs[:5] {
	// 	fmt.Println(pub.Authors)
	// 	fmt.Println(pub.UTUniqueWOSID)
	// } ////////////////////////////////////////////////////////////////////////////////_____________________________
	var publist []Pub_sht

	file, _ := ioutil.ReadFile("webofscience.json")
	//data := PaperNodes{}
	var data []Pub
	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(data))
	for i := 0; i < len(data[:20]); i++ {
		id := strings.TrimPrefix(data[i].UTUniqueWOSID, "WOS:")
		if contains(lsttxt, id) {
			pub := Pub_sht{}
			pub.Citer_Writer = author_flip(data[i].Authors)
			fmt.Println("BEFFFORE")
			fmt.Println(data[i].Authors)
			fmt.Println("AFFFTER")
			fmt.Println(pub.Citer_Writer)
			pub.Citer_Name = data[i].ArticleTitle
			pub.Citer_Year = data[i].PublicationYear
			pub.Cites_Nr = data[i].TimesCitedWoSCore
			pub.Tag = "Text Source - Full Text" //spaeter nochmal splitten nach available and not;
			pub.WOSID = strings.TrimPrefix(data[i].UTUniqueWOSID, "WOS:")
			fmt.Println("REST")
			fmt.Println(pub.Citer_Name)
			fmt.Println(pub.Citer_Year)
			fmt.Println(pub.Cites_Nr)
			fmt.Println(pub.Tag)
			fmt.Println(pub.WOSID)
			publist = append(publist, pub)
		}

		//fmt.Println(data[i].UTUniqueWOSID)

	}
	// csvContent, err := gocsv.MarshalString(&pubs) // Get all clients as CSV string
	// //err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(csvContent) // Display all clients as CSV string
	/////////////////////////aktuell auskommentiert
	// } //////////////////////////////////////////////////////////////////////////HIER
	cited_fin := make([]Relationv4, 0)
	k := 0

	// // listpdf, err := ioutil.ReadDir("C:/Users/neumu/pdffer3_wos/paper_fulltext")
	// // if err != nil {
	// // 	panic(err)
	// // }
	// //var def []map[int]string
	// //def := make([]map[int]string, 0)
	// //"C:/Users/neumu/pdffer3_wos/paper_fulltext"
	// //for _, li := range listpdf {
	for _, li := range publist[:20] {

		// 	// 	///////////////////////////////////////////////////////////////////////TEMPORARY PASSIVE////////////////////////

		// 	// 	// do something with the article
		content, err := ioutil.ReadFile("paper_fulltext/" + li.WOSID + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("paper_fulltext/" + li.WOSID + ".txt")
		text := string(content)
		// 	fmt.Println(text)

		// 	// Convert []byte to string and print to screen
		// 	//text := string(content)
		refs, phrasess := analyze_ref2(text)
		fmt.Println("TESTERR")
		if len(refs) > 0 {
			fmt.Println(refs[0][0])
			fmt.Println(refs[0][1])
			fmt.Println(refs[0][2])

			fmt.Println(phrasess[0])
			phrases_c := strings.Join(phrasess, ", ")
			for i := 0; i < len(refs); i++ {
				relation := Relationv4{}
				relation.Citer_Name = li.Citer_Name
				relation.Citer_Writer = li.Citer_Writer
				relation.Citer_Year = li.Citer_Year
				relation.Cites_Nr = li.Cites_Nr
				relation.WOSID = li.WOSID
				relation.Phrases = phrases_c
				relation.Tag = li.Tag
				//item.Cited_by = li.Name()
				relation.Citing_Writer = refs[i][0]
				relation.Citing_Name = refs[i][2]
				//item.Cites_Nr
				intVar, err := strconv.Atoi(refs[i][1])
				if err != nil {
					relation.Citing_Year = 0
				} else {
					relation.Citing_Year = intVar
				}
				k = k + 1
				//EMPTYSPACE ALSO POSSIBLE BUT THEN MUST CLEAN ALL EMPTY SPACE IN BEFORE
				// writer_mod := strings.Trim(ddref[0], " ")
				// writer_bef, _, status := strings.Cut(writer_mod, " ")
				// if !status {
				// 	writer_mod = writer_bef
				// }
				//fmt.Println("WRITERMOD" + writer_mod)
				//Count_citations := count_ref(text_re, writer_mod)  //////////////// ALTERNATIVE BERECHNUNG MIT AUTOR FUNKTION BAUEN WENN CITES 0
				relation.Cites_Nr2 = 0 ///// CHECK IF K IS ALWAYS RIGHT
				relation.Tag2 = "Text Source - References"
				//newciter := count_ref2(text, cites[k])
				// fmt.Println("NEW")
				// fmt.Println(newciter)
				// fmt.Println("VS")
				// fmt.Println("OLD")
				// fmt.Println(counts[k])
				//relation.Cites_Nr2 = newciter
				relation.Phrases = phrases_c
				auth_first := relation.Citing_Writer

				//Cut out first author
				writer_bef, _, status := strings.Cut(relation.Citing_Writer, ",")
				if status {
					writer_bef, aft, status := strings.Cut(writer_bef, " ")
					if status {
						auth_first = aft
					} else {
						auth_first = writer_bef
					}
				}

				relation.Phrases_Cit = get_citref2(text, auth_first, refs[i][1])

				fmt.Println("PHRASES_TEST")
				fmt.Println(relation.Phrases_Cit)
				fmt.Println(len(relation.Phrases_Cit))

				//ADD YEAR TO SEARCH FOR CITATION; NOCH WEITER TESTEN;
				// 	writer_bef,_,status = strings.Cut(ddref[0], "and")
				// 	if !status {
				// 		writer_bef,_,status = strings.Cut(ddref[0], "&")
				// 	}

				// }

				// do something with the article
				cited_fin = append(cited_fin, relation)
			}
		} else {
			relation := Relationv4{}
			relation.Citer_Name = li.Citer_Name
			relation.Citer_Writer = li.Citer_Writer
			relation.Citer_Year = li.Citer_Year
			relation.Cites_Nr = li.Cites_Nr
			relation.WOSID = li.WOSID
			relation.Phrases = ""
			relation.Tag = li.Tag
			//item.Cited_by = li.Name()
			relation.Citing_Writer = ""
			relation.Citing_Name = ""
			relation.Cites_Nr2 = 0
			relation.Citing_Year = 0
			relation.Tag2 = ""
			//empty := make([]string, 0)
			var empty []string
			relation.Phrases_Cit = empty
			cited_fin = append(cited_fin, relation)

		}

	} // } ////////////////////////////////////////////////////////////////////////////////_____________________________

	// 	fmt.Println("COUNTSTEST")
	// 	fmt.Println(len(counts))
	// 	fmt.Println(len(references))
	// 	fmt.Println("CITESSTEST")
	// 	fmt.Println(len(cites))

	// 	phrases_c := strings.Join(phrases, ", ")
	// 	// fmt.Println("PHRASES_TEST")
	// 	// fmt.Println(phrases_c)

	// 	// content, err := ioutil.ReadFile("C:/Users/neumu/pdffer/pdfs/10-1108_BPMJ-10-2021-0677.pdf")
	// 	// if err != nil {
	// 	// 	log.Fatal(err)
	// 	// }
	// 	//fmt.Println(counts[0])
	// 	cited_papers := make([]Relationv3, 0)
	// 	// fmt.Println(content)
	// 	//text_re := text
	// 	k := 0
	// 	for _, re := range references { // TEST (http|https|ftp): WEGLASSEN
	// 		rex := regexp.MustCompile(`(http|https|ftp):[\/]{2}([a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,4})(:[0-9]+)?\/?([a-zA-Z0-9\-\._\?\,\'\/\\\+&amp;%\$#\=~]*)`)

	// 		submatchall := rex.FindAllString(re, -1)
	// 		if len(submatchall) > 0 {
	// 			//fmt.Println("Foundyeyy")
	// 			//fmt.Println(len(submatchall))
	// 			for _, su := range submatchall {
	// 				//fmt.Println(su)
	// 				re = strings.Replace(re, su, "", -1)
	// 			}
	// 		}
	// 		//iter := 0 NOTFALLS BENUTZEN WENN LAENGE REF DEUTLICH LAENGER ALS LAENGE COUNT
	// 		// fmt.Println("TESTLENREFERENCES")
	// 		// fmt.Println(len(references))
	// 		// fmt.Println("TESTLENCOUNTS")
	// 		// fmt.Println(len(counts))
	// 		// if len(references) > len(counts) { //////////////////////////////////////////////////////////////
	// 		// 	iter = len(counts)
	// 		// } else {
	// 		// 	iter = len(references)
	// 		// }////////////////////////////////////////////////////////////////////////////////////////////////
	// 		// 	}
	// 		// for i := 0; i < iter; i++ {
	// 		// 	// fmt.Println("REFERENCE: ")
	// 		// 	// fmt.Println(references[i])
	// 		// 	// fmt.Println("Counted: ")
	// 		// 	// fmt.Println(counts[i])
	// 		// }
	// 		ddref := analyze_dref(re)
	// 		if ddref != nil {
	// 			item := Paper{}
	// 			relation := Relationv3{}
	// 			relation.Citer_Name = li.Citer_Name
	// 			relation.Citer_Writer = li.Citer_Writer
	// 			relation.Citer_Year = li.Citer_Year
	// 			relation.Cites_Nr = li.Cites_Nr
	// 			relation.Phrases = phrases_c
	// 			relation.Tag = li.Tag
	// 			//item.Cited_by = li.Name()
	// 			relation.Citing_Writer = ddref[0]
	// 			relation.Citing_Name = ddref[1]
	// 			//item.Cites_Nr
	// 			intVar, err := strconv.Atoi(ddref[2])
	// 			if err != nil {
	// 				relation.Citing_Year = 0
	// 			} else {
	// 				relation.Citing_Year = intVar
	// 			}
	// 			i = i + 1
	// 			//EMPTYSPACE ALSO POSSIBLE BUT THEN MUST CLEAN ALL EMPTY SPACE IN BEFORE
	// 			writer_mod := strings.Trim(ddref[0], " ")
	// 			writer_bef, _, status := strings.Cut(writer_mod, " ")
	// 			if !status {
	// 				writer_mod = writer_bef
	// 			}
	// 			//fmt.Println("WRITERMOD" + writer_mod)
	// 			//Count_citations := count_ref(text_re, writer_mod)  //////////////// ALTERNATIVE BERECHNUNG MIT AUTOR FUNKTION BAUEN WENN CITES 0
	// 			relation.Cites_Nr2 = counts[k] ///// CHECK IF K IS ALWAYS RIGHT
	// 			relation.Tag2 = "Text Source - References"
	// 			newciter := count_ref2(text, cites[k])
	// 			fmt.Println("NEW")
	// 			fmt.Println(newciter)
	// 			fmt.Println("VS")
	// 			fmt.Println("OLD")
	// 			fmt.Println(counts[k])
	// 			relation.Cites_Nr2 = newciter
	// 			relation.Phrases = phrases_c
	// 			relation.Phrases_Cit = get_citref(text, cites[k])
	// 			fmt.Println("PHRASES_TEST")
	// 			fmt.Println(item.Phrases_Cit)

	// 			//ADD YEAR TO SEARCH FOR CITATION; NOCH WEITER TESTEN;
	// 			// 	writer_bef,_,status = strings.Cut(ddref[0], "and")
	// 			// 	if !status {
	// 			// 		writer_bef,_,status = strings.Cut(ddref[0], "&")
	// 			// 	}

	// 			// }

	// 			// do something with the article
	// 			cited_papers = append(cited_papers, relation)
	// 		}
	// 		k = k + 1
	// 	}
	// 	cited_fin = append(cited_fin, cited_papers...)
	// 	///////////////////////////////////////////////////////////////////////TEMPORARY PASSIVE////////////////////////
	// }

	js, err := json.MarshalIndent(cited_fin, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("cited_phraseswosfin2.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
		fmt.Print(k)
		fmt.Print(" references")
	}
	//////////////////////////////////////////////////////////////////////////////AKTUELL AUSKOMMENTIERT

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

	/////////////////////////////////////////////WORD CLOUD PREPARE/////////////////////////////
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

func get_citref2(text string, ref string, year string) []string {

	preref, _, bools := strings.Cut(text, "References")
	if !bools {
		preref, _, bools = strings.Cut(text, "REFERENCES")
	}
	if bools {
		text = preref
	}
	// 	fmt.Println(el)
	// }
	//conc_cites := strings.Join(citesmatchlist, "")
	count := 0
	var sentences_citref []string
	//str := "ref+`\,|\]{1}`"
	// count = strings.Count(conc_cites, ref) [\w\&\s\.]* ([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})
	//ADD EMPTY SPC ?

	if len(ref) > 0 {
		//PARSE CITATIONS IN FULL TEXT WITH FORMAT : FIRST AUTHOR ... (YEAR)
		strf := fmt.Sprintf(`%s{1}[\w\&\s\.]*\({1}%s{1}\){1}`, regexp.QuoteMeta(ref), regexp.QuoteMeta(year)) //meta uberfl?//notfalls gleich meta in die mitte und jeweils zwei klammer darum
		// fmt.Println("strf")
		// fmt.Println(strf)
		// fmt.Println("ref")
		// fmt.Println(ref)

		rex := regexp.MustCompile(strf) /////ODER PUNKTE ENDE LINEBREAK // geht notfalls normal

		//submatchall = rex.FindAllStringIndex(text, -1)
		submatchs := rex.FindAllString(text, -1)
		submatchs_ind := rex.FindAllStringIndex(text, -1)

		fmt.Println(len(submatchs))
		for i := 0; i < len(submatchs); i++ {
			//refstr := fmt.Sprintf(`%s{1}`, regexp.QuoteMeta(ref))
			//rext := regexp.MustCompile(refstr)
			//refstr_ind := rext.FindAllStringIndex(submatchs[i], -1)
			newtext := submatchs[i]
			//for l := 0; l < len(refstr_ind); l++ {
			if submatchs_ind[i][1]+1 < len(text) && strings.Contains(submatchs[i], ref) {

				fmt.Println("TESTSUBMATCHS ")
				fmt.Println(submatchs[i])
				fmt.Println(len(newtext))
				//fmt.Println(conc_cites[submatchs_ind[i][1] : submatchs_ind[i][1]+1])
				// fmt.Println("TESTCOUNTSSUBSTR2: ")
				// fmt.Println(conc_cites[submatchs_ind[i][1]+1 : submatchs_ind[i][1]+2]) //:submatchs_ind[i][1]+1
				//if (string(newtext[refstr_ind[l][1]]) == "," || string(newtext[refstr_ind[l][1]]) == "]") && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {

				indpre := strings.LastIndex(text[:submatchs_ind[i][0]], ".")
				indpre2 := strings.LastIndex(text[:submatchs_ind[i][0]], "?")
				if indpre2 > indpre {
					indpre = indpre2
				}
				indpost := strings.Index(text[submatchs_ind[i][1]:], ".")
				indpost2 := strings.Index(text[submatchs_ind[i][1]:], "?")
				if indpost2 != -1 && indpost > indpost2 {
					indpost = indpost2
				}
				textsh := text[submatchs_ind[i][1]:]

				if indpre != -1 && indpost != -1 {
					sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				} else if indpre != -1 && indpost == -1 {
					sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				} else if indpre == -1 && indpost != -1 {
					sentence := text[:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				}

				//}
			} // } else if strings.Contains(submatchs[i], ref) && refstr_ind[l][1]+1 == len(text) && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {
			// 	indpre := strings.LastIndex(text[:submatchs_ind[i][0]], ".")
			// 	indpost := strings.Index(text[submatchs_ind[i][1]:], ".")
			// 	textsh := text[submatchs_ind[i][1]:]

			// 	if indpre != -1 && indpost != -1 {
			// 		sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
			// 		fmt.Println("TESTSENT")
			// 		fmt.Println(sentence)
			// 		sentences_citref = append(sentences_citref, sentence)
			// 		count = count + 1
			// 	} else if indpre != -1 && indpost == -1 {
			// 		sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh
			// 		fmt.Println("TESTSENT")
			// 		fmt.Println(sentence)
			// 		sentences_citref = append(sentences_citref, sentence)
			// 		count = count + 1
			// 	} else if indpre == -1 && indpost != -1 {
			// 		sentence := text[:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
			// 		fmt.Println("TESTSENT")
			// 		fmt.Println(sentence)
			// 		sentences_citref = append(sentences_citref, sentence)
			// 		count = count + 1
			// 	}
			// } ([1]{1}[89]{1}[0-9]{2})|([2]{1}[0]{1}[0-9]{2})
			//}

		}

		//PARSE CITATIONS IN FULL TEXT WITH FORMAT : (FIRST AUTHOR ..., YEAR)
		strf2 := fmt.Sprintf(`\({1}%s{1}[\w\&\s\.\,]*%s{1}[\;\w\&\s\.\,]*\){1}`, regexp.QuoteMeta(ref), regexp.QuoteMeta(year))
		rex2 := regexp.MustCompile(strf2)
		submatchs = rex2.FindAllString(text, -1)
		submatchs_ind = rex2.FindAllStringIndex(text, -1)
		fmt.Println(len(submatchs))

		for i := 0; i < len(submatchs); i++ {
			//refstr := fmt.Sprintf(`%s{1}`, regexp.QuoteMeta(ref))
			//rext := regexp.MustCompile(refstr)
			//refstr_ind := rext.FindAllStringIndex(submatchs[i], -1)
			newtext := submatchs[i]
			//for l := 0; l < len(refstr_ind); l++ {
			if submatchs_ind[i][1]+1 < len(text) && strings.Contains(submatchs[i], ref) {

				fmt.Println("VORHER ")
				fmt.Println(ref + "  " + year)

				fmt.Println("TESTSUBMATCHS ")
				fmt.Println(submatchs[i])
				fmt.Println(len(newtext))
				//fmt.Println(conc_cites[submatchs_ind[i][1] : submatchs_ind[i][1]+1])
				// fmt.Println("TESTCOUNTSSUBSTR2: ")
				// fmt.Println(conc_cites[submatchs_ind[i][1]+1 : submatchs_ind[i][1]+2]) //:submatchs_ind[i][1]+1
				//if (string(newtext[refstr_ind[l][1]]) == "," || string(newtext[refstr_ind[l][1]]) == "]") && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {

				indpre := strings.LastIndex(text[:submatchs_ind[i][0]], ".")
				indpre2 := strings.LastIndex(text[:submatchs_ind[i][0]], "?")
				if indpre2 > indpre {
					indpre = indpre2
				}
				indpost := strings.Index(text[submatchs_ind[i][1]:], ".")
				indpost2 := strings.Index(text[submatchs_ind[i][1]:], "?")
				if indpost2 != -1 && indpost > indpost2 {
					indpost = indpost2
				}
				textsh := text[submatchs_ind[i][1]:]

				if indpre != -1 && indpost != -1 {
					sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				} else if indpre != -1 && indpost == -1 {
					sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				} else if indpre == -1 && indpost != -1 {
					sentence := text[:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				}

				//}
			} // } else if strings.Contains(submatchs[i], ref) && refstr_ind[l][1]+1 == len(text) && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {
			// 	indpre := strings.LastIndex(text[:submatchs_ind[i][0]], ".")
			// 	indpost := strings.Index(text[submatchs_ind[i][1]:], ".")
			// 	textsh := text[submatchs_ind[i][1]:]

			// 	if indpre != -1 && indpost != -1 {
			// 		sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
			// 		fmt.Println("TESTSENT")
			// 		fmt.Println(sentence)
			// 		sentences_citref = append(sentences_citref, sentence)
			// 		count = count + 1
			// 	} else if indpre != -1 && indpost == -1 {
			// 		sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh
			// 		fmt.Println("TESTSENT")
			// 		fmt.Println(sentence)
			// 		sentences_citref = append(sentences_citref, sentence)
			// 		count = count + 1
			// 	} else if indpre == -1 && indpost != -1 {
			// 		sentence := text[:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
			// 		fmt.Println("TESTSENT")
			// 		fmt.Println(sentence)
			// 		sentences_citref = append(sentences_citref, sentence)
			// 		count = count + 1
			// 	}
			// }
			//}

		}
	}

	fmt.Println("count")
	fmt.Println(count)
	fmt.Println(len(sentences_citref))
	//sent_c := strings.Join(sentences_citref, " ")

	//ANSONSTNE EINFACH DEN BESTEN SENTENCE EINFACH ANZEIGEN KOMPLETT
	//phrases := phr_rank_nonoun(sent_c)
	//count = len(submatchs)
	// fmt.Println("CHECK COUNT  ")
	// IF PHRASE DETECTION DELIVERED LESS THEN 3 PHRASES BETTER REPLACE WITH FULL SENTENCE
	//phrases_c := ""
	// if len(phrases) >= 3 {
	// 	phrases_c = strings.Join(phrases, ", ")
	// } else
	// if len(sentences_citref) > 0 {
	// 	phrases_c = sentences_citref[0]
	// }
	// fmt.Println(count)
	for j := 0; j < len(sentences_citref); j++ {
		sentences_citref[j] = strings.Replace(sentences_citref[j], "\r", " ", -1)
		sentences_citref[j] = strings.Replace(sentences_citref[j], "\n", " ", -1)
	}

	return sentences_citref
}

func get_citref(text string, ref string) []string {

	//add phr logic from below
	//citref := ""

	//CREATES LIST OF ALL CITATIONS AND COUNTS CITATIONS OF PAPER X
	// nrselect2 := regexp.MustCompile(`\[{1}[0-9\,\s]*\]{1}`)
	// citesmatchlist := nrselect2.FindAllString(text, -1)

	// for _, el := range citesmatchlist {
	// 	fmt.Println("CHECK IF RIGHT READ CITES  ")

	preref, _, bools := strings.Cut(text, "References")
	if !bools {
		preref, _, bools = strings.Cut(text, "REFERENCES")
	}
	if bools {
		text = preref
	}
	// 	fmt.Println(el)
	// }
	//conc_cites := strings.Join(citesmatchlist, "")
	count := 0
	var sentences_citref []string
	//str := "ref+`\,|\]{1}`"
	// count = strings.Count(conc_cites, ref)
	//ADD EMPTY SPC ?
	strf := fmt.Sprintf(`\[{1}[0-9\,\s]*%s{1}[0-9\,\s]*\]{1}`, regexp.QuoteMeta(ref)) //meta uberfl?//notfalls gleich meta in die mitte und jeweils zwei klammer darum
	fmt.Println("strf")
	fmt.Println(strf)
	fmt.Println("ref")
	fmt.Println(ref)

	rex := regexp.MustCompile(strf) /////ODER PUNKTE ENDE LINEBREAK // geht notfalls normal
	//submatchall = rex.FindAllStringIndex(text, -1)
	submatchs := rex.FindAllString(text, -1)
	submatchs_ind := rex.FindAllStringIndex(text, -1)
	fmt.Println(len(submatchs))
	for i := 0; i < len(submatchs); i++ {
		refstr := fmt.Sprintf(`%s{1}`, regexp.QuoteMeta(ref))
		rext := regexp.MustCompile(refstr)
		refstr_ind := rext.FindAllStringIndex(submatchs[i], -1)
		newtext := submatchs[i]
		for l := 0; l < len(refstr_ind); l++ {
			if submatchs_ind[i][1]+1 < len(text) && strings.Contains(submatchs[i], ref) {

				fmt.Println("TESTSUBMATCHS ")
				fmt.Println(submatchs[i])
				fmt.Println(len(newtext))
				//fmt.Println(conc_cites[submatchs_ind[i][1] : submatchs_ind[i][1]+1])
				// fmt.Println("TESTCOUNTSSUBSTR2: ")
				// fmt.Println(conc_cites[submatchs_ind[i][1]+1 : submatchs_ind[i][1]+2]) //:submatchs_ind[i][1]+1

				// CHECK FOR FALSE HITS LIKE [6] TO EXCLUDE FORMAT HITS ON [16]
				if (string(newtext[refstr_ind[l][1]]) == "," || string(newtext[refstr_ind[l][1]]) == "]") && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {

					indpre := strings.LastIndex(text[:submatchs_ind[i][0]], ".")
					indpre2 := strings.LastIndex(text[:submatchs_ind[i][0]], "?")
					if indpre2 > indpre {
						indpre = indpre2
					}
					indpost := strings.Index(text[submatchs_ind[i][1]:], ".")
					indpost2 := strings.Index(text[submatchs_ind[i][1]:], "?")
					if indpost2 != -1 && indpost > indpost2 {
						indpost = indpost2
					}
					textsh := text[submatchs_ind[i][1]:]

					if indpre != -1 && indpost != -1 {
						sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
						fmt.Println("TESTSENT")
						fmt.Println(sentence)
						sentences_citref = append(sentences_citref, sentence)
						count = count + 1
					} else if indpre != -1 && indpost == -1 {
						sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh
						fmt.Println("TESTSENT")
						fmt.Println(sentence)
						sentences_citref = append(sentences_citref, sentence)
						count = count + 1
					} else if indpre == -1 && indpost != -1 {
						sentence := text[:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
						fmt.Println("TESTSENT")
						fmt.Println(sentence)
						sentences_citref = append(sentences_citref, sentence)
						count = count + 1
					}

				}

				// CHECK FOR FALSE HITS LIKE [6] TO EXCLUDE FORMAT HITS ON [16]
			} else if strings.Contains(submatchs[i], ref) && refstr_ind[l][1]+1 == len(text) && (string(newtext[refstr_ind[l][0]-1]) == "," || string(newtext[refstr_ind[l][0]-1]) == "[" || string(newtext[refstr_ind[l][0]-1]) == " ") {
				indpre := strings.LastIndex(text[:submatchs_ind[i][0]], ".")
				indpost := strings.Index(text[submatchs_ind[i][1]:], ".")
				textsh := text[submatchs_ind[i][1]:]

				if indpre != -1 && indpost != -1 {
					sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				} else if indpre != -1 && indpost == -1 {
					sentence := text[indpre+1:submatchs_ind[i][0]] + submatchs[i] + textsh
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				} else if indpre == -1 && indpost != -1 {
					sentence := text[:submatchs_ind[i][0]] + submatchs[i] + textsh[:indpost+1]
					fmt.Println("TESTSENT")
					fmt.Println(sentence)
					sentences_citref = append(sentences_citref, sentence)
					count = count + 1
				}
			}
		}

	}

	fmt.Println("count")
	fmt.Println(count)
	fmt.Println(len(sentences_citref))
	//sent_c := strings.Join(sentences_citref, " ")

	//ANSONSTNE EINFACH DEN BESTEN SENTENCE EINFACH ANZEIGEN KOMPLETT
	//phrases := phr_rank_nonoun(sent_c)
	//count = len(submatchs)
	// fmt.Println("CHECK COUNT  ")
	// IF PHRASE DETECTION DELIVERED LESS THEN 3 PHRASES BETTER REPLACE WITH FULL SENTENCE
	//phrases_c := ""
	// if len(phrases) >= 3 {
	// 	phrases_c = strings.Join(phrases, ", ")
	// } else
	// if len(sentences_citref) > 0 {
	// 	phrases_c = sentences_citref[0]
	// }
	// fmt.Println(count)
	for j := 0; j < len(sentences_citref); j++ {
		sentences_citref[j] = strings.Replace(sentences_citref[j], "\r", " ", -1)
		sentences_citref[j] = strings.Replace(sentences_citref[j], "\n", " ", -1)
	}

	return sentences_citref
}

func phr_rank_nonoun(rawText string) []string {
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
	if len(rankedPhrases) < 20 {
		iterator = len(rankedPhrases)
	} else {
		iterator = 20
	}

	for i := 0; i < iterator; i++ {

		nleft := strings.Replace(rankedPhrases[i].Left, "-\r", "-", -1)
		nleft = strings.Replace(nleft, "\r", " ", -1)
		nleft = strings.Replace(nleft, "\n", " ", -1)

		nright := strings.Replace(rankedPhrases[i].Right, "-\r", "-", -1)
		nright = strings.Replace(nright, "\r", " ", -1)
		nright = strings.Replace(nright, "\n", " ", -1)

		s := nleft + " " + nright //fmt.Sprint(rankedPhrases[i])
		fmt.Println("FURTTEST")
		fmt.Println(rankedPhrases[i])
		fmt.Println(rankedPhrases[i].Left)
		fmt.Println(rankedPhrases[i].Right)
		fmt.Println(s)

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
		// }lst_ranked

		// HIER TESTEN ABFANGEN VON -
		if len(nleft) > 0 && len(nright) > 0 && string(nleft[0]) != "-" && string(nleft[len(nleft)-1]) != "-" && string(nright[0]) != "-" && string(nright[len(nright)-1]) != "-" {
			lst_ranked = append(lst_ranked, s)
		}

	}

	if len(lst_ranked) > 5 {
		lst_ranked = lst_ranked[:5]
	}

	return lst_ranked
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
	if len(rankedPhrases) < 20 {
		iterator = len(rankedPhrases)
	} else {
		iterator = 20
	}

	for i := 0; i < iterator; i++ {
		nleft := strings.Replace(rankedPhrases[i].Left, "-\r", "-", -1)
		nleft = strings.Replace(nleft, "\r", " ", -1)
		nleft = strings.Replace(nleft, "\n", " ", -1)

		nright := strings.Replace(rankedPhrases[i].Right, "-\r", "-", -1)
		nright = strings.Replace(nright, "\r", " ", -1)
		nright = strings.Replace(nright, "\n", " ", -1)

		s := nleft + " " + nright //fmt.Sprint(rankedPhrases[i])
		fmt.Println("FURTTEST")
		fmt.Println(rankedPhrases[i])
		fmt.Println(rankedPhrases[i].Left)
		fmt.Println(rankedPhrases[i].Right)
		fmt.Println(s)

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
		//lk := fmt.Sprint(rankedPhrases[i])
		if len(nleft) > 0 && len(nright) > 0 && string(nleft[0]) != "-" && string(nleft[len(nleft)-1]) != "-" && string(nright[0]) != "-" && string(nright[len(nright)-1]) != "-" {
			lst_ranked = append(lst_ranked, s)
		}
	}

	//l := 0
	// if len(lst_ranked) > 5 {
	//  l = 5
	// } else {
	//  l = len(lst_ranked)
	// }
	var doc_lst []string
	for h := 0; h < len(lst_ranked); h++ {
		doc, err := prose.NewDocument(lst_ranked[h])
		if err != nil {
			log.Fatal(err)
		}
		noun_or_adj := true
		counter := 0
		fmt.Println("TEST LEN")
		fmt.Println(len(doc.Tokens()))
		//if
		// Iterate over the doc's tokens:
		var allowed_nouns = []string{"NN", "NNP", "NNPS", "NNS"}
		var allowed_adj = []string{"JJ", "JJR", "JJS"}

		tags := doc.Tokens()
		for k := 0; k < len(tags); k++ {
			// fmt.Println(tok.Text, tok.Tag, tok.Label)
			if contains(allowed_nouns, tags[k].Tag) {
				counter = counter + 1 //JJ JJR JJS RB RBR RBS RP
				//noun_or_adj = false
				//Phrase detected which does not confirm requirements
			} else if !contains(allowed_adj, tags[k].Tag) {
				noun_or_adj = false

			}
			// Go NNP B-GPE
			// is VBZ O
			// an DT O
			// ...
		}

		if noun_or_adj == true && counter >= 1 && lst_ranked[h] != "" {
			doc_lst = append(doc_lst, lst_ranked[h])
		}
		// if counter >= 1 {
		// 	min_noun = false
		// }

	}
	if len(doc_lst) > 5 {
		doc_lst = doc_lst[:5]
	}

	return doc_lst
}

// func convrtToUTF8(str string, origEncoding string) string {
// 	strBytes := []byte(str)
// 	byteReader := bytes.NewReader(strBytes)
// 	reader, _ := charset.NewReaderLabel(origEncoding, byteReader)
// 	strBytes, _ = ioutil.ReadAll(reader)
// 	return string(strBytes)
// }

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
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
