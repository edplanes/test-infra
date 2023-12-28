package tests

import "github.com/tebeka/selenium"

type TestCase struct {
	Name   string
	Before func()
	Test   func(driver selenium.WebDriver) error
	After  func()
}

type TesCaseFactory func() TestCase
