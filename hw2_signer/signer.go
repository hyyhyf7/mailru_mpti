package main

// сюда писать код

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

var mx = &sync.Mutex{}
var begin = time.Now()

func ExecutePipeline(jobs ...job) {
	outputChans := make([]chan interface{}, len(jobs)+1)

	for index, _ := range outputChans {
		outputChans[index] = make(chan interface{}, 1)
	}
	wg := &sync.WaitGroup{}

	defer fmt.Println("after pipeline")
	defer wg.Wait()

	for index, worker := range jobs {
		NewLayerWorker := func(data chan interface{}, output chan interface{}, currWorker job) {
			currWorker(data, output)
			close(output)
			wg.Done()
		}

		wg.Add(1)
		go NewLayerWorker(outputChans[index], outputChans[index+1], worker)
	}

}

func SingleHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}

	defer wg.Wait()

	for value := range in {
		wg.Add(1)
		data := ReadString(value)
		go func() {
			//fmt.Println("Time begin", time.Since(begin))
			crc32wg := &sync.WaitGroup{}

			crc32 := make(chan string, 1)
			crc32 <- data

			crc32wg.Add(1)
			go NewLayerCrc32(crc32, crc32, crc32wg)

			//md5begin := time.Now()
			mx.Lock()
			md5 := DataSignerMd5(data)
			mx.Unlock()
			//fmt.Println("time for md5", time.Since(md5begin), "time since", time.Since(begin))

			crc32md5 := make(chan string, 1)
			crc32md5 <- md5

			//fmt.Println("time since before latter", time.Since(begin))
			crc32wg.Add(1)
			go NewLayerCrc32(crc32md5, crc32md5, crc32wg)

			crc32wg.Wait()
			fmt.Println("time since after waiting for crc32", time.Since(begin))

			out <- (<-crc32) + "~" + (<-crc32md5)
			fmt.Println("time after crc32 waiting for channels to be read", time.Since(begin))

			wg.Done()
			//fmt.Println("time for SingleHash", time.Since(begin))
		}()
	}
}

func MultiHash(in chan interface{}, out chan interface{}) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	for value := range in {
		wg.Add(1)
		data := ReadString(value)

		go func() {
			fmt.Println("data is received", time.Since(begin))
			th := [6]string{"0", "1", "2", "3", "4", "5"}

			singles := make([]chan string, 6)
			for index := range th {
				singles[index] = make(chan string, 1)
			}

			threadsWg := &sync.WaitGroup{}

			for index := range th {
				threadsWg.Add(1)
				thd := make(chan string, 1)
				thd <- th[index] + data
				go NewLayerCrc32(thd, singles[index], threadsWg)
			}

			threadsWg.Wait()

			res := ""
			for _, single := range singles {
				res += <-single
			}
			fmt.Println("COMBINED", res)
			out <- res

			wg.Done()
		}()

	}
}

func CombineResults(in chan interface{}, out chan interface{}) {
	res := ""

	var slc []string

	for value := range in {
		data := ReadString(value)

		slc = append(slc, data)
	}

	sort.Strings(slc)

	for index, value := range slc {
		res += value
		if index < len(slc)-1 {
			res += "_"
		}
	}

	out <- res
}

func ReadString(dataRaw interface{}) string {
	data, ok := dataRaw.(string)
	if !ok {
		data += strconv.Itoa(dataRaw.(int))
	}

	return data
}

func NewLayerCrc32(in chan string, out chan string, wg *sync.WaitGroup) {
	out <- DataSignerCrc32(<-in)

	wg.Done()
}
