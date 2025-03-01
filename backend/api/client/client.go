package client

import (
	"cli-depend/domain/cli-dep"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func GetDependencies(packageName string, version string) (cli_dep.DependenciesDto, error) {
	encodedUrl := "https://api.deps.dev/v3/systems/go/packages/" + url.QueryEscape(packageName) + "/versions/" + version + ":dependencies"
	body, err := makeGetRequest(encodedUrl)
	if err != nil {
		//log.Fatalf("Error fetching dependencies: %v", err)
		return cli_dep.DependenciesDto{}, err
	}

	var dependencies cli_dep.DependenciesDto
	err = json.Unmarshal(body, &dependencies)
	if err != nil {
		fmt.Println("Errore durante la decodifica:", err)
		return cli_dep.DependenciesDto{}, err
	}
	return dependencies, nil
}

func GetScore(packageName string) (float32, error) {
	baseUrl := "https://api.deps.dev/v3/projects/"
	encodedParam := url.QueryEscape(packageName)

	if !strings.Contains(encodedParam, "github.com") {
		return -1, nil
	}
	body, err := makeGetRequest(baseUrl + encodedParam)
	if err != nil {
		//log.Fatalf("Error fetching dependencies: %v", err)
		return 0, err
	}

	var score cli_dep.ScoreCardDto
	err = json.Unmarshal(body, &score)
	if err != nil {
		fmt.Println("Errore durante la decodifica: dello score per il package ", err)
		fmt.Println(packageName)
		return 0, err
	}
	return score.OverallScore, nil
}

func makeGetRequest(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("errore nella creazione della richiesta: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("errore nell'esecuzione della richiesta: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("errore nella lettura della risposta: %v", err)
	}
	return body, nil
}
