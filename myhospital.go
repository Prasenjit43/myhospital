/*
SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const HOSPITAL = "HOSPITAL"
const PATIENT = "PATIENT"
const HOSPITAL_ADMIN = "HOSPITAL_ADMIN"
const DRUGGIST = "DRUGGIST"
const PATHOLOGIST = "PATHOLOGIST"
const ACCESS = "ACCESS"

const ADMIN = "admin"
const DOCTOR = "DOCTOR"
const Org1MSP = "Org1MSP"
const Org2MSP = "Org2MSP"
const PRESCRIPTION = "PRESCRIPTION"

const idDoctypeIndex = "id~doctype"

//const id_DrId_DoctypeIndex = "id~drId~doctype"

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
	Active          bool   `json:"active"`
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
	EmailId          string `json:"emailId,omitempty"`
	BloodGroup       string `json:"bloodGroup,omitempty"`
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
	EmailId          string `json:"emailId,omitempty"`
	Dob              string `json:"dob,omitempty"`
	BloodGroup       string `json:"bloodGroup,omitempty"`
	HospitalId       string `json:"hospitalId,omitempty"`
	DocType          string `json:"docType"`
	RegistrationNum  string `json:"registrationNum,omitempty"`
	Active           bool   `json:"active,omitempty"`
	// PrescriptionDeatils interface{} `json:"annexture,omitempty"`
	MetaData interface{} `json:"annexture,omitempty"`

	// Allergies        string `json:"allergies,omitempty"`
	// Comorbidity      string `json:"comorbidity,omitempty"`
	// Symptoms         string `json:"symptoms,omitempty"`
	// CurrentTreatment string `json:"currentTreatment,omitempty"`
}

type Medicine struct {
	Name      string `json:"name,omitempty"`
	Dosage    string `json:"dosage,omitempty"`
	Frequency string `json:"frequency,omitempty"`
	Remarks   string `json:"remarks,omitempty"`
	Amount    int32  `json:"amount,omitempty"`
}

type Prescription struct {
	Timestamp      int64      `json:"timestamp,omitempty"`
	DocType        string     `json:"docType,omitempty"`
	DoctorId       string     `json:"doctorId,omitempty"`
	Id             string     `json:"id"`
	Desc           string     `json:"desc"`
	MedicineRecord []Medicine `json:"medicine"`
}

type Billing struct {
	//Consultation   int32      `json:"consultation,omitempty"`
	//Timestamp      int64      `json:"timestamp,omitempty"`
	DocType string `json:"docType,omitempty"`
	//DoctorId       string     `json:"doctorId,omitempty"`
	Id string `json:"id,omitempty"`
	//Desc           string     `json:"desc,omitempty"`
	MedicineRecord []Medicine `json:"medicine,omitempty"`
}

type MedicalRecordAccess struct {
	Id       string `json:"id"`
	DoctorId string `json:"doctorId,omitempty"`
	DocType  string `json:"docType,omitempty"`
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
		return fmt.Errorf("permission denied: only super admin can call this function")
	}

	//Check if hospital is present or not
	hospitalDetailer, err := getEntityDetails(ctx, hospitalInput.Id, hospitalInput.DocType)
	if err != nil {
		return err
	}
	if hospitalDetailer != nil {
		return fmt.Errorf("Hospital already exist with id : %v", hospitalInput.Id)
	}

	//Inserting hospital record
	hospitalBytes, err := json.Marshal(hospitalInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Hospital record : %v", err.Error())

	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{hospitalInput.Id, hospitalInput.DocType})
	if err != nil {
		return fmt.Errorf("failed to create composite key for hospital %v and err is :%v", hospitalInput.Id, err.Error())
	}
	err = ctx.GetStub().PutState(compositeKey, hospitalBytes)
	if err != nil {
		return fmt.Errorf("failed to insert hospital details to couchDB : %v", err.Error())
	}
	fmt.Println("********** End of AddHospital Function******************")
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
		return fmt.Errorf("permission denied: only super admin can call this function")
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

	//Inserting hospital admin record
	hospitalAdminBytes, err := json.Marshal(hospitalAdminInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of Hospital admin record : %v", err.Error())
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{hospitalAdminInput.Id, hospitalAdminInput.DocType})
	if err != nil {
		return fmt.Errorf("failed to create composite key for hospital admin %v and err is :%v", hospitalAdminInput.Id, err.Error())
	}
	err = ctx.GetStub().PutState(compositeKey, hospitalAdminBytes)
	if err != nil {
		return fmt.Errorf("failed to insert hospital admin details to couchDB : %v", err.Error())
	}
	fmt.Println("********** End of AddHospitalAdmin Function******************")
	return nil
}

func (s *SmartContract) AddEntity(ctx contractapi.TransactionContextInterface, entityInputString string) error {
	var entityInput User
	err := json.Unmarshal([]byte(entityInputString), &entityInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", entityInput)

	//fetching cerificate attributes
	attributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization"})
	if err != nil {
		return err
	}
	fmt.Println("userRole :", attributes["userRole"])
	fmt.Println("organization :", attributes["organization"])

	if attributes["userRole"] != HOSPITAL_ADMIN {
		return fmt.Errorf("Only Hospital Admin are allowed to register %v", entityInput.DocType)
	}

	// hospitalAdminMSPOrgId, err := loggedInUserMSPID(ctx)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println("hospitalAdminMSPOrgId :", hospitalAdminMSPOrgId)
	// if hospitalAdminMSPOrgId != Org1MSP {
	// 	return fmt.Errorf("OrgMSP Id is different")
	// }

	//Check if entity-id is present or not
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
		return fmt.Errorf("failed to create composite key for %v and err is :%v", entityInput.Id, err.Error())
	}
	err = ctx.GetStub().PutState(compositeKey, entityBytes)
	if err != nil {
		return fmt.Errorf("failed to insert entity details to couchDB : %v", err.Error())
	}
	fmt.Println("********** End of AddEntity Function******************")
	return nil
}

func (s *SmartContract) SelfRegistration(ctx contractapi.TransactionContextInterface, entityInputString string) error {
	var entityInput User
	err := json.Unmarshal([]byte(entityInputString), &entityInput)
	if err != nil {
		return fmt.Errorf("Error while doing unmarshal of input string : %v", err.Error())
	}
	fmt.Println("Input String :", entityInput)

	//Check if hospital is present or not
	hospitalDetailer, err := getEntityDetails(ctx, entityInput.HospitalId, HOSPITAL)
	if err != nil {
		return err
	}
	if hospitalDetailer == nil {
		return fmt.Errorf("Hospital does not exist with Id : %v", entityInput.HospitalId)
	}

	//Check if patient-id is present or not
	entityDetailer, err := getEntityDetails(ctx, entityInput.Id, entityInput.DocType)
	if err != nil {
		return err
	}
	if entityDetailer != nil {
		return fmt.Errorf("%v is already registered with hospital.", entityInput.Id)
	}

	//Inserting entity record
	entityBytes, err := json.Marshal(entityInput)
	if err != nil {
		return fmt.Errorf("failed to marshal of entity records : %v", err.Error())
	}

	compositeKey, err := ctx.GetStub().CreateCompositeKey(idDoctypeIndex, []string{entityInput.Id, entityInput.DocType})
	if err != nil {
		return fmt.Errorf("failed to create composite key for %v and err is :%v", entityInput.Id, err.Error())
	}
	err = ctx.GetStub().PutState(compositeKey, entityBytes)
	if err != nil {
		return fmt.Errorf("failed to insert entity details to couchDB : %v", err.Error())
	}
	fmt.Println("********** End of Self Registration Function ******************")
	return nil
}

// Dr may be act as patient in hospital, pass doctype also
func (s *SmartContract) ViewOwnDetails(ctx contractapi.TransactionContextInterface) (string, error) {

	userIdentity, err := getUserIdentityName(ctx)
	fmt.Println("userIdentity :", userIdentity)
	if err != nil {
		return "", err
	}

	docType, _, err := getCertificateAttributeValue(ctx, "userRole")
	fmt.Println("Attribute docType value :", docType)

	//Check if entity Id is present or not
	objectBytes, err := getEntityDetails(ctx, userIdentity, docType)
	if err != nil {
		return "", err
	}
	if objectBytes == nil {
		return "", fmt.Errorf("Record for %v user does not exist", userIdentity)
	}

	fmt.Println("********** End of ViewOwnDetails Function******************")
	return string(objectBytes.([]byte)), nil
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

	userBytes, err := getEntityDetails(ctx, input.UserId, PATIENT)
	if err != nil {
		return "", err
	}
	if userBytes == nil {
		return "", fmt.Errorf("Record for %v patient does not exist", input.UserId)
	}

	var patientDetail User
	err = json.Unmarshal(userBytes.([]byte), &patientDetail)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal user data: %s", err.Error())
	}
	fmt.Println("User Data : ", patientDetail)

	//Validate attributes
	if drAttributes["organization"] != patientDetail.HospitalId || (drAttributes["userRole"] != DOCTOR && drAttributes["userRole"] != HOSPITAL_ADMIN) || patientDetail.DocType != PATIENT {
		return "", fmt.Errorf("Ony doctor are authorized to see patient data")
	}

	/**/
	queryString := fmt.Sprintf(`{"selector":{"id":"%s","docType":"PRESCRIPTION"},"sort":[{"timestamp":"desc"}],"limit":3}`, patientDetail.Id)
	last3PrescriptionDetails, err := s.GetPatientPrescriptionHistory(ctx, queryString)
	fmt.Println("last3PrescriptionDetails : ", last3PrescriptionDetails)

	/*****************************************/

	var buffer bytes.Buffer
	//buffer.WriteString("{")
	buffer.WriteString("{\"patientDetail\":")
	//buffer.WriteString("\"")
	buffer.WriteString(string(userBytes.([]byte)))
	//buffer.WriteString("\"")

	buffer.WriteString(", \"last3PrescriptionDetails\":")
	buffer.WriteString(last3PrescriptionDetails)
	buffer.WriteString("}")
	fmt.Println("buffer string : ", buffer.String())
	fmt.Println("********** End of ViewPatientDetails Function ******************")
	return buffer.String(), nil

}

func (s *SmartContract) ViewPatientDetailsByOtherEntity(ctx contractapi.TransactionContextInterface, inputString string) (string, error) {
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

	//Validate patient id
	userBytes, err := getEntityDetails(ctx, input.UserId, PATIENT)
	if err != nil {
		return "", err
	}
	if userBytes == nil {
		return "", fmt.Errorf("Record for %v patient does not exist", input.UserId)
	}

	var patientDetail User
	err = json.Unmarshal(userBytes.([]byte), &patientDetail)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal user data: %s", err.Error())
	}
	fmt.Println("User Data : ", patientDetail)

	isAccessFlag := false
	/*Check if already granted access or not*/
	queryString := fmt.Sprintf(`{"selector":{"id":"%s","docType":"ACCESS","$or":[{"doctorId":"%s"},{"doctorId":"%s"}]}}`, input.UserId, drAttributes["organization"], input.ViewerId)
	fmt.Println("queryString : ", queryString)
	key, err := isAccessGranted(ctx, queryString)
	if err != nil {
		return "", err
	}
	fmt.Println("key 		:", key)
	if key != "nil" {
		isAccessFlag = true
	}

	if !isAccessFlag && drAttributes["organization"] == patientDetail.HospitalId && (drAttributes["userRole"] == DOCTOR || drAttributes["userRole"] == HOSPITAL_ADMIN) {
		isAccessFlag = true
	}

	if !isAccessFlag {
		return "", fmt.Errorf("You are not authorized to see patient data")
	}

	//queryString = fmt.Sprintf(`{"selector":{"id":"%s","docType":"PRESCRIPTION"},"sort":[{"timestamp":"desc"}],"limit":3}`, patientDetail.Id)
	queryString = fmt.Sprintf(`{"selector":{"id":"%s","docType":"PRESCRIPTION"},"sort":[{"timestamp":"desc"}]}`, patientDetail.Id)
	fmt.Println("last3PrescriptionDetails queryString : ", queryString)
	//last3PrescriptionDetails, err := s.GetPatientPrescriptionHistory(ctx, queryString)
	// fmt.Println("last3PrescriptionDetails : ", last3PrescriptionDetails)
	last3PrescriptionDetails, err := getQueryResultForQueryStringWithPagination(ctx, queryString, 3, "")
	if err != nil {
		return "", err
	}

	fmt.Println("last3PrescriptionDetails err : ", last3PrescriptionDetails)

	/*****************************************/

	var buffer bytes.Buffer
	//buffer.WriteString("{")
	buffer.WriteString("{\"patientDetail\":")
	//buffer.WriteString("\"")
	buffer.WriteString(string(userBytes.([]byte)))
	//buffer.WriteString("\"")

	buffer.WriteString(", \"last3PrescriptionDetails\":")
	buffer.WriteString(last3PrescriptionDetails)
	buffer.WriteString("}")
	fmt.Println("buffer string : ", buffer.String())
	fmt.Println("********** End of ViewPatientDetails Function ******************")
	return buffer.String(), nil

}

func (s *SmartContract) CreatePrescription(ctx contractapi.TransactionContextInterface, prescriptionInputString string) (string, error) {
	var prescriptionInput Prescription
	err := json.Unmarshal([]byte(prescriptionInputString), &prescriptionInput)
	if err != nil {
		return "", fmt.Errorf("Error while doing unmarshal of prescription input string : %v", err.Error())
	}
	fmt.Println("Input String :", prescriptionInput)

	//Validate doctor attributes
	doctorIdentity, err := getUserIdentityName(ctx)
	fmt.Println("doctorIdentity :", doctorIdentity)
	if err != nil {
		return "", err
	}

	drAttributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization", "orgRole"})
	if err != nil {
		return "", err
	}
	fmt.Println("userRole 		:", drAttributes["userRole"])
	fmt.Println("organization 	:", drAttributes["organization"])
	fmt.Println("orgRole 		:", drAttributes["orgRole"])

	//fetching patient details
	patientDetailer, err := getEntityDetails(ctx, prescriptionInput.Id, PATIENT)
	if err != nil {
		return "", err
	}
	if patientDetailer == nil {
		return "", fmt.Errorf("%v patient does not registered with hospital.", prescriptionInput.Id)
	}

	var patientDetail User
	err = json.Unmarshal(patientDetailer.([]byte), &patientDetail)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal patient details: %s", err.Error())
	}
	fmt.Println("Patient Details 			:", patientDetail)

	//Validate attributes
	if drAttributes["organization"] != patientDetail.HospitalId || drAttributes["userRole"] != DOCTOR || patientDetail.DocType != PATIENT {
		return "", fmt.Errorf("You are not authorized to create prescription")
	}

	//Generating timestamp and formatting
	// var timestamp time.Time
	// var timestampString string
	// timestamp = time.Now()
	// timestampString = timestamp.Format("January 2, 2006 15:04:05")

	// //assigning prescription attributes
	// prescriptionInput.Timestamp = timestampString
	prescriptionInput.DocType = PRESCRIPTION
	prescriptionInput.DoctorId = doctorIdentity

	fmt.Println("Final Prescription Details :", prescriptionInput)

	prescriptionBytes, err := json.Marshal(prescriptionInput)
	if err != nil {
		return "", fmt.Errorf("failed to marshal of Patient Prescription records : %v", err.Error())
	}
	fmt.Println("prescriptionBytes :", string(prescriptionBytes))

	//Inserting prescription record
	// compositeKey, err := ctx.GetStub().CreateCompositeKey(id_DrId_DoctypeIndex, []string{prescriptionInput.PatientId, doctorIdentity,timestampString, PRESCRIPTION})
	// if err != nil {
	// 	return fmt.Errorf("failed to create composite key for prescription for patient %v and err is :%v", prescriptionInput.PatientId, err.Error())
	// }
	// err = ctx.GetStub().PutState(compositeKey, prescriptionBytes)
	txID := ctx.GetStub().GetTxID()
	err = ctx.GetStub().PutState(txID, prescriptionBytes)
	if err != nil {
		return "", fmt.Errorf("failed to insert prescription details to couchDB : %v", err.Error())
	}
	fmt.Println("********** End of CreatePrescription Function******************")
	return txID, nil
}

// func (s *SmartContract) GenerateBill(ctx contractapi.TransactionContextInterface, billingInputString string) error {
// 	billingInput := struct {
// 		Id string `json:"id,omitempty"`

// 	}{}
// 	err := json.Unmarshal([]byte(prescriptionInputString), &prescriptionInput)
// 	if err != nil {
// 		return "", fmt.Errorf("Error while doing unmarshal of prescription input string : %v", err.Error())
// 	}
// 	fmt.Println("Input String :", prescriptionInput)

// 	//fetching cerificate attributes
// 	attributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization"})
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Println("userRole :", attributes["userRole"])
// 	fmt.Println("organization :", attributes["organization"])

// 	if attributes["userRole"] != DRUGGIST {
// 		return fmt.Errorf("Only Druggist Admin are allowed to generate bill ")
// 	}

// 	identity, err := getUserIdentityName(ctx)
// 	fmt.Println("identity :", identity)
// 	if identity != DRUGGIST {
// 		return fmt.Errorf("permission denied: only DRUGGIST can call this function")
// 	}

// 	//Check if druggist is present or not
// 	entityDetailer, err := getEntityDetails(ctx, identity, DRUGGIST)
// 	if err != nil {
// 		return err
// 	}
// 	if entityDetailer == nil {
// 		return fmt.Errorf("%v does not registered with hospital.", identity)
// 	}

// 	//Check if prescription is present or not
// 	entityDetailer, err = getEntityDetails(ctx, prescriptionId, PRESCRIPTION)
// 	if err != nil {
// 		return err
// 	}
// 	if entityDetailer == nil {
// 		return fmt.Errorf("Prescription with %v does not exist", prescriptionId)
// 	}

// 	return nil
// }

func (s *SmartContract) AccessRights(ctx contractapi.TransactionContextInterface, accessInputString string) (string, error) {
	accessInput := struct {
		FromUserId   string `json:"fromUserId"`
		ToUserId     string `json:"toUserId"`
		ToUserIdRole string `json:"toUserIdRole"`
		AccessType   string `json:"accessType"`
	}{}

	err := json.Unmarshal([]byte(accessInputString), &accessInput)
	if err != nil {
		return "", fmt.Errorf("Error while doing unmarshal of access input string : %v", err.Error())
	}
	fmt.Println("Input String :", accessInput)

	//fetching cerificate attributes of patient
	attributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization", "orgRole"})
	if err != nil {
		return "", err
	}
	fmt.Println("userRole 		:", attributes["userRole"])
	fmt.Println("organization 	:", attributes["organization"])
	fmt.Println("orgRole 		:", attributes["orgRole"])

	if attributes["userRole"] != PATIENT {
		return "", fmt.Errorf("Only Patient are allowed to grant/revoke access ")
	}

	if accessInput.ToUserIdRole == PATIENT {
		return "", fmt.Errorf("You cannot grant access to Patient ")
	}

	//Check if entity-id is present or not
	entityDetailer, err := getEntityDetails(ctx, accessInput.ToUserId, accessInput.ToUserIdRole)
	if err != nil {
		return "", err
	}
	if entityDetailer == nil {
		return "", fmt.Errorf("%v does not registered ", accessInput.ToUserId)
	}

	/*Check if already granted access or not*/
	queryString := fmt.Sprintf(`{"selector":{"id":"%v","docType":"ACCESS","doctorId":"%v"}}`, accessInput.FromUserId, accessInput.ToUserId)
	key, err := isAccessGranted(ctx, queryString)
	if err != nil {
		return "", err
	}
	fmt.Println("key 		:", key)

	var returnVar string
	//Inserting new access record
	if key == "nil" {
		if accessInput.AccessType == "GRANT" {
			medicalRecordAccess := MedicalRecordAccess{
				Id:       accessInput.FromUserId,
				DoctorId: accessInput.ToUserId,
				DocType:  ACCESS,
			}
			accessBytes, err := json.Marshal(medicalRecordAccess)
			txID := ctx.GetStub().GetTxID()
			err = ctx.GetStub().PutState(txID, accessBytes)
			if err != nil {
				return "", fmt.Errorf("failed to grant access rights : %v", err.Error())
			}
			returnVar = txID
		} else if accessInput.AccessType == "REVOKE" {
			return "", fmt.Errorf("ACCESS already revoked")
		}
	} else if key != "nil" {
		if accessInput.AccessType == "REVOKE" {
			err = ctx.GetStub().DelState(key)
			if err != nil {
				return "", fmt.Errorf("failed to revoke access rights  : %v", err.Error())
			}
			returnVar = key
		} else if accessInput.AccessType == "GRANT" {
			return "", fmt.Errorf("ACCESS already granted")
		}
	}

	fmt.Println("********** End of AccessRights Function ******************")
	return returnVar, nil
}

// func (s *SmartContract) GetPrescription(ctx contractapi.TransactionContextInterface, queryStringInput string) (string, error) {
// 	fmt.Println("Query String : ", queryStringInput)

// 	//get loggedin userid
// 	identity, err := getUserIdentityName(ctx)
// 	fmt.Println("identity :", identity)

// 	//fetching cerificate attributes of loggedin identity
// 	attributes, err := getAllCertificateAttributes(ctx, []string{"userRole", "organization", "orgRole"})
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println("userRole 		:", attributes["userRole"])
// 	fmt.Println("organization 	:", attributes["organization"])
// 	fmt.Println("orgRole 		:", attributes["orgRole"])

// 	prescriptionDetails , err := getPrescriptionDetails(ctx,queryStringInput)

// 	/*Check if already granted access or not*/
// 	queryString := fmt.Sprintf(`{"selector":{"id":"%s","docType":"ACCESS","$or":[{"doctorId":"%s"},{"doctorId":"%s"}]}}`, identity, attributes, input.ViewerId)
// 	fmt.Println("queryString : ", queryString)
// 	key, err := isAccessGranted(ctx, queryString)
// 	if err != nil {
// 		return "", err
// 	}
// 	fmt.Println("key 		:", key)
// 	if key != "nil" {
// 		isAccessFlag = true
// 	}

// 	resultsIterator, err := ctx.GetStub().GetQueryResult(queryStringInput)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resultsIterator.Close()

// 	if !resultsIterator.HasNext() {
// 		fmt.Println(`{ "message" : "No Prescription records found"}`)
// 		return "", fmt.Errorf("No Prescription records found")
// 	}

// }

// func (s *SmartContract) getPrescriptionDetails(ctx contractapi.TransactionContextInterface, queryStringInput string) (string, error) {
// 	fmt.Println("Query String : ", queryStringInput)
// 	resultsIterator, err := ctx.GetStub().GetQueryResult(queryStringInput)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resultsIterator.Close()

// 	if !resultsIterator.HasNext() {
// 		fmt.Println(`{ "message" : "No Prescription records found"}`)
// 		return "", fmt.Errorf("No Prescription records found")
// 	}

// 	// var buffer bytes.Buffer
// 	// buffer.WriteString("[")
// 	// var prescriptions []*Prescription
// 	// fmt.Println("resultsIterator : ", resultsIterator)
// 	// fmt.Println("resultsIterator.HasNext() - Before : ", resultsIterator.HasNext())

// 	// bArrayMemberAlreadyWritten := false
// 	// for resultsIterator.HasNext() {
// 	// 	queryResult, err := resultsIterator.Next()
// 	// 	if err != nil {
// 	// 		return "", err
// 	// 	}
// 	// 	if bArrayMemberAlreadyWritten == true {
// 	// 		buffer.WriteString(",")
// 	// 	}
// 	// 	var prescription Prescription
// 	// 	err = json.Unmarshal(queryResult.Value, &prescription)
// 	// 	if err != nil {
// 	// 		return "", err
// 	// 	}
// 	// 	fmt.Println("buffer string : ", string(queryResult.Value))
// 	// 	buffer.WriteString(string(queryResult.Value))
// 	// 	prescriptions = append(prescriptions, &prescription)
// 	// 	bArrayMemberAlreadyWritten = true
// 	// }
// 	// fmt.Println("resultsIterator.HasNext() - After : ", resultsIterator.HasNext())
// 	// buffer.WriteString("]")
// 	// fmt.Println("buffer string : ", buffer.String())

// 	/******************************/
// 	prescriptions, err := constructQueryResponseFromIterator(resultsIterator)
// 	if err != nil {
// 		return "", err
// 	}

// 	fmt.Println("********** End of GetPatientPrescriptionHistory Function******************")
// 	//return buffer.String(), nil
// 	return prescriptions, nil
// }

func (s *SmartContract) GetPatientPrescriptionHistory(ctx contractapi.TransactionContextInterface, queryStringInput string) (string, error) {
	fmt.Println("Query String : ", queryStringInput)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryStringInput)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	if !resultsIterator.HasNext() {
		fmt.Println(`{ "message" : "No Prescription records found"}`)
		return "", fmt.Errorf("No Prescription records found")
	}

	// var buffer bytes.Buffer
	// buffer.WriteString("[")
	// var prescriptions []*Prescription
	// fmt.Println("resultsIterator : ", resultsIterator)
	// fmt.Println("resultsIterator.HasNext() - Before : ", resultsIterator.HasNext())

	// bArrayMemberAlreadyWritten := false
	// for resultsIterator.HasNext() {
	// 	queryResult, err := resultsIterator.Next()
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	if bArrayMemberAlreadyWritten == true {
	// 		buffer.WriteString(",")
	// 	}
	// 	var prescription Prescription
	// 	err = json.Unmarshal(queryResult.Value, &prescription)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	fmt.Println("buffer string : ", string(queryResult.Value))
	// 	buffer.WriteString(string(queryResult.Value))
	// 	prescriptions = append(prescriptions, &prescription)
	// 	bArrayMemberAlreadyWritten = true
	// }
	// fmt.Println("resultsIterator.HasNext() - After : ", resultsIterator.HasNext())
	// buffer.WriteString("]")
	// fmt.Println("buffer string : ", buffer.String())

	/******************************/
	prescriptions, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return "", err
	}

	fmt.Println("********** End of GetPatientPrescriptionHistory Function******************")
	//return buffer.String(), nil
	return prescriptions, nil
}

func (s *SmartContract) GetPatientPrescriptionHistory1(ctx contractapi.TransactionContextInterface, queryStringInput string) (string, error) {
	fmt.Println("Query String : ", queryStringInput)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryStringInput)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	if !resultsIterator.HasNext() {
		fmt.Println(`{ "message" : "No Prescription records found"}`)
		return "", fmt.Errorf("No Prescription records found")
	}

	buffer.WriteString("[")
	var prescriptions []*Prescription
	fmt.Println("resultsIterator : ", resultsIterator)
	fmt.Println("resultsIterator.HasNext() - Before : ", resultsIterator.HasNext())

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		var prescription Prescription
		err = json.Unmarshal(queryResult.Value, &prescription)
		if err != nil {
			return "", err
		}
		fmt.Println("buffer string : ", string(queryResult.Value))
		buffer.WriteString(string(queryResult.Value))
		prescriptions = append(prescriptions, &prescription)
		bArrayMemberAlreadyWritten = true
	}
	fmt.Println("resultsIterator.HasNext() - After : ", resultsIterator.HasNext())
	buffer.WriteString("]")
	fmt.Println("buffer string : ", buffer.String())

	fmt.Println("********** End of GetPatientPrescriptionHistory Function******************")
	return buffer.String(), nil
}

func (s *SmartContract) GetPatientPrescriptionHistory_bkp(ctx contractapi.TransactionContextInterface, queryStringInput string) ([]*Prescription, error) {
	fmt.Println("Query String : ", queryStringInput)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryStringInput)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var prescriptions []*Prescription

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

func getQueryResultForQueryStringWithPagination(ctx contractapi.TransactionContextInterface, queryString string, pageSize int32, bookmark string) (string, error) {

	resultsIterator, responseMetadata, err := ctx.GetStub().GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return "", err
	}
	defer resultsIterator.Close()

	prescriptions, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return "", err
	}

	fmt.Println("Records : ", prescriptions)
	fmt.Println("FetchedRecordsCount : ", responseMetadata.FetchedRecordsCount)
	fmt.Println("Bookmark : ", responseMetadata.Bookmark)
	return prescriptions, nil
}

// constructQueryResponseFromIterator constructs a slice of assets from the resultsIterator
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (string, error) {
	var prescriptions []*Prescription
	var buffer bytes.Buffer
	buffer.WriteString("[")
	fmt.Println("resultsIterator : ", resultsIterator)
	fmt.Println("resultsIterator.HasNext() - Before : ", resultsIterator.HasNext())
	bArrayMemberAlreadyWritten := false

	// for resultsIterator.HasNext() {
	// 	queryResult, err := resultsIterator.Next()
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	var prescription Prescription
	// 	err = json.Unmarshal(queryResult.Value, &prescription)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	prescriptions = append(prescriptions, &prescription)
	// }

	// var buffer bytes.Buffer
	// buffer.WriteString("[")
	// var prescriptions []*Prescription
	// fmt.Println("resultsIterator : ", resultsIterator)
	// fmt.Println("resultsIterator.HasNext() - Before : ", resultsIterator.HasNext())

	// bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return "", err
		}

		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		var prescription Prescription
		err = json.Unmarshal(queryResult.Value, &prescription)
		if err != nil {
			return "", err
		}
		fmt.Println("buffer string : ", string(queryResult.Value))
		buffer.WriteString(string(queryResult.Value))
		prescriptions = append(prescriptions, &prescription)
		bArrayMemberAlreadyWritten = true
	}
	fmt.Println("resultsIterator.HasNext() - After : ", resultsIterator.HasNext())
	buffer.WriteString("]")
	fmt.Println("buffer string : ", buffer.String())
	return buffer.String(), nil

}

func constructQueryResponseFromIterator_bkp(resultsIterator shim.StateQueryIteratorInterface) ([]*Prescription, error) {
	var prescriptions []*Prescription
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

	// Decoding the objectBytes assuming it's JSON
	var data map[string]interface{}
	_ = json.Unmarshal(objectBytes, &data)

	// Access the "active" attribute
	activeAttribute, ok := data["active"]
	if !ok {
		return nil, fmt.Errorf("Active attribute not found")
	}
	fmt.Println("activeAttribute :", activeAttribute)

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

func isAccessGranted(ctx contractapi.TransactionContextInterface, queryString string) (string, error) {
	// queryString := fmt.Sprintf(`{"selector":{"id":"%v","docType":"ACCESS","doctorId":"%v"}}`, fromUser, toUser)
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return "nil", err
	}
	defer resultsIterator.Close()
	if resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		fmt.Println("QueryResult : ", queryResult)
		if err != nil {
			return "nil", fmt.Errorf("Error on fetching access record : %v", err.Error())
		}
		return queryResult.Key, nil
	}
	return "nil", nil
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
