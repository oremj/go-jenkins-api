package jenkins

import (
	"encoding/json"
	"io"
	"net/http"
)

type Api struct {
	Username string
	Password string
	BaseURL  string // e.g., https://deploy.jenkins.com

	Client *http.Client
}

func NewApi(username, password, baseUrl string) *Api {
	return &Api{
		Username: username,
		Password: password,
		BaseURL:  baseUrl,
		Client:   &http.Client{},
	}
}

type ApiJobListResponse struct {
	Jobs []ApiJobs
}

type ApiJobs struct {
	Name     string `json:"name"`
	Property []struct {
		Parameters []struct {
			Name     string `json:"name"`
			Defaults struct {
				Value string `json:"value"`
			} `json:"defaultParameterValue"`
		} `json:"parameterDefinitions"`
	} `json:"property"`
}

func (j *Api) BuildURL(path string) string {
	return j.BaseURL + path
}

func (j *Api) Do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(j.Username, j.Password)
	return j.Client.Do(req)
}

func (j *Api) Get(v interface{}, path string) error {
	req, err := http.NewRequest("GET", j.BuildURL(path), nil)
	if err != nil {
		return err
	}

	return j.doRequest(v, req)
}

func (j *Api) Post(v interface{}, path string, body io.Reader) error {
	req, err := http.NewRequest("POST", j.BuildURL(path), body)
	if err != nil {
		return err
	}

	return j.doRequest(v, req)
}

func (j *Api) doRequest(v interface{}, req *http.Request) error {
	resp, err := j.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
}

func (j *Api) FetchJobList() (*ApiJobListResponse, error) {
	resp := new(ApiJobListResponse)

	err := j.Get(resp, "/api/json?pretty=true&tree=jobs[name,property[parameterDefinitions[name,defaultParameterValue[value]]]]")
	return resp, err
}

func (j *ApiJobListResponse) FilterByProperty(name, value string) []ApiJobs {
	jobs := make([]ApiJobs, 0)
	for _, job := range j.Jobs {
		for _, prop := range job.Property {
			for _, param := range prop.Parameters {
				if param.Name == name && param.Defaults.Value == value {
					jobs = append(jobs, job)
				}
			}
		}
	}
	return jobs
}
