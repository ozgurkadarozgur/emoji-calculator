package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type error interface {
	Error()
}

type DivisionError struct{}

type UnquoteCodePointError struct{}

func (e *DivisionError) Error() {
	emoji, err := unquoteCodePoint("\\" + division_error_emoji)
	if err != nil {
		fmt.Println("parse to emoji error:", err)
	}
	fmt.Printf("%s\n", emoji)
}

func (e *UnquoteCodePointError) Error() {
	fmt.Println("unquote code point error:", e)
}

type Operation struct {
	Symbol  string
	Index   int
	Number1 int
	Number2 int
}

const (
	division_error_emoji = "U+1F937"
)

var (
	operations      = []string{"+", "-", "*", "/"}
	operation_texts = []string{"plus", "minus", "times", "divided by"}
	emoji_list      map[string]string
	result_set      map[string]string
)

func main() {

	//s := "8 times 2"

	fmt.Println("Type your input to calculate: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()

	initEmojis()

	sorted_result_set := initResultSet()

	s = ConvertInputWords2Operation(s)

	fmt.Println("converted input : ", s)

	equation := ""
	//equation := "4 + 6 - 3 * 2"
	r := []rune(s)

	for _, item := range r {
		code := fmt.Sprintf("%U", item)
		fmt.Println(string(item), code)
		equation += emoji_list[code]
	}
	fmt.Println("initial_equation", equation)

	trimmed_equation := strings.Replace(equation, " ", "", -1)
	fmt.Println("trimmed_equation", trimmed_equation)

	equation_array := ConvertEquation2Array(equation)
	fmt.Println("equation_array", equation_array)

	priority_calculation_result, err := CalculatePriorityOperations(equation_array)
	if err != nil {
		switch e := err.(type) {
		case *DivisionError:
			{
				e.Error()
				break
			}
		default:
			break
		}
		return
	}
	fmt.Println("priority_calculation_result_array", priority_calculation_result)

	result_array, err := CalculateEquation(priority_calculation_result)
	if err != nil {
		switch e := err.(type) {
		case *DivisionError:
			{
				e.Error()
				return
			}
		default:
			break
		}
	}
	fmt.Println("result_array", result_array)

	result := result_array[0]
	fmt.Println("result : ", result)

	emoji_result := ConvertResult2Emoji(result, sorted_result_set)
	fmt.Println("emoji_result", emoji_result)
}

func initEmojis() {
	emoji_list = map[string]string{
		"0": "0",
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
		"5": "5",
		"6": "6",
		"7": "7",
		"8": "8",
		"9": "9",

		"+": "+",
		"-": "-",
		"*": "*",
		"/": "/",

		"U+0030": "0",
		"U+0031": "1",
		"U+0032": "2",
		"U+0033": "3",
		"U+0034": "4",
		"U+0035": "5",
		"U+0036": "6",
		"U+0037": "7",
		"U+0038": "8",
		"U+0039": "9",

		"U+1F51F": "10",

		"U+002B": "+",
		"U+002D": "-",
		"U+2716": "*",
		"U+002F": "/",

		"U+0078": "*",
		"U+0025": "/",

		"U+2795": "+",
		"U+2796": "-",
		"U+2797": "/",
		"U+002A": "*",

		"plus":       "+",
		"minus":      "-",
		"times":      "*",
		"divided by": "/",

		"U+1F3B1": "8",
		"U+1F4AF": "100",
	}
}

func initResultSet() []int {
	result_set = map[string]string{
		"0":   "U+0030",
		"1":   "U+0031",
		"2":   "U+0032",
		"3":   "U+0033",
		"4":   "U+0034",
		"5":   "U+0035",
		"6":   "U+0036",
		"7":   "U+0037",
		"8":   "U+0038",
		"9":   "U+0039",
		"10":  "U+1F51F",
		"100": "U+1F4AF",
	}

	sorted_arr := make([]int, 0, len(result_set))
	for key, _ := range result_set {
		int_key, _ := strconv.Atoi(key)
		sorted_arr = append(sorted_arr, int_key)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sorted_arr)))
	fmt.Println("sorted_arr", sorted_arr)
	return sorted_arr
}

func CalculatePriorityOperations(equation_array []string) ([]string, error) {
	var new_arr []string
	fmt.Println("priority_equation_array", equation_array)
	fmt.Println("new_array", new_arr)
	if hasPriorityOperation(equation_array) {
		for index, item := range equation_array {
			i := string(item)
			if isOperation(i) && (i == "*" || i == "/") {
				number1, err_number_1 := strconv.Atoi(equation_array[index-1])
				if err_number_1 != nil {
					panic(err_number_1)
				}

				number2, err_number_2 := strconv.Atoi(equation_array[index+1])
				if err_number_2 != nil {
					panic(err_number_2)
				}

				operation := Operation{
					Index:   index,
					Symbol:  i,
					Number1: number1,
					Number2: number2,
				}
				r, err := CalculateOperation(operation)
				if err != nil {
					return nil, err
				}
				result := strconv.Itoa(r)
				first_part_arr := equation_array[0 : operation.Index-1]
				second_part_arr := equation_array[operation.Index+2 : len(equation_array)]
				//new_arr = equation_array[operation.Index+2 : len(equation_array)]
				new_arr = append(new_arr, first_part_arr...)
				new_arr = append(new_arr, result)
				new_arr = append(new_arr, second_part_arr...)
				fmt.Println(new_arr)
				break
			}
		}
		return CalculatePriorityOperations(new_arr)
	}
	return equation_array, nil
}

func CalculateEquation(equation_array []string) ([]string, error) {
	var new_arr []string
	if hasOperation(equation_array) {
		for index, item := range equation_array {
			i := string(item)
			if isOperation(i) {
				number1, err_number_1 := strconv.Atoi(equation_array[index-1])
				if err_number_1 != nil {
					panic(err_number_1)
				}

				number2, err_number_2 := strconv.Atoi(equation_array[index+1])
				if err_number_2 != nil {
					panic(err_number_2)
				}

				operation := Operation{
					Index:   index,
					Symbol:  i,
					Number1: number1,
					Number2: number2,
				}
				r, err := CalculateOperation(operation)
				if err != nil {
					return nil, err
				}
				result := strconv.Itoa(r)
				new_arr = equation_array[operation.Index+2 : len(equation_array)]
				new_arr = append([]string{result}, new_arr...)
				fmt.Println(new_arr)
				break
			}
		}
		return CalculateEquation(new_arr)
	}
	return equation_array, nil
}

func ConvertResult2Emoji(s string, sorted_result_set []int) string {
	fmt.Println("myss", sorted_result_set)
	result := ""

	for key, _ := range result_set {
		if key == s {
			result = result_set[key]
			break
		}
	}

	if result != "" {
		fmt.Println("r : ", result)
		emoji, err := unquoteCodePoint("\\" + result)
		if err != nil {
			fmt.Println("parse to emoji error:", err)
		}
		fmt.Printf("%s\n", emoji)

		return emoji
	}

	for _, item := range sorted_result_set {
		//for i := len(sorted_result_set) - 1; i >= 0; i-- {
		//str_item := strconv.Itoa(sorted_result_set[i])
		str_item := strconv.Itoa(item)
		if strings.Contains(s, str_item) {
			emoji, err := unquoteCodePoint("\\" + result_set[str_item])
			if err != nil {
				fmt.Println("error while parsing result to emoji", err)
			}
			s = strings.Replace(s, str_item, emoji, -1)
		}
		//}
	}

	/*
		for key, _ := range result_set {
			if strings.Contains(s, key) {
				emoji_code := result_set[key]
				emoji, err := unquoteCodePoint("\\" + emoji_code)
				if err != nil {
					fmt.Println("parse result to emoji error", err)
				}
				//s = strings.Replace(s, key, emoji, -1)
				result += emoji
			}
		}
	*/
	return s
}

func unquoteCodePoint(s string) (string, error) {
	r, err := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	if err != nil {
		return "", &UnquoteCodePointError{}
	}
	return string(r), nil
}

func ConvertInputWords2Operation(s string) string {
	for _, item := range operation_texts {
		if strings.Contains(s, item) {
			s = strings.Replace(s, item, emoji_list[item], -1)
		}
	}
	return s
}

func ConvertEquation2Array(equation string) []string {
	var arr []string
	var number string
	for index, item := range equation {
		i := strings.Replace(string(item), " ", "", -1)
		if isDigit(i) {
			number += i
			if index == len(equation)-1 {
				arr = append(arr, number)
			}
		} else {
			arr = append(arr, number, i)
			number = ""
		}
	}
	return arr
}

func CalculateOperation(operation Operation) (int, error) {
	fmt.Println(fmt.Sprintf("operationing : %d %s %d", operation.Number1, operation.Symbol, operation.Number2))
	result := 0
	switch operation.Symbol {
	case "+":
		{
			result = operation.Number1 + operation.Number2
			break
		}
	case "-":
		{
			result = operation.Number1 - operation.Number2
			break
		}
	case "*":
		{
			result = operation.Number1 * operation.Number2
			break
		}
	case "/":
		{
			if operation.Number2 == 0 {
				return -1, &DivisionError{}
			}
			result = operation.Number1 / operation.Number2
			break
		}
	}
	return result, nil
}

func isDigit(s string) bool {
	if isOperation(s) {
		return false
	}
	return true
}

func isOperation(s string) bool {
	for _, item := range operations {
		if item == s {
			return true
		}
	}
	return false
}

func hasOperation(equation_array []string) bool {
	for _, item := range equation_array {
		if isOperation(string(item)) {
			return true
		}
	}
	return false
}

func hasPriorityOperation(equation_array []string) bool {
	for _, item := range equation_array {
		i := string(item)
		if isOperation(i) && (i == "*" || i == "/") {
			return true
		}
	}
	return false
}
