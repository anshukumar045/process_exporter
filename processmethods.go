/*********************************************************
/proc/[pid]/stat
/proc/[pid]/smaps
/proc/[pid]/io
/proc/[pid]/status
/proc/sys/fs/file-nr
/proc/meminfo
**********************************************************/

package collector

import (
        "strings"
        "strconv"
        "os"
        "log"
        "os/exec"
        "bufio"
        // "fmt"
        "regexp"
)

type CPUstat struct {
        CPU_total   int
        User       int
        System     int
        Nice       int
        Iowait     int
        Guest      int
        Minpgflt   int
        Majpgflt   int
        ThreadCnt  int
}

type smapsCounter struct{
        rss     int
        pss     int
        sdirty  int
        pdirty  int
}

type IOCounters struct {
        rdcnt      int
        wrtcnt     int
        rdbytes    int

        wrtbytes   int
}

type StatusInfo struct {
        VmData int
        VmStk  int
        VmLck  int
        VmLib  int
        VmExe  int
}
type FileStat struct {
        OpenFiles  int
        TotalFiles int
}
type MemStat struct {
        MemTotal       int
        MemFree        int
        MemAvailable   int
        MemUsed        int
}

func getMemInfo()(*MemStat, error) {
        path := "/proc/meminfo"
        file, err := os.Open(path)
        if err != nil {
                log.Println(err)
        }
        defer file.Close()
        memtotal, memfree, memavailable := 0,0,0
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                if strings.HasPrefix(scanner.Text(), "MemTotal") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        memtotal_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        memtotal = memtotal_
                }
                if strings.HasPrefix(scanner.Text(), "MemFree") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        memfree_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        memfree = memfree_
                }
                if strings.HasPrefix(scanner.Text(), "MemAvailable") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        memavailable_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        memavailable = memavailable_
                }
        }
return &MemStat{MemTotal: memtotal, MemFree: memfree, MemAvailable: memavailable, MemUsed: memtotal-memfree},nil
}

func getFileStat()(*FileStat, error) {
        path := "/proc/sys/fs/file-nr"
        file, err := os.Open(path)
        if err != nil {
                log.Println(err)
        }
        defer file.Close()
        fileOpen, fileTotal := 0,0
        var fdesc string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                strtst := removewhitespaces(scanner.Text())
                for _ ,v := range strtst {
                _, err:=  strconv.Atoi(string(v))
                        if err == nil {
                         fdesc = fdesc + string(v)
                        } else {
                                fdesc = fdesc + " "
                        }
                }
                mystr := strings.Split(fdesc, " ")

                fileopn,err := strconv.Atoi(mystr[0])  // open files
                        if err != nil {
                                return nil, err
                        }
                fileOpen = fileopn
                filetotal,err := strconv.Atoi(mystr[2])   //  total files
                        if err != nil {
                                return nil, err
                        }
                fileTotal = filetotal
        }

return &FileStat{OpenFiles: fileOpen, TotalFiles: fileTotal}, nil
}
func getStatusInfo(pid int)(*StatusInfo, error) {
        path := "/proc/" + strconv.Itoa(pid) + "/status"
        file , err := os.Open(path)
        if err != nil {
                log.Println(err)
        }
        defer file.Close()
        vmlck,vmdata,vmstk,vmexe,vmlib := 0,0,0,0,0
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                if strings.HasPrefix(scanner.Text(), "VmLck") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        vmlck_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        vmlck = vmlck_
                }
                if strings.HasPrefix(scanner.Text(), "VmData") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        vmdata_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        vmdata = vmdata_
                }
                if strings.HasPrefix(scanner.Text(), "VmStk") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        vmstk_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        vmstk = vmstk_
                }
                if strings.HasPrefix(scanner.Text(), "VmExe") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        vmexe_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                         return nil, err
                                }
                        vmexe = vmexe_
                }
                if strings.HasPrefix(scanner.Text(), "VmLib") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        vmlib_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        vmlib = vmlib_
                }


        }

return &StatusInfo{VmData: vmdata, VmStk: vmstk, VmLck: vmlck, VmLib: vmlib, VmExe: vmexe}, nil
}


func removewhitespaces(input string) string {
        re_leadclose_whtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
        re_inside_whtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
        final := re_leadclose_whtsp.ReplaceAllString(input, "")
        final = re_inside_whtsp.ReplaceAllString(final, " ")
        return final
}

func getCommandline(pid int)(string, error){
        path := "/proc/"+strconv.Itoa(pid)+"/cmdline"
        cln, err :=  exec.Command("cat", path).Output()
                if err != nil {
                      return " ", err
                  }
return string(cln), err
}


func getIOCounter(pid int)(*IOCounters, error){
        path := "/proc/" + strconv.Itoa(pid) + "/io"
        file , err := os.Open(path)
        if err != nil {
                log.Println(err)
        }
        defer file.Close()
        rc, wc, rb, wb := 0,0,0,0
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                if strings.HasPrefix(scanner.Text(), "rchar") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        rc_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        rc = rc_
                }
                if strings.HasPrefix(scanner.Text(), "wchar") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        wc_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        wc = wc_
                }
                if strings.HasPrefix(scanner.Text(), "read_bytes") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        rb_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        rb = rb_
                }
                if strings.HasPrefix(scanner.Text(), "write_bytes") {
                        tmp := strings.Split(removewhitespaces(scanner.Text()), " ")
                        wb_, err := strconv.Atoi(tmp[1])
                                if err != nil {
                                        return nil, err
                                }
                        wb = wb_
                }
        }
return &IOCounters{rdcnt: rc, wrtcnt: wc, rdbytes: rb, wrtbytes: wb}, nil
}

func getSmaps(pid int)(*smapsCounter,error) {
        path := "/proc/" + strconv.Itoa(pid) + "/smaps"
        file , err := os.Open(path)
        if err != nil {
                log.Println(err)
        }
        defer file.Close()
        sumRss, sumPss, sDirty, pDirty := 0,0,0,0
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                if strings.HasPrefix(scanner.Text(), "Rss") {
                        tst := strings.Split(removewhitespaces(scanner.Text()), " ")
                        rss , err := strconv.Atoi(tst[1])
                        if err != nil {
                                return nil, err
                        }
                        sumRss = sumRss + rss

                }
                 if strings.HasPrefix(scanner.Text(), "Pss") {
                        tst := strings.Split(removewhitespaces(scanner.Text()), " ")
                        pss , err := strconv.Atoi(tst[1])
                        if err != nil {
                                return nil, err
                        }
                        sumPss = sumPss + pss

                }
                if strings.HasPrefix(scanner.Text(), "Shared_Dirty") {
                        tst := strings.Split(removewhitespaces(scanner.Text()), " ")
                        sd , err := strconv.Atoi(tst[1])
                        if err != nil {
                                return nil, err
                        }
                        sDirty = sDirty + sd

                }
                if strings.HasPrefix(scanner.Text(), "Private_Dirty") {
                        tst := strings.Split(removewhitespaces(scanner.Text()), " ")
                        pd , err := strconv.Atoi(tst[1])

                        if err != nil {
                                return nil, err
                        }
                        pDirty = pDirty + pd

                }



        }
        return &smapsCounter{rss: sumRss, pss: sumPss, sdirty: sDirty, pdirty: pDirty},nil
}

// get CPU Info user system  nice iowait guest
func CPUinfo(pid int)(*CPUstat, error) {

        // Read the value of /proc/[pid]/stat
        path := "/proc/" + strconv.Itoa(pid) + "/stat"
        file , err := os.Open(path)
                if err != nil {
                        log.Println(err)
        }
        defer file.Close()

        cpu_usage, sys, user, nice, iowait_,guest_,minflt, majflt, thrdcnt := 0,0,0,0,0,0,0,0,0
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
                stat := strings.Split(scanner.Text(), " ")
                utime , err := strconv.Atoi(stat[13])
                        if err != nil {
                                return nil, err
                        }
                stime, err := strconv.Atoi(stat[14])
                        if err != nil {
                                return nil,err
                        }
                cutime , err := strconv.Atoi(stat[15])
                        if err != nil {
                                return nil, err
                        }
                cstime, err := strconv.Atoi(stat[16])
                        if err != nil {
                                return nil,err
                        }
                nice_ , err := strconv.Atoi(stat[18])
                        if err != nil {
                                return nil,err
                        }
                iowait__, err := strconv.Atoi(stat[41])
                         if err != nil {
                                return nil,err
                        }
                gtime, err := strconv.Atoi(stat[42])
                        if err != nil {
                                return nil,err
                        }
                cgtime, err := strconv.Atoi(stat[43])
                        if err != nil {
                                return nil,err
                        }
                 minflt_, err := strconv.Atoi(stat[9])
                        if err != nil {
                                return nil,err
                        }
                cminflt, err := strconv.Atoi(stat[10])
                        if err != nil {
                                return nil,err
                        }
                majflt_, err := strconv.Atoi(stat[12])
                        if err != nil {
                                return nil,err
                        }
                cmajflt, err := strconv.Atoi(stat[13])
                        if err != nil {
                                return nil,err
                        }
                thrdcnt_, err := strconv.Atoi(stat[19])
                        if err != nil {
                                return nil,err
                        }

                cpu_usage = (utime + stime + cutime + cstime)
                sys = stime + cstime
                user = utime + cutime
                nice = nice_
                iowait_ = iowait__
                guest_ = gtime + cgtime
                minflt = minflt_ + cminflt
                majflt = majflt_ + cmajflt
                thrdcnt = thrdcnt_
        }

return &CPUstat{CPU_total: cpu_usage, User: user, System: sys, Nice: nice, Iowait: iowait_,
                Guest: guest_,Minpgflt: minflt, Majpgflt: majflt, ThreadCnt: thrdcnt}, nil
}
