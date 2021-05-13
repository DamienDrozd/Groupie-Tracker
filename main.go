package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("templates/css")))) // The server will analyse the static folder to seach thes called files in the html
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("templates/img"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("templates/js"))))
	// http.HandleFunc("/", serveur)
	http.HandleFunc("/search", search)
	// http.HandleFunc("/result", result)
	http.HandleFunc("/artist", groupe)
	http.ListenAndServe(":8080", nil)
}

func serveur(w http.ResponseWriter, r *http.Request) {
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("templates/*.html"))
	templates.ExecuteTemplate(w, "index.html", nil)
}

type Tab struct {
	// Download string
	Name         string
	Image        string
	Members      string
	CreationDate string
	Url          string
	Id           string
	Search       []string

	// Test string
}

type Output struct {
	// Download string
	Tab             []Tab
	Tabmembers      []string
	Tabcreationdate []string
	Tabfirstalbum   []string
	Tablocations    []string
	Tabsearchbar    []string

	// Test string
}

var run = false

func search(w http.ResponseWriter, r *http.Request) {

	//--------------------------------------SearchBar--------------------------------
	timestart := time.Now()

	var result Group

	SearchBar := ""
	if r.FormValue("search") != "" {
		SearchBar = r.FormValue("search")
		for i := range SearchBar {
			if SearchBar[i] == 124 {
				SearchBar = SearchBar[:i-1]
				break
			}

		}
	}

	SearchBar = strings.ToLower(SearchBar)
	SearchBar = strings.Title(SearchBar)

	// fmt.Println(SearchBar)

	result.Name = SearchBar
	result.Members = append(result.Members, SearchBar)
	result.CreationDate, _ = strconv.ParseFloat(SearchBar, 64)
	result.FirstAlbum = SearchBar
	SearchBar = strings.ToLower(SearchBar)
	result.Locations = append(result.Members, SearchBar)

	var resultTab = make([]string, 0)

	resultgroup := findgroup(result) // -----------------------------------A optimiser------------------------

	for i := range resultgroup["ResultName"] {
		resultTab = append(resultTab, resultgroup["ResultName"][i])
	}
	for i := range resultgroup["ResultMembers"] {
		resultTab = append(resultTab, resultgroup["ResultMembers"][i])
	}
	for i := range resultgroup["ResultCreationDate"] {
		resultTab = append(resultTab, resultgroup["ResultCreationDate"][i])
	}
	for i := range resultgroup["ResultFirstAlbum"] {
		resultTab = append(resultTab, resultgroup["ResultFirstAlbum"][i])
	}
	for i := range resultgroup["ResultLocations"] {
		resultTab = append(resultTab, resultgroup["ResultLocations"][i])
	}

	// Barre de recherche

	// ---------------------------------DropDown-------------------------------------------
	var tabtest Group

	if r.FormValue("Members") != "" {
		tabtest.Members = append(result.Members, r.FormValue("Members"))
	}
	if r.FormValue("CreationDate") != "" {
		tabtest.CreationDate, _ = strconv.ParseFloat(r.FormValue("CreationDate"), 64)
	}
	if r.FormValue("FirstAlbum") != "" {
		tabtest.FirstAlbum = r.FormValue("FirstAlbum")
	}
	if r.FormValue("Locations") != "" {
		tabtest.Locations = append(result.Locations, r.FormValue("Locations"))

	}

	tab1 := findgroup(tabtest)

	resultTab = compareTab(resultTab, tab1["ResultMembers"])
	resultTab = compareTab(resultTab, tab1["ResultCreationDate"])
	resultTab = compareTab(resultTab, tab1["ResultFirstAlbum"])
	resultTab = compareTab(resultTab, tab1["resultTab"])

	// ---------------------------------------------Objet a renvoyer en requête http---------------------------------------------------

	var tab []Tab
	var tabmembers []string
	var tabcreationdate []string
	var tabfirstalbum []string
	var tablocations []string
	var tabsearchbar []string

	listGroup := listof()

	tabsearchbar = searchbar(listGroup)

	for i := range listGroup {
		tab = append(tab, Tab{})

		tab[i].Name = listGroup[i].Name
		tab[i].Image = listGroup[i].Image
		tab[i].Url = "/artist?artist=" + fmt.Sprintf("%v", i+1)

		tabcreationdate = append(tabcreationdate, fmt.Sprintf("%v", listGroup[i].CreationDate))
		tabfirstalbum = append(tabfirstalbum, fmt.Sprintf("%v", listGroup[i].FirstAlbum))

		for j := range listGroup[i].Locations {
			tablocations = append(tablocations, listGroup[i].Locations[j])
		}

		for j := range listGroup[i].Members {
			tabmembers = append(tabmembers, listGroup[i].Members[j])
		}

	}
	tabsearchbar = tritab(tabsearchbar, false)
	tablocations = tritab(tablocations, true)
	tabcreationdate = tritab(tabcreationdate, false)
	tabfirstalbum = tritab(tabfirstalbum, true)
	tabmembers = tritab(tabmembers, false)

	displaytab := resultTab
	if displaytab != nil && run == true && (SearchBar != "" || r.FormValue("Members") != "" || r.FormValue("CreationDate") != "" || r.FormValue("FirstAlbum") != "" || r.FormValue("Locations") != "") {
		run = false
		var displaytabgroup []Group

		for i := range displaytab {
			displaytabgroup = append(displaytabgroup, groupof(displaytab[i]))
		}

		tab = make([]Tab, 0)
		for i := range displaytabgroup {
			tab = append(tab, Tab{})

			tab[i].Name = displaytabgroup[i].Name
			tab[i].Image = displaytabgroup[i].Image
			tab[i].Url = "/artist?artist=" + fmt.Sprintf("%v", displaytabgroup[i].Id)
		}

	}

	p := Output{
		Tab:             tab,
		Tabmembers:      tabmembers,
		Tabcreationdate: tabcreationdate,
		Tabfirstalbum:   tabfirstalbum,
		Tablocations:    tablocations,
		Tabsearchbar:    tabsearchbar,
	}

	templates := template.New("Label de ma template")
	templates = template.Must(templates.ParseFiles("./templates/recherche.html"))
	err := templates.ExecuteTemplate(w, "recherche", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err) // If the executetemplate function cannot run, displays an error message
	}
	run = true

	t := time.Now()
	fmt.Println("time1:", t.Sub(timestart))
	timestart = time.Now()

}

func result(w http.ResponseWriter, r *http.Request) {
	var templates *template.Template
	templates = template.Must(templates.ParseGlob("templates/*.html"))
	templates.ExecuteTemplate(w, "result.html", nil)
}

func groupe(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["artist"]

	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 'key' is missing")
		return
	}

	// Query()["key"] will return an array of items,
	// we only want the single item.
	key := keys[0]

	log.Println("Url Param 'key' is: " + string(key))

	var test = string(key)

	GroupOutput := groupof(test)
	// fmt.Println(GroupOutput.Id)

	p := GroupOutput

	templates := template.New("Label de ma template")
	templates = template.Must(templates.ParseFiles("./templates/artist.html"))
	err := templates.ExecuteTemplate(w, "artist", p)

	if err != nil {
		log.Fatalf("Template execution: %s", err) // If the executetemplate function cannot run, displays an error message
	}
}

//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------
//----------------------------------------------------------------

type Group struct {
	Id           float64
	Name         string
	Image        string
	Members      []string
	CreationDate float64
	FirstAlbum   string
	Locations    []string
	ConcertDates []string
	Relations    map[string][]string
	RelationsTab [][]string
}

func tritab(tab []string, revert bool) []string {

	x := len(tab)

	for i := 0; i < x; i++ {
		for j := 0; j < x; j++ {
			if tab[i] < tab[j] {
				tab[i], tab[j] = tab[j], tab[i]
			}
		}
	}

	for i := 0; i < x; i++ {
		for j := 0; j < x; j++ {
			if tab[i] == tab[j] && i != j {
				// fmt.Println(tab[i], tab[j])
				tab = append(tab[:i], tab[i+1:]...)
				j--
				x--

			}

		}
	}

	if revert == true {
		for i := 0; i < len(tab); i++ {
			for j := 0; j < len(tab); j++ {
				var mota string
				var motb string
				for a := len(tab[i]) - 1; a > 0; a-- {
					if tab[i][a] == 45 {
						mota = tab[i][a:]
						break
					}
				}
				for b := len(tab[j]) - 1; b > 0; b-- {
					if tab[j][b] == 45 {
						motb = tab[j][b:]
						break
					}
				}

				if mota < motb {
					tab[i], tab[j] = tab[j], tab[i]
				}

			}
		}
	}

	return tab
}

func listof() []Group {
	var listGroup []Group
	var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")
	var loc = readurl("https://groupietrackers.herokuapp.com/api/locations")

	for i := range tab {

		var group Group
		group.Name = fmt.Sprintf("%v", tab[i]["name"])
		group.Image = fmt.Sprintf("%v", tab[i]["image"])
		group.FirstAlbum = fmt.Sprintf("%v", tab[i]["firstAlbum"])
		group.CreationDate = tab[i]["creationDate"].(float64)
		group.Id = 1 //tab[i]["Id"].(float64)
		for j := range tab[i]["members"].([]interface{}) {
			group.Members = append(group.Members, fmt.Sprintf("%v", tab[i]["members"].([]interface{})[j]))
		}

		// locations := readurl(fmt.Sprintf("%v", tab[i]["locations"])) //--------------A optimiser-----------------

		for j := range loc[i]["locations"].([]interface{}) {
			group.Locations = append(group.Locations, fmt.Sprintf("%v", loc[i]["locations"].([]interface{})[j]))

		}

		// dates := readurl(fmt.Sprintf("%v", tab[i]["concertDates"]))
		// for j := range dates[0]["dates"].([]interface{}) {
		// 	group.ConcertDates = append(group.ConcertDates, fmt.Sprintf("%v", dates[0]["dates"].([]interface{})[j]))
		// }
		listGroup = append(listGroup, group)
	}
	return listGroup
}

func searchbar(group []Group) []string {
	var result []string
	for i := range group {
		result = append(result, "")
		result = append(result, group[i].Name+" | Artist - Band")
		result = append(result, group[i].FirstAlbum+" | FirstAlbum")
		result = append(result, fmt.Sprintf("%v", group[i].CreationDate)+" | CreationDate")
		for j := range group[i].Members {
			result = append(result, group[i].Members[j]+" | Members")
		}
		for j := range group[i].Locations {
			result = append(result, group[i].Locations[j]+" | Locations")
		}
	}

	return result
}

func groupof(input string) Group {
	var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")
	var group Group
	for i := range tab {
		// fmt.Println(reflect.TypeOf(fmt.Sprintf("%v", tab[i]["id"])), reflect.TypeOf(input))
		if fmt.Sprintf("%v", tab[i]["id"]) == input {

			group.Name = fmt.Sprintf("%v", tab[i]["name"])
			group.Image = fmt.Sprintf("%v", tab[i]["image"])
			group.FirstAlbum = fmt.Sprintf("%v", tab[i]["firstAlbum"])
			group.CreationDate = tab[i]["creationDate"].(float64)
			group.Id = tab[i]["id"].(float64)
			for j := range tab[i]["members"].([]interface{}) {
				group.Members = append(group.Members, fmt.Sprintf("%v", tab[i]["members"].([]interface{})[j]))
			}
			locations := readurl(fmt.Sprintf("%v", tab[i]["locations"]))
			for j := range locations[0]["locations"].([]interface{}) {
				group.Locations = append(group.Locations, fmt.Sprintf("%v", locations[0]["locations"].([]interface{})[j]))
			}
			dates := readurl(fmt.Sprintf("%v", tab[i]["concertDates"]))
			for j := range dates[0]["dates"].([]interface{}) {
				group.ConcertDates = append(group.ConcertDates, fmt.Sprintf("%v", dates[0]["dates"].([]interface{})[j]))
			}
			tabconvert := readurl(fmt.Sprintf("%v", tab[i]["relations"]))
			group.Relations = makerelations(tabconvert)
			var x = 0
			for j := range group.Relations {
				group.RelationsTab = append(group.RelationsTab, make([]string, 0))
				group.RelationsTab[x] = append(group.RelationsTab[x], j)
				for k := range group.Relations[j] {
					group.RelationsTab[x] = append(group.RelationsTab[x], group.Relations[j][k])
				}
				x++
			}
		}
	}
	return group
}

func findgroup(input Group) map[string][]string {

	var ResultName = make([]string, 0)
	var ResultMembers = make([]string, 0)
	var ResultCreationDate = make([]string, 0)
	var ResultFirstAlbum = make([]string, 0)
	var ResultLocations = make([]string, 0)
	var ResultConcertDates = make([]string, 0)
	var ResultRelations = make([]string, 0)
	var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")

	if input.Name != "" {

		for i := range tab {
			if tab[i]["name"] == input.Name {
				ResultName = append(ResultName, fmt.Sprintf("%v", tab[i]["id"]))
			}
		}
	}
	if input.CreationDate != 0 {

		for i := range tab {
			if tab[i]["creationDate"] == input.CreationDate {
				ResultCreationDate = append(ResultCreationDate, fmt.Sprintf("%v", tab[i]["id"]))
			}
		}
	}
	if input.FirstAlbum != "" {

		for i := range tab {
			if tab[i]["firstAlbum"] == input.FirstAlbum {
				ResultFirstAlbum = append(ResultFirstAlbum, fmt.Sprintf("%v", tab[i]["id"]))
			}
		}
	}
	if input.Members != nil {

		for i := range tab {
			for k := 0; k < len(tab[i]["members"].([]interface{})); k++ {
				for j := range input.Members {
					arr := tab[i]["members"].([]interface{})
					if input.Members[j] == arr[k] {
						// fmt.Println(ResultMembers)
						ResultMembers = append(ResultMembers, fmt.Sprintf("%v", tab[i]["id"]))
					}
				}
			}
		}
	}
	if input.Locations != nil {
		tab := readurl("https://groupietrackers.herokuapp.com/api/locations")

		for i := range tab {
			arr := tab[i]["locations"].([]interface{})
			for k := 0; k < len(arr); k++ {
				for j := range input.Locations {
					if input.Locations[j] == arr[k] {
						ResultLocations = append(ResultLocations, fmt.Sprintf("%v", tab[i]["id"]))
					}
				}
			}
		}
	}
	if input.ConcertDates != nil {
		var tab = readurl("https://groupietrackers.herokuapp.com/api/dates")
		for i := range tab {
			// fmt.Println(tab[i])
			arr := tab[i]["dates"].([]interface{})
			// fmt.Println(tab[i]["dates"].([]interface{}))
			for k := 0; k < len(arr); k++ {
				for j := range input.ConcertDates {
					if input.ConcertDates[j] == arr[k] {
						ResultConcertDates = append(ResultConcertDates, fmt.Sprintf("%v", tab[i]["id"]))
					}
				}
			}
		}
	}
	if input.ConcertDates != nil && input.Locations != nil {
		tabconvert := readurl("https://groupietrackers.herokuapp.com/api/relation")
		tab := make([]map[string]string, 0)
		for i := range tabconvert {
			mapInterface := make(map[string]interface{})
			mapString := make(map[string]string)
			mapInterface = tabconvert[i]
			for key, value := range mapInterface {
				strKey := fmt.Sprintf("%v", key)
				strValue := fmt.Sprintf("%v", value)
				mapString[strKey] = strValue
			}
			tab = append(tab, mapString)
		}
		notes := makerelations(tabconvert)
		for x := range tab {
			for i := range input.Locations {
				for j := range input.ConcertDates {
					if notes[input.Locations[i]] != nil {
						for k := range notes[input.Locations[i]] {
							// fmt.Println(notes[input.Locations[i]][k], input.ConcertDates[j][1:])
							if notes[input.Locations[i]][k] == input.ConcertDates[j][1:] {
								ResultRelations = append(ResultRelations, tab[x]["id"])
							}
						}
					}
				}
			}
		}

	}
	var result = make(map[string][]string, 0)
	result["ResultName"] = ResultName
	result["ResultMembers"] = ResultMembers
	result["ResultFirstAlbum"] = ResultFirstAlbum
	result["ResultCreationDate"] = ResultCreationDate
	result["ResultLocations"] = ResultLocations
	result["ResultConcertDates"] = ResultConcertDates
	result["ResultRelations"] = ResultRelations
	return result

}

func makerelations(tabconvert []map[string]interface{}) map[string][]string {
	tab := make([]map[string]string, 0)
	var notes = make(map[string][]string)
	for i := range tabconvert {
		mapInterface := make(map[string]interface{})
		mapString := make(map[string]string)
		mapInterface = tabconvert[i]
		for key, value := range mapInterface {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)

			mapString[strKey] = strValue
		}
		tab = append(tab, mapString)
	}

	for x := range tab {

		arr := tab[x]["datesLocations"]

		addkey := true
		adddict := false

		key := string(arr[3:4])
		mot := ""

		for i := 4; i < len(arr)-1; i++ {
			// fmt.Println(string(arr[i]), addkey, adddict)
			if arr[i] == 93 {

				addkey = true
				// fmt.Println(string(key))
				// fmt.Println(string(mot))
				var array = make([]string, 1)
				x := 0
				for j := 1; j < len(mot); j++ {
					// fmt.Println(string(mot[j])
					// fmt.Println(x, len(array))
					if mot[j] != 32 {
						array[x] += string(mot[j])
					} else {
						x++
						array = append(array, "")
					}

				}
				notes[key[1:]] = array
				// fmt.Println(key)
				mot = ""
				key = ""

			}
			if arr[i] == 58 {
				addkey = false
				i++

			}
			if arr[i] == 91 {
				adddict = true

			}
			if arr[i] == 93 {
				adddict = false
				i++
				// mot += string(arr[i])
			}

			if addkey == true {

				key += string(arr[i])

			}
			if adddict == true {

				mot += string(arr[i])

			}

		}

	}

	return notes
}

func compareTab(tab1, tab2 []string) []string {

	var result = make([]string, 0)

	if len(tab1) == 0 {
		tab1 = tab2
	}

	if len(tab2) == 0 {
		result = tab1
	}

	for i := range tab1 {
		for j := range tab2 {
			if tab1[i] == tab2[j] {
				add := true

				for x := range result {
					if tab1[i] == result[x] {
						add = false
						break
					}
				}

				if add == true {
					result = append(result, tab1[i])
				}

			}

		}
	}

	return result
}

func readurl(url string) []map[string]interface{} {

	spaceClient := http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "spacecount-tutorial")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	data, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return transformtab(data)
}

func transformtab(data []byte) []map[string]interface{} {

	var tabstring = make([]string, 1)
	var j = 0

	if data[0] == 91 {
		for i := 1; i < len(data); i++ {

			if i < len(data)-2 {
				// fmt.Println(i, " : ", string(data[i]), data[i])
				if data[i] == 125 && data[i+1] == 44 && data[i+2] == 123 {
					tabstring[j] += string(data[i])
					i++
					j++

					tabstring = append(tabstring, "")

				} else {
					tabstring[j] += string(data[i])
				}
			}
		}
	} else {
		tabstring[0] = string(data)
	}

	if string(tabstring[0][0:9]) == "{\"index\":" {
		tabstring[0] = tabstring[0][10:]
		tabstring[0] = tabstring[0][:len(tabstring[0])-3]

		cutString := tabstring[0]
		tabstring[0] = ""

		for i := 0; i < len(cutString)-3; i++ {

			if cutString[i:i+3] == "},{" {

				tabstring[j] += string(cutString[i])

				i++
				j++

				tabstring = append(tabstring, "")

			} else {
				tabstring[j] += string(cutString[i])
			}

		}
		tabstring[len(tabstring)-1] += cutString[len(cutString)-3:]
	}

	var tabmap = make([]map[string]interface{}, 0)

	for i := 0; i < len(tabstring); i++ {
		// fmt.Println(tabstring[i])
		// Declared an empty map interface
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(tabstring[i]), &result)

		tabmap = append(tabmap, result)

	}

	return tabmap

}
