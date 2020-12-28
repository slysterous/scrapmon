package cobra_test

import (
	"bytes"
	"github.com/slysterous/scrapmon/internal/cobra"
	"io/ioutil"
	"strings"
	"testing"
)

func TestNewPurgeCommand(t *testing.T) {
	cc:=cobra.NewClient()
	purgeLogic:=func ()error{
		return nil
	}
	purgeCommand:=cc.NewPurgeCommand(purgeLogic)

	want:="purge"
	if purgeCommand.Use!=want{
		t.Errorf("want: %s, got: %s",want,purgeCommand.Use)
	}

	want="Purges db and filesystem storage"
	if purgeCommand.Short!=want{
		t.Errorf("want: %s, got: %s",want,purgeCommand.Short)
	}

	if purgeCommand.SilenceErrors != true {
		t.Errorf("want: %v, got: %v",true,purgeCommand.SilenceErrors)
	}

	if purgeCommand.RunE ==nil{
		t.Errorf("expected RunE not to be nil, nil given")
	}
}

func TestNewStartCommand(t *testing.T) {
	cc:=cobra.NewClient()
	startLogic:=func (fromCode string, iterations int, workerNumber int)error{
		return nil
	}
	startCommand,err:=cc.NewStartCommand(startLogic)
	if err!=nil{
		t.Errorf("unexpected error occured, err: %v",err)
	}
	want:="start"
	if startCommand.Use!=want{
		t.Errorf("want: %s, got: %s",want,startCommand.Use)
	}

	want="Starts scraping images from imgur"
	if startCommand.Short!=want{
		t.Errorf("want: %s, got: %s",want,startCommand.Short)
	}

	if startCommand.SilenceErrors != true {
		t.Errorf("want: %v, got: %v",true,startCommand.SilenceErrors)
	}

	if startCommand.RunE ==nil{
		t.Errorf("expected RunE not to be nil, nil given")
	}
}


func TestExecutePurgeCommand(t *testing.T){
	cc:=cobra.NewClient()
	purgeLogic:=func ()error{
		return nil
	}
	purgeCommand:=cc.NewPurgeCommand(purgeLogic)
	b := bytes.NewBufferString("")
	purgeCommand.SetOut(b)
	err:=purgeCommand.Execute()
	if err!=nil{
		t.Errorf("unexpected error occured, err: %v",err)
	}
	_, err = ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
}

var executeStartTests = []struct {
	description string
	in []string
	wantError bool
	out string
}{
	{
		"Success",
		[]string{"--workers","4","--from","a","--iterations","10"},false,""},
	{
		"Failed to parse workers",
		[]string{"--workers","text","--from","a","--iterations","10"},
		true,
		"workers provided was not a number",
	},
	{
		"Failed to parse iterations",
		[]string{"--workers","4","--from","a","--iterations","text"},
		true,
		"count provided was not a number",
	},
	{
		"No Workers provided",
		[]string{"--from","a","--iterations","10"},
		true,
		"required flag(s) \"workers\" not set",
	},
	{
		"Negative amount of workers provided",
		[]string{"--workers","-4","--from","a","--iterations","10"},
		true,
		"workers have to be at least 1",
	},
}


func TestExecuteStartCommand(t *testing.T){

	for _, tt :=range executeStartTests {
		t.Run(tt.description,func (t *testing.T){
			cc:=cobra.NewClient()
			startLogic:=func (fromCode string, iterations int, workerNumber int)error{
				return nil
			}
			startCommand,err:=cc.NewStartCommand(startLogic)
			if err!=nil{
				t.Errorf("unexpected error occured, err: %v",err)
			}
			startCommand.SetArgs(tt.in)
			b := bytes.NewBufferString("")
			startCommand.SetOut(b)
			err=startCommand.Execute()

			if tt.wantError && err==nil{
				t.Errorf("expected error got nil")
			}

			if !tt.wantError && err!=nil{
				t.Errorf("unexpected error occured, err: %v",err)
			}

			if  tt.wantError && !strings.Contains(err.Error(),tt.out){
				t.Fatalf("wanted error to contain: %s, got: %s",tt.out,err.Error())
			}
			_, err = ioutil.ReadAll(b)
			if err != nil {
				t.Fatal(err)
			}
		})
	}

}