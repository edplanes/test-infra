package tests

import (
	"strings"
	"time"

	"github.com/tebeka/selenium"
)

func TestLoginFlowFactory() TestCase {
	return TestCase{
		Test: login,
	}
}

func login(driver selenium.WebDriver) error {
	err := driver.Get("http://localhost:4200/login")
	if err != nil {
		return err
	}

	email, err := driver.FindElement(selenium.ByCSSSelector, "input[type=\"email\"]")
	if err != nil {
		return err
	}
	email.SendKeys("admin@localhost.com")
	password, err := driver.FindElement(selenium.ByCSSSelector, "input[type=\"password\"]")
	if err != nil {
		return err
	}
	password.SendKeys("admin")

	submit, err := driver.FindElement(selenium.ByCSSSelector, "button[type=\"submit\"]")
	if err != nil {
		return err
	}
	submit.Click()

	time.Sleep(1 * time.Second)

	url, err := driver.CurrentURL()
	if err != nil {
		return err
	}
	if !strings.HasSuffix(url, "/user") {
		return err
	}

	return nil
}
