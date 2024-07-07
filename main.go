package main

import (
	"bufio"
	"cloudcost/global"
	"cloudcost/text"
	"cloudcost/utils"
	"context"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"slices"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"
)

// store billing data
type Billing struct {
	company  []string
	accounts []string
	sdate    []string
	fdate    []string
	products []string
	rscename []string
	currency []string
	totalUSD float64
	totalBRL float64
	platform string
	repcloud string
}

// store csv field position
type Fields struct {
	pos_sdate    int
	pos_fdate    int
	pos_company  int
	pos_accounts int
	pos_products int
	pos_usagtype int
	pos_costusd  int
	pos_costbrl  int
	pos_currency int
	pos_resident int
	pos_repcloud int
}

type Ptax struct {
	CotacaoVenda    interface{}
	DataHoraCotacao interface{}
	Offline         bool // true = offline
}

type ProdCount struct {
	product  string
	cost_usd float64
	cost_brl float64
}

type AccountCount struct {
	account  string
	cost_brl float64
	cost_usd float64
}

type RsceCount struct {
	resource string
	cost_usd float64
	cost_brl float64
}

type RsceIdent struct {
	resource string
	cost_usd float64
	cost_brl float64
}

var accountCount []AccountCount
var prodCount []ProdCount
var rsceCount []RsceCount
var rsceIdent []RsceIdent

var index int64 // count total future usage

var Construct = map[string]map[string]int{}

var barchar = global.Barchar

func main() {
	// resource management
	if global.Max_memlimit != 0 {
		debug.SetMemoryLimit(global.Max_memlimit)
	}
	if runtime.GOOS == "windows" {
		barchar = global.BarcharWin
	}

	// enviroment variables
	var platform, repcloud string
	args := map[string]string{}
	goversion := runtime.Version()
	timeststart := time.Now()

	billing := Billing{} // initiate billing struct
	cotacao := Ptax{}    // initiate ptax struct

	// concurency control
	ch := make(chan Billing) // define main channel
	chp := make(chan Ptax)   // define ptax channel
	var wg sync.WaitGroup    // create concurrency waitgroup

	// slices variables
	var (
		sdate,
		fdate,
		company,
		accounts,
		products,
		cred_name,
		currency []string
	)

	// strings variables
	var (
		path,
		suppname,
		saviname,
		msg_conn,
		date_str,
		date_end,
		fee_sval,
		fee_msg,
		cursymbol,
		ptax_valdate string
	)

	// mathematic variables
	var (
		mpbrltotal,
		mpusdtotal,
		credtotal,
		supptotal,
		usage,
		savitotal float64
	)

	// market place resource sum
	var total_brl, total_usage, feemp_usd, fee_value, ptax_flt float64 // total usage + fee values
	var sumusd, sumbrl, checksum float64                               // main sum values + checksum

	var savings bool

	// columns tabwriter.(TabIndent|StripEscape|AlignRight|DiscardEmptyColumns|Debug)
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.StripEscape)

	// main maps
	prodBrlTotal := make(map[string]float64)
	prodUsdTotal := make(map[string]float64)
	rsceBrlTotal := make(map[string]float64)
	rsceUsdTotal := make(map[string]float64)
	mkplBrlTotal := make(map[string]float64)
	rsceBrlIdent := make(map[string]float64)
	rsceUsdIdent := make(map[string]float64)
	accountBrlTotal := make(map[string]float64)
	accountUsdTotal := make(map[string]float64)

	// set arguments using package flag
	dirFlag := flag.String(global.Flagdirpath, "", global.Msg_flagdir)
	accFlag := flag.String(global.Flagaccount, "", global.Msg_flagacc)
	schFlag := flag.String(global.Flagsearch, "", global.Msg_flagsch)

	verFlag := flag.Bool(global.Flagversion, false, global.Msg_version)
	csvFlag := flag.Bool(global.Flagsheader, false, global.Msg_flaghed)
	rscFlag := flag.Bool(global.Flagusgtype, false, global.Msg_flagrsc)
	rscIdet := flag.Bool(global.Flagresrcid, false, global.Msg_flagrid)
	rscMkpl := flag.Bool(global.Flagmarkplc, false, global.Msg_flagmkl)
	nosPipe := flag.Bool(global.Flagnopipe, false, global.Msg_nopipe)
	rscGrop := flag.Bool(global.Flagresrcgr, false, global.Msg_flagrsg)

	memLimt := flag.Int64(global.Flagmemlimt, 0, global.Msg_memlimt)

	feeCost := flag.Float64(global.Flagfee, 0, global.Msg_flagfee)
	feeMplc := flag.Float64(global.Flagfeemp, 0, global.Msg_flfeemp)

	// TO-DO
	// ptxCost := flag.Float64(global.Flagptax, 0, global.Msg_flagptx)
	// ptxMplc := flag.Float64(global.Flagptaxmp, 0, global.Msg_flptxmp)

	flag.Parse()

	// change resource management from flag
	if *memLimt > 0 && *memLimt == int64(*memLimt) {
		// checking for overflow
		if *memLimt <= math.MaxInt64 {
			debug.SetMemoryLimit(*memLimt) // set custom max memory limit
			global.Memory_alert = *memLimt // set memory alert
		}
	}

	// show version
	if *verFlag {
		fmt.Printf("%v %v %v\n%v %v\n", global.Program, global.Version, global.Codename, global.License, goversion)
		os.Exit(0)
	}

	// create valid path
	if len(*dirFlag) > 0 {
		path = *dirFlag
		err := utils.FileExists(&path)

		if !err {
			fmt.Println(global.Msg_fileerro)
			os.Exit(0)
		}

		if !strings.Contains(path, global.Msg_extcsv) {
			path = *dirFlag + "/*" + global.Msg_extcsv
		}

	} else if len(os.Args) > 1 {
		path = os.Args[1]
		err := utils.FileExists(&os.Args[1])

		if !err {
			fmt.Println(global.Msg_fileerro)
			os.Exit(0)
		}

		if !strings.Contains(os.Args[1], global.Msg_extcsv) {
			path = os.Args[1] + "/*.csv"
		}
	}

	// check is path contains malformat string "//"
	if ok := strings.Contains(path, "//"); ok {
		path = strings.Replace(path, "//", "/", 1)
	}

	// treat path case insensitive
	files, err := filepath.Glob(utils.CaseInsensitive(&path))

	// show usage program if path is wrong
	if err != nil || files == nil {
		fmt.Printf("%v\n\n", global.Msg_fileerro)
		fmt.Printf("%v\n\n", global.Progdesc)
		flag.Usage()
		utils.ShowExamples(&os.Args[0])
		return
	}

	// account arguments #miau
	// FIX-ME need have a interface to string and slices
	if utils.SlcContainStr(os.Args, global.Flagaccount) {
		for i := 0; i < len(os.Args); i++ {
			if strings.Contains(os.Args[i], "--"+global.Flagaccount) {
				*schFlag = os.Args[i+1]
			}
		}
	}

	// check for account flag
	if len(*accFlag) > 0 {
		args[global.Flagaccount] = *accFlag
	}

	// search arguments #miau
	// FIX-ME need have a interface to string and slices
	if utils.SlcContainStr(os.Args, global.Flagsearch) {
		for i := 0; i < len(os.Args); i++ {
			if strings.Contains(os.Args[i], "--"+global.Flagsearch) {
				*schFlag = os.Args[i+1]
			}
		}
	}

	// check for search flag
	if len(*schFlag) > 0 {
		args[global.Flagsearch] = *schFlag
	}

	// show csv file header
	if utils.SlicesContainArg(*csvFlag, &os.Args, global.Flagsheader) {
		args[global.Flagsheader] = "true"
	}

	// show flag by resource
	if utils.SlicesContainArg(*rscFlag, &os.Args, global.Flagusgtype) {
		args[global.Flagusgtype] = "true"
	}

	// show flag by resource id (arn)
	if utils.SlicesContainArg(*rscIdet, &os.Args, global.Flagresrcid) {
		args[global.Flagresrcid] = "true"
	}

	// show flag by market place
	if utils.SlicesContainArg(*rscMkpl, &os.Args, global.Flagmarkplc) {
		args[global.Flagmarkplc] = "true"
	}

	// disable pipe line if nopipe arg is present
	if utils.SlicesContainArg(*nosPipe, &os.Args, global.Flagnopipe) {
		global.ShowPipe = false // used in utils.ShowPipeData()
		global.UseColor = false // used to bold text and colors
	}

	// set flag fee
	if *feeCost != 0 {
		args[global.Flagfee] = fmt.Sprintf("%.4f", *feeCost)
	}
	// set flag fee market place
	if *feeMplc != 0 {
		args[global.Flagfeemp] = fmt.Sprintf("%.4f", *feeMplc)
	}

	// TO-DO
	// set flag ptax
	// if *ptxCost != 0 {
	// 	args[global.Flagptax] = fmt.Sprintf("%.4f", *ptxCost)
	// }
	// // set flag ptax market place
	// if *ptxMplc != 0 {
	// 	args[global.Flagptaxmp] = fmt.Sprintf("%.4f", *ptxMplc)
	// }

	// create context to cancel progress bar showProgressBar
	ctx, cancelWait := context.WithCancel(context.Background())

	// check if progress bar is enabled globally
	progress_bar := utils.CheckProgressBar(args)

	if progress_bar {
		// start progress bar parallel
		go utils.ShowProgressBar(ctx)
	}

	// MAIN GOROUTINE PROCCESS
	for i, filename := range files {
		// limit max files to analysis
		if global.Max_limitfiles > 0 && i >= global.Max_limitfiles {
			defer fmt.Printf("%v (%v)\n", global.Msg_maxreached, global.Max_limitfiles)
			break
		}

		// identify default rune (delimiter) from csv
		runecsv := csvDelimiter(&filename)

		// try discover cloud origin
		cloudcsv, mhead := cloudFound(&filename, &runecsv)

		// set waitgroup counter
		wg.Add(1)

		// load csv file in concurrency mode
		go readCSV(&filename, args, ch, &wg, &runecsv, cloudcsv, mhead)

		// receive channel with struct data from readCSV()
		billing = <-ch

		// sum usd/brl from each file
		sumusd += billing.totalUSD
		sumbrl += billing.totalBRL

		// just store file quantity (not used yet)
		i++
	}

	// load main data from struct channel
	if len(billing.company) > 0 {
		company = utils.RemoveDuplicate(&billing.company)
		sdate = utils.RemoveDuplicate(&billing.sdate)
		fdate = utils.RemoveDuplicate(&billing.fdate)
		accounts = utils.RemoveDuplicate(&billing.accounts)
		products = utils.RemoveDuplicate(&billing.products)
		currency = utils.RemoveDuplicate(&billing.currency)
		platform = billing.platform
		repcloud = billing.repcloud
	}

	// cancel progress bar
	if progress_bar {
		timestop := time.Since(timeststart)

		// create context to show wait bar
		go func() {
			cancelWait()
		}()

		// replace last progress bar from showProgressBar
		fmt.Print("\x1B[0G[ ")
		for i := 0; i < global.Barsize; i++ {
			fmt.Print(barchar)
		}
		fmt.Print(" ] [ ", timestop, " ] \n\n")

	}

	// try to get Ptax from API BC in json with concurrency
	for i, dt_finish := range fdate {

		var ptaxvalue []string

		wg.Add(1)
		go GetPtax(&dt_finish, global.SymbolUSD, chp, &wg)
		cotacao = <-chp

		// take care with this s*
		date_str = sdate[i][:10]

		// cut initial date using length bytes
		date_end = fdate[i][:10] // cut end date using length bytes

		// proccess ptax response in Online mode (cotacao.Offline == false)
		if !cotacao.Offline {
			// control ptax custom values
			msg_conn = global.Msg_connon
			ptaxstr := fmt.Sprintf("%v %v", cotacao.CotacaoVenda, cotacao.DataHoraCotacao)
			ptaxvalue = append(ptaxvalue, ptaxstr)
			ptax_valdate = ptaxvalue[i][:18] // cut ptax + date using length bytes

			// convert ptax string to float64
			ptax_flt, _ = strconv.ParseFloat(ptax_valdate[:6], 64)

			// try to emulate fee brl
			total_usage = ptax_flt * sumusd // multiply ptax times usd total

		}
	}

	if !slices.Contains(currency, global.SymbolUSD) {
		cursymbol = global.Msg_SymblBR
	} else {
		cursymbol = global.Msg_SymblUS
	}

	// show account not found and exit
	if len(args[global.Flagaccount]) > 0 && len(accounts) == 0 {
		fmt.Printf("%v (%v)\n", global.Msg_anotfoud, args[global.Flagaccount])
		os.Exit(0)
	}

	// show search msg not found and exit
	if len(products) < 1 {
		fmt.Printf("%v (%v)\n", global.Msg_snotfoud, args[global.Flagsearch])
		os.Exit(0)
	}

	// show company information
	for _, str := range company {
		fmt.Fprintf(w, "%v\t%v\t\n\n", global.Msg_contrt, str)
	}

	// sum values by product
	for k, value := range prodCount {
		prodBrlTotal[prodCount[k].product] += value.cost_brl
		prodUsdTotal[prodCount[k].product] += value.cost_usd
	}

	// sum values by accounts
	for k, value := range accountCount {
		accountBrlTotal[accountCount[k].account] += value.cost_brl
		accountUsdTotal[accountCount[k].account] += value.cost_usd
	}

	// load resource cost values into map
	for k, value := range rsceCount {
		rsceUsdTotal[rsceCount[k].resource] += value.cost_usd
		rsceBrlTotal[rsceCount[k].resource] += value.cost_brl
	}

	// load resource ident values into map
	for k, value := range rsceIdent {
		rsceUsdIdent[rsceIdent[k].resource] += value.cost_usd
		rsceBrlIdent[rsceIdent[k].resource] += value.cost_brl
	}

	// prints platform
	fmt.Fprintf(w, "%v\t%v-%v\t\n\n", global.Msg_source, platform, repcloud)

	// prints account
	fmt.Printf("%v\n\n", global.Msg_accout)
	fmt.Fprintf(w, global.Account_Header)
	for key, value := range accountBrlTotal {
		if repcloud == global.RepAzure {
			value = value / ptax_flt
		}
		vtotal := accountUsdTotal[key] // get usd value for each account
		fmt.Fprintf(w, "\t%v\t %v %.4f\t\t%v %.4f\t\n", key, global.Msg_SymblUS, value, cursymbol, vtotal)
	}
	w.Flush()

	fmt.Printf("\n%v\n\n", global.Msg_produt)
	fmt.Fprintf(w, global.Product_Header)

	// prints products with max char control
	for pkey, value := range prodBrlTotal {

		if len(pkey) > global.Max_charprint {
			pkey = pkey[:global.Max_charprint]
		}

		vtotal := prodUsdTotal[pkey] // get usd value for each resource

		if repcloud == global.RepAzure {
			vtotal = vtotal / ptax_flt
		}
		fmt.Fprintf(w, "\t%v\t %v %.4f\t\t%v %.4f\t\n", pkey, global.Msg_SymblUS, vtotal, cursymbol, value)

		checksum += value
		// run slice to check string (credprefix) into string (key)
		for _, credprefix := range global.Credits {
			if strings.Contains(pkey, credprefix) {
				credtotal += value
				cred_name = append(cred_name, pkey)
			}
		}
		// get support total from product name
		for _, sprefix := range global.Support {
			if strings.Contains(pkey, sprefix) {
				supptotal += value
				suppname = pkey
			}
		}
	}

	w.Flush()

	// run UsageType slice to get total foreach resource
	for key, value := range rsceBrlTotal {

		var mpblankvalue float64

		if key == "" {
			mpblankvalue = value
		}
		// get market place total from usage type (BRL)
		for _, mprefix := range global.Marketplace {
			if strings.Contains(key, mprefix) {
				mpbrltotal += value
				mkplBrlTotal[key] = value
			}
		}
		// add market place from blank value
		mpbrltotal += mpblankvalue
	}

	// get market place total from usage type (USD)
	for key, value := range rsceUsdTotal {
		var mpblankvalue float64
		if key == "" {
			mpblankvalue = value
		}
		for _, mprefix := range global.Marketplace {
			if strings.Contains(key, mprefix) {
				mpusdtotal += value
			}
		}
		mpusdtotal += mpblankvalue
	}

	// looks for Savings Plans
	for key, value := range prodBrlTotal {
		for _, sprefix := range global.Savings {
			if !savings {
				// run once
				if strings.Contains(key, sprefix) {
					saviname = key
					savitotal += value
					savings = true
				}
			}

		}
	}

	// SHOW USAGETYPE message with flag
	if args[global.Flagusgtype] == "true" {
		utils.ShowPipeData(rsceUsdTotal, rsceBrlTotal, &cursymbol, global.PipeLine, w, global.Msg_rsrctype)
	}

	// SHOW RESOURCEID with pipe control using io.Pipe()
	if args[global.Flagresrcid] == "true" {
		utils.ShowPipeData(rsceUsdIdent, rsceBrlIdent, &cursymbol, global.PipeLine, w, global.Msg_resrcid)
	}

	if *rscGrop {
		// load only for Azure reports
		if repcloud == global.RepAzure {
			// Create aux map to store sum
			sumById := map[string]float64{}
			var strId string
			var isKey []string
			// Running map "rsceBrlIdent"
			for id, value := range rsceBrlIdent {
				isKey = strings.Split(id, "/")
				if len(isKey) > 4 {
					strId = isKey[4]
				}

				// Check if ID exists in aux map
				_, exist := sumById[strId]

				// If ID not exists, add to aux map with original value
				if !exist {
					sumById[strId] = value
				} else {
					// If Id exists, sum value into
					sumById[strId] += value
				}
			}

			// Print map "somaPorId"
			fmt.Printf("\n%v\n\n", text.Bold(global.Msg_resrgrp))
			for id, sum := range sumById {
				fmt.Fprintf(w, "\t%v\t%v %.2f\n", id, global.Msg_SymblBR, sum)
			}
			w.Flush()
		}
	}

	// summarize total usage
	usage = sumbrl - mpbrltotal - supptotal

	// TO-DO
	// control ptax custom values
	// if *ptxCost > 0 {
	// 	ptax_flt = *ptxCost
	// 	ptax_valdate = fmt.Sprintf("%.4f", ptax_flt)
	// 	total_usage = ptax_flt * sumusd
	// }

	// validate connection
	if cotacao.Offline {
		msg_conn = global.Msg_connoff
		mpusdtotal = mpbrltotal // traversal
	}

	// control fee custom values
	switch {
	case *feeCost > 0:
		// define fee from flag
		fee_value = *feeCost
	case platform != global.Cmp:
		// define default fee from global (1)
		fee_value = global.DefaultFee
	case mpusdtotal > 0:
		// reverse enginneer when fee marketplace != fee usage
		total_us := sumusd - mpusdtotal  // sum total usd - mp total us
		total_br := total_us * ptax_flt  // total us * ptax
		total_wmp := sumbrl - mpbrltotal // sum total brl - mp total br
		fee_value = total_wmp / total_br // total wp / total br
	default:
		// define fee value from report
		fee_value = (sumbrl / total_usage)
	}

	// control feemp (market place) custom values
	switch {
	case *feeMplc > 0:
		feemp_usd = *feeMplc
	case platform != global.Cmp:
		feemp_usd = global.DefaultFee
	default:
		mpbrl := mpusdtotal * ptax_flt
		feemp_usd = mpbrltotal / mpbrl
	}

	// miau - treat marketplace total cost
	if slices.Contains(currency, global.SymbolUSD) {
		mpusdtotal = mpbrltotal
	}

	// manipule market place total cost with feemp
	mpbrltotal = mpusdtotal * ptax_flt * feemp_usd

	// check if fee_value has error or is +Inf (infinity) like math.IsInf()
	if fee_value < 0 || utils.IsInf(fee_value) {
		fee_sval = "n/a" // do not show negative fee numbers (like error)
	} else {
		fee_sval = fmt.Sprintf("%.4f", fee_value)
	}

	// SHOW DEFAULT USAGE
	fmt.Printf("\n%v\n\n", global.Msg_usage)
	fmt.Fprintf(w, "\t%v\n", global.Resource_Header)
	if repcloud == global.RepAzure {
		sumusd = sumusd / ptax_flt
	}
	fmt.Fprintf(w, "\t%v\t\t%v\t\t%v\t\t%v\t\t%v %.2f\t\t%v %.2f\n", date_str, date_end, ptax_valdate, fee_sval, global.Msg_SymblUS, sumusd, cursymbol, usage)
	w.Flush()

	// SHOW MARKET PLACE
	if mpbrltotal > 0 && mpbrltotal == float64(mpbrltotal) { // !math.IsNaN(mpbrltotal) {
		fmt.Printf("\n%v\n\n", global.Msg_mktpl)
		fmt.Fprintf(w, "\t%v\n", global.MarketP_Header)
		fmt.Fprintf(w, "\t%v\t\t%.4f\t\t%v %.2f\t\t%v %.2f\n", ptax_flt, feemp_usd, global.Msg_SymblUS, mpusdtotal, global.Msg_SymblBR, mpbrltotal)
		w.Flush()
	}

	// SHOW MARKET PLACE FROM FLAG
	if args[global.Flagmarkplc] == "true" {
		// FIX-ME Need show prodname instead resourcetype
		fmt.Printf("\n%v\n\n", global.Msg_mktpl)
		fmt.Fprintf(w, "\t%v\n", global.Detail_Header)
		for key, value := range mkplBrlTotal {
			fmt.Fprintf(w, "\t%v %.2f\t%v\n", cursymbol, value, key)
		}
		w.Flush()
	}

	// SHOW CREDITS
	if credtotal != 0 && credtotal == float64(credtotal) {
		// FIX-ME I need be itereted not only once!!!
		fmt.Printf("\n%v\n\n", global.Msg_credit)
		fmt.Fprintf(w, "\t%v\n", global.Detail_Header)
		fmt.Fprintf(w, "\t%v %.2f\t%v\n", cursymbol, credtotal, cred_name)

	}
	w.Flush()

	// SHOW SUPPORT
	if supptotal != 0 && supptotal == (supptotal) {
		fmt.Printf("\n%v\n\n", global.Msg_support)
		fmt.Fprintf(w, "\t%v\n", global.Detail_Header)
		fmt.Fprintf(w, "\t%v %.2f\t%v\n", cursymbol, supptotal, suppname)
	}
	w.Flush()

	// SHOW SAVINGS
	if savings {
		fmt.Printf("\n%v\n\n", global.Msg_savings)
		fmt.Fprintf(w, "\t%v\n", global.Detail_Header)
		fmt.Fprintf(w, "\t%v %.2f\t%v\n", cursymbol, savitotal, saviname)
	}

	w.Flush()

	// SHOW TOTAL
	fmt.Printf("\n%v\n", global.Msg_total)

	w.Flush()

	// prepare custom values
	if slices.Contains(currency, global.SymbolUSD) {
		if *feeCost > 0 {
			total_brl = (sumusd * ptax_flt) * (*feeCost)
			fee_msg = global.Msg_without // msg with fee added
		} else {
			total_brl = sumusd * ptax_flt
			fee_msg = global.Msg_withfee // msg without fee
		}
	} else {
		if *feeCost > 0 {
			total_brl = (sumusd * ptax_flt) * (*feeCost)
			fee_msg = global.Msg_without // msg with fee added
		} else {
			total_brl = sumbrl
		}
	}

	// show total values
	fmt.Fprintf(w, "\n\t%v\n", global.Total_Header)
	fmt.Fprintf(w, "\t%v %.2f\t\t%v %.2f\n", global.Msg_SymblUS, sumusd, global.Msg_SymblBR, total_brl)

	w.Flush()

	// show validation checksum & connectivity
	str1 := fmt.Sprintf("%.2f", checksum)
	str2 := fmt.Sprintf("%.2f", sumbrl)

	if str1 == str2 {
		fmt.Printf("\n%v %v\t%v\t%v\n\n", global.Msg_valid, text.Green(global.Msg_ok), msg_conn, text.Magenta(fee_msg))
	} else {
		fmt.Printf("\n%v %v %v\t%v\t%v\n", global.Msg_valid, text.Red(global.Msg_error), text.Red(global.Msg_vdiff), msg_conn, text.Magenta(fee_msg))
		fmt.Printf("%v (%v) != %v (%v)\n", global.Msg_sumtotal, str1, global.Msg_sumprodt, str2)
	}
	// Aguardar todas as goroutines terminarem
	wg.Wait()
}

// function to read csv files with goroutines waitgroup
func readCSV(filename *string, args map[string]string, ch chan Billing, wg *sync.WaitGroup, runecsv *rune, cloudcsv *string, mhead *map[string]int) {
	defer wg.Done()

	data := &Billing{}
	fp := Fields{} // csv field position
	var count int64

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', tabwriter.TabIndent) // used flagsheader

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("%v %s: %v\n", global.Msg_fileopen, *filename, err)
	}
	defer file.Close()

	reader := csv.NewReader(file) // read csv
	reader.Comma = *runecsv       // define default rune (delimiter)

	// extract field position from csv main header (*mhead)
	for k, v := range *mhead {
		switch {
		case slices.Contains(global.StartDate, k):
			fp.pos_sdate = v // start date position
		case slices.Contains(global.EndDate, k):
			fp.pos_fdate = v // finish date position
		case slices.Contains(global.CompanyName, k):
			fp.pos_company = v // company position
		case slices.Contains(global.UsageAccount, k):
			fp.pos_accounts = v // accounts position
		case slices.Contains(global.ProductName, k):
			fp.pos_products = v // products position
		case slices.Contains(global.UsageType, k):
			fp.pos_usagtype = v // usagetype position
		case slices.Contains(global.ResourceCost, k):
			if *cloudcsv == global.Cmp {
				fp.pos_costusd = v // cost usd position
			}
		case slices.Contains(global.FinalCost, k):
			if *cloudcsv != global.Cmp {
				fp.pos_costusd = v // cost usd position
			}
			fp.pos_costbrl = v // cost brl position
		case slices.Contains(global.ReportCloud, k):
			fp.pos_repcloud = v // Report Cloud
		case slices.Contains(global.CurrencyCode, k):
			fp.pos_currency = v // currency position
		case slices.Contains(global.ResourceIdent, k):
			fp.pos_resident = v // resource id position
		}

	}

	for {

		rln, err := reader.Read() // slice rln read line from file

		if err == io.EOF {
			break
		}
		// show error to read file
		if err != nil {
			log.Printf("%v %v: %v\n", global.Msg_fileread, *filename, err)
			os.Exit(0)
		}

		// show csv header file in the first line of file only and os.Exit
		if count < 1 {
			if args[global.Flagsheader] == "true" {
				fmt.Printf("%v %v\n\n", global.Msg_file, *filename)
				for k := range rln {
					fmt.Fprintf(w, "\t%v\t%v\n", rln[k], k)
				}
				w.Flush()  // release Fprint content
				os.Exit(0) // explicit terminate
			}
		}

		// main load engine ignore first row (count (line) > 0)
		if count != 0 {
			switch {
			case args[global.Flagaccount] != "":
				// FIX-ME when run --account I load data with all data fallthrough it's don´t help
				argaccount := strings.Split(args[global.Flagaccount], ",")
				for _, arg := range argaccount {
					if arg == rln[fp.pos_accounts] {
						data = appendData(data, rln, cloudcsv, fp)
					}
				}
				fallthrough // FIX-ME I don´t help (like else if)
			case len(args[global.Flagsearch]) > 2:
				SearchContent(data, args, rln, fp, cloudcsv)
			default:
				data = appendData(data, rln, cloudcsv, fp)
			}
		}

		index++
		count++
	}

	// return struct to channel
	ch <- *data
}

// append data to struct and return pointer
func appendData(data *Billing, rln []string, cloudcsv *string, fp Fields) *Billing {

	data.platform = *cloudcsv
	data.repcloud = rln[fp.pos_repcloud] // do not use append

	data.sdate = append(data.sdate, rln[fp.pos_sdate])

	data.fdate = append(data.fdate, rln[fp.pos_fdate])
	data.company = append(data.company, rln[fp.pos_company])
	data.accounts = append(data.accounts, rln[fp.pos_accounts])
	data.products = append(data.products, rln[fp.pos_products])

	data.rscename = append(data.rscename, rln[fp.pos_usagtype])
	data.currency = append(data.currency, rln[fp.pos_currency])

	susd := strings.Replace(rln[fp.pos_costusd], ",", ".", 1)
	sbrl := strings.Replace(rln[fp.pos_costbrl], ",", ".", 1)

	fusd, _ := strconv.ParseFloat(susd, 64)
	fbrl, _ := strconv.ParseFloat(sbrl, 64)

	data.totalUSD += fusd
	data.totalBRL += fbrl

	prodCount = append(prodCount, ProdCount{rln[fp.pos_products], fusd, fbrl})
	rsceCount = append(rsceCount, RsceCount{rln[fp.pos_usagtype], fusd, fbrl})
	rsceIdent = append(rsceIdent, RsceIdent{rln[fp.pos_resident], fusd, fbrl})
	accountCount = append(accountCount, AccountCount{rln[fp.pos_accounts], fusd, fbrl})

	return data
}

// check for default csv delimiter (rune - named return)
// https://www.ietf.org/rfc/rfc4180.txt
func csvDelimiter(filename *string) (runecsv rune) {

	var headercsv string

	var err error

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	fl := bufio.NewReader(file)
	// read the first line only
	for i := 0; i < 1; i++ {
		headercsv, err = fl.ReadString('\n')
		if err != nil {
			fmt.Println(global.Msg_fileread, *filename)
		}
	}
	file.Close() // immediate close no defer!

	switch {
	case strings.Contains(headercsv, ";"):
		runecsv = ';'
	case strings.Contains(headercsv, "|"):
		runecsv = '|'
	case strings.Contains(headercsv, "\t"):
		runecsv = '\t'
	default:
		runecsv = ','
	}
	return
}

func cloudFound(filename *string, runecsv *rune) (*string, *map[string]int) {
	var cloud string
	maphead := map[string]int{}

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = *runecsv

	for i := 0; i < 1; i++ {
		r, err := reader.Read()
		if err == io.EOF {
			break
		}

		switch {
		case slices.Contains(r, global.Src_cexp):
			cloud = global.Aws_cexp
		case slices.Contains(r, global.Src_cur):
			cloud = global.Aws_cur
		case slices.Contains(r, global.Src_azur):
			cloud = global.Azure
		case slices.Contains(r, global.Src_gcp):
			cloud = global.Google
		case slices.Contains(r, global.Src_Cmp):
			cloud = global.Cmp
		default:
			fmt.Println(global.Msg_fileform)
		}
		for k, v := range r {
			maphead[v] = k
		}
	}
	return &cloud, &maphead
}

// try to get dolar ptax from bc with parallelism using channel chp
// BC API https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/aplicacao#!/recursos
func GetPtax(dt_parse *string, symbol string, chp chan Ptax, wg *sync.WaitGroup) {
	defer wg.Done()

	ptax := Ptax{}
	ptax.Offline = false

	var cotacao map[string]interface{}

	y, m, d := utils.ParseDatetoTime(*dt_parse)
	sday := utils.SubDays(d)
	minday := utils.ParseDatetoString(y, m, sday)
	maxday := utils.ParseDatetoString(y, m, d)

	url := fmt.Sprintf(global.PtaxAPIUrl, symbol, minday, maxday)

	// Http request with timeout control default max 3 seconds
	client := http.Client{
		Timeout: global.PtaxTimeout * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		// turn off line mode on
		ptax.Offline = true
		chp <- ptax
		return // finish statement with ptax.Offline on channel
	}
	defer resp.Body.Close()

	// Lendo a resposta
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		defer fmt.Printf("%v %v\n", global.PtaxRespError, err)
	}

	jsonData := []byte(string(responseData))
	err = json.Unmarshal(jsonData, &cotacao)
	if err != nil {
		ptax.Offline = true
		chp <- ptax
		return // finish statement with ptax.Offline on channel
	}

	// Ignoring error test to get cotacaoVenda and data only
	value := cotacao["value"]
	valueSlice := value.([]interface{})
	fistElement := valueSlice[0]
	cotacaoMap := fistElement.(map[string]interface{})

	// Error test only in the last extract
	dataHoraCotacao, ok := cotacaoMap[global.PtaxIndexMapD] // index dataHoraCoracao
	if !ok {
		defer fmt.Println(global.PtaxDataError)
	}
	cotacaoVenda, ok := cotacaoMap[global.PtaxIndexMapV] // index cotacaoVenda
	if !ok {
		defer fmt.Println(global.PtaxPtaxError)
	}
	ptax.CotacaoVenda = cotacaoVenda
	ptax.DataHoraCotacao = dataHoraCotacao

	chp <- ptax
}

func SearchContent(data *Billing, args map[string]string, rln []string, fp Fields, cloudcsv *string) {
	argsearch := strings.Split(args[global.Flagsearch], ",")
	for _, arg := range argsearch {
		switch {
		case strings.Contains(rln[fp.pos_accounts], arg):
			data = appendData(data, rln, cloudcsv, fp)
			//fallthrough
		case strings.Contains(rln[fp.pos_products], arg):
			data = appendData(data, rln, cloudcsv, fp)
			//fallthrough
		case strings.Contains(rln[fp.pos_usagtype], arg):
			data = appendData(data, rln, cloudcsv, fp)
			//fallthrough
		case strings.Contains(rln[fp.pos_resident], arg):
			data = appendData(data, rln, cloudcsv, fp)
		}
	}
}
