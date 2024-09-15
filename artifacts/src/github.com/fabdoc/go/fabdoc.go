package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	sc "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/common/flogging"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
)

// SmartContract Define the Smart Contract structure
type SmartContract struct {
}

// Doc :  Define the doc structure, with 4 properties.  Structure tags are used by encoding/json library
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
	MotherOccupation  string `json:"motherOccupation"`
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

type docPrivateDetails struct {
	Centre string `json:"centre"`
	MainCentre string `json:"mainCentre"`
}

// Init ;  Method for initializing smart contract
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

var logger = flogging.MustGetLogger("fabdoc_cc")

// Invoke :  Method for INVOKING smart contract
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	function, args := APIstub.GetFunctionAndParameters()
	logger.Infof("Function name is:  %d", function)
	logger.Infof("Args length is : %d", len(args))

	switch function {
	case "queryDoc":
		return s.queryDoc(APIstub, args)
	case "initLedger":
		return s.initLedger(APIstub)
	case "createDoc":
		return s.createDoc(APIstub, args)
	case "queryAllDocs":
		return s.queryAllDocs(APIstub)
	case "changeDocCentre":
		return s.changeDocCentre(APIstub, args)
	case "changeDocStatus":
		return s.changeDocStatus(APIstub, args)
	case "updateDocFather":
		return s.updateDocFather(APIstub, args)	
	case "getHistoryForAsset":
		return s.getHistoryForAsset(APIstub, args)
	case "queryDocsByCentre":
		return s.queryDocsByCentre(APIstub, args)
	case "queryDocsByStatus":
		return s.queryDocsByStatus(APIstub, args)		
	case "restictedMethod":
		return s.restictedMethod(APIstub, args)
	case "test":
		return s.test(APIstub, args)
	case "createPrivateDoc":
		return s.createPrivateDoc(APIstub, args)
	case "readPrivateDoc":
		return s.readPrivateDoc(APIstub, args)
	case "updatePrivateData":
		return s.updatePrivateData(APIstub, args)
	case "readDocPrivateDetails":
		return s.readDocPrivateDetails(APIstub, args)
	case "createPrivateDocImplicitForOrg1":
		return s.createPrivateDocImplicitForOrg1(APIstub, args)
	case "createPrivateDocImplicitForOrg2":
		return s.createPrivateDocImplicitForOrg2(APIstub, args)
	case "queryPrivateDataHash":
		return s.queryPrivateDataHash(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}

	// return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(docAsBytes)
}

func (s *SmartContract) readPrivateDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	// collectionDocs, collectionDocPrivateDetails, _implicit_org_Org1MSP, _implicit_org_Org2MSP
	docAsBytes, err := APIstub.GetPrivateData(args[0], args[1])
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[1] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if docAsBytes == nil {
		jsonResp := "{\"Error\":\"Doc private details does not exist: " + args[1] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(docAsBytes)
}

func (s *SmartContract) readPrivateDocIMpleciteForOrg1(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docAsBytes, _ := APIstub.GetPrivateData("_implicit_org_Org1MSP", args[0])
	return shim.Success(docAsBytes)
}

func (s *SmartContract) readDocPrivateDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docAsBytes, err := APIstub.GetPrivateData("collectionDocPrivateDetails", args[0])

	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get private details for " + args[0] + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if docAsBytes == nil {
		jsonResp := "{\"Error\":\"Marble private details does not exist: " + args[0] + "\"}"
		return shim.Error(jsonResp)
	}
	return shim.Success(docAsBytes)
}

func (s *SmartContract) test(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(docAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	docs := []Doc{		
		Doc{
			Surname: "Lamda", 
			GivenName: "Epsilon", 
			DateBirth: "2001-01-03", 
			PlaceBirth: "Buea", 
			Gender: "Male", 	
			FatherName: "Habib Marwane",	
			FatherBornOn:"1988-07-06",	
			FatherBornAt: "Bamenda", 
			FatherOccupation: "Fonctionnaide de Police",	
			FatherResidence: "Batibo",	
			FatherNationality: "Camerounaise", 
			FatherDocument: "CNI 8610510567 du 17/09/2018",
			MotherName: "Chloe Marker",	
			MotherBornOn: "1999-06-03", 
			MotherBornAt: "Bertoua", 
			MotherOccupation: "Menagere",	
			MotherResidence: "Tiko",	
			MotherNationality: "Camerounaise", 
			MotherDocument: "CNI 000510567 du 12/05/2020",
			Declarer: "Ndourlaye Ngomnoga Wilfred, pere de l'enfant", 
			RegistrationDate: "2001-04-08", 
			Centre: "Tombel", 
			Officer: "Souleymam Bouba", 
			Secretary: "Japsin Vanessa Ashley",
			Status: "1",
			Observations:  "",
		},
	}

	i := 0
	for i < len(docs) {
		docAsBytes, _ := json.Marshal(docs[i])
		APIstub.PutState("DOC"+strconv.Itoa(i), docAsBytes)
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createPrivateDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	type docTransientInput struct {		
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
		FatherDocument		string	`json:"fatherDocument"`
		MotherName   string `json:"motherName"`
		MotherBornOn  string `json:"motherBornOn"`
		MotherBornAt  string `json:"motherBornAt"`
		MotherOccupation  string `json:"motherOccupation"`
		MotherResidence  string `json:"motherResidence"`
		MotherNationality  string `json:"motherNationality"`
		MotherDocument	string	`json:"motherDocument"`
		Declarer   string `json:"declarer"`
		RegistrationDate string `json:"registrationDate"`
		Centre  string `json:"centre"`	
		Officer  string `json:"officer"`
		Secretary  string `json:"secretary"`
		MainCentre string `json:"mainCentre"`
		Status  string `json:"status"`
		Observations  string `json:"observations"`
		Key   string `json:"key"`
		
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	docDataAsBytes, ok := transMap["doc"]
	if !ok {
		return shim.Error("doc must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(docDataAsBytes))

	if len(docDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var docInput docTransientInput
	err = json.Unmarshal(docDataAsBytes, &docInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(docDataAsBytes) + "Error is : " + err.Error())
	}

	logger.Infof("3333")	

	if len(docInput.Key) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(docInput.Surname) == 0 {
		return shim.Error("Surname field must be a non-empty string")
	}
	if len(docInput.GivenName) == 0 {
		return shim.Error("GivenName field must be a non-empty string")
	}
	if len(docInput.DateBirth) == 0 {
		return shim.Error("DateBirth field must be a non-empty string")
	}
	if len(docInput.PlaceBirth) == 0 {
		return shim.Error("PlaceBirth field must be a non-empty string")
	}
	if len(docInput.Gender) == 0 {
		return shim.Error("Gender field must be a non-empty string")
	}

	if len(docInput.MotherName) == 0 {
		return shim.Error("MotherName field must be a non-empty string")
	}
	if len(docInput.MotherBornOn) == 0 {
		return shim.Error("Mother's Born On field must be a non-empty string")
	}
	if len(docInput.MotherBornAt) == 0 {
		return shim.Error("Mother's Born At field must be a non-empty string")
	}
	if len(docInput.MotherResidence) == 0 {
		return shim.Error("Mother's Residence field must be a non-empty string")
	}
	if len(docInput.MotherOccupation) == 0 {
		return shim.Error("Mother's Occupation field must be a non-empty string")
	}
	if len(docInput.MotherNationality) == 0 {
		return shim.Error("Mother's Nationality field must be a non-empty string")
	}
	if len(docInput.MotherDocument) == 0 {
		return shim.Error("Mother's Reference Document field must be a non-empty string")
	}
	if len(docInput.Declarer) == 0 {
		return shim.Error("Declarer field must be a non-empty string")
	}
	if len(docInput.RegistrationDate) == 0 {
		return shim.Error("RegistrationDate field must be a non-empty string")
	}
	if len(docInput.Centre)== 0 {
		return shim.Error("Centre field must be a non-empty string")
	}
	if len(docInput.Officer) == 0 {
		return shim.Error("Officer field must be a non-empty string")
	}
	if len(docInput.Secretary) == 0 {
		return shim.Error("Secretary field must be a non-empty string")
	}
	if len(docInput.MainCentre) == 0 {
		return shim.Error("Main centre field must be a non-empty string")
	}
	
	logger.Infof("444444")

	// ==== Check if doc already exists ====
	docAsBytes, err := APIstub.GetPrivateData("collectionDocs", docInput.Key)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if docAsBytes != nil {
		fmt.Println("This doc already exists: " + docInput.Key)
		return shim.Error("This doc already exists: " + docInput.Key)
	}	

	logger.Infof("55555")

	var doc = Doc{
		Surname: docInput.Surname, 
		GivenName: docInput.GivenName, 
		DateBirth: docInput.DateBirth, 
		PlaceBirth: docInput.PlaceBirth, 
		Gender: docInput.Gender, 				
		FatherName: docInput.FatherName,
		FatherBornOn: docInput.FatherBornOn,
		FatherBornAt: docInput.FatherBornAt,
		FatherOccupation: docInput.FatherOccupation,
		FatherResidence: docInput.FatherResidence,
		FatherNationality: docInput.FatherNationality,	
		FatherDocument: docInput.FatherDocument,	
		MotherName: docInput.MotherName, 
		MotherBornOn: docInput.MotherBornOn, 
		MotherBornAt: docInput.MotherBornAt, 
		MotherOccupation: docInput.MotherOccupation,
		MotherResidence: docInput.MotherResidence,
		MotherNationality:  docInput.MotherNationality,	
		MotherDocument:	docInput.MotherDocument,			
		Declarer: docInput.Declarer, 
		RegistrationDate: docInput.RegistrationDate, 
		Centre: docInput.Centre, 
		Officer: docInput.Officer, 
		Secretary: docInput.Secretary,
		Status: docInput.Status,
	}

	docAsBytes, err = json.Marshal(doc)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutPrivateData("collectionDocs", docInput.Key, docAsBytes)
	if err != nil {
		logger.Infof("6666666")
		return shim.Error(err.Error())
	}

	docPrivateDetails := &docPrivateDetails{Centre: docInput.Centre, MainCentre: docInput.MainCentre}

	docPrivateDetailsAsBytes, err := json.Marshal(docPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionDocPrivateDetails", docInput.Key, docPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(docAsBytes)
}

func (s *SmartContract) updatePrivateData(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	type docTransientInput struct {
		Centre string `json:"centre"`
		MainCentre string `json:"mainCentre"`
		Key   string `json:"key"`
	}
	if len(args) != 0 {
		return shim.Error("1111111----Incorrect number of arguments. Private marble data must be passed in transient map.")
	}

	logger.Infof("11111111111111111111111111")

	transMap, err := APIstub.GetTransient()
	if err != nil {
		return shim.Error("222222 -Error getting transient: " + err.Error())
	}

	docDataAsBytes, ok := transMap["doc"]
	if !ok {
		return shim.Error("doc must be a key in the transient map")
	}
	logger.Infof("********************8   " + string(docDataAsBytes))

	if len(docDataAsBytes) == 0 {
		return shim.Error("333333 -marble value in the transient map must be a non-empty JSON string")
	}

	logger.Infof("2222222")

	var docInput docTransientInput
	err = json.Unmarshal(docDataAsBytes, &docInput)
	if err != nil {
		return shim.Error("44444 -Failed to decode JSON of: " + string(docDataAsBytes) + "Error is : " + err.Error())
	}

	docPrivateDetails := &docPrivateDetails{Centre: docInput.Centre, MainCentre: docInput.MainCentre}

	docPrivateDetailsAsBytes, err := json.Marshal(docPrivateDetails)
	if err != nil {
		logger.Infof("77777")
		return shim.Error(err.Error())
	}

	err = APIstub.PutPrivateData("collectionDocPrivateDetails", docInput.Key, docPrivateDetailsAsBytes)
	if err != nil {
		logger.Infof("888888")
		return shim.Error(err.Error())
	}

	return shim.Success(docPrivateDetailsAsBytes)

}

func (s *SmartContract) createDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 27 {
		return shim.Error("Incorrect number of arguments.Doc must Expecting 27")
	}
	  
	var doc = Doc{
		Surname: 			args[1], 
		GivenName: 			args[2], 
		DateBirth: 			args[3], 
		PlaceBirth: 		args[4], 
		Gender: 			args[5], 						
		FatherName: 		args[6],
		FatherBornOn: 		args[7],
		FatherBornAt: 		args[8],
		FatherOccupation: 	args[9],
		FatherResidence: 	args[10],
		FatherNationality: args[11],	
		FatherDocument: 	args[12],	
		MotherName: 		args[13],
		MotherBornOn: 		args[14],
		MotherBornAt: 		args[15],
		MotherOccupation: 	args[16],
		MotherResidence: 	args[17],
		MotherNationality:  args[18],
		MotherDocument:		args[19],					
		Declarer: 			args[20],
		RegistrationDate: 	args[21],
		Centre:				args[22],
		Officer: 			args[23],
		Secretary: 			args[24],
		Status:				args[25],
		Observations:		args[26],
	}
	
	docAsBytes, _ := json.Marshal(doc)
	APIstub.PutState(args[0], docAsBytes)

	indexName := "centre~key"
	docNameIndexKey, err := APIstub.CreateCompositeKey(indexName, []string{doc.Centre, args[0]})
	if err != nil {
		return shim.Error(err.Error())
	}
	value := []byte{0x00}
	APIstub.PutState(docNameIndexKey, value)

	return shim.Success(docAsBytes)
}

func (s *SmartContract) queryDocsByCentre(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Nombre incorrect d'arguments. Un seul argument attendu.")
	}
	centre := args[0]

	centreAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("centre~key", []string{centre})
	if err != nil {
		return shim.Error("Erreur lors de la récupération des documents : " + err.Error())
	}
	defer centreAndIdResultIterator.Close()

	// Structure pour contenir l'ID et les données du document
	type DocWithID struct {
		Key    string `json:"Key"`
		Record Doc    `json:"Record"`
	}

	var docs []DocWithID

	for centreAndIdResultIterator.HasNext() {
		responseRange, err := centreAndIdResultIterator.Next()
		if err != nil {
			return shim.Error("Erreur lors de la lecture des résultats : " + err.Error())
		}

		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error("Erreur lors du découpage de la clé composite : " + err.Error())
		}

		id := compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)
		if err != nil {
			return shim.Error("Erreur lors de la récupération de l'état du document : " + err.Error())
		}

		if assetAsBytes == nil {
			return shim.Error("Document non trouvé pour l'ID : " + id)
		}

		var doc Doc
		err = json.Unmarshal(assetAsBytes, &doc)
		if err != nil {
			return shim.Error("Erreur lors de la désérialisation du document : " + err.Error())
		}

		docWithID := DocWithID{
			Key:    id,
			Record: doc,
		}
		docs = append(docs, docWithID)
	}

	docsAsBytes, err := json.Marshal(docs)
	if err != nil {
		return shim.Error("Erreur lors de la sérialisation des résultats : " + err.Error())
	}

	return shim.Success(docsAsBytes)
}
// func (S *SmartContract) queryDocsByCentre(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments")
// 	}
// 	centre := args[0]

// 	centreAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("centre~key", []string{centre})
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	defer centreAndIdResultIterator.Close()

// 	var i int
// 	var id string

// 	var docs []byte
// 	bArrayMemberAlreadyWritten := false

// 	docs = append([]byte("["))

// 	for i = 0; centreAndIdResultIterator.HasNext(); i++ {
// 		responseRange, err := centreAndIdResultIterator.Next()
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}

// 		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}

// 		id = compositeKeyParts[1]
// 		assetAsBytes, err := APIstub.GetState(id)

// 		if bArrayMemberAlreadyWritten == true {
// 			newBytes := append([]byte(","), assetAsBytes...)
// 			docs = append(docs, newBytes...)

// 		} else {
// 			// newBytes := append([]byte(","), docsAsBytes...)
// 			docs = append(docs, assetAsBytes...)
// 		}

// 		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
// 		bArrayMemberAlreadyWritten = true

// 	}

// 	docs = append(docs, []byte("]")...)

// 	return shim.Success(docs)
// }

// func (S *SmartContract) queryDocsByStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments")
// 	}
// 	status := args[0]

// 	centreAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("status~key", []string{status})
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	defer centreAndIdResultIterator.Close()

// 	var i int
// 	var id string

// 	var docs []byte
// 	bArrayMemberAlreadyWritten := false

// 	docs = append([]byte("["))

// 	for i = 0; centreAndIdResultIterator.HasNext(); i++ {
// 		responseRange, err := centreAndIdResultIterator.Next()
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}

// 		objectType, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
// 		if err != nil {
// 			return shim.Error(err.Error())
// 		}

// 		id = compositeKeyParts[1]
// 		assetAsBytes, err := APIstub.GetState(id)

// 		if bArrayMemberAlreadyWritten == true {
// 			newBytes := append([]byte(","), assetAsBytes...)
// 			docs = append("{\"Key\":")
// 			docs = append("\"")
// 			docs = append(id) 
// 			docs = append("\"")

// 			docs = append(", \"Record\":")
// 			docs = append(docs, newBytes...)
// 			docs = append("}")
// 			//docs = append("\"")

// 		} else {
// 			// newBytes := append([]byte(","), docsAsBytes...)
// 			docs = append(docs, assetAsBytes...)
// 		}

// 		fmt.Printf("Found a asset for index : %s asset id : ", objectType, compositeKeyParts[0], compositeKeyParts[1])
// 		bArrayMemberAlreadyWritten = true

// 	}

// 	docs = append(docs, []byte("]")...)

// 	return shim.Success(docs)
// }

func (s *SmartContract) queryDocsByStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Nombre incorrect d'arguments. Un seul argument attendu.")
	}
	status := args[0]

	statusAndIdResultIterator, err := APIstub.GetStateByPartialCompositeKey("status~key", []string{status})
	if err != nil {
		return shim.Error("Erreur lors de la récupération des documents : " + err.Error())
	}
	defer statusAndIdResultIterator.Close()

	// Structure pour contenir l'ID et les données du document
	type DocWithID struct {
		Key    string `json:"Key"`
		Record Doc    `json:"Record"`
	}

	var docs []DocWithID

	for statusAndIdResultIterator.HasNext() {
		responseRange, err := statusAndIdResultIterator.Next()
		if err != nil {
			return shim.Error("Erreur lors de la lecture des résultats : " + err.Error())
		}

		_, compositeKeyParts, err := APIstub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error("Erreur lors du découpage de la clé composite : " + err.Error())
		}

		id := compositeKeyParts[1]
		assetAsBytes, err := APIstub.GetState(id)
		if err != nil {
			return shim.Error("Erreur lors de la récupération de l'état du document : " + err.Error())
		}

		if assetAsBytes == nil {
			return shim.Error("Document non trouvé pour l'ID : " + id)
		}

		var doc Doc
		err = json.Unmarshal(assetAsBytes, &doc)
		if err != nil {
			return shim.Error("Erreur lors de la désérialisation du document : " + err.Error())
		}

		docWithID := DocWithID{
			Key:    id,
			Record: doc,
		}
		docs = append(docs, docWithID)
	}

	docsAsBytes, err := json.Marshal(docs)
	if err != nil {
		return shim.Error("Erreur lors de la sérialisation des résultats : " + err.Error())
	}

	return shim.Success(docsAsBytes)
}

func (s *SmartContract) queryAllDocs(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "DOC0"
	endKey := "DOC9999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllDocs:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) restictedMethod(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	// get an ID for the client which is guaranteed to be unique within the MSP
	//id, err := cid.GetID(APIstub) -

	// get the MSP ID of the client's identity
	//mspid, err := cid.GetMSPID(APIstub) -

	// get the value of the attribute
	//val, ok, err := cid.GetAttributeValue(APIstub, "attr1") -

	// get the X509 certificate of the client, or nil if the client's identity was not based on an X509 certificate
	//cert, err := cid.GetX509Certificate(APIstub) -

	val, ok, err := cid.GetAttributeValue(APIstub, "role")
	if err != nil {
		// There was an error trying to retrieve the attribute
		shim.Error("Error while retriving attributes")
	}
	if !ok {
		// The client identity does not possess the attribute
		shim.Error("Client identity doesnot posses the attribute")
	}
	// Do something with the value of 'val'
	if val != "approver" {
		fmt.Println("Attribute role: " + val)
		return shim.Error("Only user with role as APPROVER have access this method!")
	} else {
		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting 1")
		}

		docAsBytes, _ := APIstub.GetState(args[0])
		return shim.Success(docAsBytes)
	}

}

func (s *SmartContract) changeDocCentre(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	docAsBytes, _ := APIstub.GetState(args[0])
	doc := Doc{}

	json.Unmarshal(docAsBytes, &doc)
	doc.Centre = args[1]

	docAsBytes, _ = json.Marshal(doc)
	APIstub.PutState(args[0], docAsBytes)

	return shim.Success(docAsBytes)
}

func (s *SmartContract) changeDocStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	docAsBytes, _ := APIstub.GetState(args[0])
	doc := Doc{}

	json.Unmarshal(docAsBytes, &doc)
	doc.Officer = args[1]
	doc.Status = args[2]
	doc.Observations = args[3]

	docAsBytes, _ = json.Marshal(doc)
	APIstub.PutState(args[0], docAsBytes)

	return shim.Success(docAsBytes)
}

func (s *SmartContract) updateDocFather(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 9 {
		return shim.Error("Incorrect number of arguments. Expecting 9")
	}

	docAsBytes, _ := APIstub.GetState(args[0])
	doc := Doc{}

	json.Unmarshal(docAsBytes, &doc)
	//doc.Status = args[1]
	doc.FatherName = args[1]
	doc.FatherBornOn = args[2]
	doc.FatherBornAt = args[3]
	doc.FatherOccupation = args[4]
	doc.FatherResidence = args[5]
	doc.FatherNationality = args[6]
	doc.FatherDocument = args[7]
	doc.Observations = args[8]

	docAsBytes, _ = json.Marshal(doc)
	APIstub.PutState(args[0], docAsBytes)

	return shim.Success(docAsBytes)
}

func (t *SmartContract) getHistoryForAsset(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docName := args[0]

	resultsIterator, err := stub.GetHistoryForKey(docName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForAsset returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) createPrivateDocImplicitForOrg1(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {	
	if len(args) != 27 {
		return shim.Error("Incorrect arguments. Expecting 27 arguments to cpOrg1")
	}
 	  
	var doc = Doc{
		Surname: 			args[1], 
		GivenName: 			args[2], 
		DateBirth: 			args[3], 
		PlaceBirth: 		args[4], 
		Gender: 			args[5], 						
		FatherName: 		args[6],
		FatherBornOn: 		args[7],
		FatherBornAt: 		args[8],
		FatherOccupation: 	args[9],
		FatherResidence: 	args[10],
		FatherNationality: 	args[11],	
		FatherDocument: 	args[12],	
		MotherName: 		args[13],
		MotherBornOn: 		args[14],
		MotherBornAt: 		args[15],
		MotherOccupation: 	args[16],
		MotherResidence: 	args[17],
		MotherNationality:  args[18],
		MotherDocument:		args[19],					
		Declarer: 			args[20],
		RegistrationDate: 	args[21],
		Centre:				args[22],
		Officer: 			args[23],
		Secretary: 			args[24],		
		Status:				args[25],
		Observations:  		args[26],
	}

	docAsBytes, _ := json.Marshal(doc)
	// APIstub.PutState(args[0], docAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_Org1MSP", args[0], docAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(docAsBytes)
}

func (s *SmartContract) createPrivateDocImplicitForOrg2(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 27 {
		return shim.Error("Incorrect arguments. Expecting 27 arguments to cpOrg2")
	}

	var doc = Doc{
		Surname: 			args[1], 
		GivenName: 			args[2], 
		DateBirth: 			args[3], 
		PlaceBirth: 		args[4], 
		Gender: 			args[5], 						
		FatherName: 		args[6],
		FatherBornOn: 		args[7],
		FatherBornAt: 		args[8],
		FatherOccupation: 	args[9],
		FatherResidence: 	args[10],
		FatherNationality: 	args[11],	
		FatherDocument: 	args[12],	
		MotherName: 		args[13],
		MotherBornOn: 		args[14],
		MotherBornAt: 		args[15],
		MotherOccupation: 	args[16],
		MotherResidence: 	args[17],
		MotherNationality:  args[18],
		MotherDocument:		args[19],					
		Declarer: 			args[20],
		RegistrationDate: 	args[21],
		Centre:				args[22],
		Officer: 			args[23],
		Secretary: 			args[24],		
		Status:				args[25],
		Observations:  		args[26],
	}

	
	docAsBytes, _ := json.Marshal(doc)
	APIstub.PutState(args[0], docAsBytes)

	err := APIstub.PutPrivateData("_implicit_org_Org2MSP", args[0], docAsBytes)
	if err != nil {
		return shim.Error("Failed to add asset: " + args[0])
	}
	return shim.Success(docAsBytes)
}

func (s *SmartContract) queryPrivateDataHash(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	docAsBytes, _ := APIstub.GetPrivateDataHash(args[0], args[1])
	return shim.Success(docAsBytes)
}

// func (s *SmartContract) CreateDocAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	if len(args) != 1 {
// 		return shim.Error("Incorrect number of arguments. Expecting 1")
// 	}

// 	var doc Doc
// 	err := json.Unmarshal([]byte(args[0]), &doc)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	docAsBytes, err := json.Marshal(doc)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	err = APIstub.PutState(doc.ID, docAsBytes)
// 	if err != nil {
// 		return shim.Error(err.Error())
// 	}

// 	return shim.Success(nil)
// }

// func (s *SmartContract) addBulkAsset(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
// 	logger.Infof("Function addBulkAsset called and length of arguments is:  %d", len(args))
// 	if len(args) >= 500 {
// 		logger.Errorf("Incorrect number of arguments in function CreateAsset, expecting less than 500, but got: %b", len(args))
// 		return shim.Error("Incorrect number of arguments, expecting 2")
// 	}

// 	var eventKeyValue []string

// 	for i, s := range args {

// 		key :=s[0];
// 		var doc = Doc{Make: s[1], Model: s[2], Colour: s[3], Centre: s[4]}

// 		eventKeyValue = strings.SplitN(s, "#", 3)
// 		if len(eventKeyValue) != 3 {
// 			logger.Errorf("Error occured, Please make sure that you have provided the array of strings and each string should be  in \"EventType#Key#Value\" format")
// 			return shim.Error("Error occured, Please make sure that you have provided the array of strings and each string should be  in \"EventType#Key#Value\" format")
// 		}

// 		assetAsBytes := []byte(eventKeyValue[2])
// 		err := APIstub.PutState(eventKeyValue[1], assetAsBytes)
// 		if err != nil {
// 			logger.Errorf("Error coocured while putting state for asset %s in APIStub, error: %s", eventKeyValue[1], err.Error())
// 			return shim.Error(err.Error())
// 		}
// 		// logger.infof("Adding value for ")
// 		fmt.Println(i, s)

// 		indexName := "Event~Id"
// 		eventAndIDIndexKey, err2 := APIstub.CreateCompositeKey(indexName, []string{eventKeyValue[0], eventKeyValue[1]})

// 		if err2 != nil {
// 			logger.Errorf("Error coocured while putting state in APIStub, error: %s", err.Error())
// 			return shim.Error(err2.Error())
// 		}

// 		value := []byte{0x00}
// 		err = APIstub.PutState(eventAndIDIndexKey, value)
// 		if err != nil {
// 			logger.Errorf("Error coocured while putting state in APIStub, error: %s", err.Error())
// 			return shim.Error(err.Error())
// 		}
// 		// logger.Infof("Created Composite key : %s", eventAndIDIndexKey)

// 	}

// 	return shim.Success(nil)
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
