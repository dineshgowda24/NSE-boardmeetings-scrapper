package main

import(
	"net/http"
	"log"
	json "github.com/rogpeppe/rjson"
	"html/template"
)

type Rows struct{
	Symbol string `json:"Symbol"`
	CompanyName string `json:"CompanyName"`
	Isin string `json:"ISIN"`
	Ind string `json:"Ind"`
	Purpose string `json:"Purpose"`
	BoardMeetingDate string `json:"BoardMeetingDate"`
	DisplayDate string `json:"DisplayDate"`
	SeqID string `json:"seqId"`
	Details string `json:"Details"`
}

type BoardMeetingResponse struct{
	Success bool `json:"success"`
	Results int `json:"results"`
	Row []Rows `json:"rows"`
}

func serverBoardMeetings(w http.ResponseWriter, req *http.Request){
	if req.Method != http.MethodGet {
		http.Error(w, "Something isn't right here!", http.StatusNotFound)
		return
	}

	err := req.ParseForm()
	if err != nil {
		log.Fatal(err)
		return
	}
	var comp string
	if len( req.Form.Get("Symbol") ) > 0	{
		comp = (req.Form["Symbol"])[0]
	}

	res, err := http.Get("http://www.nseindia.com/corporates/corpInfo/equities/getBoardMeetings.jsp?period=*&symbol=" + comp)
	if err != nil {
		log.Fatal(err)
		return
	}
	
	var s = new(BoardMeetingResponse)
 	dec := json.NewDecoder(res.Body)
 	err = dec.Decode(&s)

 	if err != nil {
 		log.Fatal(err)
 		return	
 	}

	tpl := template.Must(template.ParseFiles("boardmeetings.html"))
	err = tpl.Execute(w, s)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func main(){
	http.Handle("/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/boardmeetings", serverBoardMeetings)
	http.ListenAndServe(":8080", nil)
}