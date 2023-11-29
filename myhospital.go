/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const HOSPITAL = "HOSPITAL"
const PATIENT = "PATIENT"
const HOSPITAL_ADMIN = "HOSPITAL_ADMIN"
const ADMIN = "admin"
const DOCTOR = "DOCTOR"
const Org1MSP = "Org1MSP"
const Org2MSP = "Org2MSP"
const PRESCRIPTION = "PRESCRIPTION"
const idDoctypeIndex = "id~doctype"
const id_DrId_DoctypeIndex = "id~drId~doctype"

// SmartContract provides functions for managing a Asset and Token
type SmartContract struct {
	contractapi.Contract
}

type Hospital struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	DocType         string `json:"docType"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Pincode         string `json:"pincode"`
	RegistrationNum string `json:"registrationNum"`
}

type HospitalAdmin struct {
	Id               string `json:"id"`
	FirstName        string `json:"firstName,omitempty"`
	MiddleName       string `json:"middleName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	Dob              string `json:"dob,omitempty"`
	ContactNo        string `json:"contactNo,omitempty"`
	EmergencyNo      string `json:"emergencyNo,omitempty"`
	PermanentAddress string `json:"permanentAddress,omitempty"`
	EmailId          string `json:"mailId,omitempty"`
	BloodGroup       string `json:"bloodGrp,omitempty"`
	DocType          string `json:"docType,omitempty"`
	Type             string `json:"type,omitempty"` //Employee
	HospitalId       string `json:"hospitalId"`
	Active           bool   `json:"active"`
}

type User struct {
	Id               string `json:"id"`
	FirstName        string `json:"firstName,omitempty"`
	MiddleName       string `json:"middleName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	ContactNo        string `json:"contactNo,omitempty"`
	EmergencyNo      string `json:"emergencyNo,omitempty"`
	LocalAddress     string `json:"localAddress,omitempty"`
	PermanentAddress string `json:"permanentAddress,omitempty"`
	EmailId          string `json:"mailId,omitempty"`
	Dob              string `json:"dob,omitempty"`
	//	Age              int         `json:"age,omitempty"`
	BloodGroup string      `json:"bloodGrp,omitempty"`
	HospitalId string      `json:"hospitalId,omitempty"`
	DocType    string      `json:"docType"`
	Active     bool        `json:"active,omitempty"`
	MetaData   interface{} `json:"annexture,omitempty"`

	// Allergies        string `json:"allergies,omitempty"`
	// Comorbidity      string `json:"comorbidity,omitempty"`
	// Symptoms         string `json:"symptoms,omitempty"`
	// CurrentTreatment string `json:"currentTreatment,omitempty"`
}

type Medicine struct {
	Name      string `json:"name"`
	Dosage    string `json:"dosage"`
	Frequency string `json:"frequency"`
	Remarks   string `json:"remarks"`
}

type Prescription struct {
	//Date             string    `json:"date"`
	Timestamp      string     `json:"timestamp,omitempty"`
	DocType        string     `json:"doctype,omitempty"`
	DoctorId       string     `json:"doctorId,omitempty"`
	PatientId      string     `json:"patientId"`
	Desc           string     `json:"desc"`
	MedicineRecord []Medicine `json:"medicine"`
}

func (s *SmartContract) AddHospital(ctx contractapi.TransactionContextInterface, hospitalInputString string) error {
	var hospitalInput Hospital
	err := json.Unmarshal([]byte(hospitalInputString), &hospitalInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", hospitalInput)

	//Validate super admin
	superAdminIdentity, err := getUserIdentityName(ctx)
	fmt.Println("superAdminIdentity :", superAdminIdentity)

	if superAdminIdentity != ADMIN {
		return fmt.Errorf("permission denied: only admin can call this function")
	}

	//Check if hospital is present or not
	hospitalDetailer, err := getEntityDetails(ctx, hospitalInput.Id, hospitalInput.DocType)
	if err != nil {
		return err
	}
	if hospitalDetailer != nil {
		return fmt.Errorf("Hospital already exist with id : %v", hospitalInput.Id)
	}

	hospitalBytes, err := json.Marshal(hospitalInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Hospital records : %v", err.Error())

	}

	//Inserting hospital record
	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{hospitalInput.Id, hospitalInput.DocType})
	if err != nil {
		return fmt.Errorf("failed to create composite key for hospital %v and err is :%v", hospitalInput.Id, err.Error())
	}

	err = ctx.GetStub().PutState(compositeKey, hospitalBytes)
	if err != nil {
		return fmt.Errorf("failed to insert hospital details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")
	return nil
}

func (s *SmartContract) AddHospitalAdmin(ctx contractapi.TransactionContextInterface, hospitalAdminInputString string) error {
	var hospitalAdminInput HospitalAdmin
	err := json.Unmarshal([]byte(hospitalAdminInputString), &hospitalAdminInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", hospitalAdminInput)

	//validating super admin identity
	superAdminIdentity, err := getUserIdentityName(ctx)
	fmt.Println("superAdminIdentity :", superAdminIdentity)
	if superAdminIdentity != ADMIN {
		return fmt.Errorf("permission denied: only admin can call this function")
	}

	//Check if hospital id is present or not
	hospitalDetailer, err := getEntityDetails(ctx, hospitalAdminInput.HospitalId, HOSPITAL)
	if err != nil {
		return err
	}
	if hospitalDetailer == nil {
		return fmt.Errorf("Hospital does not exist with Id : %v", hospitalAdminInput.HospitalId)
	}

	//Check if hospital admin is present or not
	hospitalAdminDetailer, err := getEntityDetails(ctx, hospitalAdminInput.Id, hospitalAdminInput.DocType)
	if err != nil {
		return err
	}
	if hospitalAdminDetailer != nil {
		return fmt.Errorf("Hospital admin already exist with id : %v", hospitalAdminInput.Id)
	}

	hospitalAdminBytes, err := json.Marshal(hospitalAdminInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Hospital admin records : %v", err.Error())
	}

	//Inserting hospital admin record
	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{hospitalAdminInput.Id, hospitalAdminInput.DocType})
	if err != nil {
		return fmt.Errorf("failed to create composite key for hospital admin %v and err is :%v", hospitalAdminInput.Id, err.Error())
	}
	err = ctx.GetStub().PutState(compositeKey, hospitalAdminBytes)
	if err != nil {
		return fmt.Errorf("failed to insert hospital admin details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")
	return nil
}

func (s *SmartContract) AddDoctor(ctx contractapi.TransactionContextInterface, doctorInputString string) error {
	var doctorInput User
	err := json.Unmarshal([]byte(doctorInputString), &doctorInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", doctorInput)

	// hospitalAdminRole, _, err := getCertificateAttributeValue(ctx, "userRole")
	// fmt.Println("Attribute userRole value :", hospitalAdminRole)
	// if hospitalAdminRole != HOSPITAL_ADMIN {
	// 	return fmt.Errorf("Only Hospital Admin are allowed to register doctor")
	// }

	// adminHospitalName, _, err := getCertificateAttributeValue(ctx, "organizationName")
	// fmt.Println("Attribute organizationName value :", adminHospitalName)

	attributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization"})
	if err != nil {
		return err
	}
	fmt.Println("userRole :", attributes["userRole"])
	fmt.Println("organization :", attributes["organization"])

	if attributes["userRole"] != HOSPITAL_ADMIN {
		return fmt.Errorf("Only Hospital Admin are allowed to register doctor")
	}

	hospitalAdminMSPOrgId, err := loggedInUserMSPID(ctx)
	if err != nil {
		return err
	}
	fmt.Println("hospitalAdminMSPOrgId :", hospitalAdminMSPOrgId)
	if hospitalAdminMSPOrgId != Org1MSP {
		return fmt.Errorf("OrgMSP Id is different")
	}

	//Check if hospital admin is present or not
	doctorDetailer, err := getEntityDetails(ctx, doctorInput.Id, doctorInput.DocType)
	if err != nil {
		return err
	}
	if doctorDetailer != nil {
		return fmt.Errorf("Doctor %v is already registered with hospital.", doctorInput.Id)
	}

	//Assigning Hospital id
	doctorInput.HospitalId = attributes["organization"]

	doctorBytes, err := json.Marshal(doctorInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Doctor records : %v", err.Error())
	}

	//Inserting hospital admin record
	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{doctorInput.Id, doctorInput.DocType})
	if err != nil {
		return fmt.Errorf("failed to create composite key for hospital doctor %v and err is :%v", doctorInput.Id, err.Error())
	}
	err = ctx.GetStub().PutState(compositeKey, doctorBytes)
	if err != nil {
		return fmt.Errorf("failed to insert doctor details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")
	return nil
}

func (s *SmartContract) AddEntity(ctx contractapi.TransactionContextInterface, entityInputString string) error {
	var entityInput User
	err := json.Unmarshal([]byte(entityInputString), &entityInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", entityInput)

	attributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization"})
	if err != nil {
		return err
	}
	fmt.Println("userRole :", attributes["userRole"])
	fmt.Println("organization :", attributes["organization"])

	if attributes["userRole"] != HOSPITAL_ADMIN {
		return fmt.Errorf("Only Hospital Admin are allowed to register %v", entityInput.DocType)
	}

	hospitalAdminMSPOrgId, err := loggedInUserMSPID(ctx)
	if err != nil {
		return err
	}
	fmt.Println("hospitalAdminMSPOrgId :", hospitalAdminMSPOrgId)
	if hospitalAdminMSPOrgId != Org1MSP {
		return fmt.Errorf("OrgMSP Id is different")
	}

	//Check if entity is present or not
	entityDetailer, err := getEntityDetails(ctx, entityInput.Id, entityInput.DocType)
	if err != nil {
		return err
	}
	if entityDetailer != nil {
		return fmt.Errorf("%v is already registered with hospital.", entityInput.Id)
	}

	//Assigning Hospital id
	entityInput.HospitalId = attributes["organization"]

	entityBytes, err := json.Marshal(entityInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of entity records : %v", err.Error())
	}

	//Inserting entity record
	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{entityInput.Id, entityInput.DocType})
	if err != nil {
		return fmt.Errorf("failed to create composite key %v and err is :%v", entityInput.Id, err.Error())
	}
	err = ctx.GetStub().PutState(compositeKey, entityBytes)
	if err != nil {
		return fmt.Errorf("failed to insert entity details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")
	return nil
}

// Dr may be act as patient in hospital, pass doctype also
func (s *SmartContract) ViewOwnDetails(ctx contractapi.TransactionContextInterface) (string, error) {
	// userMSPOrgId, err := loggedInUserMSPID(ctx)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("userMSPOrgId :", userMSPOrgId)
	// if userMSPOrgId != Org2MSP {
	// 	return fmt.Errorf("OrgMSP Id is different")
	// }

	userIdentity, err := getUserIdentityName(ctx)
	fmt.Println("userIdentity :", userIdentity)
	if err != nil {
		return "", err
	}

	docType, _, err := getCertificateAttributeValue(ctx, "userRole")
	fmt.Println("Attribute docType value :", docType)

	//fetching user record from couchDB
	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{userIdentity, docType})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key for user %v and err is :%v", userIdentity, err.Error())
	}
	objectBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return "", fmt.Errorf("Failed to read data from world state %s", err.Error())
	}
	if objectBytes == nil {
		return "", fmt.Errorf("Record for %v user does not exist", userIdentity)
	}
	fmt.Println("objectBytes : ", string(objectBytes))

	return string(objectBytes), nil
}

func (s *SmartContract) ViewPatientDetails(ctx contractapi.TransactionContextInterface, inputString string) (string, error) {
	input := struct {
		ViewerId string `json:"viewerId"`
		UserId   string `json:"userId"`
	}{}
	err := json.Unmarshal([]byte(inputString), &input)
	if err != nil {
		return "", fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", input)

	//Validate ViewerId attributes
	drAttributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization", "orgRole"})
	if err != nil {
		return "", err
	}
	fmt.Println("userRole :", drAttributes["userRole"])
	fmt.Println("organization :", drAttributes["organization"])
	fmt.Println("orgRole :", drAttributes["orgRole"])

	//fetching user record from couchDB
	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{input.UserId, PATIENT})
	if err != nil {
		return "", fmt.Errorf("failed to create composite key for user %v and err is :%v", input.UserId, err.Error())
	}
	userBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		fmt.Println("Failed to read data from world state %s", err.Error())
		return "", fmt.Errorf("Failed to read data from world state %s", err.Error())
	}
	if userBytes == nil {
		fmt.Println("No Data found")
		return "", fmt.Errorf("Record for %v patient does not exist", input.UserId)
	}
	fmt.Println("userBytes : ", string(userBytes))

	var patientDetail User
	err = json.Unmarshal(userBytes, &patientDetail)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal user data: %s", err.Error())
	}
	fmt.Println("User Data : ", patientDetail)

	//Validate attributes
	if drAttributes["organization"] != patientDetail.HospitalId || drAttributes["userRole"] != DOCTOR || patientDetail.DocType != PATIENT {
		return "", fmt.Errorf("Ony doctor are authorized to see patient data")
	}

	return string(userBytes), nil

}

func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, prescriptionInputString string) error {

	// prescriptionInput := struct {
	// 	PatientId        string `json:"patientId"`
	// 	PrescriptionDesc string `json:"patientDesc"`
	// 	// MedName          string `json:"medicineName"`
	// 	// MedQuantity      int    `json:"medicineQuantity"`
	// 	// MedFrequency     int    `json:"medicineFrequency"`
	// 	// MedRemarks       string `json:"medicineRemarks"`
	// 	//MedicineDetails []interface{} `json:"medicineDetails"`
	// 	MedicineDetails []Medicine `json:"medicineDetails"`
	// }{}
	var prescriptionInput Prescription
	err := json.Unmarshal([]byte(prescriptionInputString), &prescriptionInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of prescription input string : %v", err.Error())
	}
	fmt.Println("Input String prescriptionInputString:", prescriptionInputString)
	//fmt.Println("Input String :", prescriptionInput.MedicineDetails)

	// Convert data to slice of Medicine structs
	//medicines := make([]Medicine, len(prescriptionInput.MedicineDetails))
	// var medicines []Medicine

	// for i, medDetail := range prescriptionInput.MedicineDetails {
	// 	medMap, ok := medDetail.(map[string]interface{})
	// 	if !ok {
	// 		fmt.Println("Invalid medicine details format: Record :%v", i+1)
	// 		return fmt.Errorf("Invalid medicine details format: Record :%v", i+1)
	// 	}
	// 	fmt.Println("medMap_%v: %v", i, medMap)
	// 	medicine := Medicine{
	// 		Name:      fmt.Sprintf("%v", medMap["name"]),
	// 		Dosage:    fmt.Sprintf("%v", medMap["dosage"]),
	// 		Frequency: fmt.Sprintf("%v", medMap["frequency"]),
	// 		Remarks:   fmt.Sprintf("%v", medMap["remarks"]),
	// 	}

	// 	medicines = append(medicines, medicine)
	// }
	// fmt.Println("medicines :", medicines)

	//Validate doctor attributes
	doctorIdentity, err := getUserIdentityName(ctx)
	fmt.Println("doctorIdentity :", doctorIdentity)
	if err != nil {
		return err
	}

	drAttributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization", "orgRole"})
	if err != nil {
		return err
	}
	fmt.Println("userRole :", drAttributes["userRole"])
	fmt.Println("organization :", drAttributes["organization"])
	fmt.Println("orgRole :", drAttributes["orgRole"])

	//fetching patient details
	patientDetailer, err := getEntityDetails(ctx, prescriptionInput.PatientId, PATIENT)
	if err != nil {
		return err
	}
	if patientDetailer == nil {
		return fmt.Errorf("%v patient does not registered with hospital.", prescriptionInput.PatientId)
	}

	var patientDetail User
	// patientDetail, ok := patientDetailer.(User)
	// if !ok {
	// 	return fmt.Errorf("Failed to convert Detailer to User type")
	// }
	// fmt.Println("Patient Details :", patientDetail)

	err = json.Unmarshal(patientDetailer.([]byte), &patientDetail)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal patient details: %s", err.Error())
	}

	fmt.Println("Patient Details :", patientDetail)
	fmt.Println("patientDetail.HospitalId :", patientDetail.HospitalId)
	fmt.Println("patientDetail.DocType :", patientDetail.DocType)

	//Validate attributes
	if drAttributes["organization"] != patientDetail.HospitalId || drAttributes["userRole"] != DOCTOR || patientDetail.DocType != PATIENT {
		return fmt.Errorf("You are not authorized to create prescription")
	}

	var timestamp time.Time
	var timestampString string
	timestamp = time.Now()
	timestampString = timestamp.Format("January 2, 2006 15:04:05")

	prescriptionInput.Timestamp = timestampString
	prescriptionInput.DocType = PRESCRIPTION
	prescriptionInput.DoctorId = doctorIdentity

	// doctorPres := Prescription{
	// 	//	Date             string    `json:"date"`
	// 	Timestamp:        timestampString,
	// 	Doctype:          PRESCRIPTION,
	// 	DoctorId:         "XXX",
	// 	PatientId:        prescriptionInput.PatientId,
	// 	PrescriptionDesc: prescriptionInput.PrescriptionDesc,
	// 	MedicineRecord:   prescriptionInput.MedicineDetails,
	// }

	fmt.Println("doctorPres :", prescriptionInput)

	prescriptionBytes, err := json.Marshal(prescriptionInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Patient Prescription records : %v", err.Error())

	}
	fmt.Println("prescriptionBytes :", string(prescriptionBytes))

	//Inserting prescription record
	// compositeKey, err := ctx.GetStub().CreateCompositeKey(id_DrId_DoctypeIndex, []string{prescriptionInput.PatientId, doctorIdentity,timestampString, PRESCRIPTION})
	// if err != nil {
	// 	return fmt.Errorf("failed to create composite key for prescription for patient %v and err is :%v", prescriptionInput.PatientId, err.Error())
	// }
	// fmt.Println("compositeKey :", compositeKey)

	// err = ctx.GetStub().PutState(compositeKey, prescriptionBytes)
	txID := ctx.GetStub().GetTxID()
	err = ctx.GetStub().PutState(txID, prescriptionBytes)
	if err != nil {
		return fmt.Errorf("failed to insert prescription details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")

	return nil

}

func (s *SmartContract) GetPatientPrescriptionHistory(ctx contractapi.TransactionContextInterface, queryStringInput string) ([]*Prescription, error) {
	//queryString := fmt.Sprintf(queryStringInput)
	fmt.Println("queryStringInput : ", queryStringInput)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryStringInput)
	var prescriptions []*Prescription

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var prescription Prescription
		err = json.Unmarshal(queryResult.Value, &prescription)
		if err != nil {
			return nil, err
		}
		prescriptions = append(prescriptions, &prescription)
	}
	return prescriptions, nil
}

func getEntityDetails(ctx contractapi.TransactionContextInterface, entityId string, docType string) (interface{}, error) {
	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{entityId, docType})
	if err != nil {
		return nil, fmt.Errorf("failed to create composite key for hospital %v and err is :%v", entityId, err.Error())
	}
	objectBytes, err := ctx.GetStub().GetState(compositeKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to read data from world state %s", err.Error())
	}
	if objectBytes == nil {
		return nil, nil
	}
	// var hospitalDetail Hospital
	// err = json.Unmarshal(objectBytes, &hospitalDetail)
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to unmarshal hospital data: %s", err.Error())
	// }
	return objectBytes, nil
}

func getUserIdentityName(ctx contractapi.TransactionContextInterface) (string, error) {
	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}

	fmt.Println("User Identity: ", string(decodeID))

	completeId := string(decodeID)
	userId := completeId[(strings.Index(completeId, "x509::CN=") + 9):strings.Index(completeId, ",")]
	fmt.Println("userId:----------", userId)

	return userId, nil
}

func getCertificateAttributeValue(ctx contractapi.TransactionContextInterface, attrName string) (string, bool, error) {
	attrValue, found, err := ctx.GetClientIdentity().GetAttributeValue(attrName)
	if err != nil {
		return "", false, fmt.Errorf("Failed to read attrValue: %v", err)
	}
	fmt.Println("Attrvalue : ", attrValue)

	return attrValue, found, nil
}

func getAllCertificateAttributes(ctx contractapi.TransactionContextInterface, attrNames []string) (map[string]string, error) {
	attributes := make(map[string]string)
	for _, attrName := range attrNames {
		attrValue, found, err := ctx.GetClientIdentity().GetAttributeValue(attrName)
		if err != nil {
			return nil, fmt.Errorf("Failed to read attrValue for %s: %v", attrName, err)
		}

		if found {
			attributes[attrName] = attrValue
		}
	}

	return attributes, nil
}

func loggedInUserMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("Failed to read mspID: %v", err)
	}
	return mspID, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
