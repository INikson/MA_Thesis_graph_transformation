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

func main() {
	//////////////////////////////////////////////////////////WIEDER NUTZEN UNIDOC
	// err, res := outputPdfText("pdfs/489426276.pdf")
	// if err != nil {
	// 	log.Fatal(err)
	// } //////////////////////////////////////////////////////////////////////////HIER
	cited_fin := make([]Paper, 0)

	listpdf, err := ioutil.ReadDir("C:/Users/neumu/pdffer3/txt")
	if err != nil {
		panic(err)
	}
	//var def []map[int]string
	//def := make([]map[int]string, 0)

	// for _, li := range listpdf {
	// 	// do something with the article
	// 	fmt.Println(li.Name())
	// 	err, pager := outputPdfText("pdf/"+li.Name(), li.Name())
	// 	if err != nil {
	// 		panic(err)
	// 	} else {
	// 		def = append(def, pager)
	// 	}

	// }

	// listpdf = listpdf[:5]
	i := 0
	for _, li := range listpdf {

		///////////////////////////////////////////////////////////////////////TEMPORARY PASSIVE////////////////////////

		// do something with the article
		content, err := ioutil.ReadFile("txt/" + li.Name())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(li.Name())

		// Convert []byte to string and print to screen
		text := string(content)
		references, counts, phrases, cites := analyze_ref(text)

		fmt.Println("COUNTSTEST")
		fmt.Println(len(counts))
		fmt.Println(len(references))
		fmt.Println("CITESSTEST")
		fmt.Println(len(cites))

		phrases_c := strings.Join(phrases, ", ")
		// fmt.Println("PHRASES_TEST")
		// fmt.Println(phrases_c)

		// content, err := ioutil.ReadFile("C:/Users/neumu/pdffer/pdfs/10-1108_BPMJ-10-2021-0677.pdf")
		// if err != nil {
		// 	log.Fatal(err)
		// }
		//fmt.Println(counts[0])
		cited_papers := make([]Paper, 0)
		// fmt.Println(content)
		//text_re := text
		k := 0
		for _, re := range references { // TEST (http|https|ftp): WEGLASSEN
			rex := regexp.MustCompile(`(http|https|ftp):[\/]{2}([a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,4})(:[0-9]+)?\/?([a-zA-Z0-9\-\._\?\,\'\/\\\+&amp;%\$#\=~]*)`)

			submatchall := rex.FindAllString(re, -1)
			if len(submatchall) > 0 {
				//fmt.Println("Foundyeyy")
				//fmt.Println(len(submatchall))
				for _, su := range submatchall {
					//fmt.Println(su)
					re = strings.Replace(re, su, "", -1)
				}
			}
			//iter := 0 NOTFALLS BENUTZEN WENN LAENGE REF DEUTLICH LAENGER ALS LAENGE COUNT
			// fmt.Println("TESTLENREFERENCES")
			// fmt.Println(len(references))
			// fmt.Println("TESTLENCOUNTS")
			// fmt.Println(len(counts))
			// if len(references) > len(counts) { //////////////////////////////////////////////////////////////
			// 	iter = len(counts)
			// } else {
			// 	iter = len(references)
			// }////////////////////////////////////////////////////////////////////////////////////////////////
			// 	}
			// for i := 0; i < iter; i++ {
			// 	// fmt.Println("REFERENCE: ")
			// 	// fmt.Println(references[i])
			// 	// fmt.Println("Counted: ")
			// 	// fmt.Println(counts[i])
			// }
			ddref := analyze_dref(re)
			if ddref != nil {
				item := Paper{}
				item.Cited_by = li.Name()
				item.Writer = ddref[0]
				item.Name = ddref[1]
				//item.Cites_Nr
				intVar, err := strconv.Atoi(ddref[2])
				if err != nil {
					item.Year = 0
				} else {
					item.Year = intVar
				}
				i = i + 1
				//EMPTYSPACE ALSO POSSIBLE BUT THEN MUST CLEAN ALL EMPTY SPACE IN BEFORE
				writer_mod := strings.Trim(ddref[0], " ")
				writer_bef, _, status := strings.Cut(writer_mod, " ")
				if !status {
					writer_mod = writer_bef
				}
				//fmt.Println("WRITERMOD" + writer_mod)
				//Count_citations := count_ref(text_re, writer_mod)  //////////////// ALTERNATIVE BERECHNUNG MIT AUTOR FUNKTION BAUEN WENN CITES 0
				item.Cites_Nr = counts[k] ///// CHECK IF K IS ALWAYS RIGHT
				newciter := count_ref2(text, cites[k])
				fmt.Println("NEW")
				fmt.Println(newciter)
				fmt.Println("VS")
				fmt.Println("OLD")
				fmt.Println(counts[k])
				item.Cites_Nr = newciter
				item.Phrases = phrases_c
				item.Phrases_Cit = get_citref(text, cites[k])
				fmt.Println("PHRASES_TEST")
				fmt.Println(item.Phrases_Cit)

				//ADD YEAR TO SEARCH FOR CITATION; NOCH WEITER TESTEN;
				// 	writer_bef,_,status = strings.Cut(ddref[0], "and")
				// 	if !status {
				// 		writer_bef,_,status = strings.Cut(ddref[0], "&")
				// 	}

				// }

				// do something with the article
				cited_papers = append(cited_papers, item)
			}
			k = k + 1
		}
		cited_fin = append(cited_fin, cited_papers...)
		///////////////////////////////////////////////////////////////////////TEMPORARY PASSIVE////////////////////////
	}

	js, err := json.MarshalIndent(cited_fin, "", "    ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Writing data to file")
	if err := os.WriteFile("cited_phrasesfin.json", js, 0664); err == nil {
		fmt.Println("Data written to file successfully")
		fmt.Print(i)
		fmt.Print(" references")
	}

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
