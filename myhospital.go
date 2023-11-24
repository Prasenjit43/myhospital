/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const HOSPITAL = "HOSPITAL"
const HOSPITAL_ADMIN = "HOSPITAL_ADMIN"
const ADMIN = "admin"
const DOCTER = "DOCTER"
const Org1MSP = "Org1MSP"
const Org2MSP = "Org2MSP"

// SmartContract provides functions for managing a Asset and Token
type SmartContract struct {
	contractapi.Contract
}

type Hospital struct {
	Name            string `json:"name"`
	DocType         string `json:"docType"`
	Address         string `json:"address"`
	City            string `json:"city"`
	Pincode         string `json:"pincode"`
	RegistrationNum string `json:"registrationNum"`
}

type HospitalAdmin struct {
	Name             string `json:"name"`
	FirstName        string `json:"firstName"`
	MiddleName       string `json:"middleName"`
	LastName         string `json:"lastName"`
	Dob              string `json:"dob"`
	ContactNo        string `json:"contactNo"`
	EmergencyNo      string `json:"emergencyNo"`
	PermanentAddress string `json:"permanentAddress"`
	EmailId          string `json:"mailId"`
	BloodGroup       int    `json:"bloodGrp"`
	DocType          string `json:"docType"`
	HospitalName     string `json:"hospitalName"`
	Active           bool   `json:"active"`
}

type Doctor struct {
	Name            string `json:"name"`
	DocType         string `json:"docType"`
	HospitalName    string `json:"hospitalName"`
	Active          bool   `json:"active"`
	Specialization  string `json:"specialization"`
	RegistrationNum string `json:"registrationNum"`
}

type Patient struct {
	FirstName        string `json:"firstName"`
	MiddleName       string `json:"middleName"`
	LastName         string `json:"lastName"`
	ContactNo        string `json:"contactNo"`
	EmergencyNo      string `json:"emergencyNo"`
	LocalAddress     string `json:"localAddress"`
	PermanentAddress string `json:"permanentAddress"`
	EmailId          string `json:"mailId"`
	Dob              string `json:"dob"`
	Age              int    `json:"age"`
	BloodGroup       int    `json:"bloodGrp"`
	Allergies        int    `json:"allergies"`
	Comorbidity      int    `json:"comorbidity"`
	Symptoms         int    `json:"symptoms"`
	CurrentTreatment int    `json:"currentTreatment"`
	HospitalName     int    `json:"hospitalName"`
	PersonId         int    `json:"personId"`
	DocType          string `json:"docType"`
	Identity         string `json:"identity"`
	Active           bool   `json:"active"`
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

	//Validate input Parameters
	// if len(strings.TrimSpace(hospitalInput.Name)) == 0 {
	// 	return fmt.Errorf("Hospital name should not be empty")
	// }
	// if hospitalInput.DocType != HOSPITAL {
	// 	return fmt.Errorf(`Doc Type for hospital should be "HOSPITAL"`)
	// }
	// if len(strings.TrimSpace(hospitalInput.Address)) == 0 {
	// 	return fmt.Errorf("Hospital address should not be empty")
	// }
	// if len(strings.TrimSpace(hospitalInput.City)) == 0 {
	// 	return fmt.Errorf("Hospital city should not be empty")
	// }
	// if len(strings.TrimSpace(hospitalInput.Pincode)) == 0 {
	// 	return fmt.Errorf("Hospital pincode should not be empty")
	// }
	// if len(strings.TrimSpace(hospitalInput.RegistrationNum)) == 0 {
	// 	return fmt.Errorf("Hospital registration num should not be empty")
	// }

	//Check if hospital is present or not
	hospitalDetailer, err := getEntityDetails(ctx, hospitalInput.Name)
	if err != nil {
		return err
	}
	if hospitalDetailer != nil {
		return fmt.Errorf("Hospital already exist with name : %v", hospitalInput.Name)
	}

	//Check if hospital ID is present or not
	// hospitalDetail, err := getHospitalDetails(ctx, hospitalInput.Id)
	// if err != nil {
	// 	return err
	// }
	// if hospitalDetail != nil {
	// 	return fmt.Errorf("Hospital already exist with ID : %v", hospitalInput.Id)
	// }

	hospitalBytes, err := json.Marshal(hospitalInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Hospital records : %v", err.Error())

	}

	//Inserting hospital record
	err = ctx.GetStub().PutState(hospitalInput.Name, hospitalBytes)
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

	superAdminIdentity, err := getUserIdentityName(ctx)
	fmt.Println("superAdminIdentity :", superAdminIdentity)

	if superAdminIdentity != ADMIN {
		return fmt.Errorf("permission denied: only admin can call this function")
	}

	//Validate input Parameters
	// if len(strings.TrimSpace(hospitalAdminInput.Name)) == 0 {
	// 	return fmt.Errorf("Admin Name should not be empty")
	// }
	// if hospitalAdminInput.DocType != HOSPITAL_ADMIN {
	// 	return fmt.Errorf(`Doc Type for Asset should be "HOSPITAL_ADMIN"`)
	// }
	// if len(strings.TrimSpace(hospitalAdminInput.HospitalName)) == 0 {
	// 	return fmt.Errorf("Hospital name should not be empty")
	// }

	//Check if hospital is present or not
	hospitalDetailer, err := getEntityDetails(ctx, hospitalAdminInput.HospitalName)
	if err != nil {
		return err
	}
	if hospitalDetailer == nil {
		return fmt.Errorf("Hospital does not exist with name : %v", hospitalAdminInput.HospitalName)
	}

	//Check if hospital admin is present or not
	hospitalAdminDetailer, err := getEntityDetails(ctx, hospitalAdminInput.Name)
	if err != nil {
		return err
	}
	if hospitalAdminDetailer != nil {
		return fmt.Errorf("Hospital admin already exist with name : %v", hospitalAdminInput.Name)
	}

	// hospitalDetail, err := getHospitaAdminDetails(ctx, hospitalAdminInput.Name)
	// if err != nil {
	// 	return err
	// }
	// if hospitalDetail != nil {
	// 	return fmt.Errorf("Hospital admin already exist with name : %v", hospitalAdminInput.Name)
	// }

	hospitalAdminBytes, err := json.Marshal(hospitalAdminInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Hospital admin records : %v", err.Error())
	}

	//Inserting hospital admin record
	err = ctx.GetStub().PutState(hospitalAdminInput.Name, hospitalAdminBytes)
	if err != nil {
		return fmt.Errorf("failed to insert hospital admin details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")
	return nil
}

func (s *SmartContract) AddDocter(ctx contractapi.TransactionContextInterface, docterInputString string) error {
	var docterInput Doctor
	err := json.Unmarshal([]byte(docterInputString), &docterInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", docterInput)

	hospitalAdminRole, _, err := getCertificateAttributeValue(ctx, "userRole")
	fmt.Println("Attribute userRole value :", hospitalAdminRole)
	if hospitalAdminRole != HOSPITAL_ADMIN {
		return fmt.Errorf("Only Hospital Admin are allowed to register docter")
	}

	adminHospitalName, _, err := getCertificateAttributeValue(ctx, "organizationName")
	fmt.Println("Attribute organizationName value :", adminHospitalName)
	if adminHospitalName != docterInput.HospitalName {
		return fmt.Errorf("Mismatch hospital name")
	}

	hospitalAdminMSPOrgId, err := getMSPID(ctx)
	if err != nil {
		return err
	}
	fmt.Println("hospitalAdminMSPOrgId :", hospitalAdminMSPOrgId)
	if hospitalAdminMSPOrgId != Org1MSP {
		return fmt.Errorf("OrgMSP Id is different")
	}

	//Validate input Parameters
	if len(strings.TrimSpace(docterInput.Name)) == 0 {
		return fmt.Errorf("Docter Name should not be empty")
	}
	if docterInput.DocType != DOCTER {
		return fmt.Errorf(`Doc Type for Asset should be "DOCTER"`)
	}
	if len(strings.TrimSpace(docterInput.HospitalName)) == 0 {
		return fmt.Errorf("Hospital name should not be empty")
	}
	if len(strings.TrimSpace(docterInput.Specialization)) == 0 {
		return fmt.Errorf("Specialization should not be empty")
	}

	if len(strings.TrimSpace(docterInput.RegistrationNum)) == 0 {
		return fmt.Errorf("RegistrationNum num should not be empty")
	}

	//Check if hospital admin is present or not
	docterDetailer, err := getHospitalDetails(ctx, docterInput.Name)
	if err != nil {
		return err
	}
	if docterDetailer != nil {
		return fmt.Errorf("Docter %v is already registered with hospital.", docterInput.Name)
	}

	// docterDetail, err := getDocterDetails(ctx, docterInput.Name)
	// if err != nil {
	// 	return err
	// }
	// if docterDetail != nil {
	// 	return fmt.Errorf("Docter %v is already registered with hospital.", docterInput.Name)
	// }

	docterBytes, err := json.Marshal(docterInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Docter records : %v", err.Error())
	}

	//Inserting hospital admin record
	err = ctx.GetStub().PutState(docterInput.Name, docterBytes)
	if err != nil {
		return fmt.Errorf("failed to insert docter details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")
	return nil
}

func (s *SmartContract) AddPatient(ctx contractapi.TransactionContextInterface, patientInputString string) error {
	var docterInput Doctor
	err := json.Unmarshal([]byte(patientInputString), &docterInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", docterInput)

	hospitalAdminRole, _, err := getCertificateAttributeValue(ctx, "userRole")
	fmt.Println("Attribute userRole value :", hospitalAdminRole)
	if hospitalAdminRole != HOSPITAL_ADMIN {
		return fmt.Errorf("Only Hospital Admin are allowed to register docter")
	}

	adminHospitalName, _, err := getCertificateAttributeValue(ctx, "organizationName")
	fmt.Println("Attribute organizationName value :", adminHospitalName)
	if adminHospitalName != docterInput.HospitalName {
		return fmt.Errorf("Mismatch hospital name")
	}

	hospitalAdminMSPOrgId, err := getMSPID(ctx)
	if err != nil {
		return err
	}
	fmt.Println("hospitalAdminMSPOrgId :", hospitalAdminMSPOrgId)
	if hospitalAdminMSPOrgId != Org1MSP {
		return fmt.Errorf("OrgMSP Id is different")
	}

	//Validate input Parameters
	if len(strings.TrimSpace(docterInput.Name)) == 0 {
		return fmt.Errorf("Docter Name should not be empty")
	}
	if docterInput.DocType != DOCTER {
		return fmt.Errorf(`Doc Type for Asset should be "DOCTER"`)
	}
	if len(strings.TrimSpace(docterInput.HospitalName)) == 0 {
		return fmt.Errorf("Hospital name should not be empty")
	}
	if len(strings.TrimSpace(docterInput.Specialization)) == 0 {
		return fmt.Errorf("Specialization should not be empty")
	}

	if len(strings.TrimSpace(docterInput.RegistrationNum)) == 0 {
		return fmt.Errorf("RegistrationNum num should not be empty")
	}

	//Check if hospital admin is present or not
	docterDetailer, err := getHospitalDetails(ctx, docterInput.Name)
	if err != nil {
		return err
	}
	if docterDetailer != nil {
		return fmt.Errorf("Docter %v is already registered with hospital.", docterInput.Name)
	}

	// docterDetail, err := getDocterDetails(ctx, docterInput.Name)
	// if err != nil {
	// 	return err
	// }
	// if docterDetail != nil {
	// 	return fmt.Errorf("Docter %v is already registered with hospital.", docterInput.Name)
	// }

	docterBytes, err := json.Marshal(docterInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Docter records : %v", err.Error())
	}

	//Inserting hospital admin record
	err = ctx.GetStub().PutState(docterInput.Name, docterBytes)
	if err != nil {
		return fmt.Errorf("failed to insert docter details to couchDB : %v", err.Error())
	}
	fmt.Println("****************************")
	return nil
}

func getHospitalDetails(ctx contractapi.TransactionContextInterface, hospitalId string) (*Hospital, error) {
	hospitalBytes, err := ctx.GetStub().GetState(hospitalId)
	if err != nil {
		return nil, fmt.Errorf("Failed to read data from world state %s", err.Error())
	}
	if hospitalBytes == nil {
		return nil, nil
	}
	var hospitalDetail Hospital
	err = json.Unmarshal(hospitalBytes, &hospitalDetail)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal hospital data: %s", err.Error())
	}
	return &hospitalDetail, nil
}

func getEntityDetails(ctx contractapi.TransactionContextInterface, entityId string) (interface{}, error) {
	objectBytes, err := ctx.GetStub().GetState(entityId)
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

func getHospitaAdminDetails(ctx contractapi.TransactionContextInterface, hospitalAdmin string) (*HospitalAdmin, error) {
	hospitalAdminBytes, err := ctx.GetStub().GetState(hospitalAdmin)
	if err != nil {
		return nil, fmt.Errorf("Failed to read data from world state %s", err.Error())
	}
	if hospitalAdminBytes == nil {
		return nil, nil
	}
	var hospitalAdminDetail HospitalAdmin
	err = json.Unmarshal(hospitalAdminBytes, &hospitalAdminDetail)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal hospital data: %s", err.Error())
	}
	return &hospitalAdminDetail, nil
}
func getDocterDetails(ctx contractapi.TransactionContextInterface, docterName string) (*Doctor, error) {
	docterBytes, err := ctx.GetStub().GetState(docterName)
	if err != nil {
		return nil, fmt.Errorf("Failed to read data from world state %s", err.Error())
	}
	if docterBytes == nil {
		return nil, nil
	}
	var doctorDetail Doctor
	err = json.Unmarshal(docterBytes, &doctorDetail)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal docter details: %s", err.Error())
	}
	return &doctorDetail, nil
}

func getUserIdentityName(ctx contractapi.TransactionContextInterface) (string, error) {
	fmt.Printf("getUserId start-->")
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

func getMSPID(ctx contractapi.TransactionContextInterface) (string, error) {
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
