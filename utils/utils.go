package utils

import (
	"bufio"
	"cloudcost/global"
	"cloudcost/text"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"slices"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"
)

// generic function to remove duplicated data (named return)
func RemoveDuplicate[T comparable](s *[]T) (result []T) {
	inResult := make(map[T]bool)
	for _, str := range *s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return
}

// case insensitive to file name (named return)
func CaseInsensitive(path *string) (p string) {
	if runtime.GOOS == "windows" {
		return *path
	}

	for _, r := range *path {
		if unicode.IsLetter(r) {
			p += fmt.Sprintf("[%c%c]", unicode.ToLower(r), unicode.ToUpper(r))
		} else {
			p += string(r)
		}
	}
	return
}

// check if file or directory exists or not
func FileExists(f *string) bool {
	_, err := os.Stat(*f)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetNumfilesDir(path string) int {
	files, _ := os.ReadDir(path)
	return len(files)
}

func DebugMemUsage() (uint64, uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	// fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	// fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	// fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	// fmt.Printf("\tNumGC = %v\n", m.NumGC)

	return bToMb(m.Alloc), bToMb(m.Sys)
}
func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func CheckProgressBar(args map[string]string) bool {
	return global.Barenable && args[global.Flagsheader] != "true"
}

func ShowProgressBar(ctx context.Context) {
	// ASCII documentation https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797

	i := 1

	timer := time.Now()
	var color string

	// use ANSI Escape Codes
	wcolor := "\033[31m"
	ncolor := "\033[0m"
	fmt.Print("\x1B[0G")               // moves to left position 0
	fmt.Printf("%v ", global.InnerBar) // show dots bar
	fmt.Print("\x1B[3G")               // moves to left position 12
	fmt.Print("\x1B[7")                // saves cursor position (DEC)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			alloc, talloc := DebugMemUsage()

			if talloc > uint64(global.Memory_alert) {
				color = wcolor
			} else {
				color = ncolor
			}
			if i < global.Barsize {
				fmt.Print(global.Barchar) // write dash
				fmt.Print("\x1B[s")       // save cursor position
				fmt.Print("\x1B[45G")     // moves to right position
				fmt.Printf(" %v", time.Since(timer))

				fmt.Print("\x1B[55G") // moves to right position
				fmt.Printf(" Alloc %vMiB MemSys %v%vMiB %v", alloc, color, talloc, ncolor)

				fmt.Print("\x1B[u") // restore the cursor to the last saved position

			} else {
				fmt.Print("\x1B[0G") // moves to left position 0
				fmt.Printf("%v", global.InnerBar)
				fmt.Print("\x1B[3G")
				i = 0
			}
			time.Sleep(200 * time.Millisecond) // main countdown
		}
		i++
	}
}

func SubDays(d int) int {
	return d - 5 // return day minus 5
}

func ParseDatetoTime(date string) (y int, m time.Month, d int) {

	var dt time.Time

	ok := IsISO8601(date)

	if ok {
		t, _ := time.Parse(time.RFC3339, date) // Parse YYYY-MM-DDTHH:MM:SSZ
		date = t.Format(global.Datelayout)     // rewrite date in MM-DD-YYYY
	}

	date = date[:10] // cut any date to max 10 chars
	date = strings.ReplaceAll(date, "/", "-")

	if date[5:5] == "-" { // yyyy-mm-dd 2024-02-28
		date = date[9:10] + "-" + date[7:8] + "-" + date[1:4] // dd-mm-yyyy
		dt, _ = time.Parse(global.Datelayout, date)           // parse date with default layout
		//date = dt.Format(layout_ptax)
		y = dt.Year()
		m = dt.Month()
		d = dt.Day()

	} else {
		dt, _ = time.Parse(global.Datelayout, date) // parse date with default layout
		//date = dt.Format(layout_ptax)
		y = dt.Year()
		m = dt.Month()
		d = dt.Day()
	}

	return
}

func ParseDatetoString(y int, m time.Month, d int) (ptaxString string) {
	dateString := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	ptaxString = dateString.Format(global.Ptaxlayout)
	return
}

// determine if date is ISO8601 CUR format YYYY-MM-DDTHH:MM:SSZ
func IsISO8601(str string) bool {
	_, err := time.Parse(time.RFC3339, str)
	return err == nil
}

/*
* IMPORTED CODE FROM go src/runtime/float.go
* Changed to be exportable from utils.go
 */
// isInf reports whether f is an infinity.
func IsInf(f float64) bool {
	return !IsNaN(f) && !IsFinite(f)
}

// isNaN reports whether f is an IEEE 754 “not-a-number” value.
func IsNaN(f float64) (is bool) {
	// IEEE 754 says that only NaNs satisfy f != f.
	return f != f
}

// isFinite reports whether f is neither NaN nor an infinity.
func IsFinite(f float64) bool {
	return !IsNaN(f - f)
}

func PrintSize(v interface{}) {
	rv := reflect.ValueOf(v)
	len, cap := sizeof(rv)
	fmt.Printf("%v => len: %d bytes, cap: %d bytes\n", rv.Type(), len, cap)
}

func sizeof(rv reflect.Value) (int, int) {
	rt := rv.Type()

	switch rt.Kind() {
	case reflect.Slice:
		size := int(rt.Size()) // Yields 12, which is the size of the reflect.SliceHeader.
		if rv.Len() > 0 {
			// Recursively get size of the slice elements.
			l, c := sizeof(rv.Index(0))

			// Multiply by total number of elements (len and cap) and add the slice
			// header size.
			return size + (l * rv.Len()), size + (c * rv.Cap())
		}
	}

	return int(rt.Size()), int(rt.Size())
}

func ShowPipeData(dataus, databr map[string]float64, symbl *string, lines int, w *tabwriter.Writer, label string) {
	fmt.Printf("\n%v\n\n", text.Bold(label))
	if len(databr) < lines || !global.ShowPipe {
		fmt.Fprintf(w, global.Product_vHeader)

		for key, value := range databr {
			fmt.Fprintf(w, "\t%v %.2f\t\t%v %.2f\t\t%v\n", global.Msg_SymblUS, dataus[key], *symbl, value, key)
		}
		w.Flush()
	} else {
		pr, pw := io.Pipe() // create pipeReader and pipeWriter
		go func() {
			for key, value := range databr {
				brlsvalue := fmt.Sprintf("%.4f", value)
				usdsvalue := fmt.Sprintf("%.2f", dataus[key])
				pw.Write([]byte("\t" + global.Msg_SymblUS + " " + usdsvalue + "\t" + *symbl + " " + brlsvalue + "\t" + key + "\n"))
			}
			defer pw.Close()
		}()
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			for i := 0; i < lines; i++ {
				if !scanner.Scan() {
					break
				}
				fmt.Println(scanner.Text())
			}
			fmt.Println(global.Msg_prsenter) // show enter msg to continue
			fmt.Scanln()                     // wait for press enter
			fmt.Print("\x1B[2A \x1B[2K")     // move cursor up 2 lines and erase entire line
		}
		_ = pr.Close() // close pipeReader
	}
}

// FIX-ME I NEED A INTERFACE (STR CAN BE SLICE OR STRING)
func SlcContainStr(slc []string, str string) (ok bool) {
	strt := "-" + str
	strtt := "--" + str
	ok = slices.Contains(slc, str) ||
		slices.Contains(slc, strt) ||
		slices.Contains(slc, strtt)

	return
}

// check for arguments into slice/ptr
func SlicesContainArg(flag_ptr interface{}, arg *[]string, srch string) bool {
	ok_srch := "-" + srch   // -argument
	ok__srch := "--" + srch // --argument
	ok_slc := slices.Contains(*arg, ok_srch) || slices.Contains(*arg, ok__srch)
	return flag_ptr == true || ok_slc
}

func ShowExamples(cmdarg *string) {
	fmt.Printf("\n%v\n", global.Msg_exemp)
	for _, example := range global.Examples {
		fmt.Printf("\n\t%v %v", *cmdarg, example)
	}
	fmt.Printf("\n\n%v\n%v\n\n", global.Msg_endex, global.Msg_github)
}
