package jenkins

import (
	"encoding/json"
	"os"
	"testing"
)

func TestApiJobListResponseFilterByProperty(t *testing.T) {
	f, err := os.Open("fixtures/joblist.json")
	if err != nil {
		t.Fatal(err)
	}

	jobList := new(ApiJobListResponse)
	err = json.NewDecoder(f).Decode(jobList)
	if err != nil {
		t.Fatal(err)
	}

	jobs := jobList.FilterByProperty("SvcopRef", "origin/master")
	if len(jobs) != 3 {
		t.Error("job list should contain 3 jobs, contains: ", len(jobs))
	}
}
