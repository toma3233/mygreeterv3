function coverage() {
    #Create cover report, convert to html and save to staging directory
    ginkgo -r -v --trace --coverprofile=.coverage-report.out --skip-package=mock ./...
    go tool cover -html=.coverage-report.out -o $1/$2.html
    #Extract coverage percentage as an integer and compare to required threshold.
    results=$(go tool cover -func=.coverage-report.out | grep total: | awk '{print $NF}')
    number="${results%\%}"
    number=$(printf "%.0f" $number)
    if [ $number -lt $threshold ]
    then
        echo -e "${RED}$2 results were: $results and below the required threshold of: $threshold${NC}"
        return 0
    fi
    echo -e "${GREEN}$2 results were: $results and above the required threshold of: $threshold${NC}"
    return 1
}
