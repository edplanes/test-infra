package main

import (
	"fmt"
	"log"

	"github.com/edplanes/test-infra/cmd/selenium/tests"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

var testCases = []tests.TestCase{
	tests.TestLoginFlowFactory(),
}

func main() {
	service, err := selenium.NewChromeDriverService("chromedriver", 4444)
	if err != nil {
		log.Fatal("Error starting chrome driver:", err)
	}
	defer service.Stop()

	caps := selenium.Capabilities{}
	caps.AddChrome(chrome.Capabilities{Args: []string{"--headless=new"}})

	driver, err := selenium.NewRemote(caps, "")
	if err != nil {
		log.Fatal("Error creating remote driver", err)
	}

	err = driver.MaximizeWindow("")
	if err != nil {
		log.Fatal("Error:", err)
	}

	for _, tc := range testCases {
		fmt.Println("Testing:", tc.Name)

		if tc.Before != nil {
			fmt.Println(" Running before function")
			tc.Before()
		}

		fmt.Println(" Running test function")

		err := tc.Test(driver)
		passed := err == nil
		if passed {
			fmt.Println("Result: Passed")
		} else {
			fmt.Println("Result: Failed", err)
		}

		if tc.After != nil {
			fmt.Println(" Running after function")
			tc.After()
		}
	}

}
