package t

import "fmt"

func PrintHelp() {
	fmt.Printf("The following arguments are mandatory:\n")
	fmt.Printf("  -url               server url\n")
	fmt.Printf("  -form              http form eg: 'a=a&b=b' required if not GET\n\n")
	fmt.Printf("The following arguments are optional:\n")
	fmt.Printf("  -method            http method ['GET']\n")
	fmt.Printf("  -duration          seconds for requst   ['60']\n")
	fmt.Printf("  -qps               requst per second ['100'] note: 10 multiples \n")
	fmt.Printf("  -disKA             DisableKeepAlives, if true, prevents re-use of TCP connections between different HTTP requests ['true']\n")
	fmt.Printf("  -disComp           DisableCompression, if true, prevents the Transport from requesting compression with an Accept-Encoding: gzip ['true']\n")
	fmt.Printf("  -timeout           time for handshake ['0']\n")
	fmt.Printf("  -cpus              number of CPUs ['maximum']\n")
	fmt.Printf("  -bodyShow          show response body ['true']\n")
}

func printResult(resp respon, bodyShow bool) {
	result := fmt.Sprintf("[Response] statusCode-%d | costTime-%3.2fms | contentLength-%d", resp.statusCode, resp.duration, resp.contentLength)
	if bodyShow {
		result = fmt.Sprintf("[Response] statusCode-%d | costTime-%3.2fms | contentLength-%d | body-%v", resp.statusCode, resp.duration, resp.contentLength, resp.body)
	}
	fmt.Println(result)
}

func printErrorResult(resp respon) {
	result := fmt.Sprintf("[Error-Response] statusCode-%d | costTime-%3.2fms | contentLength-%d | err-[%v]", resp.statusCode, resp.duration, resp.contentLength, resp.err)
	fmt.Println(result)
}

func printStatMap(reqStat reqStatus) string {
	return fmt.Sprintf("<50ms-[%v] 50ms~99ms-[%v] 100ms~199ms-[%v] 200ms~299ms-[%v] 300ms~399ms-[%v] 400ms~499ms-[%v] >500ms-[%v]", reqStat.t50, reqStat.t50_99, reqStat.t100_199, reqStat.t200_299, reqStat.t300_399, reqStat.t400_499, reqStat.t500)
}
