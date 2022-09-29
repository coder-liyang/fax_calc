package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("输入 exit 退出程序")
	for {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("请输入税前工资,单位元:")
		sentence, err := buf.ReadBytes('\n')
		if err != nil {
			panic(err.Error())
		}
		//最后肯定是一个换行,截取掉(在windows下不行,Windows的换行符是CRLF)
		//str := string(sentence[:len(sentence)-1])
		str := strings.TrimSpace(string(sentence))
		if str == "exit" {
			return
		}
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Println(str, "?")
			fmt.Println("来!你给我算一个试试!")
		} else {
			tax(num)
		}
	}
}

//2019年1月1日起，国家推出新的个人所得税政策，起征点上调值5000元。
//也就是说税前工资扣除三险一金（三险一金数额假设是税前工资的10%）后如果不足5000元，则不交税。
//如果大于5000元，那么大于5000元的部分按梯度交税，具体梯度比例如下：
//	0 ~ 3000元的部分，交税3%
//	3000 ~ 12000元的部分，交税10%
//	12000 ~ 25000的部分 ， 交税20%
//	25000 ~ 35000的部分，交税25%
//	35000 ~ 55000的部分，交税30%
//	55000 ~ 80000的部分，交税35%
//	超过80000的部分，交税45%
func tax(num float64) {
	fmt.Printf("你输入的金额是:%.2f元\n", num)
	numAfterInsurance := num * 0.9
	fmt.Printf("扣完三险一金后剩余:%.2f\n", numAfterInsurance)

	//应交税工资部分
	var numTax float64
	//应交税总额
	var taxTotal float64

	if numAfterInsurance <= 5000 {
		fmt.Println("恭喜,你不用交税.继续保持.")
		return
	} else {
		numTax = numAfterInsurance - 5000
		fmt.Printf("扣除免税部分5000元后,应扣税部分工资为:%.2f\n", numTax)
	}

	taxTotal += faxLowToHigh(7, numTax)
	taxTotal += faxLowToHigh(6, numTax)
	taxTotal += faxLowToHigh(5, numTax)
	taxTotal += faxLowToHigh(4, numTax)
	taxTotal += faxLowToHigh(3, numTax)
	taxTotal += faxLowToHigh(2, numTax)
	taxTotal += faxLowToHigh(1, numTax)
	fmt.Printf("应交税总额为%.2f\n", taxTotal)
}

func faxLowToHigh(taxLevel int, num float64) float64 {
	var moreThanLow float64
	var moreThanLowTax float64
	taxMap := map[int]struct {
		Rate  float64
		Low   float64
		Heigh float64
	}{
		1: {0.03, 0, 3000},
		2: {0.10, 3000, 12000},
		3: {0.20, 12000, 25000},
		4: {0.25, 25000, 35000},
		5: {0.30, 35000, 55000},
		6: {0.35, 55000, 80000},
		7: {0.45, 80000, 0},
	}
	rate, ok := taxMap[taxLevel]
	if !ok {
		panic("错误的交税等级")
	}
	switch taxLevel {
	case 7:
		//这个挡位比较特殊,其他都是一个范围,只有这个上不封顶
		moreThanLow = num - rate.Low
		if moreThanLow > 0 {
			moreThanLowTax = moreThanLow * rate.Rate
			fmt.Printf("在%.0f元以上档,你应交税:%.2f\n", rate.Low, moreThanLowTax)
		} else {
			fmt.Printf("在%.0f元以上档,你不用交税.\n", rate.Low)
		}
	default:
		if num < rate.Low {
			moreThanLowTax = 0
			fmt.Printf("在%.0f元档,你不用交税.\n", rate.Low)
		} else if num > rate.Heigh {
			//收入超过了此档最大值
			moreThanLow = rate.Heigh - rate.Low
		} else {
			//收入没超过此档最大值
			moreThanLow = num - rate.Low
		}
		moreThanLowTax = moreThanLow * rate.Rate
		if moreThanLowTax > 0 {
			fmt.Printf("在%.0f - %.0f元档,你应交税:%.2f\n", rate.Low, rate.Heigh, moreThanLowTax)
		}
	}
	return moreThanLowTax
}
