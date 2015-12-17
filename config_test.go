package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func loadTestFile(content bool, comments bool) (*os.File, error) {
	f, err := ioutil.TempFile("", "testconfig")
	if err != nil {
		return nil, err
	}
	if content {
		if comments {
			_, err = f.WriteString("# COMMENT \nKEY=VALUE\n# RANDOM COMMENTS TO BE IGNORED")
			if err != nil {
				return nil, err
			}
		} else {
			_, err = f.WriteString("KEY=VALUE")
			if err != nil {
				return nil, err
			}
		}
	}
	return f, nil
}

func Test_DotEnvFile(t *testing.T) {
	f, err := loadTestFile(true, false)
	if err != nil {
		t.Fatal(err)
	}
	config := Load(f.Name())
	if v := config.Get("KEY", ""); v != "VALUE" {
		t.Fatal("Config did not return VALUE")
	} else {
		t.Logf("Config returned KEY=%s", v)
	}

	if config.Get("BAD_KEY", "DEFAULT") != "DEFAULT" {
		t.Fatal("Config did not return DEFAULT")
	}
}

func Test_DotEnvFileWithComments(t *testing.T) {
	f, err := loadTestFile(true, true)
	if err != nil {
		t.Fatal(err)
	}
	config := Load(f.Name())
	if v := config.Get("KEY", ""); v != "VALUE" {
		t.Fatal("Config did not return VALUE")
	} else {
		t.Logf("Config returned KEY=%s", v)
	}
}

func Test_DotEnvWithBadDelimiter(t *testing.T) {
	f, err := ioutil.TempFile("", "testconfig")
	if err != nil {
		t.Fatal(err)
	}

	_, err = f.WriteString("KEY:VALUE")
	if err != nil {
		t.Fatal(err)
	}

	config := Load(f.Name())
	if v := config.Get("KEY", ""); v != "" {
		t.Log("Bad delimiter ignored")
	}
}

func Test_DotEnvEnvironment(t *testing.T) {
	os.Setenv("KEY", "VALUE")
	defer os.Setenv("KEY", "")
	config := Load("")

	if v := config.Get("KEY", ""); v != "VALUE" {
		t.Fatal("Config did not return VALUE")
	} else {
		t.Logf("Config returned KEY=%s", v)
	}
	if config.Get("BAD_KEY", "DEFAULT") != "DEFAULT" {
		t.Fatal("Config did not return DEFAULT")
	}
}
