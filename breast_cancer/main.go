package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

func readFile(path string) [][]string {
	file, err := os.Open(path)

	if err != nil {
		log.Fatalf("Could not open file: ", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil {
		log.Fatalln("Could not read file: ", err)
	}

	return record
}

func main() {
	records := readFile("data.csv")

	classes := getDistinctClasses(getLastColumns(records, len(records[0]) - 1))

	var train [][]string
	var test [][]string

	for i := range classes {
		values := getElementsByClass(records, classes[i])

		trainData, testData := getDataSetPart(values, 0.9)

		train = append(train, trainData...)
		test = append(test, testData...)
	}

	var hits int
	//for each line in test lines
	for i := range test {
		result := classify(train, test[i], 11)

		columnIndex := len(test[i]) - 1

		fmt.Println("tumor: ", test[i][columnIndex], " classificado como: ", result)

		if result == test[i][columnIndex] {
			hits++
		}
	}

	fmt.Println("Total de dados: ", len(records))
	fmt.Println("Total de treinamento: ", len(train))
	fmt.Println("Total de testes: ", len(test))
	fmt.Println("Total de acertos: ", hits)
	fmt.Println("Porcentagem de acertos: ", 100 * hits / len(test), "%")
}

// All the classes are in the last column of the items
func getLastColumns(elements [][]string, columnIndex int) []string {
	columns := make([]string, 0)

	for i := range elements {
		columns = append(columns, (elements)[i][columnIndex])
	}

	return columns
}

func getDistinctClasses(items []string) []string {
	encountered := make(map[string]bool)
	result := make([]string, 0)

	for i := range items {
		if !encountered[items[i]] {
			result = append(result, items[i])

			encountered[items[i]] = true
		}
	}

	return result
}

func getElementsByClass(records [][]string, class string) [][]string {
	newRecords := make([][]string, 0)

	for i := range records {
		if records[i][len(records[0]) - 1] == class {
			newRecords = append(newRecords, records[i])
		}
	}

	return newRecords
}

func getDataSetPart(records [][]string, percentage float32) ([][]string, [][]string) {
	breakingPoint := int(float32(len(records)) * percentage)

	newRecords := make([][]string, breakingPoint)

	for i := breakingPoint - 1; i >= 0; i-- {
		newRecords[i] = records[i]
	}

	totalItems := len(records) - 1

	residue := make([][]string, 0)

	for j := totalItems; j > breakingPoint - 1; j-- {
		residue = append(residue, records[j])
	}

	return newRecords, residue
}

func getEuclidianDistance(pi, qi []string) float64 {
	sum := 0.0

	for i := len(pi) -1; i > 0; i-- {

		pif, _ := strconv.ParseFloat(pi[i - 1], 32)
		qif, _ := strconv.ParseFloat(qi[i - 1], 32)

		sum += math.Pow(pif - qif, 2)
	}

	return math.Sqrt(sum)
}


func getMapKeys(m map[float64]string) []float64 {
	keys := make([]float64, len(m))

	i := 0
	for key := range m {
		keys[i] = key
		i++
	}

	return keys
}

func getMapValues(m map[float64]string) []string {
	values := make([]string, len(m))

	i := 0
	for _, value := range m {
		values[i] = value
	}

	return values
}

func getKnn(list map[float64]string, k int) map[float64]string {
	keys := getMapKeys(list)

	sort.Float64s(keys)

	sortedMap := make(map[float64]string, len(list))

	for i, key := range keys {
		if i < k {
			//add the keyValue into a new keyValue list
			sortedMap[key] = list[key]
		} else {
			break
		}
	}

	return sortedMap
}

func getPredominantClass(knn map[float64]string) string {
	classes := getDistinctClasses(getMapValues(knn))


	var predominantClassCount int
	var predominantClass string

	for i := range classes {

		var classCount int

		for key := range knn {

			if knn[key] == classes[i] {
				classCount++
			}

		}

		if predominantClassCount < classCount {
			predominantClassCount = classCount
			predominantClass = classes[i]
		}

	}

	return predominantClass
}

func classify(train [][]string, valueToPredict []string, k int) string {
	distances := make(map[float64]string)

	for i := range train {
		class := train[i][len(train[i]) - 1]
		distance := getEuclidianDistance(train[i], valueToPredict)

		distances[distance] = class
	}

	knn := getKnn(distances, k)

	return getPredominantClass(knn)
}



