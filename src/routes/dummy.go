package routes

import (
	"math"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"web-api-searching-and-pagination/src/models"

	"github.com/gin-gonic/gin"
)

var dummies []models.Dummy = generateDummies(100)

func getRandomName(index int) string {
	names := []string{"Lorem", "ipsum", "dolor", "sit", "amet", "consecutor"}
	namesLen := len(names)
	return names[index%namesLen]
}

func generateDummies(count int) []models.Dummy {
	slice := make([]models.Dummy, 0, count)
	for i := 0; i < count; i++ {
		dummy := []models.Dummy{
			{
				Id:     i,
				Name:   getRandomName(rand.Int()),
				Number: rand.Intn(200),
			},
		}
		slice = append(slice, dummy...)
	}
	return slice
}

func getDummies(c *gin.Context) {
	startRaw := c.Query("start")
	lengthRaw := c.Query("length")
	start := 0
	length := math.MaxInt

	if startRaw != "" {
		startParsed, err := strconv.Atoi(startRaw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "could not parse 'start' query param",
			})
			return
		}
		start = startParsed
	}

	if lengthRaw != "" {
		lengthParsed, err := strconv.Atoi(lengthRaw)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "could not parse 'length' query param",
			})
			return
		}
		length = lengthParsed
	}

	nameFilter := c.Query("name")
	filteredByName, nameFilterError := filterDummiesByName(dummies, nameFilter)
	if nameFilterError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not filter by query param 'name'",
		})
		return
	}

	numberFilter := c.Query("number")
	filteredByNameAndNumber, numberFilterError := filterDummiesByNumber(filteredByName, numberFilter)
	if numberFilterError != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "could not filter by query param 'number'",
		})
		return
	}

	c.JSON(http.StatusOK, takeOf(filteredByNameAndNumber, start, length))
}

func takeOf[T any](slice []T, start int, length int) []T {
	sliceLength := len(slice)
	if start >= sliceLength {
		return []T{}
	}
	maxLength := start + int(math.Min(float64(length), float64(sliceLength-start)))
	return slice[start:maxLength]
}

func filterDummiesByName(dummies []models.Dummy, nameFilter string) ([]models.Dummy, error) {
	if nameFilter == "" {
		return dummies, nil
	}
	pattern, err := regexp.Compile(nameFilter)
	if err != nil {
		return nil, err
	}
	slice := make([]models.Dummy, 0)
	for _, dummy := range dummies {
		if pattern.MatchString(dummy.Name) {
			slice = append(slice, dummy)
		}
	}
	return slice, nil
}

func filterDummiesByNumber(dummies []models.Dummy, numberFilter string) ([]models.Dummy, error) {
	const (
		Equal        int = iota
		GreaterEqual     = iota
		Greater          = iota
		LowerEqual       = iota
		Lower            = iota
	)

	if numberFilter == "" {
		return dummies, nil
	}

	filterType := Equal
	numberText := numberFilter
	if strings.HasPrefix(numberFilter, ">=") {
		filterType = GreaterEqual
		numberText = numberFilter[2:]
	} else if strings.HasPrefix(numberFilter, ">") {
		filterType = Greater
		numberText = numberFilter[1:]
	} else if strings.HasPrefix(numberFilter, "<=") {
		filterType = LowerEqual
		numberText = numberFilter[2:]
	} else if strings.HasPrefix(numberFilter, "<") {
		filterType = Lower
		numberText = numberFilter[1:]
	} else if strings.HasPrefix(numberFilter, "=") {
		filterType = Equal
		numberText = numberFilter[1:]
	}
	numberParsed, numberParseError := strconv.Atoi(numberText)
	if numberParseError != nil {
		return nil, numberParseError
	}
	slice := make([]models.Dummy, 0)
	for _, dummy := range dummies {
		if filterType == Equal && dummy.Number == numberParsed ||
			filterType == GreaterEqual && dummy.Number >= numberParsed ||
			filterType == Greater && dummy.Number > numberParsed ||
			filterType == LowerEqual && dummy.Number <= numberParsed ||
			filterType == Lower && dummy.Number < numberParsed {
			slice = append(slice, dummy)
		}
	}
	return slice, nil
}

func RegisterDummyRoutes(group *gin.RouterGroup) {
	group.GET("", getDummies)
}
