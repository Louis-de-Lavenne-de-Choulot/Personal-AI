package main

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

func calculationRec(e string) (float64, error) {
	//remove spaces
	e = strings.ReplaceAll(e, " ", "")
	// loop through the string and find operators, add them to a stack by importance
	_, err := strconv.Atoi(e)
	priority := 0
	a := 0
	for err != nil {
		if strings.Contains(e, "^") {
			for i := len(e) - 1; i >= 0; i-- {
				//if ^ is found then calculate the result and replace the string with the result
				if e[i] == '^' {
					//find the numbers before and after the operator
					//find the number before the operator
					before, after := beforeAndAfter(e, i)
					//calculate the result
					//convert the string to float
					beforeNum, err := strconv.ParseFloat(e[before:i], 64)
					if err != nil {
						return 0, err
					}
					afterNum, err := strconv.ParseFloat(e[i+1:after+1], 64)
					if err != nil {
						return 0, err
					}

					res := math.Pow(float64(beforeNum), float64(afterNum))
					//replace the string between the brackets with the result
					e = e[:before] + strconv.FormatFloat(res, 'f', -1, 64) + e[after+1:]
					//process the stack again
					break
				}
			}
		} else {

			if !strings.Contains(e, "(") && priority == 0 {
				priority++
				a = 0
			}
			if !strings.Contains(e, "%") && !strings.Contains(e, "/") && !strings.Contains(e, "*") {
				if priority == 1 {
					priority++
					a = 0
				}
			}
			println(e, priority)
			for a = 0; a < len(e)-1; a++ {
				c := e[a]
				if c == '(' && priority == 0 {
					//find the closing bracket with indexof
					//process the stack with the string between the brackets
					ind := strings.Index(e, ")")
					//get next ( and check if it is before the closing bracket
					//if it is, get the next closing bracket
					indNext := strings.Index(e[a+1:], "(") + a + 1
					println(ind, indNext)
					//create iteBuffer that will contain ind and indNext
					iteBuffer := []int{ind, indNext}
					for indNext != -1 && indNext < ind {
						ind += strings.Index(e[ind+1:], ")") + 1
						if ind == len(e)-1 {
							break
						}
						indNext += strings.Index(e[indNext+1:], "(") + 1
						//if iteBufffer values compared to ind and indNext are the same, it means that there is no more brackets
						if iteBuffer[0] == ind && iteBuffer[1] == indNext {
							break
						}
						iteBuffer[0] = ind
						iteBuffer[1] = indNext
					}

					if ind == -1 {
						return 0, errors.New("no closing bracket")
					}
					//process the stack with the string between the brackets
					//replace the string between the brackets with the result
					//process the stack again
					res, _ := calculationRec(e[a+1 : ind])
					//if ( is not the first character and if the character before is a number then add a * before the result
					if a != 0 && int(e[a-1]) >= 48 && int(e[a-1]) <= 57 {
						e = e[:a] + "*" + strconv.FormatFloat(res, 'f', -1, 64) + e[ind+1:]
					} else {
						e = e[:a] + strconv.FormatFloat(res, 'f', -1, 64) + e[ind+1:]
					}
					break
				} else if c == '*' || c == '/' || c == '%' {
					if priority == 1 {
						if a == 0 {
							continue
						}
						before, after := beforeAndAfter(e, a)
						res := 0.0
						switch c {
						case '*':
							res, _ = strconv.ParseFloat(e[before:a], 64)
							toMult, _ := strconv.ParseFloat(e[a+1:after+1], 64)
							res *= toMult
							e = e[:before] + strconv.FormatFloat(res, 'f', -1, 64) + e[after+1:]
						case '/':
							res, _ = strconv.ParseFloat(e[before:a], 64)
							toDiv, _ := strconv.ParseFloat(e[a+1:after+1], 64)
							res /= toDiv
							e = e[:before] + strconv.FormatFloat(res, 'f', -1, 64) + e[after+1:]
						case '%':
							res, _ = strconv.ParseFloat(e[before+1:a], 64)
							toMod, _ := strconv.ParseFloat(e[a+1:after+1], 64)
							res = float64(int(res) % int(toMod))
							e = e[:before+1] + strconv.FormatFloat(res, 'f', -1, 64) + e[after+1:]
						default:
							return 0, errors.New("unknown operator")
						}

						break
					}
				} else if c == '+' || c == '-' {
					if priority == 2 {
						if a == 0 {
							continue
						}
						before, after := beforeAndAfter(e, a)
						//do the calculation with number before and after the operator
						//replace the string between the brackets with the result
						//process the stack again
						res := 0.0
						switch c {
						case '+':
							res, _ = strconv.ParseFloat(e[before:a], 64)
							toAdd, _ := strconv.ParseFloat(e[a+1:after+1], 64)
							res += toAdd
						case '-':
							res, _ = strconv.ParseFloat(e[before:a], 64)
							toSub, _ := strconv.ParseFloat(e[a+1:after+1], 64)
							res -= toSub
						default:
							return 0, errors.New("unknown operator")
						}
						e = e[:before] + strconv.FormatFloat(res, 'f', -1, 64) + e[after+1:]
						break
					}
				}

			}
		}
		_, err = strconv.ParseFloat(e, 64)
	}
	//e to float
	res, _ := strconv.ParseFloat(e, 64)
	//round to 2 decimal places
	return res, nil
}

func calculation(e string) (float64, error) {
	res, err := calculationRec(e)
	res = float64(int(res*100)) / 100
	return res, err
}

func beforeAndAfter(e string, a int) (int, int) {
	tol := 0
	//do the calculation with number before and after the operator
	//replace the string between the brackets with the result
	//process the stack again
	//find the number before the operator
	before := 0
	for b := a - 1; b >= 0; b-- {
		if e[b] == '+' || e[b] == '*' || e[b] == '/' || e[b] == '%' || e[b] == '(' || e[b] == ')' || e[b] == '^' {
			before = b + 1
			break
		} else if e[b] == '-' {
			before = b
			break
		}
	}
	tol = 0
	//find the number after the operator
	after := 0
	for b := a + 1; b < len(e); b++ {
		if b == len(e)-1 || e[b] == '+' && tol >= 1 || e[b] == '-' && tol >= 1 || e[b] == '*' || e[b] == '/' || e[b] == '%' || e[b] == ')' || e[b] == '(' || e[b] == '^' {
			after = b
			//if is not the end of the string then after --
			if len(e)-1 > b {
				after--
			}
			break
		} else if e[b] == '-' && tol == 0 || int(e[b]) >= 48 && int(e[b]) <= 57 {
			tol++
		}
	}

	return before, after
}
