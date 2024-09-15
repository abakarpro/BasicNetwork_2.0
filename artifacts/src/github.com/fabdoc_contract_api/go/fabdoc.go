package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Doc struct {	
	Surname   string `json:"surname"`
	GivenName  string `json:"givenName"`
	DateBirth string `json:"dateBirth"`
	PlaceBirth  string `json:"placeBirth"`
	Gender   string `json:"gender"`	

	FatherName   string `json:"fatherName"`
	FatherBornOn  string `json:"fatherBornOn"`
	FatherBornAt  string `json:"fatherBornAt"`
	FatherOccupation  string `json:"fatherOccupation"`
	FatherResidence  string `json:"fatherResidence"`
	FatherNationality  string `json:"fatherNationality"`
	FatherDocument  string `json:"fatherDocument"`

	MotherName   string `json:"motherName"`
	MotherBornOn  string `json:"motherBornOn"`
	MotherBornAt  string `json:"motherBornAt"`
	MotherOccupation  string `json:"motherccupation"`
	MotherResidence  string `json:"motherResidence"`
	MotherNationality  string `json:"motherNationality"`
	MotherDocument  string `json:"motherDocument"`

	Declarer   string `json:"declarer"`
	RegistrationDate string `json:"registrationDate"`
	Centre  string `json:"centre"`	
	Officer  string `json:"officer"`
	Secretary  string `json:"secretary"`
		
	Status  string `json:"status"`	
	Observations  string `json:"observations"`

	}


type QueryResult struct {
	Key    string `json:"Key"`
	Record *Doc
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	docs := []Doc{		
		Doc{Surname: "Lamda", GivenName: "Epsilon", DateBirth: "2001-01-03", PlaceBirth: "Buea", Gender: "Male", 
		FatherName:"Christopher Le",		FatherBornOn:"1956-09-12",		FatherBornAt:"fatherBornAt",
		FatherOccupation:"Enseignant",		FatherResidence:"Buea",		FatherNationality:"Graven James", FatherDocument:"Coco Jasmine",

		MotherName:"motherName",		MotherBornOn:"1968-04-02",		MotherBornAt:"Lagdo",
		MotherOccupation:"Enseignante",		MotherResidence: "Tombel",		MotherNationality:"Camerounaise", 	MotherDocument:"Laylah Amirah",
		
		Declarer: "Ndourlaye", RegistrationDate: "2001-04-08", Centre: "Tombel", Officer: "officer", Secretary: "Henriette Fontcha",  Status: "1"},
		
		Doc{Surname: "Omega", GivenName: "Alpha", DateBirth: "2002-01-07", PlaceBirth: "Bertoua", Gender: "Male", 
		FatherName:"Frendy Gaigama",		FatherBornOn:"1975-08-03",
		FatherBornAt:"Makari",		FatherOccupation:"Diplomate",
		FatherResidence:"Goulfey",		FatherNationality:"Camerounaise", FatherDocument:"Hammid Zain",

		MotherName:"Leonna Gustive",		MotherBornOn:"1988-03-07",
		MotherBornAt:"Tokombere",		MotherOccupation:"Architecte",
		MotherResidence: "Kribi",		MotherNationality:"Camerounaise", 	MotherDocument:"Farida Symouine",
		
		Declarer: "Nyangono Fils", RegistrationDate: "2001-05-10", Centre: "Tiko", Officer: "Ndouvla Jean", Secretary: "Yvonne Tchang",  Status: "0"},
		
		Doc{Surname: "Iota", GivenName: "Delta", DateBirth: "2000-04-08", PlaceBirth: "Kousseri", Gender: "Female", 
		FatherName:"Tchiroma Gregory",		FatherBornOn:"1973-04-07",
		FatherBornAt:"Maroua",		FatherOccupation:"Informaticien",
		FatherResidence:"Douala",		FatherNationality:"Camerounaise", FatherDocument:"Camerounaise",

		MotherName:"Sandrine Gold",		MotherBornOn:"1983-06-01",
		MotherBornAt:"Pointe-Noire",		MotherOccupation:"Journaliste",
		MotherResidence: "Yaounde",		MotherNationality:"Gabonaise", 	MotherDocument:"Camerounaise",
		
		Declarer: "Sambero Jean", RegistrationDate: "2002-02-11", Centre: "Kousseri", Officer: "officer", Secretary: "secretary3", Status: "1"},
	}

	for i, doc := range docs {
		docAsBytes, _ := json.Marshal(doc)
		err := ctx.GetStub().PutState("DOC"+strconv.Itoa(i), docAsBytes)

		if err != nil {
			return fmt.Errorf("Failed to put to world state. %s", err.Error())
		}
	}

	return nil
}

func (s *SmartContract) CreateDoc(ctx contractapi.TransactionContextInterface, docNumber string, surname string, givenName string, dateBirth string, placeBirth string, gender string, 
		fatherName string, fatherBornOn string, fatherBornAt string, fatherOccupation string, fatherResidence string, fatherNationality string, fatherDocument string,
		motherName string,		motherBornOn string,	motherBornAt string,	motherOccupation string, motherResidence string, motherNationality string, motherDocument string,
	declarer string, registrationDate string, centre string, officer string, secretary string, status string, observations string) error {
	doc := Doc{		
		Surname:   surname,
		GivenName:  givenName,
		DateBirth: dateBirth,
		PlaceBirth:  placeBirth,
		Gender:   gender,

		FatherName:	fatherName,
		FatherBornOn:	fatherBornOn,
		FatherBornAt:	fatherBornAt,
		FatherOccupation:	fatherOccupation,
		FatherResidence:	fatherResidence,
		FatherNationality:	fatherNationality,
		FatherDocument:		fatherDocument,

		MotherName:	motherName,
		MotherBornOn:	motherBornOn,
		MotherBornAt:	motherBornAt,
		MotherOccupation:	motherOccupation,
		MotherResidence:	motherResidence,
		MotherNationality:	motherNationality,
		MotherDocument:		motherDocument,

		Declarer:	declarer,
		RegistrationDate: registrationDate,
		Centre:  centre,
		Officer: officer,
		Secretary: secretary,

		Status: status,
		Observations: observations,
	}

	docAsBytes, _ := json.Marshal(doc)

	return ctx.GetStub().PutState(docNumber, docAsBytes)
}

func (s *SmartContract) QueryDoc(ctx contractapi.TransactionContextInterface, docNumber string) (*Doc, error) {
	docAsBytes, err := ctx.GetStub().GetState(docNumber)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if docAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", docNumber)
	}

	doc := new(Doc)
	_ = json.Unmarshal(docAsBytes, doc)

	return doc, nil
}

func (s *SmartContract) QueryAllDocs(ctx contractapi.TransactionContextInterface) ([]QueryResult, error) {
	startKey := "DOC0"
	endKey := "DOC9999"

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []QueryResult{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}

		doc := new(Doc)
		_ = json.Unmarshal(queryResponse.Value, doc)

		queryResult := QueryResult{Key: queryResponse.Key, Record: doc}
		results = append(results, queryResult)
	}

	return results, nil
}

func (s *SmartContract) ChangeDocCentre(ctx contractapi.TransactionContextInterface, docNumber string, newCentre string) error {
	doc, err := s.QueryDoc(ctx, docNumber)

	if err != nil {
		return err
	}

	doc.Centre = newCentre

	docAsBytes, _ := json.Marshal(doc)

	return ctx.GetStub().PutState(docNumber, docAsBytes)
}

func (s *SmartContract) ChangeDocStatus(ctx contractapi.TransactionContextInterface, docNumber string, newStatus string, observations string) error {
	doc, err := s.QueryDoc(ctx, docNumber)

	if err != nil {
		return err
	}

	doc.Status = newStatus
	doc.Observations = observations

	docAsBytes, _ := json.Marshal(doc)

	return ctx.GetStub().PutState(docNumber, docAsBytes)
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabdoc chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabdoc chaincode: %s", err.Error())
	}
}
