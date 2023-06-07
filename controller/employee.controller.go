package controller

import (
	"crud/dbconfig"
	"crud/models"
	repo "crud/service"
	"crud/service/employee"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Employee struct {
	repo repo.EmpRepo
}

func NewEmployeeHandler(db *dbconfig.DB) *Employee {
	return &Employee{
		repo: employee.NewEmpRepo(db.SQL),
	}
}

type TransactionBody struct {
	SenderID   int64   `json:"senderId"`
	ReceiverID int64   `json:"receiverId"`
	Amount     float64 `json:"amount"`
}

func (e *Employee) GetEmployeesList(w http.ResponseWriter, r *http.Request) {
	res, err := e.repo.Fetch(r.Context())
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondwithJSON(w, 200, res)
}

func (e *Employee) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {

		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}
	res, err := e.repo.GetByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found")
		return
	}
	respondwithJSON(w, 200, res)
}

func fileUpload(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", err
	}

	file, header, err := r.FormFile("picture")
	if err != nil {
		return "", err
	}

	defer file.Close()

	fileExt := strings.Split(header.Filename, ".")[1]

	imgName := fmt.Sprintf("uploaded-*.%s", fileExt)

	tempFile, err := ioutil.TempFile("images", imgName)
	if err != nil {
		return "", err
	}

	err = tempFile.Close()
	if err != nil {
		return "", err
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	tempFile.Write(fileBytes)
	return tempFile.Name(), nil
}

// Create Employee
func (e *Employee) CreateEmployee(w http.ResponseWriter, r *http.Request) {

	file, err := fileUpload(r)

	req := models.Employee{}
	req.Name = r.FormValue("name")
	req.Phone = r.FormValue("phone")
	req.Job = r.FormValue("job")
	req.Country = r.FormValue("country")
	req.City = r.FormValue("city")
	req.Postalcode, _ = strconv.ParseInt(r.FormValue("postalcode"), 10, 64)
	if err == nil {
		req.Picture = file
	} else {
		fmt.Errorf(err.Error(), err)
	}

	res, err := e.repo.Create(r.Context(), &req)
	if err != nil {

		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	// On succes
	respondwithJSON(w, 200, res)
}

// Update employee
func (e *Employee) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	req := models.Employee{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {

		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	res, err := e.repo.Update(r.Context(), &req, id)

	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	// On succes
	respondwithJSON(w, http.StatusAccepted, res)
}

// Delete employee
func (e *Employee) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {

		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}

	res, err := e.repo.Delete(r.Context(), id)
	if err != nil {

		respondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}
	respondwithJSON(w, http.StatusOK, res)
}

// Transaction
func (e *Employee) TransactMoney(w http.ResponseWriter, r *http.Request) {
	req := TransactionBody{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Bad request")
		return
	}
	err = e.repo.Trasaction(r.Context(), req.Amount, req.SenderID, req.ReceiverID)
	if err != nil {
		respondWithError(w, http.StatusConflict, "Conflict")
		return
	}
	respondwithJSON(w, http.StatusOK, "Transaction success")
}
