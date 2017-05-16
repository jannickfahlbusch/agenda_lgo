package agenda_lgo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseURL = "https://agenda-lgo.de/api"

// Authentication is used to login
type Authentication struct {
	Email    string
	Password string
}

// Represents one salary statement
type Document struct {
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	Name         string `json:"name"`
	DownloadPath string `json:"downloadPath"`
	Type         string `json:"type"`
	Read         bool   `json:"read"`
	CreatedAt    int64  `json:"createdAt"`
}

// Contains general information about the employee & employer
type DocumentResponse []struct {
	ID            string      `json:"id"`
	Employee      string      `json:"employee"`
	Employer      string      `json:"employer"`
	ActivationKey interface{} `json:"activationKey"`
	DocumentList  []Document  `json:"documents"`
}

// LGO represents the API of "Agenda: Lohn- und Gehaltsdokumente"
type LGO struct {
	client       *http.Client
	sessionToken string
	authFilePath string
	outDir       string
}

// URPResponse The response from "Agenda LGO" which contains the session-token
type URPResponse struct {
	URP string `json:"urp"`
}

// Instanciates a new LGO-instance
func NewLGO(authFilePath, outDir string) *LGO {
	lgo := &LGO{
		authFilePath: authFilePath,
		outDir:       outDir,
	}
	transport := &http.Transport{}
	lgo.client = &http.Client{
		Transport: transport,
	}

	return lgo
}

// Saves the document in the specified out-path
func (lgo *LGO) SaveDocument(document Document) error {
	downloadPath := lgo.generateURL(document.DownloadPath + "/" + document.Name)

	req, err := http.NewRequest("GET", downloadPath, nil)
	if err != nil {
		return err
	}
	lgo.setHeaders(req)

	resp, err := lgo.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	out, err := os.Create(fmt.Sprintf("%s/%d-%s.pdf", lgo.outDir, document.Year, time.Month(document.Month)))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// Fetches the list of all available documents
func (lgo *LGO) FetchDocumentList() ([]Document, error) {
	// Fetch all documents
	req, err := http.NewRequest("GET", lgo.generateURL("/me/e"), nil)
	if err != nil {
		return nil, err
	}
	lgo.setHeaders(req)

	resp, err := lgo.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	documentResponse := DocumentResponse{}
	err = json.NewDecoder(resp.Body).Decode(&documentResponse)
	if err != nil {
		return nil, err
	}

	return documentResponse[0].DocumentList, nil
}

// Generates the neccessary reader for the login
func (lgo *LGO) generateAuthentication() (*strings.Reader, error) {
	reader, err := os.Open(lgo.authFilePath)
	if err != nil {
		return nil, err
	}

	auth := Authentication{}
	err = json.NewDecoder(reader).Decode(&auth)
	if err != nil {
		return nil, err
	}

	authStr := fmt.Sprintf("eml=%s&pwd=%s", auth.Email, auth.Password)

	return strings.NewReader(authStr), nil
}

// Authenticates the user with "Agenda: LGO"
func (lgo *LGO) Login() error {
	// First login

	authenticationReader, err := lgo.generateAuthentication()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", lgo.generateURL("/auth"), authenticationReader)
	if err != nil {
		return err
	}
	lgo.setHeaders(req)

	resp, err := lgo.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	urpResponse := URPResponse{}

	err = json.NewDecoder(resp.Body).Decode(&urpResponse)
	if err != nil {
		return err
	}

	lgo.sessionToken = urpResponse.URP

	// Strange, but we need a second login via GET
	req, err = http.NewRequest("GET", lgo.generateURL("/auth"), nil)
	if err != nil {
		return err
	}
	lgo.setHeaders(req)

	resp, err = lgo.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}

// Sets the neccessary headers
func (lgo *LGO) setHeaders(req *http.Request) {
	req.Header.Set("Origin", "https://agenda-lgo.de")
	req.Header.Set("User-Agent", "LGO-Downloader 0.1")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
}

// Generates the URL
func (lgo *LGO) generateURL(method string) string {
	return baseURL + method + lgo.sessionToken
}
