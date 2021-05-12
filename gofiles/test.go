package gofiles

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// func main() {

// 	var input Group

// 	input.CreationDate = 1997
// 	input.Name = "SOJA"
// 	input.Members = []string{"Jacob Hemphill", "Bob Jefferson", "Ryan \"Byrd\" Berty", "Ken Brownell", "Patrick O'Shea", "Hellman Escorcia", "Rafael Rodriguez", "Trevor Young"}
// 	input.FirstAlbum = "05-06-2002"
// 	input.ConcertDates = []string{"*05-12-2019", "06-12-2019", "07-12-2019", "08-12-2019", "09-12-2019", "*16-11-2019", "*15-11-2019"}
// 	input.Locations = []string{"playa_del_carmen-mexico", "papeete-french_polynesia", "noumea-new_caledonia"}
// 	// input.Relations = map[string][]string{"dunedin-new_zealand" : ["10-02-2020"] , "georgia-usa":["22-08-2019"],"los_angeles-usa":["20-08-2019"],"nagoya-japan":["30-01-2019"],"north_carolina-usa":["23-08-2019"],"osaka-japan":["28-01-2020"],"penrose-new_zealand":["07-02-2020"],"saitama-japan":["26-01-2020"]}

// 	tabinput := findgroup(input)

// 	for i := range tabinput {
// 		groupof(tabinput[i])
// 	}

// 	// fmt.Println(readurl("https://groupietrackers.herokuapp.com/api"))
// }

type Group struct {
	Name         string
	Image        string
	Members      []string
	CreationDate float64
	FirstAlbum   string
	Locations    []string
	ConcertDates []string
	Relations    map[string][]string
}

func listof() {
	var listGroup []Group
	var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")
	for i := range tab {
		listGroup[i].Name = tab[i].Name

		listGroup[i].Name = fmt.Sprintf("%v", tab[i]["name"])
		listGroup[i].Image = fmt.Sprintf("%v", tab[i]["image"])
		listGroup[i].FirstAlbum = fmt.Sprintf("%v", tab[i]["firstAlbum"])
		listGroup[i].CreationDate = tab[i]["creationDate"].(float64)

		for j := range tab[i]["members"].([]interface{}) {
			listGroup[i].Members = append(listGroup[i].Members, fmt.Sprintf("%v", tab[i]["members"].([]interface{})[j]))
		}

		locations := readurl(fmt.Sprintf("%v", tab[i]["locations"]))
		for j := range locations[0]["locations"].([]interface{}) {
			listGroup[i].Locations = append(listGroup[i].Locations, fmt.Sprintf("%v", locations[0]["locations"].([]interface{})[j]))
		}

		dates := readurl(fmt.Sprintf("%v", tab[i]["concertDates"]))
		for j := range dates[0]["dates"].([]interface{}) {
			listGroup[i].ConcertDates = append(listGroup[i].ConcertDates, fmt.Sprintf("%v", dates[0]["dates"].([]interface{})[j]))
		}

		tabconvert := readurl(fmt.Sprintf("%v", tab[i]["relations"]))

		listGroup[i].Relations = makerelations(tabconvert)
	}

	for i := range listGroup {
		fmt.Prinln(listGroup[i].Name)
	}
}

func groupof(input string) {
	var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")
	for i := range tab {
		// fmt.Println(reflect.TypeOf(fmt.Sprintf("%v", tab[i]["id"])), reflect.TypeOf(input))
		if fmt.Sprintf("%v", tab[i]["id"]) == input {

			var group Group
			group.Name = fmt.Sprintf("%v", tab[i]["name"])
			group.Image = fmt.Sprintf("%v", tab[i]["image"])
			group.FirstAlbum = fmt.Sprintf("%v", tab[i]["firstAlbum"])
			group.CreationDate = tab[i]["creationDate"].(float64)

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

			fmt.Println("group.Name", group.Name)
			fmt.Println("group.Image", group.Image)
			fmt.Println("group.Members", group.Members)
			fmt.Println("group.CreationDate", group.CreationDate)
			fmt.Println("group.FirstAlbum", group.FirstAlbum)
			fmt.Println("group.Locations", group.Locations)
			fmt.Println("group.ConcertDates", group.ConcertDates)
			fmt.Println("group.Relations", group.Relations)
		}
	}
}

func findgroup(input Group) []string {

	var ResultName = make([]string, 0)
	var ResultMembers = make([]string, 0)
	var ResultCreationDate = make([]string, 0)
	var ResultFirstAlbum = make([]string, 0)
	var ResultLocations = make([]string, 0)
	var ResultConcertDates = make([]string, 0)
	var ResultRelations = make([]string, 0)

	if input.Name != "" {
		var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")
		for i := range tab {

			if tab[i]["name"] == input.Name {
				ResultName = append(ResultName, fmt.Sprintf("%v", tab[i]["id"]))
			}

		}
	}

	if input.CreationDate != 0 {
		var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")
		for i := range tab {
			if tab[i]["creationDate"] == input.CreationDate {

				ResultCreationDate = append(ResultCreationDate, fmt.Sprintf("%v", tab[i]["id"]))
			}
		}
	}

	if input.FirstAlbum != "" {
		var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")
		for i := range tab {
			if tab[i]["firstAlbum"] == input.FirstAlbum {
				ResultFirstAlbum = append(ResultFirstAlbum, fmt.Sprintf("%v", tab[i]["id"]))
			}
		}
	}

	if input.Members != nil {
		var tab = readurl("https://groupietrackers.herokuapp.com/api/artists")

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

	// fmt.Println("ResultName", ResultName)
	// fmt.Println("ResultMembers", ResultMembers)
	// fmt.Println("ResultFirstAlbum", ResultFirstAlbum)
	// fmt.Println("ResultCreationDate", ResultCreationDate)
	// fmt.Println("ResultLocations", ResultLocations)
	// fmt.Println("ResultConcertDates", ResultConcertDates)
	// fmt.Println("ResultRelations", ResultRelations)

	var result = make([]string, 0)
	result = compareTab(result, ResultName)
	result = compareTab(result, ResultMembers)
	result = compareTab(result, ResultCreationDate)
	result = compareTab(result, ResultFirstAlbum)
	result = compareTab(result, ResultLocations)
	result = compareTab(result, ResultConcertDates)
	// result = compareTab(ResultLocations, ResultConcertDates)
	fmt.Println("result", result)

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
